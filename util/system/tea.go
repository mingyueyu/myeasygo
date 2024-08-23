package system

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

var teaK [4]uint32


func TeaDemo() {
	// Example usage
	plaintext := "this is a very long string that we want to encrypt with the TEA algorithm.还有中文和数字123哦！" // 明文
	k := [4]uint32{0x01234567, 0x89abcdef, 0xfedcba98, 0x76543210} // 密钥

	// Encrypt
	encryptedHex := TeaEncryptStringWithKey(plaintext, k)
	fmt.Printf("Encrypted value (hex): %s\n", encryptedHex)

	// Decrypt
	decrypted := TeaDecryptStringWithKey(encryptedHex, k)
	fmt.Printf("Decrypted string: %s\n", decrypted)
}

func TeaInit(k [4]uint32) {
	teaK = k
}

func TeaEncryptString(plaintext string) string {
	if teaK == [4]uint32{} {
		return ""
	}
	return TeaEncryptStringWithKey(plaintext, teaK)
}
// TeaEncryptString 实现TEA算法对长字符串的加密
func TeaEncryptStringWithKey(plaintext string, k [4]uint32) string {
	var encryptedBlocks []byte

	// 分块加密
	for i := 0; i < len(plaintext); i += 8 {
		end := i + 8
		if end > len(plaintext) {
			end = len(plaintext)
		}
		block := strToUint32Array(plaintext[i:end])
		encryptedBlock := teaEncrypt(block, k)
		encryptedBlocks = append(encryptedBlocks, uint32ArrayToBytes(encryptedBlock)...)
	}
	return fmt.Sprintf("%x", encryptedBlocks)
}

// TeaDecryptString 实现TEA算法对长字符串的解密
func TeaDecryptString(encryptedHex string) string {
	if teaK == [4]uint32{} {
		return ""
	}
	return TeaDecryptStringWithKey(encryptedHex, teaK)
}
func TeaDecryptStringWithKey(encryptedHex string, k [4]uint32) string {
	// 将加密后的十六进制字符串转换为字节数组
	ciphertext, _ := hex.DecodeString(encryptedHex)

	var decryptedText bytes.Buffer

	// 分块解密
	for i := 0; i < len(ciphertext); i += 8 {
		block := bytesToUint32Array(ciphertext[i : i+8])
		decryptedBlock := teaDecrypt(block, k)
		decryptedText.WriteString(uint32ArrayToStr(decryptedBlock))
	}

	return decryptedText.String()
}

// TeaEncrypt 实现Tea算法加密
func teaEncrypt(v [2]uint32, k [4]uint32) [2]uint32 {
	var delta uint32 = 0x9e3779b9
	var sum uint32 = 0
	var v0 uint32 = v[0]
	var v1 uint32 = v[1]
	var k0 uint32 = k[0]
	var k1 uint32 = k[1]
	var k2 uint32 = k[2]
	var k3 uint32 = k[3]

	for i := 0; i < 32; i++ {
		sum += delta
		v0 += ((v1 << 4) + k0) ^ (v1 + sum) ^ ((v1 >> 5) + k1)
		v1 += ((v0 << 4) + k2) ^ (v0 + sum) ^ ((v0 >> 5) + k3)
	}

	return [2]uint32{v0, v1}
}

// TeaDecrypt 实现Tea算法解密
func teaDecrypt(v [2]uint32, k [4]uint32) [2]uint32 {
	var delta uint32 = 0x9e3779b9
	var sum uint32 = delta << 5
	var v0 uint32 = v[0]
	var v1 uint32 = v[1]
	var k0 uint32 = k[0]
	var k1 uint32 = k[1]
	var k2 uint32 = k[2]
	var k3 uint32 = k[3]

	for i := 0; i < 32; i++ {
		v1 -= ((v0 << 4) + k2) ^ (v0 + sum) ^ ((v0 >> 5) + k3)
		v0 -= ((v1 << 4) + k0) ^ (v1 + sum) ^ ((v1 >> 5) + k1)
		sum -= delta
	}

	return [2]uint32{v0, v1}
}


// StrToUint32Array 将8字节的字符串转换为两个32位的无符号整数数组
func strToUint32Array(s string) [2]uint32 {
	var arr [2]uint32
	b := []byte(s)
	if len(b) < 8 {
		b = append(b, bytes.Repeat([]byte{0}, 8-len(b))...)
	}
	arr[0] = binary.BigEndian.Uint32(b[:4])
	arr[1] = binary.BigEndian.Uint32(b[4:])
	return arr
}

// Uint32ArrayToStr 将两个32位的无符号整数数组转换为字符串
func uint32ArrayToStr(arr [2]uint32) string {
	b := make([]byte, 8)
	binary.BigEndian.PutUint32(b[:4], arr[0])
	binary.BigEndian.PutUint32(b[4:], arr[1])
	return string(bytes.TrimRight(b, "\x00")) // 去除填充的零字节
}


// Uint32ArrayToBytes 将两个32位无符号整数数组转换为字节数组
func uint32ArrayToBytes(arr [2]uint32) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint32(b[:4], arr[0])
	binary.BigEndian.PutUint32(b[4:], arr[1])
	return b
}

// BytesToUint32Array 将字节数组转换为两个32位无符号整数数组
func bytesToUint32Array(b []byte) [2]uint32 {
	var arr [2]uint32
	arr[0] = binary.BigEndian.Uint32(b[:4])
	arr[1] = binary.BigEndian.Uint32(b[4:])
	return arr
}

