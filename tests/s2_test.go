package main_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/cucumber/godog"
)

func iPutElementsInQueue(ctx *MyCtx, n int) error {
	client := http.DefaultClient

	for i := 1; i <= n; i++ {
		req, _ := http.NewRequest("PUT", ctx.serverBaseURL+"/"+ctx.qName+`?v=`+strconv.Itoa(i), nil)
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("not nil error response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("invalid status code in response: %d", resp.StatusCode)
		}

		_ = resp.Body.Close()
	}

	return nil
}

func subscribersGetValuesInTheFifoOrder(ctx *MyCtx) error {
	resultC := ctx.resultC

	for res := range resultC {
		if res.err != nil {
			return fmt.Errorf("error subscriber resp: %w: %+v", res.err, res)
		}

		if res.resp != strconv.Itoa(res.num) {
			return fmt.Errorf("invalid resp: %+v", res)
		}
	}

	return nil
}

func subscribersWaitingForValueInQueue(ctx *MyCtx, n int) (context.Context, error) {
	resultC := make(chan *Result, n)
	cancelC := make(chan context.CancelFunc, n)

	connectedC := make(chan struct{})

	var wg sync.WaitGroup

	for i := 1; i <= n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			client := &http.Client{
				Transport: &http.Transport{
					DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
						nDialer := &net.Dialer{}

						conn, err := nDialer.DialContext(ctx, network, address)
						if err != nil {
							return conn, err
						}

						// Wait some time for the request pass into the controller after connection
						time.Sleep(10 * time.Millisecond)
						connectedC <- struct{}{}

						return conn, err
					},
				},
			}

			ctx2, cancel := context.WithCancel(context.Background())
			cancelC <- cancel

			req, _ := http.NewRequestWithContext(ctx2, "GET", ctx.serverBaseURL+"/"+ctx.qName+"?timeout=300", nil)
			resp, err := client.Do(req)
			if err != nil {
				resultC <- &Result{num: i, err: err}
				return
			}

			if resp.StatusCode != 200 {
				resultC <- &Result{num: i, err: fmt.Errorf("status code %d", resp.StatusCode)}
				return
			}

			body, err := io.ReadAll(resp.Body)
			_ = resp.Body.Close()

			if err != nil {
				resultC <- &Result{num: i, err: err}
				return
			}

			resultC <- &Result{num: i, resp: string(body)}
		}(i)

		// After client is connected to the server, go further to the next request
		<-connectedC
	}

	go func() {
		wg.Wait()
		close(resultC)
	}()

	ctx.resultC = resultC
	ctx.cancelChan = cancelC

	return ctx, nil
}

func InitializeScenario2(ctx *godog.ScenarioContext) {
	ctx.Step(`^I put (\d+) elements in queue$`, iPutElementsInQueue)
	ctx.Step(`^subscribers get values in the fifo order$`, subscribersGetValuesInTheFifoOrder)
	ctx.Step(`^(\d+) subscribers waiting for value in queue$`, subscribersWaitingForValueInQueue)
}
