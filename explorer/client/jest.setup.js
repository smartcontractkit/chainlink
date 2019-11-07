const { configure } = require('enzyme')
const Adapter = require('enzyme-adapter-react-16')

configure({ adapter: new Adapter() })

const JavascriptTimeAgo = require('javascript-time-ago')
const en = require('javascript-time-ago/locale/en')

JavascriptTimeAgo.locale(en)
