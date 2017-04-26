package http

import (
	"github.com/sparcs-kaist/whale"

	"net/http"
)

// Server implements the whale.Server interface
type Server struct {
	BindAddress            string
	AssetsPath             string
	AuthDisabled           bool
	SSOID                  string
	SSOKey                 string
	EndpointManagement     bool
	UserService            whale.UserService
	EndpointService        whale.EndpointService
	ResourceControlService whale.ResourceControlService
	CryptoService          whale.CryptoService
	JWTService             whale.JWTService
	FileService            whale.FileService
	Settings               *whale.Settings
	TemplatesURL           string
	Handler                *Handler
}

// Start starts the HTTP server
func (server *Server) Start() error {
	middleWareService := &middleWareService{
		jwtService:   server.JWTService,
		authDisabled: server.AuthDisabled,
	}

	var authHandler = NewAuthHandler(middleWareService)
	authHandler.UserService = server.UserService
	authHandler.CryptoService = server.CryptoService
	authHandler.JWTService = server.JWTService
	authHandler.EndpointService = server.EndpointService
	authHandler.ssoClient = NewSSOClient(server.SSOID, server.SSOKey)
	authHandler.authDisabled = server.AuthDisabled
	var userHandler = NewUserHandler(middleWareService)
	userHandler.UserService = server.UserService
	userHandler.CryptoService = server.CryptoService
	userHandler.ResourceControlService = server.ResourceControlService
	var settingsHandler = NewSettingsHandler(middleWareService)
	settingsHandler.settings = server.Settings
	var templatesHandler = NewTemplatesHandler(middleWareService)
	templatesHandler.containerTemplatesURL = server.TemplatesURL
	var dockerHandler = NewDockerHandler(middleWareService, server.ResourceControlService)
	dockerHandler.EndpointService = server.EndpointService
	var websocketHandler = NewWebSocketHandler()
	websocketHandler.EndpointService = server.EndpointService
	var endpointHandler = NewEndpointHandler(middleWareService)
	endpointHandler.authorizeEndpointManagement = server.EndpointManagement
	endpointHandler.EndpointService = server.EndpointService
	endpointHandler.FileService = server.FileService
	var uploadHandler = NewUploadHandler(middleWareService)
	uploadHandler.FileService = server.FileService
	var fileHandler = newFileHandler(server.AssetsPath)

	server.Handler = &Handler{
		AuthHandler:      authHandler,
		UserHandler:      userHandler,
		EndpointHandler:  endpointHandler,
		SettingsHandler:  settingsHandler,
		TemplatesHandler: templatesHandler,
		DockerHandler:    dockerHandler,
		WebSocketHandler: websocketHandler,
		FileHandler:      fileHandler,
		UploadHandler:    uploadHandler,
	}

	return http.ListenAndServe(server.BindAddress, server.Handler)
}
