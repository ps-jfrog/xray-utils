package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/jfrog/jfrog-cli-core/v2/plugins/common"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const ignoreRulesAPI = "api/v1/ignore_rules"

func GetCreateIgnoreRuleCommand() components.Command {
	return components.Command{
		Name:        "create-ignore-rule",
		Description: "Create an Xray ignore rule from a JSON file (body of the request).",
		Aliases:     []string{"cir"},
		Arguments:   getCreateIgnoreRuleArguments(),
		Flags:       getCreateIgnoreRuleFlags(),
		Action: func(c *components.Context) error {
			return createIgnoreRuleCmd(c)
		},
	}
}

func getCreateIgnoreRuleFlags() []components.Flag {
	return []components.Flag{
		common.GetServerIdFlag(),
	}
}

func getCreateIgnoreRuleArguments() []components.Argument {
	return []components.Argument{
		{
			Name:        "json-file",
			Description: "Path to a JSON file containing the ignore rule payload (request body).",
		},
	}
}

func createIgnoreRuleCmd(c *components.Context) error {
	if len(c.Arguments) < 1 {
		return errors.New("json-file path is required")
	}
	jsonPath := c.Arguments[0]
	if jsonPath == "" {
		return errors.New("json-file path cannot be empty")
	}

	bodyBytes, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file %q: %w", jsonPath, err)
	}

	if err := validateIgnoreRuleJSON(bodyBytes); err != nil {
		return err
	}

	xrayManager, err := createXrayManager(c)
	if err != nil {
		return err
	}

	xrayDetails := xrayManager.Config().GetServiceDetails()
	httpDetails := xrayDetails.CreateHttpClientDetails()
	httpDetails.SetContentTypeApplicationJson()

	url := clientutils.AddTrailingSlashIfNeeded(xrayDetails.GetUrl()) + ignoreRulesAPI

	resp, body, err := xrayManager.Client().SendPost(url, bodyBytes, &httpDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusCreated); err != nil {
		return err
	}

	log.Info("Successfully created ignore rule.")
	log.Debug("Response: " + string(body))
	return nil
}

// validateIgnoreRuleJSON performs a lightweight pre-flight check on the JSON
// payload before sending it to the Xray API, giving users clearer errors for
// common mistakes (missing notes, empty filters).
func validateIgnoreRuleJSON(body []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(body, &m); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	if _, ok := m["notes"]; !ok {
		return errors.New("JSON must contain a \"notes\" field (required by Xray API)")
	}
	filters, ok := m["ignore_filters"].(map[string]interface{})
	if !ok || filters == nil {
		return errors.New("JSON must contain an \"ignore_filters\" object (required by Xray API)")
	}
	for _, v := range filters {
		if v != nil {
			return nil
		}
	}
	return errors.New("ignore_filters must contain at least one criterion (e.g. cves, vulnerabilities, licenses, components)")
}
