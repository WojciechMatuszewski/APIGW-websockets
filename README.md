# Serverless websockets example (Go)

This is an example of creating a simple websocket API using Api-Gateway and SAM.

## Deployment

- make sure you have SAM installed

- build the binaries: `sam build`

- deploy the application: `sam deploy --guided`

## Lambdas

There are 3 lambdas, folder structure should make them explanatory.

## Handling connections

When connection is established, lambda will be invoked with an event, which contains `connectionId`. 
In most of the examples I've seen people store that `connectionId` within `DynamoDB` and then perform `scan` operation to broadcast the message.

The `scan` operation is inefficient and should be avoided. That is why I've used `DynamoDB sets` to append / filter `connectionId` list, avoiding the `scan`.

Some hacks were employed to make it work correctly, like creating an `set` while the underlying attribute does not exist yet.

An example from `internal/room/store.go`
```go
func toStringSet(connection string) *dynamodb.AttributeValue {
	return (&dynamodb.AttributeValue{}).SetSS(aws.StringSlice([]string{connection}))
}
```

## Testing

For testing, I'm trying out something new - `testcontainers`. 
This allows me to spin up `dynamoDB` instance programmatically.
