<p align="center"><img src="art/diabuddy.webp" alt="Diabuddy Error package"></p>

## Introduction

Diabuddy Errors package  provides a library which allows us to use same error format for Api in All Diabuddy Api,

### Install
With Go's module support, go [build|run|test] automatically fetches the necessary dependencies when you add the import in your code: 
```go
    import "github.com/hbttundar/diabuddy-errors"
```
Alternatively, use go get:
```bash
   go get -u github.com/hbttundar/diabuddy-errors
```
### Use 
A basic example:

```go
package main

import (
	"encoding/json"
	errors "github.com/hbttundar/diabuddy-errors"
)

func main() {
	data := []byte("{response from an specific api}")
	specificStruct := &struct {
		User   *User  `json:"user"`
		Gender string `json:"gender"`
	}{}
	if  err := json.Unmarshal(data, &specificStruct); err != nill {
		errors.NewApiError(errors.UnprocessableEntityErrorType, "unprocess response data", diabuddyErrors.WithInternalError(err))
	}
}
```
