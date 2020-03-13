// Copyright (c) 2019 Baidu, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jwa

import "github.com/baidu/bfe/bfe_modules/mod_auth_jwt/jwk"

// algorithm used to calculate signature for JWS
type JWSAlg interface {
	Update(msg []byte) (n int, err error) // update msg
	Sign() (sig []byte)                   // get signature
	Verify(sig []byte) bool               // verify signature
}

// algorithm used to encrypt & decrypt CEK(Content Encryption Key) fro JWE
type JWEAlg interface {
	//Encrypt(msg []byte) (cipher []byte) // implement this if needed
	Decrypt(eCek []byte) (cek []byte, err error)
}

// algorithm used to encrypt & decrypt content for JWE
type JWEEnc interface {
	//Encrypt(msg []byte) (cipher []byte) // implement this if needed
	Decrypt(iv, aad, cipherText, tag []byte) (msg []byte, err error)
}

type jwsAlgFactory func(*jwk.JWK) (JWSAlg, error)
type jweEncFactory func(cek []byte) (JWEEnc, error)
type jweAlgFactory func(*jwk.JWK, map[string]interface{}) (JWEAlg, error)

// exported algorithms
var (
	JWSAlgSet = map[string]jwsAlgFactory{
		"HS256": NewHS256,
		"HS384": NewHS384,
		"HS512": NewHS512,
		"RS256": NewRS256,
		"RS384": NewRS384,
		"RS512": NewRS512,
		"ES256": NewES256,
		"ES384": NewES384,
		"ES512": NewES512,
		"PS256": NewPS256,
		"PS384": NewPS384,
		"PS512": NewPS512,
	}

	JWEEncSet = map[string]jweEncFactory{
		"A128CBC-HS256": NewA128CBCHS256,
		"A192CBC-HS384": NewA192CBCHS384,
		"A256CBC-HS512": NewA256CBCHS512,
		"A128GCM":       NewA128GCM,
		"A192GCM":       NewA192GCM,
		"A256GCM":       NewA256GCM,
	}

	JWEEncKeyLength = map[string]int{
		"A128CBC-HS256": 128,
		"A192CBC-HS384": 192,
		"A256CBC-HS512": 256,
		"A128GCM":       128,
		"A192GCM":       192,
		"A256GCM":       256,
	}

	JWEAlgSet = map[string]jweAlgFactory{
		"dir":                NewDir,
		"RSA1_5":             NewRSA15,
		"RSA-OAEP":           NewRSAOAEPDefault,
		"RSA-OAEP-256":       NewRSAOAEP256,
		"A128KW":             NewA128KW,
		"A192KW":             NewA192KW,
		"A256KW":             NewA256KW,
		"A128GCMKW":          NewA128GCMKW,
		"A192GCMKW":          NewA192GCMKW,
		"A256GCMKW":          NewA256GCMKW,
		"ECDH-ES":            NewECDHES,
		"ECDH-ES+A128KW":     NewECDHESA128KW,
		"ECDH-ES+A192KW":     NewECDHESA192KW,
		"ECDH-ES+A256KW":     NewECDHESA256KW,
		"PBES2-HS256+A128KW": NewPBES2HS256A128KW,
		"PBES2-HS384+A192KW": NewPBES2HS384A192KW,
		"PBES2-HS512+A256KW": NewPBES2HS512A256KW,
	}
)
