package loginmanager

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"golang.org/x/crypto/bcrypt"

	"github.com/mayadata-io/kubera-auth/manager/jwtmanager"
	"github.com/mayadata-io/kubera-auth/manager/usermanager"
	"github.com/mayadata-io/kubera-auth/pkg/errors"
	"github.com/mayadata-io/kubera-auth/pkg/generates"
	"github.com/mayadata-io/kubera-auth/pkg/models"
	"github.com/mayadata-io/kubera-auth/pkg/store"
)

// LocalLoginUser verifies user password
func LocalLoginUser(userStore *store.UserStore, accessGenerate *generates.JWTAccessGenerate, username, password string) (*models.Token, error) {
	tgr, err := validationAuthenticateRequest(userStore, username, password)
	if err != nil {
		return nil, err
	}

	ti, err := jwtmanager.GenerateAuthToken(accessGenerate, tgr, models.TokenLogin)
	if err != nil {
		return nil, err
	}

	storedUser, err := userStore.GetUser(bson.M{"username": username, "kind": models.LocalAuth})
	if err != nil {
		return nil, err
	}
	storedUser.LoggedIn = true
	return ti, userStore.UpdateUser(storedUser)
}

// SocialLoginUser get the user information
func SocialLoginUser(userStore *store.UserStore, accessGenerate *generates.JWTAccessGenerate, user *models.UserCredentials) (*models.Token, error) {
	query := bson.M{"social_auth_id": user.SocialAuthID}
	storedUser, err := userStore.GetUser(query)
	if err != nil && err == mgo.ErrNotFound {
		// If user does not exists
		createErr := usermanager.CreateSocialUser(userStore, user)
		if createErr != nil {
			return nil, createErr
		}
		storedUser = user
	} else if err != nil {
		// Error other than user exists
		return nil, err
	} else {
		// If user exists
		storedUser.LoggedIn = true
		err = userStore.UpdateUser(storedUser)
		if err != nil {
			return nil, err
		}
	}

	tgr := &jwtmanager.TokenGenerateRequest{
		UserInfo: storedUser.GetPublicInfo(),
	}
	return jwtmanager.GenerateAuthToken(accessGenerate, tgr, models.TokenLogin)
}

// validationAuthenticateRequest the authenticate request validation
func validationAuthenticateRequest(userStore *store.UserStore, username, password string) (*jwtmanager.TokenGenerateRequest, error) {
	user, err := userStore.GetUser(bson.M{"username": username, "kind": models.LocalAuth})
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.ErrInvalidPassword
	}

	req := &jwtmanager.TokenGenerateRequest{
		UserInfo: user.GetPublicInfo(),
	}
	return req, nil
}

// LogoutUser verifies user password
func LogoutUser(userStore *store.UserStore, id bson.ObjectId) error {
	storedUser, err := userStore.GetUserByID(id)
	if err != nil {
		return err
	}
	storedUser.LoggedIn = false
	return userStore.UpdateUser(storedUser)
}
