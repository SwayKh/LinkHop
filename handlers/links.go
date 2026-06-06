package handlers

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"urlshortner/database"
	"urlshortner/middleware"
	"urlshortner/render"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err == nil && cookie.Value != "" {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	render.Template(w, "index", &render.Data{})
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	links, err := database.GetLinksByUserID(userID)
	if err != nil {
		render.Template(w, "dashboard", &render.Data{
			User:  &render.UserInfo{ID: userID, Email: middleware.GetUserEmail(r)},
			Error: "Failed to load links",
		})
		return
	}

	var linkInfos []render.LinkInfo
	for _, l := range links {
		linkInfos = append(linkInfos, render.LinkInfo{
			ID:           l.ID,
			OriginalURL:  l.OriginalURL,
			ShortCode:    l.ShortCode,
			CustomAlias:  l.CustomAlias,
			ClickCount:   l.ClickCount,
			LastAccessed: l.LastAccessed,
			CreatedAt:    l.CreatedAt,
			UpdatedAt:    l.UpdatedAt,
		})
	}

	render.Template(w, "dashboard", &render.Data{
		User:    &render.UserInfo{ID: userID, Email: middleware.GetUserEmail(r)},
		Links:   linkInfos,
		Success: r.URL.Query().Get("success"),
		Error:   r.URL.Query().Get("error"),
	})
}

