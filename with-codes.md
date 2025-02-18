+++
title = "Methods in Go"
description =  "Understand how methods work in Go real quick!"
date = 2025-02-18
go_version = "1.23.6"

[author]
name = "Shreya P Rao"
email = "shreya@example.com"
+++

```go
type Vertex struct {
    X int
    Y int
}

// Receiver pointer
func (v *Vertex) Scaler(f float64) {
    v.X = v.X + f
    v.Y = v.Y + f
}

func main() {
    v := Vertex{3, 4}
    v.Scaler(10)
    fmt.Println("After changes: ", v.X)
}
```
