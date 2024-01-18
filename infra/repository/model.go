package repository

import "time"

// dbApplication is the database model for an application.
type dbApplication struct {
	ID          string `bson:"id"`
	Code        string `bson:"code"`
	Name        string `bson:"name,omitempty"`
	Description string `bson:"description,omitempty"`
	Enabled     bool   `bson:"enabled"`
}

// dbProvider is the database model for an authentication provider.
type dbProvider struct {
	ID            string         `bson:"id"`
	ApplicationID string         `bson:"applicationId"`
	Type          dbProviderType `bson:"type"`
	Code          string         `bson:"code"`
	Name          string         `bson:"name,omitempty"`
	Description   string         `bson:"description,omitempty"`
	Enabled       bool           `bson:"enabled"`
	ClientID      string         `bson:"clientId"`
	ClientSecret  string         `bson:"clientSecret"`
	RedirectURL   string         `bson:"redirectUrl"`
}

// dbUser is the database model for a user.
type dbUser struct {
	ID            string       `bson:"id"`
	ApplicationID string       `bson:"applicationId"`
	Username      string       `bson:"username"`
	Description   string       `bson:"description,omitempty"`
	Enabled       bool         `bson:"enabled"`
	Role          dbSystemRole `bson:"role"`
}

// dbAccount is the database model for an account.
type dbAccount struct {
	ID      string `bson:"id"`
	UserID  string `bson:"userId"`
	UserIDN string `bson:"userIdn"`
}

// dbDaemon is the database model for a daemon.
type dbDaemon struct {
	ID            string `bson:"id"`
	ApplicationID string `bson:"applicationId"`
	Code          string `bson:"code"`
	Name          string `bson:"name,omitempty"`
	Description   string `bson:"description,omitempty"`
	Enabled       bool   `bson:"enabled"`
}

// dbAPIKey is the database model for an API key.
type dbAPIKey struct {
	ID          string         `bson:"id"`
	OwnerID     string         `bson:"ownerId"`
	OwnerType   dbKeyOwnerType `bson:"ownerType"`
	Name        string         `bson:"name,omitempty"`
	Description string         `bson:"description,omitempty"`
	Enabled     bool           `bson:"enabled"`
	Key         string         `bson:"key"`
	ExpiresAt   time.Time      `bson:"expiresAt"`
}
