package main

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/openziti/fablab"
	"github.com/openziti/fablab/kernel/lib/actions"
	"github.com/openziti/fablab/kernel/lib/actions/component"
	"github.com/openziti/fablab/kernel/lib/actions/host"
	"github.com/openziti/fablab/kernel/lib/actions/semaphore"
	"github.com/openziti/fablab/kernel/lib/binding"
	"github.com/openziti/fablab/kernel/lib/runlevel/0_infrastructure/aws_ssh_key"
	"github.com/openziti/fablab/kernel/lib/runlevel/0_infrastructure/semaphore"
	"github.com/openziti/fablab/kernel/lib/runlevel/0_infrastructure/terraform"
	distribution "github.com/openziti/fablab/kernel/lib/runlevel/3_distribution"
	"github.com/openziti/fablab/kernel/lib/runlevel/3_distribution/rsync"
	aws_ssh_key2 "github.com/openziti/fablab/kernel/lib/runlevel/6_disposal/aws_ssh_key"
	"github.com/openziti/fablab/kernel/lib/runlevel/6_disposal/terraform"
	"github.com/openziti/fablab/kernel/model"
	"github.com/openziti/fablab/resources"
	"github.com/openziti/ziti/controller/persistence"
	"github.com/openziti/ziti/zititest/models/test_resources"
	"github.com/openziti/ziti/zititest/zitilab"
	"github.com/openziti/ziti/zititest/zitilab/actions/edge"
	"github.com/openziti/ziti/zititest/zitilab/models"
	"go.etcd.io/bbolt"
	"os"
	"path"
	"strings"
	"time"
)

const TargetZitiVersion = "v0.30.5"

//go:embed configs
var configResource embed.FS

type scaleStrategy struct{}

func (self scaleStrategy) IsScaled(entity model.Entity) bool {
	return entity.GetType() == model.EntityTypeHost && entity.GetScope().HasTag("scaled")
}

func (self scaleStrategy) GetEntityCount(entity model.Entity) uint32 {
	if entity.GetType() == model.EntityTypeHost && entity.GetScope().HasTag("scaled") {
		return 4
	}
	return 1
}

type dbStrategy struct{}

func (d dbStrategy) GetDbFile(m *model.Model) string {
	return m.MustStringVariable("db_file")
}

func (d dbStrategy) GetSite(router *persistence.EdgeRouter) (string, bool) {
	if strings.Contains(strings.ToLower(router.Name), "london") {
		return "eu-west-2a", true // london region
	}
	if strings.Contains(strings.ToLower(router.Name), "virginia") {
		return "us-east-1a", true // london region
	}
	if strings.Contains(strings.ToLower(router.Name), "melbourne") {
		return "ap-southeast-2a", true // sydney region
	}

	return "us-east-1a", true
}

func (d dbStrategy) PostProcess(router *persistence.EdgeRouter, c *model.Component) {
	if router.IsTunnelerEnabled {
		c.Scope.Tags = append(c.Scope.Tags, "tunneler")
	}
	c.Scope.Tags = append(c.Scope.Tags, "edge-router")
	c.Scope.Tags = append(c.Scope.Tags, "pre-created")
	c.Host.InstanceType = "c5.large"
	c.Type.(*zitilab.RouterType).Version = TargetZitiVersion
}

func (d dbStrategy) ProcessDbModel(tx *bbolt.Tx, m *model.Model, builder *models.ZitiDbBuilder) error {
	if err := builder.CreateEdgeRouterHosts(tx, m); err != nil {
		return err
	}
	return d.CreateIdentityHosts(tx, m, builder)
}

func (d dbStrategy) CreateIdentityHosts(tx *bbolt.Tx, m *model.Model, builder *models.ZitiDbBuilder) error {
	edgeStores := builder.GetEdgeStores()
	ids, _, err := edgeStores.Identity.QueryIds(tx, "true limit none")
	if err != nil {
		return err
	}

	servicesCount := 0
	hostingIdentities := map[string]int{}

	for _, identityId := range ids {
		cursorProvider := edgeStores.Identity.GetIdentityServicesCursorProvider(identityId)
		cursor := cursorProvider(tx, true)
		identityServiceCount := 0
		for cursor.IsValid() {
			serviceId := string(cursor.Current())
			if edgeStores.EdgeService.IsBindableByIdentity(tx, serviceId, identityId) {
				identityServiceCount++
			}
			cursor.Next()
		}
		if identityServiceCount > 0 {
			servicesCount += identityServiceCount
			hostingIdentities[identityId] = identityServiceCount
		}
	}

	regionCount := len(m.Regions)

	perRegion := servicesCount / regionCount
	idIdx := 0

	avgTunnelsPerHost := 15

	m.RangeSortedRegions(func(regionId string, region *model.Region) {
		regionServiceCount := 0

		var regionIdentityIds []string

		for {
			if idIdx >= len(ids) {
				break
			}
			identityId := ids[idIdx]
			idIdx++

			svcCount, found := hostingIdentities[identityId]
			if !found {
				continue
			}
			regionServiceCount += svcCount
			regionIdentityIds = append(regionIdentityIds, identityId)
			if regionServiceCount > perRegion {
				break
			}
		}

		hostCount := len(regionIdentityIds) / avgTunnelsPerHost
		var hosts []*model.Host

		for i := 0; i < hostCount; i++ {
			tunnelsHost := &model.Host{
				Scope:        model.Scope{Tags: model.Tags{}},
				Region:       region,
				Components:   model.Components{},
				InstanceType: "t3.medium",
			}
			hostId := fmt.Sprintf("%s_svc_hosts_%v", regionId, i)
			region.Hosts[hostId] = tunnelsHost
			hosts = append(hosts, tunnelsHost)
		}

		hostIdx := 0
		for _, identityId := range regionIdentityIds {
			tunnelHost := hosts[hostIdx%len(hosts)]
			hostIdx++

			svcCount := hostingIdentities[identityId]

			tunnelComponent := &model.Component{
				Scope: model.Scope{Tags: model.Tags{"sdk-tunneler", "pre-created", fmt.Sprintf("serviceCount=%v", svcCount)}},
				Type: &zitilab.ZitiTunnelType{
					Mode:    zitilab.ZitiTunnelModeHost,
					Version: TargetZitiVersion,
				},
				Host: tunnelHost,
			}
			tunnelHost.Components[identityId] = tunnelComponent
		}
	})

	return nil
}

var dbStrategyInstance = dbStrategy{}

var m = &model.Model{
	Id: "router-test",
	Scope: model.Scope{
		Defaults: model.Variables{
			"environment": "sdk-hosting-test",
			"credentials": model.Variables{
				"aws": model.Variables{
					"managed_key": true,
				},
				"ssh": model.Variables{
					"username": "ubuntu",
				},
				"edge": model.Variables{
					"username": "admin",
					"password": "admin",
				},
			},
			"metrics": model.Variables{
				"influxdb": model.Variables{
					"url": "http://localhost:8086",
					"db":  "ziti",
				},
			},
		},
	},
	StructureFactories: []model.Factory{
		model.NewScaleFactoryWithDefaultEntityFactory(scaleStrategy{}),
		&models.ZitiDbBuilder{Strategy: dbStrategyInstance},
	},
	Resources: model.Resources{
		resources.Configs:   resources.SubFolder(configResource, "configs"),
		resources.Binaries:  os.DirFS(path.Join(os.Getenv("GOPATH"), "bin")),
		resources.Terraform: test_resources.TerraformResources(),
	},
	Regions: model.Regions{
		"us-east-1": {
			Region: "us-east-1",
			Site:   "us-east-1a",
			Hosts: model.Hosts{
				"ctrl": {
					InstanceType: "c5.large",
					Components: model.Components{
						"ctrl": {
							Scope: model.Scope{Tags: model.Tags{"ctrl"}},
							Type: &zitilab.ControllerType{
								Version: TargetZitiVersion,
							},
						},
					},
				},
			},
		},
	},

	Actions: model.ActionBinders{
		"bootstrap": model.ActionBinder(func(m *model.Model) model.Action {
			workflow := actions.Workflow()

			//workflow.AddAction(component.Stop("*"))
			//workflow.AddAction(host.GroupExec("*", 25, "rm -f logs/*"))

			workflow.AddAction(component.Start("#ctrl"))
			workflow.AddAction(semaphore.Sleep(2 * time.Second))

			workflow.AddAction(edge.Login("#ctrl"))

			workflow.AddAction(edge.ReEnrollEdgeRouters(".edge-router .pre-created", 2))
			workflow.AddAction(edge.ReEnrollIdentities(".sdk-tunneler .pre-created", 10))
			return workflow
		}),
		"stop": model.Bind(component.StopInParallelHostExclusive("*", 15)),
		"clean": model.Bind(actions.Workflow(
			component.StopInParallelHostExclusive("*", 15),
			host.GroupExec("*", 25, "rm -f logs/*"),
		)),
		"login": model.Bind(edge.Login("#ctrl")),
	},

	Infrastructure: model.Stages{
		aws_ssh_key.Express(),
		&terraform_0.Terraform{
			Retries: 3,
			ReadyCheck: &semaphore_0.ReadyStage{
				MaxWait: 90 * time.Second,
			},
		},
	},

	Distribution: model.Stages{
		distribution.DistributeSshKey("*"),
		distribution.Locations("*", "logs"),
		rsync.RsyncStaged(),
		model.StageActionF(func(run model.Run) error {
			dbFile := dbStrategyInstance.GetDbFile(run.GetModel())
			deferred := rsync.NewRsyncHost("#ctrl", dbFile, "/home/ubuntu/fablab/ctrl.db")
			return deferred.Execute(run)
		}),
	},

	Disposal: model.Stages{
		terraform.Dispose(),
		aws_ssh_key2.Dispose(),
	},
}

func main() {
	m.AddActivationActions("stop", "bootstrap")

	model.AddBootstrapExtension(binding.AwsCredentialsLoader)
	model.AddBootstrapExtension(aws_ssh_key.KeyManager)

	fablab.InitModel(m)
	fablab.Run()
}
