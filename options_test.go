package call

import (
	"fmt"
	"reflect"
	"testing"
)

func TestOptions_VisitOptions(t *testing.T) {
	type fields struct {
		option map[string]func([]reflect.Value, ...string) ([]reflect.Value, error)
	}
	type args struct {
		arg string
		v   any
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       []reflect.Value
		wantErr    bool
		wantErrStr string
	}{
		{
			name: "option not found",
			fields: fields{
				option: map[string]func([]reflect.Value, ...string) ([]reflect.Value, error){
					"option1": func(v []reflect.Value, args ...string) ([]reflect.Value, error) {
						return nil, nil
					},
				},
			},
			args: args{
				arg: "test:option2",
				v:   []string{"test"},
			},
			want:       nil,
			wantErr:    true,
			wantErrStr: `option option2 not found`,
		},
		{
			name: "option error",
			fields: fields{
				option: map[string]func([]reflect.Value, ...string) ([]reflect.Value, error){
					"option1": func(v []reflect.Value, args ...string) ([]reflect.Value, error) {
						return nil, fmt.Errorf("test error")
					},
				},
			},
			args: args{
				arg: "test:option1",
				v:   []string{"test"},
			},
			want:       nil,
			wantErr:    true,
			wantErrStr: `option1; test error`,
		},
		{
			name: "value is nil",
			fields: fields{
				option: map[string]func([]reflect.Value, ...string) ([]reflect.Value, error){
					"option1": func(v []reflect.Value, args ...string) ([]reflect.Value, error) {
						return nil, nil
					},
				},
			},
			args: args{
				arg: "test:option1",
				v:   []string{"test"},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{
				option: tt.fields.option,
			}
			got, err := o.VisitOptions(tt.args.arg, tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Options.VisitOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.wantErrStr {
				t.Errorf("Options.VisitOptions() error = %v, wantErrStr %v", err, tt.wantErrStr)
				return
			}
			if tt.wantErr == false {
				for i := range got {
					if !reflect.DeepEqual(got[i].Interface(), tt.want[i].Interface()) {
						t.Errorf("Options.VisitOptions() = %v, want %v", got, tt.want)
					}
				}
			}
		})
	}
}
