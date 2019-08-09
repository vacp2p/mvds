package state

import (
	"database/sql"
	"errors"
)

// Verify that SyncState interface is implemented.
var _ SyncState = (*sqliteSyncState)(nil)

type sqliteSyncState struct {
	db *sql.DB
}

func NewPersistentSyncState(db *sql.DB) *sqliteSyncState {
	return &sqliteSyncState{db: db}
}

func (p *sqliteSyncState) Add(newState State) error {
	var groupIDBytes []byte
	if newState.GroupID != nil {
		groupIDBytes = newState.GroupID[:]
	}

	_, err := p.db.Exec(`
		INSERT INTO mvds_states 
			(type, send_count, send_epoch, group_id, peer_id, message_id) 
		VALUES 
			(?, ?, ?, ?, ?, ?)`,
		newState.Type,
		newState.SendCount,
		newState.SendEpoch,
		groupIDBytes,
		newState.PeerID[:],
		newState.MessageID[:],
	)
	return err
}

func (p *sqliteSyncState) Remove(messageID MessageID, peerID PeerID) error {
	result, err := p.db.Exec(
		`DELETE FROM mvds_states WHERE message_id = ? AND peer_id = ?`,
		messageID[:],
		peerID[:],
	)
	if err != nil {
		return err
	}
	if n, err := result.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return errors.New("state not found")
	}
	return nil
}

func (p *sqliteSyncState) All() ([]State, error) {
	var result []State

	rows, err := p.db.Query(`
		SELECT 
			type, send_count, send_epoch, group_id, peer_id, message_id 
		FROM
			mvds_states
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			state                      State
			groupID, peerID, messageID []byte
		)
		err := rows.Scan(
			&state.Type,
			&state.SendCount,
			&state.SendEpoch,
			&groupID,
			&peerID,
			&messageID,
		)
		if err != nil {
			return nil, err
		}
		if len(groupID) > 0 {
			val := GroupID{}
			copy(val[:], groupID)
			state.GroupID = &val
		}
		copy(state.PeerID[:], peerID)
		copy(state.MessageID[:], messageID)

		result = append(result, state)
	}

	return result, nil
}

func (p *sqliteSyncState) Map(epoch int64, process func(State) State) error {
	states, err := p.All()
	if err != nil {
		return err
	}

	var updated []State

	for _, state := range states {
		if state.SendEpoch > epoch {
			continue
		}
		newState := process(state)
		if newState != state {
			updated = append(updated, newState)
		}
	}

	if len(updated) == 0 {
		return nil
	}

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	for _, state := range updated {
		if err := updateInTx(tx, state); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func updateInTx(tx *sql.Tx, state State) error {
	_, err := tx.Exec(`
		UPDATE mvds_states
		SET 
			 send_count = ?,
			 send_epoch = ?
		WHERE
			message_id = ? AND 
			peer_id = ?
		`,
		state.SendCount,
		state.SendEpoch,
		state.MessageID[:],
		state.PeerID[:],
	)
	return err
}
