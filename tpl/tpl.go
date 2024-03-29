package tpl

var Header = `// !!!Don't Edit it!!!
package {{ .PackageName }}

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/go-comm/sqlxmodel"
	"github.com/jmoiron/sqlx"
)

var (
	_ context.Context
	_ sql.DB
	_ sqlx.DB
	_ sqlxmodel.SqlxModel

	_ = strings.Join
	_ = fmt.Println
	_ = log.Println
	_ = reflect.ValueOf
)
`

var Model = `
// {{ .Name | Title }}Model model of {{ .Name }}
//
// !!!Don't Edit it!!!
var {{ .Name | Title }}Model = new({{ .Name }})
`
