package usermanager

import (
	"github.com/globalsign/mgo/bson"
	"golang.org/x/crypto/bcrypt"

	"github.com/mayadata-io/kubera-auth/pkg/errors"
	"github.com/mayadata-io/kubera-auth/pkg/models"
	"github.com/mayadata-io/kubera-auth/pkg/store"
	"github.com/mayadata-io/kubera-auth/pkg/types"
)

// UpdateUserDetails get the user information
func UpdateUserDetails(userStore *store.UserStore, user *models.UserCredentials) (*models.PublicUserInfo, error) {

	storedUser, err := userStore.GetUserByID(user.GetID())
	if err != nil {
		return nil, errors.ErrInvalidUser
	}

	if user.Email != nil {
		storedUser.Email = user.Email
	}
	if user.Name != nil {
		storedUser.Name = user.Name
	}
	if user.UserName != nil {
		storedUser.UserName = user.UserName
	}
	if user.IsEmailVerified != nil {
		storedUser.IsEmailVerified = user.IsEmailVerified
	}

	exists, err := CheckUserExists(userStore, storedUser)
	if err != nil {
		return nil, err
	} else if exists == false {
		return nil, errors.ErrInvalidUser
	}

	err = userStore.UpdateUser(storedUser)
	return storedUser.GetPublicInfo(), err
}

// UpdatePassword get the user information
func UpdatePassword(userStore *store.UserStore, reset bool, oldPassword, newPassword, userID string) (*models.PublicUserInfo, error) {

	var storedUser *models.UserCredentials
	var err error
	if reset == true {
		storedUser, err = userStore.GetUser(bson.M{"username": userID, "kind": models.LocalAuth})
		if err != nil {
			return nil, err
		}
	} else {
		storedUser, err = userStore.GetUser(bson.M{"uid": userID, "kind": models.LocalAuth})
		if err != nil {
			return nil, err
		}

		err = bcrypt.CompareHashAndPassword([]byte(storedUser.GetPassword()), []byte(oldPassword))
		if err != nil {
			return storedUser.GetPublicInfo(), errors.ErrInvalidPassword
		}
		// Updating the state when the user changes the password himself. Will be useful for the first time
		storedUser.State = models.StateActive
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), types.PasswordEncryptionCost)
	if err != nil {
		return nil, err
	}

	password := string(hashedPassword)
	storedUser.Password = &password

	err = userStore.UpdateUser(storedUser)
	return storedUser.GetPublicInfo(), err
}
