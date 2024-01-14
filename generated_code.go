package 5066/message

import  github.com/ghostiam/binstruct 

type C_PDU1 struct {
	Type uint16 `yaml:"Type" `
	Vallen uint16 `yaml:"Vallen" `
	Value []byte `yaml:"Value" bin:"len:ValLen"`
}

func (*C_PDU1 p) Encode()([]byte,error) {

}

func Decode_C_PDU1(bindata []byte)(C_PDU1,error) {
	var st C_PDU1
	err := binstruct.UnmarshalBE(bindata,&st)
	return st,err
}

