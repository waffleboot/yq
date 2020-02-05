package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func (c config) read() (map[interface{}]interface{}, error) {
	file, errOpen := os.Open(c.filename)
	if errOpen != nil {
		return nil, errOpen
	}
	defer file.Close()
	ans := make(map[interface{}]interface{})
	return ans, yaml.NewDecoder(file).Decode(ans)
}

func (c config) write(m interface{}) error {
	file, errCreate := os.Create(c.filename)
	if errCreate != nil {
		return errCreate
	}
	defer file.Close()
	return yaml.NewEncoder(file).Encode(m)
}

func get(m map[interface{}]interface{}, name string) map[interface{}]interface{} {
	raw, ok := m[name]
	if ok {
		return raw.(map[interface{}]interface{})
	}
	return make(map[interface{}]interface{})
}

func (c config) update(obj map[interface{}]interface{}) {
	obj["cpu"] = c.cpu
	obj["memory"] = c.mem
}

func (c config) run() error {
	cfg, errRead := c.read()
	if errRead != nil {
		return errRead
	}
	updater := func(name string) {
		obj := get(cfg, name)
		c.update(obj)
		cfg[name] = obj
	}
	updater("systemReserved")
	updater("kubeReserved")
	return c.write(cfg)
}

type config struct {
	filename string
	cpu      string
	mem      string
}

func main() {
	c := config{
		filename: os.Args[1],
		cpu:      os.Args[2],
		mem:      os.Args[3],
	}
	if errRun := c.run(); errRun != nil {
		log.Fatal(errRun)
	}
}
