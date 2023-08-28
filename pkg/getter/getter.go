/*
Copyright © 2023 Arush Salil

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package getter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/arush-sal/gese/pkg/constant"
	"github.com/arush-sal/gese/pkg/types"
	"github.com/google/go-github/v54/github"
)

type Client struct {
	HTTPClient *http.Client
	Repo       *types.GHRepo
	Token      string
	Req        *http.Request
}

func (c *Client) GetRepoID() error {
	gr := new(github.Repository)
	err := c.GetHTTPRequest(http.MethodGet, fmt.Sprintf("%s/repos/%s/%s", constant.GITHUB_API_ENDPOINT, c.Repo.Owner, c.Repo.Name))
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(c.Req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, gr)
	if err != nil {
		return err
	}

	c.Repo.ID = strconv.Itoa(int(*gr.ID))

	return nil
}

func (c *Client) GetPubKeyEndpoint(env string) (string, error) {
	pubkey := new(types.EnvPubKey)

	err := c.GetHTTPRequest(http.MethodGet, fmt.Sprintf("%s/repositories/%s/environments/%s/secrets/public-key", constant.GITHUB_API_ENDPOINT, c.Repo.ID, env))
	if err != nil {
		return "", err
	}

	resp, err := c.HTTPClient.Do(c.Req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, pubkey)
	if err != nil {
		return "", err
	}

	fmt.Println(pubkey.KeyID)

	return pubkey.Key, err
}

func NewClient() (*Client, error) {
	c := new(Client)
	c.Repo = new(types.GHRepo)
	c.HTTPClient = http.DefaultClient
	return c, nil
}

func (c *Client) GetHTTPRequest(m, u string) error {
	r, err := http.NewRequest(m, u, nil)
	r.Header.Add("Accept", "application/vnd.github+json")
	r.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	c.Req = r
	return err
}
