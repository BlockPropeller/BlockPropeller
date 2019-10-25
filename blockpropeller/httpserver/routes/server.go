package routes

import (
	"context"

	"blockpropeller.dev/blockpropeller/httpserver/request"
	"blockpropeller.dev/blockpropeller/infrastructure"
	"blockpropeller.dev/blockpropeller/provision"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// ListServersResponse is a response to the list servers request.
type ListServersResponse struct {
	Servers []*infrastructure.Server `json:"servers"`
}

// GetServerResponse is a response to the get server request.
type GetServerResponse struct {
	Server *infrastructure.Server `json:"server"`
}

// AddAuthorizedKeyRequest is a request for adding a new authorized ssh key to a server.
type AddAuthorizedKeyRequest struct {
	PublicKey string `json:"public_key" form:"public_key" validate:"required"`
}

// Server REST Resource for accessing server information.
type Server struct {
	srvDestroyer    *provision.ServerDestroyer
	deplProvisioner *provision.DeploymentProvisioner

	srvRepo infrastructure.ServerRepository
}

// NewServerRoutes returns a new Server routes instance.
func NewServerRoutes(srvDestroyer *provision.ServerDestroyer, deplProvisioner *provision.DeploymentProvisioner, srvRepo infrastructure.ServerRepository) *Server {
	return &Server{srvDestroyer: srvDestroyer, deplProvisioner: deplProvisioner, srvRepo: srvRepo}
}

// LoadServer is a middleware for loading Server into request context
// as well as checking for correct permissions of an authenticated user.
func (s *Server) LoadServer(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authAcc := request.AuthFromContext(c)
		if authAcc == nil {
			return echo.ErrForbidden.SetInternal(errors.New("missing authenticated account"))
		}

		srvID := infrastructure.ServerID(c.Param("server_id"))
		srv, err := s.srvRepo.Find(context.Background(), srvID)
		if err != nil {
			return echo.ErrInternalServerError.SetInternal(err)
		}
		if authAcc.ID != srv.AccountID {
			return echo.ErrForbidden.
				SetInternal(errors.Errorf("unauthorized server access: authenticated %s, server %s",
					authAcc.ID, srvID))
		}

		request.WithServer(c, srv)

		return next(c)
	}
}

// List all Servers for an Account.
func (s *Server) List(c echo.Context) error {
	acc := request.AuthFromContext(c)
	if acc == nil {
		return echo.ErrForbidden
	}

	servers, err := s.srvRepo.List(context.Background(), acc.ID)
	if err != nil {
		return errors.Wrap(err, "list servers")
	}

	return c.JSON(200, &ListServersResponse{
		Servers: servers,
	})
}

// Get a Server.
func (s *Server) Get(c echo.Context) error {
	srv := request.ServerFromContext(c)
	if srv == nil {
		return echo.ErrNotFound.SetInternal(errors.New("server not found in context"))
	}

	return c.JSON(200, &GetServerResponse{Server: srv})
}

// AddAuthorizedKey to a Server.
func (s *Server) AddAuthorizedKey(c echo.Context) error {
	srv := request.ServerFromContext(c)
	if srv == nil {
		return echo.ErrNotFound.SetInternal(errors.New("server not found in context"))
	}

	var req AddAuthorizedKeyRequest
	if err := request.Parse(c, &req); err != nil {
		return err
	}

	err := s.deplProvisioner.AddAuthorizedKey(srv, req.PublicKey)
	if err != nil {
		return errors.Wrap(err, "add authorized key")
	}

	return c.NoContent(201)
}

// Delete issues a delete request for a specific Server.
func (s *Server) Delete(c echo.Context) error {
	srv := request.ServerFromContext(c)
	if srv == nil {
		return echo.ErrNotFound.SetInternal(errors.New("server not found in context"))
	}

	err := s.srvDestroyer.Destroy(context.Background(), srv)
	if err != nil {
		return errors.Wrap(err, "destroy server")
	}

	return c.NoContent(204)
}
