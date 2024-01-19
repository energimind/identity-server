package repository

import "github.com/energimind/identity-service/domain/auth"

func toID(id auth.ID) string {
	return id.String()
}

func fromID(id string) auth.ID {
	return auth.ID(id)
}

func toApplication(app auth.Application) dbApplication {
	return dbApplication{
		ID:          toID(app.ID),
		Code:        app.Code,
		Name:        app.Name,
		Description: app.Description,
		Enabled:     app.Enabled,
	}
}

func fromApplication(app dbApplication) auth.Application {
	return auth.Application{
		ID:          fromID(app.ID),
		Code:        app.Code,
		Name:        app.Name,
		Description: app.Description,
		Enabled:     app.Enabled,
	}
}

func fromApplications(apps []dbApplication) []auth.Application {
	converted := make([]auth.Application, len(apps))

	for i, app := range apps {
		converted[i] = fromApplication(app)
	}

	return converted
}

func toProvider(provider auth.Provider) dbProvider {
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

func fromProvider(provider dbProvider) auth.Provider {
	return auth.Provider{
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

func toUser(user auth.User) dbUser {
	return dbUser{
		ID:            toID(user.ID),
		ApplicationID: toID(user.ApplicationID),
		Username:      user.Username,
		Description:   user.Description,
		Enabled:       user.Enabled,
		Role:          toSystemRole(user.Role),
	}
}

func fromUser(user dbUser) auth.User {
	return auth.User{
		ID:            fromID(user.ID),
		ApplicationID: fromID(user.ApplicationID),
		Username:      user.Username,
		Description:   user.Description,
		Enabled:       user.Enabled,
		Role:          fromSystemRole(user.Role),
	}
}

func toAccount(account auth.Account) dbAccount {
	return dbAccount{
		ID:      toID(account.ID),
		UserID:  toID(account.UserID),
		UserIDN: account.UserIDN,
	}
}

func fromAccount(account dbAccount) auth.Account {
	return auth.Account{
		ID:      fromID(account.ID),
		UserID:  fromID(account.UserID),
		UserIDN: account.UserIDN,
	}
}

func toDaemon(daemon auth.Daemon) dbDaemon {
	return dbDaemon{
		ID:            toID(daemon.ID),
		ApplicationID: toID(daemon.ApplicationID),
		Code:          daemon.Code,
		Name:          daemon.Name,
		Description:   daemon.Description,
		Enabled:       daemon.Enabled,
	}
}

func fromDaemon(daemon dbDaemon) auth.Daemon {
	return auth.Daemon{
		ID:            fromID(daemon.ID),
		ApplicationID: fromID(daemon.ApplicationID),
		Code:          daemon.Code,
		Name:          daemon.Name,
		Description:   daemon.Description,
		Enabled:       daemon.Enabled,
	}
}

func toAPIKey(apiKey auth.APIKey) dbAPIKey {
	return dbAPIKey{
		ID:          toID(apiKey.ID),
		OwnerID:     toID(apiKey.OwnerID),
		OwnerType:   toKeyOwnerType(apiKey.OwnerType),
		Name:        apiKey.Name,
		Description: apiKey.Description,
		Enabled:     apiKey.Enabled,
		Key:         apiKey.Key,
		ExpiresAt:   apiKey.ExpiresAt,
	}
}

func fromAPIKey(apiKey dbAPIKey) auth.APIKey {
	return auth.APIKey{
		ID:          fromID(apiKey.ID),
		OwnerID:     fromID(apiKey.OwnerID),
		OwnerType:   fromKeyOwnerType(apiKey.OwnerType),
		Name:        apiKey.Name,
		Description: apiKey.Description,
		Enabled:     apiKey.Enabled,
		Key:         apiKey.Key,
		ExpiresAt:   apiKey.ExpiresAt,
	}
}
