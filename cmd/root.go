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
	"os"
	"path"

	"github.com/fox-one/dirtoracle/config"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yiplee/structs"
)

var (
	cfgFile   string
	cfg       config.Config
	debugMode bool

	initialized bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dirtoracle",
	Short: "dirt oracle",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ver string) {
	rootCmd.Version = ver
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initLogging, initDone)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mtg.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "toggle debug mode")
}

func initConfig() {
	if initialized {
		return
	}

	if cfgFile == "" {
		dir, err := homedir.Dir()
		if err != nil {
			panic(err)
		}

		filename := path.Join(dir, ".dirtoracle.yaml")
		info, err := os.Stat(filename)
		if !os.IsNotExist(err) && !info.IsDir() {
			cfgFile = filename
		}
	}

	if cfgFile != "" {
		logrus.Debugln("use config file", cfgFile)
	}

	if err := config.Load(cfgFile, &cfg); err != nil {
		panic(err)
	}
}

func initLogging() {
	if initialized {
		return
	}

	if debugMode {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	structs.DefaultTagName = "json"
}

func initDone() {
	initialized = true
}
