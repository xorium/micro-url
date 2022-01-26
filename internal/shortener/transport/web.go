package transport

import (
	"github.com/labstack/echo/v4"
	"micro-url/internal/shortener"
	"micro-url/internal/shortener/service"
	"net/http"
)

func Start(cfg shortener.Config, svc *service.URLShortener, errCh chan error) {
	e := echo.New()
	registerHandlers(e, svc)
	errCh <- e.Start(cfg.HTTTPAddr)
}

func registerHandlers(e *echo.Echo, svc *service.URLShortener) {
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	e.POST("/shorten", shortenHandler(svc))
	e.GET("/redirect-to/:shorten-id", redirectHandler(svc))
}

func shortenHandler(svc *service.URLShortener) func(c echo.Context) error {
	return func(c echo.Context) error {
		svc.ShortenURL(c.Request().)
	}
}

func redirectHandler(svc *service.URLShortener) func(c echo.Context) error {
	return func(c echo.Context) error {
		svc.ShortenURL(c.Request().)
	}
}
