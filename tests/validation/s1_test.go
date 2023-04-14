package main_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/cucumber/godog"
)

type ctxResponseKeyType string
const ctxResponseKey = ctxResponseKeyType("response")

func iGetStatus(ctx context.Context, status int) error {
	resp := ctx.Value(ctxResponseKey).(*http.Response)

	if resp.StatusCode != status {
		return fmt.Errorf("invalid http status code: %d, expected %d", resp.StatusCode, status)
	}

	return nil
}

func iRequestValueWithNegativeTimeout(ctx context.Context) (context.Context, error) {
	baseUrl := ctx.Value("serverBaseURL").(string)

	resp, err := http.Get(baseUrl + "/queue?timeout=-7")
	if err != nil {
		return ctx, fmt.Errorf("request negative timeout err: %w", err)
	}

	return context.WithValue(ctx, ctxResponseKey, resp), nil
}

func iPutValueInQueue(ctx context.Context, val, queue string) (context.Context, error) {
	baseUrl := ctx.Value("serverBaseURL").(string)

	req, err := http.NewRequest("PUT", baseUrl+"/"+queue+"?v="+val, nil)
	if err != nil {
		panic("build req err: " + err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ctx, fmt.Errorf("put request failed to perform: %w", err)
	}

	return context.WithValue(ctx, ctxResponseKey, resp), nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		return context.WithValue(ctx, "serverBaseURL", "http://127.0.0.1:2802"), nil
	})

	ctx.Step(`^I get (\d+) status$`, iGetStatus)
	ctx.Step(`^I request value with negative timeout$`, iRequestValueWithNegativeTimeout)
	ctx.Step(`^I put value "([^"]*)" in queue "([^"]*)"$`, iPutValueInQueue)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features/validation.feature"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
