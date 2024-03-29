package main

import (
	"lsat/macaroon"
	"lsat/mock"
	"testing"
)

var service macaroon.Service = macaroon.NewService("image", 1000)

var caveat macaroon.Caveat = macaroon.NewCaveat("image/png", "test.png")

func TestSecret(t *testing.T) {
	uid := secretStore.CreateUser()

	asecret, err := secretStore.NewSecret(uid)

	if err != nil {
		t.Error("Secret not found in the Store.")
	}

	bsecret, _ := secretStore.GetSecret(uid)

	if asecret != bsecret {
		t.Error("Mismatch between the secret from the same user.")
	}
}

func TestUid(t *testing.T) {
	store := mock.NewTestStore()

	userA := store.CreateUser()

	userB := store.CreateUser()

	if userA == userB {
		t.Error("Two users cannot have the same id.")
	}
}

func TestMacaroon(t *testing.T) {
	uid := secretStore.CreateUser()

	secret, _ := secretStore.NewSecret(uid)

	oven := macaroon.NewOven(secret)

	mac, _ := oven.WithCaveats(caveat).WithService(service).Cook()

	signaturea := mac.Signature()

	uid = secretStore.CreateUser()

	secret, _ = secretStore.NewSecret(uid)

	oven = macaroon.NewOven(secret)

	mac, _ = oven.WithCaveats(caveat).WithService(service).Cook()

	signatureb := mac.Signature()

	if signaturea == signatureb {
		t.Error("Two users cannot produce the same signature for the same given caveats.")
	}
}
