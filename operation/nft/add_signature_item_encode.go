package nft

import (
	"github.com/imfact-labs/currency-model/types"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util/encoder"
	"github.com/imfact-labs/mitum2/util/hint"
)

func (it *AddSignatureItem) unpack(
	enc encoder.Encoder,
	ht hint.Hint,
	ca string,
	nft uint64,
	cid string,
) error {
	it.BaseHinter = hint.NewBaseHinter(ht)
	it.currency = types.CurrencyID(cid)
	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return err
	default:
		it.contract = a
	}

	it.nftIdx = nft

	return nil
}
