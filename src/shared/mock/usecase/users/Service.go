// Code generated by mockery v2.4.0-beta. DO NOT EDIT.

package users

import (
	context "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
	mock "github.com/stretchr/testify/mock"

	users "github.com/cloudsrc/api.awaymail.v1.go/src/usecase/users"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// GetAuthToken provides a mock function with given fields: ctxSess, req, agent
func (_m *Service) GetAuthToken(ctxSess *context.Context, req *users.GetAuthTokenReq, agent string) (users.GetAuthTokenResp, error) {
	ret := _m.Called(ctxSess, req, agent)

	var r0 users.GetAuthTokenResp
	if rf, ok := ret.Get(0).(func(*context.Context, *users.GetAuthTokenReq, string) users.GetAuthTokenResp); ok {
		r0 = rf(ctxSess, req, agent)
	} else {
		r0 = ret.Get(0).(users.GetAuthTokenResp)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*context.Context, *users.GetAuthTokenReq, string) error); ok {
		r1 = rf(ctxSess, req, agent)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
