package backup

//package main //for testing main

import (
	"fmt"
	"io"
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
	deleteFile()
	createFile()
	writeToFile(orders)
}

func ReadFromBackup() [config.N_floors]int { //fixed version
	var orders [config.N_floors]int
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
			fmt.Printf("%#v", orders[k])
			k++
		}
		if err != nil {
			fmt.Println(err)
			return orders
		}
	}
	return orders
}

//_______________________Help functions_____________________________//

func writeToFile(orders [config.N_floors]int) {

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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

func createFile() {
	// detect if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if isError(err) {
			return
		}
		defer file.Close()
	}

	fmt.Println("==> done creating file", path)
}

func deleteFile() {
	// delete file
	var err = os.Remove(path)
	if isError(err) {
		return
	}

	fmt.Println("==> done deleting file")
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}

//UNUSED
func writeFile() {
	// open file using READ & WRITE permission
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if isError(err) {
		return
	}
	defer file.Close()

	// write some text line-by-line to file
	_, err = file.WriteString("halo\n")
	if isError(err) {
		return
	}
	_, err = file.WriteString("mari belajar golang\n")
	if isError(err) {
		return
	}

	// save changes
	err = file.Sync()
	if isError(err) {
		return
	}

	fmt.Println("==> done writing to file")
}

//UNUSED
func readFile() {
	// re-open file
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if isError(err) {
		return
	}
	defer file.Close()

	// read file, line by line
	var text = make([]byte, 1024)
	for {
		_, err = file.Read(text)

		// break if finally arrived at end of file
		if err == io.EOF {
			break
		}

		// break if error occured
		if err != nil && err != io.EOF {
			isError(err)
			break
		}
	}

	fmt.Println("==> done reading from file")
	fmt.Println(string(text))
}
