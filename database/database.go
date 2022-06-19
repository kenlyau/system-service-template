package database

import (
	"os"
	"path/filepath"
	"system-service-template/config"
	"system-service-template/utils"

	"github.com/glebarez/sqlite"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var dialector gorm.Dialector
	var err error
	databaseConf := config.GetConfig().Database
	dbType, _ := databaseConf["type"]

	switch dbType {
	case "mysql":
	default:
		dbPath, found := databaseConf["path"]
		if !found {
			dbPath = "app.db"
		}
		if !filepath.IsAbs(dbPath) {
			dbPath = utils.GetAbsPath(dbPath)
		}
		_, err = os.Stat(dbPath)
		if os.IsNotExist(err) {
			f, _ := os.Create(dbPath)
			f.Close()
		}
		dialector = sqlite.Open(dbPath)
	}

	DB, err = gorm.Open(dialector)
	if err != nil {
		log.Println("init database error:", err)
	}

	//auto migrate
	DB.AutoMigrate(&User{})
}
