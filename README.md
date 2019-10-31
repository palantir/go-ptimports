<p align="right">
<a href="https://autorelease.general.dmz.palantir.tech/palantir/go-ptimports"><img src="https://img.shields.io/badge/Perform%20an-Autorelease-success.svg" alt="Autorelease"></a>
</p>

go-ptimports
============
go-ptimports is a formatter for Go code. Functionally, it is almost identical to [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports).
However, it provides the following extra functionality:

* When the `-s` flag is specified, code is simplified in the same manner as `gofmt -s`

For example, a range statement such as:

```go
range i, _ := arr {
    _ = i
}
```

becomes:

```go
range i := arr {
    _ = i
}
```

* When the `-r` flag is specified, any non-CGo non-block imports are converted to block imports

For example:

```go
import "go/ast"
```

becomes:

```go
import (
	"go/ast"
)
```

Import statements of the form `import "C"` (and any comments associated with such import statements) are preserved.

The output of go-ptimports is compliant with goimports and gofmt.

Runing go-ptimports without any extra flags matches the behavior of goimports. Running go-ptimports with the
`--format-only` flag is roughly equivalent to the behavior of gofmt.
