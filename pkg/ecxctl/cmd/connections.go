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

	apiconnections "github.com/jxoir/go-ecxfabric/buyer/client/connections"
	"github.com/spf13/cobra"
)

var filterValues string
var deleteUUID string
var connectionMetro string

// flag to call wrapper around l2 connection
var createL2CSP bool

// vars for create connection command
var createL2ConAuthorizationKey string
var createL2ConNamedTag string
var createL2ConNotificationsEmail string
var createL2ConPrimaryName string

var createL2ConPrimaryPortUUID string
var createL2ConPrimaryVlanCTag string // should be casted to int64 - Inner tag
var createL2ConPrimaryVlanSTag int64  // should be casted to int64 - Outer tag

var createL2ConPrimaryZSidePortUUID string // should be casted to int64
var createL2ConPrimaryZSideVlanCTag string // should be casted to int64
var createL2ConPrimaryZSideVlanSTag int64  // should be casted to int64

var createL2ConSellerProfileUUID string
var purchaseOrderNumber string

var createL2ConSecondaryName string
var createL2ConSecondaryPortUUID string
var createL2ConSecondaryVlanCTag string
var createL2ConSecondaryVlanSTag int64

var createL2ConSecondaryZSidePortUUID string // should be casted to int64
var createL2ConSecondaryZSideVlanCTag int64  // should be casted to int64
var createL2ConSecondaryZSideVlanSTag int64  // should be casted to int64

var createL2ConSellerMetroCode string
var createL2ConSellerRegion string

var createL2ConSpeed int64 // should be casted to int64
var createL2ConSpeedUnit string

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

var connectionsCreateL2Cmd = &cobra.Command{
	Use:   "create",
	Short: "create L2 connection (virtual circuit) to any destination",
	Run:   connectionsCreateCommand,
}

func init() {
	rootCmd.AddCommand(connectionsCmd)
	connectionsCmd.AddCommand(connectionsListCmd)
	connectionsCmd.AddCommand(connectionsGetCmd)
	connectionsCmd.AddCommand(connectionsDeleteCmd)
	connectionsCmd.AddCommand(connectionsCreateL2Cmd)

	connectionsListCmd.PersistentFlags().StringVarP(&filterValues, "filter", "f", "", "Comma separated key-value pair of filter (eg.: filter=Key=Name,Value=ECX)")
	connectionsListCmd.PersistentFlags().StringVarP(&connectionMetro, "metro", "", "", "Filter metro code (ex.: LD)")

	connectionsDeleteCmd.Flags().StringVarP(&deleteUUID, "uuid", "u", "", "*connection* to delete")
	connectionsDeleteCmd.MarkFlagRequired("uuid")

	connectionsCreateL2Cmd.Flags().BoolVarP(&createL2CSP, "cloud", "c", false, "connect to a public cloud provider ex.: Azure, AWS, Google")
	connectionsCreateL2Cmd.Flags().StringVarP(&createL2ConAuthorizationKey, "auth-key", "", "", "service authorization key (in AWS case use AWS Account ID)")
	connectionsCreateL2Cmd.Flags().StringVarP(&createL2ConNotificationsEmail, "notifications-email", "", "", "email for notifications")
	connectionsCreateL2Cmd.Flags().StringVarP(&createL2ConNamedTag, "named-tag", "", "", "Private, Public, Microsoft, Manual (Microsoft requires special authorization, Manual forces stag)")

	connectionsCreateL2Cmd.Flags().StringVarP(&createL2ConPrimaryName, "name", "n", "", "name for the new connection")
	connectionsCreateL2Cmd.Flags().StringVarP(&createL2ConSecondaryName, "sec-name", "", "", "name for the secondary connection")

	connectionsCreateL2Cmd.Flags().StringVarP(&createL2ConPrimaryPortUUID, "port-uuid", "", "", "user port uuid")
	connectionsCreateL2Cmd.Flags().Int64VarP(&createL2ConPrimaryVlanSTag, "port-stag", "", 0, "S-Tag/Outer-tag of the primary port (vlan id for Dot1Q)")
	connectionsCreateL2Cmd.Flags().StringVarP(&createL2ConPrimaryVlanCTag, "port-ctag", "", "", "C-Tag/Inner-tag of the primary port (customer vlan id for QinQ)")

	connectionsCreateL2Cmd.Flags().StringVarP(&createL2ConPrimaryZSidePortUUID, "port-zside-uuid", "", "", "Z-side (remote) user port uuid")
	connectionsCreateL2Cmd.Flags().Int64VarP(&createL2ConPrimaryZSideVlanSTag, "port-zside-stag", "", 0, "Z-side (remote) S-Tag/Outer-tag of the primary port (vlan id for Dot1Q)")
	connectionsCreateL2Cmd.Flags().StringVarP(&createL2ConPrimaryZSideVlanCTag, "port-zside-ctag", "", "", "Z-side (remote) C-Tag/Inner-tag of the primary port (customer vlan id for QinQ)")

	connectionsCreateL2Cmd.Flags().StringVarP(&createL2ConSecondaryPortUUID, "sec-port-uuid", "", "", "Z-side (remote) secondary user port uuid")
	connectionsCreateL2Cmd.Flags().Int64VarP(&createL2ConSecondaryVlanSTag, "sec-port-stag", "", 0, "Z-side (remote) S-Tag/Outer-tag of the secondary port (vlan id for Dot1Q)")
	connectionsCreateL2Cmd.Flags().StringVarP(&createL2ConSecondaryVlanCTag, "sec-port-ctag", "", "", "Z-side (remote) C-Tag/Inner-tag of the primary port (customer vlan id for QinQ)")

	connectionsCreateL2Cmd.Flags().StringVarP(&createL2ConSecondaryZSidePortUUID, "sec-port-zside-uuid", "", "", "user port uuid")
	connectionsCreateL2Cmd.Flags().Int64VarP(&createL2ConSecondaryZSideVlanSTag, "sec-port-zside-stag", "", 0, "S-Tag/Outer-tag of the primary port (vlan id for Dot1Q)")
	connectionsCreateL2Cmd.Flags().Int64VarP(&createL2ConSecondaryZSideVlanCTag, "sec-port-zside-ctag", "", 0, "C-Tag/Inner-tag of the primary port (customer vlan id for QinQ)")

	connectionsCreateL2Cmd.Flags().StringVarP(&createL2ConSellerProfileUUID, "seller-uuid", "", "", "seller profile uuid (destination UUID)")
	connectionsCreateL2Cmd.Flags().StringVarP(&createL2ConSellerRegion, "seller-region", "", "", "seller destination region (ex. AWS: eu-west-1)")
	connectionsCreateL2Cmd.Flags().StringVarP(&createL2ConSellerMetroCode, "seller-metro", "", "", "seller destination metro code (ex.: LD)")

	connectionsCreateL2Cmd.Flags().Int64VarP(&createL2ConSpeed, "speed", "", 0, "connection speed (must be by 50 for MB ex.: 50, 100, 200, 500)")
	connectionsCreateL2Cmd.Flags().StringVarP(&createL2ConSpeedUnit, "speed-unit", "", "", "connection speed unit MB, GB")

	connectionsCreateL2Cmd.MarkFlagRequired("name")
	connectionsCreateL2Cmd.MarkFlagRequired("port-uuid")
	connectionsCreateL2Cmd.MarkFlagRequired("speed")
	connectionsCreateL2Cmd.MarkFlagRequired("speed-unit")

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

// connectionsCreateCloudCommand helper to assist in the creation of L2 connections to Cloud CSP's or other sellers in platform Equinix
func connectionsCreateCloudCommand(cmd *cobra.Command, args []string) {
	// required params
	// sellerService UUID - uuid to connect to
	// speed 50 mb
	// name (connection)
	// authorizationKey (AWS accountid as an example)
	// vlanSTag - vlan source tag
	// notifications email

	params := ConnectionsAPIClient.NewCreateL2ConnectionParams()

	if createL2ConAuthorizationKey == "" {
		log.Fatal("Authorization key required for cloud connection")
	}
	params.AuthorizationKey = createL2ConAuthorizationKey // aws account id in this case

	params.PrimaryName = createL2ConPrimaryName

	params.PrimaryPortUUID = createL2ConPrimaryPortUUID

	if createL2ConPrimaryVlanSTag == 0 {
		log.Fatal("Primary Vlan not specified (S-Tag)")
	}

	params.PrimaryVlanSTag = createL2ConPrimaryVlanSTag

	// secondary connection
	if createL2ConSecondaryName != "" {
		params.SecondaryName = createL2ConSecondaryName
	}
	if createL2ConSecondaryPortUUID != "" {
		params.SecondaryPortUUID = createL2ConSecondaryPortUUID
	}

	if createL2ConNamedTag != "" {
		params.NamedTag = createL2ConNamedTag
	}
	if createL2ConSecondaryVlanSTag != 0 {
		params.SecondaryVlanSTag = createL2ConSecondaryVlanSTag
	}
	if createL2ConSpeed == 0 {
		log.Fatal("Connection speed not specified")
	}

	params.Speed = createL2ConSpeed
	params.SpeedUnit = createL2ConSpeedUnit
	params.Notifications = []string{createL2ConNotificationsEmail}
	params.SellerRegion = createL2ConSellerRegion       //"eu-west-1" // get from seller? this should be AWS
	params.SellerMetroCode = createL2ConSellerMetroCode // provided by customer

	params.ProfileUUID = createL2ConSellerProfileUUID

	conn, err := ConnectionsAPIClient.CreateL2ConnectionToSellerProfile(params, SellerServicesAPIClient)
	if err != nil {
		switch t := err.(type) {
		case *apiconnections.CreateConnectionUsingPOSTBadRequest:
			for _, er := range t.Payload {
				log.Fatalf("Error %s with message %s\n", er.ErrorCode, er.ErrorMessage)
			}
		default:
			log.Fatalf("Error creating connection: %s\n", err.Error())
		}
	}
	fmt.Printf("Connection %s succesfully created\n", conn.Payload.PrimaryConnectionID)

	fmt.Printf("Connection %s succesfully created\n", conn.Payload.PrimaryConnectionID)
}

// connectionsCreateCommand creates a L2 connection
func connectionsCreateCommand(cmd *cobra.Command, args []string) {
	// required params
	// sellerService UUID - uuid to connect to
	// speed 50 mb
	// name (connection)
	// authorizationKey (AWS accountid as an example)
	// vlanSTag - vlan source tag
	// notifications email

	if createL2CSP {
		connectionsCreateCloudCommand(cmd, args)
		return
	}

	params := ConnectionsAPIClient.NewCreateL2ConnectionParams()

	if createL2ConNamedTag != "" {
		params.NamedTag = createL2ConNamedTag
	}

	if createL2ConAuthorizationKey != "" {
		params.AuthorizationKey = createL2ConAuthorizationKey
	}
	params.PrimaryName = createL2ConPrimaryName
	params.PrimaryPortUUID = createL2ConPrimaryPortUUID

	params.PrimaryVlanSTag = createL2ConPrimaryVlanSTag
	params.PrimaryVlanCTag = createL2ConPrimaryVlanCTag

	if createL2ConSpeed == 0 {
		log.Fatal("Connection speed not specified")
	}

	// Primary port zside
	if createL2ConPrimaryZSidePortUUID != "" {
		params.PrimaryZSidePortUUID = createL2ConPrimaryZSidePortUUID
		params.PrimaryVlanCTag = createL2ConPrimaryZSideVlanCTag
		params.PrimaryVlanSTag = createL2ConPrimaryZSideVlanSTag
	}

	// secondary port
	if createL2ConSecondaryName != "" {
		params.SecondaryName = createL2ConSecondaryName
	}
	if createL2ConSecondaryPortUUID != "" {
		params.SecondaryPortUUID = createL2ConSecondaryPortUUID
		params.SecondaryVlanSTag = createL2ConSecondaryVlanSTag
		params.SecondaryVlanCTag = createL2ConSecondaryVlanCTag
	}
	if createL2ConSecondaryZSidePortUUID != "" {
		params.SecondaryZSidePortUUID = createL2ConSecondaryZSidePortUUID
		params.SecondaryZSideVlanCTag = createL2ConSecondaryZSideVlanCTag
		params.SecondaryZSideVlanSTag = createL2ConSecondaryZSideVlanSTag
	}

	if purchaseOrderNumber != "" {
		params.PurchaseOrderNumber = purchaseOrderNumber
	}

	params.Speed = createL2ConSpeed
	params.SpeedUnit = createL2ConSpeedUnit
	params.Notifications = []string{createL2ConNotificationsEmail}
	params.SellerRegion = createL2ConSellerRegion
	params.SellerMetroCode = createL2ConSellerMetroCode

	// we're creating a connection to a seller profile

	conn, err := ConnectionsAPIClient.CreateL2Connection(params)
	if err != nil {
		switch t := err.(type) {
		case *apiconnections.CreateConnectionUsingPOSTBadRequest:
			for _, er := range t.Payload {
				log.Fatalf("Error %s with message %s\n", er.ErrorCode, er.ErrorMessage)
			}
		default:
			log.Fatalf("Error creating connection: %s\n", err.Error())
		}
	}
	fmt.Printf("Connection %s succesfully created\n", conn.Payload.PrimaryConnectionID)

}
