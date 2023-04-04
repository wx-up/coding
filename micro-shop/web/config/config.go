package config

import "github.com/spf13/viper"

// 初始化一下
var c = &SystemConfig{}

func Config() *SystemConfig {
	return c
}

type SystemConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`

	Mysql      MysqlConfig      `mapstructure:"mysql"`
	UserServer UserServerConfig `mapstructure:"user_server"`

	Jwt JwtConfig `mapstructure:"jwt"`

	Consul ConsulConfig `mapstructure:"consul"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type JwtConfig struct {
	Secret     string `mapstructure:"secret"`
	ExpireTime int64  `mapstructure:"expire_time"`
}

type UserServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Name string `mapstructure:"name"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

func InitConfig(filepath string) {
	conf := viper.New()
	conf.SetConfigFile(filepath)
	if err := conf.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := conf.Unmarshal(c); err != nil {
		panic(err)
	}
}
