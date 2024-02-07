package main

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
)

type MockDatastore struct {
	mock.Mock
}

func (m *MockDatastore) Find(dest interface{}, conds ...interface{}) (User, error) {
	args := m.Called(dest, conds)
	if tx, ok := args.Get(0).(User); ok {
		return tx, nil
	}
	return User{}, args.Error(1)
}

func (m *MockDatastore) FindAll(dest interface{}, conds ...interface{}) ([]User, error) {
	args := m.Called(dest, conds)
	if users, ok := args.Get(0).([]User); ok {
		return users, nil
	}
	return nil, args.Error(1)
}

func TestGetUser(t *testing.T) {
	mockdatastore := new(MockDatastore)
	service := NewMyService(mockdatastore)
	mockdatastore.On("Find", mock.Anything, mock.Anything).Return(&gorm.DB{})
	user, err := service.GetUser(1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if user == nil {
		t.Error("Expected user to not be nil")
	}

	mockdatastore.On("FindAll", mock.Anything, mock.Anything).Return([]User{
		{
			ID:   1,
			Name: "John",
		},
		{
			ID:   2,
			Name: "Jane",
		},
	}, nil)

	users, err := service.GetAllUsers()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}

	mockdatastore.AssertExpectations(t)
}
