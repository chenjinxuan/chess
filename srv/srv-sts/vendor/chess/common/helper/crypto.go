package helper

import (
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
)

func pkc5_pad(text string, blockSize int) string {
	padding := blockSize - len(text)%blockSize

	return text + strings.Repeat(string(padding), padding)
}

func pkcs5_unpad(in []byte) []byte {
	pad := in[len(in)-1]

	if int(pad) > len(in) {
		return in
	}
	start := len(in) - int(pad)

	return in[0:start]
}

type Des struct {
	key string
	iv  string
}

/*
* Decrypt strings
* CBC
 */
func (d *Des) Decrypt(text string) (decrypted string, err error) {
	block, c_err := des.NewCipher([]byte(d.key))
	if c_err != nil {
		err = c_err
		return
	}

	blockMode := cipher.NewCBCDecrypter(block, []byte(d.iv))

	dec_bytes, hex_error := hex.DecodeString(text)

	if hex_error != nil {
		err = hex_error
		return
	}

	dec_len := len(dec_bytes)

	//log.Println(dec_len)
	out := make([]byte, dec_len)

	blockMode.CryptBlocks(out, []byte(dec_bytes))
	decrypted = string(pkcs5_unpad(out))
	return
}

/*
* Encrypt strings
* CBC
 */
func (d *Des) Encrypt(text string) (encrypted string, err error) {
	block, c_err := des.NewCipher([]byte(d.key))
	if c_err != nil {
		err = c_err
		return
	}
	blockMode := cipher.NewCBCEncrypter(block, []byte(d.iv))

	in := []byte(pkc5_pad(text, block.BlockSize()))
	in_len := len(in)

	out := make([]byte, in_len)

	blockMode.CryptBlocks(out, []byte(in))

	encrypted = strings.ToUpper(hex.EncodeToString(out))
	return
}

/*
* Decrypt strings
* ECB
 */
func (d *Des) DecryptECB(text string) (decrypted string, err error) {

	cipher, c_err := des.NewCipher([]byte(d.key))
	if c_err != nil {
		err = c_err
		return
	}

	dec_bytes, hex_error := hex.DecodeString(text)

	if hex_error != nil {
		err = hex_error
		return
	}

	dec_len := len(dec_bytes)

	//log.Println(dec_len)
	out := make([]byte, dec_len)

	for i := 0; i < dec_len/8; i++ {
		start := i * 8
		end := start + 8

		cipher.Decrypt(out[start:end], dec_bytes[start:end])
	}
	decrypted = string(pkcs5_unpad(out))
	return
}

/*
* Encrypt strings
* ECB
 */
func (d *Des) EncryptECB(text string) (encrypted string, err error) {

	cipher, c_err := des.NewCipher([]byte(d.key))
	if c_err != nil {
		err = c_err
		return
	}

	in := []byte(pkc5_pad(text, 8))
	in_len := len(in)

	out := make([]byte, in_len)

	for i := 0; i < in_len/8; i++ {
		start := i * 8
		end := start + 8

		cipher.Encrypt(out[start:end], in[start:end])
	}

	encrypted = strings.ToUpper(hex.EncodeToString(out))
	return
}

func NewDes(key string, iv string) (des *Des, err error) {

	if len(key) != 8 {
		err = errors.New("The length of key should be 8.")
		return
	}

	if len(iv) == 0 {
		iv = key
	} else if len(iv) != 8 {
		err = errors.New("The length of iv should be 8.")
		return
	}

	des = &Des{key: key, iv: iv}
	return
}

func DesEncrypt(key string, iv string, text string) (encrypted string) {
	defer func() {
		if r := recover(); r != nil {
			encrypted = ""
		}
	}()

	//Decrypt
	d, n_err := NewDes(key, iv)
	if n_err != nil {
		return
	}
	enc, e_err := d.Encrypt(text)
	if e_err != nil {
		return
	}
	encrypted = enc
	return
}

func DesDecrypt(key string, iv string, encrypted string) (text string) {
	defer func() {
		if r := recover(); r != nil {
			text = ""
		}
	}()

	d, n_err := NewDes(key, iv)
	if n_err != nil {
		fmt.Println(n_err)
		return
	}
	decrypted, d_err := d.Decrypt(encrypted)
	if d_err != nil {
		fmt.Println(d_err)
		return
	}
	text = decrypted
	return
}

func NewDesECB(key string) (des *Des, err error) {

	if len(key) != 8 {
		err = errors.New("The length of key should be 8.")
		return
	}

	des = &Des{key: key}
	return
}

func DesEncryptECB(key string, text string) (encrypted string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("DesEncryptECB panic:", r)
		}
	}()

	//Decrypt
	d, n_err := NewDesECB(key)
	if n_err != nil {
		return
	}
	enc, e_err := d.EncryptECB(text)
	if e_err != nil {
		return
	}
	encrypted = enc
	return
}

func DesDecryptECB(key string, encrypted string) (text string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("DesDecryptECB panic:", r)
		}
	}()

	d, n_err := NewDesECB(key)
	if n_err != nil {
		return
	}
	decrypted, d_err := d.DecryptECB(encrypted)
	if d_err != nil {
		return
	}
	text = decrypted
	return
}

func Md5(data string) string {
	h := md5.New()
	io.WriteString(h, data)
	return fmt.Sprintf("%x", h.Sum(nil))
}
