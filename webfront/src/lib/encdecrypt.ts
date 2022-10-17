import aesjs from "aes-js"
import { Buffer } from "buffer"
import scrypt from "scrypt-js"

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

const validationString = "hello-goprivate "

const addClearTextValidation = (text: string): string => {
   return `${validationString}${text}`
}

const validateClearText = (
   text: string
): {
   validText: boolean
   validatedText: string
} => {
   const validationRE = new RegExp(`^${validationString}`)
   const validText = validationRE.test(text)
   const validatedText = text.replace(validationRE, "")

   return { validText, validatedText }
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

const encryptionKeyFromPassword = async (
   password: string
): Promise<Uint8Array> => {
   const buffedPassword = Buffer.from(password.normalize("NFKC"))

   /* 
    Salt is static since users should be able to
    decide on a password before using the service.
    The passwords are also never stored on the server
    so a rainbow attack is not an option.
    */
   const buffedSalt = Buffer.from("thisismysalt")

   const N = 1024,
      r = 8,
      p = 1
   const dkLen = 16

   function updateInterface(progress: number) {
      const currentProgress = Math.trunc(100 * progress)
      console.log(`Deriving key: ${currentProgress}`)
   }

   return await scrypt.scrypt(
      buffedPassword,
      buffedSalt,
      N,
      r,
      p,
      dkLen,
      updateInterface
   )
}

const _encrypt = async (
   message: string,
   keyAsBytes: aesjs.ByteSource
): Promise<{
   encryptedHex: string
   ivHex: string
   keyHex: string
}> => {
   const iv = window.crypto.getRandomValues(new Uint8Array(16))
   const messageWithValidation = addClearTextValidation(message)
   console.log({ messageWithValidation })
   const paddedMessage = padding(messageWithValidation, 16)
   console.log({ paddedMessage })

   let textBytes = aesjs.utils.utf8.toBytes(paddedMessage)

   const aesCbc = new aesjs.ModeOfOperation.cbc(keyAsBytes, iv)
   const encryptedBytes = aesCbc.encrypt(textBytes)
   const encryptedHex = aesjs.utils.hex.fromBytes(encryptedBytes)
   const ivHex = aesjs.utils.hex.fromBytes(iv)
   const keyHex = aesjs.utils.hex.fromBytes(keyAsBytes)

   return {
      encryptedHex,
      ivHex,
      keyHex
   }
}

type MessageWithValidationCheck = {
   messageIsValid: boolean
   message: string
}

const _decrypt = (
   encryptedHex: string,
   keyHex: string,
   ivHex: string
): MessageWithValidationCheck => {
   const encryptedBytes = aesjs.utils.hex.toBytes(encryptedHex)
   const iv = aesjs.utils.hex.toBytes(ivHex)
   const key = aesjs.utils.hex.toBytes(keyHex)

   const aesCbc = new aesjs.ModeOfOperation.cbc(key, iv)
   const decryptedBytes = aesCbc.decrypt(encryptedBytes)

   // Convert our bytes back into text
   const { validText, validatedText } = validateClearText(
      depadd(aesjs.utils.utf8.fromBytes(decryptedBytes))
   )
   return {
      messageIsValid: validText,
      message: validatedText
   }
}

const encrypt = async (
   message: string
): Promise<{
   encryptedHex: string
   ivHex: string
   keyHex: string
}> => {
   const key = await newEncryptionKey()
   return _encrypt(message, key)
}

const decrypt = async (
   encryptedHex: string,
   keyHex: string,
   ivHex: string
): Promise<MessageWithValidationCheck> => {
   return _decrypt(encryptedHex, keyHex, ivHex)
}

const decryptWithPassword = async (
   encryptedHex: string,
   password: string,
   ivHex: string
): Promise<MessageWithValidationCheck> => {
   const key = await encryptionKeyFromPassword(password)
   const keyHex = aesjs.utils.hex.fromBytes(key)
   return _decrypt(encryptedHex, keyHex, ivHex)
}

const encryptWithPassword = async (
   message: string,
   password: string
): Promise<{
   encryptedHex: string
   ivHex: string
   keyHex: string
}> => {
   const key = await encryptionKeyFromPassword(password)

   return _encrypt(message, key)
}

export default {
   passwordHasValidFormat: (password: string): boolean => {
      const re = /^[A-Za-z0-9!@#$%^&*()]+$/g
      return password.match(re) !== null
   },

   encrypt: encrypt,
   decrypt: decrypt,
   encryptWithPassword: encryptWithPassword,
   decryptWithPassword: decryptWithPassword
}
