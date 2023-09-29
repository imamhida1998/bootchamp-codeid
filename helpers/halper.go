package helpers

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
)

var reader *bufio.Reader

//fungsi mengubah data dari struktur apapun menjadi string
func ToJson(data interface{}) string {

	//dikonversikan ke dalam bytes terlebih dahulu
	bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	//dari bytes dikonversikan ke string, dan return
	return string(bytes)
}
func Input() string {
	if reader == nil {
		reader = bufio.NewReader(os.Stdin)
	}

	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	input = strings.TrimSpace(input)
	return input
}
func Empty(val string) bool {
	if val == "" {
		return true
	} else {
		return false
	}
}

//greater than zero
func Gtz(val int) bool {
	if val > 0 {
		return true
	} else {
		return false
	}
}

//greater than or equals to
func Float32Gtez(val float32) bool {
	if val >= 0 {
		return true
	} else {
		return false
	}
}
func Float32Gtz(val float32) bool {
	if val > 0 {
		return true
	} else {
		return false
	}
}

func InString(val string, arr []string) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] == val {
			return true
		}
	}
	return false
}
