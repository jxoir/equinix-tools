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

	"github.com/jxoir/go-ecxfabric/client/connections"

	"github.com/spf13/cobra"
)

// metrosCmd represents the metros command
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

func init() {
	rootCmd.AddCommand(connectionsCmd)
	connectionsCmd.AddCommand(connectionsListCmd)
	connectionsCmd.AddCommand(connectionsGetCmd)

}

func connectionsListCommand(cmd *cobra.Command, args []string) {

	if globalFlags.Debug {
		log.Println("Listing connections...")
	}

	connectionsOK, _, err := EcxAPIClient.Client.Connections.GetAllBuyerConnectionsUsingGET(nil, EcxAPIClient.apiToken)
	if err != nil {
		switch t := err.(type) {
		default:
			log.Fatal(err)
		case *connections.GetAllBuyerConnectionsUsingGETBadRequest:
			for _, getconnerrors := range t.Payload {
				fmt.Println(getconnerrors.ErrorMessage)
			}
		}
	}

	if err != nil {
		log.Fatal(err)
	}
	if connectionsOK != nil {
		fmt.Printf("Total connections: %v", connectionsOK.Payload.TotalCount)
		for _, connection := range connectionsOK.Payload.Content {
			connRes, _ := json.MarshalIndent(connection, "", "    ")
			fmt.Println(string(connRes))
		}
	}
}

func connectionsGetByUUIDCommand(cmd *cobra.Command, args []string) {
	params := connections.NewGetConnectionByUUIDUsingGETParams()

	for _, uuid := range args {
		if globalFlags.Debug {
			log.Println("Get connection:" + uuid)
		}

		params.ConnID = uuid

		connectionOK, connectionNC, err := EcxAPIClient.Client.Connections.GetConnectionByUUIDUsingGET(params, EcxAPIClient.apiToken)
		if err != nil {
			switch t := err.(type) {
			default:
				log.Fatal(err)
			case *connections.GetConnectionByUUIDUsingGETBadRequest:
				for _, getconnerrors := range t.Payload {
					fmt.Println(getconnerrors.ErrorMessage + ":" + uuid)
				}
			}
		}

		if connectionOK != nil {
			// for now just print the Marshal version of the payload but we need to create an output transformation
			// TODO: custom formatter
			connRes, _ := json.MarshalIndent(connectionOK, "", "    ")
			fmt.Println(string(connRes))
		} else if connectionNC != nil {
			fmt.Println("Connection " + uuid + " not found")
		}
	}
}
