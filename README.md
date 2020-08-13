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
