package aways

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	CreateAway(req *Away) (entity *Away, err error)
	GetAwayList(userID primitive.ObjectID, enable bool) (entity []*Away, err error)
	GetAway(awayID primitive.ObjectID) (entity *Away, err error)
	UpdateAway(req *Away) (err error)
	UpdateAwayMode(isEnabled bool, userID primitive.ObjectID) (total int64, err error)
	DeleteAway(awayID primitive.ObjectID) (err error)
}
