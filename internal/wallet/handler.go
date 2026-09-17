package wallet

import (
	"errors"
	"net/http"

	"github.com/itz-prashant/mini-wallet-api/internal/utils/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HandleCreateWallet(w http.ResponseWriter, r *http.Request) {
	var req CreateWalletRequest

	if err := response.ReadJson(r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	wallet, err := h.service.CreateWallet(r.Context(), req)

	if err != nil {
		if errors.Is(err, ErrEmptyOwnerName) || errors.Is(err, ErrNegativeBalance) {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJson(w, http.StatusCreated, wallet)
}
