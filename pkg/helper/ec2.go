package helper

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
