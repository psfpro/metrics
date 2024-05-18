package handler

import (
	"net"
	"net/http"
)

type TrustedSubnet struct {
	trustedSubnet string
}

func NewTrustedSubnet(trustedSubnet string) *TrustedSubnet {
	return &TrustedSubnet{trustedSubnet: trustedSubnet}
}

func (h *TrustedSubnet) Check(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		_, ipNet, _ := net.ParseCIDR(h.trustedSubnet)
		header := r.Header.Get("X-Real-IP")
		ip := net.ParseIP(header)
		if ip != nil && ipNet != nil && !ipNet.Contains(ip) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
