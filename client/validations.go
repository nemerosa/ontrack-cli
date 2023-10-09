package client

import (
	"ontrack-cli/config"
)

func ValidateWithTests(
	cfg *config.Config,
	project string,
	branch string,
	build string,
	validation string,
	description string,
	runInfo *RunInfo,
	passed int,
	skipped int,
	failed int,
) error {

	// Mutation payload
	var payload struct {
		ValidateBuildWithTests struct {
			Errors []struct {
				Message string
			}
		}
	}

	// Runs the mutation
	if err := GraphQLCall(cfg, `
			mutation ValidateBuildWithTests(
				$project: String!,
				$branch: String!,
				$build: String!,
				$validationStamp: String!,
				$description: String!,
				$runInfo: RunInfoInput,
				$passed: Int!,
				$skipped: Int!,
				$failed: Int!
			) {
				validateBuildWithTests(input: {
					project: $project,
					branch: $branch,
					build: $build,
					validation: $validationStamp,
					description: $description,
					runInfo: $runInfo,
					passed: $passed,
					skipped: $skipped,
					failed: $failed
				}) {
					errors {
						message
					}
				}
			}
		`, map[string]interface{}{
		"project":         project,
		"branch":          branch,
		"build":           build,
		"validationStamp": validation,
		"description":     description,
		"runInfo":         runInfo,
		"passed":          passed,
		"skipped":         skipped,
		"failed":          failed,
	}, &payload); err != nil {
		return err
	}

	// Checks for errors
	if err := CheckDataErrors(payload.ValidateBuildWithTests.Errors); err != nil {
		return err
	}

	// OK
	return nil
}
