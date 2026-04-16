package nft

import (
	"github.com/imfact-labs/currency-model/common"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/hint"
	"github.com/pkg/errors"
)

var ApproveItemHint = hint.MustNewHint("mitum-nft-approve-item-v0.0.1")

type ApproveItem struct {
	hint.BaseHinter
	contract base.Address
	approved base.Address
	nftIdx   uint64
}

func NewApproveItem(contract base.Address, approved base.Address, nftIdx uint64) ApproveItem {
	return ApproveItem{
		BaseHinter: hint.NewBaseHinter(ApproveItemHint),
		contract:   contract,
		approved:   approved,
		nftIdx:     nftIdx,
	}
}

func (it ApproveItem) IsValid([]byte) error {
	if it.approved.Equal(it.contract) {
		return common.ErrSelfTarget.Wrap(errors.Errorf("approved %v is same with contract contract", it.approved))
	}

	return util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
		it.approved,
	)
}

func (it ApproveItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.approved.Bytes(),
		util.Uint64ToBytes(it.nftIdx),
	)
}

func (it ApproveItem) Contract() base.Address {
	return it.contract
}

func (it ApproveItem) Approved() base.Address {
	return it.approved
}

func (it ApproveItem) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = it.approved
	return as, nil
}

func (it ApproveItem) NFTIdx() uint64 {
	return it.nftIdx
}
