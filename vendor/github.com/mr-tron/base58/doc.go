/*
Package base58 provides fast implementation of base58 encoding.

Base58 Usage

To decode a base58 string:

	encoded := "1QCaxc8hutpdZ62iKZsn1TCG3nh7uPZojq"
	buf, _ := base58.Decode(encoded)

To encode the same data:

	encoded := base58.Encode(buf)
 
With custom alphabet

  customAlphabet := base58.NewAlphabet("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")
  encoded := base58.EncodeAlphabet(buf, customAlphabet)
  
*/
package base58
