# SnowFlake
雪花算法


```
go get github.com/s84662355/SnowFlake
```


 
```

const (
	WorkerIdBitsMax           uint8 = 14            //服务器id bit 最大长度
	WorkerIdBitsMin           uint8 = 8             //服务器id bit 最小长度
	NumberBitsMax             uint8 = 14            //自增id bit 最大长度
	NumberBitsMin             uint8 = 8             //自增id  bit 最小长度
	WorkerIdBitsAndNumberBits uint8 = 22            //服务器id bit +  自增id  bit
	EpochMin                  int64 = 1666195200000 //最小起始时间
	TimeStampBits             uint8 = 41            //时间戳 bit 长度
	CharacterBits             uint8 = 1
)

```

 

```
type Options struct {
	WorkerIdBits uint8 ///服务器id长度
	NumberBits   uint8 ///自增id长度
	Epoch        int64 //起始时间戳 毫秒
	WorkId       int64 //机器id
}


```
 
```
type ID struct {
	TimeStamp int64 
	WorkId    int64
	Number    int64
	Id        int64
}
```


**Demo**

```
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

```


 