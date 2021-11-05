import 'mock-local-storage'
import promiseFinally from 'promise.prototype.finally'
import JavascriptTimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'
import '@testing-library/jest-dom'

promiseFinally.shim(Promise)
JavascriptTimeAgo.locale(en)

global.fetch = require('fetch-mock').sandbox()
global.fetch.config.overwriteRoutes = true
