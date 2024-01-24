package repository

import (
	"testing"
	"time"

	"github.com/energimind/identity-service/core/domain/auth"
	"github.com/energimind/identity-service/test/utils"
	"github.com/stretchr/testify/require"
)

func Test_allUserFieldsAreMapped(t *testing.T) {
	utils.CheckAllFieldsAreMapped(t, auth.Application{}, dbApplication{})
	utils.CheckAllFieldsAreMapped(t, auth.Provider{}, dbProvider{})
	utils.CheckAllFieldsAreMapped(t, auth.User{}, dbUser{})
	utils.CheckAllFieldsAreMapped(t, auth.Account{}, dbAccount{})
	utils.CheckAllFieldsAreMapped(t, auth.Daemon{}, dbDaemon{})
	utils.CheckAllFieldsAreMapped(t, auth.APIKey{}, dbAPIKey{})

	utils.CheckAllFieldsAreMapped(t, dbApplication{}, auth.Application{})
	utils.CheckAllFieldsAreMapped(t, dbProvider{}, auth.Provider{})
	utils.CheckAllFieldsAreMapped(t, dbUser{}, auth.User{})
	utils.CheckAllFieldsAreMapped(t, dbAccount{}, auth.Account{})
	utils.CheckAllFieldsAreMapped(t, dbDaemon{}, auth.Daemon{})
	utils.CheckAllFieldsAreMapped(t, dbAPIKey{}, auth.APIKey{})
}

func Test_mapApplication(t *testing.T) {
	t.Parallel()

	from := auth.Application{
		ID:          "app1",
		Code:        "app1",
		Name:        "Application 1",
		Description: "Application 1",
		Enabled:     true,
	}

	expected := dbApplication{
		ID:          "app1",
		Code:        "app1",
		Name:        "Application 1",
		Description: "Application 1",
		Enabled:     true,
	}

	mapped := toApplication(from)
	back := fromApplication(mapped)

	require.Equal(t, expected, mapped)
	require.Equal(t, from, back)
}

func Test_mapProvider(t *testing.T) {
	t.Parallel()

	from := auth.Provider{
		ID:            "provider1",
		ApplicationID: "app1",
		Type:          auth.ProviderTypeGoogle,
		Code:          "google",
		Name:          "Google",
		Description:   "Google Description",
		Enabled:       true,
		ClientID:      "client1",
		ClientSecret:  "secret1",
		RedirectURL:   "https://google.com",
	}

	expected := dbProvider{
		ID:            "provider1",
		ApplicationID: "app1",
		Type:          dbProviderTypeGoogle,
		Code:          "google",
		Name:          "Google",
		Description:   "Google Description",
		Enabled:       true,
		ClientID:      "client1",
		ClientSecret:  "secret1",
		RedirectURL:   "https://google.com",
	}

	mapped := toProvider(from)
	back := fromProvider(mapped)

	require.Equal(t, expected, mapped)
	require.Equal(t, from, back)
}

func Test_mapUser(t *testing.T) {
	t.Parallel()

	from := auth.User{
		ID:            "user1",
		ApplicationID: "app1",
		Username:      "user1",
		Description:   "User 1",
		Enabled:       true,
		Role:          auth.SystemRoleManager,
		Accounts:      []auth.Account{{}},
		APIKeys:       []auth.APIKey{{}},
	}

	expected := dbUser{
		ID:            "user1",
		ApplicationID: "app1",
		Username:      "user1",
		Description:   "User 1",
		Enabled:       true,
		Role:          dbSystemRoleManager,
		Accounts:      []dbAccount{{}},
		APIKeys:       []dbAPIKey{{}},
	}

	mapped := toUser(from)
	back := fromUser(mapped)

	require.Equal(t, expected, mapped)
	require.Equal(t, from, back)
}

func Test_mapAccount(t *testing.T) {
	t.Parallel()

	from := auth.Account{
		Identifier: "user1@domain.com",
		Enabled:    true,
	}

	expected := dbAccount{
		Identifier: "user1@domain.com",
		Enabled:    true,
	}

	mapped := toAccount(from)
	back := fromAccount(mapped)

	require.Equal(t, expected, mapped)
	require.Equal(t, from, back)
}

func Test_mapDaemon(t *testing.T) {
	t.Parallel()

	from := auth.Daemon{
		ID:            "daemon1",
		ApplicationID: "app1",
		Code:          "daemon1",
		Name:          "Daemon 1",
		Description:   "Daemon 1",
		Enabled:       true,
		APIKeys:       []auth.APIKey{{}},
	}

	expected := dbDaemon{
		ID:            "daemon1",
		ApplicationID: "app1",
		Code:          "daemon1",
		Name:          "Daemon 1",
		Description:   "Daemon 1",
		Enabled:       true,
		APIKeys:       []dbAPIKey{{}},
	}

	mapped := toDaemon(from)
	back := fromDaemon(mapped)

	require.Equal(t, expected, mapped)
	require.Equal(t, from, back)
}

func Test_mapAPIKey(t *testing.T) {
	t.Parallel()

	now := time.Now().Round(time.Second)

	from := auth.APIKey{
		Name:        "Key 1",
		Description: "Key 1",
		Enabled:     true,
		Key:         "key1",
		ExpiresAt:   now,
	}

	expected := dbAPIKey{
		Name:        "Key 1",
		Description: "Key 1",
		Enabled:     true,
		Key:         "key1",
		ExpiresAt:   now,
	}

	mapped := toAPIKey(from)
	back := fromAPIKey(mapped)

	require.Equal(t, expected, mapped)
	require.Equal(t, from, back)
}
