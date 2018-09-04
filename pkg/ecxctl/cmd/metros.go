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
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// metrosCmd represents the metros command
var metrosCmd = &cobra.Command{
	Use:   "metros",
	Short: "Operations related to ECX Metros",
}

// metrosCmd represents the metros command
var metrosListCmd = &cobra.Command{
	Use:   "list",
	Short: "list available metros",
	Run:   metrosListCommand,
}

func init() {
	rootCmd.AddCommand(metrosCmd)
	metrosCmd.AddCommand(metrosListCmd)

}

func metrosListCommand(cmd *cobra.Command, args []string) {
	if globalFlags.Debug {
		log.Println("Listing metros...")
	}
	respMetrosOk, _, errMetro := EcxAPIClient.Client.Metros.GetMetrosUsingGET(nil, EcxAPIClient.apiToken)
	if errMetro != nil {
		log.Fatal(errMetro)
	}
	if respMetrosOk != nil {
		for _, metro := range respMetrosOk.Payload {
			connRes, _ := json.MarshalIndent(metro, "", "    ")
			fmt.Println(string(connRes))
		}
	}
}
