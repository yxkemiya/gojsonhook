# Golang json hook

Have you ever be puzzle by the built-in json.Marshal/json.Unmarshal that you cannot do something before/after serialized/unserialized to json, and you can only write some function and call it every time you need,if so, this library is for you.

The very common use case is that, when I use ORM library for go, I always define a column with type time.Time, and when I marshal it to json in redis, I may want it to be a timestamp or something else, so I must do something before marshal, and when I want to unmarshal from json bytes, I also want the timestamp can construct to time.Time, something else must be done after unmarshal. I try to write some functions to do this, but there is some problem when you forget to call the functions.

# Installation

Install:

> go get -u "github.com/yxkemiya/gojsonhook"

Import:

> import "github.com/yxkemiya/gojsonhook"

# Howto
Please go through `json_test.go` to get an idea how to use this package