package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
	gormv1 "github.com/jinzhu/gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
)

// Base ...
type Base struct {
	ID       uuid.UUID      `gorm:"type:uuid;primary_key"`
	CreateAt int64          `gorm:"autoCreateTime"`
	UpdateAt int64          `gorm:"autoUpdateTime"`
	DeleteAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(scope *gorm.DB) (err error) {
	id, err := uuid.NewV4()
	if err != nil {
		return
	}

	scope.Statement.SetColumn("ID", id)
	return
}

var db *gorm.DB
var dbv1 *gormv1.DB
var rediscache *redis.Client

func init() {

	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DB := os.Getenv("DB_DRIVER")
	HostName := os.Getenv("DB_HOST")

	DBURL := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", DbUser, DbPassword, DbName, HostName)
	fmt.Println(DBURL)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)

	conn, err := gorm.Open(postgres.Open(DBURL), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		fmt.Print(err)
	}

	dbv1, err = gormv1.Open(DB, DBURL)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	err = db.Debug().AutoMigrate(
		&User{},
		&Partner{},
		&Post{},
		&Image{},
		&OldJob{}, &Skill{}, &Language{}, &Education{},
		&Company{},
		&Admin{},
		&TimeLine{},
	)

	if err != nil {
		log.Fatal(err)
	}

	generateLocation()

	// redishost := os.Getenv("REDIS_HOST")
	// redisport := os.Getenv("REDIS_PORT")
	// redispassword := os.Getenv("REDIS_PASSWORD")

	// ADDRESS := fmt.Sprintf("%s:%s", redishost, redisport)
	// fmt.Println("Redis-Server: [ " + ADDRESS + " ]" + "[ password=" + redispassword + " ]")
	// cache := redis.NewClient(&redis.Options{
	// 	Addr:     ADDRESS,
	// 	Password: redispassword,
	// 	DB:       0,
	// })

	// rediscache = cache
	// _, err = rediscache.Ping().Result()
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

// GetDB ...
func GetDB() *gorm.DB {
	return db
}

// GetDBV1 ...
func GetDBV1() *gormv1.DB {
	return dbv1
}

// // GetCache ... for Redis
// func GetCache() *redis.Client {
// 	return rediscache
// }

func generateLocation() {
	dbv1 = GetDBV1()

	//Fetch Locations Data
	var Country LocationList
	data := Country.GernerateCoordinate()

	// Check if Table Location is existed
	tableExist := dbv1.HasTable(&Location{})
	if !tableExist {
		err := dbv1.Debug().AutoMigrate(&Location{}).Error
		if err != nil {
			log.Fatal(err)
		}

		err = Country.InstallLocaion(dbv1, data)
		if err != nil {
			panic(err)
		}
	}

	// Check if Table Locations is Not as Expected
	var counter int
	dbv1.Model(&Location{}).Count(&counter)
	if counter != len(data) {
		fmt.Println("Location is not Valid... Generating")

		err := dbv1.Debug().DropTable(&Location{}).Error
		if err != nil {
			log.Fatalf("cannot drop table: %v", err)
		}

		err = dbv1.Debug().AutoMigrate(&Location{}).Error
		if err != nil {
			log.Fatal(err)
		}

		err = Country.InstallLocaion(dbv1, data)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Everything is OK ...")
}
