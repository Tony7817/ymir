package google

import (
	"context"
	"time"

	"cloud.google.com/go/auth/credentials/idtoken"
	"github.com/pkg/errors"
	"ymir.com/pkg/xerr"
)

type User struct {
	GoogleUserId string
	Email        string
	AvatarUtl    string
	UserName     string
}

func ValidateToken(ctx context.Context, token string) (*User, error) {
	res, err := idtoken.Validate(ctx, token, "")
	if err != nil {
		return nil, err
	}

	if !res.Claims["email_verified"].(bool) {
		return nil, xerr.NewErrCode(xerr.ErrorInvalidEmail)
	}

	if time.Now().UTC().Unix() > res.Expires {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.NotAuthorizedError), "google token has expired")
	}

	return &User{
		GoogleUserId: res.Subject,
		Email:        res.Claims["email"].(string),
		AvatarUtl:    res.Claims["picture"].(string),
		UserName:     res.Claims["name"].(string),
	}, nil
}
