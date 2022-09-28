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
	"fmt"
	"net/http"
	"time"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/dirtoracle/exchanges"
	"github.com/fox-one/dirtoracle/exchanges/bwatch"
	"github.com/fox-one/dirtoracle/exchanges/fswap"
	"github.com/fox-one/dirtoracle/handler/hc"
	"github.com/fox-one/dirtoracle/worker"
	"github.com/fox-one/dirtoracle/worker/cashier"
	"github.com/fox-one/dirtoracle/worker/oracle"
	"github.com/fox-one/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "run dirtoracle worker",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		cfg.DB.ReadHost = ""
		database, err := provideDatabase()
		if err != nil {
			cmd.PrintErrf("provideDatabase failed: %v", err)
			return
		}

		defer database.Close()
		client := provideMixinClient()

		subscribers := provideSubscriberStore(database)
		wallets := provideWalletStore(database)
		walletz := provideWalletService(client)
		system := provideSystem()
		assetz := provideAssetService(client)

		exs, _ := cmd.Flags().GetStringArray("exchanges")
		ex := loadExchange(assetz, exs)

		workers := []worker.Worker{
			oracle.New(ex, client, wallets, subscribers, system),
			cashier.New(wallets, walletz),
		}

		// worker api
		{
			mux := chi.NewMux()
			mux.Use(middleware.Recoverer)
			mux.Use(middleware.StripSlashes)
			mux.Use(cors.AllowAll().Handler)
			mux.Use(logger.WithRequestID)
			mux.Use(middleware.Logger)

			// hc
			{
				mux.Mount("/hc", hc.Handle(rootCmd.Version))
			}

			// launch server
			port, _ := cmd.Flags().GetInt("port")
			addr := fmt.Sprintf(":%d", port)

			go http.ListenAndServe(addr, mux)
		}

		cmd.Printf("dirtoracle worker with version %q launched!\n", rootCmd.Version)

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

func loadExchange(assetz core.AssetService, exchNames []string) core.Exchange {
	var exch core.Exchange
	exchs := provideAllExchanges(assetz)
	for _, exchName := range exchNames {
		if e, ok := exchs[exchName]; ok {
			if exch == nil {
				exch = e
			} else {
				exch = exchanges.Chain(exch, e)
			}
		}
	}

	exch = exchanges.Cache(exch, time.Minute)
	exch = bwatch.New(exch)
	exch = exchanges.Humanize(fswap.Lp(exch))
	exch = exchanges.Cache(exch, time.Minute)
	return exch
}

func init() {
	rootCmd.AddCommand(workerCmd)
	workerCmd.Flags().Int("port", 9245, "worker api port")
	workerCmd.Flags().StringArray("exchanges", nil, "exchanges")
}
