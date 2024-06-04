package repository

import "github.com/energimind/identity-server/internal/core/domain/admin"

func toID(id admin.ID) string {
	return id.String()
}

func fromID(id string) admin.ID {
	return admin.ID(id)
}

func toApplication(app admin.Application) dbApplication {
	return dbApplication{
		ID:          toID(app.ID),
		Code:        app.Code,
		Name:        app.Name,
		Description: app.Description,
		Enabled:     app.Enabled,
	}
}

func fromApplication(app dbApplication) admin.Application {
	return admin.Application{
		ID:          fromID(app.ID),
		Code:        app.Code,
		Name:        app.Name,
		Description: app.Description,
		Enabled:     app.Enabled,
	}
}

func toProvider(provider admin.Provider) dbProvider {
	return dbProvider{
		ID:            toID(provider.ID),
		ApplicationID: toID(provider.ApplicationID),
		Type:          toProviderType(provider.Type),
		Code:          provider.Code,
		Name:          provider.Name,
		Description:   provider.Description,
		Enabled:       provider.Enabled,
		ClientID:      provider.ClientID,
		ClientSecret:  provider.ClientSecret,
		RedirectURL:   provider.RedirectURL,
	}
}

func fromProvider(provider dbProvider) admin.Provider {
	return admin.Provider{
		ID:            fromID(provider.ID),
		ApplicationID: fromID(provider.ApplicationID),
		Type:          fromProviderType(provider.Type),
		Code:          provider.Code,
		Name:          provider.Name,
		Description:   provider.Description,
		Enabled:       provider.Enabled,
		ClientID:      provider.ClientID,
		ClientSecret:  provider.ClientSecret,
		RedirectURL:   provider.RedirectURL,
	}
}

func toUser(user admin.User) dbUser {
	return dbUser{
		ID:            toID(user.ID),
		ApplicationID: toID(user.ApplicationID),
		Username:      user.Username,
		Email:         user.Email,
		DisplayName:   user.DisplayName,
		Description:   user.Description,
		Enabled:       user.Enabled,
		Role:          toSystemRole(user.Role),
		APIKeys:       mapSlice(user.APIKeys, toAPIKey),
	}
}

func fromUser(user dbUser) admin.User {
	return admin.User{
		ID:            fromID(user.ID),
		ApplicationID: fromID(user.ApplicationID),
		Username:      user.Username,
		Email:         user.Email,
		DisplayName:   user.DisplayName,
		Description:   user.Description,
		Enabled:       user.Enabled,
		Role:          fromSystemRole(user.Role),
		APIKeys:       mapSlice(user.APIKeys, fromAPIKey),
	}
}

func toDaemon(daemon admin.Daemon) dbDaemon {
	return dbDaemon{
		ID:            toID(daemon.ID),
		ApplicationID: toID(daemon.ApplicationID),
		Code:          daemon.Code,
		Name:          daemon.Name,
		Description:   daemon.Description,
		Enabled:       daemon.Enabled,
		APIKeys:       mapSlice(daemon.APIKeys, toAPIKey),
	}
}

func fromDaemon(daemon dbDaemon) admin.Daemon {
	return admin.Daemon{
		ID:            fromID(daemon.ID),
		ApplicationID: fromID(daemon.ApplicationID),
		Code:          daemon.Code,
		Name:          daemon.Name,
		Description:   daemon.Description,
		Enabled:       daemon.Enabled,
		APIKeys:       mapSlice(daemon.APIKeys, fromAPIKey),
	}
}

func toAPIKey(apiKey admin.APIKey) dbAPIKey {
	return dbAPIKey{
		ID:          toID(apiKey.ID),
		Name:        apiKey.Name,
		Description: apiKey.Description,
		Enabled:     apiKey.Enabled,
		Key:         apiKey.Key,
		ExpiresAt:   apiKey.ExpiresAt,
	}
}

func fromAPIKey(apiKey dbAPIKey) admin.APIKey {
	return admin.APIKey{
		ID:          fromID(apiKey.ID),
		Name:        apiKey.Name,
		Description: apiKey.Description,
		Enabled:     apiKey.Enabled,
		Key:         apiKey.Key,
		ExpiresAt:   apiKey.ExpiresAt,
	}
}
