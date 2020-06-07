package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	cfDynamo "github.com/awslabs/goformation/v4/cloudformation/dynamodb"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDynamo struct {
	DB *dynamodb.DynamoDB
}

// container represents dynamoDB container initialised by the init function
var container testcontainers.Container

func init() {
	if container != nil {
		return
	}

	dynamoC, err := getTestDBContainer()
	if err != nil {
		panic(err.Error())
	}

	container = dynamoC
}

func NewDynamo(t *testing.T) TestDynamo {
	ctx := context.Background()
	t.Helper()

	sess, err := getLocalSession()
	if err != nil {
		t.Fatal(err)
	}

	if container == nil {
		t.Fatal("container is nil!")
	}
	endpoint, err := container.Endpoint(ctx, "http")
	if err != nil {
		t.Fatal(err)
	}

	db := dynamodb.New(sess, &aws.Config{Endpoint: aws.String(endpoint)})
	return TestDynamo{DB: db}
}

type Cleanup func()

func (td TestDynamo) CreateTables(t *testing.T, prefix string) Cleanup {
	tables, err := getDynamoTables()
	if err != nil {
		t.Fatal(err)
	}

	for _, tb := range tables {
		fmt.Println("creating table", prefix+tb.TableName)
		_, err = td.DB.CreateTable(&dynamodb.CreateTableInput{
			TableName:              aws.String(prefix + tb.TableName),
			KeySchema:              parseKeySchema(tb.KeySchema),
			BillingMode:            aws.String(dynamodb.BillingModePayPerRequest),
			AttributeDefinitions:   parseAttributeDefinitions(tb.AttributeDefinitions),
			GlobalSecondaryIndexes: parseGlobalSecondaryIndexes(tb.GlobalSecondaryIndexes),
			LocalSecondaryIndexes:  parseLocalSecondaryIndexes(tb.LocalSecondaryIndexes),
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	return func() {
		for _, tb := range tables {
			_, err = td.DB.DeleteTable(&dynamodb.DeleteTableInput{TableName: aws.String(prefix + tb.TableName)})
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func parseAttributeDefinitions(attrs []cfDynamo.Table_AttributeDefinition) []*dynamodb.AttributeDefinition {
	out := make([]*dynamodb.AttributeDefinition, len(attrs))

	for i, attr := range attrs {
		out[i] = &dynamodb.AttributeDefinition{
			AttributeName: aws.String(attr.AttributeName),
			AttributeType: aws.String(attr.AttributeType),
		}
	}

	return out
}

func parseGlobalSecondaryIndexes(idxs []cfDynamo.Table_GlobalSecondaryIndex) []*dynamodb.GlobalSecondaryIndex {
	if len(idxs) == 0 {
		return nil
	}

	out := make([]*dynamodb.GlobalSecondaryIndex, len(idxs))
	for i, idx := range idxs {
		out[i] = &dynamodb.GlobalSecondaryIndex{
			IndexName: aws.String(idx.IndexName),
			KeySchema: parseKeySchema(idx.KeySchema),
			Projection: &dynamodb.Projection{
				NonKeyAttributes: aws.StringSlice(idx.Projection.NonKeyAttributes),
				ProjectionType:   aws.String(idx.Projection.ProjectionType),
			},
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(idx.ProvisionedThroughput.ReadCapacityUnits),
				WriteCapacityUnits: aws.Int64(idx.ProvisionedThroughput.WriteCapacityUnits),
			},
		}
	}

	return out
}

func parseLocalSecondaryIndexes(idxs []cfDynamo.Table_LocalSecondaryIndex) []*dynamodb.LocalSecondaryIndex {
	if len(idxs) == 0 {
		return nil
	}

	out := make([]*dynamodb.LocalSecondaryIndex, len(idxs))
	for i, idx := range idxs {
		out[i] = &dynamodb.LocalSecondaryIndex{
			IndexName: aws.String(idx.IndexName),
			KeySchema: parseKeySchema(idx.KeySchema),
			Projection: &dynamodb.Projection{
				NonKeyAttributes: aws.StringSlice(idx.Projection.NonKeyAttributes),
				ProjectionType:   aws.String(idx.Projection.ProjectionType),
			},
		}
	}

	return out
}

func parseKeySchema(kschema []cfDynamo.Table_KeySchema) []*dynamodb.KeySchemaElement {
	if len(kschema) == 0 {
		return nil
	}

	out := make([]*dynamodb.KeySchemaElement, len(kschema))
	for i, schema := range kschema {
		out[i] = &dynamodb.KeySchemaElement{
			AttributeName: aws.String(schema.AttributeName),
			KeyType:       aws.String(schema.KeyType),
		}
	}

	return out
}

func getTestDBContainer() (testcontainers.Container, error) {
	ctx := context.Background()

	port, err := nat.NewPort("tcp", "8000")
	if err != nil {
		return nil, err
	}

	req := testcontainers.ContainerRequest{
		Image:        "amazon/dynamodb-local",
		Env:          nil,
		ExposedPorts: []string{"8000/tcp"},
		WaitingFor:   wait.ForListeningPort(port),
	}

	dynamoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	return dynamoC, err
}

func getLocalSession() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region:      aws.String("local"),
		Credentials: credentials.NewStaticCredentials("local", "local", "local"),
	})
}
