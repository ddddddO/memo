package models

// Code generated by xo. DO NOT EDIT.

import (
	"context"
	"database/sql"
)

// Memo represents a row from 'public.memos'.
type Memo struct {
	ID          int           `json:"id"`           // id
	Subject     string        `json:"subject"`      // subject
	Content     string        `json:"content"`      // content
	UsersID     sql.NullInt64 `json:"users_id"`     // users_id
	CreatedAt   sql.NullTime  `json:"created_at"`   // created_at
	UpdatedAt   sql.NullTime  `json:"updated_at"`   // updated_at
	NotifiedCnt sql.NullInt64 `json:"notified_cnt"` // notified_cnt
	IsExposed   sql.NullBool  `json:"is_exposed"`   // is_exposed
	ExposedAt   sql.NullTime  `json:"exposed_at"`   // exposed_at
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the Memo exists in the database.
func (m *Memo) Exists() bool {
	return m._exists
}

// Deleted returns true when the Memo has been marked for deletion from
// the database.
func (m *Memo) Deleted() bool {
	return m._deleted
}

// Insert inserts the Memo to the database.
func (m *Memo) Insert(ctx context.Context, db DB) error {
	switch {
	case m._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case m._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (manual)
	const sqlstr = `INSERT INTO public.memos (` +
		`id, subject, content, users_id, created_at, updated_at, notified_cnt, is_exposed, exposed_at` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9` +
		`)`
	// run
	logf(sqlstr, m.ID, m.Subject, m.Content, m.UsersID, m.CreatedAt, m.UpdatedAt, m.NotifiedCnt, m.IsExposed, m.ExposedAt)
	if _, err := db.ExecContext(ctx, sqlstr, m.ID, m.Subject, m.Content, m.UsersID, m.CreatedAt, m.UpdatedAt, m.NotifiedCnt, m.IsExposed, m.ExposedAt); err != nil {
		return logerror(err)
	}
	// set exists
	m._exists = true
	return nil
}

// Update updates a Memo in the database.
func (m *Memo) Update(ctx context.Context, db DB) error {
	switch {
	case !m._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case m._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with composite primary key
	const sqlstr = `UPDATE public.memos SET ` +
		`subject = $1, content = $2, users_id = $3, created_at = $4, updated_at = $5, notified_cnt = $6, is_exposed = $7, exposed_at = $8 ` +
		`WHERE id = $9`
	// run
	logf(sqlstr, m.Subject, m.Content, m.UsersID, m.CreatedAt, m.UpdatedAt, m.NotifiedCnt, m.IsExposed, m.ExposedAt, m.ID)
	if _, err := db.ExecContext(ctx, sqlstr, m.Subject, m.Content, m.UsersID, m.CreatedAt, m.UpdatedAt, m.NotifiedCnt, m.IsExposed, m.ExposedAt, m.ID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the Memo to the database.
func (m *Memo) Save(ctx context.Context, db DB) error {
	if m.Exists() {
		return m.Update(ctx, db)
	}
	return m.Insert(ctx, db)
}

// Upsert performs an upsert for Memo.
func (m *Memo) Upsert(ctx context.Context, db DB) error {
	switch {
	case m._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO public.memos (` +
		`id, subject, content, users_id, created_at, updated_at, notified_cnt, is_exposed, exposed_at` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9` +
		`)` +
		` ON CONFLICT (id) DO ` +
		`UPDATE SET ` +
		`subject = EXCLUDED.subject, content = EXCLUDED.content, users_id = EXCLUDED.users_id, created_at = EXCLUDED.created_at, updated_at = EXCLUDED.updated_at, notified_cnt = EXCLUDED.notified_cnt, is_exposed = EXCLUDED.is_exposed, exposed_at = EXCLUDED.exposed_at `
	// run
	logf(sqlstr, m.ID, m.Subject, m.Content, m.UsersID, m.CreatedAt, m.UpdatedAt, m.NotifiedCnt, m.IsExposed, m.ExposedAt)
	if _, err := db.ExecContext(ctx, sqlstr, m.ID, m.Subject, m.Content, m.UsersID, m.CreatedAt, m.UpdatedAt, m.NotifiedCnt, m.IsExposed, m.ExposedAt); err != nil {
		return logerror(err)
	}
	// set exists
	m._exists = true
	return nil
}

// Delete deletes the Memo from the database.
func (m *Memo) Delete(ctx context.Context, db DB) error {
	switch {
	case !m._exists: // doesn't exist
		return nil
	case m._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM public.memos ` +
		`WHERE id = $1`
	// run
	logf(sqlstr, m.ID)
	if _, err := db.ExecContext(ctx, sqlstr, m.ID); err != nil {
		return logerror(err)
	}
	// set deleted
	m._deleted = true
	return nil
}

// MemoByID retrieves a row from 'public.memos' as a Memo.
//
// Generated from index 'memos_pkey'.
func MemoByID(ctx context.Context, db DB, id int) (*Memo, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, subject, content, users_id, created_at, updated_at, notified_cnt, is_exposed, exposed_at ` +
		`FROM public.memos ` +
		`WHERE id = $1`
	// run
	logf(sqlstr, id)
	m := Memo{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, id).Scan(&m.ID, &m.Subject, &m.Content, &m.UsersID, &m.CreatedAt, &m.UpdatedAt, &m.NotifiedCnt, &m.IsExposed, &m.ExposedAt); err != nil {
		return nil, logerror(err)
	}
	return &m, nil
}

// MemosByUsersID retrieves a row from 'public.memos' as a Memo.
//
// Generated from index 'memos_users_id_idx'.
func MemosByUsersID(ctx context.Context, db DB, usersID sql.NullInt64) ([]*Memo, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, subject, content, users_id, created_at, updated_at, notified_cnt, is_exposed, exposed_at ` +
		`FROM public.memos ` +
		`WHERE users_id = $1`
	// run
	logf(sqlstr, usersID)
	rows, err := db.QueryContext(ctx, sqlstr, usersID)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// process
	var res []*Memo
	for rows.Next() {
		m := Memo{
			_exists: true,
		}
		// scan
		if err := rows.Scan(&m.ID, &m.Subject, &m.Content, &m.UsersID, &m.CreatedAt, &m.UpdatedAt, &m.NotifiedCnt, &m.IsExposed, &m.ExposedAt); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &m)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}

// User returns the User associated with the Memo's (UsersID).
//
// Generated from foreign key 'memos_users_id_fkey'.
func (m *Memo) User(ctx context.Context, db DB) (*User, error) {
	return UserByID(ctx, db, int(m.UsersID.Int64))
}
