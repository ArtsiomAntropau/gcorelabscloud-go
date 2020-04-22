package nodegroups

import (
	"fmt"
	"strconv"
	"strings"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/magnum/v1/types"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/magnum/v1/nodegroups"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/task/v1/tasks"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/flags"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/utils"

	"github.com/urfave/cli/v2"
)

var (
	nodeGroupIDText       = "nodegroup_id is mandatory argument"
	clusterIDText         = "cluster_id is mandatory argument"
	nodegroupsUpdateTypes = types.ClusterUpdateOperation("").StringList()
	nodegroupRoles        = types.NodegroupRole("").StringList()
)

var nodegroupListSubCommand = cli.Command{
	Name:      "list",
	Usage:     "Magnum list nodegroups",
	Category:  "nodegroup",
	ArgsUsage: "<cluster_id>",
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, clusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "list")
			return cli.NewExitError(err, 1)
		}
		client, err := utils.BuildClient(c, "magnum", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		pages, err := nodegroups.List(client, clusterID, nodegroups.ListOpts{}).AllPages()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		results, err := nodegroups.ExtractClusterNodeGroups(pages)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var nodegroupDeleteSubCommand = cli.Command{
	Name:      "delete",
	Usage:     "Magnum delete nodegroup",
	ArgsUsage: "<nodegroup_id>",
	Category:  "nodegroup",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "cluster-id",
			Aliases:  []string{"c"},
			Usage:    "Magnum cluster ID",
			Required: true,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		nodeGroupID, err := flags.GetFirstStringArg(c, nodeGroupIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return cli.NewExitError(err, 1)
		}
		clusterID := c.String("cluster-id")
		client, err := utils.BuildClient(c, "magnum", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := nodegroups.Delete(client, clusterID, nodeGroupID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, false, func(task tasks.TaskID) (interface{}, error) {
			_, err := nodegroups.Get(client, clusterID, nodeGroupID).Extract()
			if err == nil {
				return nil, fmt.Errorf("cannot delete nodegroup with ID: %s", nodeGroupID)
			}
			switch err.(type) {
			case gcorecloud.ErrDefault404:
				return nil, nil
			default:
				return nil, err
			}
		})

	},
}

var nodegroupUpdateSubCommand = cli.Command{
	Name:      "update",
	Usage:     "Magnum update nodegroup",
	ArgsUsage: "<nodegroup_id>",
	Category:  "nodegroup",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "cluster-id",
			Aliases:  []string{"c"},
			Usage:    "Magnum cluster ID",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:     "path",
			Aliases:  []string{"p"},
			Usage:    "Update json path. Example /node_count",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:     "value",
			Aliases:  []string{"v"},
			Usage:    "Update json value",
			Required: true,
		},
		&cli.GenericFlag{
			Name:    "op",
			Aliases: []string{"o"},
			Value: &utils.EnumStringSliceValue{
				Enum: nodegroupsUpdateTypes,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(nodegroupsUpdateTypes, ", ")),
			Required: true,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		nodeGroupID, err := flags.GetFirstStringArg(c, nodeGroupIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return cli.NewExitError(err, 1)
		}
		clusterID := c.String("cluster-id")
		client, err := utils.BuildClient(c, "magnum", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		paths := c.StringSlice("path")
		values := c.StringSlice("value")
		ops := utils.GetEnumStringSliceValue(c, "op")

		if len(paths) != len(values) || len(values) != len(ops) {
			_ = cli.ShowCommandHelp(c, "update")
			return cli.NewExitError(fmt.Errorf("path, value and op parameters number should be same"), 1)
		}

		var opts nodegroups.UpdateOpts

		for idx, path := range paths {
			if !strings.HasPrefix(path, "/") {
				return cli.NewExitError(fmt.Errorf("path parameter should be in format /path"), 1)
			}
			var updateValue interface{}
			value := values[idx]
			intValue, err := strconv.Atoi(value)
			if err == nil {
				updateValue = intValue
			} else if path == "/labels" {
				updateValue, err = utils.StringSliceToMap(strings.Split(value, ","))
				if err != nil {
					return cli.NewExitError(fmt.Errorf("wrong labels format. should be in format: label_one=value_one,label_two=value_two"), 1)
				}
			} else {
				updateValue = value
			}
			el := nodegroups.UpdateOptsElem{
				Path:  path,
				Value: updateValue,
				Op:    types.ClusterUpdateOperation(ops[idx]),
			}
			opts = append(opts, el)
		}
		result, err := nodegroups.Update(client, clusterID, nodeGroupID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var nodegroupGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "Magnum get nodegroup",
	ArgsUsage: "<nodegroup_id>",
	Category:  "nodegroup",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "cluster-id",
			Aliases:  []string{"c"},
			Usage:    "Magnum cluster ID",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		nodeGroupID, err := flags.GetFirstStringArg(c, nodeGroupIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return cli.NewExitError(err, 1)
		}
		clusterID := c.String("cluster-id")
		client, err := utils.BuildClient(c, "magnum", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := nodegroups.Get(client, clusterID, nodeGroupID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var nodegroupCreateSubCommand = cli.Command{
	Name:      "create",
	Usage:     "Magnum create nodegroup",
	Category:  "nodegroup",
	ArgsUsage: "<cluster_id>",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Cluster nodegroup name",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "node-count",
			Usage:    "Node count",
			Aliases:  []string{"c"},
			Value:    1,
			Required: false,
		},
		&cli.StringFlag{
			Name:     "flavor-id",
			Usage:    "Node flavor",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "image-id",
			Usage:    "Node image",
			Required: true,
		},
		&cli.GenericFlag{
			Name:    "role",
			Aliases: []string{"r"},
			Value: &utils.EnumValue{
				Enum: nodegroupRoles,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(nodegroupRoles, ", ")),
			Required: false,
		},
		&cli.IntFlag{
			Name:        "docker-volume-size",
			Usage:       "The size in GB for the local storage on each server for the Docker daemon to cache the images and host the containers",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.IntFlag{
			Name:        "min-node-count",
			Usage:       "Minimum number of nodes",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.IntFlag{
			Name:        "max-node-count",
			Usage:       "Maximum number of nodes",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.StringSliceFlag{
			Name:        "labels",
			Usage:       "Arbitrary labels. The accepted keys and valid values are defined in the nodegroup drivers. --labels one=two --labels three=four ",
			DefaultText: "nil",
			Required:    false,
		},
	}, flags.WaitCommandFlags...,
	),
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, clusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}
		client, err := utils.BuildClient(c, "magnum", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		labels, err := utils.StringSliceToMap(c.StringSlice("labels"))
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		role, err := types.NodegroupRole(c.String("role")).ValidOrNil()
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}
		opts := nodegroups.CreateOpts{
			Name:             c.String("name"),
			FlavorID:         c.String("flavor-id"),
			ImageID:          c.String("image-id"),
			Role:             role,
			NodeCount:        c.Int("node-count"),
			DockerVolumeSize: utils.IntToPointer(c.Int("docker-volume-size")),
			Labels:           &labels,
			MinNodeCount:     utils.IntToPointer(c.Int("min-node-count")),
			MaxNodeCount:     utils.IntToPointer(c.Int("max-node-count")),
		}

		results, err := nodegroups.Create(client, clusterID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if results == nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			nodegroupID, err := nodegroups.ExtractClusterNodeGroupIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve nodegroup ID from task info: %w", err)
			}
			nodegroup, err := nodegroups.Get(client, clusterID, nodegroupID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get nodegroup with ID: %s. Error: %w", nodegroupID, err)
			}
			utils.ShowResults(nodegroup, c.String("format"))
			return nil, nil
		})
	},
}

var ClusterNodeGroupCommands = cli.Command{
	Name:  "nodegroup",
	Usage: "Magnum nodegroup commands",
	Subcommands: []*cli.Command{
		&nodegroupListSubCommand,
		&nodegroupDeleteSubCommand,
		&nodegroupGetSubCommand,
		&nodegroupCreateSubCommand,
		&nodegroupUpdateSubCommand,
	},
}
