package dependency

import (
	"database/sql"

	"github.com/vacp2p/mvds/state"
)

// Verify that MessageDependency interface is implemented.
var _ MessageDependency = (*sqliteDependency)(nil)

type sqliteDependency struct {
	db *sql.DB
}

func NewPersistentDependency(db *sql.DB) *sqliteDependency {
	return &sqliteDependency{db: db}
}

func (sd *sqliteDependency) Add(msg, dependency state.MessageID) error {
	_, err := sd.db.Exec(`INSERT INTO mvds_dependencies (msg, dependency) VALUES (?, ?)`, msg[:], dependency[:])
	if err != nil {
		return err
	}

	return nil
}

func (sd *sqliteDependency) Dependants(id state.MessageID) ([]state.MessageID, error) {
	rows, err := sd.db.Query(`SELECT msg FROM mvds_dependencies WHERE dependency = ?`, id[:])
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	msgs := make([]state.MessageID, 0)

	for rows.Next() {
		var msg state.MessageID
		err := rows.Scan(&msg)
		if err != nil {
			return nil, err
		}

		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func (sd *sqliteDependency) MarkResolved(msg state.MessageID, dependency state.MessageID) error {
	result, err := sd.db.Exec(
		`DELETE FROM mvds_dependencies WHERE msg = ? AND dependency = ?`,
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

func (sd *sqliteDependency) HasUnresolvedDependencies(id state.MessageID) (bool, error) {
	result := sd.db.QueryRow(`SELECT COUNT(*) FROM mvds_dependencies WHERE msg = ?`, id[:])
	var num int64
	err := result.Scan(num)
	if err != nil {
		return false, err
	}

	if num > 0 {
		return true, nil
	}
}

