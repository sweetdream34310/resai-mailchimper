package aways_test

import (
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/logger"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	domain "github.com/cloudsrc/api.awaymail.v1.go/src/domain/aways"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/constants"
	mockRepoAways "github.com/cloudsrc/api.awaymail.v1.go/src/shared/mock/mongo/aways"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/models"
	. "github.com/cloudsrc/api.awaymail.v1.go/src/usecase/aways"
)

var (
	ctxSess = context.New(logger.NewNoopLogger())
)

func TestCreatAways(t *testing.T) {
	type mockRepo struct {
		CreateAway      *domain.Away
		CreateAwayError error
	}

	tests := []struct {
		name     string
		req      *CreateAwayReq
		user     models.UserSession
		repo     mockRepo
		wantResp *CreateAwayResp
		wantErr  bool
		wantData bool
	}{
		{
			name:    "error Request - all time schedule is empty",
			req:     &CreateAwayReq{Title: "Test aways"},
			wantErr: true,
		},
		{
			name:    "error Request - ActivateAllow and DeactivateAllow is empty",
			req:     &CreateAwayReq{Title: "Test aways", Repeat: []string{"Sunday", "Monday"}, AllDay: true},
			wantErr: true,
		},
		{
			name:    "error Request - invalid repeat date",
			req:     &CreateAwayReq{Title: "Test aways", Repeat: []string{"Sun", "Monday"}, AllDay: true, ActivateAllow: time.Now(), DeactivateAllow: time.Now().Add(3 * time.Hour)},
			wantErr: true,
		},
		{
			name:    "error create away - error database",
			req:     &CreateAwayReq{Title: "Test aways", Repeat: []string{"Sunday", "Monday"}, AllDay: true, ActivateAllow: time.Now(), DeactivateAllow: time.Now().Add(3 * time.Hour)},
			wantErr: true,
			repo: mockRepo{
				CreateAwayError: constants.ErrorDatabase,
			},
		},
		{
			name: "success",
			req:  &CreateAwayReq{Title: "Test aways", Repeat: []string{"Sunday", "Sunday", "Monday"}, AllDay: true, ActivateAllow: time.Now(), DeactivateAllow: time.Now().Add(3 * time.Hour)},
			repo: mockRepo{
				CreateAway: &domain.Away{
					ID:              "62056d328945fc00018c49b0",
					Title:           "1 day away",
					IsEnabled:       true,
					ActivateAllow:   time.Now(),
					DeactivateAllow: time.Now().Add(3 * time.Hour),
					Repeat:          []string{"Sunday", "Monday"},
					AllDay:          true,
				},
			},
			wantResp: &CreateAwayResp{
				ID:              "62056d328945fc00018c49b0",
				Title:           "1 day away",
				IsEnabled:       true,
				ActivateAllow:   time.Now(),
				DeactivateAllow: time.Now().Add(3 * time.Hour),
				Repeat:          []string{"Sunday", "Monday"},
				AllDay:          true,
			},
		},
		{
			name: "success with duplicate repeat date",
			req:  &CreateAwayReq{Title: "Test aways", Repeat: []string{"Sunday", "Monday"}, AllDay: true, ActivateAllow: time.Now(), DeactivateAllow: time.Now().Add(3 * time.Hour)},
			repo: mockRepo{
				CreateAway: &domain.Away{
					ID:              "62056d328945fc00018c49b0",
					Title:           "1 day away",
					IsEnabled:       true,
					ActivateAllow:   time.Now(),
					DeactivateAllow: time.Now().Add(3 * time.Hour),
					Repeat:          []string{"Sunday", "Monday"},
					AllDay:          true,
				},
			},
			wantResp: &CreateAwayResp{
				ID:              "62056d328945fc00018c49b0",
				Title:           "1 day away",
				IsEnabled:       true,
				ActivateAllow:   time.Now(),
				DeactivateAllow: time.Now().Add(3 * time.Hour),
				Repeat:          []string{"Sunday", "Monday"},
				AllDay:          true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAwaysRepo := &mockRepoAways.Repository{}
			mockAwaysRepo.On("CreateAway", mock.Anything).Return(tt.repo.CreateAway, tt.repo.CreateAwayError)

			s := New(mockAwaysRepo)

			resp, err := s.CreateAway(ctxSess, models.UserSession{UserID: "61f9c63d5496b70001e03d23"}, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("s.CreateAway() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantResp != nil {
				resp.ActivateAllow = tt.wantResp.ActivateAllow
				resp.DeactivateAllow = tt.wantResp.DeactivateAllow
			}
			if !reflect.DeepEqual(resp, tt.wantResp) {
				t.Errorf("s.CreateAway() = %v, want %v", resp, tt.wantResp)
			}
		})
	}
}

func TestGetAwayList(t *testing.T) {
	type mockRepo struct {
		GetAwayList      []*domain.Away
		GetAwayListError error
	}

	tests := []struct {
		name     string
		user     models.UserSession
		repo     mockRepo
		wantResp []*CreateAwayResp
		wantErr  bool
		wantData bool
	}{
		{
			name:    "error database",
			wantErr: true,
			repo: mockRepo{
				GetAwayListError: constants.ErrorDataNotFound,
			},
		},
		{
			name: "success",
			repo: mockRepo{
				GetAwayList: []*domain.Away{
					{
						ID:              "61fc3cedf5922d9d54954369",
						UserID:          "61f9c63d5496b70001e03d23",
						Title:           "1 day away",
						IsEnabled:       true,
						ActivateAllow:   time.Now(),
						DeactivateAllow: time.Now().Add(3 * time.Hour),
						Repeat:          []string{"Sunday", "Monday"},
						AllDay:          true,
					},
				},
			},
			wantResp: []*CreateAwayResp{
				{
					ID:              "61fc3cedf5922d9d54954369",
					Title:           "1 day away",
					IsEnabled:       true,
					ActivateAllow:   time.Now(),
					DeactivateAllow: time.Now().Add(3 * time.Hour),
					Repeat:          []string{"Sunday", "Monday"},
					AllDay:          true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAwaysRepo := &mockRepoAways.Repository{}
			mockAwaysRepo.On("GetAwayList", mock.Anything).Return(tt.repo.GetAwayList, tt.repo.GetAwayListError)

			s := New(mockAwaysRepo)

			resp, err := s.GetAwayList(ctxSess, models.UserSession{UserID: "61f9c63d5496b70001e03d23"})
			if (err != nil) != tt.wantErr {
				t.Errorf("s.CreateAway() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantResp != nil {
				for key, _ := range resp {
					resp[key].ActivateAllow = tt.wantResp[key].ActivateAllow
					resp[key].DeactivateAllow = tt.wantResp[key].DeactivateAllow
				}
			}
			if !reflect.DeepEqual(resp, tt.wantResp) {
				t.Errorf("s.CreateAway() = %v, want %v", resp, tt.wantResp)
			}
		})
	}
}

func TestGetAway(t *testing.T) {
	type mockRepo struct {
		GetAway      *domain.Away
		GetAwayError error
	}

	tests := []struct {
		name     string
		user     models.UserSession
		repo     mockRepo
		wantResp *CreateAwayResp
		wantErr  bool
		wantData bool
	}{
		{
			name:    "error database",
			wantErr: true,
			repo: mockRepo{
				GetAwayError: constants.ErrorDataNotFound,
			},
		},
		{
			name: "success",
			repo: mockRepo{
				GetAway: &domain.Away{
					ID:              "61fc3cedf5922d9d54954369",
					UserID:          "61f9c63d5496b70001e03d23",
					Title:           "1 day away",
					IsEnabled:       true,
					ActivateAllow:   time.Now(),
					DeactivateAllow: time.Now().Add(3 * time.Hour),
					Repeat:          []string{"Sunday", "Monday"},
					AllDay:          true,
				},
			},
			wantResp: &CreateAwayResp{
				ID:              "61fc3cedf5922d9d54954369",
				Title:           "1 day away",
				IsEnabled:       true,
				ActivateAllow:   time.Now(),
				DeactivateAllow: time.Now().Add(3 * time.Hour),
				Repeat:          []string{"Sunday", "Monday"},
				AllDay:          true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAwaysRepo := &mockRepoAways.Repository{}
			mockAwaysRepo.On("GetAway", mock.Anything).Return(tt.repo.GetAway, tt.repo.GetAwayError)

			s := New(mockAwaysRepo)

			resp, err := s.GetAway(ctxSess, "61fc3cedf5922d9d54954369")
			if (err != nil) != tt.wantErr {
				t.Errorf("s.GetAway() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantResp != nil {
				resp.ActivateAllow = tt.wantResp.ActivateAllow
				resp.DeactivateAllow = tt.wantResp.DeactivateAllow
			}
			if !reflect.DeepEqual(resp, tt.wantResp) {
				t.Errorf("s.GetAway() = %v, want %v", resp, tt.wantResp)
			}
		})
	}
}

func TestUpdateAway(t *testing.T) {
	type mockRepo struct {
		GetAway      *domain.Away
		GetAwayError error

		UpdateAway      *domain.Away
		UpdateAwayError error
	}

	tests := []struct {
		name     string
		user     models.UserSession
		req      *UpdateAwayReq
		repo     mockRepo
		wantResp *CreateAwayResp
		wantErr  bool
		wantData bool
	}{
		{
			name:    "error database",
			req:     &UpdateAwayReq{ID: "61fc3cedf5922d9d54954369", IsEnabled: false},
			wantErr: true,
			repo: mockRepo{
				GetAwayError: constants.ErrorDataNotFound,
			},
		},
		{
			name:    "error not authorized",
			req:     &UpdateAwayReq{ID: "61fc3cedf5922d9d54954369", IsEnabled: false},
			wantErr: true,
			repo: mockRepo{
				GetAway: &domain.Away{
					ID:              "61fc3cedf5922d9d54954369",
					UserID:          "62001350f5922d7a8ee34cd9",
					Title:           "1 day away",
					IsEnabled:       true,
					ActivateAllow:   time.Now(),
					DeactivateAllow: time.Now().Add(3 * time.Hour),
					Repeat:          []string{"Sunday", "Monday"},
					AllDay:          true,
				},
			},
		},
		{
			name:    "error update database",
			req:     &UpdateAwayReq{ID: "61fc3cedf5922d9d54954369", IsEnabled: false},
			wantErr: true,
			repo: mockRepo{
				GetAway: &domain.Away{
					ID:              "61fc3cedf5922d9d54954369",
					UserID:          "61f9c63d5496b70001e03d23",
					Title:           "1 day away",
					IsEnabled:       true,
					ActivateAllow:   time.Now(),
					DeactivateAllow: time.Now().Add(3 * time.Hour),
					Repeat:          []string{"Sunday", "Monday"},
					AllDay:          true,
				},
				UpdateAwayError: constants.ErrorDatabase,
			},
		},
		{
			name: "success",
			req:  &UpdateAwayReq{ID: "61fc3cedf5922d9d54954369", IsEnabled: false},
			repo: mockRepo{
				GetAway: &domain.Away{
					ID:              "61fc3cedf5922d9d54954369",
					UserID:          "61f9c63d5496b70001e03d23",
					Title:           "1 day away",
					IsEnabled:       true,
					ActivateAllow:   time.Now(),
					DeactivateAllow: time.Now().Add(3 * time.Hour),
					Repeat:          []string{"Sunday", "Monday"},
					AllDay:          true,
				},
				UpdateAway: &domain.Away{
					ID:              "61fc3cedf5922d9d54954369",
					UserID:          "61f9c63d5496b70001e03d23",
					Title:           "1 day away",
					IsEnabled:       false,
					ActivateAllow:   time.Now(),
					DeactivateAllow: time.Now().Add(3 * time.Hour),
					Repeat:          []string{"Sunday", "Monday"},
					AllDay:          true,
				},
			},
			wantResp: &CreateAwayResp{
				ID:              "61fc3cedf5922d9d54954369",
				Title:           "1 day away",
				IsEnabled:       false,
				ActivateAllow:   time.Now(),
				DeactivateAllow: time.Now().Add(3 * time.Hour),
				Repeat:          []string{"Sunday", "Monday"},
				AllDay:          true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAwaysRepo := &mockRepoAways.Repository{}
			mockAwaysRepo.On("GetAway", mock.Anything).Return(tt.repo.GetAway, tt.repo.GetAwayError)
			mockAwaysRepo.On("UpdateAway", mock.Anything).Return(tt.repo.UpdateAway, tt.repo.UpdateAwayError)

			s := New(mockAwaysRepo)

			resp, err := s.UpdateAway(ctxSess, models.UserSession{UserID: "61f9c63d5496b70001e03d23"}, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("s.UpdateAway() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantResp != nil {
				resp.ActivateAllow = tt.wantResp.ActivateAllow
				resp.DeactivateAllow = tt.wantResp.DeactivateAllow
			}
			if !reflect.DeepEqual(resp, tt.wantResp) {
				t.Errorf("s.UpdateAway() = %v, want %v", resp, tt.wantResp)
			}
		})
	}
}

func TestEnableAway(t *testing.T) {
	type mockRepo struct {
		GetAway      *domain.Away
		GetAwayError error

		UpdateAway      *domain.Away
		UpdateAwayError error
	}

	tests := []struct {
		name     string
		user     models.UserSession
		req      bool
		repo     mockRepo
		wantResp *CreateAwayResp
		wantErr  bool
		wantData bool
	}{
		{
			name:    "error database",
			req:     true,
			wantErr: true,
			repo: mockRepo{
				GetAwayError: constants.ErrorDataNotFound,
			},
		},
		{
			name:    "error not authorized",
			req:     false,
			wantErr: true,
			repo: mockRepo{
				GetAway: &domain.Away{
					ID:              "61fc3cedf5922d9d54954369",
					UserID:          "62001350f5922d7a8ee34cd9",
					Title:           "1 day away",
					IsEnabled:       true,
					ActivateAllow:   time.Now(),
					DeactivateAllow: time.Now().Add(3 * time.Hour),
					Repeat:          []string{"Sunday", "Monday"},
					AllDay:          true,
				},
			},
		},
		{
			name:    "error update database",
			req:     true,
			wantErr: true,
			repo: mockRepo{
				GetAway: &domain.Away{
					ID:              "61fc3cedf5922d9d54954369",
					UserID:          "61f9c63d5496b70001e03d23",
					Title:           "1 day away",
					IsEnabled:       true,
					ActivateAllow:   time.Now(),
					DeactivateAllow: time.Now().Add(3 * time.Hour),
					Repeat:          []string{"Sunday", "Monday"},
					AllDay:          true,
				},
				UpdateAwayError: constants.ErrorDatabase,
			},
		},
		{
			name: "success",
			req:  false,
			repo: mockRepo{
				GetAway: &domain.Away{
					ID:              "61fc3cedf5922d9d54954369",
					UserID:          "61f9c63d5496b70001e03d23",
					Title:           "1 day away",
					IsEnabled:       true,
					ActivateAllow:   time.Now(),
					DeactivateAllow: time.Now().Add(3 * time.Hour),
					Repeat:          []string{"Sunday", "Monday"},
					AllDay:          true,
				},
				UpdateAway: &domain.Away{
					ID:              "61fc3cedf5922d9d54954369",
					UserID:          "61f9c63d5496b70001e03d23",
					Title:           "1 day away",
					IsEnabled:       false,
					ActivateAllow:   time.Now(),
					DeactivateAllow: time.Now().Add(3 * time.Hour),
					Repeat:          []string{"Sunday", "Monday"},
					AllDay:          true,
				},
			},
			wantResp: &CreateAwayResp{
				ID:              "61fc3cedf5922d9d54954369",
				Title:           "1 day away",
				IsEnabled:       false,
				ActivateAllow:   time.Now(),
				DeactivateAllow: time.Now().Add(3 * time.Hour),
				Repeat:          []string{"Sunday", "Monday"},
				AllDay:          true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAwaysRepo := &mockRepoAways.Repository{}
			mockAwaysRepo.On("GetAway", mock.Anything).Return(tt.repo.GetAway, tt.repo.GetAwayError)
			mockAwaysRepo.On("UpdateAway", mock.Anything).Return(tt.repo.UpdateAway, tt.repo.UpdateAwayError)

			s := New(mockAwaysRepo)

			resp, err := s.EnableAway(ctxSess, models.UserSession{UserID: "61f9c63d5496b70001e03d23"}, tt.req, "61fc3cedf5922d9d54954369")
			if (err != nil) != tt.wantErr {
				t.Errorf("s.EnableAway() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantResp != nil {
				resp.ActivateAllow = tt.wantResp.ActivateAllow
				resp.DeactivateAllow = tt.wantResp.DeactivateAllow
			}
			if !reflect.DeepEqual(resp, tt.wantResp) {
				t.Errorf("s.EnableAway() = %v, want %v", resp, tt.wantResp)
			}
		})
	}
}

func TestEnableAllAway(t *testing.T) {
	type mockRepo struct {
		UpdateAwayMode      int
		UpdateAwayModeError error
	}

	tests := []struct {
		name     string
		user     models.UserSession
		req      bool
		repo     mockRepo
		wantResp int
		wantErr  bool
		wantData bool
	}{
		{
			name:    "error update database",
			req:     true,
			wantErr: true,
			repo: mockRepo{
				UpdateAwayModeError: constants.ErrorDatabase,
			},
		},
		{
			name: "success",
			req:  false,
			repo: mockRepo{
				UpdateAwayMode: 1,
			},
			wantResp: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAwaysRepo := &mockRepoAways.Repository{}
			mockAwaysRepo.On("UpdateAwayMode", mock.Anything, mock.Anything).Return(tt.repo.UpdateAwayMode, tt.repo.UpdateAwayModeError)

			s := New(mockAwaysRepo)

			resp, err := s.EnableAllAway(ctxSess, models.UserSession{UserID: "61f9c63d5496b70001e03d23"}, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("s.EnableAway() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(resp, tt.wantResp) {
				t.Errorf("s.EnableAway() = %v, want %v", resp, tt.wantResp)
			}
		})
	}
}

func TestDeleteAway(t *testing.T) {
	type mockRepo struct {
		GetAway      *domain.Away
		GetAwayError error

		DeleteAwayError error
	}

	tests := []struct {
		name     string
		user     models.UserSession
		req      bool
		repo     mockRepo
		wantResp *CreateAwayResp
		wantErr  bool
		wantData bool
	}{
		{
			name:    "error database",
			req:     true,
			wantErr: true,
			repo: mockRepo{
				GetAwayError: constants.ErrorDataNotFound,
			},
		},
		{
			name:    "error not authorized",
			req:     false,
			wantErr: true,
			repo: mockRepo{
				GetAway: &domain.Away{
					ID:              "61fc3cedf5922d9d54954369",
					UserID:          "62001350f5922d7a8ee34cd9",
					Title:           "1 day away",
					IsEnabled:       true,
					ActivateAllow:   time.Now(),
					DeactivateAllow: time.Now().Add(3 * time.Hour),
					Repeat:          []string{"Sunday", "Monday"},
					AllDay:          true,
				},
			},
		},
		{
			name:    "error update database",
			req:     true,
			wantErr: true,
			repo: mockRepo{
				GetAway: &domain.Away{
					ID:              "61fc3cedf5922d9d54954369",
					UserID:          "61f9c63d5496b70001e03d23",
					Title:           "1 day away",
					IsEnabled:       true,
					ActivateAllow:   time.Now(),
					DeactivateAllow: time.Now().Add(3 * time.Hour),
					Repeat:          []string{"Sunday", "Monday"},
					AllDay:          true,
				},
				DeleteAwayError: constants.ErrorDatabase,
			},
		},
		{
			name: "success",
			req:  false,
			repo: mockRepo{
				GetAway: &domain.Away{
					ID:              "61fc3cedf5922d9d54954369",
					UserID:          "61f9c63d5496b70001e03d23",
					Title:           "1 day away",
					IsEnabled:       true,
					ActivateAllow:   time.Now(),
					DeactivateAllow: time.Now().Add(3 * time.Hour),
					Repeat:          []string{"Sunday", "Monday"},
					AllDay:          true,
				},
			},
			wantResp: &CreateAwayResp{
				ID:              "61fc3cedf5922d9d54954369",
				Title:           "1 day away",
				IsEnabled:       false,
				ActivateAllow:   time.Now(),
				DeactivateAllow: time.Now().Add(3 * time.Hour),
				Repeat:          []string{"Sunday", "Monday"},
				AllDay:          true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAwaysRepo := &mockRepoAways.Repository{}
			mockAwaysRepo.On("GetAway", mock.Anything).Return(tt.repo.GetAway, tt.repo.GetAwayError)
			mockAwaysRepo.On("DeleteAway", mock.Anything).Return(tt.repo.DeleteAwayError)

			s := New(mockAwaysRepo)

			err := s.DeleteAway(ctxSess, models.UserSession{UserID: "61f9c63d5496b70001e03d23"}, "61fc3cedf5922d9d54954369")
			if (err != nil) != tt.wantErr {
				t.Errorf("s.EnableAway() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
