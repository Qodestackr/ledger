// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/ledger/pkg/client/models/components"
)

type V2DeleteAccountMetadataRequest struct {
	// Name of the ledger.
	Ledger string `pathParam:"style=simple,explode=false,name=ledger"`
	// Account address
	Address string `pathParam:"style=simple,explode=false,name=address"`
	// The key to remove.
	Key string `pathParam:"style=simple,explode=false,name=key"`
}

func (o *V2DeleteAccountMetadataRequest) GetLedger() string {
	if o == nil {
		return ""
	}
	return o.Ledger
}

func (o *V2DeleteAccountMetadataRequest) GetAddress() string {
	if o == nil {
		return ""
	}
	return o.Address
}

func (o *V2DeleteAccountMetadataRequest) GetKey() string {
	if o == nil {
		return ""
	}
	return o.Key
}

type V2DeleteAccountMetadataResponse struct {
	HTTPMeta components.HTTPMetadata `json:"-"`
}

func (o *V2DeleteAccountMetadataResponse) GetHTTPMeta() components.HTTPMetadata {
	if o == nil {
		return components.HTTPMetadata{}
	}
	return o.HTTPMeta
}
