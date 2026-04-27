package nft

import (
	"fmt"

	"github.com/imfact-labs/currency-model/common"
	"github.com/imfact-labs/currency-model/operation/extras"
	"github.com/imfact-labs/currency-model/types"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/hint"
	"github.com/imfact-labs/mitum2/util/valuehash"
	"github.com/imfact-labs/nft-model/operation/processor"
	"github.com/pkg/errors"
)

var (
	ApproveAllFactHint = hint.MustNewHint("mitum-nft-approve-all-operation-fact-v0.0.1")
	ApproveAllHint     = hint.MustNewHint("mitum-nft-approve-all-operation-v0.0.1")
)

var MaxApproveAllItems = 100

type ApproveAllFact struct {
	base.BaseFact
	sender   base.Address
	items    []ApproveAllItem
	currency types.CurrencyID
}

func NewApproveAllFact(
	token []byte,
	sender base.Address,
	items []ApproveAllItem,
	currency types.CurrencyID,
) ApproveAllFact {
	bf := base.NewBaseFact(ApproveAllFactHint, token)
	fact := ApproveAllFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
		currency: currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact ApproveAllFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if l := len(fact.items); l < 1 {
		return common.ErrFactInvalid.Wrap(errors.Errorf("empty items for DelegateFact"))
	} else if l > int(MaxApproveAllItems) {
		return common.ErrFactInvalid.Wrap(
			common.ErrValOOR.Wrap(errors.Errorf("items over allowed, %d > %d", l, MaxApproveAllItems)))
	}

	if err := util.CheckIsValiders(nil, false, fact.sender, fact.currency); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	founds := map[string]map[string]struct{}{}
	for _, item := range fact.items {
		if err := item.IsValid(nil); err != nil {
			return common.ErrFactInvalid.Wrap(err)
		}

		if fact.sender.Equal(item.contract) {
			return common.ErrFactInvalid.Wrap(common.ErrSelfTarget.Wrap(
				errors.Errorf("sender %v is same with contract account", fact.sender)))
		}

		if fact.sender.Equal(item.approved) {
			return common.ErrFactInvalid.Wrap(common.ErrSelfTarget.Wrap(
				errors.Errorf("sender %v delegates to itself", fact.sender)))
		}

		if addressMap, collectionFound := founds[item.contract.String()]; !collectionFound {
			founds[item.contract.String()] = make(map[string]struct{})
		} else if _, addressFound := addressMap[item.Approved().String()]; addressFound {
			return common.ErrFactInvalid.Wrap(
				common.ErrDupVal.Wrap(
					errors.Errorf("approved %v in contract account %v", item.Approved(), item.contract)))
		}

		founds[item.contract.String()][item.Approved().String()] = struct{}{}
	}
	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	return nil
}

func (fact ApproveAllFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact ApproveAllFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact ApproveAllFact) Bytes() []byte {
	is := make([][]byte, len(fact.items))
	for i, item := range fact.items {
		is[i] = item.Bytes()
	}

	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.currency.Bytes(),
		util.ConcatBytesSlice(is...),
	)
}

func (fact ApproveAllFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact ApproveAllFact) Sender() base.Address {
	return fact.sender
}

func (fact ApproveAllFact) Currency() types.CurrencyID {
	return fact.currency
}

func (fact ApproveAllFact) Addresses() ([]base.Address, error) {
	l := len(fact.items)

	as := make([]base.Address, l+1)

	for i, item := range fact.items {
		as[i] = item.Approved()
	}

	as[l] = fact.sender

	return as, nil
}

func (fact ApproveAllFact) FeeBase() (types.CurrencyID, int, int, bool) {
	return fact.Currency(), len(fact.items), len(fact.Bytes()), extras.HasItem
}

func (fact ApproveAllFact) FeePayer() base.Address {
	return fact.sender
}

func (fact ApproveAllFact) FactUser() base.Address {
	return fact.sender
}

func (fact ApproveAllFact) Signer() base.Address {
	return fact.sender
}

func (fact ApproveAllFact) ActiveContract() []base.Address {
	var arr []base.Address
	for i := range fact.items {
		arr = append(arr, fact.items[i].contract)
	}
	return arr
}

func (fact ApproveAllFact) DupKey() (map[types.DuplicationKeyType][]string, error) {
	r := make(map[types.DuplicationKeyType][]string)
	dupSet := make(map[string]struct{}, len(fact.items))
	for _, item := range fact.items {
		key := fmt.Sprintf("%s:%s", item.contract.String(), fact.sender.String())
		_, found := dupSet[key]
		if !found {
			r[processor.DuplicationTypeNFTApprove] = append(
				r[processor.DuplicationTypeNFTApprove],
				key,
			)
			dupSet[key] = struct{}{}
		}
	}

	return r, nil
}

func (fact ApproveAllFact) Items() []ApproveAllItem {
	return fact.items
}

type ApproveAll struct {
	extras.ExtendedOperation
}

func (op ApproveAll) DupKey() (map[types.DuplicationKeyType][]string, error) {
	r := make(map[types.DuplicationKeyType][]string)

	if err := extras.AddOperationFeePayerDupKeys(r, op); err != nil {
		return nil, err
	}

	return r, nil
}

func NewDelegate(fact ApproveAllFact) (ApproveAll, error) {
	return ApproveAll{
		ExtendedOperation: extras.NewExtendedOperation(ApproveAllHint, fact),
	}, nil
}
