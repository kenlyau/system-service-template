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

	"github.com/kardianos/service"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
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
