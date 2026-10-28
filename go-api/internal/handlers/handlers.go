package handlers

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/interseguro/go-api/internal/clients"
	"github.com/interseguro/go-api/internal/models"
	"github.com/interseguro/go-api/internal/services"
)

// Handler agrupa las dependencias de los handlers HTTP.
type Handler struct {
	nodeClient *clients.NodeClient
}

func New(nodeClient *clients.NodeClient) *Handler {
	return &Handler{nodeClient: nodeClient}
}

// QR maneja POST /api/v1/qr
// 1. Recibe la matriz original.
// 2. Calcula su factorización QR.
// 3. Envía Q y R a Node.js para obtener estadísticas.
// 4. Devuelve QR + estadísticas.
func (h *Handler) QR(c *fiber.Ctx) error {
	var req models.MatrixRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "JSON inválido: se esperaba {\"matrix\": [[...]]}",
		})
	}

	qr, err := services.FactorizeQR(req.Matrix)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 8*time.Second)
	defer cancel()

	stats, err := h.nodeClient.GetStatistics(ctx, [][][]float64{qr.Q.Data, qr.R.Data})
	if err != nil {
		log.Printf("error obteniendo estadísticas de Node: %v", err)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "no fue posible obtener estadísticas desde Node API",
		})
	}

	return c.JSON(models.QRResponse{
		Operation:  "QR_FACTORIZATION",
		Original:   services.MatrixToResponse(req.Matrix),
		QR:         *qr,
		Statistics: *stats,
	})
}

// Rotate maneja POST /api/v1/rotate
// 1. Recibe la matriz original.
// 2. La rota 90° en sentido horario.
// 3. Envía el resultado a Node.js para obtener estadísticas.
// 4. Devuelve la matriz rotada + estadísticas.
func (h *Handler) Rotate(c *fiber.Ctx) error {
	var req models.MatrixRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "JSON inválido: se esperaba {\"matrix\": [[...]]}",
		})
	}

	rotated, err := services.Rotate90CW(req.Matrix)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 8*time.Second)
	defer cancel()

	stats, err := h.nodeClient.GetStatistics(ctx, [][][]float64{rotated})
	if err != nil {
		log.Printf("error obteniendo estadísticas de Node: %v", err)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "no fue posible obtener estadísticas desde Node API",
		})
	}

	return c.JSON(models.RotateResponse{
		Operation:  "ROTATE_90_CW",
		Original:   services.MatrixToResponse(req.Matrix),
		Rotated:    services.MatrixToResponse(rotated),
		Statistics: *stats,
	})
}