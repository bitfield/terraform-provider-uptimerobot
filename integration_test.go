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

func TestHandleResourceGoneAwayIntegration(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: ".",
		Vars: map[string]interface{}{
			"uptimerobot_api_key": getAPIKey(t),
		},
	}
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)
	utr := uptimerobot.New(getAPIKey(t))
	results, err := utr.SearchMonitors("My test monitor")
	if err != nil {
		t.Fatalf("failed to find test monitor: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected one match for test monitor but got %d", len(results))
	}
	ID := results[0].ID
	// Delete the resource out-of-band to make sure that Terraform correctly
	// handles this situation when planning
	err = utr.DeleteMonitor(ID)
	if err != nil {
		t.Fatalf("failed to delete test monitor out of band: %v", err)
	}
	planPath := "./test.plan"
	exit, err := terraform.GetExitCodeForTerraformCommandE(t, terraformOptions, terraform.FormatArgs(terraformOptions, "plan", "--out="+planPath, "-input=false", "-lock=true", "-detailed-exitcode")...)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(planPath)
	if exit != terraform.TerraformPlanChangesPresentExitCode {
		t.Fatalf("want TerraformPlanChangesPresentExitCode (indicating plan would change resources), got %d", exit)
	}
}
