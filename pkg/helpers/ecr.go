package helpers

import (
	"fmt"
	"log"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// ImageArr array of image details
type ImageArr []*ecr.ImageDetail

func (s ImageArr) Len() int { return len(s) }
func (s ImageArr) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ImageArr) Less(i, j int) bool {
	t1 := s[i].ImagePushedAt.Unix()
	t2 := s[j].ImagePushedAt.Unix()
	return t1 > t2
}

// GetECRRepositories from aws region
func GetECRRepositories(region string) []*ecr.Repository {
	log.Printf("Running on region: %v", region)
	awsSession, _ := InitAwsSession(region)
	svc := ecr.New(awsSession)
	input := &ecr.DescribeRepositoriesInput{}
	result, _ := svc.DescribeRepositories(input)

	registriesSlice := result.Repositories

	// iterate over NextToken to retrive all repositories from Ecr in the region
	for result.NextToken != nil {
		input := &ecr.DescribeRepositoriesInput{
			NextToken: result.NextToken,
		}
		result, _ = svc.DescribeRepositories(input)
		registriesSlice = append(registriesSlice, result.Repositories...)
	}
	return registriesSlice
}

// CheckECRRepositoryLifecyclePolicy for a repository in a region
func CheckECRRepositoryLifecyclePolicy(repositoryName string, region string) bool {
	awsSession, _ := InitAwsSession(region)
	svc := ecr.New(awsSession)
	input := &ecr.GetLifecyclePolicyInput{
		RepositoryName: aws.String(repositoryName),
	}
	_, err := svc.GetLifecyclePolicy(input)
	if err != nil {
		// Ecr lifecyclePolicy is not set
		log.Printf("error: %v", err)
	}
	return err == nil
}

// SetEcrRepositoryLifecyclePolicy set the life time policy
func SetEcrRepositoryLifecyclePolicy(repositoryName string, days int, region string) {
	log.Printf("Setting lifecycle policy to %v for %v days", repositoryName, days)
	awsSession, _ := InitAwsSession(region)
	svc := ecr.New(awsSession)

	// tagStatus is hardcoded to prevent mistakes like deleting latest that is not maintained
	lifecyclePolicyText := fmt.Sprintf(`{
		"rules": [
			{
				"rulePriority": 1,
				"description": "Expire untagged images by awsctl",
				"selection": {
					"tagStatus": "untagged",
					"countType": "sinceImagePushed",
					"countUnit": "days",
					"countNumber": %v
				},
				"action": {
					"type": "expire"
				}
			}
		]
	  }`, days)

	input := &ecr.PutLifecyclePolicyInput{
		RepositoryName:      aws.String(repositoryName),
		LifecyclePolicyText: aws.String(lifecyclePolicyText),
	}

	_, err := svc.PutLifecyclePolicy(input)

	if err != nil {
		log.Printf("There was a problem setting lifecycle policy to %v", repositoryName)
		log.Println(err)
	}
}

// SetEcrRepoImageScanOnPush set image scan on push configuration on ecr repository
func SetEcrRepoImageScanOnPush(repositoryName string, region string, scanOnPush bool) {
	awsSession, _ := InitAwsSession(region)
	svc := ecr.New(awsSession)
	input := &ecr.PutImageScanningConfigurationInput{
		RepositoryName: aws.String(repositoryName),
		ImageScanningConfiguration: &ecr.ImageScanningConfiguration{
			ScanOnPush: aws.Bool(scanOnPush),
		},
	}
	_, err := svc.PutImageScanningConfiguration(input)
	if err != nil {
		log.Printf("There was a problem setting repository %v to ScanOnPush %v", repositoryName, scanOnPush)
	} else {
		log.Printf("scanOnPush for repository %v was set to %v", repositoryName, scanOnPush)
	}
}

// EcrDescribeImages get details on images
func EcrDescribeImages(repositoryName string, region string, keep int) ([]*ecr.ImageDetail, int) {
	awsSession, _ := InitAwsSession(region)
	svc := ecr.New(awsSession)
	input := &ecr.DescribeImagesInput{
		RepositoryName: aws.String(repositoryName),
	}
	result, _ := svc.DescribeImages(input)

	imagesDetailsSlice := result.ImageDetails

	// iterate over NextToken to retrieve all repositories from Ecr in the region
	for result.NextToken != nil {
		input := &ecr.DescribeImagesInput{
			RepositoryName: aws.String(repositoryName),
			NextToken:      result.NextToken,
		}
		result, _ = svc.DescribeImages(input)
		imagesDetailsSlice = append(imagesDetailsSlice, result.ImageDetails...)
	}
	sortedImagesToDelete := sortEcrRepos(imagesDetailsSlice, keep)
	return sortedImagesToDelete, len(imagesDetailsSlice) - len(sortedImagesToDelete)
}

func sortEcrRepos(imagesDetail []*ecr.ImageDetail, keep int) []*ecr.ImageDetail {
	var imageArray ImageArr = imagesDetail
	sort.Stable(imageArray)
	if len(imageArray) > keep {
		imageArray = imageArray[keep:]
	}
	return imageArray
}

// DeleteEcrImages delete according to image digest
func DeleteEcrImages(repo string, digest []string, region string, apply bool) {
	images := make([]*ecr.ImageIdentifier, len(digest))

	for k, v := range digest {
		t := &ecr.ImageIdentifier{}
		t.ImageDigest = aws.String(v)
		images[k] = t
	}

	bulk := 100

	for i := 0; i < len(images); i += bulk {
		batch := images[i:min(i+bulk, len(images))]
		deleteImagesBatch(repo, region, batch, apply)
	}
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func deleteImagesBatch(repo string, region string, images []*ecr.ImageIdentifier, apply bool) {
	awsSession, _ := InitAwsSession(region)
	svc := ecr.New(awsSession)
	input := &ecr.BatchDeleteImageInput{
		RepositoryName: aws.String(repo),
		ImageIds:       images,
	}

	if apply {
		log.Printf("Repo: %v, deleting bulk of %v images\n", repo, len(images))
		_, err := svc.BatchDeleteImage(input)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Printf("Repo: %v, should delete %v images, pass --yes to apply\n", repo, len(images))
	}
}
