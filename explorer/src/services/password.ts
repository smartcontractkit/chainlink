import bcrypt from 'bcrypt'

const DEFAULT_SALT_ROUNDS = 10

export function hash(
  plaintext: string,
  rounds: number = DEFAULT_SALT_ROUNDS,
): Promise<string> {
  return bcrypt.hash(plaintext, rounds)
}

export function compare(plaintext: string, hash: string): Promise<boolean> {
  return bcrypt.compare(plaintext, hash)
}
