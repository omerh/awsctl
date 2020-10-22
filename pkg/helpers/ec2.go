package helpers

import (
	"fmt"
	"log"
	"time"

	"github.com/omerh/awsctl/pkg/outputs"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Ec2Event struct
type Ec2Event struct {
	instanceID       string
	eventStatusCode  string
	eventDescription string
	eventDueDate     time.Time
}

// GetAllEc2Events will check for all events
//
func GetAllEc2Events(region string, out string) {
	input := &ec2.DescribeInstanceStatusInput{}
	config := aws.NewConfig().WithRegion(region)
	sess, err := session.NewSession(config)

	if err != nil {
		fmt.Println(err)
	}

	svc := ec2.New(sess, config)
	result, _ := svc.DescribeInstanceStatus(input)

	var ec2evnets []Ec2Event

	for _, instance := range result.InstanceStatuses {
		for _, event := range instance.Events {
			if event.Code != nil {
				ec2evnets = append(ec2evnets, Ec2Event{
					instanceID:       *instance.InstanceId,
					eventStatusCode:  *event.Code,
					eventDescription: *event.Description,
					eventDueDate:     *event.NotBefore,
				})
			}
		}
	}

	printEc2Events(ec2evnets, out, region)
}

func printEc2Events(ec2evnets []Ec2Event, out string, region string) {
	switch out {
	case "json":
		// fmt.Println(ec2evnets)
		outputs.PrintGenericJSONOutput(ec2evnets, region)
	default:
		log.Printf("Running on region: %v", region)
		if len(ec2evnets) > 0 {
			for _, ec2evnet := range ec2evnets {
				log.Printf("Instance %v has event %v", ec2evnet.instanceID, ec2evnet.eventStatusCode)
				log.Printf("Description: %v", ec2evnet.eventDescription)
				log.Printf("Handle until: %v", ec2evnet.eventDueDate)
			}
		} else {
			log.Println("None found for region")
		}
		log.Println("==============================================")
	}
}

// GetAllEC2Instances get all instances
func GetAllEC2Instances(region string, lifeCycle string, state string) []*ec2.Reservation {
	awsSession, _ := InitAwsSession(region)
	svc := ec2.New(awsSession)
	var filter []*ec2.Filter

	if state == "running" {
		filter = []*ec2.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
					aws.String("pending"),
				},
			},
		}
	} else {
		filter = []*ec2.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("stopped"),
				},
			},
		}
	}

	input := &ec2.DescribeInstancesInput{
		Filters: filter,
	}
	result, _ := svc.DescribeInstances(input)
	ec2Slice := result.Reservations

	for result.NextToken != nil {
		input := &ec2.DescribeInstancesInput{
			NextToken: result.NextToken,
		}
		result, _ = svc.DescribeInstances(input)
		for _, r := range result.Reservations {
			ec2Slice = append(ec2Slice, r)
		}
	}

	if lifeCycle == "all" {
		return ec2Slice
	}

	var filteredEC2Slice []*ec2.Reservation
	for _, i := range ec2Slice {
		if lifeCycle == "spot" {
			if i.Instances[0].InstanceLifecycle != nil {
				filteredEC2Slice = append(filteredEC2Slice, i)
			}
		} else {
			if i.Instances[0].InstanceLifecycle == nil {
				filteredEC2Slice = append(filteredEC2Slice, i)
			}
		}
	}

	return filteredEC2Slice
}

// SummeriesEC2Instances summarieses into map instances by type and count them
func SummeriesEC2Instances(ec2Slice []*ec2.Reservation) map[string]int64 {
	summary := make(map[string]int64, len(ec2Slice))
	for _, i := range ec2Slice {
		c := summary[*i.Instances[0].InstanceType]
		summary[*i.Instances[0].InstanceType] = c + 1
	}
	return summary
}
