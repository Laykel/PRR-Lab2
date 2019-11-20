package network

import (
    "reflect"
    "testing"
)

func Test_encode(t *testing.T) {
	// Describe test cases
	tests := []struct {
		name    string
		message interface{}
		want    []byte
	}{
		{
			"Test RequestCS message encoding",
			RequestCS{28, 1},
			[]byte{0, 0, 0, 28, 0, 0, 0, 1},
		},
		{
			"Test ReleaseCS message encoding",
			ReleaseCS{56, 3, 45780},
			[]byte{0, 0, 0, 56, 0, 0, 0, 3, 0, 0, 178, 212},
		},
	}

	// Run the test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := encode(test.message)

			// Compare result with wanted result
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("encode() got %v, want %v", got, test.want)
			}
		})
	}
}
