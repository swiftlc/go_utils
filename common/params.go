package common

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
)

func ReadStdin() <-chan string {
	reader := bufio.NewReader(os.Stdin)
	ch := make(chan string)
	go func() {
		for {
			line, _, err := reader.ReadLine()
			if err != nil {
				close(ch)
				if err != io.EOF {
					log.Printf("read stdin err:%+v", err)
				}
				return
			}
			ch <- string(line)
		}
	}()
	return ch
}

func GetParamOrStdin() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}

	data, _ := ioutil.ReadAll(os.Stdin)
	return string(data)
}

func JsonToForm(data []byte) (url.Values, error) {
	var d map[string]interface{}
	decoder := json.NewDecoder(bytes.NewBuffer(data))
	decoder.UseNumber()
	err := decoder.Decode(&d)
	if err != nil {
		return nil, err
	}
	q := make(url.Values)
	for k, v := range d {
		q.Set(k, Str(v))
	}
	return q, nil
}
