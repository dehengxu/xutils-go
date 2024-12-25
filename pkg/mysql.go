package pkg

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MysqlConfig struct {
	MYSQL_HOST     string
	MYSQL_PASSWORD string
	MYSQL_USER     string
	MYSQL_DB       string
	//DSN            string
	// REDIS_HOST     string
	// REDIS_PASSWORD string
	// REDIS_DB       int

	IS_DEBUG bool
}

var _dsn string = ""

func (conf *MysqlConfig) GetDSN() string {
	confWRLock.RLock()
	defer confWRLock.RUnlock()
	if mysql_config_changed {
		_dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", conf.MYSQL_USER, conf.MYSQL_PASSWORD, conf.MYSQL_HOST, conf.MYSQL_DB)
		mysql_config_changed = false
	}
	return _dsn
}

var mysql_config_once sync.Once
var mysql_config *MysqlConfig
var confWRLock sync.RWMutex
var mysql_config_changed = true

func getMysqlConfig() *MysqlConfig {
	confWRLock.RLock()
	defer confWRLock.RUnlock()

	mysql_config_once.Do(func() {
		conf := &MysqlConfig{}
		conf.MYSQL_HOST = GetEnv("MYSQL_HOST", "localhost:3306")
		conf.MYSQL_PASSWORD = GetEnv("MYSQL_PASSWORD", "root")
		conf.MYSQL_USER = GetEnv("MYSQL_USER", "root")
		conf.MYSQL_DB = GetEnv("MYSQL_DB", "pcs_db")
		conf.IS_DEBUG = strings.ToLower(GetEnv("DEBUG", "false")) == "true"
		mysql_config = conf
	})
	return mysql_config
}

func UpdateMysqlConfig(conf MysqlConfig) {
	confWRLock.Lock()
	defer confWRLock.Unlock()

	mysql_config.IS_DEBUG = conf.IS_DEBUG

	if len(conf.MYSQL_DB) > 0 {
		mysql_config_changed = true
		mysql_config.MYSQL_DB = conf.MYSQL_DB
	}
	if len(conf.MYSQL_HOST) > 0 {
		mysql_config_changed = true
		mysql_config.MYSQL_HOST = conf.MYSQL_HOST
	}
	if len(conf.MYSQL_PASSWORD) > 0 {
		mysql_config_changed = true
		mysql_config.MYSQL_PASSWORD = conf.MYSQL_PASSWORD
	}
	if len(conf.MYSQL_USER) > 0 {
		mysql_config_changed = true
		mysql_config.MYSQL_USER = conf.MYSQL_USER
	}
}

func init() {
	_ = getMysqlConfig()
	//autoCreateDB(conf)
}

var db_once sync.Once
var _db *gorm.DB

func GetDB() (*gorm.DB, error) {
	db_once.Do(func() {
		conf := getMysqlConfig()

		//db, err := gorm.Open(mysql.Open(conf.DSN), &gorm.Config{})
		db, err := OpenDB(conf)

		if err != nil {
			db, err = autoCreateDB(conf)
		}

		fmt.Printf("GetDB: %T\n", db)

		if err != nil {
			fmt.Printf("sql.Open failed: %v\n", err)
			_db = nil
			return
		}
		_db = db
	})

	return _db, nil
}

func OpenDB(conf *MysqlConfig) (*gorm.DB, error) {
	gormConfig := gorm.Config{}
	if conf.IS_DEBUG {
		gormConfig = gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		}
	}

	db, err := gorm.Open(mysql.Open(conf.GetDSN()), &gormConfig)

	if err != nil {
		fmt.Printf("sql.Open failed: %v, try to create db :%v\n", err, conf.MYSQL_DB)
		return nil, err
	}
	return db, nil
}

func autoCreateDB(conf *MysqlConfig) (*gorm.DB, error) {
	fmt.Printf("autoCreateDB: %v\n", conf)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/", conf.MYSQL_USER, conf.MYSQL_PASSWORD, conf.MYSQL_HOST)
	if sqldb, sqlerr := sql.Open("mysql", dsn); sqlerr == nil {
		fmt.Printf("sql try to create db: %v\n", conf.MYSQL_DB)
		_ = fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", conf.MYSQL_DB)
		sqlstr := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", conf.MYSQL_DB)
		if _, err_create := sqldb.Exec(sqlstr); err_create != nil {
			fmt.Printf("CREATE DATABASE failed: %v", err_create)
			return nil, err_create
		} else {
			gormConfig := gorm.Config{}
			if conf.IS_DEBUG {
				gormConfig = gorm.Config{
					Logger: logger.Default.LogMode(logger.Info),
				}
			}
			db, err := gorm.Open(mysql.Open(conf.GetDSN()), &gormConfig)
			return db, err
		}
	} else {
		fmt.Printf("auto create db failed: %v\n", sqlerr)
	}
	return nil, fmt.Errorf("sql.Open failed")
}
