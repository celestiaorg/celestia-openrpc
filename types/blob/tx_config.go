package blob

// NOTICE: In celestia-node, blob module simply imports state's TxConfig struct and aliases it.
// Here, we have to define the TxConfig struct from scratch to avoid the unnecessary dependency
// for users using celestia-openrpc in a modular way.

import (
	"encoding/json"
	"fmt"
)

const (
	// DefaultGasPrice specifies the default gas price value to be used when the user
	// wants to use the global minimal gas price, which is fetched from the celestia-app.
	DefaultGasPrice float64 = -1.0
	// gasMultiplier is used to increase gas limit in case if tx has additional cfg.
	gasMultiplier = 1.1
)

// NewSubmitOptions constructs a new SubmitOptions with the provided options.
// It starts with a DefaultGasPrice and then applies any additional
// options provided through the variadic parameter.
func NewSubmitOptions(opts ...ConfigOption) *SubmitOptions {
	options := &SubmitOptions{gasPrice: DefaultGasPrice}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

// SubmitOptions specifies additional options that will be applied to the Tx.
type SubmitOptions struct {
	// Specifies the address from the keystore that will sign transactions.
	// NOTE: Only `signerAddress` or `KeyName` should be passed.
	// signerAddress is a primary cfg. This means If both the address and the key are specified,
	// the address field will take priority.
	signerAddress string
	// Specifies the key from the keystore associated with an account that
	// will be used to sign transactions.
	// NOTE: This `Account` must be available in the `Keystore`.
	keyName string
	// gasPrice represents the amount to be paid per gas unit.
	// Negative gasPrice means user want us to use the minGasPrice
	// defined in the node.
	gasPrice float64
	// since gasPrice can be 0, it is necessary to understand that user explicitly set it.
	isGasPriceSet bool
	// 0 gas means users want us to calculate it for them.
	gas uint64
	// Specifies the account that will pay for the transaction.
	// Input format Bech32.
	feeGranterAddress string
}

func (cfg *SubmitOptions) GasPrice() float64 {
	if !cfg.isGasPriceSet {
		return DefaultGasPrice
	}
	return cfg.gasPrice
}

func (cfg *SubmitOptions) GasLimit() uint64 { return cfg.gas }

func (cfg *SubmitOptions) KeyName() string { return cfg.keyName }

func (cfg *SubmitOptions) SignerAddress() string { return cfg.signerAddress }

func (cfg *SubmitOptions) FeeGranterAddress() string { return cfg.feeGranterAddress }

type jsonTxConfig struct {
	GasPrice          float64 `json:"gas_price,omitempty"`
	IsGasPriceSet     bool    `json:"is_gas_price_set,omitempty"`
	Gas               uint64  `json:"gas,omitempty"`
	KeyName           string  `json:"key_name,omitempty"`
	SignerAddress     string  `json:"signer_address,omitempty"`
	FeeGranterAddress string  `json:"fee_granter_address,omitempty"`
}

func (cfg *SubmitOptions) MarshalJSON() ([]byte, error) {
	jsonOpts := &jsonTxConfig{
		SignerAddress:     cfg.signerAddress,
		KeyName:           cfg.keyName,
		GasPrice:          cfg.gasPrice,
		IsGasPriceSet:     cfg.isGasPriceSet,
		Gas:               cfg.gas,
		FeeGranterAddress: cfg.feeGranterAddress,
	}
	return json.Marshal(jsonOpts)
}

func (cfg *SubmitOptions) UnmarshalJSON(data []byte) error {
	var jsonOpts jsonTxConfig
	err := json.Unmarshal(data, &jsonOpts)
	if err != nil {
		return fmt.Errorf("unmarshalling TxConfig: %w", err)
	}

	cfg.keyName = jsonOpts.KeyName
	cfg.signerAddress = jsonOpts.SignerAddress
	cfg.gasPrice = jsonOpts.GasPrice
	cfg.isGasPriceSet = jsonOpts.IsGasPriceSet
	cfg.gas = jsonOpts.Gas
	cfg.feeGranterAddress = jsonOpts.FeeGranterAddress
	return nil
}

// ConfigOption is the functional option that is applied to the TxConfig instance
// to configure parameters.
type ConfigOption func(cfg *SubmitOptions)

// WithGasPrice is an option that allows to specify a GasPrice, which is needed
// to calculate the fee. In case GasPrice is not specified, the global GasPrice fetched from
// celestia-app will be used.
func WithGasPrice(gasPrice float64) ConfigOption {
	return func(cfg *SubmitOptions) {
		if gasPrice >= 0 {
			cfg.gasPrice = gasPrice
			cfg.isGasPriceSet = true
		}
	}
}

// WithGas is an option that allows to specify Gas.
// Gas will be calculated in case it wasn't specified.
func WithGas(gas uint64) ConfigOption {
	return func(cfg *SubmitOptions) {
		cfg.gas = gas
	}
}

// WithKeyName is an option that allows you to specify an KeyName, which is needed to
// sign the transaction. This key should be associated with the address and stored
// locally in the key store. Default Account will be used in case it wasn't specified.
func WithKeyName(key string) ConfigOption {
	return func(cfg *SubmitOptions) {
		cfg.keyName = key
	}
}

// WithSignerAddress is an option that allows you to specify an address, that will sign the transaction.
// This address must be stored locally in the key store. Default signerAddress will be used in case it wasn't specified.
func WithSignerAddress(address string) ConfigOption {
	return func(cfg *SubmitOptions) {
		cfg.signerAddress = address
	}
}

// WithFeeGranterAddress is an option that allows you to specify a GranterAddress to pay the fees.
func WithFeeGranterAddress(granter string) ConfigOption {
	return func(cfg *SubmitOptions) {
		cfg.feeGranterAddress = granter
	}
}
