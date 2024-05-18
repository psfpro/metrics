package handler

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTrustedSubnet_Check(t *testing.T) {
	type fields struct {
		trustedSubnet string
	}
	type args struct {
		ip string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "real ip match",
			fields: fields{
				trustedSubnet: "1.1.1.0/24",
			},
			args: args{
				ip: "1.1.1.1",
			},
			want: http.StatusOK,
		},
		{
			name: "real ip don't match",
			fields: fields{
				trustedSubnet: "1.1.1.1/24",
			},
			args: args{
				ip: "2.2.2.2",
			},
			want: http.StatusForbidden,
		},
		{
			name: "real ip empty",
			fields: fields{
				trustedSubnet: "1.1.1.0/24",
			},
			args: args{
				ip: "",
			},
			want: http.StatusOK,
		},
		{
			name: "trusted subnet empty",
			fields: fields{
				trustedSubnet: "",
			},
			args: args{
				ip: "1.1.1.1",
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &TrustedSubnet{
				trustedSubnet: tt.fields.trustedSubnet,
			}
			handler := h.Check(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			}))
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.Header.Add("X-Real-IP", tt.args.ip)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)
			assert.Equal(t, tt.want, w.Result().StatusCode)
		})
	}
}
