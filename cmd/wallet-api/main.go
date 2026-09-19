package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/itz-prashant/mini-wallet-api/internal/config"
	"github.com/itz-prashant/mini-wallet-api/internal/db"
	"github.com/itz-prashant/mini-wallet-api/internal/wallet"
)


func main() {
	cfg := config.MustLoad()

	log.Printf("[INFO] Configuration loaded for environment: %s", cfg.Env)

	sqliteDb, err := db.New(cfg)

	if err != nil {
		log.Fatalf("[FATAL] Failed to inialtize databse %v", err)
	}

	defer sqliteDb.Db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func (w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]string{"status":"ok"})
	})

	walletRepo := wallet.NewRepository(sqliteDb.Db)
	walletService := wallet.NewService(walletRepo)
	walletHanlder := wallet.NewHandler(walletService)

	mux.HandleFunc("POST /api/v1/wallets", walletHanlder.HandleCreateWallet)
	mux.HandleFunc("GET /api/v1/wallets/{id}", walletHanlder.HandleGetWallet)
	mux.HandleFunc("POST /api/v1/transfers", walletHanlder.HandleTransfer)

	server := http.Server{
		Addr: cfg.Address,
		Handler: mux,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 15 * time.Second,
	}

	go func (){
		log.Printf("[INFO] Server listening on http://localhost%s", cfg.Address)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed){
			log.Fatalf("[FATAL] Server listening failed: %v", err)
		}
	}()
	
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt ,syscall.SIGTERM)
	<-stop
	log.Println("[INFO] Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[INFO], server forced to shutdown %v", err)
	}
}