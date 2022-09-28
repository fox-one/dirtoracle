/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"encoding/json"

	"github.com/fox-one/dirtoracle/config"
	"github.com/go-yaml/yaml"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "init config",
	Run: func(cmd *cobra.Command, args []string) {
		buf := &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(config.Config{}); err != nil {
			cmd.PrintErrln(err)
			return
		}

		v := make(map[string]interface{})
		if err := json.NewDecoder(buf).Decode(&v); err != nil {
			cmd.PrintErrln(err)
			return
		}

		if err := yaml.NewEncoder(cmd.OutOrStdout()).Encode(v); err != nil {
			cmd.PrintErrln(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
