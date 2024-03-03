package main

import (
	"errors"
	"log"
)

type Store interface {
	CreateLink(link string) (string, error)
	FindURL(url string) (bool, int, error)
	GetLink(id int) (string, error)
}

type Storage struct {
	*SQLiteDB
}

func NewStorage(db *SQLiteDB) Store {
	return &Storage{
		db,
	}
}

func (s *Storage) CreateLink(link string) (string, error) {
	_, err := s.db.Exec("INSERT INTO links (url) VALUES (?)", link)
	if err != nil {
		log.Printf("%-20s error creating link in db. Err: %v", "CreateLink", err)
		return "", errors.New("error creating link")
	}
	return "", nil
}

func (s *Storage) GetLink(id int) (string, error) {
	var link string
	err := s.db.QueryRow("SELECT url FROM links WHERE id = ?", id).Scan(&link)
	if err != nil {
		log.Printf("%-20s error getting link with id: %d from db. Returning blank. Err: %v", "GetLink", id, err)
		return "", errors.New("link not found")
	}
	return link, nil
}

func (s *Storage) FindURL(url string) (isPresent bool, url_id int, err error) {
	sqlError := s.db.QueryRow("SELECT (id,url) FROM links WHERE url = ?", url).Scan(&url_id)
	if err != nil {
		log.Printf("%-20s error finding url: %v from db.  Err: %v", "FindURL", url, sqlError)
		err = errors.New("url not found")
		return
	}
	if url_id != 0 {
		isPresent = true
		return
	}
	return
}
