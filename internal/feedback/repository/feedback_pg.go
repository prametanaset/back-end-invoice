package repository

import (
	"invoice_project/internal/feedback/domain"

	"gorm.io/gorm"
)

type FeedbackRepository interface {
	SubmitFeedback(fb *domain.Feedback) error
}

type feedbackPg struct {
	db *gorm.DB
}

func NewFeedbackRepository(db *gorm.DB) FeedbackRepository {
	return &feedbackPg{db: db}
}

func (r *feedbackPg) SubmitFeedback(fb *domain.Feedback) error {
	return r.db.Create(fb).Error
}

