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

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortCode() (string, error) {
	code := make([]byte, 6)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[n.Int64()]
	}
	return string(code), nil
}

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

func CreateLinkHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	if r.Method == "GET" {
		render.Template(w, "create", &render.Data{
			User: &render.UserInfo{ID: userID, Email: middleware.GetUserEmail(r)},
		})
		return
	}

	originalURL := strings.TrimSpace(r.FormValue("original_url"))
	customAlias := strings.TrimSpace(r.FormValue("custom_alias"))

	if originalURL == "" {
		render.Template(w, "create", &render.Data{
			User:  &render.UserInfo{ID: userID, Email: middleware.GetUserEmail(r)},
			Error: "Original URL is required",
		})
		return
	}

	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		originalURL = "https://" + originalURL
	}

	shortCode := customAlias
	if shortCode == "" {
		var err error
		shortCode, err = generateShortCode()
		if err != nil {
			render.Template(w, "create", &render.Data{
				User:  &render.UserInfo{ID: userID, Email: middleware.GetUserEmail(r)},
				Error: "Failed to generate short code",
			})
			return
		}
	}

	if err := database.CreateLink(userID, originalURL, shortCode, customAlias); err != nil {
		render.Template(w, "create", &render.Data{
			User:  &render.UserInfo{ID: userID, Email: middleware.GetUserEmail(r)},
			Error: "Custom alias already taken or failed to create link",
		})
		return
	}

	http.Redirect(w, r, "/dashboard?success=Link+created+successfully", http.StatusSeeOther)
}

func EditLinkHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	link, err := database.GetLinkByID(id)
	if err != nil || link.UserID != userID {
		render.Template(w, "edit", &render.Data{
			User:  &render.UserInfo{ID: userID, Email: middleware.GetUserEmail(r)},
			Error: "Link not found",
		})
		return
	}

	render.Template(w, "edit", &render.Data{
		User: &render.UserInfo{ID: userID, Email: middleware.GetUserEmail(r)},
		Link: &render.LinkInfo{
			ID:          link.ID,
			OriginalURL: link.OriginalURL,
			ShortCode:   link.ShortCode,
			CustomAlias: link.CustomAlias,
		},
	})
}

func UpdateLinkHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	link, err := database.GetLinkByID(id)
	if err != nil || link.UserID != userID {
		http.Redirect(w, r, "/dashboard?error=Link+not+found", http.StatusSeeOther)
		return
	}

	originalURL := strings.TrimSpace(r.FormValue("original_url"))
	customAlias := strings.TrimSpace(r.FormValue("custom_alias"))

	if originalURL == "" {
		render.Template(w, "edit", &render.Data{
			User:  &render.UserInfo{ID: userID, Email: middleware.GetUserEmail(r)},
			Link:  &render.LinkInfo{ID: id, OriginalURL: originalURL, CustomAlias: customAlias},
			Error: "Original URL is required",
		})
		return
	}

	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		originalURL = "https://" + originalURL
	}

	if err := database.UpdateLink(id, originalURL, customAlias); err != nil {
		render.Template(w, "edit", &render.Data{
			User:  &render.UserInfo{ID: userID, Email: middleware.GetUserEmail(r)},
			Link:  &render.LinkInfo{ID: id, OriginalURL: originalURL, CustomAlias: customAlias},
			Error: "Failed to update link",
		})
		return
	}

	http.Redirect(w, r, "/dashboard?success=Link+updated+successfully", http.StatusSeeOther)
}

func DeleteLinkHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	link, err := database.GetLinkByID(id)
	if err != nil || link.UserID != userID {
		http.Redirect(w, r, "/dashboard?error=Link+not+found", http.StatusSeeOther)
		return
	}

	if err := database.DeleteLink(id); err != nil {
		http.Redirect(w, r, "/dashboard?error=Failed+to+delete+link", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/dashboard?success=Link+deleted+successfully", http.StatusSeeOther)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")

	link, err := database.GetLinkByShortCode(shortCode)
	if err != nil {
		render.Template(w, "index", &render.Data{Error: "Link not found"})
		return
	}

	database.IncrementClickCount(link.ID)

	http.Redirect(w, r, link.OriginalURL, http.StatusMovedPermanently)
}

func AnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	link, err := database.GetLinkByID(id)
	if err != nil || link.UserID != userID {
		render.Template(w, "analytics", &render.Data{
			User:  &render.UserInfo{ID: userID, Email: middleware.GetUserEmail(r)},
			Error: "Link not found",
		})
		return
	}

	render.Template(w, "analytics", &render.Data{
		User: &render.UserInfo{ID: userID, Email: middleware.GetUserEmail(r)},
		Link: &render.LinkInfo{
			ID:           link.ID,
			OriginalURL:  link.OriginalURL,
			ShortCode:    link.ShortCode,
			CustomAlias:  link.CustomAlias,
			ClickCount:   link.ClickCount,
			LastAccessed: link.LastAccessed,
			CreatedAt:    link.CreatedAt,
		},
	})
}
