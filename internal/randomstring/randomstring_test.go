package randomstring

import (
	"testing"
)

func TestRandomStringWithAlphabet(t *testing.T) {
	type args struct {
		n        int
		alphabet string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				n:        5,
				alphabet: "aaaaaaaaa",
			},
			want: "aaaaa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandomStringWithAlphabet(tt.args.n, tt.args.alphabet); got != tt.want {
				t.Errorf("RandomStringWithAlphabet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandomString(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				n: 5,
			},
			want: "q3TdG",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandomString(tt.args.n); got != tt.want {
				t.Errorf("RandomString() = %v, want %v", got, tt.want)
			}
		})
	}
}
