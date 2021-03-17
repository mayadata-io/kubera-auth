package usermanager

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/mayadata-io/kubera-auth/pkg/errors"
	"github.com/mayadata-io/kubera-auth/pkg/models"
	"github.com/mayadata-io/kubera-auth/pkg/store"
	"github.com/mayadata-io/kubera-auth/pkg/types"
	"github.com/mayadata-io/kubera-auth/pkg/utils/uuid"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser get the user information
func CreateUser(userStore *store.UserStore, user *models.UserCredentials) (*models.PublicUserInfo, error) {
	exists, err := IsUserExists(userStore, user)
	if err != nil {
		return nil, err
	} else if exists {
		return nil, errors.ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), types.PasswordEncryptionCost)
	if err != nil {
		return nil, err
	}

	var newUser *models.UserCredentials
	if user.Role == models.RoleAdmin {
		newUser = user
		newUser.Password = string(hashedPassword)
	} else {
		newUser = &models.UserCredentials{
			UID:             uuid.Must(uuid.NewRandom()).String(),
			UserName:        user.UserName,
			Password:        string(hashedPassword),
			Name:            user.Name,
			UnverifiedEmail: user.UserName,
			Kind:            models.LocalAuth,
			Role:            models.RoleUser,
			OnBoardingState: models.BoardingStateSignup,
		}
	}

	err = userStore.Set(newUser)
	return newUser.GetPublicInfo(), err
}

//CreateSocialUser creates a user if the user opts logging in with some oauth
func CreateSocialUser(userStore *store.UserStore, user *models.UserCredentials) error {
	query := bson.M{"email": user.Email, "kind": models.LocalAuth}
	storedUser, err := userStore.GetUser(query)
	if err != nil && err == mgo.ErrNotFound {
		user.UserName = generateUserName(user.Name)
		user.UID = uuid.Must(uuid.NewRandom()).String()
	} else if err != nil {
		return err
	} else {
		user.UserName = storedUser.UserName
		user.UID = storedUser.UID
	}
	return userStore.Set(user)
}
