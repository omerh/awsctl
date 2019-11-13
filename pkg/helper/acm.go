package helper

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/acm"
)

type acmCertificate struct {
	certificateArn string
	domainName     string
}

// GetAcmCertificates retrives all certificates
func GetAcmCertificates(region string) []*acm.CertificateSummary {
	log.Printf("Running on region: %v", region)
	awsSession, _ := InitAwsSession(region)
	svc := acm.New(awsSession)
	input := &acm.ListCertificatesInput{}
	result, err := svc.ListCertificates(input)
	if err != nil {
		log.Println(err)
	}
	certificateSlice := result.CertificateSummaryList

	// iterate over NextToken to retrive all
	for result.NextToken != nil {
		input := &acm.ListCertificatesInput{
			NextToken: result.NextToken,
		}
		result, _ := svc.ListCertificates(input)
		for _, certificate := range result.CertificateSummaryList {
			certificateSlice = append(certificateSlice, certificate)
		}
	}

	return certificateSlice
}

// DescribeAcmCertificate to list certificate metadata
func DescribeAcmCertificate(region string, arnCertificate string) *acm.DescribeCertificateOutput {
	awsSession, _ := InitAwsSession(region)
	svc := acm.New(awsSession)
	input := &acm.DescribeCertificateInput{
		CertificateArn: aws.String(arnCertificate),
	}
	certificate, err := svc.DescribeCertificate(input)
	if err != nil {
		log.Println(err)
	}
	// today := time.Now()
	// log.Printf("Checking certificate %v for domain %v\n", *certificate.Certificate.CertificateArn, *certificate.Certificate.DomainName)

	// // Check if domain is expired
	// if certificate.Certificate.NotAfter.Sub(today) <= 0 {
	// 	log.Printf("Domain is expired please renew with %v!!!", *certificate.Certificate.DomainValidationOptions[0].ValidationMethod)
	// } else {
	// 	// Domain is valid
	// }
	// // Check if the certificate was validated
	// if *certificate.Certificate.DomainValidationOptions[0].ValidationStatus == "SUCCESS" {
	// 	fmt.Println(*certificate.Certificate.DomainValidationOptions[0].ValidationMethod)
	// 	// Check if certificate is being used by any amazon resource (cloudfront or loadbalancer)
	// 	if len(certificate.Certificate.InUseBy) > 0 {
	// 		diff := certificate.Certificate.NotAfter.Sub(today)
	// 		log.Printf("life span is %v\n", diff)
	// 		log.Printf("expiring at %v and is being used\n", certificate.Certificate.NotAfter)
	// 	}
	// }

	return certificate
}

// DeleteUnusedAcmCertificates deletes all unused certificate in a region
func DeleteUnusedAcmCertificates(region string, apply bool) {
	certificates := GetAcmCertificates(region)
	for _, cert := range certificates {
		certInfo := DescribeAcmCertificate(region, *cert.CertificateArn)
		if len(certInfo.Certificate.InUseBy) == 0 {
			if apply {
				awsSession, _ := InitAwsSession(region)
				svc := acm.New(awsSession)
				input := &acm.DeleteCertificateInput{
					CertificateArn: aws.String(*certInfo.Certificate.CertificateArn),
				}
				_, err := svc.DeleteCertificate(input)
				if err != nil {
					log.Println(err)
				}
				log.Printf("Deleted certificate %v for domain %v", *cert.CertificateArn, *cert.DomainName)
			} else {
				log.Printf("Certificate %v for %v is not being used and can be deleted, pass --yes to delete", *cert.CertificateArn, *cert.DomainName)
			}
		}
	}
	fmt.Println("=====================================================================================================")
}
