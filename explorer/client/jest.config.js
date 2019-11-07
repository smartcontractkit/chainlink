module.exports = {
  preset: 'ts-jest',
  transform: {
    '^.+\\.(t|j)sx?$': 'ts-jest',
  },
  setupFilesAfterEnv: ['<rootDir>/jest.setup.js'],
  transformIgnorePatterns: ['node_modules/(?!(@chainlink/json-api-client)/)'],
  testRegex: '(/__tests__/(?!support/*)|(\\.|/)(test|spec))\\.tsx?$',
}
