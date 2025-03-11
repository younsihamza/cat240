package utils

import (
	"encoding/json"
	"fmt"
	"radar240/global"
	"radar240/parser/cat240"
)

func ParseData() {
	for {
		select {
		case data := <-global.FilteredData:
			validData, err := cat240.Parser(data)
			if err != nil {
				fmt.Println(err)
				continue 
			}
			DataBlock := cat240.Decode(validData)
			jsonData , err := json.Marshal(DataBlock)
			if jsonData == nil {
				fmt.Println(err)
				continue
			}
			global.ParsedData <- jsonData
		}
	}
}