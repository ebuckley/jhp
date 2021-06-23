package jhp

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dop251/goja"
	"io"
	"log"
	"net/http"
	"strings"
)

type RequestReader struct {
	DB *sql.DB
	Params map[string]interface{}
}

type ResponseWriter interface {
	io.Writer
}

func ToParamMap(r *http.Request) map[string]interface{} {

	params := make(map[string]interface{})
	for s, v := range r.URL.Query() {
		if len(v) == 1 {
			params[s] = v[0]
		} else {
			params[s] = v
		}
	}

	for s, v := range r.Form {
		if len(v) == 1 {
			params[s] = v[0]
		} else {
			params[s] = v
		}
	}
	return params
}

// Register will setup the vm for a request context
func Register(vm *goja.Runtime, r RequestReader, w ResponseWriter) error {

	err := vm.Set("params", vm.ToValue(r.Params))
	if err != nil {
		return err
	}

	// register echo function
	err = vm.Set("echo", func(fn goja.FunctionCall) goja.Value {
		val := fn.Argument(0)
		fmt.Fprint(w, val)
		return vm.ToValue(true)
	})
	if err != nil {
		return err
	}
	// register a way to query the database
	err = vm.Set("sql", func(fn goja.FunctionCall) goja.Value {
		query := fn.Argument(0).String()
		queryResult, err := r.DB.Query(query)
		if err != nil {
			// how to throw??
			log.Println("Error with sql: ", query, ":", err)
			return goja.Null()
		}
		rows := make([]map[string]interface{}, 0)

		cols, err := queryResult.Columns()
		if err != nil {
			log.Println("error with SQL fetching columns: ", err)
			return goja.Null()
		}
		for queryResult.Next() {
			columns := make([]interface{}, len(cols))
			columnPointers := make([]interface{}, len(cols))
			for i, _ := range columns {
				columnPointers[i] = &columns[i]
			}
			err = queryResult.Scan(columnPointers...)
			if err != nil {
				// how to throw??
				log.Println("Error with sql: ", query, ":", err)
				return goja.Null()
			}
			m := make(map[string]interface{})
			for i, colName := range cols {
				val := columnPointers[i].(*interface{})
				m[colName] = *val
			}
			rows = append(rows, m)
		}
		return vm.ToValue(rows)
	})

	if err != nil {
		return err
	}

	// register parameter fetching code

	return nil
}

func Render(vm *goja.Runtime,  r RequestReader, w ResponseWriter, content string) error {
	err := Register(vm, r, w)
	if err != nil {
		return err
	}

	_, err = vm.RunString(toEcho(content))

	return err
}

// turn a jhp page into a runnable expression
func toEcho(in string) string {
	if len(in) == 0 {
		return ""
	}
	// start with echo
	out := "echo("

	// search for the opening <?jhp
	idx := strings.Index(in, "<?jhp")
	if idx == -1 {
		return out + toJSString(in) + ");"
	}
	// TODO escape the bit in between
	out += toJSString(in[:idx]) + ");\n"
	idx += 5 // skip <?jhp
	in = in[idx:]

	// search for the close ?>
	idx = strings.IndexAny(in, "?>")
	if idx == -1 {
		return in
	}
	out += in[:idx]
	// the remaining can be continued with the same old logic
	return out + toEcho(in[idx + 2:]) // +2 to exclude ?> on the end here
}

func toJSString(in string) string {
	o, _ := json.Marshal(in)
	return string(o)
}