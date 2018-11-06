package client

import (
	"reflect"
	"strconv"
	"strings"
)

type ECXAPIResponse interface {
	GetItems() []interface{}
	AppendItems(items []interface{})
	SetItems(items []interface{})
	FilterItems(filters map[string]string)
	Count() int
}

type ECXAPIPayload interface {
	Get() []interface{}
	Count() int64
}

// ResponseFilter applies a set of map string filters to a Response and sets the new value to Response items
func ResponseFilter(response ECXAPIResponse, filters map[string]string) {
	// TODO: implement nested filtering
	// Obtain interface of slice
	r := response.GetItems()
	//make a zero lenght copy to append items later
	newitems := r[:0]
	for _, item := range r {
		t := reflect.TypeOf(item)
		// if type is a pointer then assign to the element
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		// we're accepting only structs
		if t.Kind() != reflect.Struct {
			continue
		}

		// iterate over struct fields to check tag annotations for each one
		filtermatch := false
		for i := 0; i < t.NumField(); i++ {
			if !filtermatch {
				f := t.Field(i)
				// skip unexported fields (PkgPath is the package pat that qualifies lower case - unexported - fields)
				if f.PkgPath != "" {
					continue
				}
				// try to obtain tags for json annotation will skip any other annotation
				tag := f.Tag.Get("json")
				if tag != "" && !filtermatch {
					for key, value := range filters {
						// parsetag options, tagName is the actual name optionos not used here
						tagName, _ := parseTag(tag)
						// if the actual field tag equals the key we are looking for
						if tagName == key {
							// obtain the struct by reflection
							r := reflect.ValueOf(item)
							// obtain the struct field value by reflection to compare (yep..inneficient and costly)
							valueField := reflect.Indirect(r).FieldByName(f.Name)
							// safe switch value, we need to verify each one...and cast the required ones
							switch valueField.Kind() {
							case reflect.String:
								//we are moving from a strict filter to a contains filter
								//if valueField.String() == value {
								if strings.Contains(valueField.String(), value) {
									newitems = append(newitems, item)
									filtermatch = true
								}

							case reflect.Int64:
								intv, err := strconv.ParseInt(value, 10, 32)
								if err != nil {
									// just continue, we werent able to parse the int of value
									continue
								}
								if valueField.Int() == intv {
									newitems = append(newitems, item)
									filtermatch = true
								}

							}

						}

					}
				}
			}

		}
	}
	// set filtered items to response
	response.SetItems(newitems)
}

/**
type PageResponse struct {
	Response ConnectionsResponse
	Error    error
}

// PageWalker traverses paging (traversing pagination)
func (m *ECXConnectionsAPI) PageWalker(metro *string) (*ConnectionsResponse, error) {
	//c := make(chan ConnectionsResult)
	//defer close(c)

	connectionsList, err := m.GetBuyerConnections(nil, nil, metro)
	if err != nil {
		return nil, err
	}

	totalCount := connectionsList.PageTotalCount
	pageSize := connectionsList.PageSize
	totalPages := int64(math.Ceil(float64(totalCount) / float64(pageSize)))
	// channels

	if pageSize > 0 && totalCount > 0 {
		// Start iterating from page 1 as we have "page 0" (yeah...swagger implementation of first page)
		c := make(chan PageResponse)
		var wg sync.WaitGroup
		var errorProc error

		psize := int32(pageSize)
		for p := 1; p <= int(totalPages-1); p++ {
			next := int32(p)
			wg.Add(1)
			go func() {
				defer wg.Done()
				req, err := m.GetBuyerConnections(&next, &psize, metro)
				res := PageResponse{
					Response: *req,
					Error:    err,
				}
				c <- res
			}()
		}

		go func() {
			wg.Wait()
			close(c)
		}()

		for resp := range c {
			if resp.Error != nil {
				errorProc = resp.Error
				break
			} else {
				connectionsList.AppendItems(resp.Response.GetItems())
			}
		}

		if errorProc != nil {
			return nil, errorProc
		}
	}

	return connectionsList, nil

}
**/
