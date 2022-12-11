package newssigner

import (
	"crypto/rsa"
	"io/ioutil"
	"strings"

	"i2pgit.org/idk/reseed-tools/su3"
)

type NewsSigner struct {
	SignerID   string
	SigningKey *rsa.PrivateKey
}

func (ns *NewsSigner) CreateSu3(xmldata string) error {
	su3File := su3.New()
	su3File.FileType = su3.FileTypeXML
	su3File.ContentType = su3.ContentTypeNews

	data, err := ioutil.ReadFile(xmldata)
	if nil != err {
		return err
	}
	su3File.Content = data

	su3File.SignerID = []byte(ns.SignerID)
	su3File.Sign(ns.SigningKey)

	bytes, err := su3File.MarshalBinary()
	if err != nil {
		return err
	}
	outfile := strings.Replace(xmldata, ".atom.xml", ".su3", -1)
	return ioutil.WriteFile(outfile, bytes, 0644)
}
