package masker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMasker(t *testing.T) {

	t.Run("Masker test", func(t *testing.T) {

		passwords := [11]string{"", "1", "12", "123", "1234", "12345", "123456", "1234567", "12345678", "123456789", "1234567890"}

		tests := []struct {
			name  string
			value string
			want  string
		}{
			{"test #1",
				passwords[0],
				""},
			{"test #2",
				passwords[1],
				"*"},
			{"test #3",
				passwords[2],
				"1*"},
			{"test #4",
				passwords[3],
				"1**"},
			{"test 5",
				passwords[4],
				"12**"},
			{"test #6",
				passwords[5],
				"12***"},
			{"test #7",
				passwords[6],
				"123***"},
			{"test #8",
				passwords[7],
				"123****"},
			{"test #9",
				passwords[8],
				"1234****"},
			{"test #10",
				passwords[9],
				"1234*****"},
			{"test #11",
				passwords[10],
				"12345*****"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				res := Masker(tt.value)
				assert.Equal(t, tt.want, res)
			})
		}
	})
}
