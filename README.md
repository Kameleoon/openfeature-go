# Kameleoon OpenFeature provider for Go

The Kameleoon OpenFeature provider for Go allows you to connect your OpenFeature Go implementation to Kameleoon without installing the Go Kameleoon SDK.

> [!WARNING]
> This is a beta version. Breaking changes may be introduced before general release.

## Supported Go versions

This version of the SDK is built for the following targets:

* Go 1.18 and above.

## Get started

This section explains how to install, configure, and customize the Kameleoon OpenFeature provider.

### Install dependencies

First, install the required dependencies in your application.

```sh
go get github.com/Kameleoon/openfeature-go
```

### Usage

The following example shows how to use the Kameleoon provider with the OpenFeature SDK.

```go
package main

import (
	"context"
	"fmt"
	"github.com/Kameleoon/client-go/v3"
	"github.com/open-feature/go-sdk/openfeature"
	"time"
)

func main() {
	var provider *kameleoonProvider
	visitorCode := "visitorCode"
	featureKey := "featureKey"

	clientConfig := kameleoon.KameleoonClientConfig{
		ClientID:        "clientId",
		ClientSecret:    "clientSecret",
		TopLevelDomain:  "topLevelDomain",
	}

	provider, err := NewKameleoonProvider("siteCode", &clientConfig)
	if err != nil {
		fmt.Println("Error creating provider:", err)
		return
	}

	err = openfeature.SetProvider(provider)
	if err != nil {
		fmt.Println("Error setting provider:", err)
		return
	}

	client := openfeature.NewClient("")

	dataDictionary := map[string]interface{}{
		"variableKey": "variableKey",
	}

	evalContext := openfeature.NewEvaluationContext(visitorCode, dataDictionary)

	numberOfRecommendedProducts, _ := client.IntValue(context.Background(), featureKey, 5, evalContext)
	showRecommendedProducts(numberOfRecommendedProducts)
}

func showRecommendedProducts(numberOfRecommendedProducts int64) {
	fmt.Printf("Number of recommended products: %d\n", numberOfRecommendedProducts)
}
```

#### Customize the Kameleoon provider

You can customize the Kameleoon provider by changing the `KameleoonClientConfig` object that you passed to the constructor above. For example:

```go
import (
    "github.com/Kameleoon/client-go/v3"
    "github.com/Kameleoon/openfeature-go"
)

clientConfig := kameleoon.KameleoonClientConfig{
	ClientID:        "clientId",
	ClientSecret:    "clientSecret",
	TopLevelDomain:  "topLevelDomain",
	RefreshInterval: 1 * time.Minute,    // Optional field
    SessionDuration: 5 * time.Minute,    // Optional field
}

provider, err := kameleoon.NewKameleoonProvider("siteCode", &clientConfig)
if err != nil {
	fmt.Println("Error creating provider:", err)
	return
}
```
> [!NOTE]
> For additional configuration options, see the [Kameleoon documentation](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#example-code).

## EvaluationContext and Kameleoon Data

Kameleoon uses the concept of associating `Data` to users, while the OpenFeature SDK uses the concept of an `EvaluationContext`, which is a dictionary of string keys and values. The Kameleoon provider maps the `EvaluationContext` to the Kameleoon `Data`.

> [!NOTE]
> To get the evaluation for a specific visitor, set the `targetingKey` value for the `EvaluationContext` to the visitor code (user ID). If the value is not provided, then the `defaultValue` parameter will be returned.

```go
evalContext := openfeature.NewEvaluationContext("userId", nil)
```

The Kameleoon provider provides a few predefined parameters that you can use to target a visitor from a specific audience and track each conversion. These are:

| Parameter              | Description                                                                                                                                                         |
|------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `Data.Type.CustomData` | The parameter is used to set [`CustomData`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#customdata) for a visitor.     |
| `Data.Type.Conversion` | The parameter is used to track a [`Conversion`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#conversion) for a visitor. |

### Data.Type.CustomData

Use `Data.Type.CustomData` to set [`CustomData`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#customdata) for a visitor. The `Data.Type.CustomData` field has the following parameters:

| Parameter                    | Type   | Description                                                       |
|------------------------------|--------|-------------------------------------------------------------------|
| `Data.CustomDataType.Index`  | int    | Index or ID of the custom data to store. This field is mandatory. |
| `Data.CustomDataType.Values` | string | Value of the custom data to store. This field is mandatory.       |

#### Example

```go
customDataDictionary := map[string]interface{}{
    Data.Type.CustomData: map[string]interface{}{
        Data.CustomDataType.INDEX:  1,
        Data.CustomDataType.VALUES: "10",
	},
}

evalContext := openfeature.NewEvaluationContext("userId", customDataDictionary)
```

### Data.Type.Conversion

Use `Data.Type.Conversion` to track a [`Conversion`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#conversion) for a visitor. The `Data.Type.Conversion` field has the following parameters:

| Parameter                     | Type  | Description                                                     |
|-------------------------------|-------|-----------------------------------------------------------------|
| `Data.ConversionType.goalId`  | int   | Identifier of the goal. This field is mandatory.                |
| `Data.ConversionType.Revenue` | float | Revenue associated with the conversion. This field is optional. |

#### Example

```go
conversionDictionary := map[string]interface{}{
    Data.ConversionType.GOAL_ID: 1,
    Data.ConversionType.REVENUE: 200,
}

evalContext := openfeature.NewEvaluationContext("userId", map[string]interface{}{
	Data.Type.Conversion: conversionDictionary,
})
```

### Use multiple Kameleoon Data types

You can provide many different kinds of Kameleoon data within a single `EvaluationContext` instance.

For example, the following code provides one `Data.Type.Conversion` instance and two `Data.Type.CustomData` instances.

```go
dataDictionary := map[string]interface{}{
    Data.Type.Conversion: map[string]interface{}{
        Data.ConversionType.GoalId: 1,
        Data.ConversionType.Revenue: 200,
    },
    Data.Type.CustomData: []map[string]interface{}{
        {
            Data.CustomDataType.Index:  1,
            Data.CustomDataType.Values: []string{"10", "30"},
        },
        {
            Data.CustomDataType.Index:  2,
            Data.CustomDataType.Values: "20",
        },
    },
}

evalContext := openfeature.NewEvaluationContext("userId", dataDictionary)
```
