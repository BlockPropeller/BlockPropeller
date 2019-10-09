package routes

import (
	"context"

	"chainup.dev/chainup/httpserver/request"
	"chainup.dev/chainup/infrastructure"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// GetProviderSettingsResponse is a response to the get provider settings request.
type GetProviderSettingsResponse struct {
	ProviderSettings *infrastructure.ProviderSettings `json:"provider_settings"`
}

// CreateProviderSettingsRequest holds the request payload for the create provider settings endpoint.
type CreateProviderSettingsRequest struct {
	Label        string                      `json:"label" form:"label" validate:"required"`
	ProviderType infrastructure.ProviderType `json:"provider_type" form:"provider_type" validate:"required,valid"`
	Credentials  string                      `json:"credentials" form:"credentials" validate:"required"`
}

// CreateProviderSettingsResponse is a response to the create provider settings request.
type CreateProviderSettingsResponse struct {
	ProviderSettings *infrastructure.ProviderSettings
}

// ProviderSettings REST Resource for accessing ProviderSettings resource.
type ProviderSettings struct {
	settingsRepo infrastructure.ProviderSettingsRepository
}

// NewProviderSettingsRoutes returns a new ProviderSettings routes instance.
func NewProviderSettingsRoutes(settingsRepo infrastructure.ProviderSettingsRepository) *ProviderSettings {
	return &ProviderSettings{settingsRepo: settingsRepo}
}

// LoadProviderSettings is a middleware for loading ProviderSettings into request context
// as well as checking for correct permissions of an authenticated user.
func (ps *ProviderSettings) LoadProviderSettings(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := request.AuthFromContext(c)
		if auth == nil {
			return echo.ErrUnauthorized.SetInternal(errors.New("missing authenticated account"))
		}

		settingsID := infrastructure.ProviderSettingsID(c.Param("settings_id"))
		if settingsID == infrastructure.NilProviderSettingsID {
			return next(c)
		}

		settings, err := ps.settingsRepo.Find(context.Background(), settingsID)
		if err != nil {
			return echo.ErrNotFound.SetInternal(errors.Wrap(err, "find settings"))
		}

		if settings.AccountID != auth.ID {
			return echo.ErrForbidden.
				SetInternal(errors.Errorf(
					"unauthorized provider settings access: authenticated %s, settings %s",
					auth.ID, settings.ID))
		}

		request.WithProviderSettings(c, settings)

		return next(c)
	}
}

// Get a ProviderSettings.
func (ps *ProviderSettings) Get(c echo.Context) error {
	settings := request.ProviderSettingsFromContext(c)
	if settings == nil {
		return echo.ErrNotFound.SetInternal(errors.New("provider settings not found in context"))
	}

	return c.JSON(200, &GetProviderSettingsResponse{ProviderSettings: settings})
}

// Create a ProviderSettings.
func (ps *ProviderSettings) Create(c echo.Context) error {
	var req CreateProviderSettingsRequest
	if err := request.Parse(c, &req); err != nil {
		return err
	}

	acc := request.AuthFromContext(c)

	settings := infrastructure.NewProviderSettings(acc.ID, req.Label, req.ProviderType, req.Credentials)

	err := ps.settingsRepo.Create(context.Background(), settings)
	if err != nil {
		return errors.Wrap(err, "create provider settings")
	}

	return c.JSON(201, &CreateProviderSettingsResponse{ProviderSettings: settings})
}
