package api

// Oker is an interface that returns if an error was present, and if so, what it was.
type Oker interface {
	Ok() bool
	Error() string
}
