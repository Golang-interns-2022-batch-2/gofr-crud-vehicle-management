// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/SukantArora/CRUD_Gofr/internal/store (interfaces: Vehicle)

// Package store is a generated GoMock package.
package store

import (
	gofr "developer.zopsmart.com/go/gofr/pkg/gofr"
	models "github.com/SukantArora/CRUD_Gofr/internal/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockVehicle is a mock of Vehicle interface
type MockVehicle struct {
	ctrl     *gomock.Controller
	recorder *MockVehicleMockRecorder
}

// MockVehicleMockRecorder is the mock recorder for MockVehicle
type MockVehicleMockRecorder struct {
	mock *MockVehicle
}

// NewMockVehicle creates a new mock instance
func NewMockVehicle(ctrl *gomock.Controller) *MockVehicle {
	mock := &MockVehicle{ctrl: ctrl}
	mock.recorder = &MockVehicleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockVehicle) EXPECT() *MockVehicleMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockVehicle) Create(arg0 *gofr.Context, arg1 *models.Vehicle) (*models.Vehicle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*models.Vehicle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockVehicleMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockVehicle)(nil).Create), arg0, arg1)
}

// Delete mocks base method
func (m *MockVehicle) Delete(arg0 *gofr.Context, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockVehicleMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockVehicle)(nil).Delete), arg0, arg1)
}

// Get mocks base method
func (m *MockVehicle) Get(arg0 *gofr.Context) ([]*models.Vehicle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].([]*models.Vehicle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockVehicleMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockVehicle)(nil).Get), arg0)
}

// GetByID mocks base method
func (m *MockVehicle) GetByID(arg0 *gofr.Context, arg1 int) (*models.Vehicle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0, arg1)
	ret0, _ := ret[0].(*models.Vehicle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID
func (mr *MockVehicleMockRecorder) GetByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockVehicle)(nil).GetByID), arg0, arg1)
}

// Update mocks base method
func (m *MockVehicle) Update(arg0 *gofr.Context, arg1 int, arg2 *models.Vehicle) (*models.Vehicle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(*models.Vehicle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockVehicleMockRecorder) Update(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockVehicle)(nil).Update), arg0, arg1, arg2)
}
