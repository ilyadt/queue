package queuetest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/cucumber/godog"
)

func InitializeScenario1(ctx *godog.ScenarioContext) {
	ctx.Step(`^I get (\d+) elements from queue$`, iGetElementsFromQueue)
	ctx.Step(`^next element will be (\d+)$`, nextElementWillBe)
	ctx.Step(`^there are (\d+) elements in queue in order from 1 to N$`, thereAreNElementsInQueueInOrderFromOneToN)
}

type MyCtx struct {
	context.Context
	QName         string
	ServerBaseURL string
	ResultC       chan *Result
	CancelChan    chan context.CancelFunc // cancel function for each request
}

type Result struct {
	Num  int    //
	Err  error  //
	Resp string //
}

func iGetElementsFromQueue(ctx *MyCtx, n int) error {
	for i := 0; i < n; i++ {
		resp, err := http.DefaultClient.Get(ctx.ServerBaseURL + "/" + ctx.QName)

		if err != nil {
			return fmt.Errorf("not nil error response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("invalid status code in response: %d", resp.StatusCode)
		}
	}

	return nil
}

func nextElementWillBe(ctx *MyCtx, value string) error {
	resp, err := http.DefaultClient.Get(ctx.ServerBaseURL + "/" + ctx.QName)

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
		return fmt.Errorf("invalid value came from queue(%s): `%s`, expected: `%s`", ctx.QName, string(body), value)
	}

	return nil
}

func thereAreNElementsInQueueInOrderFromOneToN(ctx *MyCtx, n int) error {
	for i := 1; i <= n; i++ {
		req, _ := http.NewRequest("PUT", ctx.ServerBaseURL+"/"+ctx.QName+`?v=`+strconv.Itoa(i), nil)
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
