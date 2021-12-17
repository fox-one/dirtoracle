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
	"encoding/json"
	"fmt"
	"net/http"

	scanner "github.com/fox-one/dirtoracle/apps/mvm/handler/pricescanner"
	"github.com/fox-one/dirtoracle/handler/hc"
	"github.com/fox-one/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run dirtoracle example server",
	Run: func(cmd *cobra.Command, args []string) {
		database := provideDatabase()
		defer database.Close()

		system := provideSystem()
		assets := provideAssetStore(database)

		mux := chi.NewMux()
		mux.Use(middleware.Recoverer)
		mux.Use(middleware.StripSlashes)
		mux.Use(cors.AllowAll().Handler)
		mux.Use(logger.WithRequestID)
		mux.Use(middleware.Logger)

		// hc
		mux.Mount("/price-requests", scanner.Handle(
			*system,
			assets,
		))
		mux.Mount("/hc", hc.Handle(rootCmd.Version))

		// launch server
		port, _ := cmd.Flags().GetInt("port")
		addr := fmt.Sprintf(":%d", port)

		http.ListenAndServe(addr, mux)
	},
}

func JSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	_ = enc.Encode(v)
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().Int("port", 9245, "worker api port")
}
