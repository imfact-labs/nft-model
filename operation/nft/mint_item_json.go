package nft

import (
	"encoding/json"

	"github.com/imfact-labs/currency-model/common"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/encoder"
	"github.com/imfact-labs/mitum2/util/hint"
	"github.com/imfact-labs/nft-model/types"
)

type MintItemJSONMarshaler struct {
	hint.BaseHinter
	Contract base.Address  `json:"contract"`
	Receiver base.Address  `json:"receiver"`
	Hash     types.NFTHash `json:"hash"`
	Uri      types.URI     `json:"uri"`
	Creators types.Signers `json:"creators"`
}

func (it MintItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(MintItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		Receiver:   it.receiver,
		Hash:       it.hash,
		Uri:        it.uri,
		Creators:   it.creators,
	})
}

type MintItemJSONUnmarshaler struct {
	Hint     hint.Hint       `json:"_hint"`
	Contract string          `json:"contract"`
	Receiver string          `json:"receiver"`
	Hash     string          `json:"hash"`
	Uri      string          `json:"uri"`
	Creators json.RawMessage `json:"creators"`
}

func (it *MintItem) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var u MintItemJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	if err := it.unpack(enc, u.Hint, u.Contract, u.Receiver, u.Hash, u.Uri, u.Creators); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	return nil
}
