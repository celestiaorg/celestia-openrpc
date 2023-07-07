package state

import (
	"math/big"
	"time"

	"cosmossdk.io/math"

	coretypes "github.com/cometbft/cometbft/types"

	"github.com/rollkit/celestia-openrpc/types/sdk"
)

// Balance is an alias to the Coin type from Cosmos-SDK.
type Balance = sdk.Coin

// Tx is an alias to the Tx type from celestia-core.
type Tx = coretypes.Tx

// TxResponse is an alias to the TxResponse type from Cosmos-SDK.
type TxResponse = sdk.TxResponse

// Address is an alias to the Address type from Cosmos-SDK.
type Address = sdk.Address

// ValAddress is an alias to the ValAddress type from Cosmos-SDK.
type ValAddress = sdk.ValAddress

// AccAddress is an alias to the AccAddress type from Cosmos-SDK.
type AccAddress = sdk.AccAddress

// Int is an alias to the Int type from Cosmos-SDK.
type Int = math.Int

// QueryDelegationResponse is response type for the Query/Delegation RPC method.
type QueryDelegationResponse struct {
	// delegation_responses defines the delegation info of a delegation.
	DelegationResponse *DelegationResponse `protobuf:"bytes,1,opt,name=delegation_response,json=delegationResponse,proto3" json:"delegation_response,omitempty"`
}

// DelegationResponse is equivalent to Delegation except that it contains a
// balance in addition to shares which is more suitable for client responses.
type DelegationResponse struct {
	Delegation Delegation `protobuf:"bytes,1,opt,name=delegation,proto3" json:"delegation"`
	Balance    sdk.Coin   `protobuf:"bytes,2,opt,name=balance,proto3" json:"balance"`
}

// Delegation represents the bond with tokens held by an account. It is
// owned by one delegator, and is associated with the voting power of one
// validator.
type Delegation struct {
	// delegator_address is the bech32-encoded address of the delegator.
	DelegatorAddress string `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
	// validator_address is the bech32-encoded address of the validator.
	ValidatorAddress string `protobuf:"bytes,2,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
	// shares define the delegation shares received.
	Shares Dec `protobuf:"bytes,3,opt,name=shares,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"shares"`
}

// NOTE: never use new(Dec) or else we will panic unmarshalling into the
// nil embedded big.Int
type Dec struct {
	// TODO(tzdybal) - this type needs more attention!
	i *big.Int //nolint:unused,structcheck
}

// QueryDelegationResponse is response type for the Query/UnbondingDelegation
// RPC method.
type QueryUnbondingDelegationResponse struct {
	// unbond defines the unbonding information of a delegation.
	Unbond UnbondingDelegation `protobuf:"bytes,1,opt,name=unbond,proto3" json:"unbond"`
}

// UnbondingDelegation stores all of a single delegator's unbonding bonds
// for a single validator in an time-ordered list.
type UnbondingDelegation struct {
	// delegator_address is the bech32-encoded address of the delegator.
	DelegatorAddress string `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
	// validator_address is the bech32-encoded address of the validator.
	ValidatorAddress string `protobuf:"bytes,2,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
	// entries are the unbonding delegation entries.
	Entries []UnbondingDelegationEntry `protobuf:"bytes,3,rep,name=entries,proto3" json:"entries"`
}

// UnbondingDelegationEntry defines an unbonding object with relevant metadata.
type UnbondingDelegationEntry struct {
	// creation_height is the height which the unbonding took place.
	CreationHeight int64 `protobuf:"varint,1,opt,name=creation_height,json=creationHeight,proto3" json:"creation_height,omitempty"`
	// completion_time is the unix time for unbonding completion.
	CompletionTime time.Time `protobuf:"bytes,2,opt,name=completion_time,json=completionTime,proto3,stdtime" json:"completion_time"`
	// initial_balance defines the tokens initially scheduled to receive at completion.
	InitialBalance Int `protobuf:"bytes,3,opt,name=initial_balance,json=initialBalance,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"initial_balance"`
	// balance defines the tokens to receive at completion.
	Balance Int `protobuf:"bytes,4,opt,name=balance,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"balance"`
}

// QueryRedelegationsResponse is response type for the Query/Redelegations RPC
// method.
type QueryRedelegationsResponse struct {
	RedelegationResponses []RedelegationResponse `protobuf:"bytes,1,rep,name=redelegation_responses,json=redelegationResponses,proto3" json:"redelegation_responses"`
	// pagination defines the pagination in the response.
	Pagination *PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

// RedelegationResponse is equivalent to a Redelegation except that its entries
// contain a balance in addition to shares which is more suitable for client
// responses.
type RedelegationResponse struct {
	Redelegation Redelegation                `protobuf:"bytes,1,opt,name=redelegation,proto3" json:"redelegation"`
	Entries      []RedelegationEntryResponse `protobuf:"bytes,2,rep,name=entries,proto3" json:"entries"`
}

// PageResponse is to be embedded in gRPC response messages where the
// corresponding request message has used PageRequest.
//
//	message SomeResponse {
//	        repeated Bar results = 1;
//	        PageResponse page = 2;
//	}
type PageResponse struct {
	// next_key is the key to be passed to PageRequest.key to
	// query the next page most efficiently. It will be empty if
	// there are no more results.
	NextKey []byte `protobuf:"bytes,1,opt,name=next_key,json=nextKey,proto3" json:"next_key,omitempty"`
	// total is total number of results available if PageRequest.count_total
	// was set, its value is undefined otherwise
	Total uint64 `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
}

// Redelegation contains the list of a particular delegator's redelegating bonds
// from a particular source validator to a particular destination validator.
type Redelegation struct {
	// delegator_address is the bech32-encoded address of the delegator.
	DelegatorAddress string `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
	// validator_src_address is the validator redelegation source operator address.
	ValidatorSrcAddress string `protobuf:"bytes,2,opt,name=validator_src_address,json=validatorSrcAddress,proto3" json:"validator_src_address,omitempty"`
	// validator_dst_address is the validator redelegation destination operator address.
	ValidatorDstAddress string `protobuf:"bytes,3,opt,name=validator_dst_address,json=validatorDstAddress,proto3" json:"validator_dst_address,omitempty"`
	// entries are the redelegation entries.
	Entries []RedelegationEntry `protobuf:"bytes,4,rep,name=entries,proto3" json:"entries"`
}

// RedelegationEntry defines a redelegation object with relevant metadata.
type RedelegationEntry struct {
	// creation_height  defines the height which the redelegation took place.
	CreationHeight int64 `protobuf:"varint,1,opt,name=creation_height,json=creationHeight,proto3" json:"creation_height,omitempty"`
	// completion_time defines the unix time for redelegation completion.
	CompletionTime time.Time `protobuf:"bytes,2,opt,name=completion_time,json=completionTime,proto3,stdtime" json:"completion_time"`
	// initial_balance defines the initial balance when redelegation started.
	InitialBalance Int `protobuf:"bytes,3,opt,name=initial_balance,json=initialBalance,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"initial_balance"`
	// shares_dst is the amount of destination-validator shares created by redelegation.
	SharesDst Dec `protobuf:"bytes,4,opt,name=shares_dst,json=sharesDst,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"shares_dst"`
}

// RedelegationEntryResponse is equivalent to a RedelegationEntry except that it
// contains a balance in addition to shares which is more suitable for client
// responses.
type RedelegationEntryResponse struct {
	RedelegationEntry RedelegationEntry `protobuf:"bytes,1,opt,name=redelegation_entry,json=redelegationEntry,proto3" json:"redelegation_entry"`
	Balance           Int               `protobuf:"bytes,4,opt,name=balance,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"balance"`
}
