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
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/pandodao/blst"
	"github.com/pandodao/blst/en256"
	"github.com/spf13/cobra"
)

// keysCmd represents the keys command
var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "generate key pairs",
	Run: func(cmd *cobra.Command, args []string) {
		alg, _ := cmd.Flags().GetString("alg")

		switch alg {
		case "blst":
			private := blst.GenerateKey()
			public := private.PublicKey()
			cmd.Println("Private", private)
			cmd.Println("Public ", public)

		case "en256":
			private := en256.GenerateKey()
			public := private.PublicKey()
			cmd.Println("Private", private)

			bts, _ := public.Bytes()
			cmd.Println("Public ", hex.EncodeToString(bts))

		case "ed25519", "ed":
			private := mixin.GenerateEd25519Key()
			public := private.Public().(ed25519.PublicKey)
			cmd.Println("Private", base64.StdEncoding.EncodeToString(private))
			cmd.Println("Public ", base64.StdEncoding.EncodeToString(public))

		default:
			cmd.PrintErr("unknown algorithm: ", alg)
		}

	},
}

func init() {
	rootCmd.AddCommand(keysCmd)

	keysCmd.Flags().String("alg", "blst", "alg")
}
