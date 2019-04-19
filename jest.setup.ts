import { clearDb } from './src/__tests__/testdatabase'

process.env.NODE_ENV = 'test'

afterEach(async () => clearDb())
