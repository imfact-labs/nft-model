package nft

import (
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util/encoder"
	"github.com/imfact-labs/mitum2/util/hint"
)

func (it *TransferItem) unpack(
	enc encoder.Encoder,
	ht hint.Hint,
	ca, rc string,
	nid uint64,
) error {
	it.BaseHinter = hint.NewBaseHinter(ht)
	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return err
	default:
		it.contract = a
	}

	receiver, err := base.DecodeAddress(rc, enc)
	if err != nil {
		return err
	}
	it.receiver = receiver
	it.nftIdx = nid

	return nil
}
