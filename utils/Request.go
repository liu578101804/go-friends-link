package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetBody(url string) ([]byte, error) {

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error GET: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read body error:", err)
		return nil, err
	}

	return body, nil
}
