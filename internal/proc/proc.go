package proc

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/lib/pq" //nolint:golint

	"github.com/vl4deee11/pg-gen/internal"
	"github.com/vl4deee11/pg-gen/internal/generator"
	"github.com/vl4deee11/pg-gen/internal/pg"
	"github.com/vl4deee11/pg-gen/internal/yamlcfg"
)

type Proc struct {
	cfg  *yamlcfg.Config
	repo *pg.Repo
}

func New(cfg *internal.Conf) (*Proc, error) {
	proc := new(Proc)
	yamlCfg, err := yamlcfg.ParseConfig(cfg.ConfPath)
	if err != nil {
		return nil, err
	}

	proc.cfg = yamlCfg
	db, err := sql.Open("postgres", cfg.PGDSN)
	if err != nil {
		return nil, err
	}

	proc.repo = pg.New(db)
	return proc, nil
}

func (p *Proc) Run() error {
	ctx := context.Background()
	// Check that all tables exists
	// TODO: create table
	for i := range p.cfg.Tables {
		err := p.repo.FindTableByName(ctx, p.cfg.Tables[i].Name)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errors.New("not found table = " + p.cfg.Tables[i].Name)
			}
			return err
		}
	}

	switch p.cfg.Type {
	case yamlcfg.Inserts:
		gen, err := generator.NewInsertGenerator(p.cfg, p.repo)
		if err != nil {
			return err
		}
		return gen.Generate(ctx)
	case yamlcfg.Memory:
		gen, err := generator.NewMemoryGenerator(p.cfg, p.repo)
		if err != nil {
			return err
		}
		return gen.Generate(ctx)
	default:
		return errors.New("not found type of generator = " + string(p.cfg.Type))
	}
}
