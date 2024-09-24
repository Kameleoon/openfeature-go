package kameleoon

import (
	"context"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/stretchr/testify/mock"
)

type MockKameleoonResolver struct {
	mock.Mock
}

func (m *MockKameleoonResolver) Resolve(
	ctx context.Context, flagKey string, defaultValue interface{}, evalCtx openfeature.FlattenedContext,
) (interface{}, *openfeature.ResolutionError, string) {
	args := m.Called(ctx, flagKey, defaultValue, evalCtx)
	return args.Get(0).(interface{}), args.Get(1).(*openfeature.ResolutionError), args.String(2)
}
