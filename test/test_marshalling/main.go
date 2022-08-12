package main

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	i := 100
	iPtr := &i
	iPtrPtr := &iPtr
	marshalled, err := json.Marshal(&iPtrPtr)
	fmt.Println("pointer address = ", &iPtrPtr)
	if err != nil {
		panic(err)
	}
	file, err := os.OpenFile(
		"test",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	defer file.Close()
	_, err = file.Write(marshalled)
	if err != nil {
		panic(err)
	}

	fmt.Println("ToGoB64(): ", ToGOB64())
	fmt.Println("FromGoB64(ToGoB64()): ", FromGOB64(ToGOB64()))
}

func init() {
	gob.Register(SX{})
	gob.Register(Session{})
}

// go binary encoder
func ToGOB64() string {
	i := 100
	m := &i
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(m)
	if err != nil {
		fmt.Println(`failed gob Encode`, err)
	}
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

// go binary decoder
func FromGOB64(str string) string {
	m := ""
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		fmt.Println(`failed base64 Decode`, err)
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(&m)
	if err != nil {
		fmt.Println(`failed gob Decode`, err)
	}
	return m
}
