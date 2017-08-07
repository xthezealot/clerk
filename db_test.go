package clerk

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"testing"
)

var (
	testFilename = "test.gob"
	testDB       = new(struct {
		DB
		Data map[int]int
	})
)

func init() {
	Init(testFilename, testDB)
}

func TestDB(t *testing.T) {
	defer func() {
		if err := testDB.Remove(); err != nil {
			panic(err)
		}
	}()
	testDB.Lock()
	defer testDB.Unlock()

	// Save random data in file.
	testDB.Data = make(map[int]int)
	for i := 0; i < 10; i++ {
		testDB.Data[rand.Int()] = rand.Int()
	}

	if err := testDB.Save(); err != nil {
		panic(err)
	}

	// Keep old data for comparison and reopen file.
	oldInts := make(map[int]int)
	for k, v := range testDB.Data {
		oldInts[k] = v
	}
	testDB.Data = nil
	if err := testDB.Rebase(); err != nil {
		panic(err)
	}

	if !reflect.DeepEqual(oldInts, testDB.Data) {
		t.Errorf("want: %v\ngot: %v", oldInts, testDB.Data)
	}
}

func benchmarkDB(n int) func(*testing.B) {
	return func(b *testing.B) {
		testDB.Data = make(map[int]int)
		for i := 0; i < n; i++ {
			testDB.Data[rand.Int()] = rand.Int()
		}

		for i := 0; i < b.N; i++ {
			testDB.Lock()
			if err := testDB.Save(); err != nil {
				panic(err)
			}
			if err := testDB.Rebase(); err != nil {
				panic(err)
			}
			testDB.Unlock()
		}
	}
}

func benchmarkPrintDBFileSize(name string) {
	fi, _ := os.Stat(testFilename)
	fmt.Printf("File size for %s: %d bytes\n", name, fi.Size())
}

//   1,000 entries   18 KB   0.9 ms/op
//  10,000 entries  180 KB   6.1 ms/op
// 100,000 entries  1.8 MB  58.4 ms/op
func BenchmarkDB(b *testing.B) {
	defer testDB.Remove()
	for _, n := range []int{1000, 10000, 100000} {
		name := fmt.Sprintf("%d entries", n)
		b.Run(name, benchmarkDB(n))
		benchmarkPrintDBFileSize(name)
	}
}
