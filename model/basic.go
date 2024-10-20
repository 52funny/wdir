package model

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

type BasicAuth struct {
	inner     map[string]string
	innerHash map[string]bool
}

func NewBasicAuth(username []string, password []string) *BasicAuth {
	m := make(map[string]string)
	mHash := make(map[string]bool)
	for i := range username {
		m[username[i]] = password[i]
		baseEnc := base64.StdEncoding.EncodeToString([]byte(username[i] + ":" + password[i]))
		sha := sha256.Sum256([]byte(baseEnc))
		hex := hex.EncodeToString(sha[:])
		mHash[hex] = true
	}
	return &BasicAuth{
		inner:     m,
		innerHash: mHash,
	}
}

func (b *BasicAuth) Len() int {
	return len(b.inner)
}

func (b *BasicAuth) Get(username string) (string, bool) {
	v, ok := b.inner[username]
	return v, ok
}

func (b *BasicAuth) Auth(auth []byte) bool {
	if b.Len() == 0 {
		return true
	}

	i := bytes.IndexByte(auth, ' ')
	if i == -1 {
		return false
	}

	if !bytes.EqualFold(auth[:i], []byte("basic")) {
		return false
	}

	decoded, err := base64.StdEncoding.DecodeString(string(auth[i+1:]))
	if err != nil {
		return false
	}

	credentials := bytes.Split(decoded, []byte(":"))
	if len(credentials) <= 1 {
		return false
	}

	user := credentials[0]
	pass := credentials[1]

	if p, ok := b.Get(string(user)); ok {
		return p == string(pass)
	}
	return false
}

func (b *BasicAuth) AuthToken(token []byte) bool {
	if b.Len() == 0 {
		return true
	}

	_, ok := b.innerHash[string(token)]
	return ok
}
