package admin

import "github.com/energimind/identity-service/core/domain/admin"

// FromApplication converts a domain application to a DTO application.
func FromApplication(app admin.Application) Application {
	return Application{
		ID:          string(app.ID),
		Code:        app.Code,
		Name:        app.Name,
		Description: app.Description,
		Enabled:     app.Enabled,
	}
}

// FromApplications converts a slice of domain applications to a slice of DTO applications.
func FromApplications(apps []admin.Application) []Application {
	dtos := make([]Application, len(apps))

	for i, app := range apps {
		dtos[i] = FromApplication(app)
	}

	return dtos
}

// ToApplication converts a DTO application to a domain application.
func ToApplication(app Application) admin.Application {
	return admin.Application{
		ID:          admin.ID(app.ID),
		Code:        app.Code,
		Name:        app.Name,
		Description: app.Description,
		Enabled:     app.Enabled,
	}
}

// FromProvider converts a domain provider to a DTO provider.
func FromProvider(provider admin.Provider) Provider {
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

// FromProviders converts a slice of domain providers to a slice of DTO providers.
func FromProviders(providers []admin.Provider) []Provider {
	dtos := make([]Provider, len(providers))

	for i, provider := range providers {
		dtos[i] = FromProvider(provider)
	}

	return dtos
}

// ToProvider converts a DTO provider to a domain provider.
func ToProvider(provider Provider) admin.Provider {
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

// FromUser converts a domain user to a DTO user.
func FromUser(user admin.User) User {
	return User{
		ID:          string(user.ID),
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Description: user.Description,
		Enabled:     user.Enabled,
		Role:        string(user.Role),
		APIKeys:     fromAPIKeys(user.APIKeys),
	}
}

// FromUsers converts a slice of domain users to a slice of DTO users.
func FromUsers(users []admin.User) []User {
	dtos := make([]User, len(users))

	for i, user := range users {
		dtos[i] = FromUser(user)
	}

	return dtos
}

// ToUser converts a DTO user to a domain user.
func ToUser(user User) admin.User {
	return admin.User{
		ID:          admin.ID(user.ID),
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Description: user.Description,
		Enabled:     user.Enabled,
		Role:        admin.SystemRole(user.Role),
		APIKeys:     toAPIKeys(user.APIKeys),
	}
}

// FromDaemon converts a domain daemon to a DTO daemon.
func FromDaemon(daemon admin.Daemon) Daemon {
	return Daemon{
		ID:          string(daemon.ID),
		Code:        daemon.Code,
		Name:        daemon.Name,
		Description: daemon.Description,
		Enabled:     daemon.Enabled,
		APIKeys:     fromAPIKeys(daemon.APIKeys),
	}
}

// FromDaemons converts a slice of domain daemons to a slice of DTO daemons.
func FromDaemons(daemons []admin.Daemon) []Daemon {
	dtos := make([]Daemon, len(daemons))

	for i, daemon := range daemons {
		dtos[i] = FromDaemon(daemon)
	}

	return dtos
}

// ToDaemon converts a DTO daemon to a domain daemon.
func ToDaemon(daemon Daemon) admin.Daemon {
	return admin.Daemon{
		ID:          admin.ID(daemon.ID),
		Code:        daemon.Code,
		Name:        daemon.Name,
		Description: daemon.Description,
		Enabled:     daemon.Enabled,
		APIKeys:     toAPIKeys(daemon.APIKeys),
	}
}

// fromAPIKey converts a domain API key to a DTO API key.
func fromAPIKey(apiKey admin.APIKey) APIKey {
	return APIKey{
		Name:        apiKey.Name,
		Description: apiKey.Description,
		Enabled:     apiKey.Enabled,
		Key:         apiKey.Key,
		ExpiresAt:   apiKey.ExpiresAt,
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
		Name:        apiKey.Name,
		Description: apiKey.Description,
		Enabled:     apiKey.Enabled,
		Key:         apiKey.Key,
		ExpiresAt:   apiKey.ExpiresAt,
	}
}

// toAPIKeys converts a slice of DTO API keys to a slice of domain API keys.
func toAPIKeys(apiKeys []APIKey) []admin.APIKey {
	dtos := make([]admin.APIKey, len(apiKeys))

	for i, apiKey := range apiKeys {
		dtos[i] = toAPIKey(apiKey)
	}

	return dtos
}
