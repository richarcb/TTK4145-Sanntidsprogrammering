package esm

import (
	"fmt"
	"os"
	"bufio"
  . "../config"
)

var path = "./backup.txt"

//Writes the orders array to the backup file backup.txt
func update_backup(orders [N_floors]int, dest Order) {
	if dest.Button == BT_Cab{
		orders[dest.Floor] = 1
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for i := 0; i < len(orders); i++ {
		w := bufio.NewWriter(f)
		_, err := fmt.Fprintln(w, orders[i])
		if err != nil {
			panic(err)
		}
		w.Flush()
		f.Sync()
	}
}

//Checks if the backup textfile exists
func backup_exist() bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//Reads from the current backup file and returns the orders in a list of integers
func read_from_backup() [N_floors]int {
	var orders [N_floors]int
	f, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	i := 0
	for  k:= 0; k<N_floors; k++{
		n, err := fmt.Fscanln(f, &i)
		if n == 1 {
			orders[k] = i
		}
		if err != nil {
			fmt.Println(err)
			return orders
		}
	}
	return orders
}
