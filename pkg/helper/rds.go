package helper

import (
	"log"

	"github.com/omerh/awsctl/pkg/outputs"

	"github.com/aws/aws-sdk-go/service/rds"
)

// GetAllRdsInstances from a region
func GetAllRdsInstances(region string, out string) {
	awsSession, _ := InitAwsSession(region)
	svc := rds.New(awsSession)
	input := &rds.DescribeDBInstancesInput{}
	rdsInstances, _ := svc.DescribeDBInstances(input)

	printRdsInstances(rdsInstances.DBInstances, region, out)
}

func printRdsInstances(rdsInstances []*rds.DBInstance, region string, out string) {
	switch out {
	case "json":
		outputs.PrintGenericJSONOutput(rdsInstances, region)
	default:
		log.Printf("Running on region: %v", region)
		if len(rdsInstances) > 0 {
			for _, rds := range rdsInstances {
				log.Printf("RDS Identifier: %v, Type: %v ", *rds.DBInstanceIdentifier, *rds.DBInstanceClass)
				log.Printf("RDS Storage: %vGB of %v encryption set to %v", *rds.AllocatedStorage, *rds.StorageType, *rds.StorageEncrypted)
				log.Printf("RDS Network: %v:%v Running in %v ", *rds.Endpoint.Address, *rds.Endpoint.Port, *rds.AvailabilityZone)
				log.Printf("RDS Backup policy: retention %v days, Backup window: %v", *rds.BackupRetentionPeriod, *rds.PreferredBackupWindow)
				log.Printf("RDS Status: %v", *rds.DBInstanceStatus)
				log.Println()
			}
		} else {
			log.Println("No RDS in region")
		}
		log.Println("==============================================")
	}
}
