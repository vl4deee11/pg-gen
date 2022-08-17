package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/vl4deee11/pg-gen/internal"
	"github.com/vl4deee11/pg-gen/internal/log"
	"github.com/vl4deee11/pg-gen/internal/proc"
)

func main() {
	cfg := new(internal.Conf)
	cfg.Parse()
	log.Logger.Info("=========== START ===========")

	defer log.Logger.Info("=========== FINISH ===========")

	g, err := proc.New(cfg)
	if err != nil {
		log.Logger.Errorf("initialize generator. ERR => %s\n", err.Error())
		return
	}
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	err = g.Run()
	if err != nil {
		log.Logger.Errorf("generator run. ERR => %s\n", err.Error())
		return
	}
}
