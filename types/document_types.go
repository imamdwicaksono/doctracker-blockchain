package types

type Document struct {
	ID     string
	Title  string
	Hash   [32]byte
	Status string
}
