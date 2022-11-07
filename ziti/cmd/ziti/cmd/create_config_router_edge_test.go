package cmd

import (
	"github.com/openziti/ziti/ziti/cmd/ziti/constants"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var defaultArgs = []string{"edge", "--routerName", "test-router"}
var testHostname, _ = os.Hostname()

func setEnvByMap[K string, V string](m map[K]V) {
	for k, v := range m {
		os.Setenv(string(k), string(v))
	}
}

func execCreateConfigCommand(args []string, keys map[string]string) {
	// Setup options
	clearOptionsAndTemplateData()
	routerOptions.Output = defaultOutput

	setEnvByMap(keys)
	// Create and run the CLI command (capture output, otherwise config prints to stdout instead of test results)
	cmd := NewCmdCreateConfigRouter()
	cmd.SetArgs(args)
	_ = captureOutput(func() {
		_ = cmd.Execute()
	})
}

func TestEdgeRouterAdvertised(t *testing.T) {
	clearOptionsAndTemplateData()
	routerAdvHostIp := "192.168.10.10"
	routerAdvHostDns := "controller01.zitinetwork.example.org"
	keys := map[string]string{
		"ZITI_CTRL_PORT":        "80",
		"ZITI_EDGE_ROUTER_PORT": "443",
	}
	execCreateConfigCommand(defaultArgs, keys)
	assert.Equal(t, testHostname, data.Router.Edge.AdvertisedHost, nil)

	keys["ZITI_EDGE_ROUTER_RAWNAME"] = routerAdvHostDns
	execCreateConfigCommand(defaultArgs, keys)
	assert.Equal(t, routerAdvHostDns, data.Router.Edge.AdvertisedHost, nil)

	keys["ZITI_EDGE_ROUTER_RAWNAME"] = ""
	keys["ZITI_EDGE_ROUTER_IP_OVERRIDE"] = routerAdvHostIp
	execCreateConfigCommand(defaultArgs, keys)
	assert.Equal(t, routerAdvHostIp, data.Router.Edge.AdvertisedHost, nil)

	keys["ZITI_EDGE_ROUTER_IP_OVERRIDE"] = ""
	keys["EXTERNAL_DNS"] = routerAdvHostDns
	execCreateConfigCommand(defaultArgs, keys)
	assert.Equal(t, routerAdvHostDns, data.Router.Edge.AdvertisedHost, nil)

	keys["ZITI_EDGE_ROUTER_ADVERTISED_HOST"] = routerAdvHostIp
	keys["ZITI_EDGE_ROUTER_IP_OVERRIDE"] = routerAdvHostIp
	execCreateConfigCommand(defaultArgs, keys)
	assert.Equal(t, routerAdvHostIp, data.Router.Edge.AdvertisedHost, nil)

	keys["ZITI_EDGE_ROUTER_ADVERTISED_HOST"] = routerAdvHostDns
	keys["EXTERNAL_DNS"] = routerAdvHostDns
	execCreateConfigCommand(defaultArgs, keys)
	assert.Equal(t, routerAdvHostDns, data.Router.Edge.AdvertisedHost, nil)
}

func TestTunnelerEnabledByDefault(t *testing.T) {
	// Setup options
	clearOptionsAndTemplateData()
	routerOptions.Output = defaultOutput

	// Create and run the CLI command without the tunnel flag
	config := createRouterConfig([]string{"edge", "--routerName", "myRouter"})

	// Expect disabled flag to be false
	assert.False(t, data.Router.TunnelerDisabled, "Expected tunneler disabled flag to be false")

	// Confirm tunneler is enabled in config output
	foundTunnel := false
	for i := 0; i < len(config.Listeners); i++ {
		if config.Listeners[i].Binding == "tunnel" {
			foundTunnel = true
		}
	}
	assert.True(t, foundTunnel, "Expected to find tunnel listener binding but it was not found")
}

func TestTunnelerDisabledFlag(t *testing.T) {
	// Setup options
	clearOptionsAndTemplateData()
	routerOptions.Output = defaultOutput

	// Create and run the CLI command with the disable tunnel flag
	config := createRouterConfig([]string{"edge", "--routerName", "myRouter", "--disableTunneler"})

	// Expect disabled flag to be true
	assert.True(t, data.Router.TunnelerDisabled, "Expected tunneler disabled flag to be true")

	// Confirm tunneler is enabled in config output
	foundTunnel := false
	for i := 0; i < len(config.Listeners); i++ {
		if config.Listeners[i].Binding == "tunnel" {
			foundTunnel = true
		}
	}
	assert.False(t, foundTunnel, "Tunnel listener binding was not expected but it was found")
}

func TestTunnelerHostModeIsDefault(t *testing.T) {
	// Setup options
	clearOptionsAndTemplateData()
	routerOptions.Output = defaultOutput

	// Create and run the CLI command without the tunnel flag
	config := createRouterConfig([]string{"edge", "--routerName", "myRouter"})

	// Expect tunneler mode to be "host" mode
	assert.Equal(t, hostMode, data.Router.TunnelerMode, "Expected tunneler mode to be %s but found %s", hostMode, data.Router.TunnelerMode)

	// Confirm tunneler mode in config is set to host mode
	for i := 0; i < len(config.Listeners); i++ {
		if config.Listeners[i].Binding == "tunnel" {
			assert.Equal(t, hostMode, config.Listeners[i].Options.Mode, "Expected tunneler mode to be %s but found %s in config output", hostMode, data.Router.TunnelerMode)
		}
	}
}

func TestTunnelerTproxyMode(t *testing.T) {
	// Setup options
	clearOptionsAndTemplateData()
	routerOptions.Output = defaultOutput

	// Create and run the CLI command without the tunnel flag
	config := createRouterConfig([]string{"edge", "--routerName", "myRouter", "--tunnelerMode", tproxyMode})

	// Expect tunneler mode to be "host" mode
	assert.Equal(t, tproxyMode, data.Router.TunnelerMode, "Expected tunneler mode to be %s but found %s", tproxyMode, data.Router.TunnelerMode)

	// Confirm tunneler mode in config is set to host mode
	for i := 0; i < len(config.Listeners); i++ {
		if config.Listeners[i].Binding == "tunnel" {
			assert.Equal(t, tproxyMode, config.Listeners[i].Options.Mode, "Expected tunneler mode to be %s but found %s in config output", tproxyMode, data.Router.TunnelerMode)
		}
	}
}

func TestTunnelerInvalidMode(t *testing.T) {
	invalidMode := "invalidMode"

	expectedErrorMsg := "Unknown tunneler mode [" + invalidMode + "] provided, should be [host] or [tproxy]"

	// Create the options with both flags set to true
	clearOptionsAndTemplateData()
	routerOptions.Output = defaultOutput
	routerOptions.TunnelerMode = invalidMode

	err := routerOptions.runEdgeRouter(&ConfigTemplateValues{})

	assert.EqualError(t, err, expectedErrorMsg, "Error does not match, expected %s but got %s", expectedErrorMsg, err)
}

func TestPrivateEdgeRouterNotAdvertising(t *testing.T) {
	clearOptionsAndTemplateData()

	// Create and run the CLI command
	config := createRouterConfig([]string{"edge", "--routerName", "myRouter", "--private"})

	// Expect that the config values are represented correctly
	assert.Equal(t, 0, len(config.Link.Listeners), "Expected zero link listeners for private edge router, found a non-zero value")
}

func TestBlankEdgeRouterNameBecomesHostname(t *testing.T) {
	hostname, _ := os.Hostname()
	blank := ""

	// Setup options with blank router name
	clearOptionsAndTemplateData()
	routerOptions.Output = defaultOutput
	routerOptions.RouterName = blank

	// Check that template values is a blank name
	assert.Equal(t, blank, data.Router.Name, "Mismatch router name, expected %s but got %s", "", hostname)

	// Create and run the CLI command
	_ = createRouterConfig([]string{"edge", "--routerName", blank})

	// Check that the blank name was replaced with hostname in the template values
	assert.Equal(t, hostname, data.Router.Name, "Mismatch router name, expected %s but got %s", "", hostname)

}

func TestEdgeRouterCannotBeWSSAndPrivate(t *testing.T) {
	expectedErrorMsg := "Flags for private and wss configs are mutually exclusive. You must choose private or wss, not both"

	// Create the options with both flags set to true
	clearOptionsAndTemplateData()
	routerOptions.Output = defaultOutput
	routerOptions.IsPrivate = true
	routerOptions.WssEnabled = true

	err := routerOptions.runEdgeRouter(&ConfigTemplateValues{})

	assert.EqualError(t, err, expectedErrorMsg, "Error does not match, expected %s but got %s", expectedErrorMsg, err)
}

func TestEdgeRouterOutputPathDoesNotExist(t *testing.T) {
	expectedErrorMsg := "stat /IDoNotExist: no such file or directory"

	// Set the router options
	clearOptionsAndTemplateData()
	routerOptions.TunnelerDisabled = true
	routerOptions.RouterName = "MyEdgeRouter"
	routerOptions.Output = "/IDoNotExist/MyEdgeRouter.yaml"

	err := routerOptions.runEdgeRouter(&ConfigTemplateValues{})

	assert.EqualError(t, err, expectedErrorMsg, "Error does not match, expected %s but got %s", expectedErrorMsg, err)
}

func TestExecuteCreateConfigRouterEdgeHasNonBlankTemplateValues(t *testing.T) {
	routerName := "MyEdgeRouter"
	expectedNonEmptyStringFields := []string{".ZitiHome", ".Hostname", ".Router.Name", ".Router.IdentityCert", ".Router.IdentityServerCert", ".Router.IdentityKey", ".Router.IdentityCA", ".Router.Edge.Hostname", ".Router.Edge.Port"}
	expectedNonEmptyStringValues := []*string{&data.ZitiHome, &data.Hostname, &data.Router.Name, &data.Router.IdentityCert, &data.Router.IdentityServerCert, &data.Router.IdentityKey, &data.Router.IdentityCA, &data.Router.Edge.Hostname, &data.Router.Edge.Port}
	expectedNonEmptyIntFields := []string{".Router.Listener.BindPort", ".Router.Listener.OutQueueSize", ".Router.Wss.ReadBufferSize", ".Router.Wss.WriteBufferSize", ".Router.Forwarder.XgressDialQueueLength", ".Router.Forwarder.XgressDialWorkerCount", ".Router.Forwarder.LinkDialQueueLength", ".Router.Forwarder.LinkDialWorkerCount"}
	expectedNonEmptyIntValues := []*int{&data.Router.Listener.BindPort, &data.Router.Listener.OutQueueSize, &data.Router.Wss.ReadBufferSize, &data.Router.Wss.WriteBufferSize, &data.Router.Forwarder.XgressDialQueueLength, &data.Router.Forwarder.XgressDialWorkerCount, &data.Router.Forwarder.LinkDialQueueLength, &data.Router.Forwarder.LinkDialWorkerCount}
	expectedNonEmptyTimeFields := []string{".Router.Listener.ConnectTimeout", "Router.Listener.GetSessionTimeout", ".Router.Wss.WriteTimeout", ".Router.Wss.ReadTimeout", ".Router.Wss.IdleTimeout", ".Router.Wss.PongTimeout", ".Router.Wss.PingInterval", ".Router.Wss.HandshakeTimeout", ".Router.Forwarder.LatencyProbeInterval"}
	expectedNonEmptyTimeValues := []*time.Duration{&data.Router.Listener.ConnectTimeout, &data.Router.Listener.GetSessionTimeout, &data.Router.Wss.WriteTimeout, &data.Router.Wss.ReadTimeout, &data.Router.Wss.IdleTimeout, &data.Router.Wss.PongTimeout, &data.Router.Wss.PingInterval, &data.Router.Wss.HandshakeTimeout, &data.Router.Forwarder.LatencyProbeInterval}

	// Create and run the CLI command
	_ = createRouterConfig([]string{"edge", "--routerName", routerName})

	// Check that the expected string template values are not blank
	for field, value := range expectedNonEmptyStringValues {
		assert.NotEqualf(t, "", *value, expectedNonEmptyStringFields[field]+" should be a non-blank value")
	}

	// Check that the expected int template values are not zero
	for field, value := range expectedNonEmptyIntValues {
		assert.NotZero(t, *value, expectedNonEmptyIntFields[field]+" should be a non-zero value")
	}

	// Check that the expected time.Duration template values are not zero
	for field, value := range expectedNonEmptyTimeValues {
		assert.NotZero(t, *value, expectedNonEmptyTimeFields[field]+" should be a non-zero value")
	}
}

func TestEdgeRouterIPOverrideIsConsumed(t *testing.T) {
	routerName := "MyFabricRouter"
	blank := ""
	externalIP := "123.456.78.9"

	// Setup options
	clearOptionsAndTemplateData()
	routerOptions.Output = defaultOutput

	// Set the env variable to non-empty value
	_ = os.Setenv(constants.ZitiEdgeRouterIPOverrideVarName, externalIP)

	// Check that template value is currently blank
	assert.Equal(t, blank, data.Router.Edge.IPOverride, "Mismatch router IP override, expected %s but got %s", blank, data.Router.Edge.IPOverride)

	// Create and run the CLI command (capture output, otherwise config prints to stdout instead of test results)
	config := createRouterConfig([]string{"edge", "--routerName", routerName})

	// Check that the template values now contains the custom external IP override value
	assert.Equal(t, externalIP, data.Router.Edge.IPOverride, "Mismatch router IP override, expected %s but got %s", externalIP, data.Router.Edge.IPOverride)

	// Check that the config output has the IP
	found := false
	for i := 1; i < len(config.Edge.Csr.Sans.Ip); i++ {
		if config.Edge.Csr.Sans.Ip[i] == externalIP {
			found = true
		}
	}
	assert.True(t, found, "Expected value not found; expected to find value of "+constants.ZitiEdgeRouterIPOverrideVarName+" in edge router config output.")
}
