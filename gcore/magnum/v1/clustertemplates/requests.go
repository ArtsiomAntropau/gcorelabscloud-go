package clustertemplates

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToClusterTemplateListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the cluster templates attributes you want to see returned. SortKey allows you to sort
// by a particular cluster templates attribute. SortDir sets the direction, and is either
// `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	Limit   int    `q:"limit"`
	Marker  string `q:"marker"`
	SortKey string `q:"sort_key"`
	SortDir string `q:"sort_dir"`
	Detail  bool   `q:"detail"`
}

// ToClusterTemplateListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToClusterTemplateListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// cluster templates. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func List(c *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToClusterTemplateListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ClusterTemplatePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAll returns all cluster templates
func ListAll(c *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]ClusterTemplate, error) {
	page, err := List(c, opts).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractClusterTemplates(page)
}

// Get retrieves a specific cluster template based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, id), &r.Body, nil)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToClusterTemplateCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a cluster template.
type CreateOpts struct {
	ImageID             string             `json:"image_id"`
	KeyPairID           string             `json:"keypair_id"`
	Name                string             `json:"name"`
	DockerVolumeSize    int                `json:"docker_volume_size"`
	Labels              *map[string]string `json:"labels,omitempty"`
	FixedSubnet         *string            `json:"fixed_subnet,omitempty"`
	MasterFlavorID      *string            `json:"master_flavor_id,omitempty"`
	FlavorID            *string            `json:"flavor_id,omitempty"`
	NoProxy             *string            `json:"no_proxy,omitempty"`
	HTTPSProxy          *string            `json:"https_proxy,omitempty"`
	HTTPProxy           *string            `json:"http_proxy,omitempty"`
	TLSDisabled         *bool              `json:"tls_disabled,omitempty"`
	Public              *bool              `json:"public,omitempty"`
	ServerType          *string            `json:"server_type,omitempty"`
	VolumeDriver        *string            `json:"volume_driver,omitempty"`
	RegistryEnabled     *bool              `json:"registry_enabled,omitempty"`
	DockerStorageDriver *string            `json:"docker_storage_driver,omitempty"`
	NetworkDriver       *string            `json:"network_driver,omitempty"`
	FixedNetwork        *string            `json:"fixed_network,omitempty"`
	MasterLbEnabled     *bool              `json:"master_lb_enabled,omitempty"`
	DNSNameServer       *string            `json:"dns_nameserver,omitempty"`
	FloatingIPEnabled   *bool              `json:"floating_ip_enabled,omitempty"`
	Hidden              *bool              `json:"hidden,omitempty"`
}

// ToClusterTemplateCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToClusterTemplateCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and creates a new cluster template using the values
// provided. This operation does not actually require a request body, i.e. the
// CreateOpts struct argument can be empty.
//
// The tenant ID that is contained in the URI is the tenant that creates the
// cluster template. An admin user, however, has the option of specifying another tenant
// ID in the CreateOpts struct.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToClusterTemplateCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToClusterTemplateUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update a network.
type UpdateOpts struct {
}

// ToClusterTemplateUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToClusterTemplateUpdateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Update accepts a UpdateOpts struct and updates an existing network using the
// values provided. For more information, see the Create function.
func Update(c *gcorecloud.ServiceClient, clusterTemplateID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToClusterTemplateUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(updateURL(c, clusterTemplateID), b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// Delete accepts a unique ID and deletes the cluster template associated with it.
func Delete(c *gcorecloud.ServiceClient, clusterTemplateID string) (r DeleteResult) {
	_, r.Err = c.Delete(deleteURL(c, clusterTemplateID), nil)
	return
}

// IDFromName is a convenience function that returns a cluster template ID, given
// its name.
func IDFromName(client *gcorecloud.ServiceClient, name string) (string, error) {
	count := 0
	id := ""

	listOpts := ListOpts{}

	pages, err := List(client, listOpts).AllPages()
	if err != nil {
		return "", err
	}

	all, err := ExtractClusterTemplates(pages)
	if err != nil {
		return "", err
	}

	for _, s := range all {
		if s.Name == name {
			count++
			id = s.UUID
		}
	}

	switch count {
	case 0:
		return "", gcorecloud.ErrResourceNotFound{Name: name, ResourceType: "clustertemplates"}
	case 1:
		return id, nil
	default:
		return "", gcorecloud.ErrMultipleResourcesFound{Name: name, Count: count, ResourceType: "clustertemplates"}
	}
}
