package mock_storage

import gomock "github.com/golang/mock/gomock"

// MockStorage mocks storage
type MockStorage struct {
	*MockUser
	*MockSubscription
	*MockInfo
}

func NewMockStorage(c *gomock.Controller) *MockStorage {
	return &MockStorage{
		NewMockUser(c),
		NewMockSubscription(c),
		NewMockInfo(c),
	}
}
