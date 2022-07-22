package actions

import (
	"time"

	"github.com/openziti/fablab/kernel/lib/actions"
	"github.com/openziti/fablab/kernel/lib/actions/component"
	"github.com/openziti/fablab/kernel/lib/actions/host"
	"github.com/openziti/fablab/kernel/lib/actions/semaphore"
	"github.com/openziti/fablab/kernel/model"
	zitilib_actions "github.com/openziti/zitilab/actions"
	"github.com/openziti/zitilab/actions/edge"
	"github.com/openziti/zitilab/models"
)

type bootstrapAction struct{}

func NewBootstrapAction() model.ActionBinder {
	action := &bootstrapAction{}
	return action.bind
}

func (a *bootstrapAction) bind(m *model.Model) model.Action {
	workflow := actions.Workflow()

	workflow.AddAction(host.GroupExec("*", 25, "rm -f logs/*"))
	workflow.AddAction(component.Stop("#ctrl"))
	workflow.AddAction(edge.InitController("#ctrl"))
	workflow.AddAction(component.Start("#ctrl"))
	workflow.AddAction(semaphore.Sleep(2 * time.Second))

	workflow.AddAction(edge.Login("#ctrl"))

	workflow.AddAction(component.StopInParallel(models.EdgeRouterTag, 25))
	workflow.AddAction(edge.InitEdgeRouters(models.EdgeRouterTag, 2))
	workflow.AddAction(edge.InitIdentities(models.SdkAppTag, 2))

	workflow.AddAction(zitilib_actions.Edge("create", "service", "echo"))

	workflow.AddAction(zitilib_actions.Edge("create", "service-policy", "echo-servers", "Bind", "--service-roles", "@echo", "--identity-roles", "#service"))
	workflow.AddAction(zitilib_actions.Edge("create", "service-policy", "echo-client", "Dial", "--service-roles", "@echo", "--identity-roles", "#client"))

	workflow.AddAction(zitilib_actions.Edge("create", "edge-router-policy", "echo-servers", "--edge-router-roles", "@router-east", "--identity-roles", "#service"))
	workflow.AddAction(zitilib_actions.Edge("create", "edge-router-policy", "echo-clients", "--edge-router-roles", "@router-west", "--identity-roles", "#client"))

	workflow.AddAction(zitilib_actions.Edge("create", "service-edge-router-policy", "echo", "--semantic", "AnyOf", "--service-roles", "@echo", "--edge-router-roles", "#all"))

	workflow.AddAction(component.Stop(models.ControllerTag))

	return workflow
}
