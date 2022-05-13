package main

import (
	"context"
	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
	app "github.com/sku4/alice-checklist"
	"github.com/sku4/alice-checklist/configs"
	"github.com/sku4/alice-checklist/internal/handler"
	"github.com/sku4/alice-checklist/internal/repository"
	"github.com/sku4/alice-checklist/internal/service"
	"github.com/sku4/alice-checklist/lang"
	"github.com/sku4/alice-checklist/pkg/boltdb"
	"os"
	"os/signal"
	"syscall"
)

// @title Alice webhook app API
// @version 1.0
// @description API Server for Alice checklist application

// @host localhost:8000
// @BasePath /

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	loc, err := lang.InitLocalize(lang.Ru)
	if err != nil {
		logrus.Fatal(err)
	}

	cfg, err := configs.Init()
	if err != nil {
		logrus.Fatalf(loc.Translate("error init config: %s"), err.Error())
	}

	db, err := initDB(cfg)
	if err != nil {
		logrus.Fatal(err)
	}

	repos := repository.NewRepository(loc, cfg, db)
	services := service.NewService(loc, repos)
	handlers := handler.NewHandler(loc, services)

	srv := new(app.Server)
	go func() {
		if err := srv.Run(cfg.Port, handlers.InitRoutes()); err != nil {
			logrus.Fatalf(loc.Translate("error occured while running http server: %s"), err.Error())
		}
	}()

	logrus.Print(loc.Translate("Checklist App Started"))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print(loc.Translate("Checklist App Shutting Down"))

	if err := db.Close(); err != nil {
		logrus.Errorf(loc.Translate("error occured on db connection close: %s"), err.Error())
	}

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf(loc.Translate("error occured on server shutting down: %s"), err.Error())
	}
}

func initDB(cfg *configs.Config) (boltdb.Storage, error) {
	db, err := bolt.Open(cfg.DBPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		for _, b := range boltdb.BucketsNote {
			_, err = tx.CreateBucketIfNotExists([]byte(b))
			if err != nil {
				return err
			}
		}
		for _, b := range boltdb.BucketsNotify {
			_, err = tx.CreateBucketIfNotExists([]byte(b))
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil
}
