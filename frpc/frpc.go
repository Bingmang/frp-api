package frpc

import (
	"fmt"
	"gopkg.in/ini.v1"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path"
)

var (
	MAX_RANDOM_PORT = 19999
	MIN_RANDOM_PORT = 10001
	RANDOM_PORT_RETRY = 3
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
	config           *ini.File
	frpWorkDirectory string
	configPath       string
	frpcExecPath     string
}

func (f *Frpc) Backup() error {
	var (
		source *os.File
		backup *os.File
		err error
	)
	source, err = os.Open(f.configPath)
	if err != nil {
		return err
	}
	defer source.Close()
	backup, err = os.Create(path.Join(f.frpWorkDirectory, "frpc.ini.backup"))
	if err != nil {
		return err
	}
	defer backup.Close()
	_, err = io.Copy(backup, source)
	return err
}

func (f *Frpc) AddProxy(p *Proxy) error {
	var (
		randomPort = 0
		randomPortExistsCount = 0
	)
	for randomPortExistsCount < RANDOM_PORT_RETRY {
		var tempRandomPort = rand.Intn(MAX_RANDOM_PORT - MIN_RANDOM_PORT) + MAX_RANDOM_PORT
		for _, section := range f.config.Sections() {
			if section.HasValue(string(tempRandomPort)) {
				randomPortExistsCount++
				break
			} else {
				randomPort = tempRandomPort
			}
		}
		if randomPort > 0 {
			break
		}
	}
	if randomPort == 0 {
		return fmt.Errorf("random port generate failed")
	}

	sec, err := f.config.NewSection(p.Name)
	if err != nil {
		return err
	}
	if _, err = sec.NewKey("local_ip", p.IP); err != nil {
		return err
	}
	if _, err = sec.NewKey("local_port", p.Port); err != nil {
		return err
	}
	if _, err = sec.NewKey("remote_port", string(randomPort)); err != nil {
		return err
	}

	if err = f.Backup(); err != nil {
		return err
	}
	return f.config.SaveTo(f.configPath)
}

func (f *Frpc) Reload() {
	cmd := exec.Command(f.frpcExecPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("frpc reload failed:", err.Error())
	}
	log.Println("frpc reload success:", out)
}

func NewFrpc(frpWorkDirectory string) *Frpc {
	var (
		configPath = path.Join(frpWorkDirectory, "frpc.ini")
		frpcExecPath = path.Join(frpWorkDirectory, "frpc")
	)
	config, err := ini.Load(configPath)
	if err != nil {
		log.Fatal("load config " + configPath + " failed:", err.Error())
	}
	if _, err := exec.LookPath(frpcExecPath); err != nil {
		log.Fatal("frpc command not found:", err.Error())
	}
	return &Frpc{
		config:           config,
		frpWorkDirectory: frpWorkDirectory,
		configPath:       configPath,
		frpcExecPath:     frpcExecPath,
	}
}
