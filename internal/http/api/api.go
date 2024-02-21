package api

import "github.com/lud0m4n/Network/internal/http/usecase"

type Handler struct {
	UseCase *usecase.UseCase
}

func NewHandler(uc *usecase.UseCase) *Handler {
	return &Handler{UseCase: uc}
}
