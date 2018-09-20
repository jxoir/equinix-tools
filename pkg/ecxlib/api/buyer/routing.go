package buyer

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	api "github.com/jxoir/equinix-tools/pkg/ecxlib/api"
	apiroutinginstance "github.com/jxoir/go-ecxfabric/client/routing_instance"
	apiroutinginstancemodel "github.com/jxoir/go-ecxfabric/models"
)

type RoutingInstanceAPIHandler interface {
	GetAllRoutingInstances() (*apiroutinginstance.GetAllRoutingInstancesUsingGETOK, error)
}

type ECXRoutingInstanceAPI struct {
	*api.EquinixAPIClient
}

type GetAllRoutingInstancesParams struct {
	MetroCode  *string
	PageSize   int32
	PageNumber int32
	States     []string
}

type CreateRoutingInstanceParams struct {
	MetroCode           string
	PrimaryName         string
	SecondaryName       string
	RequiredRedundancy  bool
	RouteType           string
	Asn                 int64
	BgpUseAuth          bool
	BgpAuthorizationKey string
	NotificationEmails  []string
}

// NewECXRoutingInstanceAPI returns instantiated ECXMetrosAPI struct
func NewECXRoutingInstanceAPI(equinixAPIClient *api.EquinixAPIClient) *ECXRoutingInstanceAPI {
	return &ECXRoutingInstanceAPI{equinixAPIClient}
}

// CreateRoutingInstance returns routing instance primary uuid or error
func (ec *ECXRoutingInstanceAPI) CreateRoutingInstance(params *CreateRoutingInstanceParams) (string, error) {
	token, err := ec.GetToken()
	if err != nil {
		log.Fatal(err)
	}

	apiParams := apiroutinginstance.NewCreateRoutingInstanceUsingPOSTParams()
	routingRequest := &apiroutinginstancemodel.RoutingInstanceCreateRequest{}
	apiParams.Request = routingRequest

	apiParams.Request.Asn = params.Asn
	apiParams.Request.BgpAuthorizationKey = params.BgpAuthorizationKey
	apiParams.Request.MetroCode = params.MetroCode
	apiParams.Request.PrimaryRIName = params.PrimaryName
	apiParams.Request.SecondaryRIName = params.SecondaryName
	apiParams.Request.NotificationEmails = params.NotificationEmails
	apiParams.Request.RouteType = params.RouteType

	routingInstanceExists, err := ec.CheckRoutingInstanceNameExists(params.PrimaryName, params.MetroCode)
	if err != nil {
		return "", err
	}

	if routingInstanceExists {
		log.Fatal("Routing instance name already exists, please choose another name")
	}

	routingInstanceOk, routingInstanceNC, err := ec.Client.RoutingInstance.CreateRoutingInstanceUsingPOST(apiParams, token)
	if err != nil {
		return "", err
	}
	if routingInstanceNC != nil {
		return "", errors.New("No content")
	}

	return routingInstanceOk.Payload.PrimaryRIUUID, nil

}

// CheckRoutingInstanceNameExists returns bool or error
func (ec *ECXRoutingInstanceAPI) CheckRoutingInstanceNameExists(name string, metroCode string) (bool, error) {

	token, err := ec.GetToken()
	if err != nil {
		log.Fatal(err)
	}

	apiParams := apiroutinginstance.NewIsRoutingInstanceExistUsingGETParams()

	apiParams.MetroCode = metroCode
	apiParams.Name = name

	apiRespOk, apiRespNC, err := ec.Client.RoutingInstance.IsRoutingInstanceExistUsingGET(apiParams, token)

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
func (ec *ECXRoutingInstanceAPI) GetAllRoutingInstances(params *GetAllRoutingInstancesParams) (*apiroutinginstance.GetAllRoutingInstancesUsingGETOK, error) {
	if params == nil {
		params = &GetAllRoutingInstancesParams{
			PageNumber: 1,
			PageSize:   10,
		}
	}

	token, err := ec.GetToken()
	if err != nil {
		log.Fatal(err)
	}

	apiParams := apiroutinginstance.NewGetAllRoutingInstancesUsingGETParams()

	apiParams.MetroCode = params.MetroCode
	apiParams.PageNumber = params.PageNumber
	apiParams.PageSize = params.PageSize
	apiParams.States = params.States

	respRoutingInstancesOk, _, err := ec.Client.RoutingInstance.GetAllRoutingInstancesUsingGET(apiParams, token)
	if err != nil {
		switch t := err.(type) {
		default:
			if ec.Debug {
				log.Println(err.Error())
			}
		case *json.UnmarshalTypeError:
			if ec.Debug {
				log.Println(err.Error())
				log.Println(t.Value)
				log.Println(t.Struct)
				log.Println(t.Field)
				log.Println(t.Offset)
			}
		case *apiroutinginstance.GetAllRoutingInstancesUsingGETBadRequest:
			if ec.Debug {
				log.Println("Bad request")
			}
		case *apiroutinginstance.GetAllRoutingInstancesUsingGETNoContent:
			if ec.Debug {
				fmt.Println(t.Error())
			}
		}
		return nil, err

	}

	return respRoutingInstancesOk, nil

}
