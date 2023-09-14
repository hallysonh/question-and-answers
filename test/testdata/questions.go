package testdata

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
	"question-and-answers/pkg/database/models"
)

func QuestionToRow(row *sqlmock.Rows, q models.Question) *sqlmock.Rows {
	return row.AddRow(
		q.ID, q.CreatedAt, q.UpdatedAt, q.DeletedAt,
		q.QuestionText, q.QuestionAuthorId,
		q.AnswerText, q.AnswerAuthorId,
		q.Published,
	)
}

var QuestionColumn = []string{
	"id", "created_at", "updated_at", "deleted_at",
	"question_text", "question_author_id",
	"answer_text", "answer_author_id", "published",
}

var answerAuthorId uint = 2

var QuestionList = []models.Question{
	{
		Model:            gorm.Model{ID: 1},
		QuestionText:     "Question 1",
		QuestionAuthorId: 1,
		AnswerText:       sql.NullString{String: "Answer 1", Valid: true},
		AnswerAuthorId:   &answerAuthorId,
		Published:        true,
	},
	{
		Model:            gorm.Model{ID: 2},
		QuestionText:     "Question 2",
		QuestionAuthorId: 1,
	},
	{
		Model:            gorm.Model{ID: 3},
		QuestionText:     "Question 3",
		QuestionAuthorId: 2,
	},
}
