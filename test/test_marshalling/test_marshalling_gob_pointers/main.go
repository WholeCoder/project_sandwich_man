package main

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
)

type SX map[string]*int

// go binary encoder
func ToGOB64(m SX) string {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(m)
	if err != nil {
		fmt.Println(`failed gob Encode`, err)
	}
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

// go binary decoder
func FromGOB64(str string) SX {
	m := SX{}
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

func main() {
	sx := SX{}
	i := 1
	sx["a"] = &i
	j := 2
	sx["b"] = &j

	fmt.Printf("sx = %#V\n", sx)

	bs1 := ToGOB64(sx)
	fmt.Printf("bs1 = %#V\n", bs1)

	sx2 := FromGOB64(bs1)
	fmt.Println("sx2[\"a\"] =", sx2["a"])
}

func init() {
	gob.Register(SX{})
}
