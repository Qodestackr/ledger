// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package components

type LedgerInfoResponse struct {
	Data *LedgerInfo `json:"data,omitempty"`
}

func (o *LedgerInfoResponse) GetData() *LedgerInfo {
	if o == nil {
		return nil
	}
	return o.Data
}