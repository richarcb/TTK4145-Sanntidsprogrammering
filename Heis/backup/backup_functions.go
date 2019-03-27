package backup

//package main //for testing main

import (
	"fmt"
	"os"

	//"time"
	"bufio"

	config "../Config"
)

var path = "./text.txt"

/*
func main(){ //testing functions

	var array [config.N_floors]int
	array[0] =1
	array[2] =1
	UpdateBackup(array) //testing for button_event
	//array = ReadFromBackup() //testing for Init_mem
	fmt.Println(array)
}
//*/

func UpdateBackup(orders [config.N_floors]int) { //CreateBackup
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
	PrintBackup()
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
	k := 0
	for {
		n, err := fmt.Fscanln(f, &i)
		if n == 1 {
			orders[k] = i
			//fmt.Println(i)
			//fmt.Printf("%#v", orders[k])
			k++
		}
		if err != nil {
			fmt.Println(err)
			return orders
		}
	}
	return orders
}

func PrintBackup() {
	fmt.Printf("Backup: %v", ReadFromBackup())
	fmt.Printf("\n")
}
