package task

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/spf13/viper"
	"github.com/zaigie/palworld-server-tool/internal/logger"
	"github.com/zaigie/palworld-server-tool/internal/tool"
	"github.com/zaigie/palworld-server-tool/service"
	"go.etcd.io/bbolt"
)

var s gocron.Scheduler

func RconSync(db *bbolt.DB) {
	playersRcon, err := tool.ShowPlayers()
	if err != nil {
		logger.Error(err)
	}
	err = service.PutPlayersRcon(db, playersRcon)
	if err != nil {
		logger.Error(err)
	}
	logger.Info("Rcon sync success\n")
}

func SavSync() {
	err := tool.ConversionLoading(viper.GetString("save.path"))
	if err != nil {
		logger.Error(err)
	}
	logger.Info("Sav sync success\n")
}

func Schedule(db *bbolt.DB) {
	s := getScheduler()

	rconSyncInterval := time.Duration(viper.GetInt("rcon.sync_interval"))
	savSyncInterval := time.Duration(viper.GetInt("save.sync_interval"))

	if rconSyncInterval > 0 {
		RconSync(db)
		_, err := s.NewJob(
			gocron.DurationJob(rconSyncInterval*time.Second),
			gocron.NewTask(RconSync, db),
		)
		if err != nil {
			logger.Error(err)
		}
	}

	if savSyncInterval > 0 {
		SavSync()
		_, err := s.NewJob(
			gocron.DurationJob(savSyncInterval*time.Second),
			gocron.NewTask(SavSync),
		)
		if err != nil {
			logger.Error(err)
		}
	}

	s.Start()
}

func Shutdown() {
	s := getScheduler()
	err := s.Shutdown()
	if err != nil {
		logger.Error(err)
	}
}

func initScheduler() gocron.Scheduler {
	s, err := gocron.NewScheduler()
	if err != nil {
		logger.Error(err)
	}
	return s
}

func getScheduler() gocron.Scheduler {
	if s == nil {
		return initScheduler()
	}
	return s
}