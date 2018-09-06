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

	apiports "github.com/jxoir/go-ecxfabric/client/ports"

	"github.com/spf13/cobra"
)

type PortsAPIHandler interface {
	GetAllMetros() (*apiports.GetPortInfoUsingGET2OK, error)
}

type ECXPortsAPI struct {
	*EquinixAPIClient
}

// metrosCmd represents the metros command
var portsCmd = &cobra.Command{
	Use:   "ports",
	Short: "Operations related to ECX ports",
}

var portsListCmd = &cobra.Command{
	Use:   "list",
	Short: "list all user virtual ports",
	Run:   portsListCommand,
}

func init() {
	rootCmd.AddCommand(portsCmd)
	portsCmd.AddCommand(portsListCmd)

}

// NewECXMetrosAPI returns instantiated ECXMetrosAPI struct
func NewECXPortsAPI(equinixAPIClient *EquinixAPIClient) *ECXPortsAPI {
	return &ECXPortsAPI{equinixAPIClient}
}

func portsListCommand(cmd *cobra.Command, args []string) {
	portsList, err := PortsAPIClient.GetAllPorts()
	if err != nil {
		log.Fatal(err)
	} else {
		if portsList != nil {
			for _, port := range portsList.Payload {
				portsRes, _ := json.MarshalIndent(port, "", "    ")
				fmt.Println(string(portsRes))
			}
		}
	}
}

// GetAllBuyerConnections returns array of GetAllBuyerConnectionsUsingGETOK with list of customer connections
func (m *ECXPortsAPI) GetAllPorts() (*apiports.GetPortInfoUsingGET2OK, error) {
	token, err := m.GetToken()
	if err != nil {
		log.Fatal(err)
	}
	respPortsOk, err := m.Client.Ports.GetPortInfoUsingGET2(nil, token)
	if err != nil {
		switch t := err.(type) {
		default:
			log.Fatal(err)
		case *apiports.GetPortInfoUsingGET2NotFound:
			if globalFlags.Debug {
				fmt.Println(t.Error())
			}
			return nil, err
		}
	}

	return respPortsOk, nil

}
