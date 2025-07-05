package usecase

import (
	"invoice_project/internal/feedback/domain"
	"invoice_project/internal/feedback/repository"

	"github.com/google/uuid"
)

type FeedbackUsecase interface {
	SubmitFeedback(score int64, comment string, userId uuid.UUID) (*domain.Feedback, error)

}

type FeedbackUC struct {
	repo repository.FeedbackRepository
}

func NewFeedbackUsecase(repo repository.FeedbackRepository) FeedbackUsecase {
	return &FeedbackUC{repo: repo}
}

func (u *FeedbackUC) SubmitFeedback(score int64, comment string, userID uuid.UUID) (*domain.Feedback, error) {
	fb := &domain.Feedback{
		Score:    		score,
		Comment:      	comment,
		CreatedByID: 	userID,
	}
	err := u.repo.SubmitFeedback(fb)
	if err != nil {
		return nil, err
	}
	return fb, nil
}