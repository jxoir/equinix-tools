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

	"github.com/jxoir/go-ecxfabric/client/ports"
	"github.com/spf13/cobra"
)

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

// List ports for the current user
func portsListCommand(cmd *cobra.Command, args []string) {

	if globalFlags.Debug {
		log.Println("Listing ports...")
	}

	portsOK, err := EcxAPIClient.Client.Ports.GetPortInfoUsingGET2(nil, EcxAPIClient.apiToken)
	if err != nil {
		switch t := err.(type) {
		default:
			log.Fatal(err)
		case *ports.GetPortInfoUsingGET2Forbidden:
			fmt.Println(t.Error())
		}
	}

	if err != nil {
		log.Fatal(err)
	}
	if portsOK != nil {
		for _, port := range portsOK.Payload {
			portRes, _ := json.MarshalIndent(port, "", "    ")
			fmt.Println(string(portRes))
		}
	}
}
