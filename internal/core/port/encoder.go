package port

// Encoder is an abstraction of a service responsible for encoding bytes.
type Encoder interface {
	EncodeToString(src []byte) string
}
