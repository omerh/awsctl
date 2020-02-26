package helpers

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/lambda"
)

// GetAllLmbdasInRegion List all lambdas in a region
func GetAllLmbdasInRegion(region string) {
	awsSession, _ := InitAwsSession(region)
	svc := lambda.New(awsSession)
	input := &lambda.ListFunctionsInput{}
	result, _ := svc.ListFunctions(input)

	fmt.Println(result.Functions)

}
