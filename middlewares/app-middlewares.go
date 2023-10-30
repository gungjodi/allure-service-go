package app_middlewares

import (
	"osp-allure/utils"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	fiberutils "github.com/gofiber/fiber/v2/utils"
)

func AddTrailingSlashes() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		originalUrl := fiberutils.CopyString(c.OriginalURL())

		// Check if the client is requesting a file extension
		extMatch, _ := regexp.MatchString("\\.[a-zA-Z0-9]+$", originalUrl)

		if !strings.HasSuffix(originalUrl, "/") && !extMatch && c.Method() == "GET" {
			c.Redirect(originalUrl + "/")
		}
		return c.Next()
	}
}

func FavIcon() fiber.Handler {
	return favicon.New(favicon.Config{
		File: filepath.Join(utils.GetAllureResourcesPath(), "favicon.ico"),
		URL:  "/favicon.ico",
	})
}

func AddIndex() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		originalUrl := fiberutils.CopyString(c.OriginalURL())

		// Check if the client is requesting a file extension
		extMatch, _ := regexp.MatchString("\\.[a-zA-Z0-9]+$", originalUrl)

		if !extMatch {
			if strings.HasSuffix(originalUrl, "/") {
				originalUrl = originalUrl + "index.html"
			} else {
				originalUrl = originalUrl + "/index.html"
			}

			c.Redirect(originalUrl)
		}
		return c.Next()
	}
}
