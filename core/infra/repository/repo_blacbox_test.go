package repository_test

import (
	"testing"

	"github.com/energimind/identity-service/test/utils"
)

var mongoEnv utils.MongoEnvironment

// TestMain sets up the MongoDB test environment for all blackbox
// tests in the repository_test package.
func TestMain(m *testing.M) {
	defer mongoEnv.Start()()

	m.Run()
}
