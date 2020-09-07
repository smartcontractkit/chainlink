module.exports = {
  testEnvironment: 'node',
  testRunner: 'jest-circus/runner',
  testPathIgnorePatterns: ['/node_modules/', 'dist/'],
  globals: {
    'ts-jest': {
      tsConfig: 'tsconfig.test.json',
    },
  },
  extraGlobals: ['Math'],
  transform: {
    '^.+\\.(t|j)sx?$': ['@swc-node/jest'],
  },
}
