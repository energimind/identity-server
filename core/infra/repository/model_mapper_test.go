package repository

import (
	"testing"
	"time"

	"github.com/energimind/identity-server/core/domain/admin"
	"github.com/energimind/identity-server/core/testutil"
	"github.com/stretchr/testify/require"
)

func Test_allUserFieldsAreMapped(t *testing.T) {
	testutil.CheckAllFieldsAreMapped(t, admin.Application{}, dbApplication{})
	testutil.CheckAllFieldsAreMapped(t, admin.Provider{}, dbProvider{})
	testutil.CheckAllFieldsAreMapped(t, admin.User{}, dbUser{})
	testutil.CheckAllFieldsAreMapped(t, admin.Daemon{}, dbDaemon{})
	testutil.CheckAllFieldsAreMapped(t, admin.APIKey{}, dbAPIKey{})

	testutil.CheckAllFieldsAreMapped(t, dbApplication{}, admin.Application{})
	testutil.CheckAllFieldsAreMapped(t, dbProvider{}, admin.Provider{})
	testutil.CheckAllFieldsAreMapped(t, dbUser{}, admin.User{})
	testutil.CheckAllFieldsAreMapped(t, dbDaemon{}, admin.Daemon{})
	testutil.CheckAllFieldsAreMapped(t, dbAPIKey{}, admin.APIKey{})
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
