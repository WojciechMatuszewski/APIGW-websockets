package room

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi/apigatewaymanagementapiiface"
)

// Broadcaster broadcast message to all connections within a given Room
type Broadcaster struct {
	apigw apigatewaymanagementapiiface.ApiGatewayManagementApiAPI
	room  Room
}

func NewBroadcaster(apigw apigatewaymanagementapiiface.ApiGatewayManagementApiAPI, room Room) Broadcaster {
	return Broadcaster{apigw: apigw, room: room}
}

func (b Broadcaster) Broadcast(ctx context.Context, data []byte) error {
	for _, c := range b.room.Connections {
		_, err := b.apigw.PostToConnectionWithContext(ctx, &apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: aws.String(c),
			Data:         data,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
