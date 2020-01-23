package Tests

import (
	"fmt"
	"github.com/dosarudaniel/CS438_Project/chord"
	"strconv"
)

func main() {
	bit_size := 8
	out, _ := chord.HashString("inputString", bit_size)  // get the sha256 hash
	fmt.Println("id hex = " + out)
	id_int, _ := strconv.ParseUint(out, 16, bit_size)
	fmt.Println(id_int)
	id_int += 1

	fmt.Println(id_int)
	fmt.Println(strconv.FormatUint(id_int, 16))
}

