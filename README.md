# gofl
[![Build Status](https://drone.io/github.com/yosisa/gofl/status.png)](https://drone.io/github.com/yosisa/gofl/latest)

gofl is a library for Go. It provides some functions that are useful for building a REST API server which supports `fl` query parameter.

`fl` query parameter often uses as a field selector, for example, `id`, `id,name` or `title,tags.name`.

## Usage
```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/yosisa/gofl"
)

type Page struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Tags  []Tag
}

type Tag struct {
	Id   int
	Name string `json:"name"`
}

func main() {
	page := Page{1, "Introduction", []Tag{{1, "sample"}, {2, "guide"}}}
	out := gofl.Pick(page, "title", "Tags.name")
	data, _ := json.Marshal(&out)
	fmt.Printf("%s\n", data)
	// {"Tags":[{"name":"sample"},{"name":"guide"}],"title":"Introduction"}
}
```

## License
The MIT License (MIT)

Copyright (c) 2014 Yoshihisa Tanaka

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
