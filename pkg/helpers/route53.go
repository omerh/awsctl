package helpers

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
)

// GetDomainHostedZoneID function get domain name and return its hosted zone id and if exists
func GetDomainHostedZoneID(domain string, region string) (domainHostedZoneID string, exists bool) {
	log.Printf("Probing route53 on %v", domain)
	awsSession, _ := InitAwsSession(region)
	svc := route53.New(awsSession)
	input := &route53.ListHostedZonesByNameInput{
		DNSName:  aws.String(domain),
		MaxItems: aws.String("1"),
	}
	result, _ := svc.ListHostedZonesByName(input)

	if strings.HasPrefix(*result.HostedZones[0].Name, domain) {
		log.Println("Domain exists in AWS")
		exists = true
	} else {
		log.Println("Domain does not in AWS")
		exists = false
	}

	return *result.HostedZones[0].Id, exists
}

// CheckIfRecoredSetValueInRoute53 Checks if a recored exists in Route53
func CheckIfRecoredSetValueInRoute53(record string, value string, hostedZoneID string, region string) bool {
	log.Printf("Probing DNS recored %v in Route53", record)
	awsSession, _ := InitAwsSession(region)
	svc := route53.New(awsSession)
	input := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(hostedZoneID),
		MaxItems:        aws.String("1"),
		StartRecordName: aws.String(record),
	}
	result, _ := svc.ListResourceRecordSets(input)
	if strings.HasPrefix(*result.ResourceRecordSets[0].ResourceRecords[0].Value, value) {
		log.Println("Recored set configured correctly in route53")
		return true
	}

	log.Println("Recored set configured correctly in route53")
	return false
}
