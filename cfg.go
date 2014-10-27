package goutils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Cfg struct {
	programName        string
	fileName           string
	uncommittedChanges map[string]string
	store              map[string]string
}

func NewCfg(programName, envVar string) (*Cfg, error) {
	envvar := os.Getenv(envVar)
	if envvar == "" {
		return nil, errors.New("The environment variable " + envVar + " is undefined")
	}

	cfg := &Cfg{
		programName:        programName,
		fileName:           envvar,
		uncommittedChanges: make(map[string]string),
		store:              make(map[string]string),
	}

	bs, err := ioutil.ReadFile(cfg.fileName)
	if err != nil {
		return nil, err
	}
	str := string(bs)
	for i, line := range strings.Split(str, "\n") {
		line := strings.TrimSpace(line)
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		v := strings.SplitN(line, "=", 2)
		if len(v) != 2 {
			return nil, errors.New("Format error at line " + strconv.Itoa(i+1))
		}

		key, value := strings.TrimSpace(v[0]), strings.TrimSpace(v[1])
		cfg.store[key] = value
	}

	return cfg, nil
}

func (cfg *Cfg) Set(key, value string) {
	cfg.uncommittedChanges[key] = value
}

func (cfg *Cfg) Get(key string) (string, bool) {
	value, ok := cfg.store[key]
	return value, ok
}

func (cfg *Cfg) Commit() error {
	// TODO: isn't there a function to do this?
	for k, v := range cfg.uncommittedChanges {
		cfg.store[k] = v
	}

	cfg.uncommittedChanges = make(map[string]string)

	file, err := os.Create(cfg.fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString("# " + cfg.programName + " configuration file\n\n")
	for k, v := range cfg.store {
		file.WriteString(fmt.Sprintf("%s = %s\n", k, v))
	}

	return nil
}

func (cfg *Cfg) Rollback() error {
	cfg.uncommittedChanges = make(map[string]string)
	return nil
}
