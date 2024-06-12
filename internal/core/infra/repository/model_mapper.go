package repository

import "github.com/energimind/identity-server/internal/core/domain/admin"

func toID(id admin.ID) string {
	return id.String()
}

func fromID(id string) admin.ID {
	return admin.ID(id)
}

func toRealm(realm admin.Realm) dbRealm {
	return dbRealm{
		ID:          toID(realm.ID),
		Code:        realm.Code,
		Name:        realm.Name,
		Description: realm.Description,
		Enabled:     realm.Enabled,
	}
}

func fromRealm(realm dbRealm) admin.Realm {
	return admin.Realm{
		ID:          fromID(realm.ID),
		Code:        realm.Code,
		Name:        realm.Name,
		Description: realm.Description,
		Enabled:     realm.Enabled,
	}
}

func toProvider(provider admin.Provider) dbProvider {
	return dbProvider{
		ID:           toID(provider.ID),
		Type:         toProviderType(provider.Type),
		Code:         provider.Code,
		Name:         provider.Name,
		Description:  provider.Description,
		Enabled:      provider.Enabled,
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		RedirectURL:  provider.RedirectURL,
	}
}

func fromProvider(provider dbProvider) admin.Provider {
	return admin.Provider{
		ID:           fromID(provider.ID),
		Type:         fromProviderType(provider.Type),
		Code:         provider.Code,
		Name:         provider.Name,
		Description:  provider.Description,
		Enabled:      provider.Enabled,
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		RedirectURL:  provider.RedirectURL,
	}
}

func toUser(user admin.User) dbUser {
	return dbUser{
		ID:          toID(user.ID),
		RealmID:     toID(user.RealmID),
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Description: user.Description,
		Enabled:     user.Enabled,
		Role:        toSystemRole(user.Role),
		APIKeys:     mapSlice(user.APIKeys, toAPIKey),
	}
}

func fromUser(user dbUser) admin.User {
	return admin.User{
		ID:          fromID(user.ID),
		RealmID:     fromID(user.RealmID),
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Description: user.Description,
		Enabled:     user.Enabled,
		Role:        fromSystemRole(user.Role),
		APIKeys:     mapSlice(user.APIKeys, fromAPIKey),
	}
}

func toDaemon(daemon admin.Daemon) dbDaemon {
	return dbDaemon{
		ID:          toID(daemon.ID),
		RealmID:     toID(daemon.RealmID),
		Code:        daemon.Code,
		Name:        daemon.Name,
		Description: daemon.Description,
		Enabled:     daemon.Enabled,
		APIKeys:     mapSlice(daemon.APIKeys, toAPIKey),
	}
}

func fromDaemon(daemon dbDaemon) admin.Daemon {
	return admin.Daemon{
		ID:          fromID(daemon.ID),
		RealmID:     fromID(daemon.RealmID),
		Code:        daemon.Code,
		Name:        daemon.Name,
		Description: daemon.Description,
		Enabled:     daemon.Enabled,
		APIKeys:     mapSlice(daemon.APIKeys, fromAPIKey),
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
