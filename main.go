package main

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"system-service-template/config"
	"system-service-template/database"
	"system-service-template/utils"
	"system-service-template/web"

	_ "github.com/alexbrainman/odbc"
	_ "github.com/axgle/mahonia"
	_ "github.com/fsnotify/fsnotify"
	_ "github.com/gin-contrib/cache"
	_ "github.com/gin-contrib/cors"
	_ "github.com/gin-contrib/logger"
	_ "github.com/gin-contrib/pprof"
	_ "github.com/gin-contrib/requestid"
	_ "github.com/gin-contrib/sessions"
	_ "github.com/gin-contrib/static"
	_ "github.com/gin-contrib/zap"
	_ "github.com/gin-gonic/gin"
	_ "github.com/glebarez/sqlite"
	_ "github.com/go-ole/go-ole"
	_ "github.com/julienschmidt/httprouter"
	"github.com/kardianos/service"
	_ "github.com/minio/minio-go/v7"
	_ "github.com/patrickmn/go-cache"
	_ "github.com/pkg/sftp"
	_ "github.com/radovskyb/watcher"
	_ "github.com/robertkrimen/otto"
	_ "github.com/robfig/cron/v3"
	_ "github.com/rs/zerolog"
	_ "github.com/sijms/go-ora"
	log "github.com/sirupsen/logrus"
	_ "github.com/spf13/viper"
	_ "github.com/thoas/go-funk"
	_ "go.etcd.io/bbolt"
	"gopkg.in/natefinch/lumberjack.v2"
	_ "gorm.io/driver/clickhouse"
	_ "gorm.io/driver/mysql"
	_ "gorm.io/driver/postgres"
	_ "gorm.io/driver/sqlserver"
	_ "gorm.io/gorm"
)

var (
	logFilePath string
	gConfig     config.Config
)

type program struct {
	logger service.Logger
	cfg    *service.Config
	srv    *http.Server
}

func (p *program) Start(s service.Service) error {
	if p.logger != nil {
		p.logger.Info("start service")
	}
	go p.run()
	return nil
}
func (p *program) Stop(s service.Service) error {
	if p.logger != nil {
		p.logger.Info("stop service")
	}
	p.srv.Shutdown(context.Background())
	return nil
}

func (p *program) run() {
	engine := web.New()
	p.srv = &http.Server{
		Addr:    ":" + gConfig.HttpPort,
		Handler: engine,
	}
	p.srv.ListenAndServe()
}

func main() {
	// setup config
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
		return
	}

	gConfig = config.GetConfig()

	// setup log
	logFilePath = gConfig.LogName
	if logFilePath == "" {
		logFilePath = "app.log"
	}
	if !filepath.IsAbs(logFilePath) {
		logFilePath = utils.GetAbsPath(logFilePath)
	}
	logger := &lumberjack.Logger{
		Filename: logFilePath,
		MaxSize:  gConfig.LogMaxSize,
		MaxAge:   gConfig.LogMaxAge,
	}
	log.SetOutput(logger)
	// setup database
	database.Init()

	pro := &program{
		cfg: &service.Config{
			Name:        gConfig.ServiceName,
			DisplayName: gConfig.ServiceDisplayName,
			Description: gConfig.ServiceDescription,
			Option: service.KeyValue{
				"DelayedAutoStart": true,
			},
		},
	}
	sys := service.ChosenSystem()
	srv, err := sys.New(pro, pro.cfg)
	if err != nil {
		log.Fatal(err)
		return
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			if err := srv.Install(); err != nil {
				log.Println("install service error:", err)
			}
		case "uninstall":
			if err := srv.Uninstall(); err != nil {
				log.Println("uninstall service error:", err)
			}
		case "stop":
			if err := srv.Stop(); err != nil {
				log.Println("stop service error:", err)
			}
		}
		return
	}

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
