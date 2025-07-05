package http

type SubmitFeedback struct {
	Score   int64  `json:"score"`
	Comment string `json:"comment"`
}