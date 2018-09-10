// Copyright © 2018 Juan Manuel Irigaray <jirigaray@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var VERSION = "0.1"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ecxctl",
	Short: "Equinix *UNOFFICIAL* ECX and ECP CLI",
	Long: `Copyright © 2018 Juan Manuel Irigaray - Licensed under the Apache License, Version 2.0
An UNOFFICIAL GO CLI for ECX and ECP Tested with Go 1.10+
	
WARNING: This CLI is NOT official,

ecxctl is a CLI for Equinix ECX Fabric within equinix-tools toolchain.`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show ecxctl client version",
	Long:  `Everything should be versioned :)`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(rootCmd.Use + " " + VERSION)
	},
}

// EcxAPIClient instance of EquinixAPI to ECX
var EcxAPIClient *EquinixAPIClient

// ConnectionsAPIClient Connections interface
var ConnectionsAPIClient *ECXConnectionsAPI

// MetrosAPIClient Metros interface
var MetrosAPIClient *ECXMetrosAPI

// PortsAPIClient Ports interface
var PortsAPIClient *ECXPortsAPI

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	globalFlags = GlobalFlags{}
)

func init() {
	cobra.OnInitialize(initConfig)

	// Define base commands
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ecxctl.yaml)")
	rootCmd.PersistentFlags().BoolVar(&globalFlags.Debug, "debug", false, "enable client-side debug logging")

	rootCmd.PersistentFlags().StringVar(&globalFlags.PlaygroundToken, "playground-token", "", "Equinix Developer Playground Token (will disable API/SECRET authentication)")
	rootCmd.PersistentFlags().StringVar(&globalFlags.PlaygroundAPIEndpoint, "playground-endpoint", "playgroundapi.equinix.com", "Equinix Developer Playground endpoint")

	rootCmd.PersistentFlags().BoolVar(&globalFlags.NoSSL, "ignore-ssl", false, "Don't verify server SSL **INSECURE**")

	rootCmd.PersistentFlags().StringVar(&globalFlags.EcxAPIHost, "ecx-api-host", os.Getenv("ECX_API_HOST"), "ECX API endpoint")
	rootCmd.PersistentFlags().StringVar(&globalFlags.UserName, "user", os.Getenv("ECX_API_USER"), "portal username")
	rootCmd.PersistentFlags().StringVar(&globalFlags.UserPassword, "password", os.Getenv("ECX_API_USER_PASSWORD"), "portal password")
	rootCmd.PersistentFlags().StringVar(&globalFlags.APICredentialGrantType, "client-grant-type", "client_credentials", "specify api client_grant type (default:client_credentials)")

	// User needs to create an API ID and Secret at Equinix developer portal https://developer.equinix.com
	rootCmd.PersistentFlags().StringVar(&globalFlags.EquinixAPIId, "equinix-api-id", os.Getenv("EQUINIX_API_ID"), "Equinix API Application ID")
	rootCmd.PersistentFlags().StringVar(&globalFlags.EquinixAPISecret, "equinix-api-secret", os.Getenv("EQUINIX_API_SECRET"), "Equinix API Application Secret")

	// Bind viper config flags
	viper.BindPFlag("ecx-api-host", rootCmd.PersistentFlags().Lookup("ecx-api-host"))
	viper.BindPFlag("equinix-api-id", rootCmd.PersistentFlags().Lookup("equinix-api-id"))
	viper.BindPFlag("equinix-api-secret", rootCmd.PersistentFlags().Lookup("equinix-api-secret"))
	viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.AddCommand(versionCmd)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".ecxctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ecxctl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	initAPIClient()

}

func initAPIClient() {
	// Setup api client and params

	// Initialize Equinix Client
	if EcxAPIClient == nil {
		if globalFlags.PlaygroundToken != "" {
			// Set-up playground environment
			globalFlags.EcxAPIHost = globalFlags.PlaygroundAPIEndpoint
		}
		clientParams := &EquinixAPIParams{
			AppID:           globalFlags.EquinixAPIId,
			AppSecret:       globalFlags.EquinixAPISecret,
			GrantType:       "client_credentials",
			UserName:        globalFlags.UserName,
			UserPassword:    globalFlags.UserPassword,
			Endpoint:        globalFlags.EcxAPIHost,
			PlaygroundToken: globalFlags.PlaygroundToken,
		}

		EcxAPIClient = NewEcxAPIClient(clientParams, globalFlags.EcxAPIHost, globalFlags.NoSSL)
		ConnectionsAPIClient = NewECXConnectionsAPI(EcxAPIClient)
		MetrosAPIClient = NewECXMetrosAPI(EcxAPIClient)
		PortsAPIClient = NewECXPortsAPI(EcxAPIClient)
	}
}
