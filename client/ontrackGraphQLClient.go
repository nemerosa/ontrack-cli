package client

import (
	"encoding/json"
	"errors"
	"fmt"
	config "ontrack-cli/config"

	resty "github.com/go-resty/resty/v2"
)

// GraphQLCall performs a GraphQL query/mutation to Ontrack
func GraphQLCall(config *config.Config, query string, variables map[string]interface{}, data interface{}) error {
	body := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	client := resty.New()
	// client.SetDebug(true)
	if config.Token != "" {
		client.SetHeader("X-Ontrack-Token", config.Token)
	} else if config.Username != "" {
		client.SetBasicAuth(config.Username, config.Password)
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(config.URL + "/graphql")
	if err != nil {
		return err
	}

	// Error returned
	var error struct {
		Status  int
		Message string
	}
	if err := json.Unmarshal(resp.Body(), &error); err == nil {
		if error.Status != 0 {
			return fmt.Errorf("HTTP %d %s", error.Status, error.Message)
		}
	}

	// Parsing
	result := &graphResponse{
		Data: data,
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return err
	}

	// Management of errors
	if result.Errors != nil {
		if len(result.Errors) > 0 {
			var message string
			for index, error := range result.Errors {
				message += fmt.Sprintf("%d) %s\n", index+1, error.Message)
			}
			return errors.New(message)
		}
	}

	// OK
	return nil
}

type graphErr struct {
	Message string
}

type graphResponse struct {
	Data   interface{}
	Errors []graphErr
}

// CheckDataErrors Given a list of errors in a data GraphQL structure (typically
// returned by a mutation), returns a GoLang error aggregating all error messages
// or returns nil if there is no error.
func CheckDataErrors(errorsList []struct{ Message string }) error {
	if errorsList != nil {
		if len(errorsList) > 0 {
			var message string
			for index, error := range errorsList {
				message += fmt.Sprintf("%d) %s\n", index+1, error.Message)
			}
			return errors.New(message)
		}
	}
	// All good
	return nil
}
