package buyer

import (
	"fmt"
	"log"
	"math"

	api "github.com/jxoir/equinix-tools/pkg/ecxlib/api"
	apiseller_services "github.com/jxoir/go-ecxfabric/client/seller_services"
	"github.com/jxoir/go-ecxfabric/models"
)

type SellerServicesAPIHandler interface {
	GetAllSellerProfiles() (*apiseller_services.GetProfilesByMetroUsingGETOK, error)
}

type ECXSellerServicesAPI struct {
	*api.EquinixAPIClient
}

type SellerProfile struct {
	*models.GetServProfServicesRespContent
}

// SellerProfiles initial wrapper for swagger GetBuyerConResContent
type L2SellerProfiles struct {
	Items      []*models.GetServProfServicesRespContent
	TotalCount int64
	PageSize   int64
}

// L3SellerServices initial wrapper for swagger GetBuyerConResContent
type L3SellerServices struct {
	Items      []*models.SellerService
	TotalCount int64
	PageSize   int64
}

// NewECXSellerServicesAPI returns instantiated ECXSellerServicesAPI struct
func NewECXSellerServicesAPI(equinixAPIClient *api.EquinixAPIClient) *ECXSellerServicesAPI {
	return &ECXSellerServicesAPI{equinixAPIClient}
}

// GetAllL2SellerProfiles returns array of ports
func (ec *ECXSellerServicesAPI) GetAllL2SellerProfiles(metroCode *[]string) (*L2SellerProfiles, error) {
	// Remember that *profiles* are a L2 service profile

	respSellProfileList, err := ec.GetL2SellerProfiles(metroCode, nil, nil)

	if err != nil {
		return nil, err
	}

	totalCount := respSellProfileList.TotalCount
	pageSize := respSellProfileList.PageSize
	totalPages := int64(math.Ceil(float64(totalCount) / float64(pageSize)))

	// Start iterating from page 1 as we have "page 0" (yeah...swagger implementation of first page)
	psize := int32(pageSize)
	for p := 1; p <= int(totalPages-1); p++ {
		next := int32(p)
		req, err := ec.GetL2SellerProfiles(metroCode, &next, &psize)
		if err != nil {
			return nil, err
		} else {
			respSellProfileList.Items = append(append(respSellProfileList.Items, req.Items...))
		}

	}

	return respSellProfileList, nil

}

// GetL2SellerProfiles retrieve list of buyer connections for a specific page number and specific page size
func (ec *ECXSellerServicesAPI) GetL2SellerProfiles(metroCode *[]string, pageNumber *int32, pageSize *int32) (*L2SellerProfiles, error) {
	token, err := ec.GetToken()
	if err != nil {
		log.Fatal(err)
	}

	params := apiseller_services.NewGetProfilesByMetroUsingGETParams()

	if metroCode != nil {
		params.MetroCode = *metroCode
	}

	if pageNumber != nil {
		params.PageNumber = pageNumber
	}

	if pageSize != nil {
		params.PageSize = pageSize
	}

	respSellPOk, respSellNC, err := ec.Client.SellerServices.GetProfilesByMetroUsingGET(params, token)
	if err != nil {
		switch t := err.(type) {
		default:
			return nil, err
		case *apiseller_services.GetProfilesByMetroUsingGETBadRequest:
			if ec.Debug {
				fmt.Println(t.Error())
			}
			return nil, err
		}
	}

	if respSellPOk == nil {
		if ec.Debug {
			fmt.Println(respSellNC)
		}
		return nil, nil
	}

	respSellerProfilesList := L2SellerProfiles{
		Items:      respSellPOk.Payload.Content,
		TotalCount: respSellPOk.Payload.TotalCount,
		PageSize:   respSellPOk.Payload.PageSize,
	}

	return &respSellerProfilesList, nil
}

// GetAllL3SellerServices returns slice of L3SellerServices profiles (L3)
func (ec *ECXSellerServicesAPI) GetAllL3SellerServices(metroCode *[]string) (*L3SellerServices, error) {
	// Remember that *profiles* are a L2 service profile

	respSellProfileList, err := ec.GetL3SellerServices(metroCode, nil, nil)

	if err != nil {
		return nil, err
	}

	totalCount := respSellProfileList.TotalCount
	pageSize := respSellProfileList.PageSize
	totalPages := int64(math.Ceil(float64(totalCount) / float64(pageSize)))

	// Start iterating from page 1 as we have "page 0" (yeah...swagger implementation of first page)
	psize := int32(pageSize)
	for p := 1; p <= int(totalPages-1); p++ {
		next := int32(p)
		req, err := ec.GetL3SellerServices(metroCode, &next, &psize)
		if err != nil {
			return nil, err
		} else {
			respSellProfileList.Items = append(append(respSellProfileList.Items, req.Items...))
		}

	}

	return respSellProfileList, nil

}

// GetL3SellerServices retrieve list of seller services profiles (L3)
func (ec *ECXSellerServicesAPI) GetL3SellerServices(metroCode *[]string, pageNumber *int32, pageSize *int32) (*L3SellerServices, error) {
	token, err := ec.GetToken()
	if err != nil {
		log.Fatal(err)
	}

	params := apiseller_services.NewGetSellerServicesUsingGETParams()

	if metroCode != nil {
		params.Metros = *metroCode
	}

	if pageNumber != nil {
		params.Page = pageNumber
	}

	if pageSize != nil {
		params.Total = pageSize
	}

	respSellPOk, err := ec.Client.SellerServices.GetSellerServicesUsingGET(params, token)
	if err != nil {
		switch t := err.(type) {
		default:
			return nil, err
		case *apiseller_services.GetProfilesByMetroUsingGETBadRequest:
			if ec.Debug {
				fmt.Println(t.Error())
			}
			return nil, err
		}
	}

	respSellerProfilesList := L3SellerServices{
		Items:      respSellPOk.Payload.SellerServices,
		TotalCount: respSellPOk.Payload.TotalCount,
		PageSize:   respSellPOk.Payload.PageSize,
	}

	return &respSellerProfilesList, nil
}

// GetSellerProfileByUUID fetch service profile by uuid
func (ec *ECXSellerServicesAPI) GetSellerProfileByUUID(uuid string) (*apiseller_services.GetProfileByIDUsingGETOK, error) {
	token, err := ec.GetToken()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(uuid)
	params := apiseller_services.NewGetProfileByIDUsingGETParams()

	params.UUID = "3708e870-51fd-4168-b89a-7e79020c630c"

	sellerProfileOK, _, err := ec.Client.SellerServices.GetProfileByIDUsingGET(params, token)
	if err != nil {
		return nil, err
	}

	if sellerProfileOK != nil {
		return sellerProfileOK, nil
	}
	return nil, nil
}
