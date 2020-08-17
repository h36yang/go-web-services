# go-web-services
Practice materials for "Creating Web Services with Go" course

### Why web services in Go?
* Small fast binaries
* Full featured standard library
* Concurrency
* Easy to learn

### Handling HTTP Requests
* Creating Basic Handlers via `net/http` package
  * 2 ways to handle HTTP requests:
    * `http.Handle(pattern string, handler Handler)`
    * `http.HandleFunc(pattern string, handler func(ResponseWriter, *Request))`
  * **ServeMux**: HTTP Request Multiplexor
    * `http.ListenAndServe(addr string, handler Handler) error`
    * `http.ListenAndServeTLS(addr, certFile, keyFile string, handler Handler) error` - less commonly used as we usually use a Load Balancer to handle the TLS encryption
* Working with JSON via `encoding/json` package
  * Marshal: `json.Marshal(v interface{}) ([]byte, error)` - only exported fields will be marshaled
  * Unmarshal: `json.Unmarshal(data []byte, v interface{}) (error)`
  * Can add metadata tags to struct fields to convert field names to lowercase and `omitempty`
* Working with Requests
  * Key properties on `http.Request` object:
    * `Method`: `string`
    * `Header`: `map[string][]string`
    * `Body`: `io.ReadCloser` - returns `EOF` when not present
* URL Path Parameters - `request.URL` struct has a `Path` property
* Middleware - authentication, logging, session management, etc.
    ```go
    func middlewareHandler(handler http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // do stuff before intended handler here
            handler.ServeHTTP(w, r)
            // do stuff after intended handler here
      })
    }

    func intendedHandlerFunc(w http.ResponseWriter, r *http.Request) {
        // business logic here
    }

    func main() {
        intendedHandler := http.HandlerFunc(intendedHandlerFunc)
        http.Handle("/foo", middlewareHandler(intendedHandler))
    }
    ```
* Enabling CORS (Cross-Origin Resource Sharing)
  * Use `w.Header().Set()` to set proper headers to enable CORS

### Persisting Data
* Connecting to a Database via `database/sql` package
  * Open connection: `sql.Open(driverName, dataSourceName string) (*sql.DB, error)`
  * Available SQL drivers for Go: https://github.com/golang/go/wiki/SQLDrivers
* Querying Data
  * Query multiple rows with `DB.Query(query string, args ...interface{}) (*Rows, error)`
    * Scan each row with `Rows.Next()` and `Rows.Scan(dest ...interface{}) error`
  * Query a single row with `DB.QueryRow(query string, args ...interface{}) *Row`
    * Scan the single row with `Row.Scan(dest ...interface{}) error`
* Executing SQL - `DB.Exec(query string, args ...interface{}) (Result, error)`
  * Structure of `sql.Result` object:
    ```go
    type Result interface {
        LastInsertId() (int64, error)
        RowsAffected() (int64, error) 
    }
    ```
* Connection Pooling
  * Configurations
    * **Connection Max Lifetime** - sets a max amount of time a connection may be used. If set to 0, a connection will never be closed.
    * **Max Idle Connections** - sets the max number of connections in the idle connection pool (default 2)
    * **Max Open Connections** - sets the max number of open connections to the database
  * Using **Context** - allows you to set a deadline, cancel a signal, or set other request-scoped values across API boundaries and between processes.
    * `QueryContext()`, `QueryRowContext()`, and `ExecContext()` methods
* Uploading and Downloading Files
  * **base64 encode** - converts the file to a string and include in JSON payload
    * Use `Encoding.DecodeString(s string) ([]byte, error)` method to consume
  * **multipart/form-data** - uses a HTTP form to submit the raw data
    * Use `Request.FormFile(key string) (multipart.File, *multipart.FileHeader, error)` method to consume

### Using WebSockets
* WebSocket Flow:
  * Client sends HTTP GET request
    * Connection: Upgrade
    * Upgrade: websocket
    * Sec-WebSocket-Key: key
  * Server responds with status code 101
    * Switching Protocols
    * Upgrade: websocket
    * Connection: Upgrade
    * Sec-WebSocket-Accept: key
* Uses for WebSockets:
  * Chat Apps
  * Multiplayer Games
  * Stock Tickers
  * Realtime Dashboards
* Creating WebSockets in Go via `golang.org/x/net/websocket` package
  * https://pkg.go.dev/golang.org/x/net/websocket?tab=doc
  * `websocket.Conn` struct:
    ```go
    type Conn struct {
        PayloadType byte
        MaxPayloadBytes int
    }
    ```
  * `websocket.Codec` struct:
    ```go
    type Codec struct {
        Marshal func(v interface{}) (data []byte, payloadType byte, err error)
        Unmarshal func(data []byte, payloadType byte, v interface{}) (err error)
    }
    ```
  * `codec.Receive()` method
    ```go
    func (cd Codec) Receive(ws *Conn, v interface{}) (err error)
    ```
  * `codec.Send()` method
    ```go
    func (cd Codec) Send(ws *Conn, v interface{}) (err error)
    ```

### Using Templating
* Template Packages:
  * `text/template`: base functionality for working with templates in Go
    * Refer to https://golang.org/pkg/text/template for more details
  * `html/template`: same interface but with added security for HTML output
    * Refer to https://golang.org/pkg/html/template for more details
* Functions:
  * `template.New()` function
    ```go
    func New(name string) *Template
    ```
  * `Template.Parse()` method
    ```go
    func (t *Template) Parse(text string) (*Template, error)
    ```
  * `Template.Execute()` method
    ```go
    func (t *Template) Execute(wr io.Writer, data interface{}) error
    ```
* Pipelines
  * Definition: *command or sequence of commands*
    * Simple value (argument)
    * Function or method call
    * Can accept arguments
  * Examples:
    * Static values or texts: `{{ "Hello" }}` or `{{ 1234 }}`
    * Struct fields: `{{ .Message }}`
    * Function calls: `{{ println "Hi" }}`
    * Method calls: `{{ .SayHello }}` or `{{ .SaySomething "Bye" }}`
    * Chained: `{{ "Hello" | .SaySomething | printf "%s %s" "World" }}`
* Pipeline Looping with `{{ range pipeline }} T1 {{ end }}`
  * The pipeline looping will output nothing if the array/slice/map structure has 0 length
  * With optional ELSE statement: `{{ range pipeline }} T1 {{ else }} T2 {{ end }}`
    * This will output T2 if the array/slice/map structure has 0 length
  * If we need to access both index and element: `{{ range $index, $element := pipeline }} T1 {{ end }}`
* Global Template Functions
  * `and`: `{{if and true true true}} {{end}}`
  * `or`: `{{if or true false true}} {{end}}`
  * `index`: `{{index . 1}}`
  * `len`: `{{len .}}`
  * `not`: `{{if not false}} {{end}}`
  * `print`, `printf`, `println`: `{{println "hey"}}`
* Global Template Operators
  * `eq`: `arg1 == arg2`
  * `ne`: `arg1 != arg2`
  * `lt`: `arg1 < arg2`
  * `le`: `arg1 <= arg2`
  * `gt`: `arg1 > arg2`
  * `ge`: `arg1 >= arg2`
* Custom Template Functions
  * Define custom functions using `Template.Funcs()` method
    ```go
    func (t *Template) Funcs(funcMap FuncMap) *Template
    
    type FuncMap map[string]interface{}
    ```
  * Custom functions can only return a single value, or a single value/error
