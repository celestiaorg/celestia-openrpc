package header

import (
	"encoding/json"
	"time"

	"github.com/celestiaorg/go-header"
	cmjson "github.com/cometbft/cometbft/libs/json"

	"github.com/celestiaorg/celestia-openrpc/types/core"
)

// RawHeader is an alias to core.Header. It is
// "raw" because it is not yet wrapped to include
// the DataAvailabilityHeader.
type RawHeader = core.Header

// ExtendedHeader represents a wrapped "raw" header that includes
// information necessary for Celestia Nodes to be notified of new
// block headers and perform Data Availability Sampling.
type ExtendedHeader struct {
	RawHeader    `json:"header"`
	Commit       *core.Commit            `json:"commit"`
	ValidatorSet *core.ValidatorSet      `json:"validator_set"`
	DAH          *DataAvailabilityHeader `json:"dah"`
}

// MarshalJSON marshals an ExtendedHeader to JSON. The ValidatorSet is wrapped with amino encoding,
// to be able to unmarshal the crypto.PubKey type back from JSON.
func (eh *ExtendedHeader) MarshalJSON() ([]byte, error) {
	type Alias ExtendedHeader
	validatorSet, err := cmjson.Marshal(eh.ValidatorSet)
	if err != nil {
		return nil, err
	}
	rawHeader, err := cmjson.Marshal(eh.RawHeader)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&struct {
		RawHeader    json.RawMessage `json:"header"`
		ValidatorSet json.RawMessage `json:"validator_set"`
		*Alias
	}{
		ValidatorSet: validatorSet,
		RawHeader:    rawHeader,
		Alias:        (*Alias)(eh),
	})
}

// UnmarshalJSON unmarshals an ExtendedHeader from JSON. The ValidatorSet is wrapped with amino
// encoding, to be able to unmarshal the crypto.PubKey type back from JSON.
func (eh *ExtendedHeader) UnmarshalJSON(data []byte) error {
	type Alias ExtendedHeader
	aux := &struct {
		RawHeader    json.RawMessage `json:"header"`
		ValidatorSet json.RawMessage `json:"validator_set"`
		*Alias
	}{
		Alias: (*Alias)(eh),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	valSet := new(core.ValidatorSet)
	if err := cmjson.Unmarshal(aux.ValidatorSet, valSet); err != nil {
		return err
	}
	rawHeader := new(RawHeader)
	if err := cmjson.Unmarshal(aux.RawHeader, rawHeader); err != nil {
		return err
	}

	eh.ValidatorSet = valSet
	eh.RawHeader = *rawHeader
	return nil
}

type DataAvailabilityHeader = core.DataAvailabilityHeader

func (eh *ExtendedHeader) New() *ExtendedHeader {
	return new(ExtendedHeader)
}

func (eh *ExtendedHeader) IsZero() bool {
	return eh == nil
}

func (eh *ExtendedHeader) ChainID() string {
	return eh.RawHeader.ChainID
}

func (eh *ExtendedHeader) Hash() header.Hash {
	return eh.Commit.BlockID.Hash.Bytes()
}

func (eh *ExtendedHeader) Height() uint64 {
	return uint64(eh.Commit.Height)
}

func (eh *ExtendedHeader) LastHeader() header.Hash {
	return eh.RawHeader.LastBlockID.Hash.Bytes()
}

func (eh *ExtendedHeader) Time() time.Time {
	return eh.RawHeader.Time
}

func (eh *ExtendedHeader) Verify(h *ExtendedHeader) error {
	panic("implement me if being used")
}

func (eh *ExtendedHeader) Validate() error {
	panic("implement me if being used")
}

func (eh *ExtendedHeader) MarshalBinary() (data []byte, err error) {
	panic("implement me if being used")
}

func (eh *ExtendedHeader) UnmarshalBinary(data []byte) error {
	panic("implement me if being used")
}
