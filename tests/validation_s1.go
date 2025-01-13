package queuetest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cucumber/godog"
)

func (s *Scenario) iGetStatus(ctx context.Context, status int) error {
	resp := ctx.Value(ResponseContextKey).(*http.Response)

	if resp.StatusCode != status {
		return fmt.Errorf("invalid http status code: %d, expected %d", resp.StatusCode, status)
	}

	return nil
}

func (s *Scenario) iRequestValueWithNegativeTimeout(ctx context.Context) (context.Context, error) {
	resp, err := http.Get(s.serverBaseURL + "/queue?timeout=-7")
	if err != nil {
		return ctx, fmt.Errorf("request negative timeout err: %w", err)
	}

	return context.WithValue(ctx, ResponseContextKey, resp), nil
}

func (s *Scenario) iPutValueInQueue(ctx context.Context, val, queue string) (context.Context, error) {
	req, err := http.NewRequest("PUT", s.serverBaseURL+"/"+queue+"?v="+val, nil)
	if err != nil {
		panic("build req err: " + err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ctx, fmt.Errorf("put request failed to perform: %w", err)
	}

	return context.WithValue(ctx, ResponseContextKey, resp), nil
}

func InitializeValidationScenario(ctx *godog.ScenarioContext, cfg *ScenarioConfig) {
	s := Scenario{serverBaseURL: cfg.ServerURL}

	ctx.Step(`^I get (\d+) status$`, s.iGetStatus)
	ctx.Step(`^I request value with negative timeout$`, s.iRequestValueWithNegativeTimeout)
	ctx.Step(`^I put value "([^"]*)" in queue "([^"]*)"$`, s.iPutValueInQueue)
}
