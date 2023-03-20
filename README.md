![call](https://raw.githubusercontent.com/rytsh/call/pages/static/logo/call.svg)

[![License](https://img.shields.io/github/license/rytsh/call?color=red&style=flat-square)](https://raw.githubusercontent.com/rytsh/call/main/LICENSE)
[![Sonar Coverage](https://img.shields.io/sonar/coverage/rytsh_call?logo=sonarcloud&server=https%3A%2F%2Fsonarcloud.io&style=flat-square)](https://sonarcloud.io/summary/overall?id=rytsh_call)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/rytsh/call/test.yml?branch=main&logo=github&style=flat-square&label=ci)](https://github.com/rytsh/call/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/rytsh/call?style=flat-square)](https://goreportcard.com/report/github.com/rytsh/call)
[![Go PKG](https://raw.githubusercontent.com/rytsh/call/pages/assets/reference.svg)](https://pkg.go.dev/github.com/rytsh/call)
[![Web](https://img.shields.io/badge/web-document-blueviolet?style=flat-square)](https://rytsh.github.io/call/)

__Call__ dependency injection library based on registry arguments and functions.

```sh
go get github.com/rytsh/call
```

## Usage

First get new registry and add own functions and arguments.

Also you can add arguments in directly function.

Registry add function and arguments not in order you can add argument later.

```go
// create registry
reg := call.NewReg().
    AddArgument("a", 6).
    AddArgument("b", 2).
    AddFunction("divide", func(a, b int) (int, error) {
        if b == 0 {
            return 0, fmt.Errorf("divide by zero")
        }

        return a / b, nil
    }, "a", "b")

// call function
returns, err := reg.Call("divide")
if err != nil {
    fmt.Println(err)

    return
}

fmt.Println(returns[0], returns[1])
// Output:
// 3 <nil>
```

It is possible to call directly with arguments also callable as variadic.

```go
returns, err := reg.CallWithArgs("divide", "a", "b")
```

If function's argument length and type is not match, it will return error.

Check documentation for more details.
