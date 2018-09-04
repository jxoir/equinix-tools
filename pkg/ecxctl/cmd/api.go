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
	"log"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	apiclient "github.com/jxoir/go-ecxfabric/client"
	"github.com/jxoir/go-ecxfabric/client/access_token"
	"github.com/jxoir/go-ecxfabric/models"
)

var defaultGrantType = "client_credentials"

// NewAPIClient returns an instantiated client with token
func NewAPIClient(params *EquinixAPIConnectionParams, endpoint string) *EcxConnection {
	ecxClient := &EcxConnection{
		Params: params,
	}
	err := ecxClient.Connect(endpoint)
	if err != nil {
		log.Fatal("There was a problem connecting to endpoint " + endpoint)
		log.Fatal(err)
	}
	return ecxClient
}

// Connect to api endpoint
func (ec *EcxConnection) Connect(endpoint string) error {
	// set default parameters
	if ec.Params.AppID == "" {
		log.Fatal("EQUINIX_API_ID not set")
	}
	if ec.Params.AppSecret == "" {
		log.Fatal("EQUINIX_API_SECRET not set")
	}
	if ec.Params.GrantType == "" {
		ec.Params.AppSecret = defaultGrantType
	}
	// create the transport
	transport := httptransport.New(endpoint, "", nil)

	// create the API client, with the transport
	ecxAPIClient := apiclient.New(transport, strfmt.Default)

	accessTokenParams := access_token.NewGetAccessTokenParams()
	accessTokenRequest := models.OAuthRequest{
		ClientID:     ec.Params.AppID,
		ClientSecret: ec.Params.AppSecret,
		GrantType:    ec.Params.GrantType,
		UserName:     ec.Params.UserName,
		UserPassword: ec.Params.UserPassword,
	}

	if globalFlags.Debug {
		log.Println("User:" + ec.Params.UserName)
		log.Println("AppId:" + ec.Params.AppID)
		log.Println("Grant Type:" + ec.Params.GrantType)
	}

	accessTokenParams.SetRequest(&accessTokenRequest)
	accessTokenParams.Authorization = "Bearer"

	accessToken, err := ecxAPIClient.AccessToken.GetAccessToken(accessTokenParams, nil)
	if err != nil {

		log.Fatal("Failed to retrieve token...")
		log.Fatal(err)
	}

	if globalFlags.Debug {
		log.Println("Token acquired...")
	}

	bearerTokenAuth := httptransport.BearerToken(accessToken.Payload.AccessToken)

	ec.Client = ecxAPIClient
	ec.apiToken = bearerTokenAuth

	return nil
}

// EquinixAPIConnectionParams struct for generic Equinix params
type EquinixAPIConnectionParams struct {
	AppID        string
	AppSecret    string
	GrantType    string
	UserName     string
	UserPassword string
}

// EcxClient containing structure for Client, params and apitoken
// TODO: Implement token refresh
type EcxConnection struct {
	Client   *apiclient.GoEcxfabric
	Params   *EquinixAPIConnectionParams
	apiToken runtime.ClientAuthInfoWriter
}
