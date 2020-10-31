package exposer

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type Config struct {
	Dsn      string
	Interval time.Duration
}

func Run(conf Config) error {
	db, err := genDB(conf.Dsn)
	if err != nil {
		return errors.Wrap(err, "generate db connection error")
	}

	ticker := time.NewTicker(conf.Interval)
	defer ticker.Stop()

	// シグナルについて(とコンテキストについて)も
	// https://text.baldanders.info/golang/ticker/
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer signal.Stop(sig)

	// 初回起動時
	if err := run(db); err != nil {
		return errors.WithStack(err)
	}
	log.Println("succeeded")

	for {
		select {
		case <-ticker.C:
			if err := run(db); err != nil {
				return errors.WithStack(err)
			}
			log.Println("succeeded")
		case s := <-sig:
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Printf("received signal: %s", s.String())
				return nil
			}
		}
	}

	return nil
}

func run(db *sql.DB) error {
	subjects, err := fetchAllExposedMemoSubject(db)
	if err != nil {
		return errors.Wrap(err, "db error")
	}

	if err := deleteMDs(subjects); err != nil {
		return errors.Wrap(err, "delete md file error")
	}

	memos, err := fetchMemos(db)
	if err != nil {
		return errors.Wrap(err, "db error")
	}

	if len(memos) == 0 {
		return nil
	}

	if err := genMDs(memos); err != nil {
		return errors.Wrap(err, "generate md file error")
	}

	if err := genSite(); err != nil {
		return errors.Wrap(err, "generate html error")
	}

	if err := uploadSite(); err != nil {
		return errors.Wrap(err, "upload site error")
	}

	if err := updateMemosExposedAt(db, memos); err != nil {
		return errors.Wrap(err, "db error")
	}

	return nil
}
