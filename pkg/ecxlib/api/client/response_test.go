package client

import (
	"encoding/json"
	"testing"

	"github.com/jxoir/go-ecxfabric/buyer/models"
)

var portUUID = "cb7cc58a-b323-473a-b40d-b2d4be61e3f1"
var portName = "EQUINIX-LD4-CX-PRI-01"
var authorizationKey = "402278354059" // AWS Account in case of AWS
var zsidePortUUID = "9350fd44-0883-4aaa-b266-613d33dd0c95"
var connectionUUID = "07c8a274-4e80-4662-8cb9-636b8b00eb26"
var sellerServiceUUID = "69ee618d-be52-468d-bc99-00566f2dd2b9" // Real one
var zSidePortName = "dxcon-f362p08u"                           // HVI in AWS

var responseJSON = `{"asideEncapsulation":"dot1q","authorizationKey":"402278354059","billingTier":"Up to 50 MB","buyerOrganizationName":"SVC EQUINIX GLOBAL","createdBy":"someuser@email.com","createdByEmail":"jxoir@github.com","createdByFullName":"Juancho","createdDate":"2018-10-29T16:40:55.470Z","lastUpdatedDate":"2018-10-29T16:51:05.061Z","metadata":{"integration_id":"AWS-DirectConnect-01","notification_emails":["jxoir@github.com"]},"metroCode":"LD","metroDescription":"London","name":"EQUINIX_TEST","notifications":["jxoir@github.com"],"portName":"EQUINIX-LD4-CX-PRI-01","portUUID":"cb7cc58a-b323-473a-b40d-b2d4be61e3f1","redundancyType":"primary","sellerMetroCode":"LD","sellerMetroDescription":"London","sellerOrganizationName":"EQUINIX-AWS","sellerServiceName":"AWS Direct Connect","sellerServiceUUID":"69ee618d-be52-468d-bc99-00566f2dd2b9","speed":50,"speedUnit":"MB","status":"PROVISIONED","uuid":"07c8a274-4e80-4662-8cb9-636b8b00eb26","vlanSTag":3022,"zSidePortName":"dxcon-f362p08u","zSidePortUUID":"9350fd44-0883-4aaa-b266-613d33dd0c95","zSideVlanSTag":333}`

// Mock of a ECXAPIResponse
// ConnectionsResponse initial wrapper for swagger GetBuyerConResContent
type ECXAPIResponseMock struct {
	Items          []interface{}
	PageTotalCount int64
	PageSize       int64
}

// AppendItems appends slice of interface items to internal Items
func (c *ECXAPIResponseMock) AppendItems(items []interface{}) {
	c.Items = append(c.Items, items...)
}

// SetItems replaces internal slice of interface items with items
func (c *ECXAPIResponseMock) SetItems(items []interface{}) {
	c.Items = items
}

// GetItems retrieves all Items
func (c *ECXAPIResponseMock) GetItems() []interface{} {
	return c.Items
}

// FilterItems applies specific filters to items and updates internal items
func (c *ECXAPIResponseMock) FilterItems(filters map[string]string) {
	ResponseFilter(c, filters)
}

// Count return total count of items
func (c *ECXAPIResponseMock) Count() int {
	return len(c.Items)
}

func TestResponseFilter(t *testing.T) {

	// create a mock of ECX API Response
	resp := ECXAPIResponseMock{}
	connAPIResponse := models.GETConnectionByUUIDResponse{}

	// Unmarshall sample json into mock up struct
	err := json.Unmarshal([]byte(responseJSON), &connAPIResponse)

	if err != nil {
		t.Errorf("Can't unmarshal test JSON into connection struct")
	}

	// Create slice of interface to represent mocked up api response content
	items := []interface{}{connAPIResponse}
	resp.SetItems(items)

	// Setup filters to PASS case
	filters := make(map[string]string)
	filters["portName"] = "EQUINIX"

	// Test applying filters, expecting
	resp.FilterItems(filters)
	if resp.Count() != 1 {
		t.Errorf("Expected 1 items after filtering, received %d", resp.Count())
	}

	// Setup filters to fail filtering case
	filters["portName"] = "METROCODE"

	// Test applying filters, expecting
	resp.FilterItems(filters)
	if resp.Count() != 0 {
		t.Errorf("Expected 0 items after filtering, received %d", resp.Count())
	}

}
