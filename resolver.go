package kameleoon

import (
	"context"
	"fmt"
	"reflect"
	"sort"

	kameleoon "github.com/Kameleoon/client-go/v3"
	"github.com/open-feature/go-sdk/openfeature"
)

// resolver interface which contains method for evalutions based on provided data
type resolver interface {
	Resolve(
		ctx context.Context, flag string, defaultValue interface{}, evalCtx openfeature.FlattenedContext,
	) (interface{}, *openfeature.ResolutionError, string)
}

// kameleoonResolver makes evalutions based on provided data, conforms to Resolver interface
type kameleoonResolver struct {
	client kameleoon.KameleoonClient
}

// newKameleoonResolver creates a new instance of KameleoonResolver.
func newKameleoonResolver(client kameleoon.KameleoonClient) *kameleoonResolver {
	return &kameleoonResolver{
		client: client,
	}
}

// Resolve is main method for getting resolution details based on provided data.
func (r *kameleoonResolver) Resolve(
	context context.Context, flag string, defaultValue interface{}, evalContext openfeature.FlattenedContext,
) (interface{}, *openfeature.ResolutionError, string) {
	// Get visitor code from context.
	visitorCode, ok := getTargetingKey(evalContext)
	if !ok {
		resError := openfeature.NewTargetingKeyMissingResolutionError(
			"The TargetingKey is required in context and cannot be omitted.")
		return defaultValue, &resError, ""
	}

	// Add targeting data from context to KameleoonClient by visitor code
	data := ToKameleoon(evalContext)
	err := r.client.AddData(visitorCode, data...)
	if err != nil {
		resError := openfeature.NewInvalidContextResolutionError(err.Error())
		return defaultValue, &resError, ""
	}

	// Get a variant
	variant, err := r.client.GetFeatureVariationKey(visitorCode, flag)
	if err != nil {
		resError := openfeature.NewFlagNotFoundResolutionError(err.Error())
		return defaultValue, &resError, variant
	}

	// Get the all variables for the variant
	variables, err := r.client.GetFeatureVariationVariables(flag, variant)
	if err != nil {
		resError := openfeature.NewFlagNotFoundResolutionError(err.Error())
		return defaultValue, &resError, variant
	}

	// Get variableKey if it's provided in context or any first in variation.
	// It's the responsibility of the client to have only one variable per variation if
	// variableKey is not provided.
	variableKey := getVariableKey(evalContext, variables)

	// Try to get value by variable key
	value, ok := variables[variableKey]
	if !ok || variableKey == "" {
		resError := openfeature.NewFlagNotFoundResolutionError(makeErrorDescription(variant, variableKey))
		return defaultValue, &resError, variant
	}

	// Check if the variable value has a required type
	if reflect.TypeOf(value) != reflect.TypeOf(defaultValue) {
		resError := openfeature.NewTypeMismatchResolutionError(
			"The type of value received is different from the requested value.")
		return defaultValue, &resError, variant
	}

	return value, nil, variant
}

// getTargetingKey retrieves the targeting key from the provided evaluation context.
func getTargetingKey(evalContext openfeature.FlattenedContext) (string, bool) {
	if targetingKey, ok := evalContext["targetingKey"].(string); ok && targetingKey != "" {
		return targetingKey, true
	}
	return "", false
}

// getVariableKey retrieves the variable key from the provided context or variables map.
func getVariableKey(context openfeature.FlattenedContext, variables map[string]interface{}) string {
	var variableKey string

	if value, ok := context["variableKey"].(string); ok && value != "" {
		variableKey = value
	} else if len(variables) > 0 {
		keys := make([]string, 0, len(variables))
		for k := range variables {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return keys[0]
	}

	return variableKey
}

// makeErrorDescription generates a descriptive error message based on the provided variant and variableKey.
func makeErrorDescription(variant, variableKey string) string {
	if variableKey == "" {
		return fmt.Sprintf("The variation '%s' has no variables", variant)
	}
	return fmt.Sprintf("The value for provided variable key '%s' isn't found in variation '%s'",
		variableKey, variant)
}
