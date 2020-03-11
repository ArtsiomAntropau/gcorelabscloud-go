package clusters

import (
	"gcloud/gcorecloud-go"
)

func resourceURL(c *gcorecloud.ServiceClient, id string) string {
	return c.ServiceURL("clusters", id)
}

func rootURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL("clusters")
}

func resourceActionURL(c *gcorecloud.ServiceClient, id, action string) string {
	return c.ServiceURL("clusters", id, "actions", action)
}

func resizeURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "resize")
}

func getURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func listURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func createURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func deleteURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}
