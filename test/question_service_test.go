package test

import (
	"context"
	"question-and-answers/pkg/api/question"
	service "question-and-answers/pkg/service/question"
	"question-and-answers/test/testdata"
	"question-and-answers/test/utils"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestServiceGetQuestionByID(t *testing.T) {
	// arrange
	s := test_utils.NewSuite(t)
	defer s.Close()

	questionService := service.NewQuestionServiceServer(s.DB)

	q := testdata.QuestionList[0]

	questionMockRows := testdata.QuestionToRow(sqlmock.NewRows(testdata.QuestionColumn), q)
	s.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "questions"`)).
		WithArgs(q.ID).
		WillReturnRows(questionMockRows)

	// act
	_, err := questionService.GetQuestionByID(
		context.Background(), &question.QuestionRequest{Id: uint64(q.ID)},
	)

	// assert
	if err != nil {
		t.Errorf("Failed to get question by ID: %v", err)
		t.FailNow()
	}
	if err := s.Mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestServiceListQuestions(t *testing.T) {
	// arrange
	s := test_utils.NewSuite(t)
	defer s.Close()

	questionService := service.NewQuestionServiceServer(s.DB)

	list := testdata.QuestionList
	questionMockRows := sqlmock.NewRows(testdata.QuestionColumn)
	for _, item := range list {
		testdata.QuestionToRow(questionMockRows, item)
	}

	s.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "questions"`)).
		WillReturnRows(questionMockRows)

	// act
	_, err := questionService.ListQuestions(
		context.Background(), &question.ListQuestionRequest{},
	)

	// assert
	if err != nil {
		t.Errorf("Failed to list questions: %v", err)
		t.FailNow()
	}
	if err := s.Mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestServiceListQuestionsWithAuthorFiltering(t *testing.T) {
	// arrange
	s := test_utils.NewSuite(t)
	defer s.Close()

	questionService := service.NewQuestionServiceServer(s.DB)

	list := testdata.QuestionList
	questionMockRows := sqlmock.NewRows(testdata.QuestionColumn)
	for _, item := range list {
		testdata.QuestionToRow(questionMockRows, item)
	}

	questionAuthorId := uint64(list[0].QuestionAuthorId)
	s.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "questions"`)).
		WithArgs(questionAuthorId).
		WillReturnRows(questionMockRows)

	// act
	_, err := questionService.ListQuestions(
		context.Background(), &question.ListQuestionRequest{AuthorId: &questionAuthorId},
	)

	// assert
	if err != nil {
		t.Errorf("Failed to list question by authorId: %v", err)
		t.FailNow()
	}
	if err := s.Mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestServiceCreateQuestion(t *testing.T) {
	// arrange
	s := test_utils.NewSuite(t)
	defer s.Close()

	questionService := service.NewQuestionServiceServer(s.DB)

	q := testdata.QuestionList[0]
	questionMockRows := testdata.QuestionToRow(sqlmock.NewRows(testdata.QuestionColumn), q)
	requestBody := &question.QuestionBody{
		QuestionText:     q.QuestionText,
		QuestionAuthorId: uint64(q.QuestionAuthorId),
	}

	s.Mock.ExpectBegin()
	s.Mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "questions"`)).
		WithArgs(
			sqlmock.AnyArg(), test_utils.NewAlmostNow(1), nil,
			requestBody.QuestionText, requestBody.QuestionAuthorId, nil, nil, false,
		).
		WillReturnRows(questionMockRows)
	s.Mock.ExpectCommit()

	// act
	_, err := questionService.CreateQuestion(
		context.Background(), &question.CreateQuestionRequest{Body: requestBody},
	)

	// assert
	if err != nil {
		t.Errorf("Failed to create question: %v", err)
		t.FailNow()
	}
	if err := s.Mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestServiceUpdateQuestion(t *testing.T) {
	// arrange
	s := test_utils.NewSuite(t)
	defer s.Close()

	var (
		answerText              = "Answer 1"
		answerQuestionId uint64 = 2
		published               = true
		q                       = testdata.QuestionList[0]
	)
	questionService := service.NewQuestionServiceServer(s.DB)
	questionMockRows := testdata.QuestionToRow(sqlmock.NewRows(testdata.QuestionColumn), q)
	requestBody := &question.QuestionBody{
		QuestionText:     q.QuestionText,
		QuestionAuthorId: uint64(q.QuestionAuthorId),
		AnswerText:       &answerText,
		AnswerAuthorId:   &answerQuestionId,
		Published:        &published,
	}

	s.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "questions"`)).
		WillReturnRows(questionMockRows)
	s.Mock.ExpectBegin()
	s.Mock.ExpectExec(regexp.QuoteMeta(`UPDATE "questions" SET`)).
		WithArgs(
			test_utils.AnyTime{}, test_utils.NewAlmostNow(1), nil,
			requestBody.QuestionText, requestBody.QuestionAuthorId,
			answerText, answerQuestionId, published, q.ID,
		).
		WillReturnResult(sqlmock.NewResult(int64(q.ID), 1))
	s.Mock.ExpectCommit()

	// act
	_, err := questionService.UpdateQuestion(
		context.Background(), &question.UpdateQuestionRequest{
			Id:   uint64(q.ID),
			Item: requestBody,
		},
	)

	// assert
	if err != nil {
		t.Errorf("Failed to update question: %v", err)
		t.FailNow()
	}
	if err := s.Mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestServiceDeleteQuestion(t *testing.T) {
	// arrange
	s := test_utils.NewSuite(t)
	defer s.Close()

	q := testdata.QuestionList[0]
	questionService := service.NewQuestionServiceServer(s.DB)
	questionMockRows := testdata.QuestionToRow(sqlmock.NewRows(testdata.QuestionColumn), q)

	s.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "questions"`)).
		WillReturnRows(questionMockRows)
	s.Mock.ExpectBegin()
	s.Mock.ExpectExec(regexp.QuoteMeta(`UPDATE "questions" SET`)).
		WithArgs(test_utils.NewAlmostNow(1), q.ID).
		WillReturnResult(sqlmock.NewResult(int64(q.ID), 1))
	s.Mock.ExpectCommit()

	// act
	_, err := questionService.DeleteQuestion(
		context.Background(), &question.DeleteQuestionRequest{Id: uint64(q.ID)},
	)

	// assert
	if err != nil {
		t.Errorf("Failed to delete question: %v", err)
		t.FailNow()
	}
	if err := s.Mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}
