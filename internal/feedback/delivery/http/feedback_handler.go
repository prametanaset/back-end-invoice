package http

import (
	"invoice_project/internal/feedback/usecase"
	"invoice_project/pkg/apperror"
	"invoice_project/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FeedbackHandler struct {
	fbUC usecase.FeedbackUsecase
}

func NewInvoiceHandler(fbUC usecase.FeedbackUsecase) *FeedbackHandler {
	return &FeedbackHandler{
		fbUC: fbUC,
	}
}

func (h *FeedbackHandler) Create(c *fiber.Ctx) error {
	var body SubmitFeedback
	if err := c.BodyParser(&body); err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	userID := c.Locals("user_id").(uuid.UUID)

	fb, err := h.fbUC.SubmitFeedback(body.Score, body.Comment, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Feedback submitted successfully",
		"data":    fb, // จะใส่หรือไม่ใส่ก็ได้
	})
}



func (h *FeedbackHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/feedback", middleware.RequireRoles("user", "admin"))
	api.Post("/submit", h.Create)
}
