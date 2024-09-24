package kameleoon

import (
	"context"
	kameleoon "github.com/Kameleoon/client-go/v3"
	"github.com/open-feature/go-sdk/openfeature"
)

const META_NAME = "Kameleoon Provider"

type kameleoonProvider struct {
	siteCode string
	client   kameleoon.KameleoonClient
	resolver resolver
}

// NewKameleoonProvider creates a new instance of kameleoonProvider with the given siteCode and client configuration.
func NewKameleoonProvider(siteCode string, config *kameleoon.KameleoonClientConfig) (*kameleoonProvider, error) {
	client, err := kameleoon.KameleoonClientFactory.Create(siteCode, config)
	if err != nil {
		return nil, openfeature.NewProviderNotReadyResolutionError(err.Error())
	}
	return &kameleoonProvider{
		siteCode: siteCode,
		client:   client,
		resolver: newKameleoonResolver(client),
	}, nil
}

// Metadata returns the metadata of the provider.
func (p *kameleoonProvider) Metadata() openfeature.Metadata {
	return openfeature.Metadata{
		Name: META_NAME,
	}
}

// BooleanEvaluation returns a boolean flag
func (p *kameleoonProvider) BooleanEvaluation(
	ctx context.Context, flag string, defaultValue bool, evalCtx openfeature.FlattenedContext,
) openfeature.BoolResolutionDetail {
	result, err, variant := p.resolver.Resolve(ctx, flag, defaultValue, evalCtx)
	providerResDetail := createProviderResolutionDetail(err, variant)
	boolResult, _ := result.(bool)
	return openfeature.BoolResolutionDetail{
		Value:                    boolResult,
		ProviderResolutionDetail: providerResDetail,
	}
}

// StringEvaluation returns a string flag
func (p *kameleoonProvider) StringEvaluation(
	ctx context.Context, flag string, defaultValue string, evalCtx openfeature.FlattenedContext,
) openfeature.StringResolutionDetail {
	result, err, variant := p.resolver.Resolve(ctx, flag, defaultValue, evalCtx)
	providerResDetail := createProviderResolutionDetail(err, variant)
	stringResult, _ := result.(string)
	return openfeature.StringResolutionDetail{
		Value:                    stringResult,
		ProviderResolutionDetail: providerResDetail,
	}
}

// FloatEvaluation returns a float flag
func (p *kameleoonProvider) FloatEvaluation(
	ctx context.Context, flag string, defaultValue float64, evalCtx openfeature.FlattenedContext,
) openfeature.FloatResolutionDetail {
	result, err, variant := p.resolver.Resolve(ctx, flag, defaultValue, evalCtx)
	providerResDetail := createProviderResolutionDetail(err, variant)
	floatResult, _ := result.(float64)
	return openfeature.FloatResolutionDetail{
		Value:                    floatResult,
		ProviderResolutionDetail: providerResDetail,
	}
}

// IntEvaluation returns an int flag
func (p *kameleoonProvider) IntEvaluation(
	ctx context.Context, flag string, defaultValue int64, evalCtx openfeature.FlattenedContext,
) openfeature.IntResolutionDetail {
	result, err, variant := p.resolver.Resolve(ctx, flag, defaultValue, evalCtx)
	providerResDetail := createProviderResolutionDetail(err, variant)
	intResult, _ := result.(int64)
	return openfeature.IntResolutionDetail{
		Value:                    intResult,
		ProviderResolutionDetail: providerResDetail,
	}
}

// ObjectEvaluation returns an object flag
func (p *kameleoonProvider) ObjectEvaluation(ctx context.Context, flag string, defaultValue interface{},
	evalCtx openfeature.FlattenedContext) openfeature.InterfaceResolutionDetail {
	result, err, variant := p.resolver.Resolve(ctx, flag, defaultValue, evalCtx)
	providerResDetail := createProviderResolutionDetail(err, variant)
	return openfeature.InterfaceResolutionDetail{
		Value:                    result,
		ProviderResolutionDetail: providerResDetail,
	}
}

// Init initializes the provider.
func (p *kameleoonProvider) Init(evaluationContext openfeature.EvaluationContext) error {
	err := p.client.WaitInit()
	return err
}

// Shutdown stops the client.
func (p *kameleoonProvider) Shutdown() {
	kameleoon.KameleoonClientFactory.Forget(p.siteCode)
}

// Status returns the current state of the provider.
func (p *kameleoonProvider) Status() openfeature.State {
	if err := p.Init(openfeature.EvaluationContext{}); err != nil {
		return openfeature.NotReadyState
	}
	return openfeature.ReadyState
}

// GetClient returns an instance of KameleoonClient SDK.
func (p *kameleoonProvider) GetClient() kameleoon.KameleoonClient {
	return p.client
}

// createProviderResolutionDetail creates a ProviderResolutionDetail based on the given error and variant.
func createProviderResolutionDetail(
	err *openfeature.ResolutionError, variant string) openfeature.ProviderResolutionDetail {
	var providerResDetail openfeature.ProviderResolutionDetail
	if err == nil {
		providerResDetail = openfeature.ProviderResolutionDetail{
			Variant: variant,
		}
	} else {
		providerResDetail = openfeature.ProviderResolutionDetail{
			ResolutionError: *err,
			Variant:         variant,
		}
	}
	return providerResDetail
}

func (p *kameleoonProvider) Hooks() []openfeature.Hook {
	return []openfeature.Hook{}
}
