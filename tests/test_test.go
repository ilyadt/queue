package main_test

import (
	"context"
	"github.com/cucumber/godog"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

type ResGetTimeout struct {
	num  int    //
	err  error  //
	resp string //
}

type MyCtx struct {
	context.Context
	qName            string
	serverBaseURL    string
	resGetTimeout    chan *ResGetTimeout
	numberOfRequests int
	cancelChan       chan context.CancelFunc // cancel function for each request
}

func init() {
	rand.Seed(time.Now().Unix())
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		return &MyCtx{
			Context:       ctx,
			qName:         "numbers_" + strconv.Itoa(rand.Int()),
			serverBaseURL: "http://127.0.0.1:2802",
		}, nil
	})

	InitializeScenario1(ctx)
	InitializeScenario2(ctx)
	InitializeScenario3(ctx)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features/queue.feature"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
