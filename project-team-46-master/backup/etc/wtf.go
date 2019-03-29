package main

import (
    "bufio"
    "fmt"
    "os"
    "time"
)

func main() {
    var orders [4]int
    orders[0] =1
    orders[1] =0
    orders[2] =2
    orders[3] =4

    f, err := os.OpenFile("./text.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        panic(err)
    }
    defer f.Close()
    for i := 0; i<len(orders); i++ {
        w := bufio.NewWriter(f)
        _, err := fmt.Fprintln(w, orders[i])
        if err != nil {
            panic(err)
        }
        w.Flush() // Flush writes any buffered data to the underlying io.Writer.
        f.Sync()  // commit the current contents of the file to stable storage.
        fmt.Println("write", orders[i])
        time.Sleep(500 * time.Millisecond)
    }
}
