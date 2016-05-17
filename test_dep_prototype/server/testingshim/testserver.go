package testingshim

// TestServerInterface defines test server functionality that tests need.
type TestServerInterface interface {
	SQLSrv() interface{}
	// Other needed stuff goes here.
}

// TestServerFactory encompasses the actual implementation of the shim
// service.
type TestServerFactory interface {
	// New instantiates a test server instance.
	New() TestServerInterface
}

var serviceImpl TestServerFactory

// InitTestServerFactory should be called once to provide the implementation
// of the service. It will be called from a xx_test package that can import the
// server package.
func InitTestServerFactory(impl TestServerFactory) {
	serviceImpl = impl
}

func NewTestServer() TestServerInterface {
	return serviceImpl.New()
}
