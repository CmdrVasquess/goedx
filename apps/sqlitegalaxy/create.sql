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
  id INTEGER PRIMARY KEY
, cmdr INTEGER NOT NULL REFERENCES cmdrs(id)
, sys INTEGER NOT NULL REFERENCES systems(id)
, arrive TEXT NOT NULL
);

CREATE INDEX visits_arrive ON visits (arrive);

