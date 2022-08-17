package internal

import (
	"flag"

	"github.com/vl4deee11/pg-gen/internal/log"
	"github.com/vl4deee11/pg-gen/internal/util"
)

type Conf struct {
	ConfPath string `desc:"path to config"`
	LogLvl   string `desc:"log level"`
	PGDSN    string `desc:"postgresql config"`
}

func (c *Conf) Parse() {
	// COMMON
	flag.StringVar(&c.ConfPath, "c", "./data/conf.yaml", "path to config file")
	flag.StringVar(&c.PGDSN, "dsn", "host=localhost  port=16406 user=picker_backend password=picker_backend dbname=picker_backend sslmode=disable binary_parameters=yes", "postgresql dsn config")
	flag.StringVar(&c.LogLvl, "v", "debug", "log level")
	flag.Parse()

	log.MakeLogger(c.LogLvl)
	c.print()
}

func (c *Conf) print() {
	util.PrintFromDesc("[COMMON CONFIG]:", *c)
}
