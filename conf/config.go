package conf

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/chxfantasy/go_bootstrap/persist/mongo"
	"github.com/chxfantasy/go_bootstrap/persist/redis"

	"github.com/ghodss/yaml"
)

// AppConfig appConfig
type AppConfig struct {
	Server          serverConfigDef  `json:"server" yaml:"server" mapstructure:"server"`
	Redis1Conf      *redis.ConfigDef `json:"redis-1" yaml:"redis-1" mapstructure:"redis-1"`
	MongoTestConf   *mongo.ConfigDef `json:"mongo_test" yaml:"mongo_test" mapstructure:"mongo_test"`
	BizLoggerConf   ConfigDef        `json:"biz_log" yaml:"biz_log" mapstructure:"biz_log"`
	TraceLoggerConf ConfigDef        `json:"trace_log" yaml:"trace_log" mapstructure:"trace_log"`
}

type serverConfigDef struct {
	Name string `json:"name" yaml:"name" mapstructure:"name"`
	Port int    `json:"port" yaml:"port" mapstructure:"port"`
	Env  string `json:"env" yaml:"env" mapstructure:"env"`
}

// LoadConfig 读取配置文件
func LoadConfig(env string) (*AppConfig, error) {
	configDir := getConfDir()
	configName := fmt.Sprintf("%s.conf.yaml", env)
	configPath := path.Join(configDir, configName)
	appConf := &AppConfig{}
	confBytes, err := os.ReadFile(configPath)
	fmt.Println(string(confBytes))
	if err != nil || confBytes == nil {
		return nil, err
	}
	err = yaml.Unmarshal(confBytes, appConf)
	return appConf, err
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// by default: ./conf/
func getConfDir() string {
	dir := "conf"
	for i := 0; i < 3; i++ {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			break
		}
		dir = filepath.Join("..", dir)
	}

	return dir
}
