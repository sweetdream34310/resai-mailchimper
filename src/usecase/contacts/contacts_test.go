package contacts_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"

	domain "github.com/cloudsrc/api.awaymail.v1.go/src/domain/contacts"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/constants"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/logger"
	mockRepoContacts "github.com/cloudsrc/api.awaymail.v1.go/src/shared/mock/mongo/contacts"
	mockRepoRedis "github.com/cloudsrc/api.awaymail.v1.go/src/shared/mock/redis"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/models"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
	. "github.com/cloudsrc/api.awaymail.v1.go/src/usecase/contacts"
)

var (
	ctxSess = context.New(logger.NewNoopLogger())
)

func TestAddContacts(t *testing.T) {
	type mockRepo struct {
		ContactAdd      *domain.Contacts
		ContactAddError error
	}

	tests := []struct {
		name     string
		req      ContactAddReq
		user     models.UserSession
		repo     mockRepo
		wantResp *domain.Contacts
		wantErr  bool
		wantData bool
	}{
		{
			name:    "error database",
			req:     ContactAddReq{Email: "test@test.com", Name: "test"},
			wantErr: true,
			repo: mockRepo{
				ContactAddError: errors.New("error database"),
			},
		},
		{
			name: "success with empty name",
			req:  ContactAddReq{Email: "test@test.com", Name: ""},
			repo: mockRepo{
				ContactAdd: &domain.Contacts{
					ID:     "61fc3cedf5922d9d54954369",
					UserID: "61f9c63d5496b70001e03d23",
					Name:   "test@test.com",
					Email:  "test@test.com",
				},
			},
			wantResp: &domain.Contacts{
				ID:     "61fc3cedf5922d9d54954369",
				UserID: "61f9c63d5496b70001e03d23",
				Name:   "test@test.com",
				Email:  "test@test.com",
			},
		},
		{
			name: "success",
			req:  ContactAddReq{Email: "test@test.com", Name: "test"},
			repo: mockRepo{
				ContactAdd: &domain.Contacts{
					ID:     "61fc3cedf5922d9d54954369",
					UserID: "61f9c63d5496b70001e03d23",
					Name:   "test",
					Email:  "test@test.com",
				},
			},
			wantResp: &domain.Contacts{
				ID:     "61fc3cedf5922d9d54954369",
				UserID: "61f9c63d5496b70001e03d23",
				Name:   "test",
				Email:  "test@test.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRedisRepo := &mockRepoRedis.Repository{}
			mockContactsRepo := &mockRepoContacts.Repository{}
			mockContactsRepo.On("Add", mock.Anything).Return(tt.repo.ContactAdd, tt.repo.ContactAddError)

			s := New(mockContactsRepo, mockRedisRepo)

			resp, err := s.AddContacts(ctxSess, tt.req, models.UserSession{UserID: "61f9c63d5496b70001e03d23"})
			if (err != nil) != tt.wantErr {
				t.Errorf("s.AddContacts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(resp, tt.wantResp) {
				t.Errorf("s.AddContacts() = %v, want %v", resp, tt.wantResp)
			}
		})
	}
}

func TestGetContactsList(t *testing.T) {
	type mockRepo struct {
		GetListAdd          []*domain.Contacts
		ContactGetListError error
	}

	tests := []struct {
		name     string
		user     models.UserSession
		repo     mockRepo
		wantResp []*domain.Contacts
		wantErr  bool
		wantData bool
	}{
		{
			name:    "error database",
			wantErr: true,
			repo: mockRepo{
				ContactGetListError: errors.New("error database"),
			},
		},
		{
			name: "success",
			repo: mockRepo{
				GetListAdd: []*domain.Contacts{
					{
						ID:     "61fc3cedf5922d9d54954369",
						UserID: "61f9c63d5496b70001e03d23",
						Name:   "test",
						Email:  "test@test.com",
					},
				},
			},
			wantResp: []*domain.Contacts{
				{
					ID:     "61fc3cedf5922d9d54954369",
					UserID: "61f9c63d5496b70001e03d23",
					Name:   "test",
					Email:  "test@test.com",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRedisRepo := &mockRepoRedis.Repository{}
			mockContactsRepo := &mockRepoContacts.Repository{}
			mockContactsRepo.On("GetList", mock.Anything).Return(tt.repo.GetListAdd, tt.repo.ContactGetListError)

			s := New(mockContactsRepo, mockRedisRepo)

			resp, err := s.GetContactsList(ctxSess, models.UserSession{UserID: "61f9c63d5496b70001e03d23"})
			if (err != nil) != tt.wantErr {
				t.Errorf("s.GetContactsList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(resp, tt.wantResp) {
				t.Errorf("s.GetContactsList() = %v, want %v", resp, tt.wantResp)
			}
		})
	}
}

func TestGetRecentContactsList(t *testing.T) {
	type mockRepo struct {
		GetPushCache      []string
		GetPushCacheError error
	}

	tests := []struct {
		name     string
		user     models.UserSession
		repo     mockRepo
		wantResp []string
		wantErr  bool
		wantData bool
	}{
		{
			name:    "error database",
			wantErr: true,
			repo: mockRepo{
				GetPushCacheError: constants.ErrorDataNotFound,
			},
		},
		{
			name: "success",
			repo: mockRepo{
				GetPushCache: []string{"test@gmail.com", "test2@gmail.com"},
			},
			wantResp: []string{"test@gmail.com", "test2@gmail.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContactsRepo := &mockRepoContacts.Repository{}
			mockRedisRepo := &mockRepoRedis.Repository{}
			mockRedisRepo.On("GetPushCache", mock.Anything, mock.Anything, mock.Anything).Return(tt.repo.GetPushCache, tt.repo.GetPushCacheError)

			s := New(mockContactsRepo, mockRedisRepo)

			resp, err := s.GetRecentContactsList(ctxSess, models.UserSession{UserID: "61f9c63d5496b70001e03d23"})
			if (err != nil) != tt.wantErr {
				t.Errorf("s.GetRecentContactsList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(resp, tt.wantResp) {
				t.Errorf("s.GetRecentContactsList() = %v, want %v", resp, tt.wantResp)
			}
		})
	}
}

func TestUpdateContacts(t *testing.T) {
	type mockRepo struct {
		ContactGet      *domain.Contacts
		ContactGetError error

		ContactUpdate      *domain.Contacts
		ContactUpdateError error
	}

	tests := []struct {
		name     string
		req      ContactUpdateReq
		user     models.UserSession
		repo     mockRepo
		wantResp *domain.Contacts
		wantErr  bool
		wantData bool
	}{
		{
			name:    "error data not found",
			req:     ContactUpdateReq{ID: "61fc3cedf5922d9d54954369", Email: "test@test.com", Name: "test"},
			wantErr: true,
			repo: mockRepo{
				ContactGetError: errors.New("data not found"),
			},
		},
		{
			name:    "email not match",
			req:     ContactUpdateReq{ID: "61fc3cedf5922d9d54954369", Email: "test123@test.com", Name: "test"},
			wantErr: true,
			repo: mockRepo{
				ContactGet: &domain.Contacts{
					ID:     "61fc3cedf5922d9d54954369",
					UserID: "61f9c63d5496b70001e03d23",
					Name:   "test",
					Email:  "test@test.com",
				},
			},
		},
		{
			name:    "error update contacts",
			req:     ContactUpdateReq{ID: "61fc3cedf5922d9d54954369", Email: "test@test.com", Name: "test"},
			wantErr: true,
			repo: mockRepo{
				ContactGet: &domain.Contacts{
					ID:     "61fc3cedf5922d9d54954369",
					UserID: "61f9c63d5496b70001e03d23",
					Name:   "test",
					Email:  "test@test.com",
				},
				ContactUpdateError: errors.New("error database"),
			},
		},
		{
			name: "success",
			req:  ContactUpdateReq{ID: "61fc3cedf5922d9d54954369", Email: "test@test.com", Name: "test 123"},
			repo: mockRepo{
				ContactGet: &domain.Contacts{
					ID:     "61fc3cedf5922d9d54954369",
					UserID: "61f9c63d5496b70001e03d23",
					Name:   "test",
					Email:  "test@test.com",
				},
				ContactUpdate: &domain.Contacts{
					ID:     "61fc3cedf5922d9d54954369",
					UserID: "61f9c63d5496b70001e03d23",
					Name:   "test 123",
					Email:  "test@test.com",
				},
			},
			wantResp: &domain.Contacts{
				ID:     "61fc3cedf5922d9d54954369",
				UserID: "61f9c63d5496b70001e03d23",
				Name:   "test 123",
				Email:  "test@test.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRedisRepo := &mockRepoRedis.Repository{}
			mockContactsRepo := &mockRepoContacts.Repository{}
			mockContactsRepo.On("Get", mock.Anything).Return(tt.repo.ContactGet, tt.repo.ContactGetError)
			mockContactsRepo.On("Update", mock.Anything).Return(tt.repo.ContactUpdate, tt.repo.ContactUpdateError)

			s := New(mockContactsRepo, mockRedisRepo)

			resp, err := s.UpdateContacts(ctxSess, tt.req, models.UserSession{UserID: "61f9c63d5496b70001e03d23"})
			if (err != nil) != tt.wantErr {
				t.Errorf("s.UpdateContacts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(resp, tt.wantResp) {
				t.Errorf("s.UpdateContacts() = %v, want %v", resp, tt.wantResp)
			}
		})
	}
}

func TestDeleteContacts(t *testing.T) {
	type mockRepo struct {
		ContactGet      *domain.Contacts
		ContactGetError error

		ContactDeleteError error
	}

	tests := []struct {
		name     string
		req      string
		user     models.UserSession
		repo     mockRepo
		wantErr  bool
		wantData bool
	}{
		{
			name:    "error data not found",
			req:     "61fc3cedf5922d9d54954369",
			wantErr: true,
			repo: mockRepo{
				ContactGetError: errors.New("data not found"),
			},
		},
		{
			name:    "user ID not match",
			req:     "61fc3cedf5922d9d54954369",
			wantErr: true,
			repo: mockRepo{
				ContactGet: &domain.Contacts{
					ID:     "61fc3cedf5922d9d54954369",
					UserID: "61f9c63d5496b70001e03d24",
					Name:   "test",
					Email:  "test@test.com",
				},
			},
		},
		{
			name:    "error delete contacts",
			req:     "61fc3cedf5922d9d54954369",
			wantErr: true,
			repo: mockRepo{
				ContactGet: &domain.Contacts{
					ID:     "61fc3cedf5922d9d54954369",
					UserID: "61f9c63d5496b70001e03d23",
					Name:   "test",
					Email:  "test@test.com",
				},
				ContactDeleteError: errors.New("error database"),
			},
		},
		{
			name: "success",
			req:  "61fc3cedf5922d9d54954369",
			repo: mockRepo{
				ContactGet: &domain.Contacts{
					ID:     "61fc3cedf5922d9d54954369",
					UserID: "61f9c63d5496b70001e03d23",
					Name:   "test",
					Email:  "test@test.com",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRedisRepo := &mockRepoRedis.Repository{}
			mockContactsRepo := &mockRepoContacts.Repository{}
			mockContactsRepo.On("Get", mock.Anything).Return(tt.repo.ContactGet, tt.repo.ContactGetError)
			mockContactsRepo.On("Delete", mock.Anything).Return(tt.repo.ContactDeleteError)

			s := New(mockContactsRepo, mockRedisRepo)

			err := s.DeleteContacts(ctxSess, tt.req, models.UserSession{UserID: "61f9c63d5496b70001e03d23"})
			if (err != nil) != tt.wantErr {
				t.Errorf("s.DeleteContacts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
