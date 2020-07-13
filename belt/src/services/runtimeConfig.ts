import fs from 'fs'
import { join } from 'path'

// Runtime configuration for belt deploy and belt exec
export interface RuntimeConfig {
  chainId: number
  mnemonic: string
  infuraProjectId: string
  gasPrice: number
  gasLimit: number
}

const RUNTIME_CONFIG = '.beltrc'
const DEFAULTS: RuntimeConfig = {
  chainId: 4,
  mnemonic: '',
  infuraProjectId: '',
  gasPrice: 40000000000, // 40 gwei
  gasLimit: 8000000,
}

/**
 * Helper for reading from and writing to .beltrc
 */
export class RuntimeConfigParser {
  path: string

  constructor(path = '.') {
    this.path = path
  }

  private exists(): boolean {
    return fs.existsSync(this.filepath())
  }

  filepath(): string {
    return join(this.path, RUNTIME_CONFIG)
  }

  load(): RuntimeConfig {
    let result
    try {
      const buffer = fs.readFileSync(this.filepath(), 'utf8')
      result = JSON.parse(buffer.toString())
    } catch (e) {
      throw Error(`Could not load .beltrc at ${this.path}`)
    }
    return result
  }

  loadWithDefaults(): RuntimeConfig {
    if (this.exists()) {
      return this.load()
    }
    return DEFAULTS
  }

  set(config: RuntimeConfig) {
    // TODO: validate config
    // assert(config.network);
    // assert(config.mnemonic);
    // assert(config.infuraProjectId);

    fs.writeFileSync(this.filepath(), JSON.stringify(config, null, 4))
  }
}