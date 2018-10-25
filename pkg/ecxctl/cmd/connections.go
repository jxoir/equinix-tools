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

// connectionsCmd represents the metros command
var connectionsCmd = &cobra.Command{
	Use:   "connections",
	Short: "Operations related to ECX connections (buyer)",
}

var connectionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "list all buyer connections",
	Run:   connectionsListCommand,
}

var connectionsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get specific connection by uuid",
	Args:  cobra.MinimumNArgs(1),
	Run:   connectionsGetByUUIDCommand,
}

/**
var connectionsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create L2 connection (virtual circuit) in Azure, AWS or other cloud services",
	Run:   connectionsCreateCommand,
}**/

func init() {
	rootCmd.AddCommand(connectionsCmd)
	connectionsCmd.AddCommand(connectionsListCmd)
	connectionsCmd.AddCommand(connectionsGetCmd)
}

func connectionsListCommand(cmd *cobra.Command, args []string) {
	connList, err := ConnectionsAPIClient.GetAllBuyerConnections()

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Total connections: %v\n", connList.TotalCount)

		for _, connection := range connList.Items {
			connRes, _ := json.MarshalIndent(connection, "", "    ")
			fmt.Println(string(connRes))
		}
	}
}

func connectionsGetByUUIDCommand(cmd *cobra.Command, args []string) {
	for _, uuid := range args {
		if globalFlags.Debug {
			log.Println("Get connection:" + uuid)
		}
		conn, err := ConnectionsAPIClient.GetByUUID(uuid)
		if err != nil {
			log.Fatal(err)
		} else {
			connRes, _ := json.MarshalIndent(conn, "", "    ")
			fmt.Println(string(connRes))
		}
	}
}

/**
func connectionsCreateCommand(cmd *cobra.Command, args []string) {

}**/
