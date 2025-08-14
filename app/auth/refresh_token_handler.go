package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type RefreshTokenRequest struct {
	RefreshToken string `cookie:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type RefreshTokenHandler struct {
	repo Repository
	ts   TokenService
}

func NewRefreshTokenHandler(repo Repository, ts TokenService) *RefreshTokenHandler {
	return &RefreshTokenHandler{
		repo: repo,
		ts:   ts,
	}
}

// @Summary		Refresh access token
// @Description	Generate a new access token using a valid refresh token.
// @Description	The API takes refresh token from the cookie and returns a new access token.
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Success		200	{object}	RefreshTokenResponse
// @Failure		401
// @Failure		500
// @Router			/auth/refresh [post]
func (h *RefreshTokenHandler) Handle(ctx context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, int, error) {

	payload, err := h.ts.ValidateAuthRefreshToken(req.RefreshToken)
	if err != nil {
		if !errors.Is(err, domain.ErrInternalServer) {
			return nil, http.StatusUnauthorized, err
		}
		return nil, http.StatusInternalServerError, domain.ErrInternalServer
	}

	exist, err := h.repo.RefreshTokenExists(ctx, req.RefreshToken)
	if err != nil {
		return nil, http.StatusInternalServerError, domain.ErrInternalServer
	}

	if !exist {
		return nil, http.StatusUnauthorized, domain.ErrNotExistRefreshToken
	}

	accessToken, err := h.ts.GenerateAuthAccessToken(payload.UserID, payload.Role)
	if err != nil {
		return nil, http.StatusInternalServerError, domain.ErrInternalServer
	}

	return &RefreshTokenResponse{
		AccessToken: accessToken,
	}, http.StatusOK, nil
}
