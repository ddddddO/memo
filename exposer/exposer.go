package exposer

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"

	"github.com/ddddddO/memo/models"
	"github.com/ddddddO/memo/repository"
	"github.com/ddddddO/memo/repository/postgres"
)

type Config struct {
	Dsn      string
	Interval time.Duration
}

func Run(conf Config) error {
	db, err := sql.Open("postgres", conf.Dsn)
	if err != nil {
		return err
	}
	memoRepository := postgres.NewMemoRepository(db)

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
	if err := run(memoRepository); err != nil {
		return errors.WithStack(err)
	}
	log.Println("succeeded")

	for {
		select {
		case <-ticker.C:
			if err := run(memoRepository); err != nil {
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

const myUserID = 1

// TODO: ここで使うメソッドとしては過剰なinterfaceだから、絞ってもいいかも
func run(repo repository.MemoRepository) error {
	memos, err := repo.FetchList(myUserID)
	if err != nil {
		return errors.WithStack(err)
	}

	subjects := []string{}
	for _, m := range memos {
		if !m.IsExposed.Valid {
			continue
		}
		if m.IsExposed.Bool {
			subjects = append(subjects, m.Subject)
		}
	}

	removedMarkdowns, err := removeMarkdwonsNotIncluded(subjects)
	if err != nil {
		return errors.Wrap(err, "remove md file error")
	}

	// 念のため。。
	time.Sleep(3 * time.Second)

	var exposeMemos []*models.Memo
	for _, m := range memos {
		if !m.IsExposed.Valid {
			continue
		}

		if m.IsExposed.Bool && !m.ExposedAt.Valid {
			exposeMemos = append(exposeMemos, m)
			continue
		}

		if !m.ExposedAt.Valid {
			continue
		}
		if !m.UpdatedAt.Valid {
			continue
		}

		if m.IsExposed.Bool && (m.UpdatedAt.Time.After(m.ExposedAt.Time)) {
			exposeMemos = append(exposeMemos, m)
		}
	}

	if err := generateMarkdowns(exposeMemos); err != nil {
		return errors.Wrap(err, "generate md file error")
	}

	if len(exposeMemos) == 0 && len(removedMarkdowns) == 0 {
		return nil
	}

	if err := generateSites(); err != nil {
		return errors.Wrap(err, "generate html error")
	}

	if err := uploadSites(); err != nil {
		return errors.Wrap(err, "upload site error")
	}

	for _, m := range exposeMemos {
		if err := repo.UpdateExposedAt(m); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
