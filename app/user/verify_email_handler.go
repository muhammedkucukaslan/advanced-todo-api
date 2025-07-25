package user

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type VerifiyEmailRequest struct {
	Language string `reqHeader:"response-language" validate:"required,oneof=tr en ar" swaggerignore:"true"`
	Token    string `json:"token" validate:"required"`
}
type VerifyEmailResponse struct{}

type VerifyEmailHandler struct {
	repo     Repository
	validate *validator.Validate
	ts       TokenService
}

func NewVerifyEmailHandler(repo Repository, validate *validator.Validate, ts TokenService) *VerifyEmailHandler {
	return &VerifyEmailHandler{
		repo:     repo,
		validate: validate,
		ts:       ts,
	}
}

// VerifyEmailHandler is responsible for handling the verification of a user's email address.
//
//	@Summary		Verify user email
//	@Description	Verifies a user's email address using a token
//	@Tags			3- User
//	@Accept			json
//	@Produce		json
//	@Param			response-language	header	string				true	"Response Language"	Enums(tr, en, ar)
//	@Param			request				body	VerifiyEmailRequest	true	"Verify Email Request"
//	@Success		204
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/users/verify-email [post]
func (h *VerifyEmailHandler) Handle(ctx context.Context, req *VerifiyEmailRequest) (*VerifyEmailResponse, int, error) {
	if err := h.validate.Struct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	email, err := h.ts.ValidateVerifyEmailToken(req.Token)
	if err != nil {
		return nil, http.StatusUnauthorized, domain.ErrUnauthorized
	}

	if err := h.repo.VerifyEmail(ctx, email); err != nil {

		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusNoContent, nil
}
