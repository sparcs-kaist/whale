package http

import (
	"github.com/portainer/portainer"

	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

// AuthHandler represents an HTTP API handler for managing authentication.
type AuthHandler struct {
	*mux.Router
	Logger          *log.Logger
	authDisabled    bool
	ssoClient       *SSOClient
	UserService     portainer.UserService
	CryptoService   portainer.CryptoService
	JWTService      portainer.JWTService
	EndpointService portainer.EndpointService
}

const (
	// ErrInvalidCredentialsFormat is an error raised when credentials format is not valid
	ErrInvalidCredentialsFormat = portainer.Error("Invalid credentials format")
	// ErrInvalidCredentials is an error raised when credentials for a user are invalid
	ErrInvalidCredentials = portainer.Error("Invalid credentials")
	// ErrAuthDisabled is an error raised when trying to access the authentication endpoints
	// when the server has been started with the --no-auth flag
	ErrAuthDisabled = portainer.Error("Authentication is disabled")
)

// NewAuthHandler returns a new instance of AuthHandler.
func NewAuthHandler(mw *middleWareService) *AuthHandler {
	h := &AuthHandler{
		Router: mux.NewRouter(),
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}
	h.Handle("/auth/local", mw.public(http.HandlerFunc(h.handlePostAuth)))
	h.Handle("/auth/sso",
		mw.public(http.HandlerFunc(h.handleSSOAuthInit))).Methods(http.MethodGet)
	h.Handle("/auth/sso",
		mw.public(http.HandlerFunc(h.handleSSOAuth))).Methods(http.MethodPost)
	return h
}

func (handler *AuthHandler) handlePostAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handleNotAllowed(w, []string{http.MethodPost})
		return
	}

	if handler.authDisabled {
		Error(w, ErrAuthDisabled, http.StatusServiceUnavailable, handler.Logger)
		return
	}

	var req postAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, handler.Logger)
		return
	}

	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		Error(w, ErrInvalidCredentialsFormat, http.StatusBadRequest, handler.Logger)
		return
	}

	var username = req.Username
	var password = req.Password

	u, err := handler.UserService.UserByUsername(username)
	if err == portainer.ErrUserNotFound {
		Error(w, err, http.StatusNotFound, handler.Logger)
		return
	} else if err != nil {
		Error(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	err = handler.CryptoService.CompareHashAndData(u.Password, password)
	if err != nil {
		Error(w, ErrInvalidCredentials, http.StatusUnprocessableEntity, handler.Logger)
		return
	}

	tokenData :=  portainer.TokenData{
		ID:       u.ID,
		Username: u.Username,
		Role:     u.Role,
	}
	token, err := handler.JWTService.GenerateToken(tokenData)
	if err != nil {
		Error(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	encodeJSON(w, &postAuthResponse{JWT: token}, handler.Logger)
}

func (handler *AuthHandler) handleSSOAuthInit(w http.ResponseWriter, r *http.Request) {
	if handler.authDisabled {
		Error(w, ErrAuthDisabled, http.StatusServiceUnavailable, handler.Logger)
		return
	}

	params, err := handler.ssoClient.GetLoginParams()
	if err != nil {
		Error(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	encodeJSON(w, &ssoAuthURI{URI: params.URI, State: params.State}, handler.Logger)
}

func (handler *AuthHandler) handleSSOAuth(w http.ResponseWriter, r *http.Request) {
	if handler.authDisabled {
		Error(w, ErrAuthDisabled, http.StatusServiceUnavailable, handler.Logger)
		return
	}

	var req ssoAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, handler.Logger)
		return
	}

	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		Error(w, ErrInvalidRequestFormat, http.StatusBadRequest, handler.Logger)
		return
	}

	info, err := handler.ssoClient.GetUserInfo(req.Code)
	if err != nil {
		Error(w, err, http.StatusBadRequest, handler.Logger)
		return
	}

	username := info["sparcs_id"].(string)
	if username == "" {
		Error(w, err, http.StatusBadRequest, handler.Logger)
		return
	}

	u, err := handler.UserService.UserByUsername(username)
	if err == portainer.ErrUserNotFound {
		u =  portainer.User{
			Username: username,
			Password: "",
			Role:     2,
		}

		err = handler.UserService.CreateUser(u)
		if err != nil {
			Error(w, err, http.StatusInternalServerError, handler.Logger)
			return
		}

		endpoints, err := handler.EndpointService.Endpoints()
		if err != nil {
			Error(w, err, http.StatusInternalServerError, handler.Logger)
			return
		}

		for _, endpoint := range endpoints {
			endpoint.AuthorizedUsers = append(endpoint.AuthorizedUsers, u.ID)
			err = handler.EndpointService.UpdateEndpoint(endpoint.ID, &endpoint)
			if err != nil {
				Error(w, err, http.StatusInternalServerError, handler.Logger)
				return
			}
		}
	} else if err != nil {
		Error(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	tokenData :=  portainer.TokenData{
		ID:       u.ID,
		Username: u.Username,
		Role:     u.Role,
	}
	token, err := handler.JWTService.GenerateToken(tokenData)
	if err != nil {
		Error(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	encodeJSON(w, &postAuthResponse{JWT: token}, handler.Logger)
}

type postAuthRequest struct {
	Username string `valid:"alphanum,required"`
	Password string `valid:"required"`
}

type postAuthResponse struct {
	JWT string `json:"jwt"`
}

type ssoAuthURI struct {
	URI   string `json:"uri"`
	State string `json:"state"`
}

type ssoAuthRequest struct {
	Code string `json:"code"`
}
