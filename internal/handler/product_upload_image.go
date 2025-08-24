package handler

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

func UploadProductImageHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Image data not found",
		})
	}

	ext := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid image format. Only .jpg, .jpeg, .png, .webp are allowed",
		})
	}

	contentType := file.Header.Get("Content-Type")
	allowedContentType := map[string]bool{
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
	}

	if !allowedContentType[contentType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid image content type. Only image/jpg, image/webp are allowed",
		})
	}

	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("product_%d%s", timestamp, filepath.Ext(file.Filename))
	uploadPath := "./storage/images/products/" + filename
	err = c.SaveFile(file, uploadPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to save image",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Image uploaded successfully",
		"data":    filename,
	})
}
