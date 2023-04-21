package aways

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Away struct {
		ID              primitive.ObjectID `json:"_id" bson:"_id"`
		UserID          primitive.ObjectID `bson:"user_id" json:"-"`
		Title           string             `bson:"title" json:"title"`
		ActivateAllow   time.Time          `bson:"activate_allow" json:"activate_allow"`
		DeactivateAllow time.Time          `bson:"deactivate_allow" json:"deactivate_allow"`
		IsEnabled       bool               `bson:"is_enabled" json:"is_enabled"`
		Repeat          []string           `bson:"repeat" json:"repeat"`
		AllDay          bool               `bson:"all_day" json:"all_day"`
		Message         string             `bson:"message" json:"message"`
		KeyDuration     string             `bson:"key_duration" json:"key_duration"`
		AllowedUsers    []Allowed          `bson:"allowed_users" json:"allowed_users"`
		AllowedSubjects []Allowed          `bson:"allowed_subjects" json:"allowed_subjects"`
		AllowedContacts []string           `bson:"allowed_contacts" json:"allowed_contacts"`
		AllowedKeywords []string           `bson:"allowed_keywords" json:"allowed_keywords"`
	}
	Allowed struct {
		ID        string `bson:"id" json:"id"`
		Name      string `bson:"name" json:"name"`
		Activated bool   `bson:"activated" json:"activated"`
	}
)
