package server

import (
	"github.com/RaduBerinde/playground/test_dep_prototype/server/testingshim"
	"github.com/RaduBerinde/playground/test_dep_prototype/sql"
)

type Server struct {
	SQLSrvImpl *sql.SQLServer
}

func (s *Server) SQLSrv() interface{} {
	return s.SQLSrvImpl
}

var _ testingshim.TestServerInterface = &Server{}

func NewServer() *Server {
	return &Server{SQLSrvImpl: sql.NewSQLServer(5)}
}

type testServerFactoryImpl struct{}

// TestServerFactory can be passed to testingshim.InitTestServerFactory
var TestServerFactory testingshim.TestServerFactory = testServerFactoryImpl{}

// New is part of TestServerInterface.
func (testServerFactoryImpl) New() testingshim.TestServerInterface {
	return NewServer()
}
