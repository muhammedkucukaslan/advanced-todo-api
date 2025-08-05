package healthcheck

import (
	"context"
	"net/http"
)

var (
	IamAlive = "I am alive, alhamdulillah"
)

type HealthcheckRequest struct{}

type HealthcheckResponse struct {
	Status string `json:"status"`
}

type HealthcheckHandler struct{}

// NewHealthcheckHandler creates a new HealthcheckHandler instance.
//
//	@Summary		Healthcheck
//	@Description	Check the health of the service
//	@Tags			Healthcheck
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		500
//	@Router			/healthcheck [get]
func NewHealthcheckHandler() *HealthcheckHandler {
	return &HealthcheckHandler{}
}

func (h *HealthcheckHandler) Handle(ctx context.Context, req *HealthcheckRequest) (*HealthcheckResponse, int, error) {
	return &HealthcheckResponse{
		Status: IamAlive,
	}, http.StatusOK, nil
}
