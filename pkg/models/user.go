package models

import (
	"os"
	"time"

	"github.com/globalsign/mgo/bson"
	log "github.com/golang/glog"
	"github.com/mayadata-io/kubera-auth/pkg/types"
	"github.com/mayadata-io/kubera-auth/pkg/utils/uuid"
)

func init() {
	if os.Getenv("ADMIN_USERNAME") == "" || os.Getenv("ADMIN_PASSWORD") == "" {
		log.Fatal("Environment variables ADMIN_USERNAME or ADMIN_PASSWORD are not set")
	}
}

//UserCredentials contains the user information
type UserCredentials struct {
	ID              bson.ObjectId   `bson:"_id,omitempty" json:"_id"`
	UID             string          `bson:"uid,omitempty" json:"uid"`
	UserName        string          `bson:"username,omitempty" json:"username"`
	Password        string          `bson:"password,omitempty" json:"password"`
	Email           string          `bson:"email,omitempty" json:"email"`
	UnverifiedEmail string          `bson:"unverified_email,omitempty" json:"unverified_email"`
	Company         string          `bson:"company,omitempty" json:"company"`
	CompanyRole     string          `bson:"company_role,omitempty" json:"company_role"`
	Name            string          `bson:"name,omitempty" json:"name"`
	Kind            AuthType        `bson:"kind,omitempty" json:"kind"`
	Role            Role            `bson:"role,omitempty" json:"role"`
	LoggedIn        bool            `bson:"logged_in,omitempty" json:"logged_in"`
	SocialAuthID    string          `bson:"social_auth_id,omitempty" json:"social_auth_id"`
	CreatedAt       *time.Time      `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt       *time.Time      `bson:"updated_at,omitempty" json:"updated_at"`
	RemovedAt       *time.Time      `bson:"removed_at,omitempty" json:"removed_at"`
	State           State           `bson:"state,omitempty" json:"state"`
	OnBoardingState OnBoardingState `bson:"onboarding_state,omitempty" json:"onboarding_state"`
	Photo           string          `bson:"pictureUrl,omitempty" json:"pictureUrl"`
}

//AuthType determines the type of authentication opted by the user for login
type AuthType string

const (
	// LocalAuth is the local authentication needs username and a password
	LocalAuth AuthType = "local"

	// GithubAuth authenticates via github OAuth
	GithubAuth AuthType = "github"

	// GoogleAuth authenticates via Google OAuth
	GoogleAuth AuthType = "google"
)

// Role states the role of the user in the portal
type Role string

const (
	//RoleAdmin gives the admin permissions to a user
	RoleAdmin Role = "admin"

	//RoleUser gives the normal user permissions to a user
	RoleUser Role = "user"
)

// OnBoardingState helps UI to define the state at which the user is currently while being onBoarded
type OnBoardingState int

// golint: unused
const (
	BoardingStateInvalid               OnBoardingState = iota // Invalid State used as zero value
	BoardingStateSignup                                       // Signup started (EmailNotVerified)
	BoardingStateEmailVerified                                // EmailVerified
	BoardingStateUnverifiedAndComplete                        // UnverifiedAndComplete
	BoardingStateVerifiedAndComplete                          // VerifiedAndComplete
)

//DefaultUser is the admin user created by default
var DefaultUser = &UserCredentials{
	UID:             uuid.Must(uuid.NewRandom()).String(),
	Name:            types.DefaultName,
	UserName:        types.DefaultUserName,
	Password:        types.DefaultUserPassword,
	Role:            RoleAdmin,
	Kind:            LocalAuth,
	OnBoardingState: BoardingStateUnverifiedAndComplete,
}

//PublicUserInfo displays the information of the user that is publicly available
type PublicUserInfo struct {
	ID              bson.ObjectId   `json:"_id"`
	UID             string          `json:"uid"`
	UserName        string          `json:"username"`
	Email           string          `json:"email"`
	UnverifiedEmail string          `json:"unverified_email"`
	Company         string          `json:"company"`
	CompanyRole     string          `json:"company_role"`
	Name            string          `json:"name"`
	Kind            AuthType        `json:"kind"`
	Role            Role            `json:"role"`
	LoggedIn        bool            `json:"logged_in"`
	CreatedAt       *time.Time      `json:"created_at"`
	UpdatedAt       *time.Time      `json:"updated_at"`
	RemovedAt       *time.Time      `json:"removed_at"`
	SocialAuthID    string          `json:"social_auth_id,omitempty"`
	State           State           `json:"state"`
	OnBoardingState OnBoardingState `json:"onboarding_state"`
}

//State is the current state of the database entry of the user
type State string

const (
	//StateCreated means admin has created the user but the user has still not logged in
	StateCreated State = "created"
	//StateActive means user has logged in successfully
	StateActive State = "active"
	//StateRemoved means user has been deleted
	// golint: unused
	StateRemoved State = "removed"
)

// GetPublicInfo fetches the pubicUserInfo from User
func (u *UserCredentials) GetPublicInfo() *PublicUserInfo {
	return &PublicUserInfo{
		Name:            u.Name,
		UserName:        u.UserName,
		Email:           u.Email,
		UnverifiedEmail: u.UnverifiedEmail,
		Company:         u.Company,
		CompanyRole:     u.CompanyRole,
		ID:              u.ID,
		UID:             u.UID,
		Kind:            u.Kind,
		Role:            u.Role,
		LoggedIn:        u.LoggedIn,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
		RemovedAt:       u.RemovedAt,
		State:           u.State,
		OnBoardingState: u.OnBoardingState,
	}
}
