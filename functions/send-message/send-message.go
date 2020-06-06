package main

import (
	"context"
	"net/http"
	"os"
	"strings"

	"websockets/internal/room"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func handler(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)
	apigw := apigatewaymanagementapi.New(sess, &aws.Config{Endpoint: aws.String(strings.Join([]string{event.RequestContext.DomainName, event.RequestContext.Stage}, "/"))})
	store := room.NewStore(db, os.Getenv("table_name"), "global")

	r, err := store.GetRoom(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Body: http.StatusText(http.StatusInternalServerError)}, err
	}

	broadcaster := room.NewBroadcaster(apigw, r)
	err = broadcaster.Broadcast(ctx, []byte("message!"))
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Body: http.StatusText(http.StatusInternalServerError)}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: http.StatusText(http.StatusOK)}, nil
}

func main() {
	lambda.Start(handler)
}
