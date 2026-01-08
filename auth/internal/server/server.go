package server

import (
	"auth/internal/config"
	"errors"
	"fmt"
)

type Server struct {
	cfg *config.Config
}

func NewServer(cfg *config.Config) (*Server, error) {
	if cfg == nil {
		return nil, errors.New("конфигурация не может быть nil")
	}

	return &Server{
		cfg: cfg,
	}, nil
}

// Stop - остановка сервера
func (s *Server) Stop() error {
	fmt.Println("Сервер остановлен")
	return nil
}

// Serve - основной метод сервера
func (s *Server) Serve() error {
	// Запускаем сервер
	address := fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port)
	fmt.Printf("Сервер готов к обработке запросов на %s...\n", address)
	return nil
}
