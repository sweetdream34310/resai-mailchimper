package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/cloudsrc/api.awaymail.v1.go/config"
	"github.com/cloudsrc/api.awaymail.v1.go/libs"
	"github.com/cloudsrc/api.awaymail.v1.go/provider"
	"github.com/cloudsrc/api.awaymail.v1.go/provider/google"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const signingSecret = "thisisaverylongbutsecuretokenstring"

type claims struct {
	ID           bson.ObjectId `json:"id"`
	RefreshToken string        `json:"refresh_token"`
	AuthToken    string        `json:"auth_token"`
	Provider     string        `json:"provider"`
	jwt.StandardClaims
}

// Session : This creates a new session instances for the routes
type Session struct {
	App *libs.App
}

// CheckToken : This checks the auth token for the user.
func (session *Session) CheckToken(ctx *gin.Context) {
	var prov provider.Provider
	authHeader := ctx.Request.Header.Get("Authorization")
	agent := ctx.Request.Header.Get("X-agent")
	authHeaderParts := strings.Split(authHeader, " ")
	res := libs.Response{
		Status: http.StatusUnauthorized,
	}
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		res.Message = "authorization header format must be bearer {token}"
		res.SendResponse(nil, ctx)
		ctx.Abort()
	} else {
		tokenString := authHeaderParts[1]
		claims := &claims{}
		token, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(signingSecret), nil
		})
		if token.Valid {
			if !checkTokenExpiry(claims.StandardClaims.ExpiresAt) {
				switch claims.Provider {
				case "gmail":
					prov = gmailProvider(
						session.App.Config,
						*session.App.Redis,
						session.App.DB.MongoDB,
						session.App.Rabbit,
					)
					prov.SetAgent(agent)
					if err := prov.ValidateToken(claims.RefreshToken); err != nil {
						res.Message = err.Error()
						res.SendResponse(nil, ctx)
						ctx.Abort()
						return
					}
				default:
					res.Message = "client not yet configured or present"
					res.Status = http.StatusUnauthorized
					res.SendResponse(nil, ctx)
					ctx.Abort()
					return
				}
				prov.SetName(claims.Provider)
				ctx.Set("provider", prov)
				ctx.Next()
			} else {
				res.Message = "token expired, please re-authenticate"
				res.SendResponse(nil, ctx)
				ctx.Abort()
			}
		} else {
			res.Message = "token invalid"
			res.SendResponse(nil, ctx)
			ctx.Abort()
		}
	}
}

func gmailProvider(config config.Config, rclient libs.RedisClient, mClient *mgo.Database, rabbitClient *libs.RabbitClient) *google.Provider {
	return google.New(&config, rclient, mClient, rabbitClient)
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
func NewBearerToken(provider, refreshToken string) (string, time.Time) {
	expiry := time.Now().Add(time.Hour * 24 * 365).Unix()
	claims := &claims{
		Provider:     provider,
		RefreshToken: refreshToken,
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
