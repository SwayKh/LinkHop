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
		"INSERT INTO links (user_id, original_url, short_code, custom_alias) VALUES ($1, $2, $3, $4)",
		userID, originalURL, shortCode, customAlias,
	)
	return err
}

func GetLinksByUserID(userID int64) ([]Link, error) {
	rows, err := DB.Query(
		`SELECT id, user_id, original_url, short_code, COALESCE(custom_alias, ''), 
		        click_count, COALESCE(last_accessed::TEXT, ''), created_at, updated_at 
		 FROM links WHERE user_id = $1 ORDER BY created_at DESC`, userID,
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
		        click_count, COALESCE(last_accessed::TEXT, ''), created_at, updated_at 
		 FROM links WHERE short_code = $1`, shortCode,
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
		        click_count, COALESCE(last_accessed::TEXT, ''), created_at, updated_at 
		 FROM links WHERE id = $1`, id,
	).Scan(&l.ID, &l.UserID, &l.OriginalURL, &l.ShortCode, &l.CustomAlias,
		&l.ClickCount, &l.LastAccessed, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func UpdateLink(id int64, originalURL, customAlias string) error {
	_, err := DB.Exec(
		"UPDATE links SET original_url = $1, custom_alias = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3",
		originalURL, customAlias, id,
	)
	return err
}

func DeleteLink(id int64) error {
	_, err := DB.Exec("DELETE FROM links WHERE id = $1", id)
	return err
}

func IncrementClickCount(id int64) error {
	_, err := DB.Exec(
		"UPDATE links SET click_count = click_count + 1, last_accessed = $1 WHERE id = $2",
		time.Now(), id,
	)
	return err
}
