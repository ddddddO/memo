package exposer

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"

	"github.com/ddddddO/tag-mng/exposer/datasource"
)

type Config struct {
	Dsn      string
	Interval time.Duration
}

func Run(conf Config) error {
	postgres, err := datasource.NewPostgres(conf.Dsn)
	if err != nil {
		return errors.Wrap(err, "db connection error")
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
	if err := run(postgres); err != nil {
		return errors.WithStack(err)
	}
	log.Println("succeeded")

	for {
		select {
		case <-ticker.C:
			if err := run(postgres); err != nil {
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

func run(ds datasource.DataSource) error {
	subjects, err := ds.FetchAllExposedMemoSubjects()
	if err != nil {
		return errors.Wrap(err, "db error")
	}

	removedMarkdowns, err := removeMarkdwonsNotIncluded(subjects)
	if err != nil {
		return errors.Wrap(err, "remove md file error")
	}

	// 念のため。。
	time.Sleep(3 * time.Second)

	memos, err := ds.FetchMemos()
	if err != nil {
		return errors.Wrap(err, "db error")
	}

	if err := generateMarkdowns(memos); err != nil {
		return errors.Wrap(err, "generate md file error")
	}

	if len(memos) == 0 && len(removedMarkdowns) == 0 {
		return nil
	}

	if err := generateSites(); err != nil {
		return errors.Wrap(err, "generate html error")
	}

	if err := uploadSites(); err != nil {
		return errors.Wrap(err, "upload site error")
	}

	if err := ds.UpdateMemosExposedAt(memos); err != nil {
		return errors.Wrap(err, "db error")
	}

	return nil
}
