package server

import (
	"fmt"

	"notes/internal/config"
	"notes/internal/handler"
)

// Структура сервера
type Server struct {
	// Конфигурация сервера
	cfg *config.Config
}

// NewServer - конструктор сервера
func NewServer(cfg *config.Config) (*Server, error) {
	// Проверяем, что конфигурация не пустая
	if cfg == nil {
		return nil, fmt.Errorf("конфигурация сервера не может быть nil")
	}
	// Создаем новый экземпляр обработчика
	handler := handler.NewHandler(cfg)
	// Проверяем, что обработчик успешно создан
	if handler == nil {
		return nil, fmt.Errorf("не удалось создать обработчик сервера")
	}
	fmt.Println("Обработчик сервера успешно создан")

	// Создаем новый экземпляр сервера
	return &Server{
		cfg: cfg,
	}, nil
}

// Start - запуск сервера
func (s *Server) Start() error {
	fmt.Printf("Сервер запускается на %s:%s\n", s.cfg.Host, s.cfg.Port)
	return nil
}

// Stop - остановка сервера
func (s *Server) Stop() error {
	fmt.Println("Сервер остановлен")
	return nil
}

// Serve - основной метод сервера
func (s *Server) Serve() error {
	if err := s.Start(); err != nil {
		return err
	}

	// Запускаем сервер через Gin
	address := fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port)
	fmt.Printf("Сервер готов к обработке запросов на %s...\n", address)
	return nil
}
