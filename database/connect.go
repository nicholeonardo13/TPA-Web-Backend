package database


import(
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const DB_NAME="staem"
const DB_HOST="127.0.0.1"
const DB_PORT="5432"				//port on installation
const DB_USER="postgres"			//default is postgres
const DB_PASS="vi1819"	//password on installation

func Connect()(*gorm.DB, error) {

	dsn := "host=127.0.0.1 user=postgres password=vi1819 dbname=staem port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}