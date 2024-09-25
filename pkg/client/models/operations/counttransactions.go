// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/stack/ledger/client/internal/utils"
	"github.com/formancehq/stack/ledger/client/models/components"
	"time"
)

// Metadata - Filter transactions by metadata key value pairs. Nested objects can be used as seen in the example below.
type Metadata struct {
}

type CountTransactionsRequest struct {
	// Name of the ledger.
	Ledger string `pathParam:"style=simple,explode=false,name=ledger"`
	// Filter transactions by reference field.
	Reference *string `queryParam:"style=form,explode=true,name=reference"`
	// Filter transactions with postings involving given account, either as source or destination (regular expression placed between ^ and $).
	Account *string `queryParam:"style=form,explode=true,name=account"`
	// Filter transactions with postings involving given account at source (regular expression placed between ^ and $).
	Source *string `queryParam:"style=form,explode=true,name=source"`
	// Filter transactions with postings involving given account at destination (regular expression placed between ^ and $).
	Destination *string `queryParam:"style=form,explode=true,name=destination"`
	// Filter transactions that occurred after this timestamp.
	// The format is RFC3339 and is inclusive (for example, "2023-01-02T15:04:01Z" includes the first second of 4th minute).
	//
	StartTime *time.Time `queryParam:"style=form,explode=true,name=startTime"`
	// Filter transactions that occurred before this timestamp.
	// The format is RFC3339 and is exclusive (for example, "2023-01-02T15:04:01Z" excludes the first second of 4th minute).
	//
	EndTime *time.Time `queryParam:"style=form,explode=true,name=endTime"`
	// Filter transactions by metadata key value pairs. Nested objects can be used as seen in the example below.
	Metadata *Metadata `queryParam:"style=deepObject,explode=true,name=metadata"`
}

func (c CountTransactionsRequest) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(c, "", false)
}

func (c *CountTransactionsRequest) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &c, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *CountTransactionsRequest) GetLedger() string {
	if o == nil {
		return ""
	}
	return o.Ledger
}

func (o *CountTransactionsRequest) GetReference() *string {
	if o == nil {
		return nil
	}
	return o.Reference
}

func (o *CountTransactionsRequest) GetAccount() *string {
	if o == nil {
		return nil
	}
	return o.Account
}

func (o *CountTransactionsRequest) GetSource() *string {
	if o == nil {
		return nil
	}
	return o.Source
}

func (o *CountTransactionsRequest) GetDestination() *string {
	if o == nil {
		return nil
	}
	return o.Destination
}

func (o *CountTransactionsRequest) GetStartTime() *time.Time {
	if o == nil {
		return nil
	}
	return o.StartTime
}

func (o *CountTransactionsRequest) GetEndTime() *time.Time {
	if o == nil {
		return nil
	}
	return o.EndTime
}

func (o *CountTransactionsRequest) GetMetadata() *Metadata {
	if o == nil {
		return nil
	}
	return o.Metadata
}

type CountTransactionsResponse struct {
	HTTPMeta components.HTTPMetadata `json:"-"`
	Headers  map[string][]string
}

func (o *CountTransactionsResponse) GetHTTPMeta() components.HTTPMetadata {
	if o == nil {
		return components.HTTPMetadata{}
	}
	return o.HTTPMeta
}

func (o *CountTransactionsResponse) GetHeaders() map[string][]string {
	if o == nil {
		return map[string][]string{}
	}
	return o.Headers
}