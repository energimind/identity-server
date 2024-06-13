package service

import (
	"strings"

	"github.com/energimind/identity-server/internal/core/domain/admin"
)

func validateRealm(realm admin.Realm) (admin.Realm, error) {
	realm.Name = strings.TrimSpace(realm.Name)
	realm.Code = strings.TrimSpace(realm.Code)

	if err := checkName(realm.Name); err != nil {
		return realm, err
	}

	if err := checkCode(realm.Code); err != nil {
		return realm, err
	}

	return realm, nil
}

func validateDaemon(daemon admin.Daemon) (admin.Daemon, error) {
	daemon.Name = strings.TrimSpace(daemon.Name)
	daemon.Code = strings.TrimSpace(daemon.Code)

	if err := checkName("name"); err != nil {
		return admin.Daemon{}, err
	}

	if err := checkCode(daemon.Code); err != nil {
		return admin.Daemon{}, err
	}

	return daemon, nil
}

func validateProvider(provider admin.Provider) (admin.Provider, error) {
	provider.Name = strings.TrimSpace(provider.Name)
	provider.Code = strings.TrimSpace(provider.Code)

	if err := checkName(provider.Name); err != nil {
		return provider, err
	}

	if err := checkCode(provider.Code); err != nil {
		return provider, err
	}

	return provider, nil
}

func validateUser(user admin.User) (admin.User, error) {
	user.BindID = strings.TrimSpace(user.BindID)
	user.Username = strings.TrimSpace(user.Username)
	user.Email = strings.TrimSpace(user.Email)
	user.DisplayName = strings.TrimSpace(user.DisplayName)

	if err := checkEmpty("bindID", user.BindID); err != nil {
		return user, err
	}

	if err := checkEmpty("username", user.Username); err != nil {
		return user, err
	}

	if err := checkEmail(user.Email); err != nil {
		return user, err
	}

	return user, nil
}

func validateAPIKey(apiKey admin.APIKey) (admin.APIKey, error) {
	apiKey.Name = strings.TrimSpace(apiKey.Name)
	apiKey.Key = strings.TrimSpace(apiKey.Key)

	if err := checkName(apiKey.Name); err != nil {
		return apiKey, err
	}

	return apiKey, nil
}
