package factory

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"platform-go/internal/infrastructure/dbmodel"
)

// UserFactory is a factory for generating User dbmodel instances for tests and seeders.
//
// Usage:
//
//	// Create a User with default values
//	user := factory.NewUserFactory().Build()
//
//	// Create a User with specific values
//	user := factory.NewUserFactory().
//	    WithName("Test Admin").
//	    WithAccountTypeAdmin().
//	    WithEmail("test@example.com").
//	    WithActive().
//	    Build()
//
//	// Insert into DB
//	user, err := factory.NewUserFactory().Insert(db)
//	// or
//	user := factory.NewUserFactory().MustInsert(db)
type UserFactory struct {
	name         string
	email        *string
	passwordHash *string
	accountType  string
	active       bool
}

// NewUserFactory creates a new UserFactory with default values
func NewUserFactory() *UserFactory {
	// Use bcrypt.MinCost to ensure fast execution in tests and seeders
	hashed, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.MinCost)
	hashedStr := string(hashed)
	email := "admin@platform.com"

	return &UserFactory{
		name:         "Administrator",
		email:        &email,
		passwordHash: &hashedStr,
		accountType:  "admin",
		active:       true, // Active by default
	}
}

// WithUser sets values from an existing dbmodel.User
func (f *UserFactory) WithUser(user dbmodel.User) *UserFactory {
	f.name = user.Name
	f.email = user.Email
	f.passwordHash = user.PasswordHash
	f.accountType = user.AccountType
	f.active = user.Active
	return f
}

// WithName sets the user's name
func (f *UserFactory) WithName(name string) *UserFactory {
	f.name = name
	return f
}

// WithEmail sets the user's email
func (f *UserFactory) WithEmail(email string) *UserFactory {
	f.email = &email
	return f
}

// WithPassword sets and hashes the password using bcrypt
func (f *UserFactory) WithPassword(password string) *UserFactory {
	// Use bcrypt.MinCost to ensure fast execution in tests and seeders
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	hashedStr := string(hashed)
	f.passwordHash = &hashedStr
	return f
}

// WithPasswordHash sets the hashed password directly
func (f *UserFactory) WithPasswordHash(hash string) *UserFactory {
	f.passwordHash = &hash
	return f
}

// WithAccountType sets the account type
func (f *UserFactory) WithAccountType(accountType string) *UserFactory {
	f.accountType = accountType
	return f
}

// WithAccountTypeAdmin sets the account type to admin
func (f *UserFactory) WithAccountTypeAdmin() *UserFactory {
	f.accountType = "admin"
	return f
}

// WithAccountTypeUser sets the account type to user
func (f *UserFactory) WithAccountTypeUser() *UserFactory {
	f.accountType = "user"
	return f
}

// WithActive sets the active status to true
func (f *UserFactory) WithActive() *UserFactory {
	f.active = true
	return f
}

// WithInactive sets the active status to false
func (f *UserFactory) WithInactive() *UserFactory {
	f.active = false
	return f
}

// WithIsActive sets the active status
func (f *UserFactory) WithIsActive(active bool) *UserFactory {
	f.active = active
	return f
}

// Build creates a new dbmodel.User instance
func (f *UserFactory) Build() *dbmodel.User {
	now := time.Now()
	return &dbmodel.User{
		Name:         f.name,
		Email:        f.email,
		PasswordHash: f.passwordHash,
		AccountType:  f.accountType,
		Active:       f.active,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// BuildWithID creates a new dbmodel.User instance with a specific ID
func (f *UserFactory) BuildWithID(id uint) *dbmodel.User {
	now := time.Now()
	return &dbmodel.User{
		ID:           id,
		Name:         f.name,
		Email:        f.email,
		PasswordHash: f.passwordHash,
		AccountType:  f.accountType,
		Active:       f.active,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// Insert persists the User to the database
func (f *UserFactory) Insert(db *gorm.DB) (*dbmodel.User, error) {
	user := f.Build()
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// MustInsert persists the User to the database, panicking on error
func (f *UserFactory) MustInsert(db *gorm.DB) *dbmodel.User {
	user, err := f.Insert(db)
	if err != nil {
		panic(err)
	}
	return user
}
