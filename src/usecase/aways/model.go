package aways

import (
	"time"

	"github.com/cloudsrc/api.awaymail.v1.go/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateAwayReq struct {
	Title           string    `json:"title" binding:"required"`
	ActivateAllow   time.Time `json:"activate_allow,omitempty"`
	DeactivateAllow time.Time `json:"deactivate_allow,omitempty"`
	Repeat          []string  `json:"repeat"`
	IsEnabled       bool      `json:"is_enabled"`
	AllDay          bool      `json:"all_day"`
	Message         string    `json:"message"`
	KeyDuration     string    `json:"key_duration"`
	AllowedContacts []string  `json:"allowed_contacts"`
	AllowedKeywords []string  `json:"allowed_keywords"`
}

type UpdateAwayReq struct {
	ID              primitive.ObjectID `json:"id" binding:"required"`
	Title           string             `json:"title" binding:"required"`
	ActivateAllow   time.Time          `json:"activate_allow,omitempty"`
	DeactivateAllow time.Time          `json:"deactivate_allow,omitempty"`
	Repeat          []string           `json:"repeat"`
	IsEnabled       bool               `json:"is_enabled"`
	AllDay          bool               `json:"all_day"`
	Message         string             `json:"message"`
	KeyDuration     string             `json:"key_duration"`
	AllowedContacts []string           `json:"allowed_contacts"`
	AllowedKeywords []string           `json:"allowed_keywords"`
}

type EnableAwayReq struct {
	IsEnabled bool `json:"is_enabled" form:"is_enabled"`
}

type CreateAwayResp struct {
	ID              primitive.ObjectID `json:"id,omitempty"`
	Title           string             `json:"title" binding:"required"`
	ActivateAllow   time.Time          `json:"activate_allow"`
	DeactivateAllow time.Time          `json:"deactivate_allow"`
	Repeat          []string           `json:"repeat"`
	IsEnabled       bool               `json:"is_enabled"`
	AllDay          bool               `json:"all_day"`
	Message         string             `json:"message"`
	KeyDuration     string             `json:"key_duration"`
	AllowedUsers    []models.Allowed   `json:"allowed_users"`
	AllowedSubjects []models.Allowed   `json:"allowed_subjects"`
	AllowedContacts []string           `json:"allowed_contacts"`
	AllowedKeywords []string           `json:"allowed_keywords"`
}
