package dto

import "github.com/energimind/identity-service/core/domain/auth"

// FromApplication converts a domain application to a DTO application.
func FromApplication(app auth.Application) Application {
	return Application{
		ID:          string(app.ID),
		Code:        app.Code,
		Name:        app.Name,
		Description: app.Description,
		Enabled:     app.Enabled,
	}
}

// FromApplications converts a slice of domain applications to a slice of DTO applications.
func FromApplications(apps []auth.Application) []Application {
	dtos := make([]Application, len(apps))

	for i, app := range apps {
		dtos[i] = FromApplication(app)
	}

	return dtos
}

// ToApplication converts a DTO application to a domain application.
func ToApplication(app Application) auth.Application {
	return auth.Application{
		ID:          auth.ID(app.ID),
		Code:        app.Code,
		Name:        app.Name,
		Description: app.Description,
		Enabled:     app.Enabled,
	}
}

// FromProvider converts a domain provider to a DTO provider.
func FromProvider(provider auth.Provider) Provider {
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
func FromProviders(providers []auth.Provider) []Provider {
	dtos := make([]Provider, len(providers))

	for i, provider := range providers {
		dtos[i] = FromProvider(provider)
	}

	return dtos
}

// ToProvider converts a DTO provider to a domain provider.
func ToProvider(provider Provider) auth.Provider {
	return auth.Provider{
		ID:           auth.ID(provider.ID),
		Type:         auth.ProviderType(provider.Type),
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
func FromUser(user auth.User) User {
	return User{
		ID:          string(user.ID),
		Username:    user.Username,
		Description: user.Description,
		Enabled:     user.Enabled,
		Role:        string(user.Role),
		Accounts:    fromAccounts(user.Accounts),
		APIKeys:     fromAPIKeys(user.APIKeys),
	}
}

// FromUsers converts a slice of domain users to a slice of DTO users.
func FromUsers(users []auth.User) []User {
	dtos := make([]User, len(users))

	for i, user := range users {
		dtos[i] = FromUser(user)
	}

	return dtos
}

// ToUser converts a DTO user to a domain user.
func ToUser(user User) auth.User {
	return auth.User{
		ID:          auth.ID(user.ID),
		Username:    user.Username,
		Description: user.Description,
		Enabled:     user.Enabled,
		Role:        auth.SystemRole(user.Role),
		Accounts:    toAccounts(user.Accounts),
		APIKeys:     toAPIKeys(user.APIKeys),
	}
}

// fromAccount converts a domain account to a DTO account.
func fromAccount(account auth.Account) Account {
	return Account{
		Identifier: account.Identifier,
		Enabled:    account.Enabled,
	}
}

// fromAccounts converts a slice of domain accounts to a slice of DTO accounts.
func fromAccounts(accounts []auth.Account) []Account {
	dtos := make([]Account, len(accounts))

	for i, account := range accounts {
		dtos[i] = fromAccount(account)
	}

	return dtos
}

// toAccount converts a DTO account to a domain account.
func toAccount(account Account) auth.Account {
	return auth.Account{
		Identifier: account.Identifier,
		Enabled:    account.Enabled,
	}
}

// toAccounts converts a slice of DTO accounts to a slice of domain accounts.
func toAccounts(accounts []Account) []auth.Account {
	dtos := make([]auth.Account, len(accounts))

	for i, account := range accounts {
		dtos[i] = toAccount(account)
	}

	return dtos
}

// FromDaemon converts a domain daemon to a DTO daemon.
func FromDaemon(daemon auth.Daemon) Daemon {
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
func FromDaemons(daemons []auth.Daemon) []Daemon {
	dtos := make([]Daemon, len(daemons))

	for i, daemon := range daemons {
		dtos[i] = FromDaemon(daemon)
	}

	return dtos
}

// ToDaemon converts a DTO daemon to a domain daemon.
func ToDaemon(daemon Daemon) auth.Daemon {
	return auth.Daemon{
		ID:          auth.ID(daemon.ID),
		Code:        daemon.Code,
		Name:        daemon.Name,
		Description: daemon.Description,
		Enabled:     daemon.Enabled,
		APIKeys:     toAPIKeys(daemon.APIKeys),
	}
}

// fromAPIKey converts a domain API key to a DTO API key.
func fromAPIKey(apiKey auth.APIKey) APIKey {
	return APIKey{
		Name:        apiKey.Name,
		Description: apiKey.Description,
		Enabled:     apiKey.Enabled,
		Key:         apiKey.Key,
		ExpiresAt:   apiKey.ExpiresAt,
	}
}

// fromAPIKeys converts a slice of domain API keys to a slice of DTO API keys.
func fromAPIKeys(apiKeys []auth.APIKey) []APIKey {
	dtos := make([]APIKey, len(apiKeys))

	for i, apiKey := range apiKeys {
		dtos[i] = fromAPIKey(apiKey)
	}

	return dtos
}

// toAPIKey converts a DTO API key to a domain API key.
func toAPIKey(apiKey APIKey) auth.APIKey {
	return auth.APIKey{
		Name:        apiKey.Name,
		Description: apiKey.Description,
		Enabled:     apiKey.Enabled,
		Key:         apiKey.Key,
		ExpiresAt:   apiKey.ExpiresAt,
	}
}

// toAPIKeys converts a slice of DTO API keys to a slice of domain API keys.
func toAPIKeys(apiKeys []APIKey) []auth.APIKey {
	dtos := make([]auth.APIKey, len(apiKeys))

	for i, apiKey := range apiKeys {
		dtos[i] = toAPIKey(apiKey)
	}

	return dtos
}
