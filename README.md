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
