CREATE TABLE systems (
  id INTEGER PRIMARY KEY
, name TEXT NOT NULL
, addr INTEGER UNIQUE
, x    REAL
, y    REAL
, z    REAL
);

CREATE INDEX systems_x ON systems (x);

CREATE INDEX systems_y ON systems (y);

CREATE INDEX systems_z ON systems (z);

CREATE TABLE cmdrs (
  id INTEGER PRIMARY KEY
, fid TEXT UNIQUE
, name TEXT NOT NULL UNIQUE
);

CREATE TABLE visits (
  cmdr INTEGER NOT NULL REFERENCES cmdrs(id)
, sys INTEGER NOT NULL REFERENCES systems(id)
, arrive TEXT NOT NULL
, UNIQUE (cmdr, arrive)
);

CREATE INDEX visits_arrive ON visits(arrive);

CREATE TABLE ports (
  id INTEGER PRIMARY KEY
, sys INTEGER NOT NULL REFERENCES systems(id)
, name TEXT NOT NULL
, type TEXT NOT NULL
, UNIQUE (sys, name)
);

CREATE TABLE docked (
  cmdr INTEGER NOT NULL REFERENCES cmdrs(id)
, port INTEGER NOT NULL REFERENCES ports(id)
, arrive TEXT NOT NULL
, UNIQUE (cmdr, arrive)
);

CREATE INDEX docked_arrive ON docked(arrive);
