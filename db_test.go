package clerk

import (
	"math/rand"
	"reflect"
	"testing"
)

var (
	testFilename = "test.gob"
	testDB       *DB
	testSource   map[int]int
)

func reset() {
	testSource = make(map[int]int)
	var err error
	testDB, err = New(testFilename, &testSource)
	if err != nil {
		panic(err)
	}
}

func TestDB(t *testing.T) {
	reset()
	defer testDB.Remove()

	// Save random data in file.
	for i := 0; i < 10; i++ {
		testSource[rand.Int()] = rand.Int()
	}
	if err := testDB.Save(); err != nil {
		panic(err)
	}

	// Keep old data for comparison and reopen file.
	oldTestSource := make(map[int]int)
	for k, v := range testSource {
		oldTestSource[k] = v
	}
	reset()

	if !reflect.DeepEqual(oldTestSource, testSource) {
		t.Fail()
	}
	if err := testDB.Remove(); err != nil {
		panic(err)
	}
}

//   1,000 entries   18 KB   0.45 ms/op
//  10,000 entries  180 KB   2.43 ms/op
// 100,000 entries  1.8 MB  22.65 ms/op
func BenchmarkSave(b *testing.B) {
	defer testDB.Remove()
	for i := 0; i < b.N; i++ {
		testDB.Save()
	}
}
