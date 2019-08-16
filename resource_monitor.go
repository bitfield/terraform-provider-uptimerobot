package main

import (
	"fmt"
	"strconv"

	uptimerobot "github.com/bitfield/uptimerobot/pkg"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceMonitorCreate,
		Read:   resourceMonitorRead,
		Delete: resourceMonitorDelete,

		Schema: map[string]*schema.Schema{
			"friendly_name": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"alert_contact": &schema.Schema{
				Type:     schema.TypeSet,
				ForceNew: true,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceMonitorCreate(d *schema.ResourceData, client interface{}) error {
	mon := uptimerobot.Monitor{
		URL:          d.Get("url").(string),
		FriendlyName: d.Get("friendly_name").(string),
		Type:         typeFromString(d.Get("type").(string)),
		Port:         80,
	}
	contacts := d.Get("alert_contact").(*schema.Set).List()
	for _, c := range contacts {
		mon.AlertContacts = append(mon.AlertContacts, c.(string))
	}
	ID, err := client.(*uptimerobot.Client).CreateMonitor(mon)
	if err != nil {
		return fmt.Errorf("API error: %v", err)
	}
	d.SetId(fmt.Sprintf("%d", ID))
	return resourceMonitorRead(d, client)
}

func resourceMonitorRead(d *schema.ResourceData, client interface{}) error {
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("bad ID %s: %v", d.Id(), err)
	}
	mon, err := client.(*uptimerobot.Client).GetMonitor(id)
	if err != nil {
		return fmt.Errorf("API error: %v", err)
	}
	d.Set("url", mon.URL)
	d.Set("type", mon.FriendlyType)
	d.Set("alert_contacts", mon.AlertContacts)
	d.SetId(fmt.Sprintf("%d", mon.ID))
	return nil
}

func resourceMonitorDelete(d *schema.ResourceData, client interface{}) error {
	ID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("bad ID %s: %v", d.Id(), err)
	}
	if err = client.(*uptimerobot.Client).DeleteMonitor(ID); err != nil {
		return fmt.Errorf("API error: %v", err)
	}
	return nil
}

// typeFromString returns an integer monitor type based on the supplied string.
func typeFromString(humanReadableType string) int {
	switch humanReadableType {
	case "HTTP":
		return uptimerobot.TypeHTTP
	case "Keyword":
		return uptimerobot.TypeKeyword
	case "Ping":
		return uptimerobot.TypePing
	case "Port":
		return uptimerobot.TypePort
	default:
		return 0
	}
}
