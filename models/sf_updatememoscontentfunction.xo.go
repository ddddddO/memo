package models

// Code generated by xo. DO NOT EDIT.

import (
	"context"
)

// UpdateMemosContentFunction calls the stored function 'public.update_memos_content_function() trigger' on db.
func UpdateMemosContentFunction(ctx context.Context, db DB) (Trigger, error) {
	// call public.update_memos_content_function
	const sqlstr = `SELECT * FROM public.update_memos_content_function()`
	// run
	var r0 Trigger
	logf(sqlstr)
	if err := db.QueryRowContext(ctx, sqlstr).Scan(&r0); err != nil {
		return Trigger{}, logerror(err)
	}
	return r0, nil
}