package kameleoon

// Data is used to add different Kameleoon data types using
// the FlattenedContext from the OpenFeature SDK.
var Data = struct {
	// Type is used to add Conversion and CustomData using FlattenedContext from the OpenFeature SDK.
	Type struct {
		Conversion string
		CustomData string
	}
	// CustomDataType is used to add CustomData using FlattenedContext from the OpenFeature SDK.
	CustomDataType struct {
		Index  string
		Values string
	}
	// ConversionType is used to add Conversion using FlattenedContext from the OpenFeature SDK.
	ConversionType struct {
		GoalId  string
		Revenue string
	}
}{
	Type: struct {
		Conversion string
		CustomData string
	}{
		Conversion: "conversion",
		CustomData: "customData",
	},
	CustomDataType: struct {
		Index  string
		Values string
	}{
		Index:  "index",
		Values: "values",
	},
	ConversionType: struct {
		GoalId  string
		Revenue string
	}{
		GoalId:  "goalId",
		Revenue: "revenue",
	},
}
