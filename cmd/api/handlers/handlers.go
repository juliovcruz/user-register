package handlers

import (
	"encoding/json"
	"errors"

	validatorv10 "github.com/go-playground/validator/v10"
	"github.com/juliovcruz/user-register/internal/mailvalidation"
	"github.com/juliovcruz/user-register/internal/security/token"
	"github.com/juliovcruz/user-register/internal/users"
	"github.com/valyala/fasthttp"
)

var validator = validatorv10.New()

type UserHandler struct {
	service      *users.Service
	tokenService *token.Service
}

func NewUserHandler(service *users.Service, tokenService *token.Service) *UserHandler {
	return &UserHandler{service: service, tokenService: tokenService}
}

// JWTMiddleware verifica a validade do token JWT
func (h *UserHandler) JWTMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		tokenString := string(ctx.Request.Header.Peek("Authorization"))
		if tokenString == "" {
			returnError(ctx, errors.New("authorization header required"), fasthttp.StatusUnauthorized)
			return
		}

		tokenString = tokenString[len("Bearer "):]

		valid, err := h.tokenService.IsValid(tokenString)
		if err != nil || !valid {
			returnError(ctx, errors.New("invalid token"), fasthttp.StatusUnauthorized)
			return
		}

		next(ctx)
	}
}

// CreateUser cria um novo usuário
// @Summary Cria um novo usuário
// @Description Cria um usuário com nome, e-mail e senha
// @Tags users
// @Accept json
// @Produce json
// @Param user body users.CreateUser true "Usuário"
// @Success 201 {object} users.User
// @Failure 400 {object} Err
// @Failure 500 {object} Err
// @Router /users [post]
func (h *UserHandler) CreateUser(ctx *fasthttp.RequestCtx) {
	var request users.CreateUser
	if err := json.Unmarshal(ctx.PostBody(), &request); err != nil {
		returnError(ctx, errors.New("invalid request"), fasthttp.StatusBadRequest)
		return
	}

	if err := validator.Struct(request); err != nil {
		returnError(ctx, errors.New("validation failed: "+err.Error()), fasthttp.StatusBadRequest)
		return
	}

	user, err := h.service.Create(ctx, request)
	if err != nil {
		if errors.Is(err, users.ErrPasswordMismatch) {
			returnError(ctx, err, fasthttp.StatusBadRequest)
			return
		}
		if errors.Is(err, users.ErrMailAlreadyExists) {
			returnError(ctx, err, fasthttp.StatusBadRequest)
			return
		}

		returnError(ctx, err, fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	if err := json.NewEncoder(ctx).Encode(user); err != nil {
		returnError(ctx, errors.New("failed to encode response"), fasthttp.StatusInternalServerError)
	}
}

// UpdatePassword atualiza a senha do usuário
// @Summary Atualiza a senha do usuário
// @Description Atualiza a senha do usuário com base no e-mail
// @Tags users
// @Accept json
// @Produce json
// @Param updatePassword body users.UpdatePassword true "Atualizar senha"
// @Success 204
// @Failure 400 {object} Err
// @Failure 500 {object} Err
// @Router /users/password [put]
func (h *UserHandler) UpdatePassword(ctx *fasthttp.RequestCtx) {
	var req users.UpdatePassword
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		returnError(ctx, errors.New("invalid request"), fasthttp.StatusBadRequest)
		return
	}

	if err := validator.Struct(req); err != nil {
		returnError(ctx, errors.New("validation failed: "+err.Error()), fasthttp.StatusBadRequest)
		return
	}

	err := h.service.UpdatePassword(ctx, req)
	if err != nil {
		if errors.Is(err, mailvalidation.ErrRecordNotFound) {
			returnError(ctx, errors.New("don't have code for this email"), fasthttp.StatusBadRequest)
			return
		}

		returnError(ctx, err, fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusNoContent)
}

// Login faz login do usuário
// @Summary Faz login do usuário
// @Description Faz login com o e-mail e a senha do usuário
// @Tags users
// @Accept json
// @Produce json
// @Param login body users.Login true "Fazer login"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} Err
// @Failure 401 {object} Err
// @Router /users/login [post]
func (h *UserHandler) Login(ctx *fasthttp.RequestCtx) {
	var request users.Login
	if err := json.Unmarshal(ctx.PostBody(), &request); err != nil {
		returnError(ctx, errors.New("invalid request"), fasthttp.StatusBadRequest)
		return
	}

	if err := validator.Struct(request); err != nil {
		returnError(ctx, errors.New("validation failed: "+err.Error()), fasthttp.StatusBadRequest)
		return
	}

	token, err := h.service.Login(ctx, request.Email, request.Password)
	if err != nil {
		returnError(ctx, err, fasthttp.StatusUnauthorized)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	response := map[string]string{"token": token}
	if err := json.NewEncoder(ctx).Encode(response); err != nil {
		returnError(ctx, errors.New("failed to encode response"), fasthttp.StatusInternalServerError)
	}
}

// ForgotPassword inicia o processo de recuperação de senha
// @Summary Inicia recuperação de senha
// @Description Envia um e-mail para recuperação de senha
// @Tags users
// @Accept json
// @Produce json
// @Param email body users.ForgotPassword true "Email para recuperação"
// @Success 204
// @Failure 400 {object} Err
// @Failure 500 {object} Err
// @Router /users/forgot_password [post]
func (h *UserHandler) ForgotPassword(ctx *fasthttp.RequestCtx) {
	var request users.ForgotPassword
	if err := json.Unmarshal(ctx.PostBody(), &request); err != nil {
		returnError(ctx, errors.New("invalid request"), fasthttp.StatusBadRequest)
		return
	}

	if err := validator.Struct(request); err != nil {
		returnError(ctx, errors.New("validation failed: "+err.Error()), fasthttp.StatusBadRequest)
		return
	}

	err := h.service.ForgotPassword(ctx, request.Email)
	if err != nil {
		returnError(ctx, err, fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusNoContent)
}

// ListUsers lista todos os usuários
// @Summary Lista usuários
// @Description Lista todos os usuários com limite e deslocamento utilizar header "Authorization": "Bearer {token}"
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Limite de usuários Padrão: 10"
// @Param offset query int false "Deslocamento - Padrão: 0"
// @Success 200 {array} users.User
// @Failure 400 {object} Err
// @Failure 500 {object} Err
// @Router /users [get]
func (h *UserHandler) ListUsers(ctx *fasthttp.RequestCtx) {
	limit, offset, err := getLimitAndOffSet(ctx)
	if err != nil {
		returnError(ctx, err, fasthttp.StatusBadRequest)
		return
	}

	users, err := h.service.List(ctx, limit, offset)
	if err != nil {
		returnError(ctx, err, fasthttp.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(ctx).Encode(users); err != nil {
		returnError(ctx, errors.New("failed to encode response"), fasthttp.StatusInternalServerError)
	}
}

func getLimitAndOffSet(ctx *fasthttp.RequestCtx) (int, int, error) {
	limit := ctx.QueryArgs().GetUintOrZero("limit")

	offset := ctx.QueryArgs().GetUintOrZero("offset")

	if limit < 0 {
		return 0, 0, errors.New("invalid limit")
	}
	if offset < 0 {
		return 0, 0, errors.New("invalid offset")
	}

	if limit == 0 {
		limit = 10
	}

	return limit, offset, nil
}

func returnError(ctx *fasthttp.RequestCtx, err error, statusCode int) {
	ctx.SetStatusCode(statusCode)
	if err := json.NewEncoder(ctx).Encode(Err{Error: err.Error()}); err != nil {
		ctx.Error("failed to encode response", fasthttp.StatusInternalServerError)
	}
}

type Err struct {
	Error string `json:"error"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func CorsMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if string(ctx.Method()) == "OPTIONS" {
			ctx.SetStatusCode(fasthttp.StatusOK)
			return
		}

		next(ctx)
	}
}
