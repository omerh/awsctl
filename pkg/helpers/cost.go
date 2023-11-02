package helpers

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/pricing"
)

type rawAmazonCloudWatch struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]string
	Terms           map[string]map[string]map[string]rawAmazonCloudWatchTerm
}

type rawAmazonCloudWatchTerm struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]string
	TermAttributes  map[string]string
}

// GetAwsServiceCost use to get a product code for getting price
//
func GetAwsServiceCost() *pricing.GetProductsOutput {
	log.Println("Cost")
	region := "us-east-1"
	awsSession, _ := InitAwsSession(region)
	svc := pricing.New(awsSession)

	pinput := &pricing.GetProductsInput{
		// Filters: []*pricing.Filter{
		// 	{
		// 		Field: aws.String("location"),
		// 		Type:  aws.String("TERM_MATCH"),
		// 		Value: aws.String("EU (Ireland)"),
		// 	},
		// },
		FormatVersion: aws.String("aws_v1"),
		ServiceCode:   aws.String("AmazonCloudwatch"),
		// ServiceCode: aws.String("AmazonEc2"),
		MaxResults: aws.Int64(1),
	}
	result, _ := svc.GetProducts(pinput)
	// fmt.Println(result.PriceList[0]["serviceCode"])
	// fmt.Println(result.PriceList[0]["terms"])
	for _, r := range result.PriceList {
		ts := make(map[string]map[string]map[string]string)

		fmt.Println(ts, r)
	}

	// fmt.Println("pricelist")
	// fmt.Println(r.PriceList)

	// input := &pricing.DescribeServicesInput{
	// 	FormatVersion: aws.String("aws_v1"),
	// 	MaxResults:    aws.Int64(1),
	// 	ServiceCode:   aws.String("AmazonCloudWatch"),
	// }

	// result, err := svc.DescribeServices(input)
	// if err != nil {
	// 	log.Println(err)
	// 	return true
	// }
	// for _, s := range result.Services {
	// 	fmt.Println(s)
	// }
	return result
}
