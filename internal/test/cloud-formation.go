package test

import (
	"errors"
	"os"

	"github.com/awslabs/goformation/v4"
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/dynamodb"
)

func getTemplate() (*cloudformation.Template, error) {
	tpath, err := getTemplatePath()
	if err != nil {
		return nil, err
	}

	return goformation.Open(tpath)
}

func getTemplatePath() (string, error) {
	path := "template.yaml"
	for i := 0; i < 10; i++ {
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			path = "../" + path
			continue
		}

		return path, nil
	}

	return "", errors.New("unable to find template")
}

func getDynamoTables() ([]*dynamodb.Table, error) {
	var tables []*dynamodb.Table

	tmpl, err := getTemplate()
	if err != nil {
		return nil, err
	}

	for _, dbTable := range tmpl.GetAllDynamoDBTableResources() {
		tables = append(tables, dbTable)
	}

	return tables, nil
}
