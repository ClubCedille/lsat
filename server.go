package main

import (
	"fmt"
	"lsat/auth"
	"lsat/macaroon"
	"lsat/mock"
	"lsat/secrets"
	"net/http"
	"strings"

	"github.com/lightningnetwork/lnd/lntypes"
)

type Handler struct {
	*auth.Minter
}

var serviceLimiter = mock.NewServiceLimiter()
var secretStore = mock.NewTestStore()
var challenger = mock.NewChallenger()

func HttpServer() {

	// Initialize your Server instance
	minter := auth.NewMinter(&serviceLimiter, &secretStore, challenger)

	// Create a Handler with access to the Minter
	handle := &Handler{Minter: &minter}

	fmt.Println("Server launched!")
	http.HandleFunc("/protected", handle.handleAuthorization)
	http.HandleFunc("/", handle.handleAuthentication) // Retourné dans le cas où la platforme reçoit un Token invalide
	err := http.ListenAndServe("localhost:8080", nil) // Ca devrait etre lié à la platforme?

	fmt.Println(err)
}

func (h *Handler) handleAuthorization(w http.ResponseWriter, r *http.Request) {
	// Extract the Authorization header
	authHeader := r.Header.Get("Authorization")

	// Check if Authorization header is present
	// Parse the Authorization header
	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || parts[0] != "L402" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Unknown request!")
		return
	}

	Macaroon, _ := macaroon.DecodeBase64(parts[1])

	err := serviceLimiter.VerifyMacaroon(&Macaroon)

	if err == nil {
		// Respond with success (for demonstration purposes)
		// We should respond with the ressource
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Request authorized: %s", Macaroon)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Authorization failed! %s", err)
	}
}

func (h *Handler) handleAuthentication(w http.ResponseWriter, r *http.Request) {
	// Extract the Authorization header
	authHeader := r.Header.Get("Authorization")

	// Check if Authorization header is present
	// Parse the Authorization header
	parts := strings.Split(authHeader, " ")

	// Extract the Macaroon and preimage from the Authorization header
	if len(parts) != 2 || parts[0] != "L402" {
		// Create a new UserId
		uid := secrets.NewUserId()

		// Invalid Authorization header format, respond with 402 Payment Required and WWW-Authenticate header
		pretoken, err := h.Minter.MintToken(uid, mock.CatService)

		if err != nil {
			fmt.Println(err)
			return
		}

		macaroon := pretoken.Macaroon.String()

		// Format Macaroon and invoice in WWW-Authenticate header
		authHeader := fmt.Sprintf("L402 macaroon=\"%s\", invoice=\"%s\"", macaroon, &pretoken.PaymentRequest)

		// Set the WWW-Authenticate header
		w.Header().Set("WWW-Authenticate", authHeader)

		w.WriteHeader(http.StatusPaymentRequired)
		fmt.Fprint(w, "Payment Required")
		return
	}

	credentials := strings.Split(parts[1], ":")

	Macaroon, _ := macaroon.DecodeBase64(credentials[0])
	Preimage, _ := lntypes.MakePreimageFromStr(credentials[1])

	token := macaroon.Token{
		Macaroon: Macaroon,
		Preimage: Preimage,
	}

	// Process the request with the extracted Macaroon and preimage
	// Your logic for handling the request goes here...
	signedMac, err := h.Minter.AuthToken(&token)

	if err == nil {
		// Respond with success (for demonstration purposes)
		// We should respond with the ressource
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Access token: %s", signedMac.ToJSON())
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Authentification failed! %s", err)
	}
}
