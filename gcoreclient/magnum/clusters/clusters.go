package clusters

import (
	"fmt"
	"strconv"
	"strings"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/magnum/v1/types"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/magnum/v1/clusters"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/task/v1/tasks"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/flags"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/utils"

	"github.com/urfave/cli/v2"
)

var (
	clusterIDText          = "cluster_id is mandatory argument"
	clusterUpdateTypes     = types.ClusterUpdateOperation("").StringList()
	k8sClusterVersionTypes = types.K8sClusterVersion("").StringList()
)

var clusterListSubCommand = cli.Command{
	Name:     "list",
	Usage:    "Magnum list clusters",
	Category: "cluster",
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "magnum", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := clusters.ListAll(client, clusters.ListOpts{})
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var clusterDeleteSubCommand = cli.Command{
	Name:      "delete",
	Usage:     "Magnum delete cluster",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Flags:     flags.WaitCommandFlags,
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstArg(c, clusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := utils.BuildClient(c, "magnum", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := clusters.Delete(client, clusterID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, false, func(task tasks.TaskID) (interface{}, error) {
			_, err := clusters.Get(client, clusterID).Extract()
			if err == nil {
				return nil, fmt.Errorf("cannot delete cluster with ID: %s", clusterID)
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

var clusterResizeSubCommand = cli.Command{
	Name:      "resize",
	Usage:     "Magnum resize cluster",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Flags: append([]cli.Flag{
		&cli.IntFlag{
			Name:     "node-count",
			Usage:    "Cluster nodes count",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:        "nodes-to-remove",
			Usage:       "Cluster nodes chose to remove",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "nodegroup",
			Usage:       "Cluster nodegroup",
			DefaultText: "nil",
			Required:    false,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstArg(c, clusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "resize")
			return err
		}
		client, err := utils.BuildClient(c, "magnum", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		nodes := c.StringSlice("nodes-to-remove")
		if len(nodes) == 0 {
			nodes = nil
		}

		opts := clusters.ResizeOpts{
			NodeCount:     c.Int("node-count"),
			NodesToRemove: nodes,
			NodeGroup:     utils.StringToPointer(c.String("nodegroup")),
		}

		results, err := clusters.Resize(client, clusterID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			clusterID, err := clusters.ExtractClusterIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve cluster ID from task info: %w", err)
			}
			cluster, err := clusters.Get(client, clusterID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get cluster with ID: %s. Error: %w", clusterID, err)
			}
			utils.ShowResults(cluster, c.String("format"))
			return nil, nil
		})

	},
}

var clusterUpgradeSubCommand = cli.Command{
	Name:      "upgrade",
	Usage:     "Magnum upgrade cluster",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "cluster-template",
			Usage:    "Cluster template",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "max-batch-size",
			Usage:    "Max batch size during upgrade",
			Required: false,
		},
		&cli.StringFlag{
			Name:        "nodegroup",
			Usage:       "Cluster nodegroup",
			DefaultText: "nil",
			Required:    false,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstArg(c, clusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "upgrade")
			return err
		}
		client, err := utils.BuildClient(c, "magnum", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := clusters.UpgradeOpts{
			ClusterTemplate: c.String("cluster-template"),
			MaxBatchSize:    utils.IntToPointer(c.Int("max-batch-size")),
			NodeGroup:       utils.StringToPointer(c.String("nodegroup")),
		}

		results, err := clusters.Upgrade(client, clusterID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			clusterID, err := clusters.ExtractClusterIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve cluster ID from task info: %w", err)
			}
			network, err := clusters.Get(client, clusterID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get cluster with ID: %s. Error: %w", clusterID, err)
			}
			utils.ShowResults(network, c.String("format"))
			return nil, nil
		})

	},
}

var clusterUpdateSubCommand = cli.Command{
	Name:      "update",
	Usage:     "Magnum update cluster",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Flags: append([]cli.Flag{
		&cli.StringSliceFlag{
			Name:     "path",
			Aliases:  []string{"p"},
			Usage:    "Update json path. Example /node/count",
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
				Enum: clusterUpdateTypes,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(clusterUpdateTypes, ", ")),
			Required: false,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstArg(c, clusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return cli.NewExitError(err, 1)
		}
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

		var opts clusters.UpdateOpts

		for idx, path := range paths {
			if !strings.HasPrefix(path, "/") {
				return cli.NewExitError(fmt.Errorf("path parameter should be in format /value"), 1)
			}
			var updateValue interface{}
			value := values[idx]
			intValue, err := strconv.Atoi(value)
			if err == nil {
				updateValue = intValue
			} else {
				updateValue = value
			}
			el := clusters.UpdateOptsElem{
				Path:  path,
				Value: updateValue,
				Op:    types.ClusterUpdateOperation(ops[idx]),
			}
			opts = append(opts, el)
		}

		results, err := clusters.Update(client, clusterID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			clusterID, err := clusters.ExtractClusterIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve cluster ID from task info: %w", err)
			}
			network, err := clusters.Get(client, clusterID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get cluster with ID: %s. Error: %w", clusterID, err)
			}
			utils.ShowResults(network, c.String("format"))
			return nil, nil
		})

	},
}

var clusterGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "Magnum get cluster",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstArg(c, clusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := utils.BuildClient(c, "magnum", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := clusters.Get(client, clusterID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var clusterConfigSubCommand = cli.Command{
	Name:      "config",
	Usage:     "Magnum get cluster config",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:     "save",
			Aliases:  []string{"s"},
			Usage:    "Save k8s config in file",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "force",
			Usage:    "Force rewrite KUBECONFIG file",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "merge",
			Aliases:  []string{"m"},
			Usage:    "Merge into existing KUBECONFIG file",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "file",
			Aliases:  []string{"c"},
			Usage:    "KUBECONFIG file",
			Value:    "~/.kube/config",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstArg(c, clusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "config")
			return err
		}
		client, err := utils.BuildClient(c, "magnum", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		result, err := clusters.GetConfig(client, clusterID).ExtractConfig()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if c.Bool("save") {
			filename := c.String("file")
			exists, err := utils.FileExists(filename)
			if err != nil {
				_ = cli.ShowCommandHelp(c, "config")
				return cli.NewExitError(err, 1)
			}
			if exists {
				merge := c.Bool("merge")
				force := c.Bool("force")
				if (!force && !merge) || (force && merge) {
					_ = cli.ShowCommandHelp(c, "config")
					return cli.NewExitError(fmt.Errorf("either --force or --merge shoud be set"), 1)
				}
				if force {
					err := utils.WriteKubeconfigFile(filename, []byte(result.Config))
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				}
				if merge {
					err := utils.MergeKubeconfigFile(filename, []byte(result.Config))
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				}
			} else {
				err := utils.WriteToFile(filename, []byte(result.Config))
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			}
		} else {
			utils.ShowResults(result, c.String("format"))
		}
		return nil
	},
}

var clusterCreateSubCommand = cli.Command{
	Name:     "create",
	Usage:    "Magnum create cluster",
	Category: "cluster",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Cluster name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "template-id",
			Aliases:  []string{"t"},
			Usage:    "Cluster template ID",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "node-count",
			Usage:    "Worker nodes count",
			Value:    1,
			Required: false,
		},
		&cli.IntFlag{
			Name:     "master-node-count",
			Usage:    "Master nodes count",
			Value:    1,
			Required: false,
		},
		&cli.StringFlag{
			Name:        "keypair",
			Aliases:     []string{"k"},
			Usage:       "The name of the SSH keypair",
			Value:       "",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "flavor",
			Usage:       "Worker node flavor",
			Value:       "",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "master-flavor",
			Usage:       "Master node flavor",
			Value:       "",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.StringSliceFlag{
			Name:        "labels",
			Usage:       "Arbitrary labels. The accepted keys and valid values are defined in the cluster drivers. --labels one=two --labels three=four ",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "fixed-subnet",
			Usage:       "Fixed subnet that are using to allocate network address for nodes in cluster.",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "fixed-network",
			Usage:       "Fixed subnet that are using to allocate network address for nodes in cluster.",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.BoolFlag{
			Name:     "floating-ip-enabled",
			Usage:    "Enable fixed IP for cluster nodes.",
			Required: false,
		},
		&cli.GenericFlag{
			Name: "version",
			Value: &utils.EnumValue{
				Enum:    k8sClusterVersionTypes,
				Default: types.K8sClusterVersion117.String(),
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(k8sClusterVersionTypes, ", ")),
			Required: false,
		},
		&cli.IntFlag{
			Name:     "create-timeout",
			Usage:    "Heat timeout to create cluster. Seconds",
			Value:    7200,
			Required: false,
		},
	}, flags.WaitCommandFlags...,
	),
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "magnum", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		labels, err := utils.StringSliceToMap(c.StringSlice("labels"))
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		opts := clusters.CreateOpts{
			Name:              c.String("name"),
			ClusterTemplateID: c.String("template-id"),
			NodeCount:         c.Int("node-count"),
			MasterCount:       c.Int("master-node-count"),
			KeyPair:           c.String("keypair"),
			FlavorID:          c.String("flavor"),
			MasterFlavorID:    c.String("master-flavor"),
			Labels:            &labels,
			FixedNetwork:      c.String("fixed-network"),
			FixedSubnet:       c.String("fixed-subnet"),
			FloatingIPEnabled: c.Bool("floating-ip-enabled"),
			CreateTimeout:     c.Int("create-timeout"),
			Version:           types.K8sClusterVersion(c.String("version")),
		}

		results, err := clusters.Create(client, opts).Extract()
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
			clusterID, err := clusters.ExtractClusterIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve cluster ID from task info: %w", err)
			}
			network, err := clusters.Get(client, clusterID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get cluster with ID: %s. Error: %w", clusterID, err)
			}
			utils.ShowResults(network, c.String("format"))
			return nil, nil
		})
	},
}

var ClusterCommands = cli.Command{
	Name:  "cluster",
	Usage: "Magnum cluster commands",
	Subcommands: []*cli.Command{
		&clusterListSubCommand,
		&clusterDeleteSubCommand,
		&clusterGetSubCommand,
		&clusterCreateSubCommand,
		&clusterResizeSubCommand,
		&clusterUpgradeSubCommand,
		&clusterUpdateSubCommand,
		&clusterConfigSubCommand,
	},
}
