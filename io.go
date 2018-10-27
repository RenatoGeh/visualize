package main

import (
	"encoding/binary"
	"github.com/RenatoGeh/gospn/spn"
	"io/ioutil"
	"os"
)

func Save(S []spn.SPN, filename string) {
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	f.Write([]byte{byte(ClassVar.Categories)})
	for _, s := range S {
		bytes := spn.Marshal(s)
		m := uint64(len(bytes))
		mb := make([]byte, 8)
		binary.LittleEndian.PutUint64(mb, m)
		f.Write(mb)
		f.Write(bytes)
	}
}

func Load(filename string) []spn.SPN {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	n := int(data[0])
	data = data[1:]
	S := make([]spn.SPN, n)
	for i := 0; i < n; i++ {
		m := binary.LittleEndian.Uint64(data[:8])
		data = data[8:]
		S[i] = spn.Unmarshal(data[:m])
		data = data[m:]
	}
	return S
}
