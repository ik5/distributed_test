package types

// RawAction convert a payload to a known structure
type RawAction struct {
	Length int
	Type   ActionType
	Body   []byte
}
