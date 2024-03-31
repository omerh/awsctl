package helpers

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// GetAllS3Buckets retrieves all the buckets in a region
func GetAllS3Buckets() []*s3.Bucket {
	awsSession, _ := InitAwsSession("us-east-1")
	svc := s3.New(awsSession)
	input := &s3.ListBucketsInput{}
	result, _ := svc.ListBuckets(input)
	return result.Buckets
}

// CheckBucketEncryption check if bucket encryption is set
func CheckBucketEncryption(bucket string, region string) bool {
	awsSession, _ := InitAwsSession(region)
	svc := s3.New(awsSession)
	input := &s3.GetBucketEncryptionInput{
		Bucket: aws.String(bucket),
	}

	_, err := svc.GetBucketEncryption(input)
	if err != nil {
		// The server side encryption configuration was not found
		log.Printf("Error getting bucket encryption\n%v\n", err)
		// return false
	}
	// fmt.Println(result)
	return err == nil

}

// GetS3BucketLocation get a bucket region
func GetS3BucketLocation(bucket string) string {
	defaultRegion := "us-east-1"
	awsSession, _ := InitAwsSession(defaultRegion)
	svc := s3.New(awsSession)
	input := &s3.GetBucketLocationInput{
		Bucket: aws.String(bucket),
	}
	// log.Printf("Checking bucket %v for location", bucket)
	result, err := svc.GetBucketLocation(input)
	if err != err {
		log.Printf("Error getting bucket location\n%v\n", err)
	}
	if result.LocationConstraint == nil {
		return defaultRegion
	}
	return *result.LocationConstraint
}

// GetS3PublicAccess get public access to s3 bucket
func GetS3PublicAccess(bucket string, region string) {
	awsSession, _ := InitAwsSession(region)
	svc := s3.New(awsSession)
	input := &s3.GetBucketAclInput{
		Bucket: aws.String(bucket),
	}
	result, _ := svc.GetBucketAcl(input)
	fmt.Printf("Checking bucket %v\n", bucket)
	for _, g := range result.Grants {
		// fmt.Println(*g.Grantee.ID)
		// fmt.Println(*result.Owner.ID)
		// fmt.Println(result)
		if g.Grantee.ID == nil {
			// its a service grants
			fmt.Printf("Service %v has %v grants\n", *g.Grantee.URI, *g.Permission)
		} else if *g.Grantee.ID != *result.Owner.ID {
			fmt.Printf("bucket %v has extra grants\n", bucket)
			// fmt.Println(g)
		}
	}
}
