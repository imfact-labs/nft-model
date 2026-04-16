package nft

import (
	"github.com/imfact-labs/currency-model/common"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/encoder"
	"github.com/imfact-labs/mitum2/util/hint"
)

type AddSignatureItemJSONMarshaler struct {
	hint.BaseHinter
	Contract base.Address `json:"contract"`
	NFTIdx   uint64       `json:"nft_idx"`
}

func (it AddSignatureItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AddSignatureItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		NFTIdx:     it.nftIdx,
	})
}

type AddSignatureItemJSONUnmarshaler struct {
	Hint     hint.Hint `json:"_hint"`
	Contract string    `json:"contract"`
	NFTIdx   uint64    `json:"nft_idx"`
}

func (it *AddSignatureItem) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var u AddSignatureItemJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	if err := it.unpack(enc, u.Hint, u.Contract, u.NFTIdx); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	return nil
}
