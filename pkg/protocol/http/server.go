package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	service "question-and-answers/pkg/service/question"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type ApiServer struct {
	server *echo.Echo
}

func NewApiServer(db *gorm.DB) *ApiServer {
	server := echo.New()

	server.HideBanner = true
	server.Logger.SetHeader("${time_rfc3339} ${level}")
	server.JSONSerializer = &GrpcJSONSerializer{}

	// Middleware
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	controller := &QuestionController{
		questionService: service.NewQuestionServiceServer(db),
	}

	// Routes
	api := server.Group("/api")
	q := api.Group("/question")
	q.GET("", controller.ListQuestions)
	q.GET("/:id", controller.GetQuestionByID)
	q.POST("", controller.CreateQuestion)
	q.PUT("/:id", controller.UpdateQuestion)
	q.DELETE("/:id", controller.DeleteQuestion)

	return &ApiServer{server: server}
}

func (s *ApiServer) Start(ctx context.Context, serverPort string) error {
	// Start server
	slog.InfoContext(ctx, "starting HTTP server...")
	if err := s.server.Start(":" + serverPort); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *ApiServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	slog.WarnContext(ctx, "shutting down HTTP server...")
	if err := s.server.Shutdown(ctx); err != nil {
		s.server.Logger.Fatal(err)
		return err
	}
	return nil
}
