/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"github.com/fox-one/dirtoracle/apps/mvm/core"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/spf13/cobra"
)

// assetsCmd represents the assets command
var assetsCmd = &cobra.Command{
	Use:   "assets",
	Short: "load top assets & chainss to database",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		cfg.DB.ReadHost = ""
		database := provideDatabase()
		defer database.Close()

		assets := provideAssetStore(database)

		var (
			allAssets []*core.Asset
			assetM    = map[string]bool{}
		)

		{
			resp, err := mixin.Request(ctx).Get("/network/chains")
			if err != nil {
				cmd.PrintErr(err)
				return
			}

			var assets []*mixin.Asset
			if err := mixin.UnmarshalResponse(resp, &assets); err != nil {
				cmd.PrintErr(err)
				return
			}
			for _, asset := range assets {
				if _, ok := assetM[asset.AssetID]; !ok {
					allAssets = append(allAssets, &core.Asset{
						AssetID:       asset.ChainID,
						Symbol:        asset.Symbol,
						PriceDuration: 600,
					})
					assetM[asset.AssetID] = true
				}
			}
		}
		{
			assets, err := mixin.ReadTopNetworkAssets(ctx)
			if err != nil {
				cmd.PrintErr(err)
				return
			}
			for _, asset := range assets {
				if _, ok := assetM[asset.AssetID]; !ok {
					allAssets = append(allAssets, &core.Asset{
						AssetID:       asset.AssetID,
						Symbol:        asset.Symbol,
						PriceDuration: 600,
					})
					assetM[asset.AssetID] = true
				}
			}
		}

		if err := assets.Create(ctx, allAssets); err != nil {
			cmd.PrintErr(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(assetsCmd)
}
