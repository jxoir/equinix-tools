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
	"os"

	apiconnections "github.com/jxoir/go-ecxfabric/buyer/client/connections"
	"github.com/spf13/cobra"
)

var filterValues string
var deleteUUID string
var connectionMetro string

// vars for create connection command
var createL2SellerConPrimaryName string
var createL2SellerConSecondaryName string
var createL2SellerConPrimaryPortUUID string
var createL2SellerConSecondaryPortUUID string
var createL2SellerConPrimaryVlanSTag int64 // should be casted to int64
var createL2SellerConSecondaryVlanSTag int64
var createL2SellerConSellerProfileUUID string
var createL2SellerConSellerRegion string
var createL2SellerConSellerMetroCode string
var createL2SellerConSpeed int64 // should be casted to int64
var createL2SellerConSpeedUnit string
var createL2SellerConNotificationsEmail string
var createL2SellerConAuthorizationKey string
var createL2SellerConNamedTag string

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

var connectionsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create L2 connection (virtual circuit) in Azure, AWS or other seller services",
	Run:   connectionsCreateCommand,
}

func init() {
	rootCmd.AddCommand(connectionsCmd)
	connectionsCmd.AddCommand(connectionsListCmd)
	connectionsCmd.AddCommand(connectionsGetCmd)
	connectionsCmd.AddCommand(connectionsDeleteCmd)
	connectionsCmd.AddCommand(connectionsCreateCmd)

	connectionsListCmd.PersistentFlags().StringVarP(&filterValues, "filter", "f", "", "Comma separated key-value pair of filter (eg.: filter=Key=Name,Value=ECX)")
	connectionsListCmd.PersistentFlags().StringVarP(&connectionMetro, "metro", "", "", "Filter metro code (ex.: LD)")

	connectionsDeleteCmd.Flags().StringVarP(&deleteUUID, "uuid", "u", "", "*connection* to delete")
	connectionsDeleteCmd.MarkFlagRequired("uuid")

	connectionsCreateCmd.Flags().StringVarP(&createL2SellerConPrimaryName, "name", "n", "", "name for the new connection")
	connectionsCreateCmd.Flags().StringVarP(&createL2SellerConSecondaryName, "name-sec", "", "", "name for the secondary connection")
	connectionsCreateCmd.Flags().StringVarP(&createL2SellerConPrimaryPortUUID, "port-uuid", "", "", "user port uuid")
	connectionsCreateCmd.Flags().StringVarP(&createL2SellerConSecondaryPortUUID, "port-uuid-sec", "", "", "secondary user port uuid")
	connectionsCreateCmd.Flags().Int64VarP(&createL2SellerConPrimaryVlanSTag, "vlan", "", 0, "user vlan")
	connectionsCreateCmd.Flags().Int64VarP(&createL2SellerConSecondaryVlanSTag, "vlan-sec", "", 0, "secondary user vlan")
	connectionsCreateCmd.Flags().StringVarP(&createL2SellerConSellerProfileUUID, "seller-uuid", "", "", "seller profile uuid (destination UUID)")
	connectionsCreateCmd.Flags().StringVarP(&createL2SellerConSellerRegion, "seller-region", "", "", "seller destination region (ex. AWS: eu-west-1)")
	connectionsCreateCmd.Flags().StringVarP(&createL2SellerConSellerMetroCode, "seller-metro", "", "", "seller destination metro code (ex.: LD)")
	connectionsCreateCmd.Flags().Int64VarP(&createL2SellerConSpeed, "speed", "", 0, "connection speed (must be by 50 for MB ex.: 50, 100, 200, 500)")
	connectionsCreateCmd.Flags().StringVarP(&createL2SellerConSpeedUnit, "speed-unit", "", "", "connection speed unit MB, GB")
	connectionsCreateCmd.Flags().StringVarP(&createL2SellerConAuthorizationKey, "auth-key", "", "", "service authorization key (in AWS case use AWS Account ID)")
	connectionsCreateCmd.Flags().StringVarP(&createL2SellerConNotificationsEmail, "notifications-email", "", "", "email for notifications")
	connectionsCreateCmd.Flags().StringVarP(&createL2SellerConNamedTag, "named-tag", "", "", "primary ZSide Vlan CTag")

	connectionsCreateCmd.MarkFlagRequired("name")
	connectionsCreateCmd.MarkFlagRequired("port-uuid")
	connectionsCreateCmd.MarkFlagRequired("vlan")
	connectionsCreateCmd.MarkFlagRequired("seller-uuid")
	connectionsCreateCmd.MarkFlagRequired("seller-region")
	connectionsCreateCmd.MarkFlagRequired("seller-metro")
	connectionsCreateCmd.MarkFlagRequired("speed")
	connectionsCreateCmd.MarkFlagRequired("speed-unit")
	connectionsCreateCmd.MarkFlagRequired("seller-region")
	connectionsCreateCmd.MarkFlagRequired("auth-key")
	connectionsCreateCmd.MarkFlagRequired("notifications-email")
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

			connections := connList.GetItems()
			connRes, err := json.MarshalIndent(connections, "", "    ")
			if err != nil {
				log.Fatal("There was an error with json response:", err)
			} else {
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

func connectionsCreateCommand(cmd *cobra.Command, args []string) {
	// required params
	// sellerService UUID - uuid to connect to
	// speed 50 mb
	// name (connection)
	// authorizationKey (AWS accountid as an example)
	// vlanSTag - vlan source tag
	// notifications email

	params := ConnectionsAPIClient.NewCreateL2SellerConnectionParams()

	params.PrimaryName = createL2SellerConPrimaryName
	if createL2SellerConSecondaryName != "" {
		params.SecondaryName = createL2SellerConSecondaryName
	}

	params.PrimaryPortUUID = createL2SellerConPrimaryPortUUID
	if createL2SellerConSecondaryPortUUID != "" {
		params.SecondaryPortUUID = createL2SellerConSecondaryPortUUID
	}

	if createL2SellerConPrimaryVlanSTag == 0 {
		panic(createL2SellerConPrimaryVlanSTag)
	}
	params.PrimaryVlanSTag = createL2SellerConPrimaryVlanSTag
	if createL2SellerConSecondaryVlanSTag != 0 {
		params.SecondaryVlanSTag = createL2SellerConSecondaryVlanSTag
	}

	if createL2SellerConNamedTag != "" {
		params.NamedTag = createL2SellerConNamedTag
	}

	if createL2SellerConSpeed == 0 {
		panic(createL2SellerConSpeed)
	}

	params.Speed = createL2SellerConSpeed
	params.SpeedUnit = createL2SellerConSpeedUnit
	params.Notifications = []string{createL2SellerConNotificationsEmail}
	params.SellerRegion = createL2SellerConSellerRegion         //"eu-west-1" // get from seller? this should be AWS
	params.SellerMetroCode = createL2SellerConSellerMetroCode   // provided by customer
	params.AuthorizationKey = createL2SellerConAuthorizationKey // aws account id in this case

	params.ProfileUUID = createL2SellerConSellerProfileUUID

	conn, err := ConnectionsAPIClient.CreateL2ConnectionSellerProfile(params, SellerServicesAPIClient)
	if err != nil {
		switch t := err.(type) {
		case *apiconnections.CreateConnectionUsingPOSTBadRequest:
			for _, er := range t.Payload {
				fmt.Printf("Error %s with message %s\n", er.ErrorCode, er.ErrorMessage)
			}
		default:
			fmt.Printf("Error creating connection: %s\n", err.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Connection %s succesfully created\n", conn.Payload.PrimaryConnectionID)
}
