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
	"strings"

	"github.com/spf13/cobra"
)

var sellerProfileMetro string
var sellerProfileUUID string

// metrosCmd represents the metros command
var sellerCmd = &cobra.Command{
	Use:   "seller",
	Short: "ECX seller L2/L3 operations",
}

var sellerL2Cmd = &cobra.Command{
	Use:   "l2",
	Short: "ECX seller L2 operations",
}

var sellerListCmd = &cobra.Command{
	Use:   "list",
	Short: "list all seller profiles for given metros",
	Run:   sellerListCommand,
}

var sellerGetCmd = &cobra.Command{
	Use:   "get",
	Short: "fetch seller profile by uuid",
	Run:   sellerGetByUUIDCommand,
}

var sellerL3Cmd = &cobra.Command{
	Use:   "l3",
	Short: "ECX seller L3 services operations",
}

var sellerServicesListCmd = &cobra.Command{
	Use:   "list",
	Short: "list all seller services (L3) for given metros",
	Run:   sellerServicesListCommand,
}

func init() {
	rootCmd.AddCommand(sellerCmd)

	// Group L2 commands
	sellerCmd.AddCommand(sellerL2Cmd)
	sellerL2Cmd.AddCommand(sellerListCmd)
	sellerL2Cmd.AddCommand(sellerGetCmd)
	sellerListCmd.Flags().StringVarP(&sellerProfileMetro, "metros", "", "", "comma separated list of metro codes")
	sellerGetCmd.Flags().StringVarP(&sellerProfileUUID, "uuid", "", "", "seller profile uuid to fetch")

	// Group L3 commands
	sellerCmd.AddCommand(sellerL3Cmd)
	sellerL3Cmd.AddCommand(sellerServicesListCmd)
	sellerServicesListCmd.Flags().StringVarP(&sellerProfileMetro, "metros", "", "", "comma separated list of metro codes")

}

func sellerListCommand(cmd *cobra.Command, args []string) {

	metros := strings.Split(sellerProfileMetro, ",")

	sellerList, err := SellerServicesAPIClient.GetAllL2SellerProfiles(&metros)
	if err != nil {
		log.Fatal(err)
	}

	if sellerList != nil && sellerList.TotalCount > 0 {
		fmt.Printf("Total service profiles (L2): %v\n", sellerList.TotalCount)

		for _, sprofile := range sellerList.Items {
			if strings.Contains(sprofile.Name, "AWS") {
				sellerRes, _ := json.MarshalIndent(sprofile, "", "    ")
				fmt.Println(string(sellerRes))
			}
		}
	} else if sellerList != nil && sellerList.TotalCount == 0 {
		fmt.Println("There are no seller profiles for specified metro")
	}

}

func sellerGetByUUIDCommand(cmd *cobra.Command, args []string) {
	if sellerProfileUUID == "" {
		log.Fatal("seller profile UUID required")
	}

	sellerProfile, err := SellerServicesAPIClient.GetSellerProfileByUUID(sellerProfileUUID)

	if err != nil {
		log.Fatal(err)
	}

	sellerRes, _ := json.MarshalIndent(sellerProfile.Payload, "", "    ")
	fmt.Println(string(sellerRes))
}

func sellerServicesListCommand(cmd *cobra.Command, args []string) {
	metros := strings.Split(sellerProfileMetro, ",")

	sellerList, err := SellerServicesAPIClient.GetAllL3SellerServices(&metros)
	if err != nil {
		log.Fatal(err)
	}

	if sellerList != nil && sellerList.TotalCount > 0 {
		fmt.Printf("Total seller services profiles (L3): %v\n", sellerList.TotalCount)

		for _, sprofile := range sellerList.Items {
			sellerRes, _ := json.MarshalIndent(sprofile, "", "    ")
			fmt.Println(string(sellerRes))
		}
	} else if sellerList != nil && sellerList.TotalCount == 0 {
		fmt.Println("There are no seller services profiles for specified metro")
	}
}
