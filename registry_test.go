package call

import (
	"reflect"
	"testing"
)

func TestReg_Arguments(t *testing.T) {
	type addArgument struct {
		name string
		v    any
	}
	type getArgument struct {
		name   string
		Want   any
		WantOk bool
	}
	type deleteArgument struct {
		name string
	}
	type getArgumentNames struct {
		want []string
	}

	type argsTest struct {
		name             string
		addArgument      *addArgument
		getArgument      *getArgument
		deleteArgument   *deleteArgument
		getArgumentNames *getArgumentNames
	}
	tests := []struct {
		name  string
		tests []argsTest
	}{
		{
			name: "test args",
			tests: []argsTest{
				{
					name: "test add",
					addArgument: &addArgument{
						name: "test",
						v:    "test-1",
					},
				},
				{
					name: "test-2 add",
					addArgument: &addArgument{
						name: "test-2",
						v:    22,
					},
				},
				{
					name: "test get",
					getArgument: &getArgument{
						name:   "test",
						Want:   "test-1",
						WantOk: true,
					},
				},
				{
					name: "test-2 get",
					getArgument: &getArgument{
						name:   "test-2",
						Want:   22,
						WantOk: true,
					},
				},
				{
					name: "test delete",
					deleteArgument: &deleteArgument{
						name: "test",
					},
				},
				{
					name: "test get arguments",
					getArgument: &getArgument{
						name:   "test",
						Want:   nil,
						WantOk: false,
					},
				},
				{
					name: "test get arguments",
					getArgumentNames: &getArgumentNames{
						want: []string{"test-2"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReg()
			for _, test := range tt.tests {
				if test.addArgument != nil {
					r.AddArgument(test.addArgument.name, test.addArgument.v)
				}
				if test.getArgument != nil {
					got, ok := r.GetArgument(test.getArgument.name)
					if ok != test.getArgument.WantOk {
						t.Errorf("GetArgument() ok = %v, want %v", ok, test.getArgument.WantOk)
					}
					if !reflect.DeepEqual(got, test.getArgument.Want) {
						t.Errorf("GetArgument() got = %v, want %v", got, test.getArgument.Want)
					}
				}
				if test.deleteArgument != nil {
					r.DeleteArgument(test.deleteArgument.name)
				}
				if test.getArgumentNames != nil {
					got := r.GetArgumentNames()
					if !reflect.DeepEqual(got, test.getArgumentNames.want) {
						t.Errorf("GetArgumentNames() got = %v, want %v", got, test.getArgumentNames.want)
					}
				}
			}
		})
	}
}

func TestReg_Functions(t *testing.T) {
	type addFunction struct {
		name string
		fn   func([]any) any
		args []string
	}
	type getFunction struct {
		name   string
		Want   func([]any) Func
		WantOk bool
	}
	type deleteFunction struct {
		name string
	}
	type getFunctionNames struct {
		want []string
	}

	type functionTest struct {
		name             string
		addFunction      *addFunction
		getFunction      *getFunction
		deleteFunction   *deleteFunction
		getFunctionNames *getFunctionNames
	}
	tests := []struct {
		name      string
		tests     []functionTest
		functions []any
		panic     bool
	}{
		{
			name: "test functions",
			functions: []any{
				func(v string) string {
					return v
				},
			},
			tests: []functionTest{
				{
					name: "test add",
					addFunction: &addFunction{
						name: "test",
						fn: func(v []any) any {
							return v[0]
						},
						args: []string{"value"},
					},
				},
				{
					name: "test add 2",
					addFunction: &addFunction{
						name: "test-2",
						fn: func(v []any) any {
							return v[0]
						},
						args: []string{"value-2"},
					},
				},
				{
					name: "get function",
					getFunction: &getFunction{
						name: "test",
						Want: func(v []any) Func {
							return Func{
								Args: []string{"value"},
								Fn:   reflect.ValueOf(v[0]),
							}
						},
						WantOk: true,
					},
				},
				{
					name: "test delete",
					deleteFunction: &deleteFunction{
						name: "test",
					},
				},
				{
					name: "get function",
					getFunction: &getFunction{
						name: "test",
						Want: func([]any) Func {
							return Func{}
						},
						WantOk: false,
					},
				},
				{
					name: "get function names",
					getFunctionNames: &getFunctionNames{
						want: []string{"test-2"},
					},
				},
			},
		},
		{
			name: "panic add functions",
			functions: []any{
				"test",
				1234,
			},
			tests: []functionTest{
				{
					name: "wrong function type",
					addFunction: &addFunction{
						name: "test",
						fn: func(v []any) any {
							return v[1]
						},
					},
				},
			},
			panic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.panic {
						t.Errorf("Reg.Functions() panic = %v, want %v", r, tt.panic)
					}
				}
			}()

			r := NewReg()
			for _, test := range tt.tests {
				if test.addFunction != nil {
					r.AddFunction(test.addFunction.name, test.addFunction.fn(tt.functions), test.addFunction.args...)
				}
				if test.getFunction != nil {
					got, gotOk := r.GetFunction(test.getFunction.name)
					if !reflect.DeepEqual(got, test.getFunction.Want(tt.functions)) {
						t.Errorf("Reg.GetFunction() got = %v, want %v", got, test.getFunction.Want(tt.functions))
					}
					if gotOk != test.getFunction.WantOk {
						t.Errorf("Reg.GetFunction() gotOk = %v, want %v", gotOk, test.getFunction.WantOk)
					}
				}
				if test.deleteFunction != nil {
					r.DeleteFunction(test.deleteFunction.name)
				}
				if test.getFunctionNames != nil {
					if got := r.GetFunctionNames(); !reflect.DeepEqual(got, test.getFunctionNames.want) {
						t.Errorf("Reg.GetFunctionNames() = %v, want %v", got, test.getFunctionNames.want)
					}
				}
			}
		})
	}
}
