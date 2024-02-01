package state

import (
	"context"

	"github.com/celestiaorg/celestia-openrpc/types/blob"
)

type API struct {
	AccountAddress    func(ctx context.Context) (Address, error)                `perm:"read"`
	IsStopped         func(ctx context.Context) bool                            `perm:"read"`
	Balance           func(ctx context.Context) (*Balance, error)               `perm:"read"`
	BalanceForAddress func(ctx context.Context, addr Address) (*Balance, error) `perm:"read"`
	Transfer          func(
		ctx context.Context,
		to AccAddress,
		amount,
		fee Int,
		gasLimit uint64,
	) (*TxResponse, error) `perm:"write"`
	SubmitTx         func(ctx context.Context, tx Tx) (*TxResponse, error) `perm:"write"`
	SubmitPayForBlob func(
		ctx context.Context,
		fee Int,
		gasLim uint64,
		blobs []*blob.Blob,
	) (*TxResponse, error) `perm:"write"`
	CancelUnbondingDelegation func(
		ctx context.Context,
		valAddr ValAddress,
		amount,
		height,
		fee Int,
		gasLim uint64,
	) (*TxResponse, error) `perm:"write"`
	BeginRedelegate func(
		ctx context.Context,
		srcValAddr,
		dstValAddr ValAddress,
		amount,
		fee Int,
		gasLim uint64,
	) (*TxResponse, error) `perm:"write"`
	Undelegate func(
		ctx context.Context,
		delAddr ValAddress,
		amount,
		fee Int,
		gasLim uint64,
	) (*TxResponse, error) `perm:"write"`
	Delegate func(
		ctx context.Context,
		delAddr ValAddress,
		amount,
		fee Int,
		gasLim uint64,
	) (*TxResponse, error) `perm:"write"`
	QueryDelegation func(
		ctx context.Context,
		valAddr ValAddress,
	) (*QueryDelegationResponse, error) `perm:"read"`
	QueryUnbonding func(
		ctx context.Context,
		valAddr ValAddress,
	) (*QueryUnbondingDelegationResponse, error) `perm:"read"`
	QueryRedelegations func(
		ctx context.Context,
		srcValAddr,
		dstValAddr ValAddress,
	) (*QueryRedelegationsResponse, error) `perm:"read"`
}
