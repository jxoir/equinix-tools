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

	apiconnections "github.com/jxoir/go-ecxfabric/client/connections"

	"github.com/spf13/cobra"
)

type ConnectionsAPIHandler interface {
	GetByUUID(uuid string) (*apiconnections.GetConnectionByUUIDUsingGETOK, error)
	GetAllBuyerConnections() (*apiconnections.GetAllBuyerConnectionsUsingGETOK, error)
}

type ECXConnectionsAPI struct {
	*EquinixAPIClient
}

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

func init() {
	rootCmd.AddCommand(connectionsCmd)
	connectionsCmd.AddCommand(connectionsListCmd)
	connectionsCmd.AddCommand(connectionsGetCmd)
}

// NewECXConnectionsAPI returns instantiated ECXConnectionsAPI struct
func NewECXConnectionsAPI(equinixAPIClient *EquinixAPIClient) *ECXConnectionsAPI {
	return &ECXConnectionsAPI{equinixAPIClient}
}

func connectionsListCommand(cmd *cobra.Command, args []string) {
	connList, err := ConnectionsAPIClient.GetAllBuyerConnections()
	fmt.Println(connList)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Total connections: %v", connList.Payload.TotalCount)
		for _, connection := range connList.Payload.Content {
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

// GetAllBuyerConnections returns array of GetAllBuyerConnectionsUsingGETOK with list of customer connections
func (m *ECXConnectionsAPI) GetAllBuyerConnections() (*apiconnections.GetAllBuyerConnectionsUsingGETOK, error) {
	token, err := m.GetToken()
	if err != nil {
		log.Fatal(err)
	}
	connectionsOK, _, err := m.Client.Connections.GetAllBuyerConnectionsUsingGET(nil, token)
	if err != nil {
		switch t := err.(type) {
		default:
			log.Fatal(err)
		case *apiconnections.GetAllBuyerConnectionsUsingGETBadRequest:
			for _, getconnerrors := range t.Payload {
				fmt.Println(getconnerrors.ErrorMessage)
				return nil, err
			}
		}
	}

	return connectionsOK, nil

}

// GetByUUID get connection by uuid
func (m *ECXConnectionsAPI) GetByUUID(uuid string) (*apiconnections.GetConnectionByUUIDUsingGETOK, error) {
	params := apiconnections.NewGetConnectionByUUIDUsingGETParams()
	params.ConnID = uuid

	token, err := m.GetToken()
	if err != nil {
		log.Fatal(err)
	}

	connectionOK, _, err := m.Client.Connections.GetConnectionByUUIDUsingGET(params, token)
	if err != nil {
		switch t := err.(type) {
		default:
			return nil, err
		case *apiconnections.GetConnectionByUUIDUsingGETBadRequest:
			for _, getconnerrors := range t.Payload {
				fmt.Println(getconnerrors.ErrorMessage + ":" + uuid)
				return nil, err
			}
		}
	}

	if connectionOK != nil {
		return connectionOK, nil
	}

	return nil, err

}
