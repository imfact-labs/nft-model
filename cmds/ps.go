package cmds

import (
	"context"

	ccmds "github.com/imfact-labs/currency-model/app/cmds"
	cprocessor "github.com/imfact-labs/currency-model/operation/processor"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/isaac"
	"github.com/imfact-labs/mitum2/launch"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/hint"
	"github.com/imfact-labs/mitum2/util/ps"
	"github.com/imfact-labs/nft-model/operation/nft"
)

var PNameOperationProcessorsMap = ps.Name("mitum-nft-operation-processors-map")

func POperationProcessorsMap(pctx context.Context) (context.Context, error) {
	var isaacParams *isaac.Params
	var db isaac.Database
	var opr *cprocessor.OperationProcessor
	var set *hint.CompatibleSet[isaac.NewOperationProcessorInternalFunc]

	if err := util.LoadFromContextOK(pctx,
		launch.ISAACParamsContextKey, &isaacParams,
		launch.CenterDatabaseContextKey, &db,
		ccmds.OperationProcessorContextKey, &opr,
		launch.OperationProcessorsMapContextKey, &set,
	); err != nil {
		return pctx, err
	}

	//err := opr.SetCheckDuplicationFunc(processor.CheckDuplication)
	//if err != nil {
	//	return pctx, err
	//}
	err := opr.SetGetNewProcessorFunc(cprocessor.GetNewProcessor)
	if err != nil {
		return pctx, err
	}
	if err := opr.SetProcessor(
		nft.RegisterModelHint,
		nft.NewRegisterModelProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		nft.UpdateModelConfigHint,
		nft.NewUpdateModelConfigProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		nft.MintHint,
		nft.NewMintProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		nft.TransferHint,
		nft.NewTransferProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		nft.ApproveAllHint,
		nft.NewDelegateProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		nft.ApproveHint,
		nft.NewApproveProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		nft.AddSignatureHint,
		nft.NewSignProcessor(),
	); err != nil {
		return pctx, err
	}

	_ = set.Add(nft.RegisterModelHint,
		func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
			return opr.New(
				height,
				getStatef,
				nil,
				nil,
			)
		})

	_ = set.Add(nft.UpdateModelConfigHint,
		func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
			return opr.New(
				height,
				getStatef,
				nil,
				nil,
			)
		})

	_ = set.Add(nft.MintHint,
		func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
			return opr.New(
				height,
				getStatef,
				nil,
				nil,
			)
		})

	_ = set.Add(nft.TransferHint,
		func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
			return opr.New(
				height,
				getStatef,
				nil,
				nil,
			)
		})

	_ = set.Add(nft.ApproveAllHint,
		func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
			return opr.New(
				height,
				getStatef,
				nil,
				nil,
			)
		})

	_ = set.Add(nft.ApproveHint,
		func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
			return opr.New(
				height,
				getStatef,
				nil,
				nil,
			)
		})

	_ = set.Add(nft.AddSignatureHint,
		func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
			return opr.New(
				height,
				getStatef,
				nil,
				nil,
			)
		})

	pctx = context.WithValue(pctx, ccmds.OperationProcessorContextKey, opr)
	pctx = context.WithValue(pctx, launch.OperationProcessorsMapContextKey, set) //revive:disable-line:modifies-parameter

	return pctx, nil
}
