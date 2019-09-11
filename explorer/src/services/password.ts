import bcrypt from 'bcrypt'

const SALT_ROUNDS: number = 10

function hash(plaintext: string): Promise<string> {
  return bcrypt.hash(plaintext, SALT_ROUNDS)
}

function compare(plaintext: string, hash: string): Promise<boolean> {
  return bcrypt.compare(plaintext, hash)
}

export { hash, compare }
