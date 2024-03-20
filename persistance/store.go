package persistance

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/salvaft/go-link-shortener/utils"
)

type Store interface {
	CreateLink(link string) (int, error)
	FindURL(url string) (int, error)
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

// TODO: Return link struct
func (s *Storage) CreateLink(link string) (int, error) {
	row, err := s.db.Exec("INSERT INTO links (url) VALUES (?)", link)
	if err != nil {
		utils.Logger.Printf("%-20s error creating link in db. Err: %v", "CreateLink", err)
		return -1, errors.New("error creating link")
	}

	id, err := row.LastInsertId()
	if err != nil {
		utils.Logger.Printf("%-20s error getting link id after insert. Err: %v", "CreateLink", err)
		return -1, errors.New("error getting link id after insert")
	}
	return int(id), nil
}

func (s *Storage) GetLink(id int) (string, error) {
	var link string
	err := s.db.QueryRow("SELECT url FROM links WHERE id = ?", id).Scan(&link)
	if err == sql.ErrNoRows {
		utils.Logger.Printf("%-20s link with id: %d not found in db. Returning blank", "GetLink", id)
		return "", &URLNotFound{fmt.Sprintf("link with id %d not found", id)}

	} else if err != nil {
		utils.Logger.Printf("%-20s error getting link with id: %d from db. Returning blank. Err: %v", "GetLink", id, err)
		return "", err
	}
	return link, nil
}

func (s *Storage) FindURL(url string) (int, error) {
	var urlID int
	row := s.db.QueryRow("SELECT id FROM links WHERE url = ? LIMIT 1", url)

	if sqlError := row.Scan(&urlID); sqlError != nil {
		if sqlError == sql.ErrNoRows {
			utils.Logger.Printf("%-20s link with id: %d not found in db. Returning blank", "FindURL", urlID)
			return -1, &URLNotFound{fmt.Sprintf("link with id %d not found", urlID)}

		} else {
			utils.Logger.Printf("%-20s error finding link %v from db. Returning -1 Err: %v", "FindURL", url, sqlError)
			return -1, sqlError
		}
	}

	utils.Logger.Printf("%-20s url: %v found in db. ID: %v", "FindURL", url, urlID)
	return urlID, nil
}
