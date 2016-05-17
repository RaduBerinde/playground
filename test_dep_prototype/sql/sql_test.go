package sql_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/RaduBerinde/playground/test_dep_prototype/server"
	"github.com/RaduBerinde/playground/test_dep_prototype/server/testingshim"
)

// This is the kind of test we want to be able to run form the sql package.
func TestSQL(t *testing.T) {
	s := server.NewServer()
	s.SQLSrvImpl.Woof()
}

func TestMain(m *testing.M) {
	fmt.Printf("TestMain!\n")
	testingshim.InitTestServerFactory(server.TestServerFactory)
	code := m.Run()
	os.Exit(code)
}
