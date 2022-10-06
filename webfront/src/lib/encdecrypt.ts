import aesjs from "aes-js"

function byteCount(s: string) {
   return encodeURI(s).split(/%..|./).length - 1
}

const padding = (text: string, multiple: number): string => {
   const numberOfBytes = byteCount(text)
   const numOfPads = multiple - (numberOfBytes % multiple)
   let paddedText = text
   for (let index = 0; index < numOfPads; index++) {
      index === 0 ? (paddedText += " ") : (paddedText += "=")
   }
   return paddedText
}

const depadd = (paddedText: string): string => {
   const re = / (=+)?$/
   return paddedText.replace(re, "")
}

const newEncryptionKey = async () => {
   const key = await window.crypto.subtle.generateKey(
      {
         name: "AES-CBC",
         length: 128
      },
      true,
      ["encrypt", "decrypt"]
   )

   const exported = await window.crypto.subtle.exportKey("raw", key)
   return new Uint8Array(exported)
}

const encrypt = async (
   message: string
): Promise<{
   encryptedHex: string
   ivHex: string
   keyHex: string
}> => {
   // An example 128-bit key
   const key = await newEncryptionKey()

   const iv = window.crypto.getRandomValues(new Uint8Array(16))
   //    console.log({ iv: `${aesjs.utils.hex.fromBytes(iv)}`, message, key:aesjs.utils.hex.fromBytes(key) })
   const paddedMessage = padding(message, 16)
   console.log({ paddedMessage })

   let textBytes = aesjs.utils.utf8.toBytes(paddedMessage)

   const aesCbc = new aesjs.ModeOfOperation.cbc(key, iv)
   const encryptedBytes = aesCbc.encrypt(textBytes)
   const encryptedHex = aesjs.utils.hex.fromBytes(encryptedBytes)
   const ivHex = aesjs.utils.hex.fromBytes(iv)
   const keyHex = aesjs.utils.hex.fromBytes(key)

   return {
      encryptedHex,
      ivHex,
      keyHex
   }
}

const decrypt = (
   encryptedHex: string,
   keyHex: string,
   ivHex: string
): string => {
   const encryptedBytes = aesjs.utils.hex.toBytes(encryptedHex)
   const iv = aesjs.utils.hex.toBytes(ivHex)
   const key = aesjs.utils.hex.toBytes(keyHex)

   const aesCbc = new aesjs.ModeOfOperation.cbc(key, iv)
   const decryptedBytes = aesCbc.decrypt(encryptedBytes)

   // Convert our bytes back into text
   return depadd(aesjs.utils.utf8.fromBytes(decryptedBytes))
}

export default {
   passwordHasValidFormat: (password: string): boolean => {
      const re = /^[A-Za-z0-9!@#$%^&*()]+$/g
      return password.match(re) !== null
   },

   encrypt: encrypt,
   decrypt: decrypt
}
