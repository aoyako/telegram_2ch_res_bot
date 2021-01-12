package mock_controller

import gomock "github.com/golang/mock/gomock"

// MockController mocks controller
type MockController struct {
	*MockInfo
	*MockUser
	*MockSubscription
}

// NewMockController constructor for mock controller
func NewMockController(c *gomock.Controller) *MockController {
	return &MockController{
		NewMockInfo(c),
		NewMockUser(c),
		NewMockSubscription(c),
	}
}
