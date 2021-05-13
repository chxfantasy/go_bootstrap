package conf

import (
	"flag"
	"fmt"
	"log"

	"github.com/chxfantasy/go_bootstrap/persist/mongo"
	"github.com/chxfantasy/go_bootstrap/persist/redis"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
)

var (
	Gin         *gin.Engine
	Redis1      *redis.Pool
	MongoTest   *mgo.Session
	AppConf     *AppConfig
	BizLogger   *Logger
	TraceLogger *Logger
)

func init() {
	var env string
	flag.StringVar(&env, "env", "dev", "set env")
	flag.Parse()
	//初始化config
	var err error
	AppConf, err = LoadConfig(env)
	if err != nil || AppConf == nil {
		log.Fatalf("load config fail, error:%s, env:%s \n", err, env)
	}

	//初始化Logger
	BizLogger = NewLogger(AppConf.BizLoggerConf)
	TraceLogger = NewLogger(AppConf.TraceLoggerConf)

	// init Redis
	Redis1, err = redis.NewRedisPool(AppConf.Redis1Conf)
	if err != nil || Redis1 == nil {
		log.Fatalf("init Redis1 fail, error:%s, env:%s \n", err, env)
	}
	//init mongo
	MongoTest, err = mongo.NewMongoSession(AppConf.MongoTestConf)
	if err != nil || MongoTest == nil {
		log.Fatalf("init MongoTest fail, error:%s, env:%s \n", err, env)
	}
	fmt.Println("init app...")

	Gin = gin.New()
	gin.SetMode(AppConf.Server.Env)
	Gin.Use(gin.Recovery())
}
