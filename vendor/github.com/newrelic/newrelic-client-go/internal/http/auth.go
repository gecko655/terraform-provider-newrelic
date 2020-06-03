package http

import (
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// RequestAuthorizer is an interface that allows customizatino of how a request is authorized.
type RequestAuthorizer interface {
	AuthorizeRequest(r *Request, c *config.Config)
}

// NerdGraphAuthorizer authorizes calls to NerdGraph.
type NerdGraphAuthorizer struct{}

// AuthorizeRequest is responsible for setting up auth for a request.
func (a *NerdGraphAuthorizer) AuthorizeRequest(r *Request, c *config.Config) {
	r.SetHeader("Api-Key", c.PersonalAPIKey)
}

// PersonalAPIKeyCapableV2Authorizer authorizes V2 endpoints that can use a personal API key.
type PersonalAPIKeyCapableV2Authorizer struct{}

// AuthorizeRequest is responsible for setting up auth for a request.
func (a *PersonalAPIKeyCapableV2Authorizer) AuthorizeRequest(r *Request, c *config.Config) {
	if c.PersonalAPIKey != "" {
		r.SetHeader("Api-Key", c.PersonalAPIKey)
		r.SetHeader("Auth-Type", "User-Api-Key")
	} else {
		r.SetHeader("X-Api-Key", c.AdminAPIKey)
	}
}

// ClassicV2Authorizer authorizes V2 endpoints that cannot use a personal API key.
type ClassicV2Authorizer struct{}

// AuthorizeRequest is responsible for setting up auth for a request.
func (a *ClassicV2Authorizer) AuthorizeRequest(r *Request, c *config.Config) {
	r.SetHeader("X-Api-Key", c.AdminAPIKey)
}

// InsightsInsertKeyAuthorizer authorizes sending custom events to New Relic.
type InsightsInsertKeyAuthorizer struct{}

func (a *InsightsInsertKeyAuthorizer) AuthorizeRequest(r *Request, c *config.Config) {
	r.SetHeader("X-Insert-Key", c.InsightsInsertKey)
}
