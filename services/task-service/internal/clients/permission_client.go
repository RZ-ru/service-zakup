package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PermissionClient struct {
	baseURL string
}

func NewPermissionClient(url string) *PermissionClient {
	return &PermissionClient{baseURL: url}
}

func (c *PermissionClient) Create(ctx context.Context, userID, taskID string) error {
	return retry(3, 200*time.Millisecond, func() error {

		body := map[string]string{
			"task_id": taskID,
		}

		data, err := json.Marshal(body)
		if err != nil {
			return err
		}

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			c.baseURL+"/permissions",
			bytes.NewBuffer(data),
		)
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")

		authHeader := ctx.Value("auth_header").(string)
		req.Header.Set("Authorization", authHeader)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("unexpected status: %d", resp.StatusCode)
		}

		return nil
	})
}

func (c *PermissionClient) Check(ctx context.Context, taskID string) (bool, error) {

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseURL+"/permissions/check?task_id="+taskID,
		nil,
	)
	if err != nil {
		return false, err
	}

	// 🔥 ОБЯЗАТЕЛЬНО
	authHeader := ctx.Value("auth_header").(string)
	req.Header.Set("Authorization", authHeader)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("status: %d", resp.StatusCode)
	}

	var result struct {
		Allowed bool `json:"allowed"`
	}

	json.NewDecoder(resp.Body).Decode(&result)

	return result.Allowed, nil
}

func retry(attempts int, sleep time.Duration, fn func() error) error {
	var err error

	for i := 0; i < attempts; i++ {
		err = fn()
		if err == nil {
			return nil
		}

		time.Sleep(sleep)
	}

	return err
}
