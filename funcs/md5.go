package funcs

import (
	"crypto/md5"
	"fmt"
)

func MD5test() {
	data1 := []byte("entity_id_financing_id_0.99")

	data2 := []byte("entity_id_financing_id_0.91")
	data3 := []byte("entity_id_financing_id_0.92")

	data4 := []byte("entity_id_financing_id_0.99")
	fmt.Println(string(data1), fmt.Sprintf("%x", md5.Sum(data1)))
	fmt.Println(string(data2), fmt.Sprintf("%x", md5.Sum(data2)))
	fmt.Println(string(data3), fmt.Sprintf("%x", md5.Sum(data3)))
	fmt.Println(string(data4), fmt.Sprintf("%x", md5.Sum(data4)))

}
