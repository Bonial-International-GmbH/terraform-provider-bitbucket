package bitbucket

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataReadGroup,

		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"auto_add": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"permission": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataReadGroup(d *schema.ResourceData, m interface{}) error {
	client := m.(Clients).httpClient

	workspace := d.Get("workspace").(string)
	slug := d.Get("slug").(string)

	groupsReq, _ := client.Get(fmt.Sprintf("1.0/groups/%s/%s", workspace, slug))

	if groupsReq.StatusCode == http.StatusNotFound {
		log.Printf("[WARN] Group (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if groupsReq.StatusCode == http.StatusNotFound {
		log.Printf("[WARN] Group (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if groupsReq.Body == nil {
		return fmt.Errorf("error reading Group (%s): empty response", d.Id())
	}

	var grp *UserGroup

	body, readerr := ioutil.ReadAll(groupsReq.Body)
	if readerr != nil {
		return readerr
	}

	log.Printf("[DEBUG] Groups Response JSON: %v", string(body))

	decodeerr := json.Unmarshal(body, &grp)
	if decodeerr != nil {
		return decodeerr
	}

	log.Printf("[DEBUG] Groups Response Decoded: %#v", grp)

	d.SetId(fmt.Sprintf("%s/%s", workspace, slug))
	d.Set("workspace", workspace)
	d.Set("slug", grp.Slug)
	d.Set("name", grp.Name)
	d.Set("auto_add", grp.AutoAdd)
	d.Set("permission", grp.Permission)

	return nil
}
