package main

import (
	"context"
	"net/http"
	"os"

	"websockets/internal/room"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func handler(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)

	store := room.NewStore(db, os.Getenv("table_name"), "global")
	err := store.RemoveFromRoom(ctx, event.RequestContext.ConnectionID)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: http.StatusText(http.StatusInternalServerError), StatusCode: http.StatusInternalServerError}, nil
	}

	return events.APIGatewayProxyResponse{Body: http.StatusText(http.StatusOK), StatusCode: http.StatusOK}, nil

}

func main() {
	lambda.Start(handler)
}
