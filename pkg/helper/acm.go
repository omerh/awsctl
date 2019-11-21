package helper

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/acm"
	"golang.org/x/net/publicsuffix"
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
		result, _ = svc.ListCertificates(input)
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
	certificateInfo, err := svc.DescribeCertificate(input)
	if err != nil {
		log.Println(err)
	}
	return certificateInfo
}

// CheckCertificateStatus checking certificate status
func CheckCertificateStatus(certificate *acm.DescribeCertificateOutput, region string) {
	today := time.Now()
	if certificate.Certificate.NotAfter.Sub(today) < 0 {
		// certificate expired
		log.Printf("Certificate for %v (%v) expired at %v !!\n", *certificate.Certificate.DomainName, *certificate.Certificate.CertificateArn, *certificate.Certificate.NotAfter)
		analyzeCertificate(certificate, region)
		fmt.Println("=====================================================================================================")
	} else if certificate.Certificate.NotAfter.Sub(today).Hours() < 168 {
		// certificate is about to expired in less than 1 week
		log.Printf("Certificate for %v (%v) is about to expire at %v\n", *certificate.Certificate.DomainName, *certificate.Certificate.CertificateArn, *certificate.Certificate.NotAfter)
		analyzeCertificate(certificate, region)
		fmt.Println("=====================================================================================================")
	} else if certificate.Certificate.NotAfter.Sub(today).Hours() < 720 {
		// certificate is about to expired in less than 30 days
		log.Printf("Certificate for %v (%v) is about to expire at %v\n", *certificate.Certificate.DomainName, *certificate.Certificate.CertificateArn, *certificate.Certificate.NotAfter)
		fmt.Println("=====================================================================================================")
	}
}

func analyzeCertificate(certificate *acm.DescribeCertificateOutput, region string) {
	log.Printf("Analyzing certificate %v (%v)\n", *certificate.Certificate.DomainName, *certificate.Certificate.CertificateArn)
	// Check if domain is being used
	if len(certificate.Certificate.InUseBy) == 0 {
		log.Printf("Certificate %v is not being used and can be deleted, run awsctl delete certificate to delete it\n", *certificate.Certificate.CertificateArn)
		return
	}

	// Check validation methud
	if *certificate.Certificate.DomainValidationOptions[0].ValidationMethod == "EMAIL" {
		awsSession, _ := InitAwsSession(region)
		svc := acm.New(awsSession)
		input := &acm.ResendValidationEmailInput{
			CertificateArn:   aws.String(*certificate.Certificate.CertificateArn),
			Domain:           aws.String(*certificate.Certificate.DomainName),
			ValidationDomain: aws.String(*certificate.Certificate.DomainName),
		}
		_, err := svc.ResendValidationEmail(input)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Certificate %v was issued using email validation, please check the email(s) %v for renewal link", *certificate.Certificate.CertificateArn, &certificate.Certificate.DomainValidationOptions[0].ValidationEmails)
	} else {
		// Validation is DNS
		log.Printf("Validation method %v", *certificate.Certificate.DomainValidationOptions[0].ValidationMethod)
		url := strings.TrimSuffix(*certificate.Certificate.DomainValidationOptions[0].ResourceRecord.Name, ".")
		value := strings.TrimSuffix(*certificate.Certificate.DomainValidationOptions[0].ResourceRecord.Value, ".")
		domain, _ := publicsuffix.EffectiveTLDPlusOne(url)

		// Check if domain is configured in aws
		hostedZoneID, domainExists := GetDomainHostedZoneID(domain, region)
		var recordConfigured bool
		if domainExists {
			// Check if certificate DNS record exists in route53
			recordConfigured = CheckIfRecoredSetValueInRoute53(url, value, hostedZoneID, region)
		}

		// Check if dns recored resolves
		cname, err := net.LookupCNAME(url)
		if err != nil {
			if recordConfigured == true {
				log.Printf("AWS Certificate CNAME %v does not resolves, but configured in route53, check your registrar", url)
			} else {
				log.Printf("AWS Certificate CNAME %v does not resolves and setup ok in route53, verify that domain is in your posetion and NS are pointing to AWS route53", url)
			}
			log.Println(err)
		} else if cname == value {
			if recordConfigured == true {
				log.Printf("CNAME resolves ok and SSL recored setup correct, check aws console")
			} else {
				log.Println("CNAME resolves ok, domain is not in route53, check with domain registrar")
			}
		} else if cname != value {
			if recordConfigured == true {
				log.Printf("CNAME %v resolves to the wrong recored %v instead of %v", url, cname, value)
			} else {
				log.Printf("CNAME %v resolves to the wrong recored %v instead of %v, but not in route53. Fix record in your registrar", url, cname, value)
			}

		}
	}
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
