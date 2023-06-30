package session

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fredele20/microservice-practice/ms.users/db/mongod"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrTokenInvalid          = errors.New("invalid token string provided")
	ErrTokenExpired          = errors.New("sorry, session has expired. Please login again to continue")
	ErrTokenSessionNotFound  = errors.New("session not found or destroyed")
	ErrInvalidUnitOfValidity = errors.New("invalid unit of validity, you must provide HOUR or MINUTE")
)

func generateToken(userId, role, email, firstName, lastName string) string {
	payload := &TokenPayload{
		Role:      role,
		UserId:    userId,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Payload: jwt.Payload{
			Issuer:   "Golang",
			Subject:  "Golang JWT",
			Audience: jwt.Audience{""},
			IssuedAt: jwt.NumericDate(time.Now()),
			JWTID:    "Golang JWT Auth",
		},
	}

	token, err := jwt.Sign(payload, jwt.NewHS256([]byte(os.Getenv("JWT_SECRET"))))
	if err != nil {
		logrus.Debugf("error generating JWT Token: %s", err)
		return ""
	}

	return string(token)
}

func GetSessionByToken(token string) (*Session, error) {
	var ctx context.Context
	if strings.TrimSpace(token) == "" {
		return nil, ErrTokenInvalid
	}

	// verify token
	_, err := verifyAuthToken(token)
	if err != nil {
		logrus.WithError(err).Error("failed to confirm session validity")
		fmt.Println(err)
		return nil, err
	}

	var session Session

	if err := mongod.SessionCollection().FindOne(ctx, bson.M{"token": token}).Decode(&session); err != nil {
		// logrus.WithError(err.Err()).Error("failed")
		return nil, ErrTokenSessionNotFound
	}

	if err = session.AssertValidity(); err != nil {
		logrus.WithError(err).Error("failed to get assert session validity")
		_ = DestroySession(session.Token) // Destroy it.
		return nil, err
	}

	fmt.Println(session)
	return &session, nil
}

func verifyAuthToken(token string) (*TokenPayload, error) {
	secret := jwt.NewHS256([]byte(os.Getenv("JWT_SECRET")))
	var payloadBody TokenPayload
	_, err := jwt.Verify([]byte(token), secret, &payloadBody)
	if err != nil {
		return nil, ErrTokenInvalid
	}
	return &payloadBody, nil
}

func newSession(userId, role, email, firstName, lastName string, validity time.Duration, unitOfValidity UnitOfValidity) *Session {
	token := generateToken(userId, role, email, firstName, lastName)
	return &Session{
		Token:          token,
		Role:           role,
		UserId:      userId,
		Validity:       validity,
		LastUsage:      time.Now(),
		UnitOfValidity: unitOfValidity,
		TimeCreated:    time.Now(),
	}
}

func CreateSession(payload Session) (string, error) {
	var ctx context.Context

	if !payload.UnitOfValidity.IsValid() {
		return "", ErrInvalidUnitOfValidity
	}

	s := newSession(payload.UserId, payload.Role, payload.Email, payload.FirstName, payload.LastName, payload.Validity, payload.UnitOfValidity)
	_, err := mongod.SessionCollection().InsertOne(ctx, s)
	if err != nil {
		logrus.WithError(err).Error("failed to store session on db")
		fmt.Println("Error: ", err)
		return "", err
	}

	return s.Token, nil
}

func DestroySession(token string) error {
	// Delete session from the DB
	var ctx context.Context
	session := mongod.SessionCollection().FindOneAndDelete(ctx, bson.M{"token": token})
	if session.Err() != nil {
		logrus.WithError(session.Err()).Error("session with the token not found")
		return session.Err()
	}

	return nil
}
