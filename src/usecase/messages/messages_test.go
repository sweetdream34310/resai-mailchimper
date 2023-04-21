package messages_test

import (
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/logger"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
	"google.golang.org/api/gmail/v1"

	"github.com/cloudsrc/api.awaymail.v1.go/libs"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/constants"
	mockRepoAways "github.com/cloudsrc/api.awaymail.v1.go/src/shared/mock/mongo/aways"
	mockRepoContacts "github.com/cloudsrc/api.awaymail.v1.go/src/shared/mock/mongo/contacts"
	mockRepoMessages "github.com/cloudsrc/api.awaymail.v1.go/src/shared/mock/mongo/messages"
	mockRedis "github.com/cloudsrc/api.awaymail.v1.go/src/shared/mock/redis"
	mockWrapperGoogle "github.com/cloudsrc/api.awaymail.v1.go/src/shared/mock/wrapper/google"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/models"
	. "github.com/cloudsrc/api.awaymail.v1.go/src/usecase/messages"
)

var (
	ctxSess = context.New(logger.NewNoopLogger())
)

func TestSentMessage(t *testing.T) {
	type mockRepo struct {
		PushCacheError error
	}

	type mockWrapper struct {
		SendMessageError               error
		SendMessageWithAttachmentError error
	}

	tests := []struct {
		name     string
		req      SendMessageReq
		user     models.UserSession
		repo     mockRepo
		wrapper  mockWrapper
		wantErr  bool
		wantData bool
	}{
		{
			name:    "error google - SendMessage",
			req:     SendMessageReq{Message: "test", Subject: "test subject", To: "test_to@gmail.com"},
			wantErr: true,
			wrapper: mockWrapper{
				SendMessageError: constants.ErrorGeneral,
			},
		},
		{
			name:    "error google - SendMessageWithAttachment",
			req:     SendMessageReq{Message: "test", Subject: "test subject", To: "test_to@gmail.com", AttachmentsURL: []string{"http://google.com"}},
			wantErr: true,
			wrapper: mockWrapper{
				SendMessageWithAttachmentError: constants.ErrorGeneral,
			},
		},
		{
			name: "success - SendMessage",
			req:  SendMessageReq{Message: "test", Subject: "test subject", To: "test_to@gmail.com"},
		},
		{
			name: "success - SendMessageWithAttachment",
			req:  SendMessageReq{Message: "test", Subject: "test subject", To: "test_to@gmail.com", AttachmentsURL: []string{"http://google.com"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGoogleWrapper := &mockWrapperGoogle.Wrapper{}
			mockGoogleWrapper.On("SendMessage", mock.Anything, mock.Anything).Return(tt.wrapper.SendMessageError)
			mockGoogleWrapper.On("SendMessageWithAttachment", mock.Anything, mock.Anything).Return(tt.wrapper.SendMessageWithAttachmentError)
			mockAwaysRepo := &mockRepoAways.Repository{}
			mockContactsRepo := &mockRepoContacts.Repository{}
			mockMessagesRepo := &mockRepoMessages.Repository{}
			mockRedisRepo := &mockRedis.Repository{}
			mockRedisRepo.On("PushCache", mock.Anything, mock.Anything).Return(tt.repo.PushCacheError)

			s := New(mockAwaysRepo, mockMessagesRepo, mockContactsRepo, mockGoogleWrapper, mockRedisRepo, &libs.RabbitClient{})

			err := s.SentMessage(ctxSess, models.UserSession{UserID: "61f9c63d5496b70001e03d23"}, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("s.GetAuthToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUpdateMessage(t *testing.T) {
	type mockWrapper struct {
		UpdateMessage      *gmail.Message
		UpdateMessageError error
	}

	tests := []struct {
		name     string
		req      UpdateRequest
		user     models.UserSession
		wrapper  mockWrapper
		wantResp *gmail.Message
		wantErr  bool
		wantData bool
	}{
		{
			name:    "error google - SendMessage",
			req:     UpdateRequest{},
			wantErr: true,
			wrapper: mockWrapper{
				UpdateMessageError: constants.ErrorGeneral,
			},
		},
		{
			name: "success - SendMessageWithAttachment",
			req:  UpdateRequest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGoogleWrapper := &mockWrapperGoogle.Wrapper{}
			mockGoogleWrapper.On("UpdateMessage", mock.Anything, mock.Anything, mock.Anything).Return(tt.wrapper.UpdateMessage, tt.wrapper.UpdateMessageError)
			mockAwaysRepo := &mockRepoAways.Repository{}
			mockContactsRepo := &mockRepoContacts.Repository{}
			mockMessagesRepo := &mockRepoMessages.Repository{}
			mockRedisRepo := &mockRedis.Repository{}

			s := New(mockAwaysRepo, mockMessagesRepo, mockContactsRepo, mockGoogleWrapper, mockRedisRepo, &libs.RabbitClient{})

			resp, err := s.UpdateMessage(ctxSess, "61f9c63d5496b70001e03d23", tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("s.UpdateMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(resp, tt.wantResp) {
				t.Errorf("s.GetAuthToken() = %v, want %v", resp, tt.wantResp)
			}
		})
	}
}

func TestDeleteMessage(t *testing.T) {
	type mockWrapper struct {
		DeleteMessageError error
	}

	tests := []struct {
		name     string
		user     models.UserSession
		wrapper  mockWrapper
		wantErr  bool
		wantData bool
	}{
		{
			name:    "error google - SendMessage",
			wantErr: true,
			wrapper: mockWrapper{
				DeleteMessageError: constants.ErrorGeneral,
			},
		},
		{
			name: "success - SendMessageWithAttachment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGoogleWrapper := &mockWrapperGoogle.Wrapper{}
			mockGoogleWrapper.On("DeleteMessage", mock.Anything, mock.Anything).Return(tt.wrapper.DeleteMessageError)
			mockAwaysRepo := &mockRepoAways.Repository{}
			mockContactsRepo := &mockRepoContacts.Repository{}
			mockMessagesRepo := &mockRepoMessages.Repository{}
			mockRedisRepo := &mockRedis.Repository{}

			s := New(mockAwaysRepo, mockMessagesRepo, mockContactsRepo, mockGoogleWrapper, mockRedisRepo, &libs.RabbitClient{})

			err := s.DeleteMessage(ctxSess, "61f9c63d5496b70001e03d23")
			if (err != nil) != tt.wantErr {
				t.Errorf("s.UpdateMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
