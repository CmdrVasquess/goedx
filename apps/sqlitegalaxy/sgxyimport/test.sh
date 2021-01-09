#!/bin/sh
DB=test.db
rm -f $DB
sqlite3 $DB < ../create.sql
go build
./sgxyimport -db $DB \
	     $HOME/doc/share/EliteDangerous/journal.bak/Journal.*.log.gz \
	     $HOME/doc/share/EliteDangerous/journal/Journal.*.log
