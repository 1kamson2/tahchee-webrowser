package View

import "fmt"

func SearchBar() string {
	fmt.Println("======= SEARCH =======")
	var url_s string
	fmt.Scanln(&url_s)
	fmt.Println("======================")
	return url_s
}
