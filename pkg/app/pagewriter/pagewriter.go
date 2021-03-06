package pagewriter

import (
	"fmt"
	"net/http"
)

// Writer is an interface for rendering html templates for both sign-in and
// error pages.
// It can also be used to write errors for the http.ReverseProxy used in the
// upstream package.
type Writer interface {
	WriteSignInPage(rw http.ResponseWriter, redirectURL string)
	WriteErrorPage(rw http.ResponseWriter, status int, redirectURL string, appError string, messages ...interface{})
	ProxyErrorHandler(rw http.ResponseWriter, req *http.Request, proxyErr error)
}

// pageWriter implements the Writer interface
type pageWriter struct {
	*errorPageWriter
	*signInPageWriter
}

// Opts contains all options required to configure the template
// rendering within OAuth2 Proxy.
type Opts struct {
	// TemplatesPath is the path from which to load custom templates for the sign-in and error pages.
	TemplatesPath string

	// ProxyPrefix is the prefix under which OAuth2 Proxy pages are served.
	ProxyPrefix string

	// Footer is the footer to be displayed at the bottom of the page.
	// If not set, a default footer will be used.
	Footer string

	// Version is the OAuth2 Proxy version to be used in the default footer.
	Version string

	// Debug determines whether errors pages should be rendered with detailed
	// errors.
	Debug bool

	// DisplayLoginForm determines whether or not the basic auth password form is displayed on the sign-in page.
	DisplayLoginForm bool

	// ProviderName is the name of the provider that should be displayed on the login button.
	ProviderName string

	// SignInMessage is the messge displayed above the login button.
	SignInMessage string
}

// NewWriter constructs a Writer from the options given to allow
// rendering of sign-in and error pages.
func NewWriter(opts Opts) (Writer, error) {
	templates, err := loadTemplates(opts.TemplatesPath)
	if err != nil {
		return nil, fmt.Errorf("error loading templates: %v", err)
	}

	errorPage := &errorPageWriter{
		template:    templates.Lookup("error.html"),
		proxyPrefix: opts.ProxyPrefix,
		footer:      opts.Footer,
		version:     opts.Version,
		debug:       opts.Debug,
	}

	signInPage := &signInPageWriter{
		template:         templates.Lookup("sign_in.html"),
		errorPageWriter:  errorPage,
		proxyPrefix:      opts.ProxyPrefix,
		providerName:     opts.ProviderName,
		signInMessage:    opts.SignInMessage,
		footer:           opts.Footer,
		version:          opts.Version,
		displayLoginForm: opts.DisplayLoginForm,
	}

	return &pageWriter{
		errorPageWriter:  errorPage,
		signInPageWriter: signInPage,
	}, nil
}
