package kameleoon

import (
	"context"
	kameleoon "github.com/Kameleoon/client-go/v3"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKameleoonProvider_Metadata(t *testing.T) {
	// Arrange
	provider := &kameleoonProvider{}

	// Act
	metadata := provider.Metadata()

	// Assert
	assert.Equal(t, "Kameleoon Provider", metadata.Name)
}

// TODO uncomment this test when will release sdk version 3.4.1
//func TestKameleoonProvider_InitWithInvalidSiteCodeThrowsError(t *testing.T) {
//	// Arrange
//	logging.SetLogLevel(logging.DEBUG)
//	siteCode := ""
//	config := kameleoon.KameleoonClientConfig{
//		ClientID:     "clientId",
//		ClientSecret: "clientSecret",
//	}
//
//	ex := errs.NewSiteCodeIsEmpty("Provided siteCode is empty")
//	expectedError := openfeature.NewProviderNotReadyResolutionError(ex.Error())
//
//	// Act
//	_, err := NewKameleoonProvider(siteCode, &config)
//
//	// Assert
//	assert.Equal(t, expectedError.Error(), err.Error())
//}

func setupResolverMock(m *MockKameleoonResolver, flagKey string, defaultValue interface{}, expectedValue interface{}) {
	m.On("Resolve", context.Background(), flagKey, defaultValue,
		openfeature.FlattenedContext(nil)).Return(expectedValue, (*openfeature.ResolutionError)(nil), "")
}

func assertResult(t *testing.T, result openfeature.ProviderResolutionDetail,
	value interface{}, expectedValue interface{}) {
	assert.Equal(t, expectedValue, value)
	assert.Nil(t, result.Error())
}

func TestKameleoonProvider_ResolveBooleanValueReturnsCorrectValue(t *testing.T) {
	// Arrange
	clientMock := new(MockKameleoonClient)
	resolverMock := new(MockKameleoonResolver)
	var provider = &kameleoonProvider{
		siteCode: "siteCode",
		client:   clientMock,
		resolver: resolverMock,
	}

	defaultValue := false
	expectedValue := true
	setupResolverMock(resolverMock, "flagKey", defaultValue, expectedValue)

	// Act
	result := provider.BooleanEvaluation(context.Background(), "flagKey", defaultValue, nil)

	// Assert
	assertResult(t, result.ProviderResolutionDetail, result.Value, expectedValue)
}

func TestKameleoonProvider_ResolveDoubleValueReturnsCorrectValue(t *testing.T) {
	// Arrange
	clientMock := new(MockKameleoonClient)
	resolverMock := new(MockKameleoonResolver)
	var provider = &kameleoonProvider{
		siteCode: "siteCode",
		client:   clientMock,
		resolver: resolverMock,
	}

	defaultValue := 0.5
	expectedValue := 2.5
	setupResolverMock(resolverMock, "flagKey", defaultValue, expectedValue)

	// Act
	result := provider.FloatEvaluation(context.Background(), "flagKey", defaultValue, nil)

	// Assert
	assertResult(t, result.ProviderResolutionDetail, result.Value, expectedValue)
}

func TestKameleoonProvider_ResolveIntegerValueReturnsCorrectValue(t *testing.T) {
	// Arrange
	clientMock := new(MockKameleoonClient)
	resolverMock := new(MockKameleoonResolver)
	var provider = &kameleoonProvider{
		siteCode: "siteCode",
		client:   clientMock,
		resolver: resolverMock,
	}

	defaultValue := int64(1)
	expectedValue := int64(2)
	setupResolverMock(resolverMock, "flagKey", defaultValue, expectedValue)

	// Act
	result := provider.IntEvaluation(context.Background(), "flagKey", int64(defaultValue), nil)

	// Assert
	assertResult(t, result.ProviderResolutionDetail, result.Value, expectedValue)
}

func TestKameleoonProvider_ResolveStringValueReturnsCorrectValue(t *testing.T) {
	// Arrange
	clientMock := new(MockKameleoonClient)
	resolverMock := new(MockKameleoonResolver)
	var provider = &kameleoonProvider{
		siteCode: "siteCode",
		client:   clientMock,
		resolver: resolverMock,
	}

	defaultValue := "1"
	expectedValue := "2"
	setupResolverMock(resolverMock, "flagKey", defaultValue, expectedValue)

	// Act
	result := provider.StringEvaluation(context.Background(), "flagKey", defaultValue, nil)

	// Assert
	assertResult(t, result.ProviderResolutionDetail, result.Value, expectedValue)
}

func TestKameleoonProvider_ResolveStructureValueReturnsCorrectValue(t *testing.T) {
	// Arrange
	clientMock := new(MockKameleoonClient)
	resolverMock := new(MockKameleoonResolver)
	var provider = &kameleoonProvider{
		siteCode: "siteCode",
		client:   clientMock,
		resolver: resolverMock,
	}

	defaultValue := map[string]interface{}{"k": 10}
	expectedValue := map[string]interface{}{"k1": 20}
	setupResolverMock(resolverMock, "flagKey", defaultValue, expectedValue)

	// Act
	result := provider.ObjectEvaluation(context.Background(), "flagKey", defaultValue, nil)

	// Assert
	assertResult(t, result.ProviderResolutionDetail, result.Value, expectedValue)
}

func TestGetStatus_ReturnsProperStatus(t *testing.T) {
	tests := []struct {
		providedTask   func() error
		expectedStatus openfeature.State
	}{
		{func() error { return nil }, openfeature.ReadyState},
		{func() error { return context.Canceled }, openfeature.NotReadyState},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			// Arrange
			clientMock := new(MockKameleoonClient)
			clientMock.On("WaitInit").Return(tt.providedTask())

			// Act
			provider := &kameleoonProvider{client: clientMock}
			status := provider.Status()

			// Assert
			assert.Equal(t, tt.expectedStatus, status)
		})
	}
}

func TestInitialize_WaitsForClientInitialization(t *testing.T) {
	// Arrange
	clientMock := new(MockKameleoonClient)
	clientMock.On("WaitInit").Return(nil)

	// Act
	provider := &kameleoonProvider{client: clientMock}
	err := provider.Init(openfeature.EvaluationContext{})

	// Assert
	assert.NoError(t, err)
	clientMock.AssertCalled(t, "WaitInit")
}

func TestShutdown_ForgetSiteCode(t *testing.T) {
	// Arrange
	siteCode := "testSiteCode"
	config := kameleoon.KameleoonClientConfig{
		ClientID:     "clientId",
		ClientSecret: "clientSecret",
	}
	provider, _ := NewKameleoonProvider(siteCode, &config)
	clientFirst := provider.GetClient()
	clientToCheck, _ := kameleoon.KameleoonClientFactory.Create(siteCode, &config)

	// Act
	provider.Shutdown()

	providerSecond, _ := NewKameleoonProvider(siteCode, &config)
	clientSecond := providerSecond.GetClient()

	// Assert
	assert.Same(t, clientToCheck, clientFirst)
	assert.NotSame(t, clientFirst, clientSecond)
}
