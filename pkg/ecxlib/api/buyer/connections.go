package buyer

import (
	"log"
	"math"

	client "github.com/jxoir/equinix-tools/pkg/ecxlib/api/client"
	apiconnections "github.com/jxoir/go-ecxfabric/buyer/client/connections"
	"github.com/jxoir/go-ecxfabric/buyer/models"
)

// ECXConnectionsAPI Connections api client container
type ECXConnectionsAPI struct {
	*client.EquinixAPIClient
}

// ConnectionsResponse initial wrapper for swagger GetBuyerConResContent
type ConnectionsResponse struct {
	Items          []interface{}
	PageTotalCount int64
	PageSize       int64
}

// AppendItems appends slice of interface items to internal Items
func (c *ConnectionsResponse) AppendItems(items []interface{}) {
	c.Items = append(c.Items, items...)
}

// SetItems replaces internal slice of interface items with items
func (c *ConnectionsResponse) SetItems(items []interface{}) {
	c.Items = items
}

// GetItems retrieves all Items
func (c *ConnectionsResponse) GetItems() []interface{} {
	return c.Items
}

// FilterItems applies specific filters to items and updates internal items
func (c *ConnectionsResponse) FilterItems(filters map[string]string) {
	client.ResponseFilter(c, filters)
}

// Count return total count of items
func (c *ConnectionsResponse) Count() int {
	return len(c.Items)
}

// parseContent copy from the returned slice of pointers GetBuyerConResContent to a []interface{}
func (c *ConnectionsResponse) parseContent(payload []*models.GetBuyerConResContent) []interface{} {
	interfaceSlice := make([]interface{}, len(payload))
	for i, d := range payload {
		//interfaceSlice[i] = d
		interfaceSlice[i] = d
	}
	return interfaceSlice
}

// NewECXConnectionsAPI returns instantiated ECXConnectionsAPI struct
func NewECXConnectionsAPI(equinixAPIClient *client.EquinixAPIClient) *ECXConnectionsAPI {
	return &ECXConnectionsAPI{equinixAPIClient}
}

// GetAllBuyerConnections get all buyer connections (traversing pagination)
func (m *ECXConnectionsAPI) GetAllBuyerConnections(metro *string) (*ConnectionsResponse, error) {

	connectionsList, err := m.GetBuyerConnections(nil, nil, metro)
	if err != nil {
		return nil, err
	}

	totalCount := connectionsList.PageTotalCount
	pageSize := connectionsList.PageSize
	totalPages := int64(math.Ceil(float64(totalCount) / float64(pageSize)))

	if pageSize > 0 && totalCount > 0 {
		// Start iterating from page 1 as we have "page 0" (yeah...swagger implementation of first page)
		psize := int32(pageSize)
		for p := 1; p <= int(totalPages-1); p++ {
			next := int32(p)
			req, err := m.GetBuyerConnections(&next, &psize, metro)
			if err != nil {
				return nil, err
			} else {
				//connectionsList.Items = append(append(connectionsList.Items, req.Items...))
				connectionsList.AppendItems(req.Items)
			}

		}
	}
	return connectionsList, nil

}

// GetBuyerConnections retrieve list of buyer connections for a specific page number and specific page size
func (m *ECXConnectionsAPI) GetBuyerConnections(pageNumber *int32, pageSize *int32, metro *string) (*ConnectionsResponse, error) {
	token, err := m.GetToken()
	if err != nil {
		log.Fatal(err)
	}

	params := apiconnections.NewGetAllBuyerConnectionsUsingGETParams()

	if metro != nil && *metro != "" {
		params.MetroCode = metro
	}

	if pageNumber != nil {
		params.PageNumber = pageNumber
	}

	if pageSize != nil {
		params.PageSize = pageSize
	}

	connectionsList := ConnectionsResponse{}

	connectionsOK, connectionsNC, err := m.Buyer.Connections.GetAllBuyerConnectionsUsingGET(params, token)
	if connectionsNC != nil {
		connectionsList.PageTotalCount = 0
		connectionsList.PageSize = 0

		return &connectionsList, nil

	}
	if err != nil {

		switch t := err.(type) {
		default:
			return nil, err
		case *apiconnections.GetAllBuyerConnectionsUsingGETBadRequest:
			// specific get bad request for envelope errors...
			return nil, t
		}
	}

	connectionsList.PageSize = connectionsOK.Payload.PageSize
	connectionsList.PageTotalCount = connectionsOK.Payload.TotalCount
	connectionsList.SetItems(connectionsList.parseContent(connectionsOK.Payload.Content))

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

	connectionOK, _, err := m.Buyer.Connections.GetConnectionByUUIDUsingGET(params, token)
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

// DeleteByUUID get connection by uuid
func (m *ECXConnectionsAPI) DeleteByUUID(uuid string) (*apiconnections.DeleteConnectionUsingDELETEOK, error) {
	params := apiconnections.NewDeleteConnectionUsingDELETEParams()
	params.SetConnID(uuid)

	token, err := m.GetToken()
	if err != nil {
		log.Fatal(err)
	}

	deleteOK, err := m.Buyer.Connections.DeleteConnectionUsingDELETE(params, token)
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

	return deleteOK, nil

}

// CreateL2Connection

/**
func (m *ECXConnectionsAPI) CreateL2Connection(params L2ConnectionParams) (*apiconnections.CreateConnectionUsingPOSTOK, error) {
	p := apiconnections.NewCreateConnectionUsingPOSTParams
	apiconnectionsmodel.PostConnectionRequest
	connectionOk, err := m.Buyer.Connections.CreateConnectionUsingPOST(p, token)
	if err != nil {
		return _, err
	}

	return connectionOk, _
}
**/
