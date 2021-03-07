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

	// Parsing
	result := &graphResponse{
		Data: data,
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return err
	}

	// TODO Management of errors

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
