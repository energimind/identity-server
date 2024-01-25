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
