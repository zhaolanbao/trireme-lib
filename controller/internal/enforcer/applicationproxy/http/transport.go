package httpproxy

import (
	"net/http"
)

// TriremeRoundTripper is the Trireme RoundTripper that will handle
// responses.
type TriremeRoundTripper struct {
	http.RoundTripper
}

// NewTriremeRoundTripper creates a new RoundTripper that handles the
// responses.
func NewTriremeRoundTripper(r http.RoundTripper) *TriremeRoundTripper {
	return &TriremeRoundTripper{
		RoundTripper: r,
	}
}

// RoundTrip implements the RoundTripper interface. It will add a cookie
// in the response in case of OIDC requests with refresh tokens.
func (t *TriremeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {

	res, err := t.RoundTripper.RoundTrip(req)
	if err != nil || res == nil {
		return res, err
	}

	data := req.Context().Value(statsContextKey)
	if data == nil {
		return res, nil
	}

	state, ok := data.(*connectionState)
	if ok && state.cookie == nil {
		return res, nil
	}

	if v := state.cookie.String(); v != "" {
		res.Header.Add("Set-Cookie", v)
	}

	return res, nil
}
