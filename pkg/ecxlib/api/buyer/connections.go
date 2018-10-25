package buyer

import (
	"log"
	"math"

	api "github.com/jxoir/equinix-tools/pkg/ecxlib/api"
	apiconnections "github.com/jxoir/go-ecxfabric/client/connections"
	"github.com/jxoir/go-ecxfabric/models"
)

// Connections initial wrapper for swagger GetBuyerConResContent
type Connections struct {
	Items      []*models.GetBuyerConResContent
	TotalCount int64
	PageSize   int64
}

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

// GetAllBuyerConnections get all buyer connections (traversing pagination)
func (m *ECXConnectionsAPI) GetAllBuyerConnections() (*Connections, error) {

	connectionsList, err := m.GetBuyerConnections(nil, nil)
	if err != nil {
		return nil, err
	}

	totalCount := connectionsList.TotalCount
	pageSize := connectionsList.PageSize
	totalPages := int64(math.Ceil(float64(totalCount) / float64(pageSize)))

	// Start iterating from page 1 as we have "page 0" (yeah...swagger implementation of first page)
	psize := int32(pageSize)
	for p := 1; p <= int(totalPages-1); p++ {
		next := int32(p)
		req, err := m.GetBuyerConnections(&next, &psize)
		if err != nil {
			return nil, err
		} else {
			connectionsList.Items = append(append(connectionsList.Items, req.Items...))
		}

	}

	return connectionsList, nil

}

// GetBuyerConnections retrieve list of buyer connections for a specific page number and specific page size
func (m *ECXConnectionsAPI) GetBuyerConnections(pageNumber *int32, pageSize *int32) (*Connections, error) {
	token, err := m.GetToken()
	if err != nil {
		log.Fatal(err)
	}

	params := apiconnections.NewGetAllBuyerConnectionsUsingGETParams()

	if pageNumber != nil {
		params.PageNumber = pageNumber
	}

	if pageSize != nil {
		params.PageSize = pageSize
	}

	connectionsOK, _, err := m.Client.Connections.GetAllBuyerConnectionsUsingGET(params, token)
	if err != nil {
		switch t := err.(type) {
		default:
			return nil, err
		case *apiconnections.GetAllBuyerConnectionsUsingGETBadRequest:
			// specific get bad request for envelope errors...
			return nil, t
		}
	}

	connectionsList := Connections{
		Items:      connectionsOK.Payload.Content,
		TotalCount: connectionsOK.Payload.TotalCount,
		PageSize:   connectionsOK.Payload.PageSize,
	}

	return &connectionsList, nil
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
			// check for envelope message in badrequest - deprecated
			/** for _, getconnerrors := range t.Payload {
				fmt.Println(getconnerrors.ErrorMessage + ":" + uuid)
				return nil, err
			}**/
			return nil, t
		}
	}

	if connectionOK != nil {
		return connectionOK, nil
	}

	return nil, err

}

// CreateL2Connection

/**
func (m *ECXConnectionsAPI) CreateL2Connection(params L2ConnectionParams) (*apiconnections.CreateConnectionUsingPOSTOK, error) {
	p := apiconnections.NewCreateConnectionUsingPOSTParams
	apiconnectionsmodel.PostConnectionRequest
	connectionOk, err := m.Client.Connections.CreateConnectionUsingPOST(p, token)
	if err != nil {
		return _, err
	}

	return connectionOk, _
}
**/
