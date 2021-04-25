package yu

import (
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

// StmtDB 嵌入sqlx.DB，加了statement缓存
type StmtDB struct {
	*sqlx.DB
	stmts  map[string]*sqlx.Stmt
	nstmts map[string]*sqlx.NamedStmt
	mutex  *sync.Mutex
}

// NewStmtDB 创建StmtDB
func NewStmtDB(db *sqlx.DB) *StmtDB {
	db.SetConnMaxLifetime(60 * time.Second)
	db.SetMaxIdleConns(30)
	db.SetMaxOpenConns(30)
	return &StmtDB{
		DB:     db,
		stmts:  make(map[string]*sqlx.Stmt, 0),
		nstmts: make(map[string]*sqlx.NamedStmt, 0),
		mutex:  new(sync.Mutex),
	}
}

// CloseDB 关闭StmtDB
func (s *StmtDB) CloseDB() {
	if s.stmts != nil {
		for _, stmt := range s.stmts {
			if stmt != nil {
				stmt.Close()
			}
		}
	}
	if s.nstmts != nil {
		for _, nstmt := range s.nstmts {
			if nstmt != nil {
				nstmt.Close()
			}
		}
	}
	s.Close()
}

// Stmt 懒加载的方式缓存sqlx.Stmt
func (s *StmtDB) Stmt(sqlStr string) (stmt *sqlx.Stmt, err error) {
	name := string(MD5(sqlStr))
	s.mutex.Lock()
	stmt, ok := s.stmts[name]
	s.mutex.Unlock()
	if ok {
		return
	}
	stmt, err = s.Preparex(sqlStr)
	if err != nil {
		LogErr(err, "Preparex")
		return
	}
	s.mutex.Lock()
	s.stmts[name] = stmt
	s.mutex.Unlock()
	return
}

// NStmt g懒加载的方式缓存sqlx.NamedStmt
func (s *StmtDB) NStmt(sqlStr string) (stmt *sqlx.NamedStmt, err error) {
	name := string(MD5(sqlStr))
	s.mutex.Lock()
	stmt, ok := s.nstmts[name]
	s.mutex.Unlock()
	if ok {
		return
	}
	stmt, err = s.PrepareNamed(sqlStr)
	if err != nil {
		LogErr(err, "Preparex")
		return
	}
	s.mutex.Lock()
	s.nstmts[name] = stmt
	s.mutex.Unlock()
	return
}

// Lock 使用锁
func (s *StmtDB) Lock() {
	s.mutex.Lock()
}

// Unlock 释放锁
func (s *StmtDB) Unlock() {
	s.mutex.Unlock()
}

// DeferLock 合并锁的使用
// example: defer s.DeferLock()()
func (s *StmtDB) DeferLock() func() {
	s.mutex.Lock()
	return s.mutex.Unlock
}
