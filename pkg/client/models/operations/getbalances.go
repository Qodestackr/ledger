// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/stack/ledger/client/internal/utils"
	"github.com/formancehq/stack/ledger/client/models/components"
)

type GetBalancesRequest struct {
	// Name of the ledger.
	Ledger string `pathParam:"style=simple,explode=false,name=ledger"`
	// Filter balances involving given account, either as source or destination.
	Address *string `queryParam:"style=form,explode=true,name=address"`
	// The maximum number of results to return per page.
	//
	PageSize *int64 `default:"15" queryParam:"style=form,explode=true,name=pageSize"`
	// Pagination cursor, will return accounts after given address, in descending order.
	After *string `queryParam:"style=form,explode=true,name=after"`
	// Parameter used in pagination requests. Maximum page size is set to 1000.
	// Set to the value of next for the next page of results.
	// Set to the value of previous for the previous page of results.
	// No other parameters can be set when this parameter is set.
	//
	Cursor *string `queryParam:"style=form,explode=true,name=cursor"`
}

func (g GetBalancesRequest) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(g, "", false)
}

func (g *GetBalancesRequest) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &g, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *GetBalancesRequest) GetLedger() string {
	if o == nil {
		return ""
	}
	return o.Ledger
}

func (o *GetBalancesRequest) GetAddress() *string {
	if o == nil {
		return nil
	}
	return o.Address
}

func (o *GetBalancesRequest) GetPageSize() *int64 {
	if o == nil {
		return nil
	}
	return o.PageSize
}

func (o *GetBalancesRequest) GetAfter() *string {
	if o == nil {
		return nil
	}
	return o.After
}

func (o *GetBalancesRequest) GetCursor() *string {
	if o == nil {
		return nil
	}
	return o.Cursor
}

type GetBalancesResponse struct {
	HTTPMeta components.HTTPMetadata `json:"-"`
	// OK
	BalancesCursorResponse *components.BalancesCursorResponse
}

func (o *GetBalancesResponse) GetHTTPMeta() components.HTTPMetadata {
	if o == nil {
		return components.HTTPMetadata{}
	}
	return o.HTTPMeta
}

func (o *GetBalancesResponse) GetBalancesCursorResponse() *components.BalancesCursorResponse {
	if o == nil {
		return nil
	}
	return o.BalancesCursorResponse
}