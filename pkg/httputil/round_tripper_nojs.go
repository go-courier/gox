//go:build !js
// +build !js

package httputil

import (
	"context"
	"net"
	"net/http"
	"time"
)

func GetShortConnClientContext(ctx context.Context, timeout time.Duration, transports ...Transport) *http.Client {
	t := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 0,
		}).DialContext,
		DisableKeepAlives:     true,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Timeout:   timeout,
		Transport: t,
	}

	for i := range transports {
		client.Transport = transports[i](client.Transport)
	}

	return client
}
