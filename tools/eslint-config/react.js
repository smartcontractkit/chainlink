module.exports = {
  root: true,
  extends: ['@chainlink/eslint-config', 'plugin:react/recommended'],
  plugins: ['react', 'react-hooks'],
  env: {
    node: true,
    browser: true,
  },
  parserOptions: {
    ecmaFeatures: {
      jsx: true,
    },
  },
  settings: {
    react: {
      version: 'detect',
    },
  },
  rules: {
    'react/prop-types': 'off',
    // TODO: enable these after removing use-react-hooks package
    'react-hooks/rules-of-hooks': 'off',
    'react-hooks/exhaustive-deps': 'off',
  },
}
