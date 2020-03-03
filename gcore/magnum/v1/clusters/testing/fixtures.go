package testing

import (
	"gcloud/gcorecloud-go"
	"gcloud/gcorecloud-go/gcore/magnum/v1/clusters"
	"gcloud/gcorecloud-go/gcore/task/v1/tasks"
	"time"
)

const ListResponse = `
{
  "clusters": [
    {
      "health_status": null,
      "master_flavor_id": "g1-standard-1-2",
      "flavor_id": "g1-standard-1-2",
      "name": "fokgkcytgg",
      "labels": {
        "gcloud_project_id": "12",
        "gcloud_region_id": "1",
        "gcloud_access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiZXhwIjoxNTg4MzM1NjM3LCJqdGkiOiJhYTE2ODhmODdmNDc1YjhmNDk3NTY5MmI5MTkyZDdmYiIsInVzZXJfdHlwZSI6ImNvbW1vbiIsInVzZXJfZ3JvdXBzIjpbIlVzZXJzIl0sImNsaWVudF9pZCI6LTIwMDMwMjEyMjAzNzUzNTI3NiwicmVnaW9uX2lkIjoxLCJwcm9qZWN0X2lkIjoxMiwidXNlcl9pZCI6MSwiaXNfYWRtaW4iOmZhbHNlfQ.jdPEAMuZOEqT9Ns1eW0IOZmo33WZsMEIs8NFXuF29iU",
        "gcloud_refresh_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImV4cCI6MTg5ODUxMTYzNywianRpIjoiYWExNjg4Zjg3ZjQ3NWI4ZjQ5NzU2OTJiOTE5MmQ3ZmIiLCJ1c2VyX3R5cGUiOiJjb21tb24iLCJ1c2VyX2dyb3VwcyI6WyJVc2VycyJdLCJjbGllbnRfaWQiOi0yMDAzMDIxMjIwMzc1MzUyNzYsInJlZ2lvbl9pZCI6MSwicHJvamVjdF9pZCI6MTIsInVzZXJfaWQiOjEsImlzX2FkbWluIjpmYWxzZX0.AS2Xv067CIxbJdjMB7Z4ydCdxEKwlRx_rLoKNheL0ks",
        "boot_volume_size": "10",
        "auto_scaling_enabled": "true"
      },
      "keypair": "keypair",
      "links": [
        {
          "href": "http://10.100.178.165:9511/clusters/e4028530-0353-494b-a84c-0230122e34ff",
          "rel": "self"
        },
        {
          "href": "http://10.100.178.165:9511/clusters/e4028530-0353-494b-a84c-0230122e34ff",
          "rel": "bookmark"
        }
      ],
      "status": "CREATE_IN_PROGRESS",
      "cluster_template_id": "0a5ce9dd-a484-4e23-80c7-7e586c80d9fc",
      "create_timeout": 360,
      "uuid": "e4028530-0353-494b-a84c-0230122e34ff",
      "master_count": 1,
      "docker_volume_size": 5,
      "node_count": 1,
      "stack_id": "78c48153-fa6c-48b8-aae3-08b5b230387a",
      "floating_ip_enabled": false
    }
  ]
}
`

const GetResponse = `
{
  "health_status": null,
  "master_flavor_id": "g1-standard-1-2",
  "discovery_url": "https://discovery.etcd.io/161d73fde241377395f481c6276b42c7",
  "flavor_id": "g1-standard-1-2",
  "name": "fokgkcytgg",
  "labels": {
    "gcloud_project_id": "12",
    "gcloud_region_id": "1",
    "gcloud_access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiZXhwIjoxNTg4MzM1NjM3LCJqdGkiOiJhYTE2ODhmODdmNDc1YjhmNDk3NTY5MmI5MTkyZDdmYiIsInVzZXJfdHlwZSI6ImNvbW1vbiIsInVzZXJfZ3JvdXBzIjpbIlVzZXJzIl0sImNsaWVudF9pZCI6LTIwMDMwMjEyMjAzNzUzNTI3NiwicmVnaW9uX2lkIjoxLCJwcm9qZWN0X2lkIjoxMiwidXNlcl9pZCI6MSwiaXNfYWRtaW4iOmZhbHNlfQ.jdPEAMuZOEqT9Ns1eW0IOZmo33WZsMEIs8NFXuF29iU",
    "gcloud_refresh_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImV4cCI6MTg5ODUxMTYzNywianRpIjoiYWExNjg4Zjg3ZjQ3NWI4ZjQ5NzU2OTJiOTE5MmQ3ZmIiLCJ1c2VyX3R5cGUiOiJjb21tb24iLCJ1c2VyX2dyb3VwcyI6WyJVc2VycyJdLCJjbGllbnRfaWQiOi0yMDAzMDIxMjIwMzc1MzUyNzYsInJlZ2lvbl9pZCI6MSwicHJvamVjdF9pZCI6MTIsInVzZXJfaWQiOjEsImlzX2FkbWluIjpmYWxzZX0.AS2Xv067CIxbJdjMB7Z4ydCdxEKwlRx_rLoKNheL0ks",
    "boot_volume_size": "10",
    "auto_scaling_enabled": "true"
  },
  "keypair": "keypair",
  "links": [
    {
      "href": "http://10.100.178.165:9511/clusters/e4028530-0353-494b-a84c-0230122e34ff",
      "rel": "self"
    },
    {
      "href": "http://10.100.178.165:9511/clusters/e4028530-0353-494b-a84c-0230122e34ff",
      "rel": "bookmark"
    }
  ],
  "fixed_subnet": null,
  "coe_version": null,
  "master_addresses": [],
  "status": "CREATE_IN_PROGRESS",
  "cluster_template_id": "0a5ce9dd-a484-4e23-80c7-7e586c80d9fc",
  "project_id": "ec0f251d-2e36-436c-9a30-7e2c33297273",
  "created_at": "2020-03-02T12:20:43+00:00",
  "container_version": null,
  "status_reason": null,
  "create_timeout": 360,
  "health_status_reason": {},
  "api_address": null,
  "uuid": "e4028530-0353-494b-a84c-0230122e34ff",
  "master_count": 1,
  "user_id": "8ba64372-1585-4808-b422-7a7aab5f3197",
  "node_addresses": [],
  "updated_at": "2020-03-02T12:20:47+00:00",
  "docker_volume_size": 5,
  "node_count": 1,
  "stack_id": "78c48153-fa6c-48b8-aae3-08b5b230387a",
  "fixed_network": null,
  "floating_ip_enabled": false
}
`

const CreateRequest = `
{
    "name": "fokgkcytgg",
    "master_count": 1,
    "cluster_template_id": "0a5ce9dd-a484-4e23-80c7-7e586c80d9fc",
    "node_count": 1,
    "create_timeout": 360, 
    "keypair": "keypair",
	"master_flavor_id": "g1-standard-1-2",
    "labels": {},
	"flavor_id": "g1-standard-1-2",
    "floating_ip_enabled": false
}
`

const CreateResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

var createdTimeString = "2020-03-02T12:20:43+00:00"
var updatedTimeString = "2020-03-02T12:20:47+00:00"
var createdTime, _ = time.Parse(time.RFC3339, createdTimeString)
var updatedTime, _ = time.Parse(time.RFC3339, updatedTimeString)

var (
	Cluster1 = clusters.Cluster{
		ClusterList: &clusters.ClusterList{
			HealthStatus:   nil,
			MasterFlavorID: "g1-standard-1-2",
			FlavorID:       "g1-standard-1-2",
			Name:           "fokgkcytgg",
			Labels: &map[string]string{
				"gcloud_project_id":    "12",
				"gcloud_region_id":     "1",
				"gcloud_access_token":  "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiZXhwIjoxNTg4MzM1NjM3LCJqdGkiOiJhYTE2ODhmODdmNDc1YjhmNDk3NTY5MmI5MTkyZDdmYiIsInVzZXJfdHlwZSI6ImNvbW1vbiIsInVzZXJfZ3JvdXBzIjpbIlVzZXJzIl0sImNsaWVudF9pZCI6LTIwMDMwMjEyMjAzNzUzNTI3NiwicmVnaW9uX2lkIjoxLCJwcm9qZWN0X2lkIjoxMiwidXNlcl9pZCI6MSwiaXNfYWRtaW4iOmZhbHNlfQ.jdPEAMuZOEqT9Ns1eW0IOZmo33WZsMEIs8NFXuF29iU",
				"gcloud_refresh_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImV4cCI6MTg5ODUxMTYzNywianRpIjoiYWExNjg4Zjg3ZjQ3NWI4ZjQ5NzU2OTJiOTE5MmQ3ZmIiLCJ1c2VyX3R5cGUiOiJjb21tb24iLCJ1c2VyX2dyb3VwcyI6WyJVc2VycyJdLCJjbGllbnRfaWQiOi0yMDAzMDIxMjIwMzc1MzUyNzYsInJlZ2lvbl9pZCI6MSwicHJvamVjdF9pZCI6MTIsInVzZXJfaWQiOjEsImlzX2FkbWluIjpmYWxzZX0.AS2Xv067CIxbJdjMB7Z4ydCdxEKwlRx_rLoKNheL0ks",
				"boot_volume_size":     "10",
				"auto_scaling_enabled": "true",
			},
			KeyPair: "keypair",
			Links: []gcorecloud.Link{{
				Href: "http://10.100.178.165:9511/clusters/e4028530-0353-494b-a84c-0230122e34ff",
				Rel:  "self",
			}, {
				Href: "http://10.100.178.165:9511/clusters/e4028530-0353-494b-a84c-0230122e34ff",
				Rel:  "bookmark",
			}},
			Status:            "CREATE_IN_PROGRESS",
			ClusterTemplateID: "0a5ce9dd-a484-4e23-80c7-7e586c80d9fc",
			CreateTimeout:     360,
			UUID:              "e4028530-0353-494b-a84c-0230122e34ff",
			MasterCount:       1,
			DockerVolumeSize:  5,
			NodeCount:         1,
			StackID:           "78c48153-fa6c-48b8-aae3-08b5b230387a",
		},
		DiscoveryURL:       "https://discovery.etcd.io/161d73fde241377395f481c6276b42c7",
		FixedSubnet:        nil,
		CoeVersion:         nil,
		MasterAddresses:    []string{},
		ProjectId:          "ec0f251d-2e36-436c-9a30-7e2c33297273",
		CreatedAt:          createdTime,
		ContainerVersion:   nil,
		StatusReason:       nil,
		HealthStatusReason: map[string]string{},
		ApiAddress:         nil,
		UserId:             "8ba64372-1585-4808-b422-7a7aab5f3197",
		NodeAddresses:      []string{},
		UpdatedAt:          &updatedTime,
		FixedNetwork:       nil,
		FloatingIpEnabled:  false,
	}
	ClusterList1 = clusters.ClusterList{
		UUID:              "e4028530-0353-494b-a84c-0230122e34ff",
		Name:              "fokgkcytgg",
		ClusterTemplateID: "0a5ce9dd-a484-4e23-80c7-7e586c80d9fc",
		KeyPair:           "keypair",
		NodeCount:         1,
		MasterCount:       1,
		DockerVolumeSize:  5,
		Labels: &map[string]string{
			"gcloud_project_id":    "12",
			"gcloud_region_id":     "1",
			"gcloud_access_token":  "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiZXhwIjoxNTg4MzM1NjM3LCJqdGkiOiJhYTE2ODhmODdmNDc1YjhmNDk3NTY5MmI5MTkyZDdmYiIsInVzZXJfdHlwZSI6ImNvbW1vbiIsInVzZXJfZ3JvdXBzIjpbIlVzZXJzIl0sImNsaWVudF9pZCI6LTIwMDMwMjEyMjAzNzUzNTI3NiwicmVnaW9uX2lkIjoxLCJwcm9qZWN0X2lkIjoxMiwidXNlcl9pZCI6MSwiaXNfYWRtaW4iOmZhbHNlfQ.jdPEAMuZOEqT9Ns1eW0IOZmo33WZsMEIs8NFXuF29iU",
			"gcloud_refresh_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImV4cCI6MTg5ODUxMTYzNywianRpIjoiYWExNjg4Zjg3ZjQ3NWI4ZjQ5NzU2OTJiOTE5MmQ3ZmIiLCJ1c2VyX3R5cGUiOiJjb21tb24iLCJ1c2VyX2dyb3VwcyI6WyJVc2VycyJdLCJjbGllbnRfaWQiOi0yMDAzMDIxMjIwMzc1MzUyNzYsInJlZ2lvbl9pZCI6MSwicHJvamVjdF9pZCI6MTIsInVzZXJfaWQiOjEsImlzX2FkbWluIjpmYWxzZX0.AS2Xv067CIxbJdjMB7Z4ydCdxEKwlRx_rLoKNheL0ks",
			"boot_volume_size":     "10",
			"auto_scaling_enabled": "true",
		},
		MasterFlavorID: "g1-standard-1-2",
		FlavorID:       "g1-standard-1-2",
		CreateTimeout:  360,
		Links: []gcorecloud.Link{{
			Href: "http://10.100.178.165:9511/clusters/e4028530-0353-494b-a84c-0230122e34ff",
			Rel:  "self",
		}, {
			Href: "http://10.100.178.165:9511/clusters/e4028530-0353-494b-a84c-0230122e34ff",
			Rel:  "bookmark",
		}},
		Status:       "CREATE_IN_PROGRESS",
		StackID:      "78c48153-fa6c-48b8-aae3-08b5b230387a",
		HealthStatus: nil,
	}
	Tasks1 = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}

	ExpectedClusterSlice = []clusters.ClusterList{ClusterList1}
)