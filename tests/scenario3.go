package queuetest

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cucumber/godog"
)

type Scenario3 struct {
	serverBaseURL string
	queue string
}

func (s3 *Scenario3) elementsPushedToQueue(ctx context.Context, n int) error {
	for i := 1; i <= n; i++ {
		req, _ := http.NewRequest("PUT", s3.serverBaseURL+"/"+s3.queue+`?v=`+strconv.Itoa(i), nil)
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

func (s3 *Scenario3) queueIsEmpty(ctx context.Context) error {
	resp, err := http.Get(s3.serverBaseURL + "/" + s3.queue)
	if err != nil {
		return fmt.Errorf("not nil error response: %w", err)
	}

	// Empty queue indicated by 404 response
	if resp.StatusCode != 404 {
		return fmt.Errorf("invalid status code in response: %d", resp.StatusCode)
	}

	return nil
}

func (*Scenario3) subscribersCancelRequest(ctx context.Context, x int) error {
	cancelC := ctx.Value(CancelChanContextKey).(chan context.CancelFunc)

	for i := 0; i < x; i++ {
		cancel := <-cancelC
		cancel()

		// Ensure request ended
		time.Sleep(10 * time.Millisecond)
	}

	return nil
}

func (s3 *Scenario3) subscribersGotValues(ctx context.Context, y int) error {
	resultC := ctx.Value(ResultChanContextKey).(chan *Result)

	i := 0
	for r := range resultC {
		if r.Err == nil {
			i++ // successful requests
			fmt.Printf("RequestNO:%d resp:%s\n", r.Num, r.Resp)
		} else {
			fmt.Printf("RequestNO:%d Error: %v\n", r.Num, r.Err.Error())
		}

	}

	if i != y {
		return fmt.Errorf("invalid Number of success request, got=%d, expected=%d", i, y)
	}

	return nil
}

func InitializeScenario3(ctx *godog.ScenarioContext, cfg *ScenarioConfig) {
	s3 := Scenario3{cfg.ServerURL, cfg.QName}

	ctx.Step(`^(\d+) elements pushed to queue$`, s3.elementsPushedToQueue)
	ctx.Step(`^Queue is empty$`, s3.queueIsEmpty)
	ctx.Step(`^(\d+) subscribers cancel request$`, s3.subscribersCancelRequest)
	ctx.Step(`^(\d+) subscribers got values$`, s3.subscribersGotValues)
}
