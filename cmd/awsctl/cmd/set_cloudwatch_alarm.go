package cmd

import (
	"log"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/omerh/awsctl/pkg/helpers"
	"github.com/spf13/cobra"
)

var cmdCloudwatchAlarm = &cobra.Command{
	Use:   "cloudwatchalarm [region]",
	Short: "Find and Set cloudwatch alarms to aws resources",
	Run: func(cmd *cobra.Command, args []string) {
		// Command params
		apply, _ := cmd.Flags().GetBool("yes")
		region, _ := cmd.Flags().GetString("region")
		resource, _ := cmd.Flags().GetString("resource")
		metric, _ := cmd.Flags().GetString("metric")
		arn, _ := cmd.Flags().GetString("arn")
		threshold, _ := cmd.Flags().GetFloat64("threshold")
		action, _ := cmd.Flags().GetString("action")

		var namespace, metricName string
		// Resource type
		if resource == "lambda" {
			namespace = "AWS/Lambda"
		} else {
			log.Println("Currently support only resource type lambda")
			return
		}

		// Metric type
		if metric == "errors" {
			metricName = "Errors"
			// } else if metric == "latency" {
			// metricName = "Latency"
		} else {
			log.Println("Currently supporting metrics: errors")
		}

		// Start function
		if region == "all" {
			awsRegions, _ := helpers.GetAllAwsRegions()
			for _, r := range awsRegions {
				// All alarms in cloudwatch
				allAlarms := helpers.ListCloudwatchAlarms(r)
				// Filtered alarms according to the namespace of the alarms
				alarms := filterCloudwatchAlarmsForNamespace(allAlarms, namespace)
				// All lambdas in a region
				lambdasSlice := helpers.GetAllLambdaInRegion(r, arn)
				checkOrCreateCloudwatchAlarm(alarms, lambdasSlice, metricName, namespace, threshold, action, r, apply)
			}
			return
		}

		if region == "" {
			region = helpers.GetDefaultAwsRegion()
		}

		// All alarms in cloudwatch
		allAlarms := helpers.ListCloudwatchAlarms(region)
		// Filtered alarms according to the namespace of the alarms
		alarms := filterCloudwatchAlarmsForNamespace(allAlarms, namespace)
		// All lambdas in a region
		lambdasSlice := helpers.GetAllLambdaInRegion(region, arn)
		checkOrCreateCloudwatchAlarm(alarms, lambdasSlice, metricName, namespace, threshold, action, region, apply)
	},
}

func init() {
	// Flags
	cmdCloudwatchAlarm.Flags().String("resource", "", "Type of resource to set the alarm")
	cmdCloudwatchAlarm.MarkFlagRequired("resource")
	cmdCloudwatchAlarm.Flags().String("metric", "", "Type for alarm to set")
	cmdCloudwatchAlarm.MarkFlagRequired("metric")
	cmdCloudwatchAlarm.Flags().String("arn", "", "Resource ARN")
	cmdCloudwatchAlarm.Flags().Float64("threshold", 0, "Threshold float")
	cmdCloudwatchAlarm.MarkFlagRequired("threshold")
	cmdCloudwatchAlarm.Flags().String("action", "", "The action to do once alert is on, like sns arn")
}

func checkIfLambdaAlarmExistsInCloudwatch(alarms []*cloudwatch.MetricAlarm, lambda *lambda.FunctionConfiguration, metricName string) bool {
	for _, a := range alarms {
		if *a.Dimensions[0].Value == *lambda.FunctionName && *a.MetricName == metricName {
			log.Printf("Lambda alarm for %v on metric %v already exists\n", *lambda.FunctionName, metricName)
			return true
		}
	}
	return false
}

func filterCloudwatchAlarmsForNamespace(allalarms []*cloudwatch.MetricAlarm, namespace string) (alarms []*cloudwatch.MetricAlarm) {
	for _, a := range allalarms {
		if *a.Namespace == namespace {
			alarms = append(alarms, a)
		}
	}
	return alarms
}

func checkOrCreateCloudwatchAlarm(alarms []*cloudwatch.MetricAlarm, lambdasSlice []*lambda.FunctionConfiguration, metricName string, namespace string, threshold float64, action string, region string, apply bool) {
	for _, l := range lambdasSlice {
		alarmExists := checkIfLambdaAlarmExistsInCloudwatch(alarms, l, metricName)
		if !alarmExists {
			// fmt.Printf("No alarm for %v, Creating one for %v\n", *l.FunctionName, metricName)
			if apply {
				helpers.CreateLambdaCloudwatchAlarm(region, *l.FunctionName, metricName, namespace, threshold, action)
				log.Printf("Alarm was created for %v on metric %v", *l.FunctionName, metricName)
			} else {
				log.Printf("Pass --yes to set monitoring alarm on lambda: %v", *l.FunctionName)
			}
		}
	}
}
