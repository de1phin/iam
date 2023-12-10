package yc

import (
	"encoding/json"
	"io"
	"net/http"
)

type TokenGetter interface {
	GetToken() (string, error)
}

type yccli struct{}

func YcCli() TokenGetter {
	return &yccli{}
}

func (*yccli) GetToken() (string, error) {
	out, err := ycExecute("iam", "create-token")
	return string(out), err
}

type computeMetadata struct{}

func ComputeMetadata() TokenGetter {
	return &computeMetadata{}
}

type metadataTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (*computeMetadata) GetToken() (string, error) {
	url := "http://169.254.169.254/computeMetadata/v1/instance/service-accounts/default/token"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Metadata-Flavor", "Google")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	token := &metadataTokenResponse{}
	err = json.Unmarshal(body, token)
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}
