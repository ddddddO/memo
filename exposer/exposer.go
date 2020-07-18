package exposer

import (
	"log"

	"github.com/pkg/errors"
)

func Run(dsn string) error {
	db, err := genDB(dsn)
	if err != nil {
		return errors.Wrap(err, "generate db connection error")
	}

	memos, err := fetchMemos(db)
	if err != nil {
		return errors.Wrap(err, "db error")
	}

	if err := genMDs(memos); err != nil {
		return errors.Wrap(err, "generate md file error")
	}

	if err := genSite(); err != nil {
		return errors.Wrap(err, "generate html error")
	}

	// if err := uploadSite(); err != nil {
	// 	return errors.Wrap(err, "upload site error")
	// }

	log.Println("succeeded")
	return nil
}
