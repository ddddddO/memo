package models

// Code generated by xo. DO NOT EDIT.

import (
	"context"
)

// UpdateMemosSubjectFunction calls the stored function 'public.update_memos_subject_function() trigger' on db.
func UpdateMemosSubjectFunction(ctx context.Context, db DB) (Trigger, error) {
	// call public.update_memos_subject_function
	const sqlstr = `SELECT * FROM public.update_memos_subject_function()`
	// run
	var r0 Trigger
	logf(sqlstr)
	if err := db.QueryRowContext(ctx, sqlstr).Scan(&r0); err != nil {
		return Trigger{}, logerror(err)
	}
	return r0, nil
}