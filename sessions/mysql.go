package sessions

import (
	"fmt"
	"time"

	"github.com/spf13/viper"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // Need this to connect to MySQL
)

var (
	sqlClient *gorm.DB
)

func newSqlClient() *gorm.DB {
	sqlClient, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true", viper.GetString("sql.username"),
		viper.GetString("sql.password"), viper.GetString("sql.host"), viper.GetInt64("sql.port"), viper.GetString("sql.dbname")))
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	sqlClient.DB().SetConnMaxLifetime(time.Duration(viper.GetInt64("sql.connection_lifetime")) * time.Second)
	sqlClient.DB().SetMaxOpenConns(80)
	sqlClient.DB().SetMaxIdleConns(20)
	return sqlClient
}

// InitSqlClient initializes all common components
func InitSqlClient() {
	if sqlClient == nil {
		sqlClient = newSqlClient()
	}
}

// GetSqlClient returns the instance of sqlClient that have
// already been initialized through InitSqlClient
func GetSqlClient() *gorm.DB {
	InitSqlClient()
	return sqlClient
}
