package common

import "github.com/micro/go-micro/v2/config"

type MysqlConfig struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Pwd      string `json:"pwd"`
	Database string `json:"database"`
	Port     string `json:"port"`
}

//get mysql optiop

func GetMysqlFromConsul(config config.Config, path ...string) *MysqlConfig {
	MysqlConfig := &MysqlConfig{}
	config.Get(path...).Scan(MysqlConfig)
	return MysqlConfig
}
