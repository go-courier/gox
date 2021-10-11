package httputil

import "net/http"

type Transport = func(next http.RoundTripper) http.RoundTripper
