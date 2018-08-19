package main

import (
	"github.com/xuebing1110/rtbus/server/app"
	_ "github.com/xuebing1110/rtbus/server/app/v3"
)

func main() {
	err := app.Get().Run()
	if err != nil {
		panic(err)
	}
}
