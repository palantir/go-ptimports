go-ptimports
============
go-ptimports is a formatter for Go code. It applies the formatting operation performed by [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports)
and also performs the following operations:

* Converts import statements into the group format

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

* Groups import statements between stdlib packages, non-stdlib packages and packages that are part of the project.

For example, a file in `github.com/palantir/go-ptimports` that imports `testing`, `github.com/stretchr/testify/assert`
and `github.com/palantir/go-ptimports/ptimports` would organize its imports as follows:

```go
import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/palantir/go-ptimports/ptimports"
)
```

The output of `go-ptimports` is compliant with `goimports` and `gofmt`.
