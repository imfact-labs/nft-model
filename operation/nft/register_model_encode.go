package nft

import (
	ctypes "github.com/imfact-labs/currency-model/types"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util/encoder"
	"github.com/imfact-labs/nft-model/types"
)

func (fact *RegisterModelFact) unmarshal(
	enc encoder.Encoder,
	sd string,
	ca string,
	nm string,
	ry uint,
	uri string,
	bws []string,
	cid string,
) error {
	fact.currency = ctypes.CurrencyID(cid)

	sender, err := base.DecodeAddress(sd, enc)
	if err != nil {
		return err
	}
	fact.sender = sender

	fact.name = types.CollectionName(nm)
	fact.royalty = types.PaymentParameter(ry)
	fact.uri = types.URI(uri)

	contract, err := base.DecodeAddress(ca, enc)
	if err != nil {
		return err
	}
	fact.contract = contract

	whitelist := make([]base.Address, len(bws))
	for i, bw := range bws {
		white, err := base.DecodeAddress(bw, enc)
		if err != nil {
			return err
		}
		whitelist[i] = white

	}
	fact.minterWhitelist = whitelist

	return nil
}
