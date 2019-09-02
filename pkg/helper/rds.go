package helper

import (
	"fmt"
	"log"

	"github.com/omerh/awsctl/pkg/outputs"

	"github.com/aws/aws-sdk-go/service/rds"
)

// GetAllRds from a region
func GetAllRds(region string, rdsType string, out string) {
	switch rdsType {
	case "instance":
		getAllRDSDBInstances(region, out)
	case "cluster":
		GetAllRdsDBClusters(region, out)
	default:
		fmt.Println("Error RdsDbType")
	}
}

func getAllRDSDBInstances(region string, out string) {
	awsSession, _ := InitAwsSession(region)
	svc := rds.New(awsSession)
	input := &rds.DescribeDBInstancesInput{
		// DBInstanceIdentifier: aws.String("smorgasbord"),
	}
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
				if rds.DBClusterIdentifier != nil {
					log.Printf("RDS is a part of a cluster %v", *rds.DBClusterIdentifier)
				}
				log.Printf("RDS Storage: %vGB of %v encryption set to %v", *rds.AllocatedStorage, *rds.StorageType, *rds.StorageEncrypted)
				if rds.Endpoint != nil {
					log.Printf("RDS Network: %v:%v Running in %v ", *rds.Endpoint.Address, *rds.Endpoint.Port, *rds.AvailabilityZone)
				}
				log.Printf("RDS Backup policy: retention %v days, Backup window: %v", *rds.BackupRetentionPeriod, *rds.PreferredBackupWindow)
				log.Printf("RDS Status: %v", *rds.DBInstanceStatus)
				log.Printf("RDS Engine: %v %v", *rds.Engine, *rds.EngineVersion)
				log.Println()
			}
		} else {
			log.Println("No RDS in region")
		}
		log.Println("==============================================")
	}
}

// GetAllRdsDBClusters get all rds db clusters
func GetAllRdsDBClusters(region string, out string) {
	awsSession, _ := InitAwsSession(region)
	svc := rds.New(awsSession)
	input := &rds.DescribeDBClustersInput{}
	rdsClusters, _ := svc.DescribeDBClusters(input)
	printRdsClusters(rdsClusters.DBClusters, region, out)
}

// GetRDSSnapshots get all snapshot for instance(s) or clusters
func GetRDSSnapshots(rdsInstance string, rdsType string, region string, out string) {
	awsSession, _ := InitAwsSession(region)
	svc := rds.New(awsSession)
	input := &rds.DescribeDBSnapshotsInput{}

	rdsSnapshots, _ := svc.DescribeDBSnapshots(input)
	fmt.Println(rdsSnapshots)
}

func printRdsClusters(rdsClusters []*rds.DBCluster, region string, out string) {
	switch out {
	case "json":
		outputs.PrintGenericJSONOutput(rdsClusters, region)
	default:
		log.Printf("Running on region: %v", region)
		if len(rdsClusters) > 0 {
			for _, rds := range rdsClusters {
				log.Printf("RDS Identifier: %v ", *rds.DBClusterIdentifier)
				log.Printf("RDS Storage: %vGB encryption set to %v", *rds.AllocatedStorage, *rds.StorageEncrypted)
				if rds.Endpoint != nil {
					log.Printf("RDS Writer Network: %v:%v", *rds.Endpoint, *rds.Port)
					log.Printf("RDS Reader Network: %v:%v", *rds.ReaderEndpoint, *rds.Port)
				}
				var azs []string
				for i := range rds.AvailabilityZones {
					azs = append(azs, *rds.AvailabilityZones[i])
				}
				log.Printf("RDS AZs: %v", azs)
				log.Printf("RDS Backup policy: retention %v days, Backup window: %v", *rds.BackupRetentionPeriod, *rds.PreferredBackupWindow)
				log.Printf("RDS Status: %v", *rds.Status)
				log.Printf("RDS Engine: %v %v", *rds.Engine, *rds.EngineVersion)
				log.Println()
			}
		} else {
			log.Println("No RDS Clusters in region")
		}
		log.Println("==============================================")
	}
}
