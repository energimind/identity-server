package repository

import (
	"testing"
	"time"

	"github.com/energimind/go-kit/testutil/mapping"
	"github.com/energimind/identity-server/core/domain/admin"
	"github.com/stretchr/testify/require"
)

func Test_allUserFieldsAreMapped(t *testing.T) {
	mapping.CheckAllFieldsAreMapped(t, admin.Application{}, dbApplication{})
	mapping.CheckAllFieldsAreMapped(t, admin.Provider{}, dbProvider{})
	mapping.CheckAllFieldsAreMapped(t, admin.User{}, dbUser{})
	mapping.CheckAllFieldsAreMapped(t, admin.Daemon{}, dbDaemon{})
	mapping.CheckAllFieldsAreMapped(t, admin.APIKey{}, dbAPIKey{})

	mapping.CheckAllFieldsAreMapped(t, dbApplication{}, admin.Application{})
	mapping.CheckAllFieldsAreMapped(t, dbProvider{}, admin.Provider{})
	mapping.CheckAllFieldsAreMapped(t, dbUser{}, admin.User{})
	mapping.CheckAllFieldsAreMapped(t, dbDaemon{}, admin.Daemon{})
	mapping.CheckAllFieldsAreMapped(t, dbAPIKey{}, admin.APIKey{})
}

func Test_mapApplication(t *testing.T) {
	t.Parallel()

	from := admin.Application{
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

	from := admin.Provider{
		ID:            "provider1",
		ApplicationID: "app1",
		Type:          admin.ProviderTypeGoogle,
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

	from := admin.User{
		ID:            "user1",
		ApplicationID: "app1",
		Username:      "user1",
		Email:         "user@somedomain.com",
		Description:   "User 1",
		Enabled:       true,
		Role:          admin.SystemRoleManager,
		APIKeys:       []admin.APIKey{{}},
	}

	expected := dbUser{
		ID:            "user1",
		ApplicationID: "app1",
		Username:      "user1",
		Email:         "user@somedomain.com",
		Description:   "User 1",
		Enabled:       true,
		Role:          dbSystemRoleManager,
		APIKeys:       []dbAPIKey{{}},
	}

	mapped := toUser(from)
	back := fromUser(mapped)

	require.Equal(t, expected, mapped)
	require.Equal(t, from, back)
}

func Test_mapDaemon(t *testing.T) {
	t.Parallel()

	from := admin.Daemon{
		ID:            "daemon1",
		ApplicationID: "app1",
		Code:          "daemon1",
		Name:          "Daemon 1",
		Description:   "Daemon 1",
		Enabled:       true,
		APIKeys:       []admin.APIKey{{}},
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

	from := admin.APIKey{
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
