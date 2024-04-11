package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

func EncryptedPassword(password, sessionKey, ivKey string) string {
	return base64.StdEncoding.EncodeToString(AesEncrypt(password, ivKey, sessionKey))
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func AesEncrypt(src, ivKey, sessionKey string) []byte {
	block, err := aes.NewCipher([]byte(sessionKey))
	if err != nil {
		fmt.Println("key error1", err)
	}
	if src == "" {
		fmt.Println("plain content empty")
	}
	ecb := cipher.NewCBCEncrypter(block, []byte(ivKey))
	content := []byte(src)
	content = pkcs5Padding(content, block.BlockSize())
	crypt := make([]byte, len(content))
	ecb.CryptBlocks(crypt, content)
	return crypt
}

func AesDecrypt(crypt []byte, sessionKey, ivKey string) []byte {
	block, err := aes.NewCipher([]byte(sessionKey))
	if err != nil {
		fmt.Println("key error1", err)
	}
	if len(crypt) == 0 {
		fmt.Println("plain content empty")
	}
	ecb := cipher.NewCBCDecrypter(block, []byte(ivKey))
	decrypted := make([]byte, len(crypt))
	ecb.CryptBlocks(decrypted, crypt)
	return pkcs5Trimming(decrypted)
}

func GenerateChecksum(body interface{}, sessionKey, ivKey string) string {
	bytearr, _ := json.Marshal(body)
	hash := sha256.Sum256(bytearr)

	fmt.Println("body:", string(bytearr))
	hexstring := hex.EncodeToString(hash[:])
	checkSum := base64.StdEncoding.EncodeToString(AesEncrypt(hexstring, ivKey, sessionKey))
	return checkSum
}

func decodeBase32(str string) ([]byte, error) {
	numBytes := (len(str)*5 + 7) / 8
	result := make([]byte, numBytes)
	resultIndex := 0
	var which, working int

	for i := 0; i < len(str); i++ {
		ch := str[i]
		var val int
		switch {
		case ch >= 'a' && ch <= 'z':
			val = int(ch - 'a')
		case ch >= 'A' && ch <= 'Z':
			val = int(ch - 'A')
		default:
			if ch < '2' || ch > '7' {
				if ch != '=' {
					return nil, errors.New("Invalid base-32 character: " + string(ch))
				}
				which = 0
				break
			}
			val = 26 + int(ch-'2')
		}

		switch which {
		case 0:
			working = (val & 31) << 3
			which = 1
		case 1:
			working |= (val & 28) >> 2
			result[resultIndex] = byte(working)
			resultIndex++
			working = (val & 3) << 6
			which = 2
		case 2:
			working |= (val & 31) << 1
			which = 3
		case 3:
			working |= (val & 16) >> 4
			result[resultIndex] = byte(working)
			resultIndex++
			working = (val & 15) << 4
			which = 4
		case 4:
			working |= (val & 30) >> 1
			result[resultIndex] = byte(working)
			resultIndex++
			working = (val & 1) << 7
			which = 5
		case 5:
			working |= (val & 31) << 2
			which = 6
		case 6:
			working |= (val & 24) >> 3
			result[resultIndex] = byte(working)
			resultIndex++
			working = (val & 7) << 5
			which = 7
		case 7:
			working |= val & 31
			result[resultIndex] = byte(working)
			resultIndex++
			which = 0
		}
	}

	if which != 0 {
		result[resultIndex] = byte(working)
		resultIndex++
	}
	if resultIndex != len(result) {
		result = result[:resultIndex]
	}
	return result, nil
}

func GenerateAccessCode(base32Secret string) string {

	timeMillis := time.Now().UnixNano() / int64(time.Millisecond)

	key, err := decodeBase32(base32Secret)
	if err != nil {
		return ""
	}

	data := make([]byte, 8)
	value := timeMillis / 1000 / int64(30)

	for i := 7; value > 0; i-- {
		data[i] = byte(value & 255)
		value >>= 8
	}

	mac := hmac.New(sha1.New, key)
	mac.Write(data)
	hash := mac.Sum(nil)

	offset := int(hash[len(hash)-1] & 15)

	truncatedHash := binary.BigEndian.Uint32(hash[offset:]) & 2147483647
	truncatedHash %= 1000000
	return strconv.Itoa(int(truncatedHash))
}
