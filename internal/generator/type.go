package generator

import "context"

type Generator interface {
	Generate(ctx context.Context) error
}

type empty struct{}
