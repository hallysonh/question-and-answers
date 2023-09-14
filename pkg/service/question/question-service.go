package question

import (
	"context"
	"errors"
	"log/slog"

	"question-and-answers/pkg/api/question"
	"question-and-answers/pkg/database/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// questionServiceServer is implementation of v1.QuestionServiceServer proto interface
type questionServiceServer struct {
	db *gorm.DB
	question.UnimplementedQuestionServiceServer
}

// NewQuestionServiceServer creates Question service
func NewQuestionServiceServer(db *gorm.DB) question.QuestionServiceServer {
	return &questionServiceServer{db: db}
}

func (s *questionServiceServer) GetQuestionByID(ctx context.Context, request *question.QuestionRequest) (*question.Question, error) {
	// get SQL connection from pool
	var questionDao models.Question
	tx := s.db.WithContext(ctx)
	result := tx.First(&questionDao, request.Id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "question not found")
		}
		slog.ErrorContext(ctx, result.Error.Error())
		return nil, status.Error(codes.Internal, "internal error")
	}
	return questionDao.ToGrpcModel(), nil
}

func (s *questionServiceServer) ListQuestions(ctx context.Context, request *question.ListQuestionRequest) (*question.QuestionList, error) {
	// get SQL connection from pool
	var questionsDao []models.Question
	tx := s.db.WithContext(ctx)
	tx = tx.Order("id asc")
	if request.AuthorId != nil {
		tx = tx.Where("question_author_id = ?", request.AuthorId)
	}
	result := tx.Find(&questionsDao)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "question not found")
		}
		slog.ErrorContext(ctx, result.Error.Error())
		return nil, status.Error(codes.Internal, "internal error")
	}
	var items = make([]*question.Question, len(questionsDao))
	for idx, item := range questionsDao {
		items[idx] = item.ToGrpcModel()
	}
	return &question.QuestionList{Items: items}, nil
}

func (s *questionServiceServer) CreateQuestion(ctx context.Context, request *question.CreateQuestionRequest) (*question.Question, error) {
	questionModel := models.FromGrpcQuestionBody(request.Body)
	tx := s.db.WithContext(ctx)
	result := tx.Create(questionModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, status.Error(codes.AlreadyExists, "question already exists")
		}
		slog.ErrorContext(ctx, result.Error.Error())
		return nil, status.Error(codes.Internal, "internal error")
	}
	return questionModel.ToGrpcModel(), nil
}

func (s *questionServiceServer) UpdateQuestion(ctx context.Context, request *question.UpdateQuestionRequest) (*question.Question, error) {
	// Get by ID from database
	var questionDao models.Question
	tx := s.db.WithContext(ctx)
	result := tx.First(&questionDao, request.Id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "question not found")
		}
		slog.ErrorContext(ctx, result.Error.Error())
		return nil, status.Error(codes.Internal, "internal error")
	}
	// Apply updates
	questionDao.UpdateWith(request.Item)
	// Save
	tx.Save(&questionDao)
	return questionDao.ToGrpcModel(), nil
}

func (s *questionServiceServer) DeleteQuestion(ctx context.Context, request *question.DeleteQuestionRequest) (*question.Question, error) {
	var questionDao models.Question
	tx := s.db.WithContext(ctx)
	result := tx.First(&questionDao, request.Id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "question not found")
		}
		slog.ErrorContext(ctx, result.Error.Error())
		return nil, status.Error(codes.Internal, "internal error")
	}
	tx.Delete(&questionDao)
	return questionDao.ToGrpcModel(), nil
}
