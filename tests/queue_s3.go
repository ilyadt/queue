package queuetest

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cucumber/godog"
)

func (s *Scenario) elementsPushedToQueue(ctx context.Context, n int) error {
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

func (s *Scenario) queueIsEmpty(ctx context.Context) error {
	resp, err := http.Get(s.serverBaseURL + "/" + s.queue)
	if err != nil {
		return fmt.Errorf("not nil error response: %w", err)
	}

	// Empty queue indicated by 404 response
	if resp.StatusCode != 404 {
		return fmt.Errorf("invalid status code in response: %d", resp.StatusCode)
	}

	return nil
}

func (*Scenario) subscribersCancelRequest(ctx context.Context, x int) error {
	cancelC := ctx.Value(CancelChanContextKey).(chan context.CancelFunc)

	for i := 0; i < x; i++ {
		cancel := <-cancelC
		cancel()

		// Ensure request ended
		time.Sleep(10 * time.Millisecond)
	}

	return nil
}

func (s *Scenario) subscribersGotValues(ctx context.Context, y int) error {
	resultC := ctx.Value(ResultChanContextKey).(chan *Result)

	i := 0
	for r := range resultC {
		// successful requests
		if r.Err == nil {
			i++
		}
	}

	if i != y {
		return fmt.Errorf("invalid Number of success request, got=%d, expected=%d", i, y)
	}

	return nil
}

func InitializeScenario3(ctx *godog.ScenarioContext, cfg *ScenarioConfig) {
	s := Scenario{cfg.ServerURL, cfg.QName}

	ctx.Step(`^(\d+) elements pushed to queue$`, s.elementsPushedToQueue)
	ctx.Step(`^Queue is empty$`, s.queueIsEmpty)
	ctx.Step(`^(\d+) subscribers cancel request$`, s.subscribersCancelRequest)
	ctx.Step(`^(\d+) subscribers got values$`, s.subscribersGotValues)
}
