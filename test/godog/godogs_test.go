package godog

import (
	"github.com/cucumber/godog"
	"github.com/kelseyhightower/envconfig"
)

func InitializeScenario(ctx *godog.ScenarioContext) {
	var cfg Config
	envconfig.MustProcess("", &cfg)
	steps := NewSteps(cfg)

	// pulse
	ctx.Step(`^I run the pulse github command on user "([^"]*)"$`, steps.iRunThePulseGithubCommandOnUser)
	ctx.Step(`^the pulse output should equal "([^"]*)"$`, steps.thePulseOutputShouldEqual)

	// rest assured
	ctx.Step(`^I have a clean rest assured environment$`, steps.iHaveACleanRestAssuredEnvironment)
	ctx.Step(`^rest assured returns "([^"]*)" on a "([^"]*)" to "([^"]*)"$`, steps.restAssuredReturns)
}
