// Code generated by mockery v2.4.0-beta. DO NOT EDIT.

package aways

import (
	aways "github.com/cloudsrc/api.awaymail.v1.go/src/usecase/aways"
	bson "gopkg.in/mgo.v2/bson"

	context "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"

	mock "github.com/stretchr/testify/mock"

	models "github.com/cloudsrc/api.awaymail.v1.go/src/shared/models"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// CreateAway provides a mock function with given fields: ctxSess, user, req
func (_m *Service) CreateAway(ctxSess *context.Context, user models.UserSession, req *aways.CreateAwayReq) (*aways.CreateAwayResp, error) {
	ret := _m.Called(ctxSess, user, req)

	var r0 *aways.CreateAwayResp
	if rf, ok := ret.Get(0).(func(*context.Context, models.UserSession, *aways.CreateAwayReq) *aways.CreateAwayResp); ok {
		r0 = rf(ctxSess, user, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*aways.CreateAwayResp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*context.Context, models.UserSession, *aways.CreateAwayReq) error); ok {
		r1 = rf(ctxSess, user, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAway provides a mock function with given fields: ctxSess, user, awayID
func (_m *Service) DeleteAway(ctxSess *context.Context, user models.UserSession, awayID string) error {
	ret := _m.Called(ctxSess, user, awayID)

	var r0 error
	if rf, ok := ret.Get(0).(func(*context.Context, models.UserSession, string) error); ok {
		r0 = rf(ctxSess, user, awayID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EnableAllAway provides a mock function with given fields: ctxSess, user, isEnabled
func (_m *Service) EnableAllAway(ctxSess *context.Context, user models.UserSession, isEnabled bool) (int, error) {
	ret := _m.Called(ctxSess, user, isEnabled)

	var r0 int
	if rf, ok := ret.Get(0).(func(*context.Context, models.UserSession, bool) int); ok {
		r0 = rf(ctxSess, user, isEnabled)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*context.Context, models.UserSession, bool) error); ok {
		r1 = rf(ctxSess, user, isEnabled)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EnableAway provides a mock function with given fields: ctxSess, user, enable, awayID
func (_m *Service) EnableAway(ctxSess *context.Context, user models.UserSession, enable bool, awayID bson.ObjectId) (*aways.CreateAwayResp, error) {
	ret := _m.Called(ctxSess, user, enable, awayID)

	var r0 *aways.CreateAwayResp
	if rf, ok := ret.Get(0).(func(*context.Context, models.UserSession, bool, bson.ObjectId) *aways.CreateAwayResp); ok {
		r0 = rf(ctxSess, user, enable, awayID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*aways.CreateAwayResp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*context.Context, models.UserSession, bool, bson.ObjectId) error); ok {
		r1 = rf(ctxSess, user, enable, awayID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAway provides a mock function with given fields: ctxSess, awayID
func (_m *Service) GetAway(ctxSess *context.Context, awayID string) (*aways.CreateAwayResp, error) {
	ret := _m.Called(ctxSess, awayID)

	var r0 *aways.CreateAwayResp
	if rf, ok := ret.Get(0).(func(*context.Context, string) *aways.CreateAwayResp); ok {
		r0 = rf(ctxSess, awayID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*aways.CreateAwayResp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*context.Context, string) error); ok {
		r1 = rf(ctxSess, awayID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAwayList provides a mock function with given fields: ctxSess, user
func (_m *Service) GetAwayList(ctxSess *context.Context, user models.UserSession) ([]*aways.CreateAwayResp, error) {
	ret := _m.Called(ctxSess, user)

	var r0 []*aways.CreateAwayResp
	if rf, ok := ret.Get(0).(func(*context.Context, models.UserSession) []*aways.CreateAwayResp); ok {
		r0 = rf(ctxSess, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*aways.CreateAwayResp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*context.Context, models.UserSession) error); ok {
		r1 = rf(ctxSess, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateAway provides a mock function with given fields: ctxSess, user, req
func (_m *Service) UpdateAway(ctxSess *context.Context, user models.UserSession, req *aways.UpdateAwayReq) (*aways.CreateAwayResp, error) {
	ret := _m.Called(ctxSess, user, req)

	var r0 *aways.CreateAwayResp
	if rf, ok := ret.Get(0).(func(*context.Context, models.UserSession, *aways.UpdateAwayReq) *aways.CreateAwayResp); ok {
		r0 = rf(ctxSess, user, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*aways.CreateAwayResp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*context.Context, models.UserSession, *aways.UpdateAwayReq) error); ok {
		r1 = rf(ctxSess, user, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
