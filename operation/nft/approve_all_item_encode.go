package nft

import (
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util/encoder"
	"github.com/imfact-labs/mitum2/util/hint"
)

func (it *ApproveAllItem) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	cAdr, dAdr, md string,
) error {
	it.BaseHinter = hint.NewBaseHinter(ht)

	it.mode = ApproveAllMode(md)

	switch a, err := base.DecodeAddress(cAdr, enc); {
	case err != nil:
		return err
	default:
		it.contract = a
	}

	delegatee, err := base.DecodeAddress(dAdr, enc)
	if err != nil {
		return err
	}
	it.approved = delegatee

	return nil
}
