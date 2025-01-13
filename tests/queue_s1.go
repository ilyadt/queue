package queuetest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/cucumber/godog"
)

func (s *Scenario) iGetElementsFromQueue(ctx context.Context, n int) error {
	for i := 0; i < n; i++ {
		resp, err := http.DefaultClient.Get(s.serverBaseURL + "/" + s.queue)

		if err != nil {
			return fmt.Errorf("not nil error response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("invalid status code in response: %d", resp.StatusCode)
		}
	}

	return nil
}

func (s *Scenario) nextElementWillBe(ctx context.Context, value string) error {
	resp, err := http.DefaultClient.Get(s.serverBaseURL + "/" + s.queue)

	if err != nil {
		return fmt.Errorf("not nil error response: %w", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("invalid status code in response: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body: %w", err)
	}

	if string(body) != value {
		return fmt.Errorf("invalid value came from queue(%s): `%s`, expected: `%s`", s.queue, string(body), value)
	}

	return nil
}

func (s *Scenario) thereAreNElementsInQueueInOrderFromOneToN(ctx context.Context, n int) error {
	for i := 1; i <= n; i++ {
		req, _ := http.NewRequest("PUT", s.serverBaseURL + "/" + s.queue+`?v=`+strconv.Itoa(i), nil)
		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			return fmt.Errorf("not nil error response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("invalid status code in response: %d", resp.StatusCode)
		}
	}

	return nil
}


func InitializeScenario1(ctx *godog.ScenarioContext, cfg *ScenarioConfig) {
  s1 := Scenario{cfg.ServerURL, cfg.QName}

  ctx.Step(`^I get (\d+) elements from queue$`, s1.iGetElementsFromQueue)
	ctx.Step(`^next element will be (\d+)$`, s1.nextElementWillBe)
	ctx.Step(`^there are (\d+) elements in queue in order from 1 to N$`, s1.thereAreNElementsInQueueInOrderFromOneToN)
}
