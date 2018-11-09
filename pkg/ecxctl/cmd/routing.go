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

	"github.com/jxoir/equinix-tools/pkg/ecxlib/api/buyer"
	"github.com/spf13/cobra"
)

var routingInstanceStates string
var routingInstanceMetro string
var routingInstanceName string

var routingInstanceSecondaryName string
var routingInstanceRequiredRedundancy bool
var routingInstanceRouteType string
var routingInstanceAsn int64
var routingInstanceBgpUseAuth bool
var routingInstanceBgpAuthorizationKey string
var routingInstanceNotificationEmails []string

// metrosCmd represents the metros command
var routingInstanceCmd = &cobra.Command{
	Use:   "routing-instance",
	Short: "Operations related to ECX Metros",
}

// metrosCmd represents the metros command
var routingInstanceListCmd = &cobra.Command{
	Use:   "list",
	Short: "list available routing instances",
	Run:   routingInstancesListCommand,
}

// metrosCmd represents the metros command
var routingInstanceCheckNameCmd = &cobra.Command{
	Use:   "check-name",
	Short: "check Routing Instance name exists or not",
	Run:   routingInstancesCheckRoutingInstanceNameExistsCommand,
}

var routingInstanceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create new Routing Instance",
	Run:   routingInstanceCreateCommand,
}

var routingInstanceDeleteCmd cobra.Command
var routingInstanceUpdateCmd cobra.Command

func init() {
	rootCmd.AddCommand(routingInstanceCmd)
	routingInstanceCmd.AddCommand(routingInstanceListCmd)
	routingInstanceCmd.AddCommand(routingInstanceCheckNameCmd)
	routingInstanceCmd.AddCommand(routingInstanceCreateCmd)

	routingInstanceListCmd.Flags().StringVarP(&routingInstanceMetro, "metro", "", "", "metro code")
	routingInstanceListCmd.Flags().StringVarP(&routingInstanceStates, "state", "", "PROVISIONED", "routing instances states")

	// Routing instances check name exists command definition and flags
	routingInstanceCheckNameCmd.Flags().StringVarP(&routingInstanceMetro, "metro", "", "", "metro code")
	routingInstanceCheckNameCmd.Flags().StringVarP(&routingInstanceName, "instance-name", "", "", "routing instance name")
	routingInstanceCheckNameCmd.MarkFlagRequired("metro")
	routingInstanceCheckNameCmd.MarkFlagRequired("instance-name")

	// routeInstanceCreateCmd
	routingInstanceCreateCmd.Flags().StringVarP(&routingInstanceMetro, "metro", "", "", "metro code")
	routingInstanceCreateCmd.Flags().StringVarP(&routingInstanceName, "name", "", "", "routing instance name (primary)")
	routingInstanceCreateCmd.Flags().StringVarP(&routingInstanceSecondaryName, "secondary-name", "", "", "routing instance secondary name")
	routingInstanceCreateCmd.Flags().BoolVarP(&routingInstanceRequiredRedundancy, "redundancy", "", false, "required redundancy")
	routingInstanceCreateCmd.Flags().Int64VarP(&routingInstanceAsn, "asn-number", "", 0, "asn number")
	routingInstanceCreateCmd.Flags().BoolVarP(&routingInstanceBgpUseAuth, "bgp-auth", "", false, "required ASN BGP authentication")
	routingInstanceCreateCmd.Flags().StringVarP(&routingInstanceBgpAuthorizationKey, "bgp-key", "", "", "bgp auth key if required")
	routingInstanceCreateCmd.Flags().StringArrayVarP(&routingInstanceNotificationEmails, "notification-emails", "", []string{}, "notification emails (comma separated)")
	routingInstanceCreateCmd.Flags().StringVarP(&routingInstanceRouteType, "type", "", "Private", "route type Private or Public")

	routingInstanceCreateCmd.MarkFlagRequired("metro")
	routingInstanceCreateCmd.MarkFlagRequired("name")
	routingInstanceCreateCmd.MarkFlagRequired("secondary-name")
	routingInstanceCreateCmd.MarkFlagRequired("asn")
	routingInstanceCreateCmd.MarkFlagRequired("notification-emails")
	routingInstanceCreateCmd.MarkFlagRequired("type")

}

func routingInstancesCheckRoutingInstanceNameExistsCommand(cmd *cobra.Command, args []string) {
	resp, err := RoutingInstanceAPIClient.CheckRoutingInstanceNameExists(routingInstanceName, routingInstanceMetro)
	if err != nil {
		log.Fatal(err)
	}
	if resp {
		fmt.Println("Routing instance name exists")
	} else {
		fmt.Println("Routing instance name doesn't exists")
	}
}

func routingInstancesListCommand(cmd *cobra.Command, args []string) {
	// Call GetRoutingInstance without params

	metro := ""
	if routingInstanceMetro != "" {
		metro = routingInstanceMetro
	}
	states := strings.Split(routingInstanceStates, ",")
	params := buyer.GetAllRoutingInstancesParams{
		PageNumber: 1,
		PageSize:   10,
		States:     states,
		MetroCode:  &metro,
	}
	routingInstanceList, err := RoutingInstanceAPIClient.GetAllRoutingInstances(&params)
	if err != nil {
		log.Fatal(err)
	} else {
		if routingInstanceList != nil {

			routingInstances := routingInstanceList.Payload.RoutingInstances
			routingInstancesRes, err := json.MarshalIndent(routingInstances, "", "    ")
			if err != nil {
				log.Fatal("There was an error with json response:", err)
			} else {
				fmt.Println(string(routingInstancesRes))
			}
		}

	}
}

func routingInstanceCreateCommand(cmd *cobra.Command, args []string) {

	// params = CreateRoutingInstanceParams{}

	if routingInstanceBgpUseAuth && routingInstanceBgpAuthorizationKey == "" {
		fmt.Println("BGP authorization key required")
		return
	}

	params := buyer.CreateRoutingInstanceParams{
		MetroCode:           routingInstanceMetro,
		PrimaryName:         routingInstanceName,
		SecondaryName:       routingInstanceSecondaryName,
		RequiredRedundancy:  routingInstanceRequiredRedundancy,
		RouteType:           routingInstanceRouteType,
		Asn:                 routingInstanceAsn,
		BgpUseAuth:          routingInstanceBgpUseAuth,
		BgpAuthorizationKey: routingInstanceBgpAuthorizationKey,
		NotificationEmails:  routingInstanceNotificationEmails,
	}

	riUUID, err := RoutingInstanceAPIClient.CreateRoutingInstance(&params)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Routing instance " + routingInstanceName + " created:" + riUUID)

}
