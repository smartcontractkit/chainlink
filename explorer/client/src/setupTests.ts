import Enzyme from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import JavascriptTimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'

JavascriptTimeAgo.locale(en)

Enzyme.configure({ adapter: new Adapter() })
