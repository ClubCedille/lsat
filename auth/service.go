package auth

import (
	"context"
	"lsat/macaroon"
)

// ServiceLimiter is an interface defining methods for managing services and their capabilities.
type ServiceLimiter interface {
	// Services retrieves information about services with the provided names.
	Services(context.Context, ...string) ([]macaroon.Service, error)

	// Capabilities retrieves the capabilities associated with the provided services.
	// It returns a list of caveats representing the capabilities of the services.
	Capabilities(context.Context, ...macaroon.Service) ([]macaroon.Caveat, error)

	// VerifyCaveats verifies the validity of a set of macaroon caveats.
	// It returns an error if any of the caveats are invalid or do not meet the specified criteria.
	VerifyCaveats(...macaroon.Caveat) error
}
