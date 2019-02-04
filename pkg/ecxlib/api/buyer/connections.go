package buyer

import (
	"errors"
	"fmt"
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

type CreateL2SellerConnectionParams struct {
	// authorization key
	AuthorizationKey string `json:"authorizationKey,omitempty"`

	// named tag
	NamedTag string `json:"namedTag,omitempty"`

	// notifications
	Notifications []string `json:"notifications"`

	// primary name
	PrimaryName string `json:"primaryName,omitempty"`

	// primary port UUID
	PrimaryPortUUID string `json:"primaryPortUUID,omitempty"`

	// primary vlan c tag
	PrimaryVlanCTag string `json:"primaryVlanCTag,omitempty"`

	// primary vlan s tag
	PrimaryVlanSTag int64 `json:"primaryVlanSTag,omitempty"`

	// primary z side port UUID
	PrimaryZSidePortUUID string `json:"primaryZSidePortUUID,omitempty"`

	// primary z side vlan c tag
	PrimaryZSideVlanCTag int64 `json:"primaryZSideVlanCTag,omitempty"`

	// primary z side vlan s tag
	PrimaryZSideVlanSTag int64 `json:"primaryZSideVlanSTag,omitempty"`

	// profile UUID
	ProfileUUID string `json:"profileUUID,omitempty"`

	// purchase order number
	PurchaseOrderNumber string `json:"purchaseOrderNumber,omitempty"`

	// secondary name
	SecondaryName string `json:"secondaryName,omitempty"`

	// secondary port UUID
	SecondaryPortUUID string `json:"secondaryPortUUID,omitempty"`

	// secondary vlan c tag
	SecondaryVlanCTag string `json:"secondaryVlanCTag,omitempty"`

	// secondary vlan s tag
	SecondaryVlanSTag int64 `json:"secondaryVlanSTag,omitempty"`

	// secondary z side port UUID
	SecondaryZSidePortUUID string `json:"secondaryZSidePortUUID,omitempty"`

	// secondary z side vlan c tag
	SecondaryZSideVlanCTag int64 `json:"secondaryZSideVlanCTag,omitempty"`

	// secondary z side vlan s tag
	SecondaryZSideVlanSTag int64 `json:"secondaryZSideVlanSTag,omitempty"`

	// seller metro code
	SellerMetroCode string `json:"sellerMetroCode,omitempty"`

	// seller region
	SellerRegion string `json:"sellerRegion,omitempty"`

	// speed
	Speed int64 `json:"speed,omitempty"`

	// speed unit
	SpeedUnit string `json:"speedUnit,omitempty"`
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

func (m *ECXConnectionsAPI) NewCreateL2SellerConnectionParams() *CreateL2SellerConnectionParams {
	return &CreateL2SellerConnectionParams{}
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

// DeleteByUUID delete connection by uuid
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

// CreateL2Connection creates an L2 connection to a specific service profile
func (m *ECXConnectionsAPI) CreateL2ConnectionSellerProfile(params *CreateL2SellerConnectionParams, ecxseller *ECXSellerServicesAPI) (*apiconnections.CreateConnectionUsingPOSTOK, error) {
	if params == nil {
		return nil, errors.New("Parameters to create L2 connection not provided")
	}

	if params.ProfileUUID == "" {
		return nil, errors.New("must provide seller profile UUID")
	}

	if params.PrimaryPortUUID == "" {
		return nil, errors.New("must provide a port id for the connection")
	}

	if m.Debug {
		log.Printf("Trying to obtain seller profile for UUID %s\n", params.ProfileUUID)
	}

	// first we obtain the seller profile
	seller, err := ecxseller.GetSellerProfileByUUID(params.ProfileUUID)
	if err != nil {
		s := fmt.Sprintf("can't obtain seller profile for %s UUID", params.ProfileUUID)
		return nil, errors.New(s)
	}

	if m.Debug {
		log.Printf("Trying to validate integration ID %s\n", seller.Payload.IntegrationID)
	}
	// validate the integrationId
	integrationIDOk, err := ecxseller.ValidateIntegrationID(seller.Payload.IntegrationID)
	if err != nil {
		return nil, err
	}

	if !integrationIDOk {
		s := fmt.Sprintf("Can't validate ontegration ID %s for seller profile %s UUID", seller.Payload.IntegrationID, params.ProfileUUID)
		return nil, errors.New(s)

	}

	// Validate that all the required information for the secondary port comes
	if params.SecondaryName != "" || params.SecondaryPortUUID != "" || params.SecondaryVlanSTag != 0 || params.NamedTag != "" {
		if params.SecondaryName == "" {
			return nil, errors.New("must provide a name for the secondary connection")
		}

		if params.SecondaryPortUUID == "" {
			return nil, errors.New("must provide the port id for the secondary connection")
		}

		if params.SecondaryVlanSTag == 0 {
			return nil, errors.New("must provide the vlan for the secondary connection")
		}

		// Validate that ports are not the same
		if params.PrimaryPortUUID == params.SecondaryPortUUID {
			return nil, errors.New("must provide a different port id for the secondary connection")
		}
	}

	// seller.Payload.IntegrationID

	ecxAPIParams := apiconnections.NewCreateConnectionUsingPOSTParams()
	request := &models.PostConnectionRequest{
		PrimaryName:       params.PrimaryName,
		PrimaryPortUUID:   params.PrimaryPortUUID,
		PrimaryVlanSTag:   params.PrimaryVlanSTag,
		SecondaryName:     params.SecondaryName,
		SecondaryPortUUID: params.SecondaryPortUUID,
		SecondaryVlanSTag: params.SecondaryVlanSTag,
		Speed:             params.Speed,
		SpeedUnit:         params.SpeedUnit,
		Notifications:     params.Notifications,
		SellerRegion:      params.SellerRegion,     //"eu-west-1" // get from seller? this should be AWS
		SellerMetroCode:   params.SellerMetroCode,  // provided by customer
		AuthorizationKey:  params.AuthorizationKey, // aws account id in this case
		ProfileUUID:       seller.Payload.UUID,
		NamedTag:          params.NamedTag,
	}

	ecxAPIParams.Request = request

	token, err := m.GetToken()
	if err != nil {
		log.Fatal(err)
	}

	connOk, err := m.Buyer.Connections.CreateConnectionUsingPOST(ecxAPIParams, token)
	if err != nil {
		// TODO create an APIError response interface and struct with a getmessages method
		/**
		switch t := err.(type) {
		case *apiconnections.CreateConnectionUsingPOSTBadRequest:
			return nil, errors.New(t.Payload)
		default:
			return nil, err
		}**/
		return nil, err
	}

	return connOk, nil

}
