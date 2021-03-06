package mods

import "barista.run/pango"

var (
	spacer            = pango.Text(" ").XXSmall()
	gsuiteOauthConfig = []byte(`{"installed": {
	"client_id":"%%GOOGLE_CLIENT_ID%%",
	"project_id":"i3-barista",
	"auth_uri":"https://accounts.google.com/o/oauth2/auth",
	"token_uri":"https://www.googleapis.com/oauth2/v3/token",
	"auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs",
	"client_secret":"%%GOOGLE_CLIENT_SECRET%%",
	"redirect_uris":["urn:ietf:wg:oauth:2.0:oob","http://localhost"]
	}}`)
)