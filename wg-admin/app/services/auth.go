package services

import (
	"context"
	"errors"
	"mp2720/wg-admin/wg-admin/storage"
	"mp2720/wg-admin/wg-admin/storage/data"
	"mp2720/wg-admin/wg-admin/transaction"
	"mp2720/wg-admin/wg-admin/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrUserAuthIsNotAllowed = errors.New("Authentication is not allowed for this user")
	ErrInvalidToken         = errors.New("Invalid auth token")
)

type AuthService interface {
	IssueAuthTokenForUser(ctx context.Context, userUUID uuid.UUID) (string, error)

	AuthenticateUser(ctx context.Context, token string) (data.User, error)

	// AuthenticateRouter(ctx context.Context, token string) (data.Router, error)
}

type authService struct {
	userRepo   storage.UserRepo
	tm         *transaction.Manager
	clock      utils.Clock
	signingKey string
	issuer     string
}

type jwtClaims struct {
	jwt.RegisteredClaims
	// "user" or "router"
	SubjectType string `json:"sub_type"`
}

func NewAuthService(
	userRepo storage.UserRepo,
	tm *transaction.Manager,
	clock utils.Clock,
	signingKey string,
	issuer string,
) AuthService {
	return authService{userRepo, tm, clock, signingKey, issuer}
}

func (as authService) signToken(sub string, subjectType string, issuedAt time.Time) (string, error) {
	issuedAtJwt := jwt.NumericDate{Time: issuedAt}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:  sub,
			IssuedAt: &issuedAtJwt,
			Issuer:   as.issuer,
		},
		SubjectType: subjectType,
	})
	return token.SignedString([]byte(as.signingKey))
}

func (as authService) parseToken(strToken string) (jwtClaims, error) {
	var claims jwtClaims
	token, err := jwt.ParseWithClaims(strToken, &claims, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != "HS256" {
			return nil, errors.New("Invalid signing method")
		}
		return []byte(as.signingKey), nil
	})
	if err != nil {
		return jwtClaims{}, err
	}
	if !token.Valid {
		return jwtClaims{}, errors.New("Token is invalid")
	}

	return claims, nil
}

func (as authService) IssueAuthTokenForUser(
	ctx context.Context,
	userUUID uuid.UUID,
) (string, error) {
	var err error
	var user data.User

	now := as.clock.Now()

	err = as.tm.InTransaction(ctx, func(ctx context.Context) error {
		user, err = as.userRepo.GetByUUIDLocked(ctx, userUUID)
		if err != nil {
			return err
		}

		if !user.CanHaveAuthenticationToken() {
			return ErrUserAuthIsNotAllowed
		}

		err = user.Update(data.UserPatch{TokenIssuedAt: &now})
		if err != nil {
			return err
		}

		err = as.userRepo.Save(ctx, user)
		return err
	})
	if err != nil {
		return "", err
	}

	return as.signToken(user.UUID.String(), "user", now)
}

func (as authService) AuthenticateUser(ctx context.Context, token string) (data.User, error) {
	claims, err := as.parseToken(token)
	if err != nil {
		return data.User{}, ErrInvalidToken
	}
	if claims.SubjectType != "user" {
		return data.User{}, ErrInvalidToken
	}

	uuid, err := uuid.Parse(claims.Subject)
	if err != nil {
		return data.User{}, ErrInvalidToken
	}

	user, err := as.userRepo.GetByUUID(ctx, uuid)
	if err != nil {
		return data.User{}, err
	}

	if !user.CanBeAuthenticatedWithToken(claims.IssuedAt.Time) {
		return data.User{}, ErrUserAuthIsNotAllowed
	}

	return user, nil
}
