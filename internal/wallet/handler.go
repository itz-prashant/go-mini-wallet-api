package wallet

import (
	"errors"
	"net/http"
	"strconv"

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

func (h *Handler) HandleGetWallet(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid wallet id")
		return
	}

	wallet, err := h.service.GetWallet(r.Context(), id)

	if err != nil {
		if errors.Is(err, ErrWalletNotFound) {
			response.WriteError(w, http.StatusNotFound, "wallet not found")
			return
		}

		if errors.Is(err, ErrInvalidId) {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		// 3. Unexpected Server Error (500)
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJson(w, http.StatusOK, wallet)
}

func (h *Handler) HandleTransfer(w http.ResponseWriter, r *http.Request) {
	var req TransferRequest

	if err := response.ReadJson(r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	err := h.service.Transfer(r.Context(), req)

	if err != nil {
		if errors.Is(err, ErrWalletNotFound) {
			response.WriteError(w, http.StatusNotFound, "wallet not found")
			return
		}

		if errors.Is(err, ErrInsufficientBalance) ||
			errors.Is(err, ErrInvalidAmount) ||
			errors.Is(err, ErrSameWalletTransfer) ||
			errors.Is(err, ErrInvalidId) {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJson(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "fund transfer successfully",
	})
}

func (h *Handler) HandleGetTransactions(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	walletId, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid wallet id")
		return
	}

	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	limit, _ := strconv.Atoi(query.Get("limit"))
	txType := query.Get("type")

	filter := TransactionFilter{
		WalletID: walletId,
		Type:     txType,
		Page:     page,
		Limit:    limit,
	}

	transactions, err := h.service.GetTransactions(r.Context(), filter)

	if err != nil {
		if errors.Is(err, ErrInvalidId) {
			response.WriteError(w, http.StatusBadRequest, err.Error())
		}
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJson(w, http.StatusOK, transactions)
}
