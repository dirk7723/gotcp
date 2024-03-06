package untils

import (
	"encoding/json"
	"io/ioutil"

	"hbq.com/ggame/hinx/hiface"
)

//存储hinx框架的全局参数，供其他模块使用

type Globalobj struct {
	TcpServer hiface.Iserver

	Host string

	TcpPort int

	Name string

	Version string

	MaxConn int

	MaxPackageSize uint32

	WorkerPoolSize uint32

	MaxWorkerTaskLen uint32
}

var GlobalObject *Globalobj

func (g *Globalobj) Reload() {
	data, err := ioutil.ReadFile("config/hinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	GlobalObject = &Globalobj{
		Name:             "HINX_SERVER",
		Version:          "V0.4",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}
	GlobalObject.Reload()
}
