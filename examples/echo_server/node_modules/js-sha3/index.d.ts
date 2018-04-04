type Message = string | number[] | ArrayBuffer | Uint8Array;

interface Hasher {
  /**
   * Update hash
   *
   * @param message The message you want to hash.
   */
  update(message: Message): Hasher;

  /**
   * Return hash in hex string.
   */
  hex(): string;

  /**
   * Return hash in hex string.
   */
  toString(): string;

  /**
   * Return hash in ArrayBuffer.
   */
  arrayBuffer(): ArrayBuffer;

  /**
   * Return hash in integer array.
   */
  digest(): number[];

  /**
   * Return hash in integer array.
   */
  array(): number[];
}

interface Hash {
  /**
   * Hash and return hex string.
   *
   * @param message The message you want to hash.
   */
  (message: Message): string;

  /**
   * Create a hash object.
   */
  create(): Hasher;

  /**
   * Create a hash object and hash message.
   *
   * @param message The message you want to hash.
   */
  update(message: Message): Hasher;
}

interface ShakeHash {
  /**
   * Hash and return hex string.
   *
   * @param message The message you want to hash.
   * @param outputBits The length of output.
   */
  (message: Message, outputBits: number): string;

  /**
   * Create a hash object.
   *
   * @param outputBits The length of output.
   */
  create(outputBits: number): Hasher;

  /**
   * Create a hash object and hash message.
   *
   * @param message The message you want to hash.
   * @param outputBits The length of output.
   */
  update(message: Message, outputBits: number): Hasher;
}

export var sha3_512: Hash;
export var sha3_384: Hash;
export var sha3_256: Hash;
export var sha3_224: Hash;
export var keccak_512: Hash;
export var keccak_384: Hash;
export var keccak_256: Hash;
export var keccak_224: Hash;
export var shake_128: ShakeHash;
export var shake_256: ShakeHash;
