package call

import (
	"reflect"
	"testing"
)

func TestReg_CallWithArgs(t *testing.T) {
	type args struct {
		name string
		args []string
	}
	tests := []struct {
		name       string
		args       args
		modify     func(*Reg)
		want       []any
		wantErr    bool
		wantErrStr string
	}{
		{
			name: "function test not found",
			args: args{
				name: "test",
				args: []string{"test-1", "test-2"},
			},
			want:       nil,
			wantErr:    true,
			wantErrStr: "function test not found",
		},
		{
			name: "argument test-1 not found",
			modify: func(r *Reg) {
				r.AddFunction("test", func() {})
			},
			args: args{
				name: "test",
				args: []string{"test-1", "test-2"},
			},
			want:       nil,
			wantErr:    true,
			wantErrStr: "argument test-1 not found",
		},
		{
			name: "optionX not found",
			modify: func(r *Reg) {
				r.AddFunction("test", func() {})
				r.AddArgument("test-1", "test-1")
			},
			args: args{
				name: "test",
				args: []string{"test-1:optionX"},
			},
			want:       nil,
			wantErr:    true,
			wantErrStr: "failed VisitOption option optionX not found",
		},
		{
			name: "not enough arguments",
			modify: func(r *Reg) {
				r.AddFunction("test", func(string, float64, ...int) {})
				r.AddArgument("test-1", "test-1")
			},
			args: args{
				name: "test",
				args: []string{"test-1"},
			},
			want:       nil,
			wantErr:    true,
			wantErrStr: "not enough arguments",
		},
		{
			name: "argument count mismatch",
			modify: func(r *Reg) {
				r.AddFunction("test", func(string, float64) {})
				r.AddArgument("test-1", "test-1")
			},
			args: args{
				name: "test",
				args: []string{"test-1"},
			},
			want:       nil,
			wantErr:    true,
			wantErrStr: "argument count mismatch",
		},
		{
			name: "variadic argument type check",
			modify: func(r *Reg) {
				r.AddFunction("test", func(string, string, string, ...float64) {})
				r.AddArgument("test-1", "test-1")
			},
			args: args{
				name: "test",
				args: []string{"test-1", "test-1", "test-1", "test-1"},
			},
			want:       nil,
			wantErr:    true,
			wantErrStr: "variadic function: index 3 argument string type mismatch with function float64 type",
		},
		{
			name: "argument type check",
			modify: func(r *Reg) {
				r.AddFunction("test", func(string, string, string) {})
				r.AddArgument("test-1", "test-1").AddArgument("test-2", 22)
			},
			args: args{
				name: "test",
				args: []string{"test-1", "test-1", "test-2"},
			},
			want:       nil,
			wantErr:    true,
			wantErrStr: "function: index 2 argument int type mismatch with function string type",
		},
		{
			name: "argument nil",
			modify: func(r *Reg) {
				r.AddFunction("test", func(v string, _ interface{}) string { return v })
				r.AddArgument("test-1", "test-1").AddArgument("test-2", nil)
			},
			args: args{
				name: "test",
				args: []string{"test-1", "test-2"},
			},
			want: []any{
				"test-1",
			},
		},
		{
			name: "variadic functions nil",
			modify: func(r *Reg) {
				r.AddFunction("test", func(v string, _ ...interface{}) string { return v })
				r.AddArgument("test-1", "test-1").AddArgument("test-2", nil)
			},
			args: args{
				name: "test",
				args: []string{"test-1", "test-2"},
			},
			want: []any{
				"test-1",
			},
		},
		{
			name: "variadic functions nil without value",
			modify: func(r *Reg) {
				r.AddFunction("test", func(v string, _ ...interface{}) string { return v })
				r.AddArgument("test-1", "test-1")
			},
			args: args{
				name: "test",
				args: []string{"test-1"},
			},
			want: []any{
				"test-1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReg()
			if tt.modify != nil {
				tt.modify(r)
			}

			got, err := r.CallWithArgs(tt.args.name, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reg.CallWithArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.wantErrStr {
				t.Errorf("Reg.CallWithArgs() error = %v, wantErrStr %v", err, tt.wantErrStr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reg.CallWithArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
