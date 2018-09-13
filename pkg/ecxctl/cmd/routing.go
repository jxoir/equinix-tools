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
	"errors"
	"fmt"
	"log"
	"strings"

	apiroutinginstance "github.com/jxoir/go-ecxfabric/client/routing_instance"
	"github.com/spf13/cobra"
)

type RoutingInstanceAPIHandler interface {
	GetAllRoutingInstances() (*apiroutinginstance.GetAllRoutingInstancesUsingGETOK, error)
}

type ECXRoutingInstanceAPI struct {
	*EquinixAPIClient
}

type GetAllRoutingInstancesParams struct {
	MetroCode  *string
	PageSize   int32
	PageNumber int32
	States     []string
}

var routingInstanceStates string
var routingInstanceMetro string
var routingInstanceName string

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
	Run:   routingInstancesCheckRoutingInstanceNameExists,
}

var routingInstanceDeleteCmd cobra.Command
var routingInstanceUpdateCmd cobra.Command
var routingInstanceCreateCmd cobra.Command

func init() {
	rootCmd.AddCommand(routingInstanceCmd)
	routingInstanceCmd.AddCommand(routingInstanceListCmd)
	routingInstanceCmd.AddCommand(routingInstanceCheckNameCmd)

	routingInstanceListCmd.Flags().StringVarP(&routingInstanceMetro, "metro", "", "", "metro code")
	routingInstanceListCmd.Flags().StringVarP(&routingInstanceStates, "state", "", "PROVISIONED", "routing instances states")

	// Routing instances check name exists command definition and flags
	routingInstanceCheckNameCmd.Flags().StringVarP(&routingInstanceMetro, "metro", "", "", "metro code")
	routingInstanceCheckNameCmd.Flags().StringVarP(&routingInstanceName, "instance-name", "", "", "routing instance name")
	routingInstanceCheckNameCmd.MarkFlagRequired("metro")
	routingInstanceCheckNameCmd.MarkFlagRequired("instance-name")

}

// NewECXRoutingInstanceAPI returns instantiated ECXMetrosAPI struct
func NewECXRoutingInstanceAPI(equinixAPIClient *EquinixAPIClient) *ECXRoutingInstanceAPI {
	return &ECXRoutingInstanceAPI{equinixAPIClient}
}

func routingInstancesCheckRoutingInstanceNameExists(cmd *cobra.Command, args []string) {
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
	params := &GetAllRoutingInstancesParams{
		PageNumber: 1,
		PageSize:   10,
		States:     states,
		MetroCode:  &metro,
	}
	routingInstanceList, err := RoutingInstanceAPIClient.GetAllRoutingInstances(params)
	if err != nil {
		log.Fatal(err)
	} else {
		if routingInstanceList != nil {
			connRes, _ := json.MarshalIndent(routingInstanceList.Payload, "", "    ")
			fmt.Println(string(connRes))
		}
	}
}

// CheckRoutingInstanceNameExists returns bool or error
func (m *ECXRoutingInstanceAPI) CheckRoutingInstanceNameExists(name string, metroCode string) (bool, error) {

	token, err := m.GetToken()
	if err != nil {
		log.Fatal(err)
	}

	apiParams := apiroutinginstance.NewIsRoutingInstanceExistUsingGETParams()

	apiParams.MetroCode = metroCode
	apiParams.Name = name

	apiRespOk, apiRespNC, err := m.Client.RoutingInstance.IsRoutingInstanceExistUsingGET(apiParams, token)

	if err != nil {
		return false, err
	}

	if apiRespNC != nil {
		return false, errors.New("No content")
	}

	if apiRespOk.Payload.Exist {
		return true, nil
	}

	return false, nil
}

// GetAllRoutingInstances returns array of GetAllRoutingInstancesUsingGETOK with list of routing instances
func (m *ECXRoutingInstanceAPI) GetAllRoutingInstances(params *GetAllRoutingInstancesParams) (*apiroutinginstance.GetAllRoutingInstancesUsingGETOK, error) {
	if params == nil {
		params = &GetAllRoutingInstancesParams{
			PageNumber: 1,
			PageSize:   10,
		}
	}

	token, err := m.GetToken()
	if err != nil {
		log.Fatal(err)
	}

	apiParams := apiroutinginstance.NewGetAllRoutingInstancesUsingGETParams()

	apiParams.MetroCode = params.MetroCode
	apiParams.PageNumber = params.PageNumber
	apiParams.PageSize = params.PageSize
	apiParams.States = params.States

	respRoutingInstancesOk, _, err := m.Client.RoutingInstance.GetAllRoutingInstancesUsingGET(apiParams, token)
	if err != nil {
		switch t := err.(type) {
		default:
			if globalFlags.Debug {
				log.Println(err.Error())
			}
		case *json.UnmarshalTypeError:
			if globalFlags.Debug {
				log.Println(err.Error())
				log.Println(t.Value)
				log.Println(t.Struct)
				log.Println(t.Field)
				log.Println(t.Offset)
			}
		case *apiroutinginstance.GetAllRoutingInstancesUsingGETBadRequest:
			if globalFlags.Debug {
				log.Println("Bad request")
			}
		case *apiroutinginstance.GetAllRoutingInstancesUsingGETNoContent:
			if globalFlags.Debug {
				fmt.Println(t.Error())
			}
		}
		return nil, err

	}

	return respRoutingInstancesOk, nil

}
