package pkg

import (
	"fmt"
	"testing"
)

type TaskInfo struct {
}

func TestMysql(t *testing.T) {
	version := GetMySQLVersion()
	fmt.Printf("Mysql version: %v", version)

	conf := getMysqlConfig()
	if conf == nil {
		t.Error("Expected MysqlConfig, got nil")
	}
	fmt.Printf("conf: %v\n", conf)
	UpdateMysqlConfig(MysqlConfig{
		MYSQL_DB: "test_db",
	})
	db, _ := GetDB()
	_ = db

}
