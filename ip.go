package hermes

import (
	"net"
	"net/http"
	"strings"
)

var (
	xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
	xRealIP       = http.CanonicalHeaderKey("X-Real-IP")
)

func GetIP(r *http.Request) net.IP {
	// cloudflare always prioritized
	cloudflareIP := r.Header.Get("CF-Connecting-IP")
	if cloudflareIP != "" {
		return net.ParseIP(cloudflareIP)
	}

	if xff := r.Header.Get(xForwardedFor); xff != "" {
		i := strings.Index(xff, ", ")

		if i == -1 {
			i = len(xff)
		}

		return net.ParseIP(xff[:i])
	}

	if ip := r.Header.Get(xRealIP); ip != "" {
		return net.ParseIP(ip)
	}

	h, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return net.IP{}
	}

	return net.ParseIP(h)
}
