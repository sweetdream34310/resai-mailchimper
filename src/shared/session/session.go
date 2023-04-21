package session

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cloudsrc/api.awaymail.v1.go/libs"
	"github.com/cloudsrc/api.awaymail.v1.go/src/domain/redis"
	googleWrapper "github.com/cloudsrc/api.awaymail.v1.go/src/infrastructure/google"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/constants"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/models"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils"
	ctxSess "github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/gmail/v1"
)

const signingSecret = "thisisaverylongbutsecuretokenstring"

type claims struct {
	UserID       primitive.ObjectID `json:"user_id"`
	Email        string             `json:"email"`
	RefreshToken string             `json:"refresh_token"`
	AuthToken    string             `json:"auth_token"`
	Provider     string             `json:"provider"`
	Name         string             `json:"name"`
	Photo        string             `json:"photo"`
	SwiftToken   string             `json:"swift_token"`
	DevTokenKey  bool               `json:"dev_token_key"`
	jwt.StandardClaims
}

// Session : This creates a new session instances for the routes
type Session struct {
	App           *libs.App
	GoogleWrapper googleWrapper.Wrapper
	RedisRepo     redis.Repository
}

// CheckToken : This checks the auth token for the user.
func (session *Session) CheckToken(ctx *gin.Context) {
	data, _ := ctx.Get(ctxSess.AppSession)
	ctxSess := data.(*ctxSess.Context)
	authHeader := ctx.Request.Header.Get("Authorization")
	agent := ctx.Request.Header.Get("X-agent")
	authHeaderParts := strings.Split(authHeader, " ")
	res := libs.Response{
		Status: http.StatusUnauthorized,
	}
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		res.Message = "authorization header format must be bearer {token}"
		res.SendResponse(ctxSess, ctx)
		ctx.Abort()
	} else {
		tokenString := authHeaderParts[1]
		claimsToken := &claims{}
		token, _ := jwt.ParseWithClaims(tokenString, claimsToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(signingSecret), nil
		})
		//key := fmt.Sprintf(constants.UserSession, claimsToken.UserID.Hex(), claimsToken.SwiftToken)
		//userSession := session.RedisRepo.GetKey(key)
		//if token.Valid && userSession == claimsToken.UserID.Hex() {
		if token.Valid {
			if !checkTokenExpiry(claimsToken.StandardClaims.ExpiresAt) {
				profile := googleWrapper.UserProfile{}
				decRefreshToken, _ := utils.Decrypt(session.App.Config.Salt, claimsToken.RefreshToken)
				switch claimsToken.Provider {
				case constants.GmailClient:
					profile, _ = session.GoogleWrapper.GetProfile(ctxSess, agent, decRefreshToken)

					watchKey := fmt.Sprintf(constants.UserWatch, profile.Email)
					watchSess := session.RedisRepo.GetKey(watchKey)
					if watchSess == "" {
						err := session.GoogleWrapper.StopWatchPushNotification(ctxSess, ctxSess.UserSession.RefreshToken)
						if err != nil {
							ctxSess.ErrorMessage = err.Error()
							ctxSess.Lv4()
						}
						watchRes, _ := session.GoogleWrapper.WatchPushNotification(ctxSess, decRefreshToken, &gmail.WatchRequest{
							LabelIds:  []string{"INBOX"},
							TopicName: fmt.Sprintf("projects/%s/topics/%s", session.App.Config.Gpubsub.ProjectName, session.App.Config.Gpubsub.Topic),
						})
						if watchRes != nil {
							session.RedisRepo.SetKey(watchKey, watchRes.HistoryId, "EX", 86400)
							session.RedisRepo.SetKey(fmt.Sprintf(constants.UserHistory, profile.Email), watchRes.HistoryId)
						}
					}
				default:
					res.Message = "client not yet configured or present"
					res.Status = http.StatusUnauthorized
					res.SendResponse(ctxSess, ctx)
					ctx.Abort()
					return
				}
				user := models.UserSession{
					UserID:       claimsToken.UserID,
					Email:        claimsToken.Email,
					SwiftToken:   claimsToken.SwiftToken,
					RefreshToken: decRefreshToken,
					AuthToken:    profile.AuthToken,
					Name:         profile.Name,
					Photo:        profile.Photo,
					DevTokenKey:  claimsToken.DevTokenKey,
				}
				ctxSess.UserSession = user
				ctx.Next()
			} else {
				res.Message = "token expired, please re-authenticate"
				res.SendResponse(ctxSess, ctx)
				ctx.Abort()
			}
		} else {
			res.Message = "token invalid"
			res.SendResponse(ctxSess, ctx)
			ctx.Abort()
		}
	}
}

func checkTokenExpiry(timestamp interface{}) bool {
	if validity, ok := timestamp.(int64); ok {
		tm := time.Unix(int64(validity), 0)
		remainder := tm.Sub(time.Now())
		if remainder > 0 {
			return false
		}
	}
	return true
}

// NewBearerToken : this method is used to create a new bearer token for the user
func NewBearerToken(userID primitive.ObjectID, email, provider, refreshToken, accessToken, name, photo, swiftToken string, devTokenKey bool) (string, time.Time) {
	expiry := time.Now().Add(time.Hour * 24 * 365).Unix()
	claims := &claims{
		UserID:       userID,
		Email:        email,
		Name:         name,
		Photo:        photo,
		Provider:     provider,
		AuthToken:    accessToken,
		RefreshToken: refreshToken,
		SwiftToken:   swiftToken,
		DevTokenKey:  devTokenKey,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiry,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if tokenString, err := token.SignedString([]byte(signingSecret)); err == nil {
		return tokenString, time.Unix(expiry, 0)
	}
	return "", time.Unix(expiry, 0)
}
