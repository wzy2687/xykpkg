package xutils

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log/slog"
)

func MustNewLeveldb(dbPath string) *leveldb.DB {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		slog.Error("init leveldb err", "path", dbPath, "err", err.Error())
		panic(err)
	}
	return db
}
