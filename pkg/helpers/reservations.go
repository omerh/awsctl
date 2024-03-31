package helpers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// type riSummary struct {
// 	instanceType  string
// 	instanceCount int
// }

// GetAllReservations retrieve all reservations
func GetAllReservations(region string, state string) map[string]int64 {
	awsSession, _ := InitAwsSession(region)
	svc := ec2.New(awsSession)
	filter := []*ec2.Filter{
		{
			Name: aws.String("state"),
			Values: []*string{
				aws.String(state),
			},
		},
	}

	input := &ec2.DescribeReservedInstancesInput{
		Filters: filter,
	}

	result, _ := svc.DescribeReservedInstances(input)
	summaryResult := summariesRI(result)
	return summaryResult
}

func summariesRI(ri *ec2.DescribeReservedInstancesOutput) map[string]int64 {
	summary := make(map[string]int64, len(ri.ReservedInstances))
	for _, r := range ri.ReservedInstances {
		i := summary[*r.InstanceType]
		summary[*r.InstanceType] = i + *r.InstanceCount
	}
	return summary
}
