package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	rnd "math/rand"
	"os"
	"path"
	"reflect"
	"strings"
	"time"
)

type Page struct {
	PageNo     int
	PageSize   int
	TotalPage  int
	TotalCount int
	FirstPage  bool
	LastPage   bool
	List       interface{}
}

type Result struct {
	Code        int
	Description string
	Detail      interface{}
}

func PageUtil(count int, pageNo int, pageSize int, list interface{}) Page {
	tp := count / pageSize
	if count%pageSize > 0 {
		tp = count/pageSize + 1
	}
	return Page{PageNo: pageNo, PageSize: pageSize, TotalPage: tp, TotalCount: count, FirstPage: pageNo == 1, LastPage: pageNo == tp, List: list}
}

func NoHtml(str string) string {
	return strings.Replace(strings.Replace(str, "<script", "&lt;script", -1), "script>", "script&gt;", -1)
}

func CreateDir(filepath string) {
	dir := path.Dir(filepath)
	exist, _ := PathExists(dir)
	if !exist {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

func PathExists(filepath string) (bool, error) {
	_, err := os.Stat(filepath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func Guid() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return Md5(base64.URLEncoding.EncodeToString(b))
}

func Md5(val string) string {
	h := md5.New()
	h.Write([]byte(val))
	return hex.EncodeToString(h.Sum(nil))
}

func DesEncrypt(val string) string {
	//key := []byte("cdb7f4f0")
	key := []byte{99, 100, 98, 55, 102, 52, 102, 48}
	result, err := encryptDes([]byte(val), key)
	if err != nil {
		panic(err) //抛异常
	}
	return base64.StdEncoding.EncodeToString(result)
}

func DencryptDES(val string) string {
	//key := []byte("cdb7f4f0")
	key := []byte{99, 100, 98, 55, 102, 52, 102, 48}
	byteval, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		panic(err) //抛异常
	}
	origData, err := decryptDes(byteval, key)
	if err != nil {
		panic(err) //抛异常
	}
	return string(origData)
}

func encryptDes(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())
	//iv := make([]byte, len("3384c73b"))
	iv := []byte{51, 51, 56, 52, 99, 55, 51, 98}
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func decryptDes(crypted, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//iv := make([]byte, len("3384c73b"))
	iv := []byte{51, 51, 56, 52, 99, 55, 51, 98}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

func EncryptAES(val string) string {
	// AES-128。key长度：16, 24, 32 bytes 对应 AES-128, AES-192, AES-256
	key := []byte("2wwydqaoqkxkcle5gsbxeuhorew9kkza")
	result, err := aesEncrypt([]byte(val), key)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(result)
}

func DencryptAES(val string) string {
	// AES-128。key长度：16, 24, 32 bytes 对应 AES-128, AES-192, AES-256
	key := []byte("2wwydqaoqkxkcle5gsbxeuhorew9kkza")
	origData, err := aesDecrypt([]byte(val), key)
	if err != nil {
		panic(err)
	}
	return string(origData)
}

func aesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func aesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func GenerateShortUrl(length int) string {
	illegal := "suck|ylx1|admin|aids|dick|duowan|penis|sex|shit|wg|bignews|falun|fapiao|freenet|fuck|hongzhi|hrichina|huanet|incest|minghui|paper64|playboy|safeweb|tibetalk|unixbox|ustibet|wstaiji|xinsheng|yuming|appledog|gn|chinamz|creaders|dafa|dajiyuan|bignews|bitch|ustibet|wstaiji|urlmap"
	arr := strings.Split(illegal, "|")

	shorturl := ""
	for {
		shorturl = BuildShortUrl(length)
		exist := false
		for i := 0; i < len(arr); i++ {
			if strings.Contains(shorturl, arr[i]) {
				exist = true
				break
			}
		}
		if !exist {
			break
		}
	}
	return shorturl
}

//生成随机字符串
func BuildShortUrl(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	arr := []byte(str)
	result := []byte{}
	r := rnd.New(rnd.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, arr[r.Intn(len(arr))])
	}
	return string(result)
}

// 判断obj是否在target中，target支持的类型arrary,slice,map
func Contain(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

func StringContains(src string, desc ...string) bool {
	src = strings.ToLower(src)
	for _, val := range desc {
		if strings.Contains(src, strings.ToLower(val)) {
			return true
		}
	}
	return false
}

func StringReplace(src, newstr string, oldstr ...string) string {
	for _, val := range oldstr {
		src = strings.Replace(src, val, newstr, -1)
	}
	return src
}

func GetPlatform(useragent string) string {
	if StringContains(useragent, "iphone", "ios", "ipod", "ipad") {
		return "ios"
	} else if StringContains(useragent, "android") {
		return "android"
	}
	return "pc"
}
