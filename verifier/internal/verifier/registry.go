package verifier

import (
	"github.com/game-of-nfts/gon-toolbox/verifier/internal/chain"
)

type Registry struct {
	vs map[string]Verifier
}

func NewRegistry(r *chain.Registry) *Registry {
	vs := map[string]Verifier{
		"A1": A1Verifier{r},
	}
	return &Registry{vs}
}

func (r *Registry) Get(key string) Verifier {
	return r.vs[key]
}
