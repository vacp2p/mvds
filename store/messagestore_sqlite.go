package store

import (
	"database/sql"
	"errors"

	"github.com/vacp2p/mvds/state"

	"github.com/vacp2p/mvds/protobuf"
)

var (
	ErrMessageNotFound = errors.New("message not found")
)

type persistentMessageStore struct {
	db *sql.DB
}

func NewPersistentMessageStore(db *sql.DB) *persistentMessageStore {
	return &persistentMessageStore{db: db}
}

func (p *persistentMessageStore) Add(message *protobuf.Message) error {
	id := message.ID()

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`INSERT INTO mvds_messages (id, group_id, timestamp, body)
		VALUES (?, ?, ?, ?)`,
		id[:],
		message.GroupId,
		message.Timestamp,
		message.Body,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	if message.Metadata != nil && len(message.Metadata.Parents) > 0 {
		query := "INSERT INTO mvds_parents(message, parent) VALUES "
		vals := make([]interface{}, 0)

		for _, row := range message.Metadata.Parents {
			query += "(?, ?),"
			vals = append(vals, id[:], row[:])
		}

		stmt, err := tx.Prepare(query[0:len(query)-1])
		if err != nil {
			return err
		}

		_, err = stmt.Exec(vals...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (p *persistentMessageStore) Get(id state.MessageID) (*protobuf.Message, error) {
	var message protobuf.Message
	row := p.db.QueryRow(
		`SELECT group_id, timestamp, body FROM mvds_messages WHERE id = ?`,
		id[:],
	)
	if err := row.Scan(
		&message.GroupId,
		&message.Timestamp,
		&message.Body,
	); err != nil {
		return nil, err
	}

	message.Metadata = &protobuf.Metadata{Ephemeral: false, Parents: make([][]byte, 0)}

	rows, err := p.db.Query(`SELECT parent FROM mvds_parents WHERE message = ?`, id[:])
	if err != nil {
		return nil, err
	}

	var parent []byte

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&parent)
		if err != nil {
			return nil, err
		}

		message.Metadata.Parents = append(message.Metadata.Parents, parent)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (p *persistentMessageStore) Has(id state.MessageID) (bool, error) {
	var result bool
	err := p.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM mvds_messages WHERE id = ?)`,
		id[:],
	).Scan(&result)
	switch err {
	case sql.ErrNoRows:
		return false, ErrMessageNotFound
	case nil:
		return result, nil
	default:
		return false, err
	}
}

func (p *persistentMessageStore) GetMessagesWithoutChildren(id state.GroupID) ([]state.MessageID, error) {
	result := make([]state.MessageID, 0)
	rows, err := p.db.Query(
		`SELECT id FROM mvds_messages WHERE group_id = ? AND id NOT IN (SELECT parent FROM mvds_parents)`,
		id[:],
	)

	if err != nil {
		return nil, err
	}

	var parent state.MessageID

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&parent)
		if err != nil {
			return nil, err
		}

		result = append(result, parent)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}
