package usermanager

import (
	"github.com/globalsign/mgo/bson"
	"golang.org/x/crypto/bcrypt"

	"github.com/mayadata-io/kubera-auth/pkg/models"
	"github.com/mayadata-io/kubera-auth/pkg/store"
	"github.com/mayadata-io/kubera-auth/pkg/types"
)

// UpdateUserDetails updates the user information
func UpdateUserDetails(userStore *store.UserStore, user *models.UserCredentials) (*models.PublicUserInfo, error) {
	// There will be following possible transitions of OnboardingState
	// BoardingStateSignup -> BoardingStateEmailVerified -> BoardingStateVerifiedAndComplete
	// BoardingStateSignup -> BoardingStateUnverifiedAndComplete -> BoardingStateVerifiedAndComplete
	switch user.OnBoardingState {
	case models.BoardingStateSignup:
		{
			if user.Email != "" {
				user.OnBoardingState = models.BoardingStateEmailVerified
			} else if user.Company != "" {
				user.OnBoardingState = models.BoardingStateUnverifiedAndComplete
			}
		}
	case models.BoardingStateEmailVerified:
		{
			if user.Company != "" {
				user.OnBoardingState = models.BoardingStateVerifiedAndComplete
			}
		}
	case models.BoardingStateUnverifiedAndComplete:
		{
			if user.Email != "" {
				user.OnBoardingState = models.BoardingStateVerifiedAndComplete
			}
		}
	}

	err := userStore.UpdateUser(user)
	return user.GetPublicInfo(), err
}

// UpdatePassword sets the new user password
func UpdatePassword(userStore *store.UserStore, newPassword, userID string) (*models.PublicUserInfo, error) {
	var storedUser *models.UserCredentials
	var err error

	storedUser, err = userStore.GetUser(bson.M{"username": userID, "kind": models.LocalAuth})
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), types.PasswordEncryptionCost)
	if err != nil {
		return nil, err
	}
	storedUser.Password = string(hashedPassword)

	err = userStore.UpdateUser(storedUser)
	return storedUser.GetPublicInfo(), err
}
