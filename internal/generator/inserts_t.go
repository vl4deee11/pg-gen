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

type insertGenerator struct {
	cfg    *yamlcfg.Config
	repo   *pg.Repo
	graph  [][]*empty
	refCnt []int
}

func NewInsertGenerator(cfg *yamlcfg.Config, repo *pg.Repo) (Generator, error) {
	name2Idx := buildIndex(cfg)
	gr, refCnt, err := buildGraphOnEmpty(cfg, name2Idx)
	if err != nil {
		return nil, err
	}
	return &insertGenerator{
		cfg:    cfg,
		repo:   repo,
		graph:  gr,
		refCnt: refCnt,
	}, nil
}

func (g *insertGenerator) Generate(ctx context.Context) error {
	i := 0
	for i < len(g.refCnt) {
		if g.refCnt[i] == 0 {
			g.refCnt[i] = -1

			log.Logger.Infof("start to generate %d rows data for table = %s", g.cfg.Tables[i].Count, g.cfg.Tables[i].Name)
			err := g.batchInsert(ctx, i)
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
					log.Logger.Infof("graph: disconnect table = %s and table = %s", g.cfg.Tables[j].Name, g.cfg.Tables[i].Name)
				}
			}
			i = 0
			continue
		}
		i++
	}
	return nil
}

func (g *insertGenerator) batchInsert(ctx context.Context, ti int) error {
	var (
		batchSize = g.cfg.InsertBatchSize
		completed = 0
		t         = g.cfg.Tables[ti]
		queryChan = make(chan string, g.cfg.MaxConcurrency)
		errChan   = make(chan error, g.cfg.MaxConcurrency)
		wg        sync.WaitGroup
		i         int64 = 0
	)

	for ; i < g.cfg.MaxConcurrency; i++ {
		wg.Add(1)
		go g.exec(&wg, ctx, queryChan, errChan)
	}

	for completed < t.Count {
		if completed+batchSize > t.Count {
			batchSize = t.Count - completed
		}
		q, err := g.buildQuery(ti, batchSize)
		if err != nil {
			return err
		}

		queryChan <- q
		completed += batchSize
		log.Logger.Debugf("generated %d rows data for table = %s", completed, t.Name)
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

func (g *insertGenerator) buildQuery(ti int, cnt int) (string, error) {
	var (
		t     = g.cfg.Tables[ti]
		li    = cnt - 1
		tName = t.Name
		query strings.Builder
	)

	query.WriteString(fmt.Sprintf(`INSERT INTO %s(%s) VALUES `, tName, fieldsForInsert(g.cfg.Tables[ti].Fields)))

	for i := 0; i < cnt; i++ {
		query.WriteRune('(')
		fli := len(t.Fields) - 1
		for j := 0; j < len(t.Fields); j++ {
			f, ok := typegen.GenMapInsertT[t.Fields[j].Type]
			if !ok {
				return "", fmt.Errorf("type %s for table %s not recognized", t.Fields[j].Type, tName)
			}
			v := f(t.Fields[j].Nullable, t.Fields[j].GParams...)
			query.WriteString(v)
			if j != fli {
				query.WriteString(", ")
			}
		}
		query.WriteRune(')')
		if i != li {
			query.WriteString(", ")
		}
	}
	// TODO: support for unique pkeys
	query.WriteString(" ;")
	return query.String(), nil
}

func (g *insertGenerator) exec(wg *sync.WaitGroup, ctx context.Context, qChan chan string, errChan chan error) {
	for q := range qChan {
		err := g.repo.Exec(ctx, q)
		if err != nil {
			errChan <- err
		}
	}
	wg.Done()
}
