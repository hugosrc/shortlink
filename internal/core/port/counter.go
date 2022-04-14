package port

// Counter is an abstraction of a service responsible for the
// ordered generation of integer numbers through incrementation.
type Counter interface {
	Inc() (int, error)
}
