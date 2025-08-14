package auth

import (
	"context"
	"net/http"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type LogoutRequest struct {
	RefreshToken string `cookie:"refresh_token"`
}

type LogoutResponse struct {
}

type LogoutHandler struct {
	repo Repository
	cs   CookieService
}

func NewLogoutHandler(repo Repository, cs CookieService) *LogoutHandler {
	return &LogoutHandler{
		repo: repo,
		cs:   cs,
	}
}

// Logout removes the refresh token from the database and clears the refresh token cookie.
//
//	@Summary		Logout user
//	@Description	Removes the refresh token from the database and clears the refresh token cookie.
//	@Description	The API takes refresh token from cookies.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Success		204
//	@Failure		500
//	@Router			/auth/logout [post]
func (h *LogoutHandler) Handle(ctx context.Context, req *LogoutRequest) (*LogoutResponse, int, error) {
	if req.RefreshToken == "" {
		return nil, http.StatusBadRequest, domain.ErrInvalidRequest
	}

	if err := h.repo.DeleteRefreshToken(ctx, req.RefreshToken); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	h.cs.RemoveRefreshToken(ctx)
	return nil, http.StatusNoContent, nil
}
