package node

// Info contains information related to the administrative
// node.
type Info struct {
	Type       Type   `json:"type"`
	APIVersion string `json:"api_version"`
}

// Type defines the Node type (e.g. `light`, `bridge`) for identity purposes.
// The zero value for Type is invalid.
type Type uint8
