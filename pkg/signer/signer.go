package signer

type Signer interface {
}

type signer struct {
}

func New() Signer {
	return &signer{}
}
