package room_test

import (
	"context"
	"errors"
	"testing"

	"websockets/internal/room"
	"websockets/internal/room/mock"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"

	"github.com/golang/mock/gomock"
)

func TestBroadcaster_Broadcast(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		apigwAPI := mock.NewMockApiGatewayManagementApiAPI(ctrl)
		r := room.Room{Connections: []string{"a", "b"}}
		broadcaster := room.NewBroadcaster(apigwAPI, r)

		data := []byte("foo")

		gomock.InOrder(
			apigwAPI.EXPECT().PostToConnectionWithContext(ctx, &apigatewaymanagementapi.PostToConnectionInput{
				ConnectionId: aws.String(r.Connections[0]),
				Data:         data,
			}).Return(nil, nil), apigwAPI.EXPECT().PostToConnectionWithContext(ctx, &apigatewaymanagementapi.PostToConnectionInput{
				ConnectionId: aws.String(r.Connections[1]),
				Data:         data,
			}).Return(nil, nil),
		)

		err := broadcaster.Broadcast(ctx, data)
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		apigwAPI := mock.NewMockApiGatewayManagementApiAPI(ctrl)
		r := room.Room{Connections: []string{"a"}}
		broadcaster := room.NewBroadcaster(apigwAPI, r)

		data := []byte("foo")

		apigwAPI.EXPECT().PostToConnectionWithContext(ctx, &apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: aws.String(r.Connections[0]),
			Data:         data,
		}).Return(nil, errors.New("boom"))

		err := broadcaster.Broadcast(ctx, data)
		assert.Error(t, err)
	})
}
