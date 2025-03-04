// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package components

import (
	"github.com/formancehq/ledger/pkg/client/internal/utils"
	"time"
)

type PostTransactionScript struct {
	Plain string         `json:"plain"`
	Vars  map[string]any `json:"vars,omitempty"`
}

func (o *PostTransactionScript) GetPlain() string {
	if o == nil {
		return ""
	}
	return o.Plain
}

func (o *PostTransactionScript) GetVars() map[string]any {
	if o == nil {
		return nil
	}
	return o.Vars
}

type PostTransaction struct {
	Timestamp *time.Time             `json:"timestamp,omitempty"`
	Postings  []Posting              `json:"postings,omitempty"`
	Script    *PostTransactionScript `json:"script,omitempty"`
	Reference *string                `json:"reference,omitempty"`
	Metadata  map[string]any         `json:"metadata,omitempty"`
}

func (p PostTransaction) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(p, "", false)
}

func (p *PostTransaction) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &p, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *PostTransaction) GetTimestamp() *time.Time {
	if o == nil {
		return nil
	}
	return o.Timestamp
}

func (o *PostTransaction) GetPostings() []Posting {
	if o == nil {
		return nil
	}
	return o.Postings
}

func (o *PostTransaction) GetScript() *PostTransactionScript {
	if o == nil {
		return nil
	}
	return o.Script
}

func (o *PostTransaction) GetReference() *string {
	if o == nil {
		return nil
	}
	return o.Reference
}

func (o *PostTransaction) GetMetadata() map[string]any {
	if o == nil {
		return nil
	}
	return o.Metadata
}
