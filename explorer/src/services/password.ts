import argon2 from 'argon2'

export function hash(plaintext: string): Promise<string> {
  return argon2.hash(plaintext, { type: argon2.argon2id })
}

export function compare(plaintext: string, hash: string): Promise<boolean> {
  return argon2.verify(hash, plaintext)
}
