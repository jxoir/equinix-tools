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
	metrosList, err := MetrosAPIClient.GetAllMetros()
	if err != nil {
		log.Fatal(err)
	} else {
		if metrosList != nil {
			metros := metrosList.Payload
			metrosRes, err := json.MarshalIndent(metros, "", "    ")
			if err != nil {
				log.Fatal("There was an error with json response:", err)
			} else {
				fmt.Println(string(metrosRes))
			}

		}
	}
}
