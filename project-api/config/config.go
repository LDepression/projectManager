/**
 * @Author: lenovo
 * @Description:
 * @File:  config
 * @Version: 1.0.0
 * @Date: 2023/07/15 20:03
 */

package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"projectManager/project-common/logs"
)

var C = InitConfig()

type Config struct {
	viper      *viper.Viper
	SC         *ServeConfig
	GC         *GrpcConfig
	EtcdConfig *EtcdConfig
}

type ServeConfig struct {
	Name string
	Addr string
}

type GrpcConfig struct {
	Name string
	Addr string
}
type EtcdConfig struct {
	Addrs []string
}

func InitConfig() *Config {
	conf := &Config{viper: viper.New()}
	workDir, _ := os.Getwd()
	conf.viper.SetConfigName("config")
	conf.viper.SetConfigType("yaml")
	conf.viper.AddConfigPath("/etc/project/user")
	conf.viper.AddConfigPath(workDir + "/config")
	if err := conf.viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
	conf.ReadServerConfig()
	conf.InitZapLog()
	conf.ReadEtcdConfig()
	return conf
}

func (c *Config) InitZapLog() {
	//从配置中读取日志配置，初始化日志
	lc := &logs.LogConfig{
		DebugFileName: c.viper.GetString("zap.debugFileName"),
		InfoFileName:  c.viper.GetString("zap.infoFileName"),
		WarnFileName:  c.viper.GetString("zap.warnFileName"),
		MaxSize:       c.viper.GetInt("maxSize"),
		MaxAge:        c.viper.GetInt("maxAge"),
		MaxBackups:    c.viper.GetInt("maxBackups"),
	}
	err := logs.InitLogger(lc)
	if err != nil {
		log.Fatalln(err)
	}
}

func (c *Config) ReadServerConfig() {
	sc := &ServeConfig{}
	sc.Name = c.viper.GetString("server.name")
	sc.Addr = c.viper.GetString("server.addr")
	c.SC = sc
}

func (c *Config) ReadEtcdConfig() {
	ec := &EtcdConfig{}
	var addrs []string
	if err := c.viper.UnmarshalKey("etcd.addrs", &addrs); err != nil {
		log.Println(err)
	}
	ec.Addrs = addrs
	c.EtcdConfig = ec
}
