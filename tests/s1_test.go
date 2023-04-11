package main_test

import (
	"fmt"
	"github.com/cucumber/godog"
	"io"
	"net/http"
	"strconv"
)

func iGetElementsFromQueue(ctx *MyCtx, x int) error {
	client := http.DefaultClient
	for i := 0; i < x; i++ {
		resp, err := client.Get(ctx.serverBaseURL + "/" + ctx.qName)

		if err != nil {
			return fmt.Errorf("not nil error response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("invalid status code in response: %d", resp.StatusCode)
		}
	}

	return nil
}

func nextElementWillBe(ctx *MyCtx, value int) error {
	resp, err := http.DefaultClient.Get(ctx.serverBaseURL + "/" + ctx.qName)

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

	if strconv.Itoa(value) != string(body) {
		return fmt.Errorf("invalid value came from queue(%s): `%s`, expected: `%s`", ctx.qName, string(body), strconv.Itoa(value))
	}

	return nil
}

func thereAreNElementsInQueueInOrderFromOneToN(ctx *MyCtx, n int) error {
	client := http.DefaultClient

	for i := 1; i <= n; i++ {
		req, _ := http.NewRequest("PUT", ctx.serverBaseURL + "/" +ctx.qName+`?v=`+strconv.Itoa(i), nil)
		resp, err := client.Do(req)

		if err != nil {
			return fmt.Errorf("not nil error response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("invalid status code in response: %d", resp.StatusCode)
		}
	}

	return nil
}

func InitializeScenario1(ctx *godog.ScenarioContext) {
	ctx.Step(`^I get (\d+) elements from queue$`, iGetElementsFromQueue)
	ctx.Step(`^next element will be (\d+)$`, nextElementWillBe)
	ctx.Step(`^there are (\d+) elements in queue in order from one to N$`, thereAreNElementsInQueueInOrderFromOneToN)
}
