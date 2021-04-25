package yu

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DSN 连接数据库路由参数
type DSN struct {
	Username string
	Password string
	URL      string
	Port     string
	DBName   string
	SkipTran bool
	PreStmt  bool
	LogLevel logger.LogLevel
}

// MySQL 生成MySQL的dsn
func (d DSN) MySQL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		d.Username, d.Password, d.URL, d.Port, d.DBName)
}

// Postgres 生成Postgres的dsn
func (d DSN) Postgres() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		d.Username, d.Password, d.URL, d.Port, d.DBName)
}

// GormDB 内嵌*gorm.DB，添加封装方法
type GormDB struct {
	*gorm.DB
}

// MustOpen 无错的连接数据库并返回*GormDB
func MustOpen(dial gorm.Dialector, d DSN) *GormDB {
	gdb, err := gorm.Open(
		dial,
		&gorm.Config{
			SkipDefaultTransaction: d.SkipTran,
			PrepareStmt:            d.PreStmt,
			Logger:                 logger.Default.LogMode(d.LogLevel),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	db, err := gdb.DB()
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(60 * time.Second)
	db.SetMaxIdleConns(30)
	db.SetMaxOpenConns(30)
	return &GormDB{gdb}
}

// MustOpenMySQL 无错的连接MySQL返回*gorm.DB
func MustOpenMySQL(d DSN) *GormDB {
	return MustOpen(mysql.Open(d.MySQL()), d)
}

// MustOpenPostgres 无错的连接Postgres返回*gorm.DB
func MustOpenPostgres(d DSN) *GormDB {
	return MustOpen(postgres.Open(d.Postgres()), d)
}

// CloseDB 关闭数据库连接
func (gdb *GormDB) CloseDB() {
	db, err := gdb.DB.DB()
	if err != nil {
		return
	}
	db.Close()
}

// ExecSQL 直接执行sql
func (db *GormDB) ExecSQL(sql string, values ...interface{}) (int64, error) {
	tx := db.Exec(sql, values...)
	return tx.RowsAffected, tx.Error
}

// Query 查询数据
func (db *GormDB) Query(dest interface{}, stmt string, values ...interface{}) error {
	tx := db.Raw(stmt, values...).Scan(dest)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// QueryRow 获取*sql.Row
func (db *GormDB) QueryRow(stmt string, values ...interface{}) *sql.Row {
	return db.Raw(stmt, values...).Row()
}

// Insert 插入数据
func (db *GormDB) Insert(tb string, value interface{}) error {
	return db.Table(tb).Create(value).Error
}

// BatchInsert 批量插入数据
func (db *GormDB) BatchInsert(tb string, value interface{}, batchSize int) error {
	if batchSize == 0 || batchSize > 2500 {
		batchSize = 2500
	}
	return db.Table(tb).CreateInBatches(value, batchSize).Error
}

// QueryRows 批量查询使用rows行扫描处理
func (db *GormDB) QueryRows(hand func(*sql.Rows) error, sql string, values ...interface{}) (err error) {
	rows, err := db.Raw(sql, values...).Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = hand(rows); err != nil {
			return
		}
	}
	return
}
