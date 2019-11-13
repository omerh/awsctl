package helper

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/acm"
)

// GetAcmCertificates retrives all certificates
func GetAcmCertificates(region string) {
	log.Printf("Running on region: %v", region)
	awsSession, _ := InitAwsSession(region)
	svc := acm.New(awsSession)
	input := &acm.ListCertificatesInput{}
	result, err := svc.ListCertificates(input)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(result)
}
