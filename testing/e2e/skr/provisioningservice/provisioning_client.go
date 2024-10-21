package provisioningservice

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	interval              = time.Minute
	provisioningTimeout   = 30 * time.Minute
	deprovisioningTimeout = 90 * time.Minute
	environmentType       = "kyma"
	serviceName           = "kymaruntime"
	accessTokenPath       = "/oauth/token"
	environmentsPath      = "/provisioning/v1/environments"
)

type ProvisioningConfig struct {
	URL          string
	ClientID     string
	ClientSecret string
	UAA_URL      string
	PlanName     string
	User         string
	InstanceName string
	Region       string
}

type ProvisioningClient struct {
	cfg         ProvisioningConfig
	logger      *slog.Logger
	cli         *http.Client
	ctx         context.Context
	accessToken string
}

func NewProvisioningClient(cfg ProvisioningConfig, logger *slog.Logger, ctx context.Context, timeoutSeconds time.Duration) *ProvisioningClient {
	cli := &http.Client{Timeout: timeoutSeconds * time.Second}
	return &ProvisioningClient{
		cfg:    cfg,
		logger: logger,
		cli:    cli,
		ctx:    ctx,
	}
}

func (p *ProvisioningClient) CreateEnvironment() (CreatedEnvironmentResponse, error) {
	requestBody := CreateEnvironment{
		EnvironmentType: environmentType,
		PlanName:        p.cfg.PlanName,
		ServiceName:     serviceName,
		User:            p.cfg.User,
		Parameters: EnvironmentParameters{
			Name:   p.cfg.InstanceName,
			Region: p.cfg.Region,
		},
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return CreatedEnvironmentResponse{}, fmt.Errorf("failed to marshal request body: %v", err)
	}

	resp, err := p.sendRequest(http.MethodPost, p.cfg.URL+environmentsPath, http.StatusAccepted, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return CreatedEnvironmentResponse{}, err
	}

	var environment CreatedEnvironmentResponse
	err = p.unmarshallResponse(resp, &environment)
	if err != nil {
		return CreatedEnvironmentResponse{}, err
	}

	return environment, nil
}

func (p *ProvisioningClient) GetEnvironment(environmentID string) (Environment, error) {
	resp, err := p.sendRequest(http.MethodGet, p.cfg.URL+environmentsPath+"/"+environmentID, http.StatusOK, nil)
	if err != nil {
		return Environment{}, err
	}

	var environment Environment
	err = p.unmarshallResponse(resp, &environment)
	if err != nil {
		return Environment{}, err
	}

	return environment, nil
}

func (p *ProvisioningClient) DeleteEnvironment(environmentID string) (Environment, error) {
	resp, err := p.sendRequest(http.MethodDelete, p.cfg.URL+environmentsPath+"/"+environmentID, http.StatusAccepted, nil)
	if err != nil {
		return Environment{}, err
	}

	var environment Environment
	err = p.unmarshallResponse(resp, &environment)
	if err != nil {
		return Environment{}, err
	}

	return environment, nil
}

func (p *ProvisioningClient) GetAccessToken() error {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	clientCredentials := fmt.Sprintf("%s:%s", p.cfg.ClientID, p.cfg.ClientSecret)
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(clientCredentials))

	req, err := http.NewRequest(http.MethodPost, p.cfg.UAA_URL+accessTokenPath, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+encodedCredentials)

	resp, err := p.cli.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var accessToken AccessToken
	err = p.unmarshallResponse(resp, &accessToken)
	if err != nil {
		return err
	}

	p.accessToken = accessToken.Token

	return nil
}

func (p *ProvisioningClient) AwaitEnvironmentCreated(environmentID string) error {
	err := wait.PollUntilContextTimeout(p.ctx, interval, provisioningTimeout, true, func(ctx context.Context) (bool, error) {
		environment, err := p.GetEnvironment(environmentID)
		if err != nil {
			p.logger.Warn(fmt.Sprintf("error getting environment: %v", err))
			return false, nil
		}
		p.logger.Info("Received environment state", "environmentID", environmentID, "state", environment.State)
		if environment.State == OK {
			return true, nil
		}
		return false, nil
	})

	if err != nil {
		return fmt.Errorf("failed to wait for environment creation: %v", err)
	}

	return nil
}

func (p *ProvisioningClient) AwaitEnvironmentDeleted(environmentID string) error {
	err := wait.PollUntilContextTimeout(p.ctx, interval, deprovisioningTimeout, true, func(ctx context.Context) (bool, error) {
		req, err := http.NewRequest(http.MethodGet, p.cfg.URL+environmentsPath+"/"+environmentID, nil)
		if err != nil {
			p.logger.Warn(fmt.Sprintf("failed to create HTTP request: %v", err))
			return false, nil
		}

		req.Header.Set("Authorization", "Bearer "+p.accessToken)

		resp, err := p.cli.Do(req)
		if err != nil {
			p.logger.Warn(fmt.Sprintf("failed to send request: %v", err))
			return false, nil
		}

		p.logger.Info("Received response with status code", "environmentID", environmentID, "status code", resp.StatusCode)
		if resp.StatusCode == http.StatusNotFound {
			return true, nil
		}

		return false, nil
	})

	if err != nil {
		return fmt.Errorf("failed to wait for environment deletion: %v", err)
	}

	return nil
}

func (p *ProvisioningClient) sendRequest(method string, url string, expectedStatus int, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.accessToken)
	if method == http.MethodPost || method == http.MethodPatch || method == http.MethodPut {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := p.cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	if resp.StatusCode != expectedStatus {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, nil
}

func (p *ProvisioningClient) unmarshallResponse(resp *http.Response, output any) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if err := json.Unmarshal(body, output); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return nil
}
