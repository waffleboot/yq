package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func read(filename string) (map[interface{}]interface{}, error) {
	file, errOpen := os.Open(filename)
	if errOpen != nil {
		return nil, errOpen
	}
	defer file.Close()
	ans := make(map[interface{}]interface{})
	return ans, yaml.NewDecoder(file).Decode(ans)
}

func write(filename string, m interface{}) error {
	file, errCreate := os.Create(filename)
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

func update(obj map[interface{}]interface{}) {
	obj["cpu"] = "250m"
	obj["memory"] = "250M"
}

func run(filename string) error {
	cfg, errRead := read(filename)
	if errRead != nil {
		return errRead
	}
	updater := func(name string) {
		obj := get(cfg, name)
		update(obj)
		cfg[name] = obj
	}
	updater("systemReserved")
	updater("kubeReserved")
	return write(filename, cfg)
}

func main() {
	if errRun := run(os.Args[1]); errRun != nil {
		log.Fatal(errRun)
	}
}
