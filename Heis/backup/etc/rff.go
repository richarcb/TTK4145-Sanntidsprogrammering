package main

import (
    "fmt"
    "os"
    //"time"
)

func main()
func readFromFile(orders []int){

    f, err := os.OpenFile("./text.txt", os.O_RDWR, 0666)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    i := 0
    k := 0
    for {
        n, err := fmt.Fscanln(f, &i)
        if n == 1 {
            orders[k] = i
            //fmt.Println(i)
            fmt.Printf("%#v",orders[k])
            k++
        }
        if err != nil {
            fmt.Println(err)
            return
        }
        //time.Sleep(500 * time.Millisecond)
    }
    //fmt.Println(orders)
}
