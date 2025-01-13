package queuetest_test

import (
	"math/rand"
	"queuetest"
	"strconv"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			cfg := &queuetest.ScenarioConfig{
				ServerURL: "http://127.0.0.1:2802",
				QName:     "numbers_" + strconv.Itoa(rand.Int()),
			}

			queuetest.InitializeScenario1(ctx, cfg)
			queuetest.InitializeScenario2(ctx, cfg)
			queuetest.InitializeScenario3(ctx, cfg)
			queuetest.InitializeValidationScenario(ctx, cfg)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
