package helpers

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// ListCloudwatchAlarms list all cloudwatch alarms
func ListCloudwatchAlarms(region string) []*cloudwatch.MetricAlarm {
	awsSession, _ := InitAwsSession(region)
	svc := cloudwatch.New(awsSession)
	input := &cloudwatch.DescribeAlarmsInput{}
	result, _ := svc.DescribeAlarms(input)
	cloudwatchAlarms := result.MetricAlarms

	for result.NextToken != nil {
		input = &cloudwatch.DescribeAlarmsInput{
			NextToken: result.NextToken,
		}
		result, _ = svc.DescribeAlarms(input)
		for _, c := range result.MetricAlarms {
			cloudwatchAlarms = append(cloudwatchAlarms, c)
		}
	}

	return cloudwatchAlarms
}

// CreateLambdaCloudwatchAlarm will create a new alarm for cloudwatch
func CreateLambdaCloudwatchAlarm(region string, functionName string, metricName string, namespace string, threshold float64, action string) bool {
	awsSession, _ := InitAwsSession(region)
	svc := cloudwatch.New(awsSession)
	input := &cloudwatch.PutMetricAlarmInput{
		AlarmName:          aws.String(fmt.Sprintf("%v on %v", metricName, functionName)),
		ComparisonOperator: aws.String(cloudwatch.ComparisonOperatorGreaterThanOrEqualToThreshold),
		EvaluationPeriods:  aws.Int64(1),
		MetricName:         aws.String(metricName),
		Namespace:          aws.String(namespace),
		Period:             aws.Int64(60),
		Statistic:          aws.String(cloudwatch.StatisticSum),
		Threshold:          aws.Float64(threshold),
		ActionsEnabled:     aws.Bool(true),
		AlarmDescription:   aws.String(fmt.Sprintf("%v on %v greater than %v", metricName, functionName, threshold)),

		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("FunctionName"),
				Value: aws.String(functionName),
			},
		},

		AlarmActions: []*string{
			aws.String(action),
		},
	}

	// Debug input
	// fmt.Println(input)

	_, err := svc.PutMetricAlarm(input)
	if err != nil {
		fmt.Println("Error", err)
		return false
	}

	return true
}
