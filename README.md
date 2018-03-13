go-ptimports
============
go-ptimports is a formatter for Go code. Functionally, it is almost identical to [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports).
The only difference is that, if the "refactor" flag is set to true, any non-CGo non-block imports are converted to block
imports.

For example:

```go
import "go/ast"
```

becomes

```go
import (
	"go/ast"
)
```

Import statements of the form `import "C"` (and any comments associated with such import statements) are preserved.

The output of `go-ptimports` is compliant with `goimports` and `gofmt`.
