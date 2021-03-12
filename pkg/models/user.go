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
	ID              bson.ObjectId   `bson:"_id,omitempty"`
	UID             *string         `bson:"uid,omitempty"`
	UserName        *string         `bson:"username,omitempty"`
	Password        *string         `bson:"password,omitempty"`
	Email           *string         `bson:"email,omitempty"`
	UnverifiedEmail *string         `bson:"unverified_email,omitempty"`
	Company         *string         `bson:"company,omitempty"`
	CompanyRole     *string         `bson:"company_role,omitempty"`
	Name            *string         `bson:"name,omitempty"`
	Kind            AuthType        `bson:"kind,omitempty"`
	Role            Role            `bson:"role,omitempty"`
	LoggedIn        *bool           `bson:"logged_in,omitempty"`
	SocialAuthID    *int64          `bson:"social_auth_id,omitempty"`
	CreatedAt       *time.Time      `bson:"created_at,omitempty"`
	UpdatedAt       *time.Time      `bson:"updated_at,omitempty"`
	RemovedAt       *time.Time      `bson:"removed_at,omitempty"`
	State           State           `bson:"state,omitempty"`
	OnBoardingState OnBoardingState `bson:"onboarding_state,omitempty"`
}

//AuthType determines the type of authentication opted by the user for login
type AuthType string

const (
	//LocalAuth is the local authentication needs username and a password
	LocalAuth AuthType = "local"

	//GithubAuth authenticates via github
	GithubAuth AuthType = "github"

	//GmailAuth authenticates via gmail
	GmailAuth AuthType = "gmail"
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

const (
	BoardingStateInvalid               OnBoardingState = iota // Invalid State used as zero value
	BoardingStateSignup                                       // Signup started (EmailNotVerified)
	BoardingStateEmailVerified                                // EmailVerified
	BoardingStateUnverifiedAndComplete                        // UnverifiedAndComplete
	BoardingStateVerifiedAndComplete                          // VerifiedAndComplete
)

var adminUID = uuid.Must(uuid.NewRandom()).String()

//DefaultUser is the admin user created by default
var DefaultUser = &UserCredentials{
	UID:             &adminUID,
	Name:            &types.DefaultName,
	UserName:        &types.DefaultUserName,
	Password:        &types.DefaultUserPassword,
	Role:            RoleAdmin,
	Kind:            LocalAuth,
	OnBoardingState: BoardingStateUnverifiedAndComplete,
}

//PublicUserInfo displays the information of the user that is publicly available
type PublicUserInfo struct {
	ID              bson.ObjectId   `json:"_id"`
	UID             *string         `json:"uid"`
	UserName        *string         `json:"username"`
	Email           *string         `json:"email"`
	UnverifiedEmail *string         `json:"unverified_email,omitempty"`
	Company         *string         `json:"company,omitempty"`
	CompanyRole     *string         `json:"company_role,omitempty"`
	Name            *string         `json:"name"`
	Kind            AuthType        `json:"kind"`
	Role            Role            `json:"role"`
	LoggedIn        *bool           `json:"logged_in"`
	SocialAuthID    *int64          `json:"social_auth_id,omitempty"`
	CreatedAt       *time.Time      `json:"created_at"`
	UpdatedAt       *time.Time      `json:"updated_at"`
	RemovedAt       *time.Time      `json:"removed_at"`
	State           State           `json:"state"`
	OnBoardingState OnBoardingState `json:"onboarding_state,omitempty"`
}

//State is the current state of the database entry of the user
type State string

const (
	//StateCreated means admin has created the user but the user has still not logged in
	StateCreated State = "created"
	//StateActive means user has logged in successfully
	StateActive State = "active"
	//StateRemoved means user has been deleted
	StateRemoved State = "removed"
)

// GetID user id
func (u *UserCredentials) GetID() bson.ObjectId {
	return u.ID
}

// GetUID user password
func (u *UserCredentials) GetUID() string {
	if u == nil || u.UID == nil {
		return ""
	}
	return *u.UID
}

// GetUserName user username
func (u *UserCredentials) GetUserName() string {
	if u == nil || u.UserName == nil {
		return ""
	}
	return *u.UserName
}

// GetPassword user password
func (u *UserCredentials) GetPassword() string {
	if u == nil || u.Password == nil {
		return ""
	}
	return *u.Password
}

// GetEmail user email
func (u *UserCredentials) GetEmail() string {
	if u == nil || u.Email == nil {
		return ""
	}
	return *u.Email
}

// GetUnverifiedEmail gets the email that needs to be verified
func (u *UserCredentials) GetUnverifiedEmail() string {
	if u == nil || u.UnverifiedEmail == nil {
		return ""
	}
	return *u.UnverifiedEmail
}

// GetCompany gets user company name
func (u *UserCredentials) GetCompany() string {
	if u == nil || u.Company == nil {
		return ""
	}
	return *u.Company
}

// GetCompanyRole get user role in the company
func (u *UserCredentials) GetCompanyRole() string {
	if u == nil || u.CompanyRole == nil {
		return ""
	}
	return *u.CompanyRole
}

// GetName returns user name
func (u *UserCredentials) GetName() string {
	if u == nil || u.Name == nil {
		return ""
	}
	return *u.Name
}

// GetKind user password
func (u *UserCredentials) GetKind() AuthType {
	if u == nil {
		return ""
	}
	return u.Kind
}

// GetRole user password
func (u *UserCredentials) GetRole() Role {
	if u == nil {
		return ""
	}
	return u.Role
}

// GetLoggedIn user password
func (u *UserCredentials) GetLoggedIn() bool {
	if u == nil || u.LoggedIn == nil {
		return false
	}
	return *u.LoggedIn
}

// GetSocialAuthID returns all the social authentications of the user
func (u *UserCredentials) GetSocialAuthID() int64 {
	if u == nil || u.SocialAuthID == nil {
		return 0
	}
	return *u.SocialAuthID
}

// GetCreatedAt defines the time at which this user was created
func (u *UserCredentials) GetCreatedAt() time.Time {
	if u == nil || u.CreatedAt == nil {
		return time.Time{}
	}
	return *u.CreatedAt
}

// GetUpdatedAt defines the time at which user was last updated
func (u *UserCredentials) GetUpdatedAt() time.Time {
	if u == nil || u.UpdatedAt == nil {
		return time.Time{}
	}
	return *u.UpdatedAt
}

// GetRemovedAt defines the time at which this user was removed
func (u *UserCredentials) GetRemovedAt() time.Time {
	if u == nil || u.RemovedAt == nil {
		return time.Time{}
	}
	return *u.RemovedAt
}

// GetState gets the user state
func (u *UserCredentials) GetState() State {
	if u == nil {
		return ""
	}
	return u.State
}

// GetOnBoardingState gets the state the user is presnt for onboarding process
func (u *UserCredentials) GetOnBoardingState() OnBoardingState {
	if u == nil {
		return BoardingStateInvalid
	}
	return u.OnBoardingState
}

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

// GetID user id
func (u *PublicUserInfo) GetID() bson.ObjectId {
	return u.ID
}

// GetUID user password
func (u *PublicUserInfo) GetUID() string {
	if u == nil || u.UID == nil {
		return ""
	}
	return *u.UID
}

// GetUserName user username
func (u *PublicUserInfo) GetUserName() string {
	if u == nil || u.UserName == nil {
		return ""
	}
	return *u.UserName
}

// GetEmail user email
func (u *PublicUserInfo) GetEmail() string {
	if u == nil || u.Email == nil {
		return ""
	}
	return *u.Email
}

// GetUnverifiedEmail gets the email that needs to be verified
func (u *PublicUserInfo) GetUnverifiedEmail() string {
	if u == nil || u.UnverifiedEmail == nil {
		return ""
	}
	return *u.UnverifiedEmail
}

// GetCompany gets user company name
func (u *PublicUserInfo) GetCompany() string {
	if u == nil || u.Company == nil {
		return ""
	}
	return *u.Company
}

// GetCompanyRole get user role in the company
func (u *PublicUserInfo) GetCompanyRole() string {
	if u == nil || u.CompanyRole == nil {
		return ""
	}
	return *u.CompanyRole
}

// GetName returns user name
func (u *PublicUserInfo) GetName() string {
	if u == nil || u.Name == nil {
		return ""
	}
	return *u.Name
}

// GetKind user password
func (u *PublicUserInfo) GetKind() AuthType {
	if u == nil {
		return ""
	}
	return u.Kind
}

// GetRole user password
func (u *PublicUserInfo) GetRole() Role {
	if u == nil {
		return ""
	}
	return u.Role
}

// GetLoggedIn user password
func (u *PublicUserInfo) GetLoggedIn() bool {
	if u == nil || u.LoggedIn == nil {
		return false
	}
	return *u.LoggedIn
}

// GetSocialAuthID returns all the social authentications of the user
func (u *PublicUserInfo) GetSocialAuthID() int64 {
	if u == nil || u.SocialAuthID == nil {
		return 0
	}
	return *u.SocialAuthID
}

// GetCreatedAt defines the time at which this user was created
func (u *PublicUserInfo) GetCreatedAt() time.Time {
	if u == nil || u.CreatedAt == nil {
		return time.Time{}
	}
	return *u.CreatedAt
}

// GetUpdatedAt defines the time at which user was last updated
func (u *PublicUserInfo) GetUpdatedAt() time.Time {
	if u == nil || u.UpdatedAt == nil {
		return time.Time{}
	}
	return *u.UpdatedAt
}

// GetRemovedAt defines the time at which this user was removed
func (u *PublicUserInfo) GetRemovedAt() time.Time {
	if u == nil || u.RemovedAt == nil {
		return time.Time{}
	}
	return *u.RemovedAt
}

// GetState user password
func (u *PublicUserInfo) GetState() State {
	if u == nil {
		return ""
	}
	return u.State
}

// GetOnBoardingState gets the state the user is presnt for onboarding process
func (u *PublicUserInfo) GetOnBoardingState() OnBoardingState {
	if u == nil {
		return BoardingStateInvalid
	}
	return u.OnBoardingState
}

// GetUserCredentials converts the struct into UserCredentials
func (u *PublicUserInfo) GetUserCredentials() *UserCredentials {
	return &UserCredentials{
		ID:              u.ID,
		UID:             u.UID,
		UserName:        u.UserName,
		Email:           u.Email,
		UnverifiedEmail: u.UnverifiedEmail,
		Company:         u.Company,
		CompanyRole:     u.CompanyRole,
		Name:            u.Name,
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
