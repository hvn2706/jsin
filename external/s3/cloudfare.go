package s3

import (
	"context"
	"fmt"
	s3config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"io"
	"jsin/config"
	"jsin/logger"
	"net/url"
	"os"
)

type IClient interface {
}

type ClientImpl struct {
	s3Client *s3.Client
	cfg      config.S3Config
}

type resolverV2 struct {
	uri       string
	accountID string
	bucket    string
}

func (r *resolverV2) ResolveEndpoint(_ context.Context, _ s3.EndpointParameters) (
	smithyendpoints.Endpoint, error,
) {
	u, err := url.Parse(fmt.Sprintf("https://%s.%s/%s", r.uri, r.accountID, r.bucket))
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
		options.EndpointResolverV2 = &resolverV2{}
	})
	return &ClientImpl{
		s3Client: client,
		cfg:      cfg,
	}
}

func (c *ClientImpl) GetImage(ctx context.Context) (string, error) {
	getObjectOutput, err := c.s3Client.GetObject(
		ctx,
		&s3.GetObjectInput{
			Bucket: &c.cfg.Cloudflare.Bucket,
		})
	if err != nil {
		logger.Errorf("===== Get image from s3 failed: %+v", err.Error())
		return "", err
	}
	// save image
	outFile, err := os.Create("image.png")
	if err != nil {
		logger.Errorf("===== Create file failed: %+v", err.Error())
		return "", err
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, getObjectOutput.Body)
	if err != nil {
		logger.Errorf("===== Copy file failed: %+v", err.Error())
		return "", err
	}

	return "", nil
}
