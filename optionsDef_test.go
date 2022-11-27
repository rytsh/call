package call

import (
	"reflect"
	"testing"
)

func TestOptionGetIndex(t *testing.T) {
	type args struct {
		v    []reflect.Value
		args []string
	}
	tests := []struct {
		name       string
		args       args
		want       []reflect.Value
		wantErr    bool
		wantErrStr string
	}{
		{
			name: "no value",
			args: args{
				v:    []reflect.Value{},
				args: []string{"1"},
			},
			want:       nil,
			wantErr:    true,
			wantErrStr: "no value",
		},
		{
			name: "index is empty",
			args: args{
				v:    []reflect.Value{reflect.ValueOf("test")},
				args: nil,
			},
			want:       nil,
			wantErr:    true,
			wantErrStr: "index is empty",
		},
		{
			name: "not related type for index",
			args: args{
				v:    []reflect.Value{reflect.ValueOf("test")},
				args: []string{"1"},
			},
			want:       nil,
			wantErr:    true,
			wantErrStr: "not related type for index",
		},
		{
			name: "index is not a number",
			args: args{
				v:    []reflect.Value{reflect.ValueOf([]string{"test"})},
				args: []string{"test"},
			},
			want:       nil,
			wantErr:    true,
			wantErrStr: `index is not a number; strconv.Atoi: parsing "test": invalid syntax`,
		},
		{
			name: "index out of range",
			args: args{
				v:    []reflect.Value{reflect.ValueOf([]string{"test"})},
				args: []string{"20"},
			},
			want:       nil,
			wantErr:    true,
			wantErrStr: `index out of range`,
		},
		{
			name: "correct value",
			args: args{
				v:    []reflect.Value{reflect.ValueOf([]string{"test"})},
				args: []string{"0"},
			},
			want: []reflect.Value{reflect.ValueOf("test")},
		},
		{
			name: "correct value",
			args: args{
				v:    []reflect.Value{reflect.ValueOf([]any{"test", 123, 123.123})},
				args: []string{"0", "2"},
			},
			want: []reflect.Value{reflect.ValueOf("test"), reflect.ValueOf(123.123)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OptionGetIndex(tt.args.v, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("OptionGetIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.wantErrStr {
				t.Errorf("OptionGetIndex() error = %v, wantErrStr %v", err, tt.wantErrStr)
				return
			}
			if tt.wantErr == false {
				for i := range got {
					if !reflect.DeepEqual(got[i].Interface(), tt.want[i].Interface()) {
						t.Errorf("OptionGetIndex() = %v, want %v", got, tt.want)
					}
				}
			}
		})
	}
}

func TestOptionVariadic(t *testing.T) {
	type args struct {
		v []reflect.Value
	}
	tests := []struct {
		name       string
		args       args
		want       []reflect.Value
		wantErr    bool
		wantErrStr string
	}{
		{
			name: "not related type for variadic",
			args: args{
				v: nil,
			},
			want:       nil,
			wantErr:    true,
			wantErrStr: "no value",
		},
		{
			name: "not related type for variadic",
			args: args{
				v: []reflect.Value{reflect.ValueOf("test")},
			},
			want:       nil,
			wantErr:    true,
			wantErrStr: "not related type for variadic",
		},
		{
			name: "return directly",
			args: args{
				v: []reflect.Value{reflect.ValueOf("test"), reflect.ValueOf("test2")},
			},
			want: []reflect.Value{reflect.ValueOf("test"), reflect.ValueOf("test2")},
		},
		{
			name: "works well",
			args: args{
				v: []reflect.Value{reflect.ValueOf([]any{"test", 123, 123.123})},
			},
			want: []reflect.Value{reflect.ValueOf("test"), reflect.ValueOf(123), reflect.ValueOf(123.123)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OptionVariadic(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("OptionVariadic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.wantErrStr {
				t.Errorf("OptionVariadic() error = %v, wantErrStr %v", err, tt.wantErrStr)
				return
			}
			if tt.wantErr == false {
				for i := range got {
					if !reflect.DeepEqual(got[i].Interface(), tt.want[i].Interface()) {
						t.Errorf("OptionVariadic() = %v, want %v", got, tt.want)
					}
				}
			}
		})
	}
}
