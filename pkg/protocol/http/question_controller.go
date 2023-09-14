package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"question-and-answers/pkg/api/question"

	"github.com/labstack/echo/v4"
)

type QuestionController struct {
	questionService question.QuestionServiceServer
}

func (c *QuestionController) GetQuestionByID(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid ID")
	}
	result, err := c.questionService.GetQuestionByID(
		ctx.Request().Context(),
		&question.QuestionRequest{Id: id},
	)
	if err := HandleGRPCError(err); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, result)
}

func (c *QuestionController) ListQuestions(ctx echo.Context) error {
	var err error
	var authorId *uint64
	authorIdQuery := ctx.QueryParam("authorId")
	if authorIdQuery != "" {
		if id, err := strconv.ParseUint(authorIdQuery, 10, 32); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid authorId query")
		} else {
			authorId = &id
		}
	}
	var response *question.QuestionList
	response, err = c.questionService.ListQuestions(
		ctx.Request().Context(),
		&question.ListQuestionRequest{
			AuthorId: authorId,
		},
	)
	if err := HandleGRPCError(err); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, response)
}

func (c *QuestionController) CreateQuestion(ctx echo.Context) error {
	var body question.QuestionBody
	if err := json.NewDecoder(ctx.Request().Body).Decode(&body); err != nil {
		slog.ErrorContext(ctx.Request().Context(), err.Error())
		return echo.ErrBadRequest
	}
	response, err := c.questionService.CreateQuestion(
		ctx.Request().Context(),
		&question.CreateQuestionRequest{
			Body: &body,
		},
	)
	if err := HandleGRPCError(err); err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, response)
}

func (c *QuestionController) UpdateQuestion(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid ID")
	}
	var body question.QuestionBody
	if err := json.NewDecoder(ctx.Request().Body).Decode(&body); err != nil {
		slog.ErrorContext(ctx.Request().Context(), err.Error())
		return echo.ErrBadRequest
	}
	response, err := c.questionService.UpdateQuestion(
		ctx.Request().Context(),
		&question.UpdateQuestionRequest{
			Id:   id,
			Item: &body,
		},
	)
	if err := HandleGRPCError(err); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, response)
}

func (c *QuestionController) DeleteQuestion(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid ID")
	}
	response, err := c.questionService.DeleteQuestion(
		ctx.Request().Context(),
		&question.DeleteQuestionRequest{Id: id},
	)
	if err := HandleGRPCError(err); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, response)
}
