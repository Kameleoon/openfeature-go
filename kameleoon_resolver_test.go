package kameleoon

import (
	"context"
	"fmt"
	"github.com/Kameleoon/client-go/v3/errs"
	"github.com/Kameleoon/client-go/v3/types"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolve_WithNilContext_ReturnsErrorForMissingTargetingKey(t *testing.T) {
	// Arrange
	clientMock := new(MockKameleoonClient)
	resolver := newKameleoonResolver(clientMock)
	flagKey := "testFlag"
	defaultValue := "defaultValue"
	expectedError := openfeature.NewTargetingKeyMissingResolutionError(
		"The TargetingKey is required in context and cannot be omitted.")

	// Act
	result, err, variant := resolver.Resolve(context.Background(), flagKey, defaultValue, nil)

	// Assert
	assert.Equal(t, defaultValue, result)
	assert.Equal(t, expectedError.Error(), err.Error())
	assert.Empty(t, variant)
}

func TestResolve_NoMatchVariables_ReturnsErrorForFlagNotFound(t *testing.T) {
	// Arrange
	flagKey := "testFlag"
	defaultValue := 42

	testCases := []struct {
		variant          string
		addVariableKey   bool
		variables        map[string]interface{}
		expectedErrorMsg string
	}{
		{"on", false, map[string]interface{}{},
			"The variation 'on' has no variables"},
		{"var", true, map[string]interface{}{"key": new(interface{})},
			"The value for provided variable key 'variableKey' isn't found in variation 'var'"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("variant=%s", tc.variant), func(t *testing.T) {
			clientMock := new(MockKameleoonClient)
			visitorCode := "testVisitor"
			clientMock.On("GetFeatureVariationKey", visitorCode, flagKey, []bool(nil)).Return(tc.variant, nil)
			clientMock.On("GetFeatureVariationVariables", flagKey, tc.variant).Return(tc.variables, nil)
			clientMock.On("AddData", visitorCode, []types.Data(nil)).Return(nil)

			resolver := newKameleoonResolver(clientMock)

			evalContext := openfeature.FlattenedContext{
				"targetingKey": visitorCode,
			}
			if tc.addVariableKey {
				evalContext["variableKey"] = "variableKey"
			}

			expectedError := openfeature.NewFlagNotFoundResolutionError(tc.expectedErrorMsg)

			// Act
			result, err, variant := resolver.Resolve(context.Background(), flagKey, defaultValue, evalContext)

			// Assert
			assert.Equal(t, defaultValue, result)
			assert.Equal(t, expectedError.Error(), err.Error())
			assert.Equal(t, tc.variant, variant)
		})
	}
}

func TestResolve_MismatchType_ReturnsErrorTypeMismatch(t *testing.T) {
	// Arrange
	flagKey := "testFlag"
	expectedVariant := "on"
	defaultValue := 42

	testCases := []interface{}{
		true,
		"string",
		10.0,
	}

	for _, returnValue := range testCases {
		t.Run(fmt.Sprintf("returnValue=%v", returnValue), func(t *testing.T) {
			clientMock := new(MockKameleoonClient)
			visitorCode := "testVisitor"
			clientMock.On("GetFeatureVariationKey", visitorCode, flagKey, []bool(nil)).Return(expectedVariant, nil)
			clientMock.On("GetFeatureVariationVariables", flagKey, expectedVariant).Return(map[string]any{
				"key": returnValue,
			}, nil)
			clientMock.On("AddData", visitorCode, []types.Data(nil)).Return(nil)

			resolver := newKameleoonResolver(clientMock)
			evalContext := openfeature.FlattenedContext{
				"targetingKey": visitorCode,
			}

			expectedError := openfeature.NewTypeMismatchResolutionError(
				"The type of value received is different from the requested value.")

			// Act
			result, err, variant := resolver.Resolve(context.Background(), flagKey, defaultValue, evalContext)

			// Assert
			assert.Equal(t, defaultValue, result)
			assert.Equal(t, expectedError.Error(), err.Error())
			assert.Equal(t, expectedVariant, variant)
		})
	}
}

func TestResolve_KameleoonException_FlagNotFound(t *testing.T) {
	// Arrange
	flagKey := "testFlag"
	visitorCode := "testVisitor"
	defaultValue := 42

	exception := errs.NewFeatureNotFound("featureException")

	clientMock := new(MockKameleoonClient)
	clientMock.On("AddData", visitorCode, []types.Data(nil)).Return(nil)
	clientMock.On("GetFeatureVariationKey", visitorCode, flagKey, []bool(nil)).Return("", exception)

	resolver := newKameleoonResolver(clientMock)
	evalContext := openfeature.FlattenedContext{
		"targetingKey": visitorCode,
	}

	expectedError := openfeature.NewFlagNotFoundResolutionError(exception.Error())

	// Act
	result, err, variant := resolver.Resolve(context.Background(), flagKey, defaultValue, evalContext)

	// Assert
	assert.Equal(t, defaultValue, result)
	assert.Equal(t, expectedError.Error(), err.Error())
	assert.Empty(t, variant)
}

func TestResolve_KameleoonException_VisitorCodeInvalid(t *testing.T) {
	// Arrange
	flagKey := "testFlag"
	visitorCode := "testVisitor"
	defaultValue := 42

	exception := errs.NewVisitorCodeInvalid("visitorCodeInvalid")

	clientMock := new(MockKameleoonClient)
	clientMock.On("AddData", visitorCode, []types.Data(nil)).Return(exception)

	resolver := newKameleoonResolver(clientMock)
	evalContext := openfeature.FlattenedContext{
		"targetingKey": visitorCode,
	}

	expectedError := openfeature.NewInvalidContextResolutionError(exception.Error())

	// Act
	result, err, variant := resolver.Resolve(context.Background(), flagKey, defaultValue, evalContext)

	// Assert
	assert.Equal(t, defaultValue, result)
	assert.Equal(t, expectedError.Error(), err.Error())
	assert.Empty(t, variant)
}

func TestResolve_ReturnsResultDetails(t *testing.T) {
	// Arrange
	flagKey := "testFlag"
	visitorCode := "testVisitor"
	expectedVariant := "variant"

	testCases := []struct {
		variableKey   string
		variables     map[string]interface{}
		expectedValue interface{}
		defaultValue  interface{}
	}{
		{"", map[string]interface{}{"k": 10}, 10, 9},
		{"", map[string]interface{}{"k1": "str"}, "str", "st"},
		{"", map[string]interface{}{"k2": true}, true, false},
		{"", map[string]interface{}{"k3": 10.0}, 10.0, 11.0},
		{"varKey", map[string]interface{}{"varKey": 10.0}, 10.0, 11.0},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("variableKey=%s", tc.variableKey), func(t *testing.T) {
			clientMock := new(MockKameleoonClient)
			clientMock.On("AddData", visitorCode, []types.Data(nil)).Return(nil)
			clientMock.On("GetFeatureVariationKey", visitorCode, flagKey, []bool(nil)).Return(expectedVariant, nil)
			clientMock.On("GetFeatureVariationVariables", flagKey, expectedVariant).Return(tc.variables, nil)

			resolver := newKameleoonResolver(clientMock)
			evalContext := openfeature.FlattenedContext{
				"targetingKey": visitorCode,
			}
			if tc.variableKey != "" {
				evalContext["variableKey"] = tc.variableKey
			}

			// Act
			result, err, variant := resolver.Resolve(context.Background(), flagKey, tc.defaultValue, evalContext)

			// Assert
			assert.Equal(t, tc.expectedValue, result)
			assert.Equal(t, expectedVariant, variant)
			assert.Nil(t, err)
		})
	}
}
