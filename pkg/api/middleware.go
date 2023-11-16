// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/sassoftware/event-provenance-registry/pkg/metrics"
)

// using this to get the status code from the request so we can pass it to Prometheus.
type responseWrapper struct {
	http.ResponseWriter
	status int
}

func (r *responseWrapper) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func requestCounter() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			wrapper := responseWrapper{w, 200}
			next.ServeHTTP(&wrapper, r)
			logger.V(1).Info(fmt.Sprintf("returned status %v", wrapper.status))
			metrics.Requests.WithLabelValues(strconv.Itoa(wrapper.status), r.Method).Inc()
		}
		return http.HandlerFunc(fn)
	}
}

func requestTimer() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			elapsed := time.Since(start).Seconds()
			metrics.ResponseTimes.Observe(elapsed)
		}
		return http.HandlerFunc(fn)
	}
}

/*
Set the security headers on the response. This is primarily to protect against XSS attacks.
	First don't do inline javascript...just stop.
	Second you have to generate a hash for all the javascript files you plan to call
	Third just avoid javascript and you won't have to do this
	Fourth whitespace and line breaks matter
	Avoid adding the following directives
	   'unsafe-inline' 	Allows the usage of inline scripts or styles.
	   'unsafe-eval' 	Allows the usage of eval in scripts.
	https://cheatsheetseries.owasp.org/cheatsheets/Content_Security_Policy_Cheat_Sheet.html

# Generate HASH for script src

JS=$(find ./resources -type f  -name "*.js")

echo -en "script-src 'self' "

for j in $JS;do
	HASH=$(cat ${j} | openssl sha256 -binary | openssl base64)
	echo -en "'sha256-${HASH}' "
done

echo -en ";"

# Generate HASH for style-strconv

CSS=$(find ./resources -type f  -name "*.css")

echo -en "style-src 'self' "

for j in $CSS;do
	HASH=$(cat ${j} | openssl sha256 -binary | openssl base64)
	echo -en "'sha256-${HASH}' "
done

echo -en ";"
*/

func securityHeaders() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// CSP whitelists resources from places that we retrieve them from.
			// TODO: add back in proper Content Security Policy hashes. if we decide to keep the graphql UI we need to add them
			// to the allowed list of sources
			// w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'sha256-pyVPiLlnqL9OWVoJPs/E6VVF5hBecRzM2gBiarnaqAo=';"+
			// "  script-src 'self' 'sha512-Wr9OKCTtq1anK0hq5bY3X/AvDI5EflDSAh0mE9gma+4hl+kXdTJPKZ3TwLMBcrgUeoY0s3dq9JjhCQc7vddtFg==' 'sha512-Vf2xGDzpqUOEIKO+X2rgTLWPY+65++WPwCHkX2nFMu9IcstumPsf/uKKRd5prX3wOu8Q0GBylRpsDB26R6ExOg==' 'sha256-Xlt9flxaXkphLAysOndSSWNzqJqv2IMqjKSHJD4oyGU=' 'sha256-DRtIv6ccKNt2rE4esiPGKe0IuOC07WKQTprwIE3WfTk=' 'sha256-CTX/EsBKST60M8ohR15xRvYYiWO0ryobO8uHPLdeBTQ=' 'sha256-CTX/EsBKST60M8ohR15xRvYYiWO0ryobO8uHPLdeBTQ='; img-src 'self' https://unpkg.com/graphiql@2.3.0/graphiql.min.css  https://unpkg.com/graphiql@2.3.0/graphiql.min.js w3.org *.swagger.io data:;")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("X-Frame-Options", "deny")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// logOrigin log the origin of our httprequest
func LogOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hdr := "X-Real-Ip"
		origin := r.Header.Get("X-Real-Ip")

		if origin == "" {
			hdr = "X-Forwarded-For"
			origin = r.Header.Get("X-Forwarded-For")
		}
		if origin == "" {
			hdr = "RemoteAddr"
			origin = r.RemoteAddr
		}
		logger.Info(fmt.Sprintf("request origin %s : %s", hdr, origin))
		next.ServeHTTP(w, r)
	})
}
