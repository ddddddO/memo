package postgres

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/ddddddO/memo/adapter"
)

type memoRepository struct {
	db *sql.DB
}

func NewMemoRepository(db *sql.DB) *memoRepository {
	return &memoRepository{
		db: db,
	}
}

func (pg *memoRepository) FetchList(userID int, tagID int) ([]adapter.Memo, error) {
	var (
		rows  *sql.Rows
		memos []adapter.Memo
		err   error
	)
	// NOTE: tagIdが設定されていない場合
	if tagID == -1 {
		query := "SELECT id, subject, notified_cnt FROM memos WHERE users_id=$1 ORDER BY id"
		rows, err = pg.db.Query(query, userID)
		if err != nil {
			return nil, err
		}
	} else {
		// NOTE: tagIdが設定されている場合
		query := "SELECT id, subject, notified_cnt FROM memos WHERE users_id=$1 AND id IN (SELECT memos_id FROM memo_tag WHERE tags_id=$2) ORDER BY id"
		rows, err = pg.db.Query(query, userID, tagID)
		if err != nil {
			return nil, err
		}
	}

	for rows.Next() {
		var memo adapter.Memo
		if err := rows.Scan(&memo.ID, &memo.Subject, &memo.NotifiedCnt); err != nil {
			return nil, err
		}
		setColor(&memo)
		memos = append(memos, memo)
	}
	// NOTE: NotifiedCntでメモを昇順にソート
	sort.SliceStable(memos,
		func(i, j int) bool {
			return memos[i].NotifiedCnt < memos[j].NotifiedCnt
		},
	)

	return memos, nil
}

// TODO: repositoryはデータの永続化・復元が責務だから、この変換はここではない。
func setColor(m *adapter.Memo) {
	switch m.NotifiedCnt {
	case 0:
		m.RowVariant = "danger"
	case 1:
		m.RowVariant = "warning"
	case 2:
		m.RowVariant = "primary"
	case 3:
		m.RowVariant = "info"
	case 4:
		m.RowVariant = "secondary"
	case 5:
		m.RowVariant = "success"
	}
}

func (pg *memoRepository) Fetch(userID int, memoID int) (adapter.Memo, error) {
	// TODO: 見直す
	const memoDetailQuery = `
	SELECT
	    m.id AS id,
	    m.subject AS subject,
		m.content AS content,
		m.is_exposed AS is_exposed,
		(SELECT jsonb_agg(DISTINCT(t.id))
		  FROM memos m
		  JOIN memo_tag mt
		  ON m.id = mt.memos_id
		  JOIN tags t
		  ON mt.tags_id = t.id
	      WHERE m.id = $1 AND m.users_id = $2
		) AS tag_ids,
		(SELECT jsonb_agg(DISTINCT(t.name))
		  FROM memos m
		  JOIN memo_tag mt
		  ON m.id = mt.memos_id
		  JOIN tags t
		  ON mt.tags_id = t.id
	      WHERE m.id = $1 AND m.users_id = $2
		) AS tag_names
		   FROM memos m
		   JOIN memo_tag mt
		   ON m.id = mt.memos_id
		   JOIN tags t
		   ON mt.tags_id = t.id
		WHERE m.id = $1 AND m.users_id = $2
		GROUP BY m.id
	`

	var (
		memo     adapter.Memo
		tagIDs   string
		tagNames string
	)
	if err := pg.db.QueryRow(memoDetailQuery, memoID, userID).Scan(
		&memo.ID, &memo.Subject, &memo.Content, &memo.IsExposed,
		&tagIDs, &tagNames,
	); err != nil {
		return adapter.Memo{}, err
	}

	memo.Tags = toTags(tagIDs, tagNames)

	return memo, nil
}

// TODO: 以下３つの変換用関数もここではない
func toTags(ids, names string) []adapter.Tag {
	convertedIDs := toInts(ids)
	convertedNames := toStrings(names)

	var tags []adapter.Tag
	for i := range convertedIDs {
		tag := adapter.Tag{}
		tag.ID = convertedIDs[i]
		tag.Name = convertedNames[i]

		tags = append(tags, tag)
	}
	return tags
}

func toInts(s string) []int {
	if len(s) <= 2 {
		return []int{}
	}

	ss := strings.Split(s[1:len(s)-1], ",")
	var nums []int
	for _, strNum := range ss {
		strNum = strings.TrimSpace(strNum)
		num, err := strconv.Atoi(strNum)
		if err != nil {
			panic(err)
		}
		nums = append(nums, num)
	}
	return nums
}

func toStrings(s string) []string {
	if len(s) <= 2 {
		return []string{}
	}
	return strings.Split(s[1:len(s)-1], ",")
}

func (pg *memoRepository) Update(memo adapter.Memo) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	const updateMemoQuery = `
	UPDATE memos SET subject=$1, content=$2, is_exposed=$3
	 WHERE id=$4 AND users_id=$5
	`

	result, err := tx.Exec(updateMemoQuery,
		memo.Subject, memo.Content, memo.IsExposed, memo.ID, memo.UserID,
	)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("unexpected")
	}

	const deleteMemoTagQuery = `
	DELETE FROM memo_tag WHERE memos_id=$1 AND tags_id <> 1
	`

	_, err = tx.Exec(deleteMemoTagQuery, memo.ID)
	if err != nil {
		return err
	}

	var insertMemoTagQuery = `
	INSERT INTO memo_tag(memos_id, tags_id) VALUES
	%s
	`

	var values string
	for _, tag := range memo.Tags {
		values += fmt.Sprintf("(%v, %d),", memo.ID, tag.ID)
	}
	insertMemoTagQuery = fmt.Sprintf(insertMemoTagQuery, values[:len(values)-1])

	_, err = tx.Exec(insertMemoTagQuery)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (pg *memoRepository) Create(memo adapter.Memo) error {
	var createMemoQuery = `
	WITH inserted AS (INSERT INTO memos(subject, content, users_id, is_exposed) VALUES($1, $2, $3, $4) RETURNING id)
	INSERT INTO memo_tag(memos_id, tags_id) VALUES
		((SELECT id FROM inserted), 1)
		%s;
	`

	var values string
	for _, tag := range memo.Tags {
		values += fmt.Sprintf(",((SELECT id FROM inserted), %d)", tag.ID)
	}
	createMemoQuery = fmt.Sprintf(createMemoQuery, values)

	_, err := pg.db.Exec(createMemoQuery,
		memo.Subject, memo.Content, memo.UserID, memo.IsExposed,
	)
	if err != nil {
		return err
	}

	return nil
}

func (pg *memoRepository) Delete(memo adapter.Memo) error {
	const deleteMemoQuery = `
	DELETE FROM memos WHERE users_id = $1 AND id = $2;
	`

	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec(deleteMemoQuery,
		memo.UserID, memo.ID,
	)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("unexpected")
	}

	const deleteMemoTagQuery = `
	DELETE FROM memo_tag WHERE memos_id = $1
	`

	_, err = tx.Exec(deleteMemoTagQuery, memo.ID)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}
