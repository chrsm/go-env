package env_test

import (
	"os"
	"testing"

	"github.com/chrsm/env"
)

func TestString(t *testing.T) {
	type structtype struct {
		Val string `env:"STRUCT_TEST"`
	}

	os.Clearenv()
	os.Setenv("STRUCT_TEST", "test!!")

	x := &structtype{}
	err := env.Decode(x)
	if err != nil {
		t.Errorf("error in env.Decode: %s", err)
	}

	if x.Val != "test!!" {
		t.Errorf("Expected to set st.Val, got %s", x.Val)
	}
}

func TestMap(t *testing.T) {
	type structtype struct {
		Val map[string]string `env:"STRUCT_MAP"`
	}

	os.Clearenv()
	os.Setenv("STRUCT_MAP", "a=b,x=y,1=2,!!!!=n")

	{
		x := &structtype{}
		err := env.Decode(x)
		if err != nil {
			t.Errorf("error in env.Decode: %s", err)
		}

		if x.Val["a"] != "b" || x.Val["x"] != "y" || x.Val["1"] != "2" || x.Val["!!!!"] != "n" {
			t.Errorf("env.Decode failed to set map, got %#v", x.Val)
		}
	}

	// ensure things convertible to map<string,string> are allowed still, despite being a distinct type
	type mapstr map[string]string
	type structtype_convertible struct {
		Val mapstr `env:"STRUCT_MAP"`
	}
	{
		x := &structtype{}
		err := env.Decode(x)
		if err != nil {
			t.Errorf("error in env.Decode: %s", err)
		}

		if x.Val["a"] != "b" || x.Val["x"] != "y" || x.Val["1"] != "2" || x.Val["!!!!"] != "n" {
			t.Errorf("env.Decode failed to set map, got %#v", x.Val)
		}
	}

	// ensure we can't parse an incorrect map type..
	type invalidstructtype struct {
		Val map[string]int `env:"STRUCT_MAP"`
	}

	{
		x := &invalidstructtype{}
		err := env.Decode(x)
		if err == nil {
			t.Errorf("error in env.Decode: expected to be unable to decode to unsupported map type")
		}
	}
}

func TestIntegers(t *testing.T) {
	type structtype struct {
		Int   int   `env:"V_INT"`
		Int32 int32 `env:"V_INT32"`
		Int64 int64 `env:"V_INT64"`

		Uint   uint   `env:"V_UINT"`
		Uint32 uint32 `env:"V_UINT32"`
		Uint64 uint64 `env:"V_UINT64"`
	}

	os.Clearenv()
	os.Setenv("V_INT", "255")
	os.Setenv("V_INT32", "392930")
	os.Setenv("V_INT64", "39292309209302")
	os.Setenv("V_UINT", "255")
	os.Setenv("V_UINT32", "392930")
	os.Setenv("V_UINT64", "39292309209302")

	x := &structtype{}
	err := env.Decode(x)
	if err != nil {
		t.Errorf("error in env.Decode: %s", err)
	}

	if x.Int != 255 || x.Uint != 255 {
		t.Errorf("Expected Int and Uint to be 255, got %vi - %vu", x.Int, x.Uint)
	}

	if x.Int32 != 392930 || x.Uint32 != 392930 {
		t.Errorf("Expected Int32 and Uint32 to be 392930, got %vi - %vu", x.Int32, x.Uint32)
	}

	if x.Int64 != 39292309209302 || x.Uint64 != 39292309209302 {
		t.Errorf("Expected Int64 and Uint64 to be 39292309209302, got %vi - %vu", x.Int64, x.Uint64)
	}
}

func TestSlice(t *testing.T) {
	type structtype struct {
		Val []string `env:"STRUCT_TEST"`
	}

	os.Clearenv()
	os.Setenv("STRUCT_TEST", "test!!,no,really")

	{
		x := &structtype{}
		err := env.Decode(x)
		if err != nil {
			t.Errorf("error in env.Decode: %s", err)
		}

		if x.Val[0] != "test!!" || x.Val[1] != "no" || x.Val[2] != "really" {
			t.Errorf("Expected to set st.Val, got %s", x.Val)
		}
	}

	type structtype_int struct {
		Val []int `env:"STRUCT_TEST"`
	}

	os.Setenv("STRUCT_TEST", "1,2,9,102,0")
	{
		x := &structtype_int{}
		err := env.Decode(x)
		if err != nil {
			t.Errorf("error in env.Decode: %s", err)
		}

		if x.Val[0] != 1 || x.Val[1] != 2 || x.Val[2] != 9 || x.Val[3] != 102 || x.Val[4] != 0 {
			t.Errorf("Expected to set st.Val, got %v", x.Val)
		}
	}
}

func TestEmbedded(t *testing.T) {
	type Structptr struct {
		PtrVal string `env:"PTR_VAL"`
	}

	type structtype struct {
		Val string `env:"STRUCT_TEST"`

		Embedded struct {
			EmbeddedVal string `env:"EMBEDDED_TEST"`
		}

		Ptr *Structptr
	}

	os.Clearenv()
	os.Setenv("STRUCT_TEST", "test!!")
	os.Setenv("EMBEDDED_TEST", "embedded text")
	os.Setenv("PTR_VAL", "ptr value text")

	x := &structtype{}
	err := env.Decode(x)
	if err != nil {
		t.Errorf("error in env.Decode: %s", err)
	}

	if x.Val != "test!!" {
		t.Errorf("Expected to set st.Val, got %s", x.Val)
	}

	if x.Embedded.EmbeddedVal != "embedded text" {
		t.Errorf("Expected to set x.Embedded.Val, got %s", x.Embedded.EmbeddedVal)
	}

	if x.Ptr.PtrVal != "ptr value text" {
		t.Errorf("Expected to set x.Ptr.PtrVal, got %s", x.Ptr.PtrVal)
	}
}
