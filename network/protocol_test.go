package network

import (
    "reflect"
    "testing"
)

func TestEncode(t *testing.T) {
	// Describe test cases
	tests := []struct {
		name    string
		message interface{}
		want    []byte
	}{
		{
			"Test RequestCS message encoding",
			RequestCS{ReqType, 28, 1},
			[]byte{0, 28, 0, 0, 0, 1},
		},
		{
			"Test ReleaseCS message encoding",
			ReleaseCS{OkType, 56, 3},
			[]byte{1, 56, 0, 0, 0, 3},
		},
		{
			"Test SetVariable message encoding",
			SetVariable{ValType, 456},
			[]byte{2, 0, 0, 1, 200},
		},
	}

	// Run the test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Encode(test.message)

			// Compare result with wanted result
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("encode() got %v, want %v", got, test.want)
			}
		})
	}
}

func TestDecodeRequest(t *testing.T) {
	// Describe test cases
	tests := []struct {
		name   string
		buffer []byte
		want   RequestCS
	}{
		{
			"Test decoding request message",
			[]byte{0, 28, 0, 0, 0, 12},
			RequestCS{ReqType, 28, 12},
		},
	}

	// Run test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := DecodeRequest(test.buffer)

			// Compare result with wanted result
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("DecodeRequest() got %v, want %v", got, test.want)
			}
		})
	}
}

func TestDecodeRelease(t *testing.T) {
	// Describe test cases
	tests := []struct {
		name   string
		buffer []byte
		want   ReleaseCS
	}{
		{
			"Test decoding release message",
			[]byte{1, 34, 0, 0, 0, 12},
			ReleaseCS{OkType, 34, 12},
		},
	}

	// Run test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := DecodeRelease(test.buffer)

			// Compare result with wanted result
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("DecodeRelease() got %v, want %v", got, test.want)
			}
		})
	}
}

func TestDecodeSetVariable(t *testing.T) {
	// Describe test cases
	tests := []struct {
		name   string
		buffer []byte
		want   SetVariable
	}{
		{
			"Test decoding SetVariable message",
			[]byte{2, 0, 0, 0, 12},
			SetVariable{ValType, 12},
		},
	}

	// Run test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := DecodeSetVariable(test.buffer)

			// Compare result with wanted result
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("DecodeRelease() got %v, want %v", got, test.want)
			}
		})
	}
}
