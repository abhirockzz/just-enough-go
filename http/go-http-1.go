Welcome to another edition of **Just Enough Go** - a series of articles about the [Go programming language](https://golang.org/) which covers some of the most commonly used [Go standard library packages](https://golang.org/pkg) e.g. `encoding/json`, `io`, `net/http`, `sync` etc. I plan to keep these relatively short and example driven.

More than happy to get your feedback via [Twitter](https://twitter.com/abhi_tweeter) or just drop a comment 🙏🏻

In this post, we will explore [`net/http`](https://golang.org/pkg/net/http/) package which provides the server and client side APIs for HTTP services. This part will provide an overview of the important server components (client APIs will be covered in another post)

> Code examples are [available on GitHub](https://github.com/abhirockzz/just-enough-go)

Let's start off with the fundamental building blocks - `ServeMux` and `Server`

![](https://media.giphy.com/media/ZvLUtG6BZkBi0/giphy.gif)

### ServeMux (multiplexer)

Simply put, `ServeMux` is an HTTP request multiplexer which is responsible for matching the URL in the request to an appropriate handler and executing it. You can create one by calling `http.NewServeMux`. The next thing you do is attach URLs and their respective handler implementations to a `ServeMux` instance using `Handle` and `HandleFunc` methods. 

Let's see how you would use the [`Handle` method](https://golang.org/pkg/net/http/#ServeMux.Handle) - it accepts a `String` and an `http.Handler`

```
func (mux *ServeMux) Handle(pattern string, handler Handler)
```

[`http.Handler` is an interface](https://golang.org/pkg/net/http/#Handler) (second parameter in the `Handle` method) with the `ServeHTTP` method

```
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

We can simply use a `struct` to provide the implementation. e.g.

```
type home struct{}

func (h home) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Welcome to the Just Enough Go! blog series!"))
}
```

... and attach this to the multiplexer as follows:

```
mux := http.NewServeMux()
mux.Handle("/", home{})
```

Let's add another handler to our `ServeMux` (`mux`) - this time, we'll use the [`HandleFunc`](https://golang.org/pkg/net/http/#ServeMux.HandleFunc) variant whose signature is the following:

```
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request))
```

Unlike the `Handle` method, `HandleFunc` accepts the handler implementation in the form of a function (along with the path for which it is to be invoked). You can use it as such

```
mux.HandleFunc("/posts", func(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("Visit http://bit.ly/just-enough-go to get started"))
})
```

### Server (HTTP server)

Now we have a multiplexer which can respond if a user navigates to the root of our service i.e. `/` as well as `/posts`. Let's tie it all together with a `Server`. It's very easy to create a new instance of a `Server`

```
server := http.Server{Addr: ":8080", Handler: mux}
```

There are a bunch of parameters which we can define for our HTTP server, but let's look at a couple of important ones i.e. `Addr` and `Handler` (highlighted above) - `Addr` is the address on which the server listens e.g. `http://localhost:8080` and `Handler` is actually an `http.Handler` instance.

The `Handler` bit is interesting because we just saw how the `Handle` method in `ServeMux` also accepts an `http.Handler`. So, do we pass the same instance here as we did for the `Handle` method in `ServeMux`, and what's the point of doing that (again)? 

If you just had route or path which you wanted to handle, you can pass an instance of an `http.Handler` (e.g. `home{}` in this case) and skip the `ServeMux` altogether. Otherwise, for most cases, you can/should pass an instance of a `ServeMux` so that you can handle multiple routes/paths (e.g. `/home`, `/items` etc.) - and this is possible because it implements `http.Handler`. Internally, it works by dispatching or routing to the appropriate handler based on the path (URL) in `http.Request`.

It defines a `ServeHTTP` method as required by the `http.Handler` interface

```
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request)
```

> The `Handler` can be `nil` - this scenario is discussed later in this post

Great! So far, have a `ServeMux` with two handlers and we have associated the Server with the multiplexer and defined where it will listen at. Finally, you just need to `start` it using the [`ListenAndServe`](https://golang.org/pkg/net/http/#Server.ListenAndServe) method

```
server.ListenAndServe()
```

That's it. Here is the consolidated code (pretty small!)

```
package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", home{})

	mux.HandleFunc("/posts", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("Visit http://bit.ly/just-enough-go to get started"))

	})

	server := http.Server{Addr: ":8080", Handler: mux}
	log.Fatal(server.ListenAndServe())
}

type home struct{}

func (h home) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Welcome to the \"Just Enough Go\" blog series!!"))
}
```

To try this

- simply save the code in a file (e.g. `go-http-1.go`) and 
- run it - `go run go-http-1.go`
- access the endpoints - `curl http://localhost:8080/` and `curl http://localhost:8080/posts`

### Default multiplexer

To make things simpler,  there is a ready-to-use multiplexer `DefaultServeMux`. You don't need to use an explicit `ServeMux`. The `Handle` and `HandleFunc` methods available in a `ServeMux` are also exposed as global functions in `net/http` package for this purpose - you can use them the same way!

```
http.Handle("/users",myHandler{})

http.HandleFunc("/items",func(rw http.ResponseWriter, req *http.Request){
    //handler logic
})
```

To start the HTTP server, you can use `http.ListenAndServe` function, just like you would with a Server instance.

```
func ListenAndServe(addr string, handler Handler) error
```

The `handler` parameter can be `nil` if you have used `http.Handle` and/or `http.HandleFunc` to specify the handler implementations for the respective routes.

### Functions as handlers

So far, we saw how to use a struct in order to implement `http.Handler` interface and use it in `HandleFunc`. You might want to use a standalone function without declaring a struct. `net/http` package defines a function type for this [`http.HandlerFunc`](https://golang.org/pkg/net/http/#HandlerFunc)

```
type HandlerFunc func(ResponseWriter, *Request)
```

`HandlerFunc` allows you to use ordinary functions as HTTP handlers e.g.

```
func welcome(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("Welcome to Just Enough Go"))
}
```

`welcome` is a standalone function with the required signature. You can use this in `Handle` method which accepts a `http.Handler` as follows:

```
http.ListenAndServe(":8080", http.HandlerFunc(welcome))
```

> `HandlerFunc(f)` is a Handler that calls the function `f`

The code looks like

```
package main

import "net/http"

func main() {
	http.Handle("/welcome", http.HandlerFunc(welcome))
	http.ListenAndServe(":8080", nil)
}

func welcome(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("Welcome to Just Enough Go"))
}
```

To try this:

- simply save the code in a file (e.g. `go-http-2.go`) and 
- run it - `go run go-http-2.go`
- access the endpoint - `curl http://localhost:8080/welcome`


Thats all for this blog where we covered the basic constructs of the server side HTTP API offered by the `net/http` package. 

I really hope you enjoyed and learned something from this 🙌 Please like and follow if you did!