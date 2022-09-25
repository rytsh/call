package call

import (
	"reflect"
	"testing"
)

func testfunc() {}

func TestGetFunctionName(t *testing.T) {
	type args struct {
		fn reflect.Value
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "simple function",
			args: args{
				fn: reflect.ValueOf(func() {}),
			},
			want: "func1",
		},
		{
			name: "simple function 2",
			args: args{
				fn: reflect.ValueOf(func() {
					_ = func() {}
				}),
			},
			want: "func2",
		},
		{
			name: "testfunc function",
			args: args{
				fn: reflect.ValueOf(testfunc),
			},
			want: "testfunc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFunctionName(tt.args.fn); got != tt.want {
				t.Errorf("GetFunctionName() = %v, want %v", got, tt.want)
			}
		})
	}
}
