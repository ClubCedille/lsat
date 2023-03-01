package main

import (
	"lsat/secrets"
	"testing"
)

func TestSecret(t *testing.T) {
	store := secrets.NewTestStore()

	uid := store.CreateUser()

	asecret, err := store.Secret(uid)

	if err != nil {
		t.Error("Secret not found in the Store.")
	}

	bsecret, _ := store.Secret(uid)

	if asecret != bsecret {
		t.Error("Mismatch between the secret from the same user.")
	}
}
