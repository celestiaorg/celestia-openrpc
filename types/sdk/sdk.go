package sdk

import (
	"fmt"

	"cosmossdk.io/math"
	"github.com/gogo/protobuf/types"
)

// Address is a common interface for different types of addresses used by the SDK
type Address interface {
	Equals(Address) bool
	Empty() bool
	Marshal() ([]byte, error)
	MarshalJSON() ([]byte, error)
	Bytes() []byte
	String() string
	Format(s fmt.State, verb rune)
}

// AccAddress a wrapper around bytes meant to represent an account address.
// When marshaled to a string or JSON, it uses Bech32.
type AccAddress []byte

// ValAddress defines a wrapper around bytes meant to present a validator's
// operator. When marshaled to a string or JSON, it uses Bech32.
type ValAddress []byte

// Coin defines a token with a denomination and an amount.
//
// NOTE: The amount field is an Int which implements the custom method
// signatures required by gogoproto.
type Coin struct {
	Denom  string   `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	Amount math.Int `protobuf:"bytes,2,opt,name=amount,proto3,customtype=Int" json:"amount"`
}

// TxResponse defines a structure containing relevant tx data and metadata. The
// tags are stringified and the log is JSON decoded.
type TxResponse struct {
	// The block height
	Height int64 `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	// The transaction hash.
	TxHash string `protobuf:"bytes,2,opt,name=txhash,proto3" json:"txhash,omitempty"`
	// Namespace for the Code
	Codespace string `protobuf:"bytes,3,opt,name=codespace,proto3" json:"codespace,omitempty"`
	// Response code.
	Code uint32 `protobuf:"varint,4,opt,name=code,proto3" json:"code,omitempty"`
	// Result bytes, if any.
	Data string `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
	// The output of the application's logger (raw string). May be
	// non-deterministic.
	RawLog string `protobuf:"bytes,6,opt,name=raw_log,json=rawLog,proto3" json:"raw_log,omitempty"`
	// The output of the application's logger (typed). May be non-deterministic.
	Logs ABCIMessageLogs `protobuf:"bytes,7,rep,name=logs,proto3,castrepeated=ABCIMessageLogs" json:"logs"`
	// Additional information. May be non-deterministic.
	Info string `protobuf:"bytes,8,opt,name=info,proto3" json:"info,omitempty"`
	// Amount of gas requested for transaction.
	GasWanted int64 `protobuf:"varint,9,opt,name=gas_wanted,json=gasWanted,proto3" json:"gas_wanted,omitempty"`
	// Amount of gas consumed by transaction.
	GasUsed int64 `protobuf:"varint,10,opt,name=gas_used,json=gasUsed,proto3" json:"gas_used,omitempty"`
	// The request transaction bytes.
	Tx *types.Any `protobuf:"bytes,11,opt,name=tx,proto3" json:"tx,omitempty"`
	// Time of the previous block. For heights > 1, it's the weighted median of
	// the timestamps of the valid votes in the block.LastCommit. For height == 1,
	// it's genesis time.
	Timestamp string `protobuf:"bytes,12,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// Events defines all the events emitted by processing a transaction. Note,
	// these events include those emitted by processing all the messages and those
	// emitted from the ante. Whereas Logs contains the events, with
	// additional metadata, emitted only by processing the messages.
	//
	// Since: cosmos-sdk 0.42.11, 0.44.5, 0.45
	Events []Event `protobuf:"bytes,13,rep,name=events,proto3" json:"events"`
}

// ABCIMessageLogs represents a slice of ABCIMessageLog.
type ABCIMessageLogs []ABCIMessageLog

// ABCIMessageLog defines a structure containing an indexed tx ABCI message log.
type ABCIMessageLog struct {
	MsgIndex uint32 `protobuf:"varint,1,opt,name=msg_index,json=msgIndex,proto3" json:"msg_index"`
	Log      string `protobuf:"bytes,2,opt,name=log,proto3" json:"log,omitempty"`
	// Events contains a slice of Event objects that were emitted during some
	// execution.
	Events StringEvents `protobuf:"bytes,3,rep,name=events,proto3,castrepeated=StringEvents" json:"events"`
}

// StringAttributes defines a slice of StringEvents objects.
type StringEvents []StringEvent

// StringEvent defines en Event object wrapper where all the attributes
// contain key/value pairs that are strings instead of raw bytes.
type StringEvent struct {
	Type       string      `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Attributes []Attribute `protobuf:"bytes,2,rep,name=attributes,proto3" json:"attributes"`
}

// Attribute defines an attribute wrapper where the key and value are
// strings instead of raw bytes.
type Attribute struct {
	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

// Event allows application developers to attach additional information to
// ResponseBeginBlock, ResponseEndBlock, ResponseCheckTx and ResponseDeliverTx.
// Later, transactions may be queried using these events.
type Event struct {
	Type       string           `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Attributes []EventAttribute `protobuf:"bytes,2,rep,name=attributes,proto3" json:"attributes,omitempty"`
}

// EventAttribute is a single key-value pair, associated with an event.
type EventAttribute struct {
	Key   []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Index bool   `protobuf:"varint,3,opt,name=index,proto3" json:"index,omitempty"`
}
