package fiber

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type FiberCookieService struct{}

func NewCookieService() *FiberCookieService {
	return &FiberCookieService{}
}

func (s *FiberCookieService) SetRefreshToken(ctx context.Context, token string) {
	fiberCtx, _ := ctx.Value(FiberContextKey{}).(*fiber.Ctx)

	fiberCtx.Cookie(&fiber.Cookie{
		Name:     domain.RefreshTokenCookieName,
		Value:    token,
		Expires:  time.Now().Add(domain.RefreshTokenCookieMaxAge),
		MaxAge:   int(domain.RefreshTokenCookieMaxAge.Seconds()),
		HTTPOnly: true,
		Secure:   domain.CookieSecure,
		SameSite: fiber.CookieSameSiteStrictMode,
		Path:     domain.RefreshCookiePath,
	})
}

func (s *FiberCookieService) SetAccessToken(ctx context.Context, token string) {
	fiberCtx, _ := ctx.Value(FiberContextKey{}).(*fiber.Ctx)

	fiberCtx.Cookie(&fiber.Cookie{
		Name:     domain.AccessTokenCookieName,
		Value:    token,
		Expires:  time.Now().Add(domain.AccessTokenCookieMaxAge),
		MaxAge:   int(domain.AccessTokenCookieMaxAge.Seconds()),
		HTTPOnly: true,
		Secure:   domain.CookieSecure,
		SameSite: fiber.CookieSameSiteStrictMode,
		Path:     domain.AccessTokenCookiePath,
	})
}

func (s *FiberCookieService) RemoveTokens(ctx context.Context) {
	fiberCtx, _ := ctx.Value(FiberContextKey{}).(*fiber.Ctx)

	fiberCtx.Cookie(&fiber.Cookie{
		Name:     domain.RefreshTokenCookieName,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour * 24),
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   domain.CookieSecure,
		SameSite: fiber.CookieSameSiteStrictMode,
		Path:     domain.RefreshCookiePath,
	})

	fiberCtx.Cookie(&fiber.Cookie{
		Name:     domain.AccessTokenCookieName,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour * 24),
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   domain.CookieSecure,
		SameSite: fiber.CookieSameSiteStrictMode,
		Path:     domain.AccessTokenCookiePath,
	})
}
