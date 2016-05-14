package sql_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/RaduBerinde/playground/test_dep_prototype/server"
	"github.com/RaduBerinde/playground/test_dep_prototype/server/testserver"
)

func TestSQL(t *testing.T) {
	s := server.NewServer()
	s.SqlSrv.Woof()
}

func TestMain(m *testing.M) {
	fmt.Printf("TestMain!\n")
	s := server.NewServer()
	testserver.TestSrvInstance.SqlSrv = s.SqlSrv
	code := m.Run()
	os.Exit(code)
}
