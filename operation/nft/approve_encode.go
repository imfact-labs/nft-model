package nft

import (
	"github.com/imfact-labs/currency-model/common"
	"github.com/imfact-labs/currency-model/types"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util/encoder"
	"github.com/pkg/errors"
)

func (fact *ApproveFact) unmarshal(enc encoder.Encoder, sd string, bit []byte, cid string) error {
	sender, err := base.DecodeAddress(sd, enc)
	if err != nil {
		return err
	}
	fact.sender = sender
	fact.currency = types.CurrencyID(cid)

	hit, err := enc.DecodeSlice(bit)
	if err != nil {
		return err
	}

	items := make([]ApproveItem, len(hit))
	for i, hinter := range hit {
		item, ok := hinter.(ApproveItem)
		if !ok {
			return common.ErrTypeMismatch.Wrap(errors.Errorf("expected ApproveItem, not %T", hinter))
		}

		items[i] = item
	}
	fact.items = items

	return nil
}
