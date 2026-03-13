package commands

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	clicommands "github.com/jfrog/jfrog-cli-core/v2/common/commands"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/common"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const (
	licensesPrioritiesPath = "/xray/api/v1/licensesNames/priorities"
)

func GetSetLicensePriorityCommand() components.Command {
	return components.Command{
		Name:        "set-license-priority",
		Description: "Set the priority for a license name via the Xray API.",
		Aliases:     []string{"slp"},
		Arguments:   getSetLicensePriorityArguments(),
		Flags:       getSetLicensePriorityFlags(),
		Action: func(c *components.Context) error {
			return setLicensePriorityCmd(c)
		},
	}
}

func getSetLicensePriorityFlags() []components.Flag {
	return []components.Flag{
		common.GetServerIdFlag(),
	}
}

func getSetLicensePriorityArguments() []components.Argument {
	return []components.Argument{
		{
			Name:        "license-name",
			Description: "The name of the license to set priority for.",
		},
		{
			Name:        "license-priority",
			Description: "The priority value to set for the license (e.g. 1, 2, 3).",
		},
	}
}

type licensePriorityRequest struct {
	LicenseName     string `json:"key"`
	LicensePriority int    `json:"priority"`
}

// getBaseURLAndToken returns the Xray/platform base URL and Authorization header value for API calls.
// It first tries the JFrog CLI configuration (from "jf c add"), then falls back to
// JFROG_URL/JF_URL and JFROG_ACCESS_TOKEN/JF_ACCESS_TOKEN environment variables.
// For config: if AccessToken is set, returns "Bearer "+token; if only User/Password are set, returns "Basic "+base64(user:password).
func getBaseURLAndToken(c *components.Context) (baseURL, authHeader string, err error) {
	serverID := c.GetStringFlagValue("server-id")

	// Try JFrog CLI config first
	details, configErr := clicommands.GetConfig(serverID, false)
	if configErr == nil && details != nil {
		if err = config.CreateInitialRefreshableTokensIfNeeded(details); err != nil {
			return "", "", fmt.Errorf("failed to refresh tokens from config: %w", err)
		}
		baseURL = details.GetUrl()
		baseURL = strings.TrimSuffix(baseURL, "/")
		if baseURL == "" {
			goto fallback
		}
		if details.GetAccessToken() != "" {
			return baseURL, "Bearer " + details.GetAccessToken(), nil
		}
		if details.GetUser() != "" && details.GetPassword() != "" {
			basic := base64.StdEncoding.EncodeToString([]byte(details.GetUser() + ":" + details.GetPassword()))
			return baseURL, "Basic " + basic, nil
		}
	}
fallback:

	// Fall back to environment variables
	baseURL = os.Getenv("JFROG_URL")
	if baseURL == "" {
		baseURL = os.Getenv("JF_URL")
	}
	if baseURL == "" {
		return "", "", errors.New(
			"no JFrog URL found. Configure the CLI with 'jf c add' and pass --server-id, or set JFROG_URL (or JF_URL) to the JFrog platform URL (e.g. https://your-instance.jfrog.io)")
	}
	baseURL = strings.TrimSuffix(baseURL, "/")

	token := os.Getenv("JFROG_ACCESS_TOKEN")
	if token == "" {
		token = os.Getenv("JF_ACCESS_TOKEN")
	}
	if token == "" {
		return "", "", errors.New(
			"no access token found. Configure the CLI with 'jf c add' and pass --server-id, or set JFROG_ACCESS_TOKEN (or JF_ACCESS_TOKEN) for authentication")
	}
	return baseURL, "Bearer " + token, nil
}

func setLicensePriorityCmd(c *components.Context) error {
	if len(c.Arguments) < 2 {
		return errors.New("both license-name and license-priority are required")
	}
	licenseName := c.Arguments[0]
	licensePriorityStr := c.Arguments[1]

	if licenseName == "" {
		return errors.New("license-name cannot be empty")
	}
	if licensePriorityStr == "" {
		return errors.New("license-priority cannot be empty")
	}

	licensePriority, err := strconv.Atoi(licensePriorityStr)
	if err != nil {
		return fmt.Errorf("license-priority must be a valid integer: %w", err)
	}

	baseURL, authHeader, err := getBaseURLAndToken(c)
	if err != nil {
		return err
	}

	reqBody := licensePriorityRequest{
		LicenseName:     licenseName,
		LicensePriority: licensePriority,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to build request body: %w", err)
	}

	url := baseURL + licensesPrioritiesPath
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)
	fmt.Println("Request: ", req)
	fmt.Println("Body: ", string(bodyBytes))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("Response: ", string(respBody))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	log.Info(fmt.Sprintf("Successfully set priority for license '%s' to %d.", licenseName, licensePriority))
	return nil
}
