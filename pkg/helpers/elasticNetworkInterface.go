package helpers

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// GetAllElasticNetworkInterfaces get all NICs
func GetAllElasticNetworkInterfaces(region string, filter string) []*ec2.NetworkInterface {
	config := aws.NewConfig().WithRegion(region)
	sess, err := session.NewSession(config)
	if err != nil {
		log.Printf("Error creating session\n%v", err)
	}
	svc := ec2.New(sess, config)
	networkInterfaceFilter := []*ec2.Filter{
		{
			Name: aws.String("status"),
			Values: []*string{
				aws.String(filter),
			},
		},
	}
	input := &ec2.DescribeNetworkInterfacesInput{Filters: networkInterfaceFilter}

	result, _ := svc.DescribeNetworkInterfaces(input)

	networkInterfacesSlice := result.NetworkInterfaces

	for result.NextToken != nil {
		input = &ec2.DescribeNetworkInterfacesInput{
			NextToken: result.NextToken,
			Filters:   networkInterfaceFilter,
		}
		result, _ = svc.DescribeNetworkInterfaces(input)
		for _, n := range result.NetworkInterfaces {
			networkInterfacesSlice = append(networkInterfacesSlice, n)
		}
	}

	return networkInterfacesSlice
}

// DeleteNetworkInterface delete by id
func DeleteNetworkInterface(region string, networkInterfaceID string, apply bool) bool {
	config := aws.NewConfig().WithRegion(region)
	sess, err := session.NewSession(config)
	if err != nil {
		log.Printf("Error creating session\n%v", err)
	}
	if !apply {
		log.Printf("awsctl would delete %v pass --yes to delete", networkInterfaceID)
	} else {
		svc := ec2.New(sess, config)
		input := &ec2.DeleteNetworkInterfaceInput{
			NetworkInterfaceId: &networkInterfaceID,
		}
		_, err = svc.DeleteNetworkInterface(input)
		if err != nil {
			log.Printf("There was a problme deleting network interface\n%v", err)
			return false
		}
	}
	return true
}
