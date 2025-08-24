package main

import (
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aldngrha/ecommerce-be/internal/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func handleGetFilename(c *fiber.Ctx) error {
	fileNameParams := c.Params("filename")
	filePath := filepath.Join("storage", "images", "products", fileNameParams)

	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return c.Status(http.StatusNotFound).SendString("File not found")
		}

		return c.Status(http.StatusInternalServerError).SendString("Internal server error")
	}

	// open file
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).SendString("Internal server error")
	}

	ext := filepath.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)

	c.Set("Content-Type", mimeType)
	return c.SendStream(file)
}

func main() {

	app := fiber.New()

	app.Use(cors.New())
	app.Get("/storage/images/products/:filename", handleGetFilename)
	app.Post("/products/upload", handler.UploadProductImageHandler)

	app.Listen(":3000")

}
