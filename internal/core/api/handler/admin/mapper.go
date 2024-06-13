package admin

import (
	"time"

	"github.com/energimind/identity-server/internal/core/domain/admin"
)

// fromRealm converts a domain realm to a DTO realm.
func fromRealm(realm admin.Realm) Realm {
	return Realm{
		ID:          string(realm.ID),
		Code:        realm.Code,
		Name:        realm.Name,
		Description: realm.Description,
		Enabled:     realm.Enabled,
	}
}

// fromRealms converts a slice of domain realms to a slice of DTO realms.
func fromRealms(realms []admin.Realm) []Realm {
	dtos := make([]Realm, len(realms))

	for i, realm := range realms {
		dtos[i] = fromRealm(realm)
	}

	return dtos
}

// toRealm converts a DTO realm to a domain realm.
func toRealm(realm Realm) admin.Realm {
	return admin.Realm{
		ID:          admin.ID(realm.ID),
		Code:        realm.Code,
		Name:        realm.Name,
		Description: realm.Description,
		Enabled:     realm.Enabled,
	}
}

// fromProvider converts a domain provider to a DTO provider.
func fromProvider(provider admin.Provider) Provider {
	return Provider{
		ID:           string(provider.ID),
		Type:         string(provider.Type),
		Code:         provider.Code,
		Name:         provider.Name,
		Description:  provider.Description,
		Enabled:      provider.Enabled,
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		RedirectURL:  provider.RedirectURL,
	}
}

// fromProviders converts a slice of domain providers to a slice of DTO providers.
func fromProviders(providers []admin.Provider) []Provider {
	dtos := make([]Provider, len(providers))

	for i, provider := range providers {
		dtos[i] = fromProvider(provider)
	}

	return dtos
}

// toProvider converts a DTO provider to a domain provider.
func toProvider(provider Provider) admin.Provider {
	return admin.Provider{
		ID:           admin.ID(provider.ID),
		Type:         admin.ProviderType(provider.Type),
		Code:         provider.Code,
		Name:         provider.Name,
		Description:  provider.Description,
		Enabled:      provider.Enabled,
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		RedirectURL:  provider.RedirectURL,
	}
}

// fromUser converts a domain sessionUser to a DTO sessionUser.
func fromUser(user admin.User) User {
	return User{
		ID:          string(user.ID),
		BindID:      user.BindID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Description: user.Description,
		Enabled:     user.Enabled,
		Role:        string(user.Role),
	}
}

// fromUsers converts a slice of domain users to a slice of DTO users.
func fromUsers(users []admin.User) []User {
	dtos := make([]User, len(users))

	for i, user := range users {
		dtos[i] = fromUser(user)
	}

	return dtos
}

// toUser converts a DTO sessionUser to a domain sessionUser.
func toUser(user User) admin.User {
	return admin.User{
		ID:          admin.ID(user.ID),
		BindID:      user.BindID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Description: user.Description,
		Enabled:     user.Enabled,
		Role:        admin.SystemRole(user.Role),
	}
}

// fromDaemon converts a domain daemon to a DTO daemon.
func fromDaemon(daemon admin.Daemon) Daemon {
	return Daemon{
		ID:          string(daemon.ID),
		Code:        daemon.Code,
		Name:        daemon.Name,
		Description: daemon.Description,
		Enabled:     daemon.Enabled,
	}
}

// fromDaemons converts a slice of domain daemons to a slice of DTO daemons.
func fromDaemons(daemons []admin.Daemon) []Daemon {
	dtos := make([]Daemon, len(daemons))

	for i, daemon := range daemons {
		dtos[i] = fromDaemon(daemon)
	}

	return dtos
}

// toDaemon converts a DTO daemon to a domain daemon.
func toDaemon(daemon Daemon) admin.Daemon {
	return admin.Daemon{
		ID:          admin.ID(daemon.ID),
		Code:        daemon.Code,
		Name:        daemon.Name,
		Description: daemon.Description,
		Enabled:     daemon.Enabled,
	}
}

// fromAPIKey converts a domain API key to a DTO API key.
func fromAPIKey(apiKey admin.APIKey) APIKey {
	return APIKey{
		ID:          string(apiKey.ID),
		Name:        apiKey.Name,
		Description: apiKey.Description,
		Enabled:     apiKey.Enabled,
		Key:         apiKey.Key,
		ExpiresAt:   fromDate(apiKey.ExpiresAt),
	}
}

// fromAPIKeys converts a slice of domain API keys to a slice of DTO API keys.
func fromAPIKeys(apiKeys []admin.APIKey) []APIKey {
	dtos := make([]APIKey, len(apiKeys))

	for i, apiKey := range apiKeys {
		dtos[i] = fromAPIKey(apiKey)
	}

	return dtos
}

// toAPIKey converts a DTO API key to a domain API key.
func toAPIKey(apiKey APIKey) admin.APIKey {
	return admin.APIKey{
		ID:          admin.ID(apiKey.ID),
		Name:        apiKey.Name,
		Description: apiKey.Description,
		Enabled:     apiKey.Enabled,
		Key:         apiKey.Key,
		ExpiresAt:   toDate(apiKey.ExpiresAt),
	}
}

func fromDate(t time.Time) *string {
	if t.IsZero() {
		return nil
	}

	value := t.Format(time.DateOnly)

	return &value
}

func toDate(s *string) time.Time {
	if s == nil || *s == "" {
		return time.Time{}
	}

	t, err := time.Parse(time.DateOnly, *s)
	if err != nil {
		return time.Time{}
	}

	return t
}
