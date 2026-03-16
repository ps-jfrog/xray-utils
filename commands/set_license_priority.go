package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jfrog/jfrog-cli-core/v2/plugins/common"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const licensesPrioritiesAPI = "api/v1/licensesNames/priorities"

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

	xrayManager, err := createXrayManager(c)
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

	xrayDetails := xrayManager.Config().GetServiceDetails()
	httpDetails := xrayDetails.CreateHttpClientDetails()
	httpDetails.SetContentTypeApplicationJson()

	url := clientutils.AddTrailingSlashIfNeeded(xrayDetails.GetUrl()) + licensesPrioritiesAPI

	resp, body, err := xrayManager.Client().SendPost(url, bodyBytes, &httpDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return err
	}

	log.Info(fmt.Sprintf("Successfully set priority for license '%s' to %d.", licenseName, licensePriority))
	return nil
}
