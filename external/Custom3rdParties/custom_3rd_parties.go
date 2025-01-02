package Custom3rdParties

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"jsin/config"
	"jsin/logger"
	"jsin/pkg/common"
	error2 "jsin/pkg/common/error"
)

type IClient interface {
	GetRandomImageFrom3rdParties(ctx context.Context) ([]byte, error)
}

type ClientImpl struct {
	cfg    config.Custom3rdPartiesConfig
	client *http.Client
}

var _ IClient = &ClientImpl{}

func NewClient(cfg config.Custom3rdPartiesConfig) *ClientImpl {
	return &ClientImpl{
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (c *ClientImpl) GetRandomImageFrom3rdParties(ctx context.Context) ([]byte, error) {
	imgURL, err := c.getRandomImageURL(ctx)
	if err != nil {
		return nil, err
	}

	img, err := c.getImageFromURL(ctx, imgURL)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (c *ClientImpl) getRandomImageURL(ctx context.Context) (string, error) {
	partyID, err := rand.Int(rand.Reader, big.NewInt(int64(len(c.cfg.Parties))))
	if err != nil {
		logger.Errorf("Failed to select a random party: %v", err)
		return "", err
	}

	party := c.cfg.Parties[partyID.Int64()]
	req, err := http.NewRequest(party.Method, party.URL, nil)
	if err != nil {
		logger.Errorf("Failed to create request: %v", err)
		return "", err
	}

	req.Header.Set(party.Header, party.Token)

	resp, err := c.client.Do(req)
	if err != nil {
		logger.Errorf("Failed to send request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Error reading response: %v", err)
		return "", err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		logger.Errorf("Error parsing JSON: %v", err)
		return "", err
	}

	imgURL, ok := common.FindKeyInMap(result, party.JSONKey)
	if !ok {
		return "", error2.ErrNotFound("no image found")
	}

	return fmt.Sprintf("%v", imgURL), nil
}

func (c *ClientImpl) getImageFromURL(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Errorf("Failed to create request: %v", err)
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Error reading response: %v", err)
		return nil, err
	}

	return body, nil
}
