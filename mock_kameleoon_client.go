package kameleoon

import (
	kameleoon "github.com/Kameleoon/client-go/v3"
	"github.com/Kameleoon/client-go/v3/types"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
	"time"
)

// MockKameleoonClient is a mock for the KameleoonClient interface
type MockKameleoonClient struct {
	mock.Mock
}

func (m *MockKameleoonClient) SetLegalConsent(visitorCode string, consent bool, response ...*fasthttp.Response) error {
	panic("implement me")
}

func (m *MockKameleoonClient) AddData(visitorCode string, allData ...types.Data) error {
	args := m.Called(visitorCode, allData)
	return args.Error(0)
}

func (m *MockKameleoonClient) TrackConversion(visitorCode string, goalID int,
	isUniqueIdentifier ...bool) error {
	panic("implement me")
}

func (m *MockKameleoonClient) TrackConversionRevenue(
	visitorCode string, goalID int, revenue float64, isUniqueIdentifier ...bool) error {
	panic("implement me")
}

func (m *MockKameleoonClient) FlushVisitor(visitorCode string, isUniqueIdentifier ...bool) error {
	panic("implement me")
}

func (m *MockKameleoonClient) FlushVisitorInstantly(visitorCode string) error {
	panic("implement me")
}

func (m *MockKameleoonClient) FlushAll(instant ...bool) {
	panic("implement me")
}

func (m *MockKameleoonClient) GetFeatureVariationKey(
	visitorCode string, featureKey string, isUniqueIdentifier ...bool) (string, error) {
	args := m.Called(visitorCode, featureKey, isUniqueIdentifier)
	return args.String(0), args.Error(1)
}

func (m *MockKameleoonClient) GetFeatureVariable(
	visitorCode string, featureKey string, variableKey string, isUniqueIdentifier ...bool) (interface{}, error) {
	panic("implement me")
}

func (m *MockKameleoonClient) IsFeatureActive(
	visitorCode string, featureKey string, isUniqueIdentifier ...bool) (bool, error) {
	panic("implement me")
}

func (m *MockKameleoonClient) GetFeatureVariationVariables(
	featureKey string, variationKey string) (map[string]interface{}, error) {
	args := m.Called(featureKey, variationKey)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockKameleoonClient) GetRemoteData(key string, timeout ...time.Duration) ([]byte, error) {
	panic("implement me")
}

func (m *MockKameleoonClient) GetVisitorWarehouseAudience(
	params kameleoon.VisitorWarehouseAudienceParams) (*types.CustomData, error) {
	panic("implement me")
}

func (m *MockKameleoonClient) GetVisitorWarehouseAudienceWithOptParams(
	visitorCode string, customDataIndex int, params ...kameleoon.VisitorWarehouseAudienceOptParams,
) (*types.CustomData, error) {
	panic("implement me")
}

func (m *MockKameleoonClient) GetRemoteVisitorData(
	visitorCode string, addData bool, timeout ...time.Duration) ([]types.Data, error) {
	panic("implement me")
}

func (m *MockKameleoonClient) GetRemoteVisitorDataWithOptParams(
	visitorCode string, addData bool, filter types.RemoteVisitorDataFilter,
	params ...kameleoon.RemoteVisitorDataOptParams) ([]types.Data, error) {
	panic("implement me")
}

func (m *MockKameleoonClient) GetRemoteVisitorDataWithFilter(
	visitorCode string, addData bool, filter types.RemoteVisitorDataFilter,
	params ...kameleoon.RemoteVisitorDataOptParams) ([]types.Data, error) {
	panic("implement me")
}

func (m *MockKameleoonClient) OnUpdateConfiguration(handler func()) {
	panic("implement me")
}

func (m *MockKameleoonClient) GetFeatureList() []string {
	panic("implement me")
}

func (m *MockKameleoonClient) GetActiveFeatureListForVisitor(visitorCode string) ([]string, error) {
	panic("implement me")
}

func (m *MockKameleoonClient) GetActiveFeatures(visitorCode string) (map[string]types.Variation, error) {
	panic("implement me")
}

func (m *MockKameleoonClient) GetEngineTrackingCode(visitorCode string) string {
	panic("implement me")
}

func (m *MockKameleoonClient) WaitInit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockKameleoonClient) GetVisitorCode(
	request *fasthttp.Request, response *fasthttp.Response, defaultVisitorCode ...string) (string, error) {
	args := m.Called(request, response, defaultVisitorCode)
	return args.String(0), args.Error(1)
}
