package helpers

import (
	"fmt"
	"log"
	"time"

	"github.com/omerh/awsctl/pkg/outputs"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
)

// RdsSnapshotInfo rds snapshot struct
type RdsSnapshotInfo struct {
	dbIdentifier         string
	dbSnapshotIdentifier string
	snapshotType         string
	snapshotCreatedTime  time.Time
	storageEncrypted     bool
}

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
	var rdsFilteredInstances []*rds.DBInstance

	for _, rds := range rdsInstances.DBInstances {
		if *rds.Engine != "neptune" {
			rdsFilteredInstances = append(rdsFilteredInstances, rds)
		}

	}
	// printRdsInstances(rdsInstances.DBInstances, region, out)
	printRdsInstances(rdsFilteredInstances, region, out)
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
				log.Printf("RDS Storage: %vGB, Storage type: %v, Encryption set to %v", *rds.AllocatedStorage, *rds.StorageType, *rds.StorageEncrypted)
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
func GetRDSSnapshots(resourceName string, rdsType string, region string, out string) []RdsSnapshotInfo {
	var rdsSnapshotInfoSlice []RdsSnapshotInfo
	switch rdsType {
	case "instance":
		var response []*rds.DBSnapshot
		rdsSnapshotInfoSlice, response = getRdsInstanceSnapshot(resourceName, region, out)
		if out == "json" {
			outputs.PrintGenericJSONOutput(response, region)
		}
	case "cluster":
		var response []*rds.DBClusterSnapshot
		rdsSnapshotInfoSlice, response = getRdsDBClusterSnapshot(resourceName, region, out)
		if out == "json" {
			outputs.PrintGenericJSONOutput(response, region)
		}
	default:
		fmt.Println("Error RdsDbType")
	}
	return rdsSnapshotInfoSlice
}

func getRdsDBClusterSnapshot(resourceName string, region string, out string) ([]RdsSnapshotInfo, []*rds.DBClusterSnapshot) {
	awsSession, _ := InitAwsSession(region)
	svc := rds.New(awsSession)
	var input *rds.DescribeDBClusterSnapshotsInput
	if resourceName != "" {
		input = &rds.DescribeDBClusterSnapshotsInput{
			DBClusterIdentifier: aws.String(resourceName),
		}
	} else {
		input = &rds.DescribeDBClusterSnapshotsInput{}
	}

	rdsSnapshots, _ := svc.DescribeDBClusterSnapshots(input)
	var rdsSnapshotSlice []RdsSnapshotInfo

	for _, r := range rdsSnapshots.DBClusterSnapshots {
		rdsSnapshotSlice = append(rdsSnapshotSlice, RdsSnapshotInfo{
			dbIdentifier:         *r.DBClusterIdentifier,
			dbSnapshotIdentifier: *r.DBClusterSnapshotIdentifier,
			snapshotType:         *r.SnapshotType,
			snapshotCreatedTime:  *r.SnapshotCreateTime,
			storageEncrypted:     *r.StorageEncrypted,
		})
	}
	return rdsSnapshotSlice, rdsSnapshots.DBClusterSnapshots
}

// PrintRdsSnapshotInformation print the needed snapshot information
func PrintRdsSnapshotInformation(rdsSnapshotInformation []RdsSnapshotInfo, region string, out string) {
	log.Printf("Running on region: %v", region)
	if len(rdsSnapshotInformation) > 0 {
		for _, i := range rdsSnapshotInformation {
			log.Printf("RDS Identifier: %v", i.dbIdentifier)
			log.Printf("RDS Snapshot identifier: %v", i.dbSnapshotIdentifier)
			log.Printf("RDS Snapshot type: %v", i.snapshotType)
			log.Printf("RDS Snapshot created date: %v", i.snapshotCreatedTime)
			log.Printf("RDS Snapshot encryption set to %v", i.storageEncrypted)
			log.Println()
		}
	} else {
		log.Println("No RDS Snapshots in region")
	}
	log.Println("==============================================")
}

func getRdsInstanceSnapshot(resourceName string, region string, out string) ([]RdsSnapshotInfo, []*rds.DBSnapshot) {
	awsSession, _ := InitAwsSession(region)
	svc := rds.New(awsSession)
	var input *rds.DescribeDBSnapshotsInput
	if resourceName != "" {
		input = &rds.DescribeDBSnapshotsInput{
			DBInstanceIdentifier: aws.String(resourceName),
		}
	} else {
		input = &rds.DescribeDBSnapshotsInput{}
	}

	rdsSnapshots, _ := svc.DescribeDBSnapshots(input)
	var rdsSnapshotSlice []RdsSnapshotInfo

	for _, r := range rdsSnapshots.DBSnapshots {
		rdsSnapshotSlice = append(rdsSnapshotSlice, RdsSnapshotInfo{
			dbIdentifier:         *r.DBInstanceIdentifier,
			dbSnapshotIdentifier: *r.DBSnapshotIdentifier,
			snapshotType:         *r.SnapshotType,
			snapshotCreatedTime:  *r.SnapshotCreateTime,
			storageEncrypted:     *r.Encrypted,
		})
	}
	return rdsSnapshotSlice, rdsSnapshots.DBSnapshots
}

func printRdsClusters(rdsClusters []*rds.DBCluster, region string, out string) {
	switch out {
	case "json":
		outputs.PrintGenericJSONOutput(rdsClusters, region)
	default:
		log.Printf("Running on region: %v", region)
		if len(rdsClusters) > 0 {
			for _, rds := range rdsClusters {
				fmt.Println("In cluster print")
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

// DeleteRdsSnapshots for deleting snapshots by date
func DeleteRdsSnapshots(rdsSnapshots []RdsSnapshotInfo, older int, region string, apply bool, rdsType string, out string) {
	log.Printf("Running on region: %v", region)
	log.Printf("Deleting snapshots older than %v days", older)

	now := time.Now()
	deleteDate := now.AddDate(0, 0, -(older))

	for _, s := range rdsSnapshots {
		if s.snapshotCreatedTime.Before(deleteDate) {
			log.Printf("Deleting %v that was created at %v for db %v", s.dbSnapshotIdentifier, s.snapshotCreatedTime, s.dbIdentifier)
			if apply == true {
				deleteSnapshot(s.dbSnapshotIdentifier, region, rdsType)
			} else {
				log.Println("Add -y/--yes to confirm delete")
			}
		}
	}
	log.Println("==============================================")
}

func deleteSnapshot(dbSnapshotIdentifier string, region string, rdsType string) {
	awsSession, _ := InitAwsSession(region)
	svc := rds.New(awsSession)
	switch rdsType {
	case "instance":
		input := &rds.DeleteDBSnapshotInput{
			DBSnapshotIdentifier: aws.String(dbSnapshotIdentifier),
		}
		_, err := svc.DeleteDBSnapshot(input)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("Snapshot %v was deleted", dbSnapshotIdentifier)
		}
	case "cluster":
		input := &rds.DeleteDBClusterSnapshotInput{
			DBClusterSnapshotIdentifier: aws.String(dbSnapshotIdentifier),
		}
		_, err := svc.DeleteDBClusterSnapshot(input)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("Snapshot %v was deleted", dbSnapshotIdentifier)
		}
	}
}
