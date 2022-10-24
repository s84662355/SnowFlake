package main

import "github.com/s84662355/SnowFlake"
import "fmt"

func main() {

	options := SnowFlake.Options{
		WorkerIdBits: 9,
		NumberBits:   13,
		Epoch:        1666576805000,
		WorkId:       6,
	}

	app, err := SnowFlake.New(&options)
	if err != nil {
		fmt.Println(err)
		return
	}

	i := 0
	for i < 1000 {
		id := app.NextId()
		fmt.Println(id)
		fmt.Println(app.DecodeID(id.Id))
		fmt.Println("-----------------")
		i++
	}

}
