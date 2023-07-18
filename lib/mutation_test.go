package passgen_test

import (
	"log"
	"os"
	"testing"

	"github.com/gtramontina/ooze"
)

func TestMutation(t *testing.T) {
	if os.PathSeparator == '\\' {
		log.Println("TestMutation skipped because ooze does not support windows")
		return
	}
	ooze.Release(t)
}
