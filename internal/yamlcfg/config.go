package yamlcfg

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/vl4deee11/pg-gen/internal/typegen"
)

type Type string

const (
	Inserts Type = "inserts"
	Memory  Type = "memory"
)

type Config struct {
	InsertBatchSize int      `yaml:"insertBatchSize"`
	MaxConcurrency  int64    `yaml:"maxConcurrency"`
	Type            Type     `yaml:"type"`
	Tables          []*Table `yaml:"tables"`
}

type Table struct {
	Name   string   `yaml:"name"`
	Count  int      `yaml:"rowCount"`
	Fields []*Field `yaml:"fields"`
}

type Field struct {
	Name     string         `yaml:"name"`
	Type     typegen.PGType `yaml:"type"`
	Pk       string         `yaml:"pk"`
	GParams  []interface{}  `yaml:"genParams"`
	Nullable bool           `yaml:"nullable"`
}

func ParseConfig(path string) (*Config, error) {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	err = yaml.Unmarshal(fileBytes, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
