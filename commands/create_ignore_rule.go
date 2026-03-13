package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/jfrog/jfrog-cli-core/v2/plugins/common"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const ignoreRulesPath = "/xray/api/v1/ignore_rules"

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

	// Validate JSON and required keys
	if err := validateIgnoreRuleJSON(bodyBytes); err != nil {
		return err
	}

	baseURL, authHeader, err := getBaseURLAndToken(c)
	if err != nil {
		return err
	}

	url := baseURL + ignoreRulesPath
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	log.Info("Successfully created ignore rule.")
	return nil
}

// validateIgnoreRuleJSON checks that the body is valid JSON and matches the Xray Create Ignore Rule schema.
// Request body must have top-level "notes" and "ignore_filters" (all filter criteria inside ignore_filters).
// See https://jfrog.com/help/r/xray-rest-apis/create-ignore-rule
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
	filterKeys := []string{"vulnerabilities", "cves", "licenses", "components", "builds", "artifacts", "docker_layers", "release_bundles", "operational_risk", "exposures", "watches", "policies"}
	hasFilter := false
	for _, k := range filterKeys {
		if v, ok := filters[k]; ok && v != nil {
			hasFilter = true
			break
		}
	}
	if !hasFilter {
		return errors.New("ignore_filters must contain at least one criterion (e.g. cves, vulnerabilities, licenses, components, builds, artifacts, docker_layers)")
	}
	return nil
}
