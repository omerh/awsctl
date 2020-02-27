package helpers

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/lambda"
)

// GetAllLmbdasInRegion List all lambdas in a region
func GetAllLmbdasInRegion(region string, arn string) []*lambda.FunctionConfiguration {
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

	// Filter out lambda
	if len(arn) > 0 {
		for _, l := range lambdas {
			if *l.FunctionArn == arn {
				fmt.Println("match")
				lambdas = nil
				lambdas = append(lambdas, l)
				return lambdas
			}
		}
	}

	return lambdas
}
