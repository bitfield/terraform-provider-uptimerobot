// +build integration

package main

import (
	"os"
	"testing"

	uptimerobot "github.com/bitfield/uptimerobot/pkg"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

var client uptimerobot.Client

func getAPIKey(t *testing.T) string {
	key := os.Getenv("UPTIMEROBOT_API_KEY")
	if key == "" {
		t.Fatal("'UPTIMEROBOT_API_KEY' must be set for integration tests")
	}
	return key
}

func TestUptimeRobotTerraformIntegration(t *testing.T) {
	t.Parallel()
	terraformOptions := &terraform.Options{
		TerraformDir: ".",
		Vars: map[string]interface{}{
			"uptimerobot_api_key": getAPIKey(t),
		},
	}
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)
	planPath := "./test.plan"
	exit, err := terraform.GetExitCodeForTerraformCommandE(t, terraformOptions, terraform.FormatArgs(terraformOptions, "plan", "--out="+planPath, "-input=false", "-lock=true", "-detailed-exitcode")...)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(planPath)
	if exit != terraform.DefaultSuccessExitCode {
		t.Fatalf("want DefaultSuccessExitCode (indicating plan is a no-op), got %d", exit)
	}
}
