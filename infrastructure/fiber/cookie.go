package fiber

import (
	"context"
	"fmt"
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

	fmt.Printf("RefreshCookiePath: %s\n", domain.RefreshCookiePath)
	fmt.Printf("RefreshTokenCookieMaxAge: %v\n", domain.RefreshTokenCookieMaxAge)
	fmt.Printf("CookieSecure: %v\n", domain.CookieSecure)

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

	fmt.Printf("AccessCookiePath: %s\n", domain.AccessTokenCookiePath)
	fmt.Printf("AccessTokenCookieMaxAge: %v\n", domain.AccessTokenCookieMaxAge)
	fmt.Printf("CookieSecure: %v\n", domain.CookieSecure)

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

func (s *FiberCookieService) RemoveRefreshToken(ctx context.Context) {
	fiberCtx, _ := ctx.Value(FiberContextKey{}).(*fiber.Ctx)

	fiberCtx.ClearCookie(domain.RefreshTokenCookieName)
}
