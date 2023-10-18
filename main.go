package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"github.com/caarlos0/env/v9"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type config struct {
	TelegramToken string `env:"TELE_TOKEN,required"`
	ChatId        int64  `env:"CHAT_ID,required"`
}

type Body struct {
	Msg string `json:"msg"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	slog.Info("telegram bot authorized successfully", slog.String("account", bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello friend."))
		w.WriteHeader(http.StatusOK)
	})

	r.Post("/msg", func(w http.ResponseWriter, r *http.Request) {
		var body Body

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			log.Printf("failed to decode request body: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		msg := tgbotapi.NewMessage(
			cfg.ChatId,
			body.Msg,
		)
		if _, err := bot.Send(msg); err != nil {
			log.Printf("failed to send message: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Printf("msg %s sent successfully", body.Msg)
		w.WriteHeader(http.StatusOK)
	})

	if err := http.ListenAndServe(":8765", r); err != nil {
		log.Panicf("cannot start server: %v", err)
	}
}
