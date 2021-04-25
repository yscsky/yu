package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/yscsky/yu"
	"gorm.io/gorm/logger"
)

type Demo struct {
	UID        string    `gorm:"uid"`
	Title      string    `gorm:"title"`
	Content    string    `gorm:"content"`
	ReadCount  int       `gorm:"read_count"`
	CreateTime yu.TStamp `gorm:"create_time"`
	UpdateTime yu.TStamp `gorm:"update_time"`
}

const uid = "5259dbc3a15b4ef7865539568574cdde"

var gdb *yu.GormDB

func init() {
	gdb = yu.MustOpenMySQL(yu.DSN{
		Username: "root",
		Password: "123456",
		URL:      "127.0.0.1",
		Port:     "3306",
		DBName:   "gormdemo",
		LogLevel: logger.Info,
	})
}

func main() {
	defer gdb.CloseDB()
	if _, err := gdb.ExecSQL(`create table if not exists demo(
		uid varchar(32) not null,
		title varchar(32) not null,
		content text not null,
		read_count int not null,
		create_time bigint not null,
		update_time bigint not null,
		primary key (uid),
		key idx_title (title)
	)`); err != nil {
		yu.LogErr(err, "create table")
		return
	}

	// insertUsage()
	insert()
	selectOne()
	selectList()
	selectOneDemo()
	updateOne()
	count()
	queryAll()
}

func insert() {
	d := &Demo{
		UID:        "0c72bfc0c7c04011988e0e03084e57d7",
		Title:      "title",
		Content:    "content",
		ReadCount:  0,
		CreateTime: yu.NewStampTime(time.Now()),
		UpdateTime: yu.NewStampTime(time.Now()),
	}
	yu.LogErr(gdb.Insert("demo", d), "Insert")
	yu.LogErr(gdb.Insert("demo", d), "Insert")
}

func count() {
	count := 0
	if err := gdb.Query(&count, `select count(uid) from demo`); err != nil {
		yu.LogErr(err, "count")
		return
	}
	yu.Logf("count: %d", count)
}

func updateOne() {
	q := &Demo{UID: uid, ReadCount: 100, UpdateTime: yu.NewStampTime(time.Now())}
	if _, err := gdb.ExecSQL("update demo set read_count=?,update_time=? where uid=?", q.ReadCount, q.UpdateTime, q.UID); err != nil {
		yu.LogErr(err, "update")
		return
	}
	selectOne()
}

func selectOneDemo() {
	q := &Demo{UID: uid}
	d := &Demo{}
	if err := gdb.Query(&d, "select * from demo where uid=@uid", sql.Named("uid", q.UID)); err != nil {
		yu.LogErr(err, "select")
		return
	}
	yu.Logf("select one demo: %v", d)
}

func selectList() {
	list := make([]*Demo, 0)
	if err := gdb.Query(&list, "select * from demo where create_time=?", 1611915765); err != nil {
		yu.LogErr(err, "select")
		return
	}
	for _, d := range list {
		yu.Logf("select list: %v", d)
	}
}

func selectOne() {
	d := &Demo{}
	if err := gdb.Query(&d, "select * from demo where uid=?", uid); err != nil {
		yu.LogErr(err, "select")
		return
	}
	yu.Logf("select one: %v", d)
}

func insertUsage() {
	list := make([]*Demo, 300)
	for i := 0; i < 300; i++ {
		list[i] = &Demo{
			UID:        yu.UUID(),
			Title:      fmt.Sprintf("title-%d", i),
			Content:    fmt.Sprintf("content-%d", i),
			CreateTime: yu.NewStampTime(time.Now()),
			UpdateTime: yu.NewStampTime(time.Now()),
		}
	}
	gormInsert(list[:150])
	gormBatchInsert(list[150:])
}

func gormInsert(list []*Demo) {
	defer yu.Trace("gormInsert")()
	yu.Logf("gorm list len: %d", len(list))
	for i := 0; i < len(list); i++ {
		gdb.Insert("demo", list[i])
	}
}

func gormBatchInsert(list []*Demo) {
	defer yu.Trace("BatchInsert")()
	yu.Logf("gorm tx list len: %d", len(list))
	gdb.BatchInsert("demo", list, 0)
}

func queryAll() {
	if err := gdb.QueryRows(func(r *sql.Rows) error {
		d := &Demo{}
		if err := gdb.ScanRows(r, &d); err != nil {
			return err
		}
		yu.Logf("uid: %s", d.UID)
		return nil
	}, "select * from demo"); err != nil {
		yu.LogErr(err, "QueryRows")
	}
}
