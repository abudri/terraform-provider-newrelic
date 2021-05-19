package newrelic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/plugins"
)

func dataSourceNewRelicPlugin() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "As of December 2, 2020, plugin access has been limited to accounts that have accessed a legacy plugin in the past 30 days. The legacy plugin experience will reach end of life (EoL) as of June 16, 2021.",
		Read:               dataSourceNewRelicPluginRead,

		Schema: map[string]*schema.Schema{
			"guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The GUID of the plugin in New Relic.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the installed plugin instance.",
			},
		},
	}
}

func dataSourceNewRelicPluginRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Plugins")

	guid := d.Get("guid").(string)

	params := plugins.ListPluginsParams{
		GUID: guid,
	}

	ps, err := client.Plugins.ListPlugins(&params)
	if err != nil {
		return err
	}

	var plugin *plugins.Plugin

	for _, p := range ps {
		if p.GUID == guid {
			plugin = p
			break
		}
	}

	if plugin == nil {
		return fmt.Errorf("the GUID '%s' does not match any New Relic plugins", guid)
	}

	d.SetId(strconv.Itoa(plugin.ID))
	d.Set("id", strconv.Itoa(plugin.ID))

	return nil
}
