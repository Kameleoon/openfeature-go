package kameleoon

import (
	"github.com/Kameleoon/client-go/v3/types"
	"github.com/open-feature/go-sdk/openfeature"
)

// dataConverter is used to convert a data from OpenFeature to Kameleoon.
type dataConverter struct {
	conversionMethods map[string]func(interface{}) types.Data
}

// newDataConverter creates a new instance of DataConverter.
func newDataConverter() *dataConverter {
	return &dataConverter{
		conversionMethods: map[string]func(interface{}) types.Data{
			Data.Type.Conversion: makeConversion,
			Data.Type.CustomData: makeCustomData,
		},
	}
}

// Private instance of dataConverter
var dc = newDataConverter()

// ToKameleoon converts FlattenedContext to Kameleoon SDK data types.
func ToKameleoon(context openfeature.FlattenedContext) []types.Data {
	if len(context) == 0 {
		return []types.Data{}
	}

	var data []types.Data
	for key, value := range context {
		var values []interface{}
		if v, ok := value.([]map[string]interface{}); ok {
			for _, item := range v {
				values = append(values, item)
			}
		} else {
			values = []interface{}{value}
		}

		if conversionMethod, ok := dc.conversionMethods[key]; ok {
			for _, val := range values {
				data = append(data, conversionMethod(val))
			}
		}
	}
	return data
}

// makeConversion creates a Conversion object from the value.
func makeConversion(value interface{}) types.Data {
	structData, ok := value.(map[string]interface{})
	if !ok {
		return nil
	}

	goalID, _ := structData[Data.ConversionType.GoalId].(int)
	revenue, ok := structData[Data.ConversionType.Revenue].(float64)
	if !ok {
		if intRevenue, ok := structData[Data.ConversionType.Revenue].(int); ok {
			revenue = float64(intRevenue)
		} else {
			revenue = 0
		}
	}
	return types.NewConversionWithRevenue(goalID, revenue, false)
}

// makeCustomData creates a CustomData object from the value.
func makeCustomData(value interface{}) types.Data {
	structData, ok := value.(map[string]interface{})
	if !ok {
		return nil
	}

	index, _ := structData[Data.CustomDataType.Index].(int)
	var values []string
	if val, ok := structData[Data.CustomDataType.Values].([]string); ok {
		values = val
	} else if s, ok := structData[Data.CustomDataType.Values].(string); ok {
		values = append(values, s)
	}
	return types.NewCustomData(index, values...)
}
