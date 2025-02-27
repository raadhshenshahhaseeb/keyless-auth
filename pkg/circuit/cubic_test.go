package circuit

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("basic circuit", func(t *testing.T) {
		newCubic, err := New(&Circuit{
			X: 2,
			Y: 8,
		})

		if err != nil {
			t.Fatal(err)
		}

		if newCubic == nil {
			t.Fatal("nil cubic object")
		}

		fmt.Println("basic circuit passed")
	})
}
