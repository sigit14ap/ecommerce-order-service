package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sigit14ap/order-service/internal/delivery/dto"
)

type UserService struct {
	BaseURL      string
	ServiceToken string
}

func NewUserService(baseURL, token string) *UserService {
	return &UserService{
		BaseURL:      baseURL,
		ServiceToken: token,
	}
}

func (c *UserService) CallUserService(method, endpoint string, authToken string, payload interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, endpoint)

	var jsonData []byte
	var err error
	if payload != nil {
		jsonData, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authToken)
	req.Header.Set("X-Service-Token", c.ServiceToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call shop service: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return resp, nil
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("unauthorized access")
	case http.StatusNotFound:
		return nil, fmt.Errorf("resource not found")
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("internal server error")
	default:
		return nil, fmt.Errorf("unexpected status code %d received from shop service", resp.StatusCode)
	}
}

func (c *UserService) UserDetail(context *gin.Context) (*dto.UserDetailResponse, error) {
	token := context.GetHeader("Authorization")
	if token == "" {
		return nil, errors.New("authorization required")
	}

	url := "api/v1/users/me"
	response, err := c.CallUserService("GET", url, token, nil)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var apiResponse dto.ApiResponse
	if err := json.NewDecoder(response.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	userData, ok := apiResponse.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format")
	}

	productDetail, err := json.Marshal(userData["user"])
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user data: %w", err)
	}

	var user dto.UserDetailResponse
	if err := json.Unmarshal(productDetail, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	return &user, nil
}
