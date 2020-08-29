package models

import (
	"database/sql"
)

const (
	imgurBaseLink = "https://i.imgur.com/"
	nameSuffix    = ".jpg"
)

// MemeModel defines the database which the functions operate on.
type MemeModel struct {
	DB *sql.DB
}

// MemeEntry represents an entry of a meme in the database.
type MemeEntry struct {
	Name string
	Link string
}

// GetAll returns a list of all memes.
func (m *MemeModel) GetAll() ([]MemeEntry, error) {
	res := []MemeEntry{}

	stmt := `SELECT name, url FROM memes ORDER BY name ASC`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		entry := MemeEntry{}
		err = rows.Scan(&entry.Name, &entry.Link)
		if err != nil {
			return res, err
		}

		entry.Name += nameSuffix
		entry.Link = imgurBaseLink + entry.Link

		res = append(res, entry)
	}

	return res, nil
}

// Get returns the image URL of the meme if it exists.
func (m *MemeModel) Get(name string) (string, error) {
	var res string
	stmt := `SELECT url FROM memes WHERE name = $1`
	row := m.DB.QueryRow(stmt, name)
	err := row.Scan(&res)
	if err != nil {
		return "", err
	}

	return imgurBaseLink + res, nil
}

// Insert inserts a meme entry to the database.
func (m *MemeModel) Insert(name string, url string) error {
	stmt := `INSERT INTO memes (name, url) VALUES ($1, $2)`
	_, err := m.DB.Exec(stmt, name, url)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes a meme entry from the database.
func (m *MemeModel) Delete(name string) error {
	stmt := `DELETE FROM memes WHERE name = $1`
	_, err := m.DB.Exec(stmt, name)
	if err != nil {
		return err
	}

	return nil
}
