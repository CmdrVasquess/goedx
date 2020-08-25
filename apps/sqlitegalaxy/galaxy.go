package sqlitegalaxy

import (
	"database/sql"
	"time"

	"git.fractalqb.de/fractalqb/sqlize"
	"git.fractalqb.de/fractalqb/sqlize/null"
	"github.com/CmdrVasquess/goedx"
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

type Galaxy struct {
	db sqlize.DB
}

func Open(dbfile string) (*Galaxy, error) {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return nil, err
	}
	return &Galaxy{db: db}, nil
}

func (g *Galaxy) Close() error {
	return g.db.Close()
}

func (gxy *Galaxy) EdgxSystem(
	addr uint64,
	name string,
	coos []float32,
	touch time.Time,
) (*goedx.System, interface{}) {
	var (
		id      int
		x, y, z sql.NullFloat64
	)
	updateCoos := func() error {
		if len(coos) < 3 {
			return nil
		}
		if x.Valid && y.Valid && z.Valid {
			return nil
		}
		x.Float64, y.Float64, z.Float64 = float64(coos[0]), float64(coos[1]), float64(coos[3])
		return sqlize.Transact(gxy.db, func(tx sqlize.Tx) error {
			_, err := tx.Exec(SqlSetSysCoos.MustSQL(),
				coos[0], coos[1], coos[2],
				id)
			return err
		})
	}
	if addr != sqlize.NoID {
		err := gxy.db.QueryRow(SqlSysByAddr.MustSQL(), addr).
			Scan(&id, &name, &x, &y, &z)
		if err == nil {
			updateCoos()
			return &goedx.System{
				Addr: addr,
				Name: name,
				Coos: goedx.SysCoos{
					goedx.ChgF32(x.Float64),
					goedx.ChgF32(y.Float64),
					goedx.ChgF32(z.Float64),
				},
			}, nil
		} // TODO log error?
	}
	err := gxy.db.QueryRow(SqlSysByName.MustSQL(), name).
		Scan(&id, &name, &addr, &x, &y, &z)
	if err == nil {
		updateCoos()
		return &goedx.System{
			Addr: addr,
			Name: name,
			Coos: goedx.SysCoos{
				goedx.ChgF32(x.Float64),
				goedx.ChgF32(y.Float64),
				goedx.ChgF32(z.Float64),
			},
		}, nil
	}
	if err != sql.ErrNoRows {
		return nil, nil
	}
	err = sqlize.Transact(gxy.db, func(tx sqlize.Tx) error {
		if len(coos) >= 3 {
			x.Valid, x.Float64 = true, float64(coos[0])
			y.Valid, y.Float64 = true, float64(coos[1])
			z.Valid, z.Float64 = true, float64(coos[2])
		}
		_, err := tx.Exec(SqlNewSystem.MustSQL(),
			name,
			null.UInt64{P: &addr},
			&x, &y, &z)
		return err
	})
	if err != nil {
		return nil, nil
	}
	return &goedx.System{
		Addr: addr,
		Name: name,
		Coos: goedx.SysCoos{
			goedx.ChgF32(x.Float64),
			goedx.ChgF32(y.Float64),
			goedx.ChgF32(z.Float64),
		},
	}, nil
}

var gxyEvents = map[events.Type]bool{
	journal.ScanEvent:      true,
	journal.FSDJumpEvent:   true,
	journal.DockedEvent:    true,
	journal.CommanderEvent: true,
	journal.LoadGameEvent:  true,
}

func (gxy *Galaxy) PrepareEDEvent(e events.Event) (token interface{}) {
	if _, ok := gxyEvents[e.EventType()]; ok {
		return e
	}
	return nil
}

func (gxy *Galaxy) FinishEDEvent(token interface{}, e events.Event, chg goedx.Change) {

}
