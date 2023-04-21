package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	Away struct {
		ID              bson.ObjectId `json:"_id" bson:"_id"`
		DB              *mgo.Database `bson:"-" json:"-"`
		UserID          bson.ObjectId `bson:"user_id" json:"-"`
		Title           string        `bson:"title" json:"title"`
		ActivateAllow   time.Time     `bson:"activate_allow" json:"activate_allow"`
		DeactivateAllow time.Time     `bson:"deactivate_allow" json:"deactivate_allow"`
		IsEnabled       *bool         `bson:"is_enabled" json:"is_enabled"`
		Repeat          []string      `bson:"repeat" json:"repeat"`
		AllDay          *bool         `bson:"all_day" json:"all_day"`
		AllowedUsers    []Allowed     `bson:"allowed_users" json:"allowed_users"`
		AllowedSubjects []Allowed     `bson:"allowed_subjects" json:"allowed_subjects"`
	}
	Allowed struct {
		ID        string `bson:"id" json:"id"`
		Name      string `bson:"name" json:"name"`
		Activated *bool  `bson:"activated" json:"activated"`
	}
)

// GetAways : get all the aways for the user
func (a *Away) GetAways(collection string) ([]Away, error) {
	var aways []Away
	if err := a.DB.C(collection).
		Find(bson.M{"user_id": a.UserID}).
		All(&aways); err != nil {
		return nil, err
	}
	return aways, nil
}

// GetAways : get all the aways for the user
func (a *Away) GetAway(collection string) (interface{}, error) {
	if err := a.DB.C(collection).FindId(a.ID).One(&a); err != nil {
		return nil, err
	}
	return a, nil
}

// CreateAway : this will create a new if not already exists else will update.
func (a *Away) CreateAway(collection string) (interface{}, error) {
	if err := a.DB.C(collection).Find(
		bson.M{
			"$and": []bson.M{
				{"title": a.Title},
				{"user_id": a.UserID},
			}}).One(&a); err == mgo.ErrNotFound {
		a.ID = bson.NewObjectId()
		if err := a.DB.C(collection).Insert(a); err != nil {
			return nil, err
		}
	}
	return a, nil
}

// UpdateAway : this will update the away created via Create command.
func (a *Away) UpdateAway(collection string) (interface{}, error) {
	if _, err := a.DB.C(collection).FindId(a.ID).Apply(mgo.Change{
		Update: bson.M{"$set": bson.M{
			"title":            a.Title,
			"repeat":           a.Repeat,
			"all_day":          a.AllDay,
			"activate_allow":   a.ActivateAllow,
			"deactivate_allow": a.DeactivateAllow,
			"allowed_users":    a.AllowedUsers,
			"allowed_subjects": a.AllowedSubjects,
		}},
		ReturnNew: true,
	}, nil); err != nil {
		return nil, err
	}
	return a, nil
}

// DeleteAway : this will update the away created via Create command.
func (a *Away) DeleteAway(collection string) error {
	if err := a.DB.C(collection).RemoveId(a.ID); err != nil {
		return err
	}
	return nil
}
