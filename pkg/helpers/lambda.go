package helpers

import (
	"github.com/aws/aws-sdk-go/service/lambda"
)

// GetAllLmbdasInRegion List all lambdas in a region
func GetAllLmbdasInRegion(region string) []*lambda.FunctionConfiguration {
	awsSession, _ := InitAwsSession(region)
	svc := lambda.New(awsSession)
	input := &lambda.ListFunctionsInput{}
	result, _ := svc.ListFunctions(input)

	lambdas := result.Functions

	// Iterate over NextMarker
	for result.NextMarker != nil {
		input = &lambda.ListFunctionsInput{
			Marker: result.NextMarker,
		}
		result, _ = svc.ListFunctions(input)
		for _, l := range result.Functions {
			lambdas = append(lambdas, l)
		}
	}

	return lambdas
}
