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

var filterValues string
var deleteUUID string
var connectionMetro string

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

var connectionsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete specific connection by uuid",
	//Args:  cobra.MinimumNArgs(1),
	Run: connectionsDeleteByUUIDCommand,
}

/**
var connectionsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create L2 connection (virtual circuit) in Azure, AWS or other cloud services",
	Run:   connectionsCreateCommand,
}
**/
func init() {
	rootCmd.AddCommand(connectionsCmd)
	connectionsCmd.AddCommand(connectionsListCmd)
	connectionsCmd.AddCommand(connectionsGetCmd)
	connectionsCmd.AddCommand(connectionsDeleteCmd)

	connectionsListCmd.PersistentFlags().StringVarP(&filterValues, "filter", "f", "", "Comma separated key-value pair of filter (eg.: filter=Key=Name,Value=ECX)")
	connectionsListCmd.PersistentFlags().StringVarP(&connectionMetro, "metro", "", "", "Filter metro code (ex.: LD)")

	connectionsDeleteCmd.PersistentFlags().StringVarP(&deleteUUID, "uuid", "u", "", "UUID of the specific connection to delete")
	connectionsDeleteCmd.MarkFlagRequired("uuid")

}

func connectionsListCommand(cmd *cobra.Command, args []string) {

	metro := connectionMetro
	connList, err := ConnectionsAPIClient.GetAllBuyerConnections(&metro)

	if err != nil {
		log.Fatal(err)
	} else {
		if connList.Count() > 0 {
			if filterValues != "" {
				filter := parseFilteringAttributes(filterValues)
				connList.FilterItems(filter)

			}

			for _, connection := range connList.GetItems() {
				connRes, _ := json.MarshalIndent(connection, "", "    ")
				fmt.Println(string(connRes))
			}
		} else {
			if metro != "" {
				// error message changed, thanks Mischa :)
				fmt.Printf("No connections found for %s metro\n", metro)
			} else {
				fmt.Println("No connections found")
			}
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

func connectionsDeleteByUUIDCommand(cmd *cobra.Command, args []string) {

	if deleteUUID != "" {
		del, err := ConnectionsAPIClient.DeleteByUUID(deleteUUID)
		if err != nil {
			fmt.Printf("Error deleting connection: %s\n", deleteUUID)
		} else {
			fmt.Printf("Connection %s succesfully deleted\n", del.Payload.PrimaryConnectionID)
			//fmt.Println(del.Payload.Message)
		}
	} else {
		fmt.Println("Specify connection UUID to delete")
	}
}
