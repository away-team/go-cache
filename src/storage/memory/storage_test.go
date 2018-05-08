package memory

import (
	"reflect"
	"testing"
	"time"
)

func Test_Storage(t *testing.T) {
	s := NewStorage(time.Minute, time.Minute)

	type testcase struct {
		name  string
		val   interface{}
		ret   interface{}
		wait  time.Duration
		isErr bool
	}

	var str string
	var i int
	var i32 int32
	var i64 int64
	var b bool
	var f32 float32
	var f64 float64
	var complex testcase

	testcases := []testcase{
		{
			name: "string",
			val:  "foo",
			ret:  str,
		},
		{
			name: "integer",
			val:  123,
			ret:  i,
		},
		{
			name: "integer 32",
			val:  int32(1234),
			ret:  i32,
		},
		{
			name: "integer 64",
			val:  int64(1234),
			ret:  i64,
		},
		{
			name: "bool",
			val:  true,
			ret:  b,
		},
		{
			name: "float 32",
			val:  float32(1234.123),
			ret:  f32,
		},
		{
			name: "float 64",
			val:  float64(1234.123),
			ret:  f64,
		},
		{
			name: "complex",
			val:  testcase{name: "foo"},
			ret:  complex,
		},
		{
			name:  "wait",
			val:   testcase{name: "foo"},
			ret:   complex,
			wait:  time.Second * 2,
			isErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			key := tc.name
			err := s.Set(key, tc.val, time.Second)
			if err != nil {
				t.Fatalf("Error setting cache: %v", err)
			}

			time.Sleep(tc.wait)

			err = s.Get(key, &tc.ret)
			if !tc.isErr && err != nil {
				t.Fatalf("Error getting cache: %v", err)
			}

			if tc.isErr && err == nil {
				t.Fatalf("Expected error did not occur")
			}
			if !tc.isErr && !reflect.DeepEqual(tc.ret, tc.val) {
				t.Fatalf("Actual (%v) did not match expected (%v)", tc.ret, tc.val)
			}
		})
	}
}
