package nft

import (
	"fmt"
	"strconv"

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
	TransferFactHint = hint.MustNewHint("mitum-nft-transfer-operation-fact-v0.0.1")
	TransferHint     = hint.MustNewHint("mitum-nft-transfer-operation-v0.0.1")
)

var MaxTransferItems = 100

type TransferFact struct {
	base.BaseFact
	sender   base.Address
	items    []TransferItem
	currency types.CurrencyID
}

func NewTransferFact(
	token []byte,
	sender base.Address,
	items []TransferItem,
	currency types.CurrencyID,
) TransferFact {
	bf := base.NewBaseFact(TransferFactHint, token)

	fact := TransferFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
		currency: currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact TransferFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if l := len(fact.items); l < 1 {
		return common.ErrFactInvalid.Wrap(common.ErrArrayLen.Wrap(errors.Errorf("empty items for TransferFact")))
	} else if l > int(MaxTransferItems) {
		return common.ErrFactInvalid.Wrap(
			common.ErrArrayLen.Wrap(errors.Errorf("items over allowed, %d > %d", l, MaxTransferItems)))
	}

	if err := util.CheckIsValiders(nil, false, fact.sender, fact.currency); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	founds := map[string]struct{}{}
	for _, item := range fact.items {
		if err := item.IsValid(nil); err != nil {
			return common.ErrFactInvalid.Wrap(err)
		}

		if fact.sender.Equal(item.contract) {
			return common.ErrFactInvalid.Wrap(
				common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with contract account", fact.sender)))
		}

		n := strconv.FormatUint(item.NFT(), 10)

		if _, found := founds[item.contract.String()+"-"+n]; found {
			return common.ErrFactInvalid.Wrap(
				common.ErrDupVal.Wrap(errors.Errorf("nft idx %v in contract account %v", n, item.contract)))
		}

		founds[n] = struct{}{}
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	return nil
}

func (fact TransferFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact TransferFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact TransferFact) Bytes() []byte {
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

func (fact TransferFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact TransferFact) Sender() base.Address {
	return fact.sender
}

func (fact TransferFact) Items() []TransferItem {
	return fact.items
}

func (fact TransferFact) Currency() types.CurrencyID {
	return fact.currency
}

func (fact TransferFact) Addresses() ([]base.Address, error) {
	var as []base.Address

	for i := range fact.items {
		if ads, err := fact.items[i].Addresses(); err != nil {
			return nil, err
		} else {
			as = append(as, ads...)
		}
	}

	as = append(as, fact.Sender())

	return as, nil
}

func (fact TransferFact) FeeBase() (types.CurrencyID, int, int, bool) {
	return fact.Currency(), len(fact.items), len(fact.Bytes()), extras.HasItem
}

func (fact TransferFact) FeePayer() base.Address {
	return fact.sender
}

func (fact TransferFact) FactUser() base.Address {
	return fact.sender
}

func (fact TransferFact) Signer() base.Address {
	return fact.sender
}

func (fact TransferFact) ActiveContract() []base.Address {
	var arr []base.Address
	for i := range fact.items {
		arr = append(arr, fact.items[i].contract)
	}
	return arr
}

func (fact TransferFact) DupKey() (map[types.DuplicationKeyType][]string, error) {
	r := make(map[types.DuplicationKeyType][]string)
	for _, item := range fact.items {
		r[processor.DuplicationTypeContractNFT] = append(
			r[processor.DuplicationTypeContractNFT], fmt.Sprintf("%s:%v", item.contract.String(), item.nftIdx))
	}

	return r, nil
}

type Transfer struct {
	extras.ExtendedOperation
}

func (op Transfer) DupKey() (map[types.DuplicationKeyType][]string, error) {
	r := make(map[types.DuplicationKeyType][]string)

	if err := extras.AddOperationFeePayerDupKeys(r, op); err != nil {
		return nil, err
	}

	return r, nil
}

func NewTransfer(fact TransferFact) (Transfer, error) {
	return Transfer{
		ExtendedOperation: extras.NewExtendedOperation(TransferHint, fact),
	}, nil
}
