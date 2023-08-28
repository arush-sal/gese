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
package runner

import (
	"github.com/arush-sal/gese/pkg/encrypter"
	"github.com/arush-sal/gese/pkg/getter"
)

func GetEncryptedValue(token, repo, owner, secret, env string) (string, error) {
	c, err := getter.NewClient()
	if err != nil {
		return "", err
	}

	c.Token = token
	c.Repo.Name = repo
	c.Repo.Owner = owner

	err = c.GetRepoID()
	if err != nil {
		return "", err
	}

	pubkey, err := c.GetPubKeyEndpoint(env)
	if err != nil {
		return "", err
	}

	return encrypter.Encrypt(pubkey, secret)

}
