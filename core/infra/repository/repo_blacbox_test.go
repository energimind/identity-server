package repository_test

import (
	"testing"

	"github.com/energimind/go-kit/testutil/mongodb"
)

var mongoEnv mongodb.MongoEnvironment

// TestMain sets up the MongoDB test environment for all blackbox
// tests in the repository_test package.
func TestMain(m *testing.M) {
	cleanUp, err := mongoEnv.Start()
	defer cleanUp()

	if err != nil {
		panic(err)
	}

	m.Run()
}
