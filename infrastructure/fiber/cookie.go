package fiber

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/muhammedkucukaslan/advanced-todo-api/app/auth"
)

type FiberCookieService struct{}

func NewCookieService() *FiberCookieService {
	return &FiberCookieService{}
}

func (s *FiberCookieService) SetRefreshToken(ctx context.Context, claims *auth.RefreshTokenCookieClaims) {
	fiberCtx, _ := ctx.Value(FiberContextKey{}).(*fiber.Ctx)

	fiberCtx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    claims.Token,
		Expires:  time.Now().Add(claims.Duration),
		MaxAge:   int(claims.Duration.Seconds()),
		HTTPOnly: true,
		Secure:   claims.Secure,
		SameSite: fiber.CookieSameSiteStrictMode,
		Path:     "/",
	})
}
