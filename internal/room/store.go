package room

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type Store struct {
	db        dynamodbiface.DynamoDBAPI
	tableName string
	roomName  string
}

func NewStore(db dynamodbiface.DynamoDBAPI, tableName string, roomName string) Store {
	return Store{db: db, tableName: tableName, roomName: roomName}
}

func (s Store) GetRoom(ctx context.Context) (Room, error) {
	var room Room

	exp, err := expression.NewBuilder().WithProjection(expression.NamesList(expression.Name("connections"))).Build()
	if err != nil {
		return room, err
	}

	out, err := s.db.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"room": {
				S: aws.String(s.roomName),
			},
		},
		ExpressionAttributeNames: exp.Names(),
		ProjectionExpression:     exp.Projection(),
		TableName:                aws.String(s.tableName),
	})
	if err != nil {
		return room, err
	}

	err = dynamodbattribute.UnmarshalMap(out.Item, &room)
	return room, err
}

func (s Store) AddToRoom(ctx context.Context, connection string) error {
	in := toStringSet(connection)
	expr, err := expression.NewBuilder().WithUpdate(expression.Add(expression.Name("connections"), expression.Value(in))).Build()
	if err != nil {
		return err
	}

	_, err = s.db.UpdateItemWithContext(ctx, &dynamodb.UpdateItemInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Key: map[string]*dynamodb.AttributeValue{
			"room": {
				S: aws.String(s.roomName),
			},
		},
		TableName:        aws.String(s.tableName),
		UpdateExpression: expr.Update(),
	})

	return err
}

func (s Store) RemoveFromRoom(ctx context.Context, connection string) error {
	in := toStringSet(connection)
	expr, err := expression.NewBuilder().WithUpdate(expression.Delete(expression.Name("connections"), expression.Value(in))).Build()
	if err != nil {
		return err
	}

	_, err = s.db.UpdateItemWithContext(ctx, &dynamodb.UpdateItemInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Key: map[string]*dynamodb.AttributeValue{
			"room": {
				S: aws.String(s.roomName),
			},
		},
		TableName:        aws.String(s.tableName),
		UpdateExpression: expr.Update(),
	})

	return err
}

func toStringSet(connection string) *dynamodb.AttributeValue {
	return (&dynamodb.AttributeValue{}).SetSS(aws.StringSlice([]string{connection}))
}
