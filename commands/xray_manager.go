package commands

import (
	"github.com/jfrog/jfrog-cli-core/v2/plugins/common"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	xrayutils "github.com/jfrog/jfrog-cli-core/v2/utils/xray"
	"github.com/jfrog/jfrog-client-go/xray"
)

// createXrayManager builds an authenticated XrayServicesManager from the
// plugin context. It reads the --server-id flag (or uses the default CLI
// config), refreshes tokens, and wires up the managed HTTP client with
// retries, TLS, and auth interceptors.
func createXrayManager(c *components.Context) (*xray.XrayServicesManager, error) {
	details, err := common.GetServerDetails(c)
	if err != nil {
		return nil, err
	}
	return xrayutils.CreateXrayServiceManager(details)
}
