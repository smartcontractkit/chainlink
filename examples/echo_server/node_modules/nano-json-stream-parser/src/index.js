// NJSP is just a straightforward monadic parser for JSON.

const Read = {};

const error = () => {
  throw "NoParse";
}

const readNonWhite = c => bind =>
  bind(c || Read, c => /\s/.test(c)
    ? readNonWhite("")(bind)
    : c);

const parseDigits = c => bind =>
  bind(c || Read, c => /[0-9]/.test(c)
    ? bind(parseDigits("")(bind), cs => [cs[0], c + cs[1]])
    : [c, ""]);

// Well fuck, this is why the do-notation exists
const parseNumber = c => bind =>
  // Parses sign
  bind(c || Read, c => bind(c === "-"
    ? ["","-"]
    : [c,""],
    ([c,sign]) =>
  // Parses the base digits
  bind(c || Read, c => bind(c === "0"
    ? ["","0"]
    : /[1-9]/.test(c)
      ? parseDigits(c)(bind)
      : error(),
    ([c, base]) =>
  // Parses the decimal digits
  bind(c || Read, c => bind(c === "."
    ? bind(parseDigits("")(bind), ([c,ds]) => [c, "." + ds])
    : [c, ""],
    ([c, frac]) =>
  // Parses the exponent
  bind(c || Read, c => bind(/[eE]/.test(c)
    ? bind(bind(Read, c => /[+-]/.test(c) ? ["", "e" + c] : [c, "e"]), ([c,es]) =>
      bind(parseDigits(c)(bind), ([c,ds]) => [c, es + ds]))
    : [c, ""],
    ([c, exp]) =>
  // Returns the number
  [c, Number(sign + base + frac + exp)]))))))));

const parseHex = bind =>
  bind(Read, h =>
    /[0-9a-fA-F]/.test(h)
      ? parseInt(h,16)
      : error());

const parseString = bind => 
  bind(bind(Read, c => c !== "\\" ? c :
    bind(Read, c => 
      /[\\"\/bfnrt]/.test(c)
        ? ({b:"\b",f:"\f",n:"\n",r:"\r",t:"\t","\\":"\\","/":"/",'"':''})[c]
        : /u/.test(c)
          ? bind(parseHex(bind), a =>
            bind(parseHex(bind), b =>
            bind(parseHex(bind), c =>
            bind(parseHex(bind), d =>
              String.fromCharCode(a*4096+b*256+c*16+d)))))
          : error())), c =>
  c === '"' ? "" : bind(parseString(bind), s => (c||'"') + s));

const parseArray = bind => {
  let array = [];
  const go = c => bind =>
    bind(c || Read, c => c === "]" ? array : 
    bind(parseValue(c)(bind), ([c,value]) => (array.push(value),
    bind(readNonWhite(c)(bind), c => /[,\]]/.test(c)
      ? go(c === "," ? "" : c)(bind)
      : error()))));
  return go("")(bind);
}

const parseObject = bind => {
  let object = {};
  const go = c => bind =>
    bind(c || Read, c => c === "}" ? object : 
    bind(readNonWhite(c)(bind), c => c !== '"' ? error() :
    bind(parseString(bind), key =>
    bind(readNonWhite("")(bind), c => c !== ':' ? error() :
    bind(parseValue("")(bind), ([c,val]) => (object[key] = val,
    bind(readNonWhite(c)(bind), c => /[,}]/.test(c)
      ? go(c === "," ? "" : c)(bind)
      : error())))))));
  return go("")(bind);
}

const parseExact = str => ret => bind => !str
  ? ret : bind(Read, c => c !== str[0]
    ? error()
    : parseExact(str.slice(1))(ret)(bind));

const parseValue = c => bind =>
  bind(readNonWhite(c)(bind), c => 
    c === "[" ? bind(parseArray(bind), v => ["",v]) :
    c === "{" ? bind(parseObject(bind), v => ["",v]) :
    c === '"' ? bind(parseString(bind), v => ["",v]) :
    c === 't' ? bind(parseExact("rue")(true)(bind), v => ["",v]) :
    c === 'f' ? bind(parseExact("alse")(false)(bind), v => ["",v]) :
    parseNumber(c)(bind));

const parseStream = parser => {
  const bind = (a, b) =>
    a === Read
      ? c => b(c)
      : typeof a === "function"
        ? c => bind(a(c), b)
        : b(a);
  return parser(bind);
}

const parser = parser => getJSON => {
  let s = parseStream(parser);
  const feed = str => { 
    for (let i = 0, l = str.length; i < l; ++i) {
      try {
        s = s(str[i]);
      } catch (e) {
        s = parseStream(parser);
        break;
      }
      if (typeof s !== "function") {
        getJSON(s[1]);
        s = parseStream(parser);
      }
    }
    return feed;
  }
  return feed;
}

module.exports = parser(parseValue(""));
