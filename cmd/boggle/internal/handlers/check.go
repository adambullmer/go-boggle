package handlers

import (
	"net/http"

	"github.com/adambullmer/go-boggle/internal/platform/web"
)

func health(w http.ResponseWriter, r *http.Request) {
	var health struct {
		Status string `json:"status"`
	}

	health.Status = "ok"
	web.Respond(w, health, 200)
}
