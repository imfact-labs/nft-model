package nft

import (
	"encoding/json"

	"github.com/imfact-labs/currency-model/common"
	"github.com/imfact-labs/currency-model/operation/extras"
	"github.com/imfact-labs/currency-model/types"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/encoder"
)

type ApproveAllFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Sender   base.Address     `json:"sender"`
	Items    []ApproveAllItem `json:"items"`
	Currency types.CurrencyID `json:"currency"`
}

func (fact ApproveAllFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ApproveAllFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Items:                 fact.items,
		Currency:              fact.currency,
	})
}

type ApproveAllFactJSONUnmarshaler struct {
	base.BaseFactJSONUnmarshaler
	Sender   string          `json:"sender"`
	Items    json.RawMessage `json:"items"`
	Currency string          `json:"currency"`
}

func (fact *ApproveAllFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var u ApproveAllFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(u.BaseFactJSONUnmarshaler)

	if err := fact.unmarshal(enc, u.Sender, u.Items, u.Currency); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

func (op ApproveAll) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(OperationMarshaler{
		BaseOperationJSONMarshaler:           op.BaseOperation.JSONMarshaler(),
		BaseOperationExtensionsJSONMarshaler: op.BaseOperationExtensions.JSONMarshaler(),
	})
}

func (op *ApproveAll) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperation = ubo

	var ueo extras.BaseOperationExtensions
	if err := ueo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperationExtensions = &ueo

	return nil
}
