package kameleoon

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckTypeValues_ProperValues(t *testing.T) {
	// Assert
	assert.Equal(t, "conversion", Data.Type.Conversion)
	assert.Equal(t, "customData", Data.Type.CustomData)

	assert.Equal(t, "index", Data.CustomDataType.Index)
	assert.Equal(t, "values", Data.CustomDataType.Values)

	assert.Equal(t, "goalId", Data.ConversionType.GoalId)
	assert.Equal(t, "revenue", Data.ConversionType.Revenue)
}
