package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/charoit/fileupdater/cache"
	"github.com/charoit/fileupdater/settings"
)

func main() {

	fmt.Print("Config loading...")
	set, err := settings.Load("settings.json")
	if err != nil {
		panic(err)
	}
	fmt.Println("OK")

	fmt.Print("Conecting redis...")
	c, err := cache.GetCache(set)
	if err != nil {
		panic(err)
	}
	fmt.Println("OK")

	fmt.Print("Caching files...")
	err = c.Update()
	if err != nil {
		panic(err)
	}
	fmt.Println("OK")

	//for _,v := range set.Paths {
	//	fmt.Println("-------------------------------------------------------------------------------------")
	//	lst, _ := c.Get(v.Path)
	//	for i := 0; i < len(lst); i++ {
	//		fmt.Println(lst[i])
	//	}
	//}

	//for _,v := range set.Paths {
	//	v, _ := c.Raw(v.Path)
	//	fmt.Println(string(v))
	//	fmt.Println("-------------------------------------------------------------------------------------")
	//}

	http.HandleFunc("/:id", func(w http.ResponseWriter, r *http.Request) {
		n := r.FormValue("id")
		log.Println(n)
		d, _ := c.Raw(n)
		w.Write(d)
	})

	http.ListenAndServe(":3000", nil)
}
