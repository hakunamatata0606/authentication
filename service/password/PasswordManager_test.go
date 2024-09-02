package password

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPasswordManager(t *testing.T) {
	manager := NewSha256Hash("secret")

	s := "aloha"
	hashed := manager.HashPassword(s)
	ok := manager.VerifyPassword(s, hashed)
	require.True(t, ok)
}

func testBasic(manager PasswordManager, base string, idx int, resultChan chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	s := base + strconv.Itoa(idx)
	hashed := manager.HashPassword(s)
	ok := manager.VerifyPassword(s, hashed)
	resultChan <- ok
}

func TestPasswordManagerConcurrent(t *testing.T) {
	resultChan := make(chan bool, 20)
	manager := NewSha256Hash("secret")
	max_run := 20
	var wg sync.WaitGroup
	for i := 0; i < max_run; i++ {
		wg.Add(1)
		go testBasic(manager, "alohalala", i, resultChan, &wg)
	}
	wg.Wait()
	for i := 0; i < max_run; i++ {
		res := <-resultChan
		require.True(t, res)
	}
}
