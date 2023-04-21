package users_test

import (
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/logger"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/cloudsrc/api.awaymail.v1.go/config"
	domain "github.com/cloudsrc/api.awaymail.v1.go/src/domain/users"
	"github.com/cloudsrc/api.awaymail.v1.go/src/infrastructure/google"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/constants"
	mockRepoUsers "github.com/cloudsrc/api.awaymail.v1.go/src/shared/mock/mongo/users"
	mockWrapperGoogle "github.com/cloudsrc/api.awaymail.v1.go/src/shared/mock/wrapper/google"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/models"
	. "github.com/cloudsrc/api.awaymail.v1.go/src/usecase/users"
)

var (
	ctxSess = context.New(logger.NewNoopLogger())
)

func TestGetAuthToken(t *testing.T) {
	type mockRepo struct {
		GetUser      *domain.User
		GetUserError error

		AddUser      *domain.User
		AddUserError error

		UpdateUserError error
	}

	type mockWrapper struct {
		UserProfile      google.UserProfile
		UserProfileError error
	}

	tests := []struct {
		name     string
		req      *GetAuthTokenReq
		user     models.UserSession
		repo     mockRepo
		wrapper  mockWrapper
		wantResp GetAuthTokenResp
		wantErr  bool
		wantData bool
	}{
		{
			name:    "error provider - no supported",
			req:     &GetAuthTokenReq{Provider: "privateEmail"},
			wantErr: true,
			wrapper: mockWrapper{
				UserProfileError: constants.ErrorClientNotSupported,
			},
		},
		{
			name:    "error google - GetProfile",
			req:     &GetAuthTokenReq{Provider: "gmail"},
			wantErr: true,
			wrapper: mockWrapper{
				UserProfileError: constants.ErrorGeneral,
			},
		},
		{
			name:    "error database - GetUser",
			req:     &GetAuthTokenReq{Provider: "gmail"},
			wantErr: true,
			wrapper: mockWrapper{
				UserProfile: google.UserProfile{
					Email: "test@gmail.com",
				},
			},
			repo: mockRepo{
				GetUserError: constants.ErrorDatabase,
			},
		},
		{
			name:    "error database - AddUser",
			req:     &GetAuthTokenReq{Provider: "gmail"},
			wantErr: true,
			wrapper: mockWrapper{
				UserProfile: google.UserProfile{
					Email: "test@gmail.com",
				},
			},
			repo: mockRepo{
				AddUserError: constants.ErrorDatabase,
			},
		},
		{
			name:    "error database - UpdateUser",
			req:     &GetAuthTokenReq{Provider: "gmail"},
			wantErr: true,
			wrapper: mockWrapper{
				UserProfile: google.UserProfile{
					Email: "test@gmail.com",
				},
			},
			repo: mockRepo{
				GetUser: &domain.User{
					ID:    "62056d328945fc00018c49b0",
					Email: "test@gmail.com",
				},
				UpdateUserError: constants.ErrorDatabase,
			},
		},
		{
			name: "success - AddUser",
			req:  &GetAuthTokenReq{Provider: "gmail"},
			wrapper: mockWrapper{
				UserProfile: google.UserProfile{
					Email: "test@gmail.com",
				},
			},
			repo: mockRepo{
				AddUser: &domain.User{
					ID:    "62056d328945fc00018c49b0",
					Email: "test@gmail.com",
				},
			},
			wantResp: GetAuthTokenResp{
				Token:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjIwMDEzNTBmNTkyMmQ3YThlZTM0Y2Q5IiwiZW1haWwiOiJ2aWNreUBjbG91ZHNvdXJjZS5pbyIsInJlZnJlc2hfdG9rZW4iOiIxLy8wZ1ljSDZMMzB2Yk9zQ2dZSUFSQUFHQkFTTndGLUw5SXJTVDhiSi14ejNUUTE5YmlBZER3QXJkWGE2akluanE3RWR3NU1KZDlfd1ZVQkVsemhWWGwybmxtMkw4cFUweFIxRDVzIiwiYXV0aF90b2tlbiI6IiIsInByb3ZpZGVyIjoiZ21haWwiLCJleHAiOjE2NzYyMzQ1NDB9.WlmPoGozCICa2ZSrfHqSiS4KsTVDPhoRAyWThHE9lmo",
				Expiry: time.Now().Add(2 * time.Hour),
			},
		},
		{
			name: "success - UpdateUser",
			req:  &GetAuthTokenReq{Provider: "gmail"},
			wrapper: mockWrapper{
				UserProfile: google.UserProfile{
					Email: "test@gmail.com",
				},
			},
			repo: mockRepo{
				GetUser: &domain.User{
					ID:    "62056d328945fc00018c49b0",
					Email: "test@gmail.com",
				},
			},
			wantResp: GetAuthTokenResp{
				Token:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjIwMDEzNTBmNTkyMmQ3YThlZTM0Y2Q5IiwiZW1haWwiOiJ2aWNreUBjbG91ZHNvdXJjZS5pbyIsInJlZnJlc2hfdG9rZW4iOiIxLy8wZ1ljSDZMMzB2Yk9zQ2dZSUFSQUFHQkFTTndGLUw5SXJTVDhiSi14ejNUUTE5YmlBZER3QXJkWGE2akluanE3RWR3NU1KZDlfd1ZVQkVsemhWWGwybmxtMkw4cFUweFIxRDVzIiwiYXV0aF90b2tlbiI6IiIsInByb3ZpZGVyIjoiZ21haWwiLCJleHAiOjE2NzYyMzQ1NDB9.WlmPoGozCICa2ZSrfHqSiS4KsTVDPhoRAyWThHE9lmo",
				Expiry: time.Now().Add(2 * time.Hour),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGoogleWrapper := &mockWrapperGoogle.Wrapper{}
			mockGoogleWrapper.On("GetProfile", mock.Anything, mock.Anything, mock.Anything).Return(tt.wrapper.UserProfile, tt.wrapper.UserProfileError)
			mockUsersRepo := &mockRepoUsers.Repository{}
			mockUsersRepo.On("GetUser", mock.Anything).Return(tt.repo.GetUser, tt.repo.GetUserError)
			mockUsersRepo.On("AddUser", mock.Anything).Return(tt.repo.AddUser, tt.repo.AddUserError)
			mockUsersRepo.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(tt.repo.UpdateUserError)

			s := New(config.Config{}, mockUsersRepo, mockGoogleWrapper)

			resp, err := s.GetAuthToken(ctxSess, tt.req, "ios")
			if (err != nil) != tt.wantErr {
				t.Errorf("s.GetAuthToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantResp.Token != "" {
				resp.Token = tt.wantResp.Token
				resp.Expiry = tt.wantResp.Expiry
			}
			if !reflect.DeepEqual(resp, tt.wantResp) {
				t.Errorf("s.GetAuthToken() = %v, want %v", resp, tt.wantResp)
			}
		})
	}
}
