package types

import (
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util/encoder"
	"github.com/imfact-labs/mitum2/util/hint"
)

func (sgn *Signer) unpack(
	enc encoder.Encoder,
	ht hint.Hint,
	ac string,
	sh uint,
	sg bool,
) error {
	sgn.BaseHinter = hint.NewBaseHinter(ht)
	sgn.share = sh
	sgn.signed = sg

	account, err := base.DecodeAddress(ac, enc)
	if err != nil {
		return err
	}
	sgn.address = account

	return nil
}
