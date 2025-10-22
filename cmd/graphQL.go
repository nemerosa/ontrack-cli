/*
Copyright © 2021 Damien Coraboeuf <damien.coraboeuf@nemerosa.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"yontrack/client"
	config "yontrack/config"

	"github.com/spf13/cobra"
)

// graphQLCmd represents the graphQL command
var graphQLCmd = &cobra.Command{
	Use:   "graphql",
	Short: "Performs a GraphQL command",
	Long: `Performs a GraphQL command.
	
For example:

    yontrack graphql --query 'query ProjectList($name: String!) { projects(name: $name) { id name branches { name } } }' --var name=ontrack
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		query, err := cmd.Flags().GetString("query")
		if err != nil {
			return err
		}

		vars, err := cmd.Flags().GetStringSlice("var")
		if err != nil {
			return err
		}

		var variables = map[string]interface{}{}
		for _, token := range vars {
			name, value, err := parseVar(token)
			if err != nil {
				return err
			}
			variables[name] = value
		}

		cfg, err := config.GetSelectedConfiguration()
		if err != nil {
			return err
		}

		var data interface{}

		if err := client.GraphQLCall(cfg, query, variables, &data); err != nil {
			return err
		}

		res, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(res))

		// OK
		return nil
	},
}

func parseVar(token string) (string, string, error) {
	re := regexp.MustCompile(`^(.+)=(.*)$`)
	match := re.FindStringSubmatch(token)
	if match == nil {
		return "", "", errors.New("Variable " + token + " must match name=value")
	}

	name := match[1]
	value := match[2]

	return name, value, nil
}

func init() {
	rootCmd.AddCommand(graphQLCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// graphQLCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// graphQLCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	graphQLCmd.Flags().StringP("query", "q", "", "GraphQL query")
	graphQLCmd.Flags().StringSliceP("var", "v", []string{}, "GraphQL variable, in the form name=value")

	graphQLCmd.MarkFlagRequired("query")
}
