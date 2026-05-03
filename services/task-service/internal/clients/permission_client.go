package clients

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type PermissionClient struct {
	baseURL string
}

func NewPermissionClient(url string) *PermissionClient {
	return &PermissionClient{baseURL: url}
}

func (c *PermissionClient) Create(userID, taskID string) error {
	body := map[string]string{
		"user_id": userID,
		"task_id": taskID,
	}

	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post(
		c.baseURL+"/permissions",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return err
	}

	return nil
}

func (c *PermissionClient) Check(userID, taskID string) (bool, error) {
	resp, err := http.Get(
		c.baseURL + "/permissions/check?user_id=" + userID + "&task_id=" + taskID,
	)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result struct {
		Allowed bool `json:"allowed"`
	}

	json.NewDecoder(resp.Body).Decode(&result)

	return result.Allowed, nil
}
