package client

import (
	"bytes"
	"text/template"

	config "ontrack-cli/config"
)

func SetupValidationStamp(
	cfg *config.Config,
	project string,
	branch string,
	validation string,
	description string,
	dataType string,
	dataTypeConfig string,
) error {

	tmpl, err := template.New("mutation").Parse(`
		mutation SetupValidationStamp(
			$project: String!,
			$branch: String!,
			$validation: String!,
			$description: String,
			$dataType: String
		) {
			setupValidationStamp(input: {
				project: $project,
				branch: $branch,
				validation: $validation,
				description: $description,
				dataType: $dataType,
				dataTypeConfig: {{ .DataTypeConfig }}
			}) {
				errors {
					message
				}
			}
		}
	`)
	if err != nil {
		return err
	}

	var tmplInput struct {
		DataTypeConfig string
	}

	if dataTypeConfig != "" {
		tmplInput.DataTypeConfig = dataTypeConfig
	} else {
		tmplInput.DataTypeConfig = "null"
	}

	var query bytes.Buffer
	if err := tmpl.Execute(&query, tmplInput); err != nil {
		return err
	}

	var data struct {
		SetupValidationStamp struct {
			Errors []struct {
				Message string
			}
		}
	}
	if err := GraphQLCall(cfg, query.String(), map[string]interface{}{
		"project":        project,
		"branch":         branch,
		"validation":     validation,
		"description":    description,
		"dataType":       dataType,
		"dataTypeConfig": dataTypeConfig,
	}, &data); err != nil {
		return err
	}

	// Checks for errors
	if err := CheckDataErrors(data.SetupValidationStamp.Errors); err != nil {
		return err
	}

	// OK
	return nil
}
