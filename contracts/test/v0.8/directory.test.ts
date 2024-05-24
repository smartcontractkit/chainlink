import fs from 'fs'
import path from 'path'
import { expect } from 'chai'

// Directories that start with a number do not currently work with typechain (https://github.com/dethcrypto/TypeChain/issues/794)
describe('Directory', () => {
  it('Should not have a file or directory starting with a number in contracts/src', () => {
    const srcPath = path.join(__dirname, '..', '..', 'src')

    const noNumbersAsFirstChar = (dirPath: string): boolean => {
      const entries = fs.readdirSync(dirPath, { withFileTypes: true })

      for (const entry of entries) {
        if (/^\d/.test(entry.name)) {
          throw new Error(
            `${path.join(dirPath, entry.name)} starts with a number`,
          )
        }

        if (entry.isDirectory()) {
          const newPath = path.join(dirPath, entry.name)
          noNumbersAsFirstChar(newPath)
        }
      }

      return true
    }

    expect(noNumbersAsFirstChar(srcPath)).to.be.true
  })
})
