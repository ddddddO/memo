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

	for {
		if err := run(memoRepository); err != nil {
			return errors.WithStack(err)
		}
		log.Println("succeeded")

		select {
		case <-ticker.C:
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

type memoRepository interface {
	FetchList(userID int) ([]*models.Memo, error)
	UpdateExposedAt(memo *models.Memo) error
}

const myUserID = 1

func run(repo memoRepository) error {
	memos, err := repo.FetchList(myUserID)
	if err != nil {
		return errors.WithStack(err)
	}

	dir, err := os.Getwd()
	if err != nil {
		return errors.WithStack(err)
	}
	gce := newGCE(dir)

	subjects := filterExposedSubjects(memos)
	removedMarkdowns, err := gce.removeMarkdwonsNotIncludedInDB(subjects)
	if err != nil {
		return errors.Wrap(err, "remove md file error")
	}

	exposeMemos := filterExposeMemos(memos)
	// 既に同名のmdファイルが存在していた場合、hugo new fuga.mdは失敗する。なので、削除する。
	if err := gce.removeExistingMarkdowns(exposeMemos); err != nil {
		return errors.WithStack(err)
	}
	if err := gce.generateMarkdowns(exposeMemos); err != nil {
		return errors.Wrap(err, "generate md file error")
	}

	if len(exposeMemos) == 0 && len(removedMarkdowns) == 0 {
		return nil
	}

	if err := gce.generateSites(); err != nil {
		return errors.Wrap(err, "generate html error")
	}
	if err := gce.uploadSites(); err != nil {
		return errors.Wrap(err, "upload site error")
	}

	for _, m := range exposeMemos {
		if err := repo.UpdateExposedAt(m); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func filterExposedSubjects(memos []*models.Memo) []string {
	subjects := []string{}
	for _, m := range memos {
		if !m.IsExposed.Valid {
			continue
		}
		if m.IsExposed.Bool {
			subjects = append(subjects, m.Subject)
		}
	}
	return subjects
}

func filterExposeMemos(memos []*models.Memo) []*models.Memo {
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

		// NOTE: すべて消してしまったら以下をコメントイン
		// exposeMemos = append(exposeMemos, m)
		// NOTE: すべて消してしまったら以下をコメントアウト
		if m.IsExposed.Bool && (m.UpdatedAt.Time.After(m.ExposedAt.Time) || m.UpdatedAt.Time.Equal(m.ExposedAt.Time)) {
			exposeMemos = append(exposeMemos, m)
		}
	}
	return exposeMemos
}
