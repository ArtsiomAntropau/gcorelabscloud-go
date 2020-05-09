package tokens

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

func processToken(c *gcorecloud.ServiceClient, opts gcorecloud.AuthOptionsBuilder, url string) (r TokenResult) {
	b := opts.ToMap()
	resp, err := c.Post(url, b, &r.Body, &gcorecloud.RequestOpts{})
	r.Err = err
	if resp != nil {
		r.Header = resp.Header
	}
	return
}

// Create authenticates and either generates a new token
func Create(c *gcorecloud.ServiceClient, opts gcorecloud.AuthOptionsBuilder) (r TokenResult) {
	return processToken(c, opts, tokenURL(c))
}

// RefreshPlatform token with GCore platform API
func RefreshPlatform(c *gcorecloud.ServiceClient, opts gcorecloud.TokenOptionsBuilder) (r TokenResult) {
	return processToken(c, opts, refreshURL(c))
}

// RefreshPlatform token with gcloud API
func RefreshGCloud(c *gcorecloud.ServiceClient, opts gcorecloud.TokenOptionsBuilder) (r TokenResult) {
	return processToken(c, opts, refreshGCloudURL(c))
}
