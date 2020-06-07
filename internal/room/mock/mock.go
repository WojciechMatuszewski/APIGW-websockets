package mock

//go:generate mockgen -destination=apigw.go -package=mock github.com/aws/aws-sdk-go/service/apigatewaymanagementapi/apigatewaymanagementapiiface ApiGatewayManagementApiAPI
