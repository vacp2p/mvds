package dependency

import (
	"database/sql"

	"github.com/vacp2p/mvds/state"
)

// Verify that Tracker interface is implemented.
var _ Tracker = (*sqliteTracker)(nil)

type sqliteTracker struct {
	db *sql.DB
}

func NewPersistentTracker(db *sql.DB) *sqliteTracker {
	return &sqliteTracker{db: db}
}

func (sd *sqliteTracker) Add(msg, dependency state.MessageID) error {
	_, err := sd.db.Exec(`INSERT INTO mvds_dependencies (msg_id, dependency) VALUES (?, ?)`, msg[:], dependency[:])
	if err != nil {
		return err
	}

	return nil
}

func (sd *sqliteTracker) Dependants(id state.MessageID) ([]state.MessageID, error) {
	rows, err := sd.db.Query(`SELECT msg_id FROM mvds_dependencies WHERE dependency = ?`, id[:])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []state.MessageID

	for rows.Next() {
		var msg []byte
		err := rows.Scan(&msg)
		if err != nil {
			return nil, err
		}

		msgs = append(msgs, state.ToMessageID(msg))
	}

	return msgs, nil
}

func (sd *sqliteTracker) Resolve(msg state.MessageID, dependency state.MessageID) error {
	result, err := sd.db.Exec(
		`DELETE FROM mvds_dependencies WHERE msg_id = ? AND dependency = ?`,
		msg[:],
		dependency[:],
	)

	if err != nil {
		return err
	}

	if _, err := result.RowsAffected(); err != nil {
		return err
	}

	return nil
}

func (sd *sqliteTracker) IsResolved(id state.MessageID) (bool, error) {
	result := sd.db.QueryRow(`SELECT COUNT(*) FROM mvds_dependencies WHERE msg_id = ?`, id[:])
	var num int64
	err := result.Scan(&num)
	if err != nil {
		return false, err
	}

	return num == 0, nil
}

