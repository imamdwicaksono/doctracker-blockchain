package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	ctx    context.Context
	client *s3.Client
	bucket string
}

var S3 *S3Storage

// InitializeS3Storage creates an S3Storage instance for Storj
func InitializeS3Storage(ctx context.Context) *S3Storage {
	endpoint := os.Getenv("STORJ_S3_ENDPOINT")
	bucket := os.Getenv("STORJ_S3_BUCKET")
	accessKey := os.Getenv("STORJ_ACCESS_KEY")
	secretKey := os.Getenv("STORJ_SECRET_KEY")

	customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           endpoint,
			SigningRegion: "us-east-1",
		}, nil
	})

	fmt.Printf("S3 Endpoint: %s\n", endpoint)
	fmt.Printf("S3 Bucket: %s\n", bucket)
	fmt.Printf("S3 Access Key: %s\n", accessKey)

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithEndpointResolver(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	// ✅ Tambahkan o.UsePathStyle = true di sini
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &S3Storage{
		ctx:    ctx,
		client: client,
		bucket: bucket,
	}
}

func (s *S3Storage) UploadS3File(fileName string, fileContent []byte) error {
	_, err := s.client.PutObject(s.ctx, &s3.PutObjectInput{
		Bucket: &s.bucket,
		Key:    &fileName,
		Body:   bytes.NewReader(fileContent),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *S3Storage) DownloadS3File(fileName string) ([]byte, error) {
	resp, err := s.client.GetObject(s.ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &fileName,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
