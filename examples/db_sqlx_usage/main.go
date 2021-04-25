package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yscsky/yu"
)

var sdb *yu.StmtDB

var pgSQLs = []string{
	`create table if not exists demo1(
		demo_key text primary key not null,
		create_time bigint not null
	)`,
}

func init() {
	sdb = yu.NewStmtDB(sqlx.MustOpen("postgres", "postgres://admin:admin@localhost:5432/demodb?sslmode=disable"))
	for _, stmt := range pgSQLs {
		if _, err := sdb.Exec(stmt); err != nil {
			yu.LogErr(err, "init sdb exec")
		}
	}
}

func main() {
	defer func() {
		if sdb != nil {
			sdb.CloseDB()
		}
	}()
	for i := 0; i < 100; i++ {
		yu.LogErr(insertOne(fmt.Sprintf("demo%d", i), yu.NewStampTime(time.Now())), "insertOne")
	}
	createTime, err := queryOne("demo10")
	if err != nil {
		yu.LogErr(err, "queryOne")
		return
	}
	log.Println("[info] - demo10 create at:", createTime)
	yu.LogErr(deleteOne("demo10"), "deleteOne")
}

func insertOne(demoKey string, createTime yu.TStamp) (err error) {
	stmt, err := sdb.Stmt(`insert into demo1 values($1,$2) on conflict(demo_key) do update set create_time=$2`)
	if err != nil {
		return
	}
	_, err = stmt.Exec(demoKey, createTime)
	return
}

func queryOne(demoKey string) (createTime yu.TStamp, err error) {
	stmt, err := sdb.Stmt(`select block_time from demo1 where demo_key=$1`)
	if err != nil {
		return
	}
	err = stmt.QueryRow(demoKey).Scan(&createTime)
	return
}

func deleteOne(demoKey string) (err error) {
	stmt, err := sdb.Stmt(`delete from demo1 where demo_key=$1`)
	if err != nil {
		return
	}
	_, err = stmt.Exec(demoKey)
	return
}

func queryByTime(start, end yu.TStamp) (list []string, err error) {
	stmt, err := sdb.Stmt(`select demo_key from demo1 where create_time>$1 and create_time<$2`)
	if err != nil {
		return
	}
	err = stmt.Select(&list, start, end)
	return
}
