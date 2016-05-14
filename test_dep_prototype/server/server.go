package server

import "github.com/RaduBerinde/playground/test_dep_prototype/sql"

type Server struct {
	SqlSrv *sql.SQLServer
}

func NewServer() *Server {
	return &Server{SqlSrv: sql.NewSQLServer(5)}
}
