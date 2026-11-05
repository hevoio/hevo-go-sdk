package pipelines

import "maps"

// SourceAuth carries the credentials of your own OAuth app and a token it issued. Hevo exchanges
// them for an access token when the Pipeline is created and refreshes it from then on. Build it
// with the client's NewSourceAuth.
type SourceAuth struct {
	clientID     *string
	clientSecret *string
	refreshToken *string
	accessToken  *string
	extra        map[string]string
}

type sourceAuthRequest struct {
	ClientAccess *clientAccessRequest `json:"client_access,omitempty"`
	RefreshToken *string              `json:"refresh_token,omitempty"`
	AccessToken  *string              `json:"access_token,omitempty"`
}

type clientAccessRequest map[string]string

// ClientID sets the OAuth app's client id.
func (a *SourceAuth) ClientID(value string) *SourceAuth {
	a.clientID = &value

	return a
}

// ClientSecret sets the OAuth app's client secret.
func (a *SourceAuth) ClientSecret(value string) *SourceAuth {
	a.clientSecret = &value

	return a
}

// ClientAccessValue sets another app-level credential the provider needs beside the client id
// and secret, for example Google Ads' developer_token.
func (a *SourceAuth) ClientAccessValue(key, value string) *SourceAuth {
	if a.extra == nil {
		a.extra = map[string]string{}
	}

	a.extra[key] = value

	return a
}

// RefreshToken sets the refresh token the app issued. Most providers use it.
func (a *SourceAuth) RefreshToken(value string) *SourceAuth {
	a.refreshToken = &value

	return a
}

// AccessToken sets the access token, for a provider that refreshes by exchanging the live
// access token or does not refresh at all.
func (a *SourceAuth) AccessToken(value string) *SourceAuth {
	a.accessToken = &value

	return a
}

// request returns the credentials as the API expects them.
func (a *SourceAuth) request() *sourceAuthRequest {
	if a == nil {
		return nil
	}

	access := clientAccessRequest{}
	maps.Copy(access, a.extra)

	if a.clientID != nil {
		access["client_id"] = *a.clientID
	}

	if a.clientSecret != nil {
		access["client_secret"] = *a.clientSecret
	}

	return &sourceAuthRequest{
		ClientAccess: &access,
		RefreshToken: a.refreshToken,
		AccessToken:  a.accessToken,
	}
}
