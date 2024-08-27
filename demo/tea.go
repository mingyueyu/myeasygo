package demo

import (
	"fmt"

	"github.com/mingyueyu/myeasygo/util/system"
)

func TeaDemo() {
	// Example usage
	plaintext := "this is a very long string that we want to encrypt with the TEA algorithm.还有中文和数字123哦！" // 明文
	k := [4]uint32{0x01234567, 0x89abcdef, 0xfedcba98, 0x76543210} // 密钥

	// Encrypt
	encryptedHex := system.TeaEncryptStringWithKey(plaintext, k)
	fmt.Printf("Encrypted value (hex): %s\n", encryptedHex)

	// Decrypt
	decrypted := system.TeaDecryptStringWithKey(encryptedHex, k)
	fmt.Printf("Decrypted string: %s\n", decrypted)
}
