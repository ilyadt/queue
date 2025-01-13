package queuetest

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

func (s *Scenario) iPutElementsInQueue(ctx context.Context, n int) error {
	for i := 1; i <= n; i++ {
		req, _ := http.NewRequest("PUT", s.serverBaseURL+"/"+s.queue+`?v=`+strconv.Itoa(i), nil)
		resp, err := http.DefaultClient.Do(req)
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

func (s *Scenario) subscribersGetValuesInTheFifoOrder(ctx context.Context) error {
	resultC := ctx.Value(ResultChanContextKey).(chan *Result)

	for res := range resultC {
		if res.Err != nil {
			return fmt.Errorf("error subscriber resp: %w: %+v", res.Err, res)
		}

		if res.Resp != strconv.Itoa(res.Num) {
			return fmt.Errorf("invalid resp: %+v", res)
		}
	}

	return nil
}

func (s *Scenario) subscribersWaitingForValueInQueue(ctx context.Context, n int) (context.Context, error) {
	ResultC := make(chan *Result, n)
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
						defer func() {
							connectedC <- struct{}{}
						}()

						nDialer := &net.Dialer{}

						conn, err := nDialer.DialContext(ctx, network, address)
						if err != nil {
							return nil, err
						}

						// Wait some time for the request pass into the controller after connection
						time.Sleep(10 * time.Millisecond)

						return conn, nil
					},
				},
			}

			ctx2, cancel := context.WithCancel(context.Background())
			cancelC <- cancel

			req, _ := http.NewRequestWithContext(ctx2, "GET", s.serverBaseURL+"/"+s.queue+"?timeout=300", nil)
			resp, err := client.Do(req)
			if err != nil {
				ResultC <- &Result{Num: i, Err: err}
				return
			}

			if resp.StatusCode != 200 {
				ResultC <- &Result{Num: i, Err: fmt.Errorf("status code %d", resp.StatusCode)}
				return
			}

			body, err := io.ReadAll(resp.Body)
			_ = resp.Body.Close()

			if err != nil {
				ResultC <- &Result{Num: i, Err: err}
				return
			}

			ResultC <- &Result{Num: i, Resp: string(body)}
		}(i)

		// After client is connected to the server, go further to the next request
		<-connectedC
	}

	go func() {
		wg.Wait()
		close(ResultC)
	}()

	ctx = context.WithValue(ctx, ResultChanContextKey, ResultC)
	ctx = context.WithValue(ctx, CancelChanContextKey, cancelC)

	return ctx, nil
}

func InitializeScenario2(ctx *godog.ScenarioContext, cfg *ScenarioConfig) {
  s := Scenario{cfg.ServerURL, cfg.QName}

	ctx.Step(`^I put (\d+) elements in queue$`, s.iPutElementsInQueue)
	ctx.Step(`^subscribers get values in the fifo order$`, s.subscribersGetValuesInTheFifoOrder)
	ctx.Step(`^(\d+) subscribers waiting for value in queue$`, s.subscribersWaitingForValueInQueue)
}
