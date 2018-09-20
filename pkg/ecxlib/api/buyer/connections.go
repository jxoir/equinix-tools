package buyer

import (
	"fmt"
	"log"

	api "github.com/jxoir/equinix-tools/pkg/ecxlib/api"
	apiconnections "github.com/jxoir/go-ecxfabric/client/connections"
)

type ConnectionsAPIHandler interface {
	GetByUUID(uuid string) (*apiconnections.GetConnectionByUUIDUsingGETOK, error)
	GetAllBuyerConnections() (*apiconnections.GetAllBuyerConnectionsUsingGETOK, error)
}

type ECXConnectionsAPI struct {
	*api.EquinixAPIClient
}

// NewECXConnectionsAPI returns instantiated ECXConnectionsAPI struct
func NewECXConnectionsAPI(equinixAPIClient *api.EquinixAPIClient) *ECXConnectionsAPI {
	return &ECXConnectionsAPI{equinixAPIClient}
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
