//go:build behavior

package godog

import (
	"fmt"
	"testing"

	"github.com/cucumber/godog"
	"github.com/joho/godotenv"
)

func TestFeatures(t *testing.T) {
	godotenv.Load("testdata/local.env")
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}
	fmt.Println("Running feature tests...")

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
