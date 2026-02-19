package types

import (
	"github.com/imfact-labs/currency-model/utils/bsonenc"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (sgns Signers) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":   sgns.Hint().String(),
			"signers": sgns.signers,
		})
}

type SignersBSONUnmarshaler struct {
	Hint    string   `bson:"_hint"`
	Signers bson.Raw `bson:"signers"`
}

func (sgns *Signers) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of Signers")

	var u SignersBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return sgns.unpack(enc, ht, u.Signers)
}
