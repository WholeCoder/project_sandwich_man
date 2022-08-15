package main

import (
	"fmt"
	"sort"

	"github.com/iancoleman/orderedmap"
)

func main() {
	var o *orderedmap.OrderedMap
	o = orderedmap.New()
	o.Set("z", 111)
	v, _ := o.Get("z")
	o.Set("z", v.(int)+1)
	o.Set("a", 222)

	for _, k := range o.Keys() {
		v, _ := o.Get(k)
		fmt.Println("h[", k, "] =", v)
	}

	//o.Delete("z")

	o.SortKeys(sort.Strings)

	fmt.Println("o =", o)
}
