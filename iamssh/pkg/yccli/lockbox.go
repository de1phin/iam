package yccli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type LockboxSecret struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func ListLockboxSecrets(folderId string) ([]LockboxSecret, error) {
	secretsRaw, err := ycExecute("lockbox", "secret", "list", "--folder-id", folderId, "--format", "json")
	if err != nil {
		return nil, fmt.Errorf("failed to list lockbox secrets: %w", err)
	}

	secrets := []LockboxSecret{}
	err = json.Unmarshal(secretsRaw, &secrets)
	if err != nil {
		return nil, fmt.Errorf("failed to parse lockbox secrets: %w", err)
	}

	return secrets, nil
}

type LockboxSecretEntry struct {
	Key       string `json:"key"`
	TextValue string `json:"textValue"`
}

type LockboxSecretGetResponse struct {
	Entries   []LockboxSecretEntry `json:"entries"`
	VersionId string               `json:"versionId"`
}

func LockboxSecretGet(secretId string) (*LockboxSecretGetResponse, error) {
	token, err := GetIamToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get iam token: %w", err)
	}

	url := fmt.Sprintf("https://payload.lockbox.api.cloud.yandex.net/lockbox/v1/secrets/%s/payload", secretId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	token = strings.Trim(token, " \n\t")
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	result := &LockboxSecretGetResponse{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, fmt.Errorf("faield to unmarshal lockbox secret: %w", err)
	}

	return result, nil
}
