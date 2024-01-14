package main

import (
	"fmt"
	"log"
	
	"github.com/ghostiam/binstruct"
)

func main() {
	//test tlv data 
	//okdata.
	//data := []byte{0x01, 0x02, 0x00, 0x04, 0x05, 0x06, 0x07, 0x08}
	//data := []byte{0x01, 0x02, 0x00, 0x08, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C}
	
	//ngdata
	data := []byte{0x01, 0x02, 0x00, 0x08, 0x05, 0x06, 0x07, 0x08}


	type dataStruct struct {
		Type uint16
		ValLen  uint16
		Value []byte `bin:"len:ValLen"`
	}
	
	var datastruct dataStruct

	err  := binstruct.UnmarshalBE(data, &datastruct)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%+v",datastruct)
}
