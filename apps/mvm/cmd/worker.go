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
	"github.com/fox-one/dirtoracle/apps/mvm/worker/payee"
	"github.com/fox-one/dirtoracle/worker"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use: "worker",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		cfg.DB.ReadHost = ""
		database := provideDatabase()
		defer database.Close()

		client := provideMixinClient()
		assets := provideAssetStore(database)
		walletz := provideWalletService(client)
		system := provideSystem()

		cmd.Printf("worker with version %q launched!\n", rootCmd.Version)

		workers := []worker.Worker{
			payee.New(*system, assets, walletz),
		}

		g, ctx := errgroup.WithContext(ctx)
		for idx := range workers {
			w := workers[idx]
			g.Go(func() error {
				return w.Run(ctx)
			})
		}

		if err := g.Wait(); err != nil {
			cmd.PrintErrln("run worker", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
	workerCmd.Flags().IntP("port", "p", 8123, "server port")
}
