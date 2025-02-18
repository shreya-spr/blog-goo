## Methods in Go

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
