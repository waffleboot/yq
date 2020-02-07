package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type kubeletConfig map[interface{}]interface{}

func (c appConfig) readKubeletConfig() (kubeletConfig, error) {
	file, errOpen := os.Open(c.kubeletConfig)
	if errOpen != nil {
		return nil, errOpen
	}
	defer file.Close()
	k := make(kubeletConfig)
	return k, yaml.NewDecoder(file).Decode(k)
}

func (c appConfig) writeKubeletConfig(k kubeletConfig) error {
	file, errCreate := os.Create(c.kubeletConfig)
	if errCreate != nil {
		return errCreate
	}
	defer file.Close()
	return yaml.NewEncoder(file).Encode(k)
}

func (config kubeletConfig) update(section, key, value string) {
	raw, ok := config[section]
	if !ok {
		raw = make(kubeletConfig)
		config[section] = raw
	}
	raw.(kubeletConfig)[key] = value
}

func (c appConfig) run() error {
	cfg, errRead := c.readKubeletConfig()
	if errRead != nil {
		return errRead
	}
	cfg.update("kubeReserved", "cpu", c.cpu)
	cfg.update("systemReserved", "cpu", c.cpu)
	cfg.update("kubeReserved", "memory", c.mem)
	cfg.update("systemReserved", "memory", c.mem)
	return c.writeKubeletConfig(cfg)
}

type appConfig struct {
	kubeletConfig string
	cpu           string
	mem           string
}

func main() {
	c := appConfig{
		kubeletConfig: os.Args[1],
		cpu:           os.Args[2],
		mem:           os.Args[3],
	}
	if errRun := c.run(); errRun != nil {
		log.Fatal(errRun)
	}
}
