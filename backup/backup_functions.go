package backup

//package main //for testing main

import (
	"fmt"
	"os"
	"bufio"
  . "../config"
)

var path = "./backup.txt"

func Update_backup(orders [N_floors]int, dest Order) { //CreateBackup
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
		w.Flush() // Flush writes any buffered data to the underlying io.Writer.
		f.Sync()  // commit the current contents of the file to stable storage.
		//fmt.Println("write", orders[i])
		//time.Sleep(500 * time.Millisecond)
	}
}

func BackupExists() bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func ReadFromBackup() [4]int { //fixed version
	var orders [4]int
	f, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	i := 0
	for  k:= 0; k<4; k++{
		n, err := fmt.Fscanln(f, &i)
		if n == 1 {
			orders[k] = i
			//fmt.Println(i)
		}
		if err != nil {
			fmt.Println(err)
			return orders
		}
	}
	return orders
}
