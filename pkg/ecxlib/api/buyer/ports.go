package buyer

import (
	"fmt"
	"log"

	api "github.com/jxoir/equinix-tools/pkg/ecxlib/api/client"
	apiports "github.com/jxoir/go-ecxfabric/buyer/client/ports"
)

type PortsAPIHandler interface {
	GetAllMetros() (*apiports.GetPortInfoUsingGET2OK, error)
}

type ECXPortsAPI struct {
	*api.EquinixAPIClient
}

// NewECXPortsAPI returns instantiated ECXMetrosAPI struct
func NewECXPortsAPI(equinixAPIClient *api.EquinixAPIClient) *ECXPortsAPI {
	return &ECXPortsAPI{equinixAPIClient}
}

// GetAllPorts returns array of ports
func (ec *ECXPortsAPI) GetAllPorts() (*apiports.GetPortInfoUsingGET2OK, error) {
	token, err := ec.GetToken()
	if err != nil {
		log.Fatal(err)
	}
	respPortsOk, err := ec.Buyer.Ports.GetPortInfoUsingGET2(nil, token)
	if err != nil {
		switch t := err.(type) {
		default:
			log.Fatal(err)
		case *apiports.GetPortInfoUsingGET2NotFound:
			if ec.Debug {
				fmt.Println(t.Error())
			}
			return nil, err
		}
	}

	return respPortsOk, nil

}
