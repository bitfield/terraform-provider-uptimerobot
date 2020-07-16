package main

import (
	"fmt"
	"os"

	uptimerobot "github.com/bitfield/uptimerobot/pkg"
	"github.com/hashicorp/terraform/helper/schema"
)

// Provider makes the provider available to Terraform.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("UPTIMEROBOT_API_KEY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"uptimerobot_monitor": resourceMonitor(),
		},
		ConfigureFunc: func(r *schema.ResourceData) (interface{}, error) {
			client := uptimerobot.New(r.Get("api_key").(string))
			debugLog := os.Getenv("UPTIMEROBOT_DEBUG_LOG")
			if debugLog != "" {
				debugFile, err := os.OpenFile(debugLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
				if err != nil {
					panic(fmt.Sprintf("can't write to debug log file: %v", err))
				}
				client.Debug = debugFile
			}
			return &client, nil
		},
	}
}
