package buyer

import (
	"fmt"
	"log"

	api "github.com/jxoir/equinix-tools/pkg/ecxlib/api/client"
	apimetros "github.com/jxoir/go-ecxfabric/buyer/client/metros"
)

type MetrosAPIHandler interface {
	GetAllMetros() (*apimetros.GetMetrosUsingGETOK, error)
}

type ECXMetrosAPI struct {
	*api.EquinixAPIClient
}

// NewECXMetrosAPI returns instantiated ECXMetrosAPI struct
func NewECXMetrosAPI(equinixAPIClient *api.EquinixAPIClient) *ECXMetrosAPI {
	return &ECXMetrosAPI{equinixAPIClient}
}

// GetAllMetros returns array of GetMetrosUsingGETOK with list of customer connections
func (ec *ECXMetrosAPI) GetAllMetros() (*apimetros.GetMetrosUsingGETOK, error) {
	token, err := ec.GetToken()
	if err != nil {
		log.Fatal(err)
	}
	respMetrosOk, _, err := ec.Buyer.Metros.GetMetrosUsingGET(nil, token)
	if err != nil {
		switch t := err.(type) {
		default:
			log.Fatal(err)
		case *apimetros.GetMetrosUsingGETNoContent:
			if ec.Debug {
				fmt.Println(t.Error())
			}
			return nil, err
		}
	}

	return respMetrosOk, nil

}
