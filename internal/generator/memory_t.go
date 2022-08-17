package generator

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/vl4deee11/pg-gen/internal/log"
	"github.com/vl4deee11/pg-gen/internal/pg"
	"github.com/vl4deee11/pg-gen/internal/typegen"
	"github.com/vl4deee11/pg-gen/internal/yamlcfg"
)

type table struct {
	name        string
	rows        [][]string
	colName2Idx map[string]int
}

type memoryGenerator struct {
	cfg       *yamlcfg.Config
	repo      *pg.Repo
	graph     [][]*empty
	refCnt    []int
	memTables []*table
	name2Idx  map[string]int
}

func NewMemoryGenerator(cfg *yamlcfg.Config, repo *pg.Repo) (Generator, error) {
	name2Idx := buildIndex(cfg)
	gr, refCnt, err := buildGraphOnEmpty(cfg, name2Idx)
	if err != nil {
		return nil, err
	}

	return &memoryGenerator{
		cfg:       cfg,
		repo:      repo,
		graph:     gr,
		refCnt:    refCnt,
		memTables: make([]*table, 0, len(cfg.Tables)),
		name2Idx:  name2Idx,
	}, nil
}

func (g *memoryGenerator) Generate(ctx context.Context) error {
	i := 0
	for i < len(g.refCnt) {
		if g.refCnt[i] == 0 {
			g.refCnt[i] = -1

			log.Logger.Infof("start to generate %d in memory rows data for table = %s", g.cfg.Tables[i].Count, g.cfg.Tables[i].Name)
			t, err := g.createTable(i)
			if err != nil {
				return err
			}
			g.memTables = append(g.memTables, t)

			log.Logger.Infof("start to insert %d rows data for table = %s", g.cfg.Tables[i].Count, g.cfg.Tables[i].Name)

			err = g.batchInsert(ctx, t, i)
			if err != nil {
				return err
			}

			for j := range g.graph {
				if i == j {
					continue
				}
				if g.graph[j][i] != nil {
					g.graph[j][i] = nil
					g.refCnt[j]--
					log.Logger.Debugf("graph: disconnect table = %s and table = %s", g.cfg.Tables[j].Name, g.cfg.Tables[i].Name)
				}
			}
			i = 0
			continue
		}
		i++
	}
	return nil
}

func (g *memoryGenerator) createTable(ti int) (*table, error) {
	var (
		t        = g.cfg.Tables[ti]
		tName    = t.Name
		i        = 0
		memTable = &table{
			name:        tName,
			rows:        make([][]string, 0, t.Count),
			colName2Idx: map[string]int{},
		}
	)

	for j := 0; j < len(t.Fields); j++ {
		memTable.colName2Idx[t.Fields[j].Name] = j
	}

	for ; i < t.Count; i++ {
		row := make([]string, len(t.Fields))
		for j := 0; j < len(t.Fields); j++ {
			f, ok := typegen.GenMapMemoryT[t.Fields[j].Type]
			if !ok {
				return nil, fmt.Errorf("type %s for table %s not recognized", t.Fields[j].Type, tName)
			}

			switch t.Fields[j].Type {
			case typegen.PK:
				val, err := g.handlePK(t, j, f)
				if err != nil {
					return nil, err
				}
				row[j] = val
			case typegen.UniqPK:
				val, err := g.handleUniqPK(t, j, f)
				if err != nil {
					return nil, err
				}
				row[j] = val
			default:
				row[j] = f(t.Fields[j].Nullable, t.Fields[j].GParams...)
			}
		}
		memTable.rows = append(memTable.rows, row)
	}

	return memTable, nil
}

func (g *memoryGenerator) handlePK(t *yamlcfg.Table, j int, f func(nullable bool, ps ...interface{}) string) (string, error) {
	if len(t.Fields[j].GParams) != 2 {
		return "", fmt.Errorf("type %s for table %s should have genParams : [table, pk_fields_name]", t.Fields[j].Type, t.Name)
	}
	pkIdx := g.name2Idx[t.Fields[j].GParams[0].(string)]
	if pkIdx >= len(g.memTables) {
		pkIdx = findIdxByTableName(g.memTables, t.Fields[j].GParams[0].(string))
		if pkIdx == -1 {
			return "", fmt.Errorf("type %s for table %s pk table not generated yet", t.Fields[j].Type, t.Name)
		}
	}
	pkFieldIdx := g.memTables[pkIdx].colName2Idx[t.Fields[j].GParams[1].(string)]

	return f(t.Fields[j].Nullable, g.memTables[pkIdx].rows, pkFieldIdx), nil
}

func (g *memoryGenerator) handleUniqPK(t *yamlcfg.Table, j int, f func(nullable bool, ps ...interface{}) string) (string, error) {
	if len(t.Fields[j].GParams) != 2 {
		return "", fmt.Errorf("type %s for table %s should have genParams : [table, pk_fields_name]", t.Fields[j].Type, t.Name)
	}
	pkIdx := g.name2Idx[t.Fields[j].GParams[0].(string)]
	if pkIdx >= len(g.memTables) {
		pkIdx = findIdxByTableName(g.memTables, t.Fields[j].GParams[0].(string))
		if pkIdx == -1 {
			return "", fmt.Errorf("type %s for table %s pk table not generated yet", t.Fields[j].Type, t.Name)
		}
	}
	pkFieldIdx := g.memTables[pkIdx].colName2Idx[t.Fields[j].GParams[1].(string)]
	val := f(t.Fields[j].Nullable, g.memTables[pkIdx].rows, pkFieldIdx, pkIdx, j)
	if val == "-1" {
		return "", fmt.Errorf("too small table %s for generate unique pkey for table %s", t.Fields[j].GParams[0].(string), t.Name)
	}

	return val, nil
}

func (g *memoryGenerator) batchInsert(ctx context.Context, t *table, ti int) error {
	var (
		batchSize = g.cfg.InsertBatchSize
		completed = 0
		count     = len(t.rows)
		queryChan = make(chan string, g.cfg.MaxConcurrency)
		errChan   = make(chan error, g.cfg.MaxConcurrency)
		wg        sync.WaitGroup
		i         int64 = 0
	)

	for ; i < g.cfg.MaxConcurrency; i++ {
		wg.Add(1)
		go g.exec(&wg, ctx, queryChan, errChan)
	}

	for completed < count {
		if completed+batchSize > count {
			batchSize = count - completed
		}

		q, err := g.buildQuery(completed, t, ti)
		if err != nil {
			return err
		}

		queryChan <- q
		completed += batchSize
		log.Logger.Debugf("generated %d rows data for table = %s", completed, t.name)
	}
	close(queryChan)
	wg.Wait()
	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

func (g *memoryGenerator) buildQuery(from int, t *table, ti int) (string, error) {
	to := from + g.cfg.InsertBatchSize
	query := strings.Builder{}
	if to > len(t.rows) {
		to = len(t.rows)
	}
	li := to - 1
	query.WriteString(fmt.Sprintf(`INSERT INTO %s(%s) VALUES `, t.name, fieldsForInsert(g.cfg.Tables[ti].Fields)))
	for i := from; i < to; i++ {
		query.WriteRune('(')
		row := t.rows[i]
		fli := len(row) - 1
		for j := 0; j < len(row); j++ {
			query.WriteString(row[j])
			if j != fli {
				query.WriteString(", ")
			}
		}
		query.WriteRune(')')
		if i != li {
			query.WriteString(", ")
		}
	}
	query.WriteString(" ;")
	return query.String(), nil
}

func (g *memoryGenerator) exec(wg *sync.WaitGroup, ctx context.Context, qChan chan string, errChan chan error) {
	for q := range qChan {
		err := g.repo.Exec(ctx, q)
		if err != nil {
			errChan <- err
		}
	}
	wg.Done()
}
