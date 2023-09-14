package main

import (
	"context"
	"flag"
	"log"
	"question-and-answers/pkg/api/question"
	"question-and-answers/pkg/config"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var testNum int
	flag.IntVar(&testNum, "test", 1, "Number of the test to run")
	flag.Parse()

	appConfig := config.LoadAppConfiguration()

	// Set up a connection to the server.
	conn, err := grpc.DialContext(
		context.Background(),
		"localhost:"+appConfig.GRPCServerPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	defer func() { _ = conn.Close() }()

	testingQuestionServer(conn, testNum)
}

func testingQuestionServer(conn *grpc.ClientConn, testNum int) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Set up a connection to the server.
	c := question.NewQuestionServiceClient(conn)
	switch testNum {
	case 1:
		testingGetQuestionByID(ctx, c)
	case 2:
		testingListQuestions(ctx, c)
	case 3:
		testingCreateQuestion(ctx, c)
	case 4:
		testingUpdateQuestion(ctx, c)
	case 5:
		testingDeleteQuestion(ctx, c)
	}
}

func testingGetQuestionByID(ctx context.Context, c question.QuestionServiceClient) {
	req := question.QuestionRequest{Id: 1}
	res, err := c.GetQuestionByID(ctx, &req)
	if err != nil {
		log.Fatalf("GetQuestionByID failed: %v", err)
	}
	log.Printf("GetQuestionByID result: <%+v>\n\n", res)
}

func testingListQuestions(ctx context.Context, c question.QuestionServiceClient) {
	req := question.ListQuestionRequest{}
	res, err := c.ListQuestions(ctx, &req)
	if err != nil {
		log.Fatalf("ListQuestions failed: %v", err)
	}
	log.Printf("ListQuestions result: <%+v>\n\n", res)
}

func testingCreateQuestion(ctx context.Context, c question.QuestionServiceClient) {
	req := question.CreateQuestionRequest{
		Body: &question.QuestionBody{
			QuestionText:     "Question 1",
			QuestionAuthorId: 1,
		},
	}
	res, err := c.CreateQuestion(ctx, &req)
	if err != nil {
		log.Fatalf("CreateQuestion failed: %v", err)
	}
	log.Printf("CreateQuestion result: <%+v>\n\n", res)
}

func testingUpdateQuestion(ctx context.Context, c question.QuestionServiceClient) {
	published := true
	req := question.UpdateQuestionRequest{
		Id: 1,
		Item: &question.QuestionBody{
			QuestionText:     "Question 1",
			QuestionAuthorId: 1,
			Published:        &published,
		},
	}
	res, err := c.UpdateQuestion(ctx, &req)
	if err != nil {
		log.Fatalf("UpdateQuestion failed: %v", err)
	}
	log.Printf("UpdateQuestion result: <%+v>\n\n", res)
}

func testingDeleteQuestion(ctx context.Context, c question.QuestionServiceClient) {
	req := question.DeleteQuestionRequest{Id: 1}
	res, err := c.DeleteQuestion(ctx, &req)
	if err != nil {
		log.Fatalf("DeleteQuestion failed: %v", err)
	}
	log.Printf("DeleteQuestion result: <%+v>\n\n", res)
}
