package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyendpoints "github.com/aws/smithy-go/endpoints"

	"jsin/config"
	"jsin/logger"
)

type IClient interface {
	UploadObject(ctx context.Context, body io.Reader, objectKey string) error
	GetImage(ctx context.Context, objectKey string) ([]byte, error)
}

type ClientImpl struct {
	s3Client *s3.Client
	cfg      config.S3Config
}

var _ IClient = &ClientImpl{}

type resolverV2 struct {
	uri       string
	accountID string
	bucket    string
}

func (r *resolverV2) ResolveEndpoint(_ context.Context, _ s3.EndpointParameters) (
	smithyendpoints.Endpoint, error,
) {
	u, err := url.Parse(fmt.Sprintf("https://%s.%s/%s", r.accountID, r.uri, r.bucket))
	if err != nil {
		return smithyendpoints.Endpoint{}, err
	}
	return smithyendpoints.Endpoint{
		URI: *u,
	}, nil
}

func NewClient(cfg config.S3Config) *ClientImpl {
	configS3, err := s3config.LoadDefaultConfig(context.TODO(),
		s3config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.Cloudflare.AccessKeyID,
			cfg.Cloudflare.SecretAccessKey,
			"",
		)),
		s3config.WithRegion("auto"),
	)
	if err != nil {
		logger.Errorf("===== Load s3 config failed: %+v", err.Error())
	}

	client := s3.NewFromConfig(configS3, func(options *s3.Options) {
		options.EndpointResolverV2 = &resolverV2{
			uri:       cfg.Cloudflare.Uri,
			accountID: cfg.Cloudflare.AccountId,
			bucket:    cfg.Cloudflare.Bucket,
		}
	})
	return &ClientImpl{
		s3Client: client,
		cfg:      cfg,
	}
}

func (c *ClientImpl) UploadObject(ctx context.Context, body io.Reader, objectKey string) error {
	response, err := c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &c.cfg.Cloudflare.Bucket,
		Key:         &objectKey,
		Body:        body,
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		logger.Errorf("===== Upload image to s3 failed: %+v", err.Error())
		return err
	}
	logger.Infof("===== Upload image to s3 success: %+v", response)
	return nil
}

func (c *ClientImpl) GetImage(ctx context.Context, objectKey string) ([]byte, error) {
	getObjectOutput, err := c.s3Client.GetObject(
		ctx,
		&s3.GetObjectInput{
			Bucket: &c.cfg.Cloudflare.Bucket,
			Key:    &objectKey,
		})
	if err != nil {
		logger.Errorf("===== Get image from s3 failed: %+v", err.Error())
		return nil, err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(getObjectOutput.Body)
	if err != nil {
		logger.Errorf("===== Read image from s3 failed: %+v", err.Error())
		return nil, err
	}

	return buf.Bytes(), nil
}
