package server

type Server interface {
	SetConf(map[string]interface{})
	Init()
	Run() error
}