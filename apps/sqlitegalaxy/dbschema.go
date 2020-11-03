package sqlitegalaxy

import bsq "git.fractalqb.de/fractalqb/sqlize/buildsq"

func init() {
	bsq.InitDecls("", bsq.CamelToSnake(),
		&TSystems,
		&TBodies,
		&TPrntBdys,
		&TCmdrs,
		&TVisits,
		&TPorts,
		&TDocked,
	)
}

var TSystems = struct {
	bsq.Table
	ID, Name, Addr, X, Y, Z bsq.Column
}{
	Table: bsq.Table{TableName: "systems", TableAlias: "s"},
	ID:    bsq.Column{Name: "id"},
}

var TBodyType = struct {
	bsq.Table
	ID, Name bsq.Column
}{
	Table: bsq.Table{TableName: "bodytype", TableAlias: "bt"},
	ID:    bsq.Column{Name: "id"},
}

var TBodyClass = struct {
	bsq.Table
	ID, Typ, Name bsq.Column
}{
	Table: bsq.Table{TableName: "bodyclass", TableAlias: "bc"},
	ID:    bsq.Column{Name: "id"},
}

var TBodies = struct {
	bsq.Table
	Sys, BodyID, Name, Type, Class, DistFromArvl bsq.Column
}{
	Table:        bsq.Table{TableName: "bodies", TableAlias: "b"},
	BodyID:       bsq.Column{Name: "bid"},
	DistFromArvl: bsq.Column{Name: "distfa"},
}

var TPrntBdys = struct {
	bsq.Table
	Sys, PB, CB bsq.Column
}{
	Table: bsq.Table{TableName: "prntbdys", TableAlias: "pb"},
	PB:    bsq.Column{Name: "pb"},
	CB:    bsq.Column{Name: "cb"},
}

var TCmdrs = struct {
	bsq.Table
	ID, FID, Name bsq.Column
}{
	Table: bsq.Table{TableName: "cmdrs", TableAlias: "c"},
	ID:    bsq.Column{Name: "id"},
	FID:   bsq.Column{Name: "fid"},
}

var TVisits = struct {
	bsq.Table
	Cmdr, Sys, Arrive bsq.Column
}{
	Table: bsq.Table{TableName: "visits", TableAlias: "v"},
}

var TPorts = struct {
	bsq.Table
	ID, Sys, Name, Type bsq.Column
}{
	Table: bsq.Table{TableName: "ports", TableAlias: "p"},
	ID:    bsq.Column{Name: "id"},
}

var TDocked = struct {
	bsq.Table
	Cmdr, Port, Arrive bsq.Column
}{
	Table: bsq.Table{TableName: "docked", TableAlias: "d"},
}

var (
	SqlCmdrByName = bsq.Lazy{
		Query: bsq.Select{From: &TCmdrs, SelectBy: bsq.Cols(&TCmdrs.Name)}}

	SqlNewCmdr = bsq.Lazy{
		Query: bsq.CreateStatement{Table: &TCmdrs, ID: &TCmdrs.ID}}

	SqlCmdrSetFID = bsq.Lazy{
		Query: bsq.Update{
			Set:      bsq.Cols(&TCmdrs.FID),
			SelectBy: bsq.Cols(&TCmdrs.ID),
		}}

	SqlAddVisit = bsq.Lazy{Query: bsq.Insert{bsq.ColsOf(&TVisits)}}

	SqlAddDocked = bsq.Lazy{Query: bsq.Insert{bsq.ColsOf(&TDocked)}}

	SqlPortBySysNm = bsq.Lazy{
		Query: bsq.Select{
			Columns:  bsq.Cols(&TPorts.ID),
			SelectBy: bsq.Cols(&TPorts.Sys, &TPorts.Name),
		}}

	SqlNewPort = bsq.Lazy{
		Query: bsq.CreateStatement{Table: &TPorts, ID: &TPorts.ID}}

	SqlSetSysCoos = bsq.Lazy{
		Query: bsq.Update{
			Set:      bsq.Cols(&TSystems.X, &TSystems.Y, &TSystems.Z),
			SelectBy: bsq.Cols(&TSystems.ID),
		}}

	SqlSysByAddr = bsq.Lazy{
		Query: bsq.Select{
			From:     &TSystems,
			SelectBy: bsq.Cols(&TSystems.Addr),
		}}

	SqlSysByName = bsq.Lazy{
		Query: bsq.Concat{
			"SELECT ", bsq.ColsOf(&TSystems),
			" FROM ", &TSystems,
			" WHERE lower(", &TSystems.Name, ") = lower(", bsq.Arg{"name"}, ")",
		}}

	SqlNewSystem = bsq.Lazy{
		Query: bsq.CreateStatement{Table: &TSystems, ID: &TSystems.ID}}

	SqlBodyByName = bsq.Lazy{
		Query: bsq.Select{
			From:     &TBodies,
			SelectBy: bsq.Cols(&TBodies.Sys, &TBodies.Name),
		}}

	SqlInsertBody = bsq.Lazy{
		Query: bsq.Insert{Columns: bsq.ColsOf(&TBodies)}}
)
