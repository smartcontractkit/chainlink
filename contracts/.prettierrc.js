module.exports = {
  semi: false,
  singleQuote: true,
  printWidth: 80,
  endOfLine: 'auto',
  tabWidth: 2,
  trailingComma: 'all',
  plugins: ['prettier-plugin-solidity'],
  overrides: [
    {
      files: '*.sol',
      options: {
        parser: 'solidity-parse',
        printWidth: 120,
        tabWidth: 2,
        useTabs: false,
        singleQuote: false,
        bracketSpacing: false,
        explicitTypes: 'always',
      },
    },
  ],
}
