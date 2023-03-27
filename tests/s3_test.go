package main_test

import (
	"fmt"
	"github.com/cucumber/godog"
	"net/http"
	"strconv"
	"time"
)

func elementsPushedToQueue(ctx *MyCtx, n int) error {
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

func queueIsEmpty(ctx *MyCtx) error {
	resp, err := http.Get(ctx.serverBaseURL + "/" + ctx.qName)
	if err != nil {
		return fmt.Errorf("not nil error response: %w", err)
	}

	// Empty queue indicated by 404 response
	if resp.StatusCode != 404 {
		return fmt.Errorf("invalid status code in response: %d", resp.StatusCode)
	}

	return nil
}

func subscribersCancelRequest(ctx *MyCtx, x int) error {
	for i := 0; i < x; i++ {
		cancel := <-ctx.cancelChan
		cancel()

		// Ensure request ended
		time.Sleep(10 * time.Millisecond)
	}

	return nil
}

func subscribersGotValues(ctx *MyCtx, y int) error {
	i := 0
	for r := range ctx.resGetTimeout {
		if r.err == nil {
			i++ // successful requests
			fmt.Printf("RequestNO:%d resp:%s\n", r.num, r.resp)
		} else {
			fmt.Printf("RequestNO:%d Error: %v\n", r.num, r.err.Error())
		}

	}

	if i != y {
		return fmt.Errorf("invalid number of success request, got=%d, expected=%d", i, y)
	}

	return nil
}

func InitializeScenario3(ctx *godog.ScenarioContext) {
	ctx.Step(`^(\d+) elements pushed to queue$`, elementsPushedToQueue)
	ctx.Step(`^Queue is empty$`, queueIsEmpty)
	ctx.Step(`^(\d+) subscribers cancel request$`, subscribersCancelRequest)
	ctx.Step(`^(\d+) subscribers got values$`, subscribersGotValues)
}
