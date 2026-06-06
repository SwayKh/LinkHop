package database

import "time"

type Link struct {
	ID           int64
	UserID       int64
	OriginalURL  string
	ShortCode    string
	CustomAlias  string
	ClickCount   int
	LastAccessed string
	CreatedAt    string
	UpdatedAt    string
}

func CreateLink(userID int64, originalURL, shortCode, customAlias string) error {
	_, err := DB.Exec(
		"INSERT INTO links (user_id, original_url, short_code, custom_alias) VALUES (?, ?, ?, ?)",
		userID, originalURL, shortCode, customAlias,
	)
	return err
}

func GetLinksByUserID(userID int64) ([]Link, error) {
	rows, err := DB.Query(
		`SELECT id, user_id, original_url, short_code, COALESCE(custom_alias, ''), 
		        click_count, COALESCE(last_accessed, ''), created_at, updated_at 
		 FROM links WHERE user_id = ? ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []Link
	for rows.Next() {
		var l Link
		if err := rows.Scan(&l.ID, &l.UserID, &l.OriginalURL, &l.ShortCode, &l.CustomAlias,
			&l.ClickCount, &l.LastAccessed, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		links = append(links, l)
	}
	return links, nil
}

func GetLinkByShortCode(shortCode string) (*Link, error) {
	l := &Link{}
	err := DB.QueryRow(
		`SELECT id, user_id, original_url, short_code, COALESCE(custom_alias, ''), 
		        click_count, COALESCE(last_accessed, ''), created_at, updated_at 
		 FROM links WHERE short_code = ?`, shortCode,
	).Scan(&l.ID, &l.UserID, &l.OriginalURL, &l.ShortCode, &l.CustomAlias,
		&l.ClickCount, &l.LastAccessed, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func GetLinkByID(id int64) (*Link, error) {
	l := &Link{}
	err := DB.QueryRow(
		`SELECT id, user_id, original_url, short_code, COALESCE(custom_alias, ''), 
		        click_count, COALESCE(last_accessed, ''), created_at, updated_at 
		 FROM links WHERE id = ?`, id,
	).Scan(&l.ID, &l.UserID, &l.OriginalURL, &l.ShortCode, &l.CustomAlias,
		&l.ClickCount, &l.LastAccessed, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func UpdateLink(id int64, originalURL, customAlias string) error {
	_, err := DB.Exec(
		"UPDATE links SET original_url = ?, custom_alias = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		originalURL, customAlias, id,
	)
	return err
}

func DeleteLink(id int64) error {
	_, err := DB.Exec("DELETE FROM links WHERE id = ?", id)
	return err
}
