package models

import (
	"database/sql"
	"question-and-answers/pkg/api/question"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type Question struct {
	gorm.Model
	QuestionText     string
	QuestionAuthorId uint
	AnswerText       sql.NullString
	AnswerAuthorId   *uint
	Published        bool `gorm:"index;default:false"`
}

func (q *Question) ToGrpcModel() *question.Question {
	var answerAuthorId *uint64
	if q.AnswerAuthorId != nil {
		id := uint64(*q.AnswerAuthorId)
		answerAuthorId = &id
	}
	var answerText *string
	if q.AnswerText.Valid {
		answerText = &q.AnswerText.String
	}
	return &question.Question{
		Id:               uint64(q.ID),
		DateCreated:      timestamppb.New(q.CreatedAt),
		QuestionText:     q.QuestionText,
		QuestionAuthorId: uint64(q.QuestionAuthorId),
		AnswerText:       answerText,
		AnswerAuthorId:   answerAuthorId,
		Published:        &q.Published,
	}
}

func (q *Question) UpdateWith(model *question.QuestionBody) {
	var answerAuthorId *uint
	if model.AnswerAuthorId != nil {
		id := uint(*model.AnswerAuthorId)
		answerAuthorId = &id
	}
	var answerText string
	if model.AnswerText != nil {
		answerText = *model.AnswerText
	}
	var published bool
	if model.Published != nil {
		published = *model.Published
	}
	q.QuestionText = model.QuestionText
	q.QuestionAuthorId = uint(model.QuestionAuthorId)
	q.AnswerText = sql.NullString{
		String: answerText,
		Valid:  model.AnswerText != nil && answerText != "",
	}
	q.AnswerAuthorId = answerAuthorId
	q.Published = published
}

func FromGrpcQuestionBody(model *question.QuestionBody) *Question {
	var q Question
	q.UpdateWith(model)
	return &q
}
