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
