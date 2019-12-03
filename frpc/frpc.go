package frpc

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log"
	"path"
)

type Proxy struct {
	IP string
	Port string
	Name string
}

func (p *Proxy) String() string {
	return fmt.Sprintf("%+v\n", *p)
}

type Frpc struct {
	Config *ini.File
	FrpWorkDirectory string
}

func (f *Frpc) AddProxy() error {

}

func NewFrpc(frpWorkDirectory string) *Frpc {
	var frpc = new(Frpc)
	var configPath = path.Join(frpWorkDirectory, "frpc.ini")
	config, err := ini.Load(configPath)
	if err != nil {
		log.Fatal("load config " + configPath + " failed:", err.Error())
	}
	frpc.FrpWorkDirectory = frpWorkDirectory
	frpc.Config = config
	return frpc
}