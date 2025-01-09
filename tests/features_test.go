package queuetest_test

import (
	"context"
	"math/rand"
	"strconv"
	"testing"
	"time"
	"queuetest"

	"github.com/cucumber/godog"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
				return &queuetest.MyCtx{
					Context:       ctx,
					QName:         "numbers_" + strconv.Itoa(rand.Int()),
					ServerBaseURL: "http://127.0.0.1:2802",
				}, nil
			})
			queuetest.InitializeScenario1(ctx)
			InitializeScenario2(ctx)
			InitializeScenario3(ctx)
		},
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
