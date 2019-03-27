package main

import "./backup2"


func main(){
    var order_list [4]int
    order_list[0] = 1
    order_list[1] = 1
    order_list[2] = 2
    order_list[3] = 3
    backup.WriteToFile(order_list)
}
