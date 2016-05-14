package sql

import "fmt"

type SQLServer struct {
	lol int
}

func NewSQLServer(x int) *SQLServer {
	return &SQLServer{lol: x}
}

func (s *SQLServer) Woof() {
	fmt.Printf("Woof %d\n", s.lol)
}
