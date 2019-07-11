package etcder

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"time"
)

var (
	rander = rand.New(rand.NewSource(time.Now().UnixNano()))
	getRandomCjk = func() rune {return 0x4e00 + rander.Int31n(0x9fff-0x4e00)}
	getRandomAscii = func() rune {return 0x21 + rander.Int31n(0x7e-0x21)}
	getRandomGreek = func() rune {return 0x391 + rander.Int31n(0x3c9-0x391)}
)

func GetRandomKv(num, kmin, kmax, vmin, vmax int) [][]string {
	strs := make([][]string, num)
	for i:=0; i<num; i++ {
		klen := kmin + rander.Intn(kmax-kmin)
		vlen := vmin + rander.Intn(vmax-vmin)
		krunes := make([]rune, klen)
		vrunes := make([]rune, vlen)
		for j:=0; j<klen; j++ {
			krunes[j] = getRandomAscii()
		}
		for j:=0; j<vlen; j++ {
			vrunes[j] = getRandomCjk()
		}
		strs[i] = []string{string(krunes), string(vrunes)}
	}
	return strs
}

func ToClose(closer io.Closer) {
	if closer != nil {
		err := closer.Close()
		if err != nil {
			fmt.Println("close err:", err)
		}
	}
}

func GetCertFromFile(filename string) (*x509.Certificate, error) {
	var err error
	byter, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("read file err:", err)
		return nil, err
	}
	block, _ := pem.Decode(byter)
	return x509.ParseCertificate(block.Bytes)
}