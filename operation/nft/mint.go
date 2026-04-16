package nft

import (
	"github.com/imfact-labs/currency-model/common"
	"github.com/imfact-labs/currency-model/operation/extras"
	"github.com/imfact-labs/currency-model/types"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/hint"
	"github.com/imfact-labs/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var MaxMintItems = 100

var (
	MintFactHint = hint.MustNewHint("mitum-nft-mint-operation-fact-v0.0.1")
	MintHint     = hint.MustNewHint("mitum-nft-mint-operation-v0.0.1")
)

type MintFact struct {
	base.BaseFact
	sender   base.Address
	items    []MintItem
	currency types.CurrencyID
}

func NewMintFact(
	token []byte,
	sender base.Address,
	items []MintItem,
	currency types.CurrencyID,
) MintFact {
	bf := base.NewBaseFact(MintFactHint, token)
	fact := MintFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
		currency: currency,
	}
	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact MintFact) IsValid(b []byte) error {
	if err := util.CheckIsValiders(nil, false,
		fact.BaseHinter,
		fact.sender,
		fact.currency,
	); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if l := len(fact.items); l < 1 {
		return common.ErrArrayLen.Wrap(errors.Errorf("empty items for MintFact"))
	} else if l > int(MaxMintItems) {
		return common.ErrFactInvalid.Wrap(common.ErrArrayLen.Wrap(errors.Errorf("items over allowed, %d > %d", l, MaxMintItems)))
	}

	for _, item := range fact.items {
		if err := item.IsValid(nil); err != nil {
			return common.ErrFactInvalid.Wrap(err)
		}
		if fact.sender.Equal(item.contract) {
			return common.ErrFactInvalid.Wrap(
				common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with contract account", fact.sender)))
		}
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	return nil
}

func (fact MintFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact MintFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact MintFact) Bytes() []byte {
	is := make([][]byte, len(fact.items))

	for i := range fact.items {
		is[i] = fact.items[i].Bytes()
	}

	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.currency.Bytes(),
		util.ConcatBytesSlice(is...),
	)
}

func (fact MintFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact MintFact) Sender() base.Address {
	return fact.sender
}

func (fact MintFact) Currency() types.CurrencyID {
	return fact.currency
}

func (fact MintFact) Addresses() ([]base.Address, error) {
	as := []base.Address{}

	for _, item := range fact.items {
		if ads, err := item.Addresses(); err != nil {
			return nil, err
		} else {
			as = append(as, ads...)
		}
	}

	as = append(as, fact.sender)

	return as, nil
}

func (fact MintFact) FeeBase() (types.CurrencyID, int, int, bool) {
	return fact.Currency(), len(fact.items), len(fact.Bytes()), extras.HasItem
}

func (fact MintFact) FeePayer() base.Address {
	return fact.sender
}

func (fact MintFact) FactUser() base.Address {
	return fact.sender
}

func (fact MintFact) Signer() base.Address {
	return fact.sender
}

func (fact MintFact) ActiveContract() []base.Address {
	var arr []base.Address
	for i := range fact.items {
		arr = append(arr, fact.items[i].contract)
	}
	return arr
}

func (fact MintFact) DupKey() (map[types.DuplicationKeyType][]string, error) {
	r := make(map[types.DuplicationKeyType][]string)
	r[extras.DuplicationKeyTypeSender] = []string{fact.sender.String()}

	dupSet := make(map[string]struct{}, len(fact.items))
	for _, item := range fact.items {
		_, found := dupSet[item.contract.String()]
		if !found {
			r[extras.DuplicationKeyTypeContractStatus] = append(
				r[extras.DuplicationKeyTypeContractStatus],
				item.contract.String(),
			)
			dupSet[item.contract.String()] = struct{}{}
		}
	}

	return r, nil
}

func (fact MintFact) Items() []MintItem {
	return fact.items
}

type Mint struct {
	extras.ExtendedOperation
}

func NewMint(fact MintFact) (Mint, error) {
	return Mint{
		ExtendedOperation: extras.NewExtendedOperation(MintHint, fact),
	}, nil
}
