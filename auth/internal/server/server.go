package server

import (
	"auth/internal/config"
	"auth/internal/errors"
	"auth/internal/handler"
	"auth/internal/routes"
	"auth/internal/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

// Структура сервера
type Server struct {
	// Конфигурация сервера
	cfg *config.Config
	// router - маршрутизатор Gin
	router *gin.Engine // Новое поле для маршрутизатора
}

// NewServer - конструктор сервера
func NewServer(cfg *config.Config) (*Server, error) {
	// Проверяем, что конфигурация не пустая
	if cfg == nil {
		return nil, fmt.Errorf("конфигурация сервера не может быть nil")
	}
	service, err := service.NewService(cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrServiceCreation, err)
	}

	// Создаем новый экземпляр обработчика с базой данных и конфигурацией
	handler := handler.NewHandler(service, cfg)

	// Проверяем, что обработчик успешно создан
	if handler == nil {
		return nil, fmt.Errorf("не удалось создать обработчик сервера")
	}
	fmt.Println("Обработчик сервера успешно создан")
	// Создаем новый экземпляр маршрутизатора
	router := routes.SetupRouter(handler) // Новое
	// Создаем новый экземпляр сервера
	return &Server{
		router: router, // Новое
		cfg:    cfg,
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
	// Используем s.router для запуска сервера
	// Это позволяет нам использовать маршрутизатор, созданный в NewServer
	return s.router.Run(address) // Новое
}
