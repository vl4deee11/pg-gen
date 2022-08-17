package generator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/vl4deee11/pg-gen/internal/yamlcfg"
)

func fieldsForInsert(fs []*yamlcfg.Field) string {
	r := make([]string, len(fs))
	for i := range fs {
		r[i] = fs[i].Name
	}
	return strings.Join(r, ", ")
}

func buildIndex(cfg *yamlcfg.Config) map[string]int {
	name2Idx := make(map[string]int)
	for i := range cfg.Tables {
		name2Idx[cfg.Tables[i].Name] = i
	}
	return name2Idx
}

func buildGraphOnEmpty(cfg *yamlcfg.Config, name2Idx map[string]int) ([][]*empty, []int, error) {
	graph := make([][]*empty, len(cfg.Tables))
	for i := range graph {
		graph[i] = make([]*empty, len(cfg.Tables))
	}

	refCnt := make([]int, len(cfg.Tables))
	for i := range cfg.Tables {
		for j := range cfg.Tables[i].Fields {
			if cfg.Tables[i].Fields[j].Type == "pk" {
				sl := cfg.Tables[i].Fields[j].GParams
				if len(cfg.Tables[i].Fields[j].GParams) != 2 {
					return nil, nil, errors.New("format pk should be '['table','field']")
				}
				ti, ok := name2Idx[sl[0].(string)]
				if !ok {
					return nil, nil, fmt.Errorf("not found table for pk = %s", cfg.Tables[i].Fields[j].Pk)
				}

				graph[i][ti] = &empty{}
				refCnt[i]++
			}
		}
	}
	return graph, refCnt, nil
}

func findIdxByTableName(tables []*table, name string) int {
	for i := range tables {
		if tables[i].name == name {
			return i
		}
	}
	return -1
}
