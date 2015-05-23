package selector

import (
	"golang.org/x/net/context"
	"gopkg.in/cogger/cogger.v1"
	"gopkg.in/cogger/cogger.v1/cogs"
	"gopkg.in/cogger/cogger.v1/limiter"
)

type Selector interface {
	cogger.Cog
	Case(func(context.Context) bool, cogger.Cog) Selector
	Default(cogger.Cog) Selector
}

func New() Selector {
	return &defaultSelector{
		defaultCog: cogs.NoOp(),
	}
}

type defaultSelector struct {
	defaultCog cogger.Cog
	limit      limiter.Limit
	cases      []selectCase
}

type selectCase struct {
	test func(context.Context) bool
	cog  cogger.Cog
}

func (ds *defaultSelector) Case(test func(context.Context) bool, cog cogger.Cog) Selector {
	ds.cases = append(ds.cases, selectCase{
		test: test,
		cog:  cog,
	})
	return ds
}

func (ds *defaultSelector) Default(cog cogger.Cog) Selector {
	ds.defaultCog = cog
	return ds
}

func (ds *defaultSelector) Do(ctx context.Context) chan error {

	cog := ds.defaultCog
	for _, c := range ds.cases {
		if c.test(ctx) {
			cog = c.cog
			break
		}
	}

	if ds.limit != nil {
		cog = cog.SetLimit(ds.limit)
	}

	return cog.Do(ctx)
}

func (ds *defaultSelector) SetLimit(lim limiter.Limit) cogger.Cog {
	ds.limit = lim
	return ds
}
