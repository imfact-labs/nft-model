package nft

import (
	"github.com/imfact-labs/currency-model/common"
	"github.com/imfact-labs/currency-model/utils/bsonenc"
	"github.com/imfact-labs/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (it ApproveAllItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":    it.Hint().String(),
			"contract": it.contract,
			"approved": it.approved,
			"mode":     it.mode,
			"currency": it.currency,
		},
	)
}

type DelegateItemBSONUnmarshaler struct {
	Hint     string `bson:"_hint"`
	Contract string `bson:"contract"`
	Approved string `bson:"approved"`
	Mode     string `bson:"mode"`
	Currency string `bson:"currency"`
}

func (it *ApproveAllItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	var u DelegateItemBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}

	if err := it.unmarshal(enc, ht, u.Contract, u.Approved, u.Mode, u.Currency); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}
	return nil
}
