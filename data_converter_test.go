package kameleoon

import (
	"math/rand"
	"testing"
	"time"

	"github.com/Kameleoon/client-go/v3/types"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/stretchr/testify/assert"
)

func TestToKameleoon_NullContext_ReturnsEmpty(t *testing.T) {
	// Arrange
	var context openfeature.FlattenedContext // nil context

	// Act
	result := ToKameleoon(context)

	// Assert
	assert.Empty(t, result)
}

func TestToKameleoon_WithConversionData_ReturnsConversionData(t *testing.T) {
	tests := []struct {
		name       string
		addRevenue bool
	}{
		{name: "WithRevenue", addRevenue: true},
		{name: "WithoutRevenue", addRevenue: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			rand.Seed(time.Now().UnixNano())
			expectedGoalId := rand.Int()
			expectedRevenue := rand.Float64()

			conversionData := map[string]interface{}{
				Data.ConversionType.GoalId: expectedGoalId,
			}

			if tt.addRevenue {
				conversionData[Data.ConversionType.Revenue] = expectedRevenue
			}

			context := openfeature.FlattenedContext{
				Data.Type.Conversion: conversionData,
			}

			// Act
			result := ToKameleoon(context)

			// Assert
			assert.Len(t, result, 1)
			conversion, ok := result[0].(*types.Conversion)
			assert.True(t, ok)
			assert.Equal(t, expectedGoalId, conversion.GoalId())

			if tt.addRevenue {
				assert.Equal(t, expectedRevenue, conversion.Revenue())
			}
		})
	}
}

func TestToKameleoon_WithCustomData_ReturnsCustomData(t *testing.T) {
	tests := []struct {
		name           string
		expectedIndex  int
		expectedValues []string
	}{
		{name: "EmptyValues", expectedIndex: rand.Int(), expectedValues: []string{}},
		{name: "SingleValue", expectedIndex: rand.Int(), expectedValues: []string{"v1"}},
		{name: "MultipleValues", expectedIndex: rand.Int(), expectedValues: []string{"v1", "v2", "v3"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			customData := map[string]interface{}{
				Data.CustomDataType.Index:  tt.expectedIndex,
				Data.CustomDataType.Values: tt.expectedValues,
			}

			context := openfeature.FlattenedContext{
				Data.Type.CustomData: customData,
			}

			// Act
			result := ToKameleoon(context)

			// Assert
			assert.Len(t, result, 1)
			customDataObj, ok := result[0].(*types.CustomData)
			assert.True(t, ok)
			assert.Equal(t, tt.expectedIndex, customDataObj.ID())
			assert.Equal(t, tt.expectedValues, customDataObj.Values())
		})
	}
}

func TestToKameleoonData_AllTypes_ReturnsAllData(t *testing.T) {
	// Arrange
	rand.Seed(time.Now().UnixNano())
	goalId1 := rand.Int()
	goalId2 := rand.Int()
	index1 := rand.Int()
	index2 := rand.Int()

	context := openfeature.FlattenedContext{
		Data.Type.Conversion: []map[string]interface{}{
			{
				Data.ConversionType.GoalId: goalId1,
			},
			{
				Data.ConversionType.GoalId: goalId2,
			},
		},
		Data.Type.CustomData: []map[string]interface{}{
			{
				Data.CustomDataType.Index: index1,
			},
			{
				Data.CustomDataType.Index: index2,
			},
		},
	}

	// Act
	result := ToKameleoon(context)

	var conversions []*types.Conversion
	var customData []*types.CustomData

	for _, item := range result {
		switch v := item.(type) {
		case *types.Conversion:
			conversions = append(conversions, v)
		case *types.CustomData:
			customData = append(customData, v)
		}
	}

	// Assert
	assert.Len(t, result, 4)

	assert.Equal(t, goalId1, conversions[0].GoalId())
	assert.Equal(t, goalId2, conversions[1].GoalId())
	assert.Equal(t, index1, customData[0].ID())
	assert.Equal(t, index2, customData[1].ID())
}
