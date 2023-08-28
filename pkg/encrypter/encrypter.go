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
package encrypter

import (
	"encoding/base64"
	"io"

	crypto_rand "crypto/rand"

	"github.com/arush-sal/gese/pkg/util"
	"golang.org/x/crypto/nacl/box"
)

func Encrypt(recipientPublicKey string, secretValue string) (string, error) {
	recipientPublicKeyBase64, err := base64.StdEncoding.DecodeString(recipientPublicKey)
	util.ErrorCheck(err)

	var recipientPublicKeyBase64Array [32]byte
	copy(recipientPublicKeyBase64Array[:], recipientPublicKeyBase64)

	// You must use a different nonce for each message you encrypt with the
	// same key. Since the nonce here is 192 bits long, a random value
	// provides a sufficiently small probability of repeats.
	var nonce [24]byte
	if _, err := io.ReadFull(crypto_rand.Reader, nonce[:]); err != nil {
		panic(err)
	}

	// This encrypts msg and appends the result to the nonce.
	encrypted, err := box.SealAnonymous(nonce[:], []byte(secretValue), &recipientPublicKeyBase64Array, crypto_rand.Reader)
	util.ErrorCheck(err)

	encryptedBase64 := base64.StdEncoding.EncodeToString(encrypted)
	return encryptedBase64, nil
}
