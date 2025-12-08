// examples/workflows/ts/cron/mietekstefan2.js
var __defProp = Object.defineProperty;
var __export = (target, all) => {
  for (var name in all)
    __defProp(target, name, {
      get: all[name],
      enumerable: true,
      configurable: true,
      set: (newValue) => all[name] = () => newValue
    });
};
function isMessage(arg, schema) {
  const isMessage2 = arg !== null && typeof arg == "object" && "$typeName" in arg && typeof arg.$typeName == "string";
  if (!isMessage2) {
    return false;
  }
  if (schema === undefined) {
    return true;
  }
  return schema.typeName === arg.$typeName;
}
var ScalarType;
(function(ScalarType2) {
  ScalarType2[ScalarType2["DOUBLE"] = 1] = "DOUBLE";
  ScalarType2[ScalarType2["FLOAT"] = 2] = "FLOAT";
  ScalarType2[ScalarType2["INT64"] = 3] = "INT64";
  ScalarType2[ScalarType2["UINT64"] = 4] = "UINT64";
  ScalarType2[ScalarType2["INT32"] = 5] = "INT32";
  ScalarType2[ScalarType2["FIXED64"] = 6] = "FIXED64";
  ScalarType2[ScalarType2["FIXED32"] = 7] = "FIXED32";
  ScalarType2[ScalarType2["BOOL"] = 8] = "BOOL";
  ScalarType2[ScalarType2["STRING"] = 9] = "STRING";
  ScalarType2[ScalarType2["BYTES"] = 12] = "BYTES";
  ScalarType2[ScalarType2["UINT32"] = 13] = "UINT32";
  ScalarType2[ScalarType2["SFIXED32"] = 15] = "SFIXED32";
  ScalarType2[ScalarType2["SFIXED64"] = 16] = "SFIXED64";
  ScalarType2[ScalarType2["SINT32"] = 17] = "SINT32";
  ScalarType2[ScalarType2["SINT64"] = 18] = "SINT64";
})(ScalarType || (ScalarType = {}));
function varint64read() {
  let lowBits = 0;
  let highBits = 0;
  for (let shift = 0;shift < 28; shift += 7) {
    let b = this.buf[this.pos++];
    lowBits |= (b & 127) << shift;
    if ((b & 128) == 0) {
      this.assertBounds();
      return [lowBits, highBits];
    }
  }
  let middleByte = this.buf[this.pos++];
  lowBits |= (middleByte & 15) << 28;
  highBits = (middleByte & 112) >> 4;
  if ((middleByte & 128) == 0) {
    this.assertBounds();
    return [lowBits, highBits];
  }
  for (let shift = 3;shift <= 31; shift += 7) {
    let b = this.buf[this.pos++];
    highBits |= (b & 127) << shift;
    if ((b & 128) == 0) {
      this.assertBounds();
      return [lowBits, highBits];
    }
  }
  throw new Error("invalid varint");
}
function varint64write(lo, hi, bytes) {
  for (let i = 0;i < 28; i = i + 7) {
    const shift = lo >>> i;
    const hasNext = !(shift >>> 7 == 0 && hi == 0);
    const byte = (hasNext ? shift | 128 : shift) & 255;
    bytes.push(byte);
    if (!hasNext) {
      return;
    }
  }
  const splitBits = lo >>> 28 & 15 | (hi & 7) << 4;
  const hasMoreBits = !(hi >> 3 == 0);
  bytes.push((hasMoreBits ? splitBits | 128 : splitBits) & 255);
  if (!hasMoreBits) {
    return;
  }
  for (let i = 3;i < 31; i = i + 7) {
    const shift = hi >>> i;
    const hasNext = !(shift >>> 7 == 0);
    const byte = (hasNext ? shift | 128 : shift) & 255;
    bytes.push(byte);
    if (!hasNext) {
      return;
    }
  }
  bytes.push(hi >>> 31 & 1);
}
var TWO_PWR_32_DBL = 4294967296;
function int64FromString(dec) {
  const minus = dec[0] === "-";
  if (minus) {
    dec = dec.slice(1);
  }
  const base = 1e6;
  let lowBits = 0;
  let highBits = 0;
  function add1e6digit(begin, end) {
    const digit1e6 = Number(dec.slice(begin, end));
    highBits *= base;
    lowBits = lowBits * base + digit1e6;
    if (lowBits >= TWO_PWR_32_DBL) {
      highBits = highBits + (lowBits / TWO_PWR_32_DBL | 0);
      lowBits = lowBits % TWO_PWR_32_DBL;
    }
  }
  add1e6digit(-24, -18);
  add1e6digit(-18, -12);
  add1e6digit(-12, -6);
  add1e6digit(-6);
  return minus ? negate(lowBits, highBits) : newBits(lowBits, highBits);
}
function int64ToString(lo, hi) {
  let bits = newBits(lo, hi);
  const negative = bits.hi & 2147483648;
  if (negative) {
    bits = negate(bits.lo, bits.hi);
  }
  const result = uInt64ToString(bits.lo, bits.hi);
  return negative ? "-" + result : result;
}
function uInt64ToString(lo, hi) {
  ({ lo, hi } = toUnsigned(lo, hi));
  if (hi <= 2097151) {
    return String(TWO_PWR_32_DBL * hi + lo);
  }
  const low = lo & 16777215;
  const mid = (lo >>> 24 | hi << 8) & 16777215;
  const high = hi >> 16 & 65535;
  let digitA = low + mid * 6777216 + high * 6710656;
  let digitB = mid + high * 8147497;
  let digitC = high * 2;
  const base = 1e7;
  if (digitA >= base) {
    digitB += Math.floor(digitA / base);
    digitA %= base;
  }
  if (digitB >= base) {
    digitC += Math.floor(digitB / base);
    digitB %= base;
  }
  return digitC.toString() + decimalFrom1e7WithLeadingZeros(digitB) + decimalFrom1e7WithLeadingZeros(digitA);
}
function toUnsigned(lo, hi) {
  return { lo: lo >>> 0, hi: hi >>> 0 };
}
function newBits(lo, hi) {
  return { lo: lo | 0, hi: hi | 0 };
}
function negate(lowBits, highBits) {
  highBits = ~highBits;
  if (lowBits) {
    lowBits = ~lowBits + 1;
  } else {
    highBits += 1;
  }
  return newBits(lowBits, highBits);
}
var decimalFrom1e7WithLeadingZeros = (digit1e7) => {
  const partial = String(digit1e7);
  return "0000000".slice(partial.length) + partial;
};
function varint32write(value, bytes) {
  if (value >= 0) {
    while (value > 127) {
      bytes.push(value & 127 | 128);
      value = value >>> 7;
    }
    bytes.push(value);
  } else {
    for (let i = 0;i < 9; i++) {
      bytes.push(value & 127 | 128);
      value = value >> 7;
    }
    bytes.push(1);
  }
}
function varint32read() {
  let b = this.buf[this.pos++];
  let result = b & 127;
  if ((b & 128) == 0) {
    this.assertBounds();
    return result;
  }
  b = this.buf[this.pos++];
  result |= (b & 127) << 7;
  if ((b & 128) == 0) {
    this.assertBounds();
    return result;
  }
  b = this.buf[this.pos++];
  result |= (b & 127) << 14;
  if ((b & 128) == 0) {
    this.assertBounds();
    return result;
  }
  b = this.buf[this.pos++];
  result |= (b & 127) << 21;
  if ((b & 128) == 0) {
    this.assertBounds();
    return result;
  }
  b = this.buf[this.pos++];
  result |= (b & 15) << 28;
  for (let readBytes = 5;(b & 128) !== 0 && readBytes < 10; readBytes++)
    b = this.buf[this.pos++];
  if ((b & 128) != 0)
    throw new Error("invalid varint");
  this.assertBounds();
  return result >>> 0;
}
var protoInt64 = /* @__PURE__ */ makeInt64Support();
function makeInt64Support() {
  const dv = new DataView(new ArrayBuffer(8));
  const ok = typeof BigInt === "function" && typeof dv.getBigInt64 === "function" && typeof dv.getBigUint64 === "function" && typeof dv.setBigInt64 === "function" && typeof dv.setBigUint64 === "function" && (typeof process != "object" || typeof process.env != "object" || process.env.BUF_BIGINT_DISABLE !== "1");
  if (ok) {
    const MIN = BigInt("-9223372036854775808");
    const MAX = BigInt("9223372036854775807");
    const UMIN = BigInt("0");
    const UMAX = BigInt("18446744073709551615");
    return {
      zero: BigInt(0),
      supported: true,
      parse(value) {
        const bi = typeof value == "bigint" ? value : BigInt(value);
        if (bi > MAX || bi < MIN) {
          throw new Error(`invalid int64: ${value}`);
        }
        return bi;
      },
      uParse(value) {
        const bi = typeof value == "bigint" ? value : BigInt(value);
        if (bi > UMAX || bi < UMIN) {
          throw new Error(`invalid uint64: ${value}`);
        }
        return bi;
      },
      enc(value) {
        dv.setBigInt64(0, this.parse(value), true);
        return {
          lo: dv.getInt32(0, true),
          hi: dv.getInt32(4, true)
        };
      },
      uEnc(value) {
        dv.setBigInt64(0, this.uParse(value), true);
        return {
          lo: dv.getInt32(0, true),
          hi: dv.getInt32(4, true)
        };
      },
      dec(lo, hi) {
        dv.setInt32(0, lo, true);
        dv.setInt32(4, hi, true);
        return dv.getBigInt64(0, true);
      },
      uDec(lo, hi) {
        dv.setInt32(0, lo, true);
        dv.setInt32(4, hi, true);
        return dv.getBigUint64(0, true);
      }
    };
  }
  return {
    zero: "0",
    supported: false,
    parse(value) {
      if (typeof value != "string") {
        value = value.toString();
      }
      assertInt64String(value);
      return value;
    },
    uParse(value) {
      if (typeof value != "string") {
        value = value.toString();
      }
      assertUInt64String(value);
      return value;
    },
    enc(value) {
      if (typeof value != "string") {
        value = value.toString();
      }
      assertInt64String(value);
      return int64FromString(value);
    },
    uEnc(value) {
      if (typeof value != "string") {
        value = value.toString();
      }
      assertUInt64String(value);
      return int64FromString(value);
    },
    dec(lo, hi) {
      return int64ToString(lo, hi);
    },
    uDec(lo, hi) {
      return uInt64ToString(lo, hi);
    }
  };
}
function assertInt64String(value) {
  if (!/^-?[0-9]+$/.test(value)) {
    throw new Error("invalid int64: " + value);
  }
}
function assertUInt64String(value) {
  if (!/^[0-9]+$/.test(value)) {
    throw new Error("invalid uint64: " + value);
  }
}
function scalarZeroValue(type, longAsString) {
  switch (type) {
    case ScalarType.STRING:
      return "";
    case ScalarType.BOOL:
      return false;
    case ScalarType.DOUBLE:
    case ScalarType.FLOAT:
      return 0;
    case ScalarType.INT64:
    case ScalarType.UINT64:
    case ScalarType.SFIXED64:
    case ScalarType.FIXED64:
    case ScalarType.SINT64:
      return longAsString ? "0" : protoInt64.zero;
    case ScalarType.BYTES:
      return new Uint8Array(0);
    default:
      return 0;
  }
}
function isScalarZeroValue(type, value) {
  switch (type) {
    case ScalarType.BOOL:
      return value === false;
    case ScalarType.STRING:
      return value === "";
    case ScalarType.BYTES:
      return value instanceof Uint8Array && !value.byteLength;
    default:
      return value == 0;
  }
}
var IMPLICIT = 2;
var unsafeLocal = Symbol.for("reflect unsafe local");
function unsafeOneofCase(target, oneof) {
  const c = target[oneof.localName].case;
  if (c === undefined) {
    return c;
  }
  return oneof.fields.find((f) => f.localName === c);
}
function unsafeIsSet(target, field) {
  const name = field.localName;
  if (field.oneof) {
    return target[field.oneof.localName].case === name;
  }
  if (field.presence != IMPLICIT) {
    return target[name] !== undefined && Object.prototype.hasOwnProperty.call(target, name);
  }
  switch (field.fieldKind) {
    case "list":
      return target[name].length > 0;
    case "map":
      return Object.keys(target[name]).length > 0;
    case "scalar":
      return !isScalarZeroValue(field.scalar, target[name]);
    case "enum":
      return target[name] !== field.enum.values[0].number;
  }
  throw new Error("message field with implicit presence");
}
function unsafeIsSetExplicit(target, localName) {
  return Object.prototype.hasOwnProperty.call(target, localName) && target[localName] !== undefined;
}
function unsafeGet(target, field) {
  if (field.oneof) {
    const oneof = target[field.oneof.localName];
    if (oneof.case === field.localName) {
      return oneof.value;
    }
    return;
  }
  return target[field.localName];
}
function unsafeSet(target, field, value) {
  if (field.oneof) {
    target[field.oneof.localName] = {
      case: field.localName,
      value
    };
  } else {
    target[field.localName] = value;
  }
}
function unsafeClear(target, field) {
  const name = field.localName;
  if (field.oneof) {
    const oneofLocalName = field.oneof.localName;
    if (target[oneofLocalName].case === name) {
      target[oneofLocalName] = { case: undefined };
    }
  } else if (field.presence != IMPLICIT) {
    delete target[name];
  } else {
    switch (field.fieldKind) {
      case "map":
        target[name] = {};
        break;
      case "list":
        target[name] = [];
        break;
      case "enum":
        target[name] = field.enum.values[0].number;
        break;
      case "scalar":
        target[name] = scalarZeroValue(field.scalar, field.longAsString);
        break;
    }
  }
}
function isObject(arg) {
  return arg !== null && typeof arg == "object" && !Array.isArray(arg);
}
function isReflectList(arg, field) {
  var _a, _b, _c, _d;
  if (isObject(arg) && unsafeLocal in arg && "add" in arg && "field" in arg && typeof arg.field == "function") {
    if (field !== undefined) {
      const a = field;
      const b = arg.field();
      return a.listKind == b.listKind && a.scalar === b.scalar && ((_a = a.message) === null || _a === undefined ? undefined : _a.typeName) === ((_b = b.message) === null || _b === undefined ? undefined : _b.typeName) && ((_c = a.enum) === null || _c === undefined ? undefined : _c.typeName) === ((_d = b.enum) === null || _d === undefined ? undefined : _d.typeName);
    }
    return true;
  }
  return false;
}
function isReflectMap(arg, field) {
  var _a, _b, _c, _d;
  if (isObject(arg) && unsafeLocal in arg && "has" in arg && "field" in arg && typeof arg.field == "function") {
    if (field !== undefined) {
      const a = field, b = arg.field();
      return a.mapKey === b.mapKey && a.mapKind == b.mapKind && a.scalar === b.scalar && ((_a = a.message) === null || _a === undefined ? undefined : _a.typeName) === ((_b = b.message) === null || _b === undefined ? undefined : _b.typeName) && ((_c = a.enum) === null || _c === undefined ? undefined : _c.typeName) === ((_d = b.enum) === null || _d === undefined ? undefined : _d.typeName);
    }
    return true;
  }
  return false;
}
function isReflectMessage(arg, messageDesc) {
  return isObject(arg) && unsafeLocal in arg && "desc" in arg && isObject(arg.desc) && arg.desc.kind === "message" && (messageDesc === undefined || arg.desc.typeName == messageDesc.typeName);
}
function isWrapper(arg) {
  return isWrapperTypeName(arg.$typeName);
}
function isWrapperDesc(messageDesc) {
  const f = messageDesc.fields[0];
  return isWrapperTypeName(messageDesc.typeName) && f !== undefined && f.fieldKind == "scalar" && f.name == "value" && f.number == 1;
}
function isWrapperTypeName(name) {
  return name.startsWith("google.protobuf.") && [
    "DoubleValue",
    "FloatValue",
    "Int64Value",
    "UInt64Value",
    "Int32Value",
    "UInt32Value",
    "BoolValue",
    "StringValue",
    "BytesValue"
  ].includes(name.substring(16));
}
var EDITION_PROTO3 = 999;
var EDITION_PROTO2 = 998;
var IMPLICIT2 = 2;
function create(schema, init) {
  if (isMessage(init, schema)) {
    return init;
  }
  const message = createZeroMessage(schema);
  if (init !== undefined) {
    initMessage(schema, message, init);
  }
  return message;
}
function initMessage(messageDesc, message, init) {
  for (const member of messageDesc.members) {
    let value = init[member.localName];
    if (value == null) {
      continue;
    }
    let field;
    if (member.kind == "oneof") {
      const oneofField = unsafeOneofCase(init, member);
      if (!oneofField) {
        continue;
      }
      field = oneofField;
      value = unsafeGet(init, oneofField);
    } else {
      field = member;
    }
    switch (field.fieldKind) {
      case "message":
        value = toMessage(field, value);
        break;
      case "scalar":
        value = initScalar(field, value);
        break;
      case "list":
        value = initList(field, value);
        break;
      case "map":
        value = initMap(field, value);
        break;
    }
    unsafeSet(message, field, value);
  }
  return message;
}
function initScalar(field, value) {
  if (field.scalar == ScalarType.BYTES) {
    return toU8Arr(value);
  }
  return value;
}
function initMap(field, value) {
  if (isObject(value)) {
    if (field.scalar == ScalarType.BYTES) {
      return convertObjectValues(value, toU8Arr);
    }
    if (field.mapKind == "message") {
      return convertObjectValues(value, (val) => toMessage(field, val));
    }
  }
  return value;
}
function initList(field, value) {
  if (Array.isArray(value)) {
    if (field.scalar == ScalarType.BYTES) {
      return value.map(toU8Arr);
    }
    if (field.listKind == "message") {
      return value.map((item) => toMessage(field, item));
    }
  }
  return value;
}
function toMessage(field, value) {
  if (field.fieldKind == "message" && !field.oneof && isWrapperDesc(field.message)) {
    return initScalar(field.message.fields[0], value);
  }
  if (isObject(value)) {
    if (field.message.typeName == "google.protobuf.Struct" && field.parent.typeName !== "google.protobuf.Value") {
      return value;
    }
    if (!isMessage(value, field.message)) {
      return create(field.message, value);
    }
  }
  return value;
}
function toU8Arr(value) {
  return Array.isArray(value) ? new Uint8Array(value) : value;
}
function convertObjectValues(obj, fn) {
  const ret = {};
  for (const entry of Object.entries(obj)) {
    ret[entry[0]] = fn(entry[1]);
  }
  return ret;
}
var tokenZeroMessageField = Symbol();
var messagePrototypes = new WeakMap;
function createZeroMessage(desc) {
  let msg;
  if (!needsPrototypeChain(desc)) {
    msg = {
      $typeName: desc.typeName
    };
    for (const member of desc.members) {
      if (member.kind == "oneof" || member.presence == IMPLICIT2) {
        msg[member.localName] = createZeroField(member);
      }
    }
  } else {
    const cached = messagePrototypes.get(desc);
    let prototype;
    let members;
    if (cached) {
      ({ prototype, members } = cached);
    } else {
      prototype = {};
      members = new Set;
      for (const member of desc.members) {
        if (member.kind == "oneof") {
          continue;
        }
        if (member.fieldKind != "scalar" && member.fieldKind != "enum") {
          continue;
        }
        if (member.presence == IMPLICIT2) {
          continue;
        }
        members.add(member);
        prototype[member.localName] = createZeroField(member);
      }
      messagePrototypes.set(desc, { prototype, members });
    }
    msg = Object.create(prototype);
    msg.$typeName = desc.typeName;
    for (const member of desc.members) {
      if (members.has(member)) {
        continue;
      }
      if (member.kind == "field") {
        if (member.fieldKind == "message") {
          continue;
        }
        if (member.fieldKind == "scalar" || member.fieldKind == "enum") {
          if (member.presence != IMPLICIT2) {
            continue;
          }
        }
      }
      msg[member.localName] = createZeroField(member);
    }
  }
  return msg;
}
function needsPrototypeChain(desc) {
  switch (desc.file.edition) {
    case EDITION_PROTO3:
      return false;
    case EDITION_PROTO2:
      return true;
    default:
      return desc.fields.some((f) => f.presence != IMPLICIT2 && f.fieldKind != "message" && !f.oneof);
  }
}
function createZeroField(field) {
  if (field.kind == "oneof") {
    return { case: undefined };
  }
  if (field.fieldKind == "list") {
    return [];
  }
  if (field.fieldKind == "map") {
    return {};
  }
  if (field.fieldKind == "message") {
    return tokenZeroMessageField;
  }
  const defaultValue = field.getDefaultValue();
  if (defaultValue !== undefined) {
    return field.fieldKind == "scalar" && field.longAsString ? defaultValue.toString() : defaultValue;
  }
  return field.fieldKind == "scalar" ? scalarZeroValue(field.scalar, field.longAsString) : field.enum.values[0].number;
}
var errorNames = [
  "FieldValueInvalidError",
  "FieldListRangeError",
  "ForeignFieldError"
];

class FieldError extends Error {
  constructor(fieldOrOneof, message, name = "FieldValueInvalidError") {
    super(message);
    this.name = name;
    this.field = () => fieldOrOneof;
  }
}
function isFieldError(arg) {
  return arg instanceof Error && errorNames.includes(arg.name) && "field" in arg && typeof arg.field == "function";
}
var symbol = Symbol.for("@bufbuild/protobuf/text-encoding");
function getTextEncoding() {
  if (globalThis[symbol] == undefined) {
    const te = new globalThis.TextEncoder;
    const td = new globalThis.TextDecoder;
    globalThis[symbol] = {
      encodeUtf8(text) {
        return te.encode(text);
      },
      decodeUtf8(bytes) {
        return td.decode(bytes);
      },
      checkUtf8(text) {
        try {
          encodeURIComponent(text);
          return true;
        } catch (_) {
          return false;
        }
      }
    };
  }
  return globalThis[symbol];
}
var WireType;
(function(WireType2) {
  WireType2[WireType2["Varint"] = 0] = "Varint";
  WireType2[WireType2["Bit64"] = 1] = "Bit64";
  WireType2[WireType2["LengthDelimited"] = 2] = "LengthDelimited";
  WireType2[WireType2["StartGroup"] = 3] = "StartGroup";
  WireType2[WireType2["EndGroup"] = 4] = "EndGroup";
  WireType2[WireType2["Bit32"] = 5] = "Bit32";
})(WireType || (WireType = {}));
var FLOAT32_MAX = 340282346638528860000000000000000000000;
var FLOAT32_MIN = -340282346638528860000000000000000000000;
var UINT32_MAX = 4294967295;
var INT32_MAX = 2147483647;
var INT32_MIN = -2147483648;

class BinaryWriter {
  constructor(encodeUtf8 = getTextEncoding().encodeUtf8) {
    this.encodeUtf8 = encodeUtf8;
    this.stack = [];
    this.chunks = [];
    this.buf = [];
  }
  finish() {
    if (this.buf.length) {
      this.chunks.push(new Uint8Array(this.buf));
      this.buf = [];
    }
    let len = 0;
    for (let i = 0;i < this.chunks.length; i++)
      len += this.chunks[i].length;
    let bytes = new Uint8Array(len);
    let offset = 0;
    for (let i = 0;i < this.chunks.length; i++) {
      bytes.set(this.chunks[i], offset);
      offset += this.chunks[i].length;
    }
    this.chunks = [];
    return bytes;
  }
  fork() {
    this.stack.push({ chunks: this.chunks, buf: this.buf });
    this.chunks = [];
    this.buf = [];
    return this;
  }
  join() {
    let chunk = this.finish();
    let prev = this.stack.pop();
    if (!prev)
      throw new Error("invalid state, fork stack empty");
    this.chunks = prev.chunks;
    this.buf = prev.buf;
    this.uint32(chunk.byteLength);
    return this.raw(chunk);
  }
  tag(fieldNo, type) {
    return this.uint32((fieldNo << 3 | type) >>> 0);
  }
  raw(chunk) {
    if (this.buf.length) {
      this.chunks.push(new Uint8Array(this.buf));
      this.buf = [];
    }
    this.chunks.push(chunk);
    return this;
  }
  uint32(value) {
    assertUInt32(value);
    while (value > 127) {
      this.buf.push(value & 127 | 128);
      value = value >>> 7;
    }
    this.buf.push(value);
    return this;
  }
  int32(value) {
    assertInt32(value);
    varint32write(value, this.buf);
    return this;
  }
  bool(value) {
    this.buf.push(value ? 1 : 0);
    return this;
  }
  bytes(value) {
    this.uint32(value.byteLength);
    return this.raw(value);
  }
  string(value) {
    let chunk = this.encodeUtf8(value);
    this.uint32(chunk.byteLength);
    return this.raw(chunk);
  }
  float(value) {
    assertFloat32(value);
    let chunk = new Uint8Array(4);
    new DataView(chunk.buffer).setFloat32(0, value, true);
    return this.raw(chunk);
  }
  double(value) {
    let chunk = new Uint8Array(8);
    new DataView(chunk.buffer).setFloat64(0, value, true);
    return this.raw(chunk);
  }
  fixed32(value) {
    assertUInt32(value);
    let chunk = new Uint8Array(4);
    new DataView(chunk.buffer).setUint32(0, value, true);
    return this.raw(chunk);
  }
  sfixed32(value) {
    assertInt32(value);
    let chunk = new Uint8Array(4);
    new DataView(chunk.buffer).setInt32(0, value, true);
    return this.raw(chunk);
  }
  sint32(value) {
    assertInt32(value);
    value = (value << 1 ^ value >> 31) >>> 0;
    varint32write(value, this.buf);
    return this;
  }
  sfixed64(value) {
    let chunk = new Uint8Array(8), view = new DataView(chunk.buffer), tc = protoInt64.enc(value);
    view.setInt32(0, tc.lo, true);
    view.setInt32(4, tc.hi, true);
    return this.raw(chunk);
  }
  fixed64(value) {
    let chunk = new Uint8Array(8), view = new DataView(chunk.buffer), tc = protoInt64.uEnc(value);
    view.setInt32(0, tc.lo, true);
    view.setInt32(4, tc.hi, true);
    return this.raw(chunk);
  }
  int64(value) {
    let tc = protoInt64.enc(value);
    varint64write(tc.lo, tc.hi, this.buf);
    return this;
  }
  sint64(value) {
    const tc = protoInt64.enc(value), sign = tc.hi >> 31, lo = tc.lo << 1 ^ sign, hi = (tc.hi << 1 | tc.lo >>> 31) ^ sign;
    varint64write(lo, hi, this.buf);
    return this;
  }
  uint64(value) {
    const tc = protoInt64.uEnc(value);
    varint64write(tc.lo, tc.hi, this.buf);
    return this;
  }
}

class BinaryReader {
  constructor(buf, decodeUtf8 = getTextEncoding().decodeUtf8) {
    this.decodeUtf8 = decodeUtf8;
    this.varint64 = varint64read;
    this.uint32 = varint32read;
    this.buf = buf;
    this.len = buf.length;
    this.pos = 0;
    this.view = new DataView(buf.buffer, buf.byteOffset, buf.byteLength);
  }
  tag() {
    let tag = this.uint32(), fieldNo = tag >>> 3, wireType = tag & 7;
    if (fieldNo <= 0 || wireType < 0 || wireType > 5)
      throw new Error("illegal tag: field no " + fieldNo + " wire type " + wireType);
    return [fieldNo, wireType];
  }
  skip(wireType, fieldNo) {
    let start = this.pos;
    switch (wireType) {
      case WireType.Varint:
        while (this.buf[this.pos++] & 128) {}
        break;
      case WireType.Bit64:
        this.pos += 4;
      case WireType.Bit32:
        this.pos += 4;
        break;
      case WireType.LengthDelimited:
        let len = this.uint32();
        this.pos += len;
        break;
      case WireType.StartGroup:
        for (;; ) {
          const [fn, wt] = this.tag();
          if (wt === WireType.EndGroup) {
            if (fieldNo !== undefined && fn !== fieldNo) {
              throw new Error("invalid end group tag");
            }
            break;
          }
          this.skip(wt, fn);
        }
        break;
      default:
        throw new Error("cant skip wire type " + wireType);
    }
    this.assertBounds();
    return this.buf.subarray(start, this.pos);
  }
  assertBounds() {
    if (this.pos > this.len)
      throw new RangeError("premature EOF");
  }
  int32() {
    return this.uint32() | 0;
  }
  sint32() {
    let zze = this.uint32();
    return zze >>> 1 ^ -(zze & 1);
  }
  int64() {
    return protoInt64.dec(...this.varint64());
  }
  uint64() {
    return protoInt64.uDec(...this.varint64());
  }
  sint64() {
    let [lo, hi] = this.varint64();
    let s = -(lo & 1);
    lo = (lo >>> 1 | (hi & 1) << 31) ^ s;
    hi = hi >>> 1 ^ s;
    return protoInt64.dec(lo, hi);
  }
  bool() {
    let [lo, hi] = this.varint64();
    return lo !== 0 || hi !== 0;
  }
  fixed32() {
    return this.view.getUint32((this.pos += 4) - 4, true);
  }
  sfixed32() {
    return this.view.getInt32((this.pos += 4) - 4, true);
  }
  fixed64() {
    return protoInt64.uDec(this.sfixed32(), this.sfixed32());
  }
  sfixed64() {
    return protoInt64.dec(this.sfixed32(), this.sfixed32());
  }
  float() {
    return this.view.getFloat32((this.pos += 4) - 4, true);
  }
  double() {
    return this.view.getFloat64((this.pos += 8) - 8, true);
  }
  bytes() {
    let len = this.uint32(), start = this.pos;
    this.pos += len;
    this.assertBounds();
    return this.buf.subarray(start, start + len);
  }
  string() {
    return this.decodeUtf8(this.bytes());
  }
}
function assertInt32(arg) {
  if (typeof arg == "string") {
    arg = Number(arg);
  } else if (typeof arg != "number") {
    throw new Error("invalid int32: " + typeof arg);
  }
  if (!Number.isInteger(arg) || arg > INT32_MAX || arg < INT32_MIN)
    throw new Error("invalid int32: " + arg);
}
function assertUInt32(arg) {
  if (typeof arg == "string") {
    arg = Number(arg);
  } else if (typeof arg != "number") {
    throw new Error("invalid uint32: " + typeof arg);
  }
  if (!Number.isInteger(arg) || arg > UINT32_MAX || arg < 0)
    throw new Error("invalid uint32: " + arg);
}
function assertFloat32(arg) {
  if (typeof arg == "string") {
    const o = arg;
    arg = Number(arg);
    if (Number.isNaN(arg) && o !== "NaN") {
      throw new Error("invalid float32: " + o);
    }
  } else if (typeof arg != "number") {
    throw new Error("invalid float32: " + typeof arg);
  }
  if (Number.isFinite(arg) && (arg > FLOAT32_MAX || arg < FLOAT32_MIN))
    throw new Error("invalid float32: " + arg);
}
function checkField(field, value) {
  const check = field.fieldKind == "list" ? isReflectList(value, field) : field.fieldKind == "map" ? isReflectMap(value, field) : checkSingular(field, value);
  if (check === true) {
    return;
  }
  let reason;
  switch (field.fieldKind) {
    case "list":
      reason = `expected ${formatReflectList(field)}, got ${formatVal(value)}`;
      break;
    case "map":
      reason = `expected ${formatReflectMap(field)}, got ${formatVal(value)}`;
      break;
    default: {
      reason = reasonSingular(field, value, check);
    }
  }
  return new FieldError(field, reason);
}
function checkListItem(field, index, value) {
  const check = checkSingular(field, value);
  if (check !== true) {
    return new FieldError(field, `list item #${index + 1}: ${reasonSingular(field, value, check)}`);
  }
  return;
}
function checkMapEntry(field, key, value) {
  const checkKey = checkScalarValue(key, field.mapKey);
  if (checkKey !== true) {
    return new FieldError(field, `invalid map key: ${reasonSingular({ scalar: field.mapKey }, key, checkKey)}`);
  }
  const checkVal = checkSingular(field, value);
  if (checkVal !== true) {
    return new FieldError(field, `map entry ${formatVal(key)}: ${reasonSingular(field, value, checkVal)}`);
  }
  return;
}
function checkSingular(field, value) {
  if (field.scalar !== undefined) {
    return checkScalarValue(value, field.scalar);
  }
  if (field.enum !== undefined) {
    if (field.enum.open) {
      return Number.isInteger(value);
    }
    return field.enum.values.some((v) => v.number === value);
  }
  return isReflectMessage(value, field.message);
}
function checkScalarValue(value, scalar) {
  switch (scalar) {
    case ScalarType.DOUBLE:
      return typeof value == "number";
    case ScalarType.FLOAT:
      if (typeof value != "number") {
        return false;
      }
      if (Number.isNaN(value) || !Number.isFinite(value)) {
        return true;
      }
      if (value > FLOAT32_MAX || value < FLOAT32_MIN) {
        return `${value.toFixed()} out of range`;
      }
      return true;
    case ScalarType.INT32:
    case ScalarType.SFIXED32:
    case ScalarType.SINT32:
      if (typeof value !== "number" || !Number.isInteger(value)) {
        return false;
      }
      if (value > INT32_MAX || value < INT32_MIN) {
        return `${value.toFixed()} out of range`;
      }
      return true;
    case ScalarType.FIXED32:
    case ScalarType.UINT32:
      if (typeof value !== "number" || !Number.isInteger(value)) {
        return false;
      }
      if (value > UINT32_MAX || value < 0) {
        return `${value.toFixed()} out of range`;
      }
      return true;
    case ScalarType.BOOL:
      return typeof value == "boolean";
    case ScalarType.STRING:
      if (typeof value != "string") {
        return false;
      }
      return getTextEncoding().checkUtf8(value) || "invalid UTF8";
    case ScalarType.BYTES:
      return value instanceof Uint8Array;
    case ScalarType.INT64:
    case ScalarType.SFIXED64:
    case ScalarType.SINT64:
      if (typeof value == "bigint" || typeof value == "number" || typeof value == "string" && value.length > 0) {
        try {
          protoInt64.parse(value);
          return true;
        } catch (_) {
          return `${value} out of range`;
        }
      }
      return false;
    case ScalarType.FIXED64:
    case ScalarType.UINT64:
      if (typeof value == "bigint" || typeof value == "number" || typeof value == "string" && value.length > 0) {
        try {
          protoInt64.uParse(value);
          return true;
        } catch (_) {
          return `${value} out of range`;
        }
      }
      return false;
  }
}
function reasonSingular(field, val, details) {
  details = typeof details == "string" ? `: ${details}` : `, got ${formatVal(val)}`;
  if (field.scalar !== undefined) {
    return `expected ${scalarTypeDescription(field.scalar)}` + details;
  }
  if (field.enum !== undefined) {
    return `expected ${field.enum.toString()}` + details;
  }
  return `expected ${formatReflectMessage(field.message)}` + details;
}
function formatVal(val) {
  switch (typeof val) {
    case "object":
      if (val === null) {
        return "null";
      }
      if (val instanceof Uint8Array) {
        return `Uint8Array(${val.length})`;
      }
      if (Array.isArray(val)) {
        return `Array(${val.length})`;
      }
      if (isReflectList(val)) {
        return formatReflectList(val.field());
      }
      if (isReflectMap(val)) {
        return formatReflectMap(val.field());
      }
      if (isReflectMessage(val)) {
        return formatReflectMessage(val.desc);
      }
      if (isMessage(val)) {
        return `message ${val.$typeName}`;
      }
      return "object";
    case "string":
      return val.length > 30 ? "string" : `"${val.split('"').join("\\\"")}"`;
    case "boolean":
      return String(val);
    case "number":
      return String(val);
    case "bigint":
      return String(val) + "n";
    default:
      return typeof val;
  }
}
function formatReflectMessage(desc) {
  return `ReflectMessage (${desc.typeName})`;
}
function formatReflectList(field) {
  switch (field.listKind) {
    case "message":
      return `ReflectList (${field.message.toString()})`;
    case "enum":
      return `ReflectList (${field.enum.toString()})`;
    case "scalar":
      return `ReflectList (${ScalarType[field.scalar]})`;
  }
}
function formatReflectMap(field) {
  switch (field.mapKind) {
    case "message":
      return `ReflectMap (${ScalarType[field.mapKey]}, ${field.message.toString()})`;
    case "enum":
      return `ReflectMap (${ScalarType[field.mapKey]}, ${field.enum.toString()})`;
    case "scalar":
      return `ReflectMap (${ScalarType[field.mapKey]}, ${ScalarType[field.scalar]})`;
  }
}
function scalarTypeDescription(scalar) {
  switch (scalar) {
    case ScalarType.STRING:
      return "string";
    case ScalarType.BOOL:
      return "boolean";
    case ScalarType.INT64:
    case ScalarType.SINT64:
    case ScalarType.SFIXED64:
      return "bigint (int64)";
    case ScalarType.UINT64:
    case ScalarType.FIXED64:
      return "bigint (uint64)";
    case ScalarType.BYTES:
      return "Uint8Array";
    case ScalarType.DOUBLE:
      return "number (float64)";
    case ScalarType.FLOAT:
      return "number (float32)";
    case ScalarType.FIXED32:
    case ScalarType.UINT32:
      return "number (uint32)";
    case ScalarType.INT32:
    case ScalarType.SFIXED32:
    case ScalarType.SINT32:
      return "number (int32)";
  }
}
function reflect(messageDesc, message, check = true) {
  return new ReflectMessageImpl(messageDesc, message, check);
}

class ReflectMessageImpl {
  get sortedFields() {
    var _a;
    return (_a = this._sortedFields) !== null && _a !== undefined ? _a : this._sortedFields = this.desc.fields.concat().sort((a, b) => a.number - b.number);
  }
  constructor(messageDesc, message, check = true) {
    this.lists = new Map;
    this.maps = new Map;
    this.check = check;
    this.desc = messageDesc;
    this.message = this[unsafeLocal] = message !== null && message !== undefined ? message : create(messageDesc);
    this.fields = messageDesc.fields;
    this.oneofs = messageDesc.oneofs;
    this.members = messageDesc.members;
  }
  findNumber(number) {
    if (!this._fieldsByNumber) {
      this._fieldsByNumber = new Map(this.desc.fields.map((f) => [f.number, f]));
    }
    return this._fieldsByNumber.get(number);
  }
  oneofCase(oneof) {
    assertOwn(this.message, oneof);
    return unsafeOneofCase(this.message, oneof);
  }
  isSet(field) {
    assertOwn(this.message, field);
    return unsafeIsSet(this.message, field);
  }
  clear(field) {
    assertOwn(this.message, field);
    unsafeClear(this.message, field);
  }
  get(field) {
    assertOwn(this.message, field);
    const value = unsafeGet(this.message, field);
    switch (field.fieldKind) {
      case "list":
        let list = this.lists.get(field);
        if (!list || list[unsafeLocal] !== value) {
          this.lists.set(field, list = new ReflectListImpl(field, value, this.check));
        }
        return list;
      case "map":
        let map = this.maps.get(field);
        if (!map || map[unsafeLocal] !== value) {
          this.maps.set(field, map = new ReflectMapImpl(field, value, this.check));
        }
        return map;
      case "message":
        return messageToReflect(field, value, this.check);
      case "scalar":
        return value === undefined ? scalarZeroValue(field.scalar, false) : longToReflect(field, value);
      case "enum":
        return value !== null && value !== undefined ? value : field.enum.values[0].number;
    }
  }
  set(field, value) {
    assertOwn(this.message, field);
    if (this.check) {
      const err = checkField(field, value);
      if (err) {
        throw err;
      }
    }
    let local;
    if (field.fieldKind == "message") {
      local = messageToLocal(field, value);
    } else if (isReflectMap(value) || isReflectList(value)) {
      local = value[unsafeLocal];
    } else {
      local = longToLocal(field, value);
    }
    unsafeSet(this.message, field, local);
  }
  getUnknown() {
    return this.message.$unknown;
  }
  setUnknown(value) {
    this.message.$unknown = value;
  }
}
function assertOwn(owner, member) {
  if (member.parent.typeName !== owner.$typeName) {
    throw new FieldError(member, `cannot use ${member.toString()} with message ${owner.$typeName}`, "ForeignFieldError");
  }
}

class ReflectListImpl {
  field() {
    return this._field;
  }
  get size() {
    return this._arr.length;
  }
  constructor(field, unsafeInput, check) {
    this._field = field;
    this._arr = this[unsafeLocal] = unsafeInput;
    this.check = check;
  }
  get(index) {
    const item = this._arr[index];
    return item === undefined ? undefined : listItemToReflect(this._field, item, this.check);
  }
  set(index, item) {
    if (index < 0 || index >= this._arr.length) {
      throw new FieldError(this._field, `list item #${index + 1}: out of range`);
    }
    if (this.check) {
      const err = checkListItem(this._field, index, item);
      if (err) {
        throw err;
      }
    }
    this._arr[index] = listItemToLocal(this._field, item);
  }
  add(item) {
    if (this.check) {
      const err = checkListItem(this._field, this._arr.length, item);
      if (err) {
        throw err;
      }
    }
    this._arr.push(listItemToLocal(this._field, item));
    return;
  }
  clear() {
    this._arr.splice(0, this._arr.length);
  }
  [Symbol.iterator]() {
    return this.values();
  }
  keys() {
    return this._arr.keys();
  }
  *values() {
    for (const item of this._arr) {
      yield listItemToReflect(this._field, item, this.check);
    }
  }
  *entries() {
    for (let i = 0;i < this._arr.length; i++) {
      yield [i, listItemToReflect(this._field, this._arr[i], this.check)];
    }
  }
}

class ReflectMapImpl {
  constructor(field, unsafeInput, check = true) {
    this.obj = this[unsafeLocal] = unsafeInput !== null && unsafeInput !== undefined ? unsafeInput : {};
    this.check = check;
    this._field = field;
  }
  field() {
    return this._field;
  }
  set(key, value) {
    if (this.check) {
      const err = checkMapEntry(this._field, key, value);
      if (err) {
        throw err;
      }
    }
    this.obj[mapKeyToLocal(key)] = mapValueToLocal(this._field, value);
    return this;
  }
  delete(key) {
    const k = mapKeyToLocal(key);
    const has = Object.prototype.hasOwnProperty.call(this.obj, k);
    if (has) {
      delete this.obj[k];
    }
    return has;
  }
  clear() {
    for (const key of Object.keys(this.obj)) {
      delete this.obj[key];
    }
  }
  get(key) {
    let val = this.obj[mapKeyToLocal(key)];
    if (val !== undefined) {
      val = mapValueToReflect(this._field, val, this.check);
    }
    return val;
  }
  has(key) {
    return Object.prototype.hasOwnProperty.call(this.obj, mapKeyToLocal(key));
  }
  *keys() {
    for (const objKey of Object.keys(this.obj)) {
      yield mapKeyToReflect(objKey, this._field.mapKey);
    }
  }
  *entries() {
    for (const objEntry of Object.entries(this.obj)) {
      yield [
        mapKeyToReflect(objEntry[0], this._field.mapKey),
        mapValueToReflect(this._field, objEntry[1], this.check)
      ];
    }
  }
  [Symbol.iterator]() {
    return this.entries();
  }
  get size() {
    return Object.keys(this.obj).length;
  }
  *values() {
    for (const val of Object.values(this.obj)) {
      yield mapValueToReflect(this._field, val, this.check);
    }
  }
  forEach(callbackfn, thisArg) {
    for (const mapEntry of this.entries()) {
      callbackfn.call(thisArg, mapEntry[1], mapEntry[0], this);
    }
  }
}
function messageToLocal(field, value) {
  if (!isReflectMessage(value)) {
    return value;
  }
  if (isWrapper(value.message) && !field.oneof && field.fieldKind == "message") {
    return value.message.value;
  }
  if (value.desc.typeName == "google.protobuf.Struct" && field.parent.typeName != "google.protobuf.Value") {
    return wktStructToLocal(value.message);
  }
  return value.message;
}
function messageToReflect(field, value, check) {
  if (value !== undefined) {
    if (isWrapperDesc(field.message) && !field.oneof && field.fieldKind == "message") {
      value = {
        $typeName: field.message.typeName,
        value: longToReflect(field.message.fields[0], value)
      };
    } else if (field.message.typeName == "google.protobuf.Struct" && field.parent.typeName != "google.protobuf.Value" && isObject(value)) {
      value = wktStructToReflect(value);
    }
  }
  return new ReflectMessageImpl(field.message, value, check);
}
function listItemToLocal(field, value) {
  if (field.listKind == "message") {
    return messageToLocal(field, value);
  }
  return longToLocal(field, value);
}
function listItemToReflect(field, value, check) {
  if (field.listKind == "message") {
    return messageToReflect(field, value, check);
  }
  return longToReflect(field, value);
}
function mapValueToLocal(field, value) {
  if (field.mapKind == "message") {
    return messageToLocal(field, value);
  }
  return longToLocal(field, value);
}
function mapValueToReflect(field, value, check) {
  if (field.mapKind == "message") {
    return messageToReflect(field, value, check);
  }
  return value;
}
function mapKeyToLocal(key) {
  return typeof key == "string" || typeof key == "number" ? key : String(key);
}
function mapKeyToReflect(key, type) {
  switch (type) {
    case ScalarType.STRING:
      return key;
    case ScalarType.INT32:
    case ScalarType.FIXED32:
    case ScalarType.UINT32:
    case ScalarType.SFIXED32:
    case ScalarType.SINT32: {
      const n = Number.parseInt(key);
      if (Number.isFinite(n)) {
        return n;
      }
      break;
    }
    case ScalarType.BOOL:
      switch (key) {
        case "true":
          return true;
        case "false":
          return false;
      }
      break;
    case ScalarType.UINT64:
    case ScalarType.FIXED64:
      try {
        return protoInt64.uParse(key);
      } catch (_a) {}
      break;
    default:
      try {
        return protoInt64.parse(key);
      } catch (_b) {}
      break;
  }
  return key;
}
function longToReflect(field, value) {
  switch (field.scalar) {
    case ScalarType.INT64:
    case ScalarType.SFIXED64:
    case ScalarType.SINT64:
      if ("longAsString" in field && field.longAsString && typeof value == "string") {
        value = protoInt64.parse(value);
      }
      break;
    case ScalarType.FIXED64:
    case ScalarType.UINT64:
      if ("longAsString" in field && field.longAsString && typeof value == "string") {
        value = protoInt64.uParse(value);
      }
      break;
  }
  return value;
}
function longToLocal(field, value) {
  switch (field.scalar) {
    case ScalarType.INT64:
    case ScalarType.SFIXED64:
    case ScalarType.SINT64:
      if ("longAsString" in field && field.longAsString) {
        value = String(value);
      } else if (typeof value == "string" || typeof value == "number") {
        value = protoInt64.parse(value);
      }
      break;
    case ScalarType.FIXED64:
    case ScalarType.UINT64:
      if ("longAsString" in field && field.longAsString) {
        value = String(value);
      } else if (typeof value == "string" || typeof value == "number") {
        value = protoInt64.uParse(value);
      }
      break;
  }
  return value;
}
function wktStructToReflect(json) {
  const struct = {
    $typeName: "google.protobuf.Struct",
    fields: {}
  };
  if (isObject(json)) {
    for (const [k, v] of Object.entries(json)) {
      struct.fields[k] = wktValueToReflect(v);
    }
  }
  return struct;
}
function wktStructToLocal(val) {
  const json = {};
  for (const [k, v] of Object.entries(val.fields)) {
    json[k] = wktValueToLocal(v);
  }
  return json;
}
function wktValueToLocal(val) {
  switch (val.kind.case) {
    case "structValue":
      return wktStructToLocal(val.kind.value);
    case "listValue":
      return val.kind.value.values.map(wktValueToLocal);
    case "nullValue":
    case undefined:
      return null;
    default:
      return val.kind.value;
  }
}
function wktValueToReflect(json) {
  const value = {
    $typeName: "google.protobuf.Value",
    kind: { case: undefined }
  };
  switch (typeof json) {
    case "number":
      value.kind = { case: "numberValue", value: json };
      break;
    case "string":
      value.kind = { case: "stringValue", value: json };
      break;
    case "boolean":
      value.kind = { case: "boolValue", value: json };
      break;
    case "object":
      if (json === null) {
        const nullValue = 0;
        value.kind = { case: "nullValue", value: nullValue };
      } else if (Array.isArray(json)) {
        const listValue = {
          $typeName: "google.protobuf.ListValue",
          values: []
        };
        if (Array.isArray(json)) {
          for (const e of json) {
            listValue.values.push(wktValueToReflect(e));
          }
        }
        value.kind = {
          case: "listValue",
          value: listValue
        };
      } else {
        value.kind = {
          case: "structValue",
          value: wktStructToReflect(json)
        };
      }
      break;
  }
  return value;
}
function base64Decode(base64Str) {
  const table = getDecodeTable();
  let es = base64Str.length * 3 / 4;
  if (base64Str[base64Str.length - 2] == "=")
    es -= 2;
  else if (base64Str[base64Str.length - 1] == "=")
    es -= 1;
  let bytes = new Uint8Array(es), bytePos = 0, groupPos = 0, b, p = 0;
  for (let i = 0;i < base64Str.length; i++) {
    b = table[base64Str.charCodeAt(i)];
    if (b === undefined) {
      switch (base64Str[i]) {
        case "=":
          groupPos = 0;
        case `
`:
        case "\r":
        case "\t":
        case " ":
          continue;
        default:
          throw Error("invalid base64 string");
      }
    }
    switch (groupPos) {
      case 0:
        p = b;
        groupPos = 1;
        break;
      case 1:
        bytes[bytePos++] = p << 2 | (b & 48) >> 4;
        p = b;
        groupPos = 2;
        break;
      case 2:
        bytes[bytePos++] = (p & 15) << 4 | (b & 60) >> 2;
        p = b;
        groupPos = 3;
        break;
      case 3:
        bytes[bytePos++] = (p & 3) << 6 | b;
        groupPos = 0;
        break;
    }
  }
  if (groupPos == 1)
    throw Error("invalid base64 string");
  return bytes.subarray(0, bytePos);
}
var encodeTableStd;
var encodeTableUrl;
var decodeTable;
function getEncodeTable(encoding) {
  if (!encodeTableStd) {
    encodeTableStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/".split("");
    encodeTableUrl = encodeTableStd.slice(0, -2).concat("-", "_");
  }
  return encoding == "url" ? encodeTableUrl : encodeTableStd;
}
function getDecodeTable() {
  if (!decodeTable) {
    decodeTable = [];
    const encodeTable = getEncodeTable("std");
    for (let i = 0;i < encodeTable.length; i++)
      decodeTable[encodeTable[i].charCodeAt(0)] = i;
    decodeTable[45] = encodeTable.indexOf("+");
    decodeTable[95] = encodeTable.indexOf("/");
  }
  return decodeTable;
}
function protoCamelCase(snakeCase) {
  let capNext = false;
  const b = [];
  for (let i = 0;i < snakeCase.length; i++) {
    let c = snakeCase.charAt(i);
    switch (c) {
      case "_":
        capNext = true;
        break;
      case "0":
      case "1":
      case "2":
      case "3":
      case "4":
      case "5":
      case "6":
      case "7":
      case "8":
      case "9":
        b.push(c);
        capNext = false;
        break;
      default:
        if (capNext) {
          capNext = false;
          c = c.toUpperCase();
        }
        b.push(c);
        break;
    }
  }
  return b.join("");
}
var reservedObjectProperties = new Set([
  "constructor",
  "toString",
  "toJSON",
  "valueOf"
]);
function safeObjectProperty(name) {
  return reservedObjectProperties.has(name) ? name + "$" : name;
}
function restoreJsonNames(message) {
  for (const f of message.field) {
    if (!unsafeIsSetExplicit(f, "jsonName")) {
      f.jsonName = protoCamelCase(f.name);
    }
  }
  message.nestedType.forEach(restoreJsonNames);
}
function parseTextFormatEnumValue(descEnum, value) {
  const enumValue = descEnum.values.find((v) => v.name === value);
  if (!enumValue) {
    throw new Error(`cannot parse ${descEnum} default value: ${value}`);
  }
  return enumValue.number;
}
function parseTextFormatScalarValue(type, value) {
  switch (type) {
    case ScalarType.STRING:
      return value;
    case ScalarType.BYTES: {
      const u = unescapeBytesDefaultValue(value);
      if (u === false) {
        throw new Error(`cannot parse ${ScalarType[type]} default value: ${value}`);
      }
      return u;
    }
    case ScalarType.INT64:
    case ScalarType.SFIXED64:
    case ScalarType.SINT64:
      return protoInt64.parse(value);
    case ScalarType.UINT64:
    case ScalarType.FIXED64:
      return protoInt64.uParse(value);
    case ScalarType.DOUBLE:
    case ScalarType.FLOAT:
      switch (value) {
        case "inf":
          return Number.POSITIVE_INFINITY;
        case "-inf":
          return Number.NEGATIVE_INFINITY;
        case "nan":
          return Number.NaN;
        default:
          return parseFloat(value);
      }
    case ScalarType.BOOL:
      return value === "true";
    case ScalarType.INT32:
    case ScalarType.UINT32:
    case ScalarType.SINT32:
    case ScalarType.FIXED32:
    case ScalarType.SFIXED32:
      return parseInt(value, 10);
  }
}
function unescapeBytesDefaultValue(str) {
  const b = [];
  const input = {
    tail: str,
    c: "",
    next() {
      if (this.tail.length == 0) {
        return false;
      }
      this.c = this.tail[0];
      this.tail = this.tail.substring(1);
      return true;
    },
    take(n) {
      if (this.tail.length >= n) {
        const r = this.tail.substring(0, n);
        this.tail = this.tail.substring(n);
        return r;
      }
      return false;
    }
  };
  while (input.next()) {
    switch (input.c) {
      case "\\":
        if (input.next()) {
          switch (input.c) {
            case "\\":
              b.push(input.c.charCodeAt(0));
              break;
            case "b":
              b.push(8);
              break;
            case "f":
              b.push(12);
              break;
            case "n":
              b.push(10);
              break;
            case "r":
              b.push(13);
              break;
            case "t":
              b.push(9);
              break;
            case "v":
              b.push(11);
              break;
            case "0":
            case "1":
            case "2":
            case "3":
            case "4":
            case "5":
            case "6":
            case "7": {
              const s = input.c;
              const t = input.take(2);
              if (t === false) {
                return false;
              }
              const n = parseInt(s + t, 8);
              if (Number.isNaN(n)) {
                return false;
              }
              b.push(n);
              break;
            }
            case "x": {
              const s = input.c;
              const t = input.take(2);
              if (t === false) {
                return false;
              }
              const n = parseInt(s + t, 16);
              if (Number.isNaN(n)) {
                return false;
              }
              b.push(n);
              break;
            }
            case "u": {
              const s = input.c;
              const t = input.take(4);
              if (t === false) {
                return false;
              }
              const n = parseInt(s + t, 16);
              if (Number.isNaN(n)) {
                return false;
              }
              const chunk = new Uint8Array(4);
              const view = new DataView(chunk.buffer);
              view.setInt32(0, n, true);
              b.push(chunk[0], chunk[1], chunk[2], chunk[3]);
              break;
            }
            case "U": {
              const s = input.c;
              const t = input.take(8);
              if (t === false) {
                return false;
              }
              const tc = protoInt64.uEnc(s + t);
              const chunk = new Uint8Array(8);
              const view = new DataView(chunk.buffer);
              view.setInt32(0, tc.lo, true);
              view.setInt32(4, tc.hi, true);
              b.push(chunk[0], chunk[1], chunk[2], chunk[3], chunk[4], chunk[5], chunk[6], chunk[7]);
              break;
            }
          }
        }
        break;
      default:
        b.push(input.c.charCodeAt(0));
    }
  }
  return new Uint8Array(b);
}
function* nestedTypes(desc) {
  switch (desc.kind) {
    case "file":
      for (const message of desc.messages) {
        yield message;
        yield* nestedTypes(message);
      }
      yield* desc.enums;
      yield* desc.services;
      yield* desc.extensions;
      break;
    case "message":
      for (const message of desc.nestedMessages) {
        yield message;
        yield* nestedTypes(message);
      }
      yield* desc.nestedEnums;
      yield* desc.nestedExtensions;
      break;
  }
}
function createFileRegistry(...args) {
  const registry = createBaseRegistry();
  if (!args.length) {
    return registry;
  }
  if ("$typeName" in args[0] && args[0].$typeName == "google.protobuf.FileDescriptorSet") {
    for (const file of args[0].file) {
      addFile(file, registry);
    }
    return registry;
  }
  if ("$typeName" in args[0]) {
    let recurseDeps = function(file) {
      const deps = [];
      for (const protoFileName of file.dependency) {
        if (registry.getFile(protoFileName) != null) {
          continue;
        }
        if (seen.has(protoFileName)) {
          continue;
        }
        const dep = resolve(protoFileName);
        if (!dep) {
          throw new Error(`Unable to resolve ${protoFileName}, imported by ${file.name}`);
        }
        if ("kind" in dep) {
          registry.addFile(dep, false, true);
        } else {
          seen.add(dep.name);
          deps.push(dep);
        }
      }
      return deps.concat(...deps.map(recurseDeps));
    };
    const input = args[0];
    const resolve = args[1];
    const seen = new Set;
    for (const file of [input, ...recurseDeps(input)].reverse()) {
      addFile(file, registry);
    }
  } else {
    for (const fileReg of args) {
      for (const file of fileReg.files) {
        registry.addFile(file);
      }
    }
  }
  return registry;
}
function createBaseRegistry() {
  const types = new Map;
  const extendees = new Map;
  const files = new Map;
  return {
    kind: "registry",
    types,
    extendees,
    [Symbol.iterator]() {
      return types.values();
    },
    get files() {
      return files.values();
    },
    addFile(file, skipTypes, withDeps) {
      files.set(file.proto.name, file);
      if (!skipTypes) {
        for (const type of nestedTypes(file)) {
          this.add(type);
        }
      }
      if (withDeps) {
        for (const f of file.dependencies) {
          this.addFile(f, skipTypes, withDeps);
        }
      }
    },
    add(desc) {
      if (desc.kind == "extension") {
        let numberToExt = extendees.get(desc.extendee.typeName);
        if (!numberToExt) {
          extendees.set(desc.extendee.typeName, numberToExt = new Map);
        }
        numberToExt.set(desc.number, desc);
      }
      types.set(desc.typeName, desc);
    },
    get(typeName) {
      return types.get(typeName);
    },
    getFile(fileName) {
      return files.get(fileName);
    },
    getMessage(typeName) {
      const t = types.get(typeName);
      return (t === null || t === undefined ? undefined : t.kind) == "message" ? t : undefined;
    },
    getEnum(typeName) {
      const t = types.get(typeName);
      return (t === null || t === undefined ? undefined : t.kind) == "enum" ? t : undefined;
    },
    getExtension(typeName) {
      const t = types.get(typeName);
      return (t === null || t === undefined ? undefined : t.kind) == "extension" ? t : undefined;
    },
    getExtensionFor(extendee, no) {
      var _a;
      return (_a = extendees.get(extendee.typeName)) === null || _a === undefined ? undefined : _a.get(no);
    },
    getService(typeName) {
      const t = types.get(typeName);
      return (t === null || t === undefined ? undefined : t.kind) == "service" ? t : undefined;
    }
  };
}
var EDITION_PROTO22 = 998;
var EDITION_PROTO32 = 999;
var TYPE_STRING = 9;
var TYPE_GROUP = 10;
var TYPE_MESSAGE = 11;
var TYPE_BYTES = 12;
var TYPE_ENUM = 14;
var LABEL_REPEATED = 3;
var LABEL_REQUIRED = 2;
var JS_STRING = 1;
var IDEMPOTENCY_UNKNOWN = 0;
var EXPLICIT = 1;
var IMPLICIT3 = 2;
var LEGACY_REQUIRED = 3;
var PACKED = 1;
var DELIMITED = 2;
var OPEN = 1;
var featureDefaults = {
  998: {
    fieldPresence: 1,
    enumType: 2,
    repeatedFieldEncoding: 2,
    utf8Validation: 3,
    messageEncoding: 1,
    jsonFormat: 2,
    enforceNamingStyle: 2,
    defaultSymbolVisibility: 1
  },
  999: {
    fieldPresence: 2,
    enumType: 1,
    repeatedFieldEncoding: 1,
    utf8Validation: 2,
    messageEncoding: 1,
    jsonFormat: 1,
    enforceNamingStyle: 2,
    defaultSymbolVisibility: 1
  },
  1000: {
    fieldPresence: 1,
    enumType: 1,
    repeatedFieldEncoding: 1,
    utf8Validation: 2,
    messageEncoding: 1,
    jsonFormat: 1,
    enforceNamingStyle: 2,
    defaultSymbolVisibility: 1
  }
};
function addFile(proto, reg) {
  var _a, _b;
  const file = {
    kind: "file",
    proto,
    deprecated: (_b = (_a = proto.options) === null || _a === undefined ? undefined : _a.deprecated) !== null && _b !== undefined ? _b : false,
    edition: getFileEdition(proto),
    name: proto.name.replace(/\.proto$/, ""),
    dependencies: findFileDependencies(proto, reg),
    enums: [],
    messages: [],
    extensions: [],
    services: [],
    toString() {
      return `file ${proto.name}`;
    }
  };
  const mapEntriesStore = new Map;
  const mapEntries = {
    get(typeName) {
      return mapEntriesStore.get(typeName);
    },
    add(desc) {
      var _a2;
      assert(((_a2 = desc.proto.options) === null || _a2 === undefined ? undefined : _a2.mapEntry) === true);
      mapEntriesStore.set(desc.typeName, desc);
    }
  };
  for (const enumProto of proto.enumType) {
    addEnum(enumProto, file, undefined, reg);
  }
  for (const messageProto of proto.messageType) {
    addMessage(messageProto, file, undefined, reg, mapEntries);
  }
  for (const serviceProto of proto.service) {
    addService(serviceProto, file, reg);
  }
  addExtensions(file, reg);
  for (const mapEntry of mapEntriesStore.values()) {
    addFields(mapEntry, reg, mapEntries);
  }
  for (const message of file.messages) {
    addFields(message, reg, mapEntries);
    addExtensions(message, reg);
  }
  reg.addFile(file, true);
}
function addExtensions(desc, reg) {
  switch (desc.kind) {
    case "file":
      for (const proto of desc.proto.extension) {
        const ext = newField(proto, desc, reg);
        desc.extensions.push(ext);
        reg.add(ext);
      }
      break;
    case "message":
      for (const proto of desc.proto.extension) {
        const ext = newField(proto, desc, reg);
        desc.nestedExtensions.push(ext);
        reg.add(ext);
      }
      for (const message of desc.nestedMessages) {
        addExtensions(message, reg);
      }
      break;
  }
}
function addFields(message, reg, mapEntries) {
  const allOneofs = message.proto.oneofDecl.map((proto) => newOneof(proto, message));
  const oneofsSeen = new Set;
  for (const proto of message.proto.field) {
    const oneof = findOneof(proto, allOneofs);
    const field = newField(proto, message, reg, oneof, mapEntries);
    message.fields.push(field);
    message.field[field.localName] = field;
    if (oneof === undefined) {
      message.members.push(field);
    } else {
      oneof.fields.push(field);
      if (!oneofsSeen.has(oneof)) {
        oneofsSeen.add(oneof);
        message.members.push(oneof);
      }
    }
  }
  for (const oneof of allOneofs.filter((o) => oneofsSeen.has(o))) {
    message.oneofs.push(oneof);
  }
  for (const child of message.nestedMessages) {
    addFields(child, reg, mapEntries);
  }
}
function addEnum(proto, file, parent, reg) {
  var _a, _b, _c, _d, _e;
  const sharedPrefix = findEnumSharedPrefix(proto.name, proto.value);
  const desc = {
    kind: "enum",
    proto,
    deprecated: (_b = (_a = proto.options) === null || _a === undefined ? undefined : _a.deprecated) !== null && _b !== undefined ? _b : false,
    file,
    parent,
    open: true,
    name: proto.name,
    typeName: makeTypeName(proto, parent, file),
    value: {},
    values: [],
    sharedPrefix,
    toString() {
      return `enum ${this.typeName}`;
    }
  };
  desc.open = isEnumOpen(desc);
  reg.add(desc);
  for (const p of proto.value) {
    const name = p.name;
    desc.values.push(desc.value[p.number] = {
      kind: "enum_value",
      proto: p,
      deprecated: (_d = (_c = p.options) === null || _c === undefined ? undefined : _c.deprecated) !== null && _d !== undefined ? _d : false,
      parent: desc,
      name,
      localName: safeObjectProperty(sharedPrefix == undefined ? name : name.substring(sharedPrefix.length)),
      number: p.number,
      toString() {
        return `enum value ${desc.typeName}.${name}`;
      }
    });
  }
  ((_e = parent === null || parent === undefined ? undefined : parent.nestedEnums) !== null && _e !== undefined ? _e : file.enums).push(desc);
}
function addMessage(proto, file, parent, reg, mapEntries) {
  var _a, _b, _c, _d;
  const desc = {
    kind: "message",
    proto,
    deprecated: (_b = (_a = proto.options) === null || _a === undefined ? undefined : _a.deprecated) !== null && _b !== undefined ? _b : false,
    file,
    parent,
    name: proto.name,
    typeName: makeTypeName(proto, parent, file),
    fields: [],
    field: {},
    oneofs: [],
    members: [],
    nestedEnums: [],
    nestedMessages: [],
    nestedExtensions: [],
    toString() {
      return `message ${this.typeName}`;
    }
  };
  if (((_c = proto.options) === null || _c === undefined ? undefined : _c.mapEntry) === true) {
    mapEntries.add(desc);
  } else {
    ((_d = parent === null || parent === undefined ? undefined : parent.nestedMessages) !== null && _d !== undefined ? _d : file.messages).push(desc);
    reg.add(desc);
  }
  for (const enumProto of proto.enumType) {
    addEnum(enumProto, file, desc, reg);
  }
  for (const messageProto of proto.nestedType) {
    addMessage(messageProto, file, desc, reg, mapEntries);
  }
}
function addService(proto, file, reg) {
  var _a, _b;
  const desc = {
    kind: "service",
    proto,
    deprecated: (_b = (_a = proto.options) === null || _a === undefined ? undefined : _a.deprecated) !== null && _b !== undefined ? _b : false,
    file,
    name: proto.name,
    typeName: makeTypeName(proto, undefined, file),
    methods: [],
    method: {},
    toString() {
      return `service ${this.typeName}`;
    }
  };
  file.services.push(desc);
  reg.add(desc);
  for (const methodProto of proto.method) {
    const method = newMethod(methodProto, desc, reg);
    desc.methods.push(method);
    desc.method[method.localName] = method;
  }
}
function newMethod(proto, parent, reg) {
  var _a, _b, _c, _d;
  let methodKind;
  if (proto.clientStreaming && proto.serverStreaming) {
    methodKind = "bidi_streaming";
  } else if (proto.clientStreaming) {
    methodKind = "client_streaming";
  } else if (proto.serverStreaming) {
    methodKind = "server_streaming";
  } else {
    methodKind = "unary";
  }
  const input = reg.getMessage(trimLeadingDot(proto.inputType));
  const output = reg.getMessage(trimLeadingDot(proto.outputType));
  assert(input, `invalid MethodDescriptorProto: input_type ${proto.inputType} not found`);
  assert(output, `invalid MethodDescriptorProto: output_type ${proto.inputType} not found`);
  const name = proto.name;
  return {
    kind: "rpc",
    proto,
    deprecated: (_b = (_a = proto.options) === null || _a === undefined ? undefined : _a.deprecated) !== null && _b !== undefined ? _b : false,
    parent,
    name,
    localName: safeObjectProperty(name.length ? safeObjectProperty(name[0].toLowerCase() + name.substring(1)) : name),
    methodKind,
    input,
    output,
    idempotency: (_d = (_c = proto.options) === null || _c === undefined ? undefined : _c.idempotencyLevel) !== null && _d !== undefined ? _d : IDEMPOTENCY_UNKNOWN,
    toString() {
      return `rpc ${parent.typeName}.${name}`;
    }
  };
}
function newOneof(proto, parent) {
  return {
    kind: "oneof",
    proto,
    deprecated: false,
    parent,
    fields: [],
    name: proto.name,
    localName: safeObjectProperty(protoCamelCase(proto.name)),
    toString() {
      return `oneof ${parent.typeName}.${this.name}`;
    }
  };
}
function newField(proto, parentOrFile, reg, oneof, mapEntries) {
  var _a, _b, _c;
  const isExtension = mapEntries === undefined;
  const field = {
    kind: "field",
    proto,
    deprecated: (_b = (_a = proto.options) === null || _a === undefined ? undefined : _a.deprecated) !== null && _b !== undefined ? _b : false,
    name: proto.name,
    number: proto.number,
    scalar: undefined,
    message: undefined,
    enum: undefined,
    presence: getFieldPresence(proto, oneof, isExtension, parentOrFile),
    listKind: undefined,
    mapKind: undefined,
    mapKey: undefined,
    delimitedEncoding: undefined,
    packed: undefined,
    longAsString: false,
    getDefaultValue: undefined
  };
  if (isExtension) {
    const file = parentOrFile.kind == "file" ? parentOrFile : parentOrFile.file;
    const parent = parentOrFile.kind == "file" ? undefined : parentOrFile;
    const typeName = makeTypeName(proto, parent, file);
    field.kind = "extension";
    field.file = file;
    field.parent = parent;
    field.oneof = undefined;
    field.typeName = typeName;
    field.jsonName = `[${typeName}]`;
    field.toString = () => `extension ${typeName}`;
    const extendee = reg.getMessage(trimLeadingDot(proto.extendee));
    assert(extendee, `invalid FieldDescriptorProto: extendee ${proto.extendee} not found`);
    field.extendee = extendee;
  } else {
    const parent = parentOrFile;
    assert(parent.kind == "message");
    field.parent = parent;
    field.oneof = oneof;
    field.localName = oneof ? protoCamelCase(proto.name) : safeObjectProperty(protoCamelCase(proto.name));
    field.jsonName = proto.jsonName;
    field.toString = () => `field ${parent.typeName}.${proto.name}`;
  }
  const label = proto.label;
  const type = proto.type;
  const jstype = (_c = proto.options) === null || _c === undefined ? undefined : _c.jstype;
  if (label === LABEL_REPEATED) {
    const mapEntry = type == TYPE_MESSAGE ? mapEntries === null || mapEntries === undefined ? undefined : mapEntries.get(trimLeadingDot(proto.typeName)) : undefined;
    if (mapEntry) {
      field.fieldKind = "map";
      const { key, value } = findMapEntryFields(mapEntry);
      field.mapKey = key.scalar;
      field.mapKind = value.fieldKind;
      field.message = value.message;
      field.delimitedEncoding = false;
      field.enum = value.enum;
      field.scalar = value.scalar;
      return field;
    }
    field.fieldKind = "list";
    switch (type) {
      case TYPE_MESSAGE:
      case TYPE_GROUP:
        field.listKind = "message";
        field.message = reg.getMessage(trimLeadingDot(proto.typeName));
        assert(field.message);
        field.delimitedEncoding = isDelimitedEncoding(proto, parentOrFile);
        break;
      case TYPE_ENUM:
        field.listKind = "enum";
        field.enum = reg.getEnum(trimLeadingDot(proto.typeName));
        assert(field.enum);
        break;
      default:
        field.listKind = "scalar";
        field.scalar = type;
        field.longAsString = jstype == JS_STRING;
        break;
    }
    field.packed = isPackedField(proto, parentOrFile);
    return field;
  }
  switch (type) {
    case TYPE_MESSAGE:
    case TYPE_GROUP:
      field.fieldKind = "message";
      field.message = reg.getMessage(trimLeadingDot(proto.typeName));
      assert(field.message, `invalid FieldDescriptorProto: type_name ${proto.typeName} not found`);
      field.delimitedEncoding = isDelimitedEncoding(proto, parentOrFile);
      field.getDefaultValue = () => {
        return;
      };
      break;
    case TYPE_ENUM: {
      const enumeration = reg.getEnum(trimLeadingDot(proto.typeName));
      assert(enumeration !== undefined, `invalid FieldDescriptorProto: type_name ${proto.typeName} not found`);
      field.fieldKind = "enum";
      field.enum = reg.getEnum(trimLeadingDot(proto.typeName));
      field.getDefaultValue = () => {
        return unsafeIsSetExplicit(proto, "defaultValue") ? parseTextFormatEnumValue(enumeration, proto.defaultValue) : undefined;
      };
      break;
    }
    default: {
      field.fieldKind = "scalar";
      field.scalar = type;
      field.longAsString = jstype == JS_STRING;
      field.getDefaultValue = () => {
        return unsafeIsSetExplicit(proto, "defaultValue") ? parseTextFormatScalarValue(type, proto.defaultValue) : undefined;
      };
      break;
    }
  }
  return field;
}
function getFileEdition(proto) {
  switch (proto.syntax) {
    case "":
    case "proto2":
      return EDITION_PROTO22;
    case "proto3":
      return EDITION_PROTO32;
    case "editions":
      if (proto.edition in featureDefaults) {
        return proto.edition;
      }
      throw new Error(`${proto.name}: unsupported edition`);
    default:
      throw new Error(`${proto.name}: unsupported syntax "${proto.syntax}"`);
  }
}
function findFileDependencies(proto, reg) {
  return proto.dependency.map((wantName) => {
    const dep = reg.getFile(wantName);
    if (!dep) {
      throw new Error(`Cannot find ${wantName}, imported by ${proto.name}`);
    }
    return dep;
  });
}
function findEnumSharedPrefix(enumName, values) {
  const prefix = camelToSnakeCase(enumName) + "_";
  for (const value of values) {
    if (!value.name.toLowerCase().startsWith(prefix)) {
      return;
    }
    const shortName = value.name.substring(prefix.length);
    if (shortName.length == 0) {
      return;
    }
    if (/^\d/.test(shortName)) {
      return;
    }
  }
  return prefix;
}
function camelToSnakeCase(camel) {
  return (camel.substring(0, 1) + camel.substring(1).replace(/[A-Z]/g, (c) => "_" + c)).toLowerCase();
}
function makeTypeName(proto, parent, file) {
  let typeName;
  if (parent) {
    typeName = `${parent.typeName}.${proto.name}`;
  } else if (file.proto.package.length > 0) {
    typeName = `${file.proto.package}.${proto.name}`;
  } else {
    typeName = `${proto.name}`;
  }
  return typeName;
}
function trimLeadingDot(typeName) {
  return typeName.startsWith(".") ? typeName.substring(1) : typeName;
}
function findOneof(proto, allOneofs) {
  if (!unsafeIsSetExplicit(proto, "oneofIndex")) {
    return;
  }
  if (proto.proto3Optional) {
    return;
  }
  const oneof = allOneofs[proto.oneofIndex];
  assert(oneof, `invalid FieldDescriptorProto: oneof #${proto.oneofIndex} for field #${proto.number} not found`);
  return oneof;
}
function getFieldPresence(proto, oneof, isExtension, parent) {
  if (proto.label == LABEL_REQUIRED) {
    return LEGACY_REQUIRED;
  }
  if (proto.label == LABEL_REPEATED) {
    return IMPLICIT3;
  }
  if (!!oneof || proto.proto3Optional) {
    return EXPLICIT;
  }
  if (isExtension) {
    return EXPLICIT;
  }
  const resolved = resolveFeature("fieldPresence", { proto, parent });
  if (resolved == IMPLICIT3 && (proto.type == TYPE_MESSAGE || proto.type == TYPE_GROUP)) {
    return EXPLICIT;
  }
  return resolved;
}
function isPackedField(proto, parent) {
  if (proto.label != LABEL_REPEATED) {
    return false;
  }
  switch (proto.type) {
    case TYPE_STRING:
    case TYPE_BYTES:
    case TYPE_GROUP:
    case TYPE_MESSAGE:
      return false;
  }
  const o = proto.options;
  if (o && unsafeIsSetExplicit(o, "packed")) {
    return o.packed;
  }
  return PACKED == resolveFeature("repeatedFieldEncoding", {
    proto,
    parent
  });
}
function findMapEntryFields(mapEntry) {
  const key = mapEntry.fields.find((f) => f.number === 1);
  const value = mapEntry.fields.find((f) => f.number === 2);
  assert(key && key.fieldKind == "scalar" && key.scalar != ScalarType.BYTES && key.scalar != ScalarType.FLOAT && key.scalar != ScalarType.DOUBLE && value && value.fieldKind != "list" && value.fieldKind != "map");
  return { key, value };
}
function isEnumOpen(desc) {
  var _a;
  return OPEN == resolveFeature("enumType", {
    proto: desc.proto,
    parent: (_a = desc.parent) !== null && _a !== undefined ? _a : desc.file
  });
}
function isDelimitedEncoding(proto, parent) {
  if (proto.type == TYPE_GROUP) {
    return true;
  }
  return DELIMITED == resolveFeature("messageEncoding", {
    proto,
    parent
  });
}
function resolveFeature(name, ref) {
  var _a, _b;
  const featureSet = (_a = ref.proto.options) === null || _a === undefined ? undefined : _a.features;
  if (featureSet) {
    const val = featureSet[name];
    if (val != 0) {
      return val;
    }
  }
  if ("kind" in ref) {
    if (ref.kind == "message") {
      return resolveFeature(name, (_b = ref.parent) !== null && _b !== undefined ? _b : ref.file);
    }
    const editionDefaults = featureDefaults[ref.edition];
    if (!editionDefaults) {
      throw new Error(`feature default for edition ${ref.edition} not found`);
    }
    return editionDefaults[name];
  }
  return resolveFeature(name, ref.parent);
}
function assert(condition, msg) {
  if (!condition) {
    throw new Error(msg);
  }
}
function boot(boot2) {
  const root = bootFileDescriptorProto(boot2);
  root.messageType.forEach(restoreJsonNames);
  const reg = createFileRegistry(root, () => {
    return;
  });
  return reg.getFile(root.name);
}
function bootFileDescriptorProto(init) {
  const proto = Object.create({
    syntax: "",
    edition: 0
  });
  return Object.assign(proto, Object.assign(Object.assign({ $typeName: "google.protobuf.FileDescriptorProto", dependency: [], publicDependency: [], weakDependency: [], optionDependency: [], service: [], extension: [] }, init), { messageType: init.messageType.map(bootDescriptorProto), enumType: init.enumType.map(bootEnumDescriptorProto) }));
}
function bootDescriptorProto(init) {
  var _a, _b, _c, _d, _e, _f, _g, _h;
  const proto = Object.create({
    visibility: 0
  });
  return Object.assign(proto, {
    $typeName: "google.protobuf.DescriptorProto",
    name: init.name,
    field: (_b = (_a = init.field) === null || _a === undefined ? undefined : _a.map(bootFieldDescriptorProto)) !== null && _b !== undefined ? _b : [],
    extension: [],
    nestedType: (_d = (_c = init.nestedType) === null || _c === undefined ? undefined : _c.map(bootDescriptorProto)) !== null && _d !== undefined ? _d : [],
    enumType: (_f = (_e = init.enumType) === null || _e === undefined ? undefined : _e.map(bootEnumDescriptorProto)) !== null && _f !== undefined ? _f : [],
    extensionRange: (_h = (_g = init.extensionRange) === null || _g === undefined ? undefined : _g.map((e) => Object.assign({ $typeName: "google.protobuf.DescriptorProto.ExtensionRange" }, e))) !== null && _h !== undefined ? _h : [],
    oneofDecl: [],
    reservedRange: [],
    reservedName: []
  });
}
function bootFieldDescriptorProto(init) {
  const proto = Object.create({
    label: 1,
    typeName: "",
    extendee: "",
    defaultValue: "",
    oneofIndex: 0,
    jsonName: "",
    proto3Optional: false
  });
  return Object.assign(proto, Object.assign(Object.assign({ $typeName: "google.protobuf.FieldDescriptorProto" }, init), { options: init.options ? bootFieldOptions(init.options) : undefined }));
}
function bootFieldOptions(init) {
  var _a, _b, _c;
  const proto = Object.create({
    ctype: 0,
    packed: false,
    jstype: 0,
    lazy: false,
    unverifiedLazy: false,
    deprecated: false,
    weak: false,
    debugRedact: false,
    retention: 0
  });
  return Object.assign(proto, Object.assign(Object.assign({ $typeName: "google.protobuf.FieldOptions" }, init), { targets: (_a = init.targets) !== null && _a !== undefined ? _a : [], editionDefaults: (_c = (_b = init.editionDefaults) === null || _b === undefined ? undefined : _b.map((e) => Object.assign({ $typeName: "google.protobuf.FieldOptions.EditionDefault" }, e))) !== null && _c !== undefined ? _c : [], uninterpretedOption: [] }));
}
function bootEnumDescriptorProto(init) {
  const proto = Object.create({
    visibility: 0
  });
  return Object.assign(proto, {
    $typeName: "google.protobuf.EnumDescriptorProto",
    name: init.name,
    reservedName: [],
    reservedRange: [],
    value: init.value.map((e) => Object.assign({ $typeName: "google.protobuf.EnumValueDescriptorProto" }, e))
  });
}
function messageDesc(file, path, ...paths) {
  return paths.reduce((acc, cur) => acc.nestedMessages[cur], file.messages[path]);
}
var file_google_protobuf_descriptor = /* @__PURE__ */ boot({ name: "google/protobuf/descriptor.proto", package: "google.protobuf", messageType: [{ name: "FileDescriptorSet", field: [{ name: "file", number: 1, type: 11, label: 3, typeName: ".google.protobuf.FileDescriptorProto" }], extensionRange: [{ start: 536000000, end: 536000001 }] }, { name: "FileDescriptorProto", field: [{ name: "name", number: 1, type: 9, label: 1 }, { name: "package", number: 2, type: 9, label: 1 }, { name: "dependency", number: 3, type: 9, label: 3 }, { name: "public_dependency", number: 10, type: 5, label: 3 }, { name: "weak_dependency", number: 11, type: 5, label: 3 }, { name: "option_dependency", number: 15, type: 9, label: 3 }, { name: "message_type", number: 4, type: 11, label: 3, typeName: ".google.protobuf.DescriptorProto" }, { name: "enum_type", number: 5, type: 11, label: 3, typeName: ".google.protobuf.EnumDescriptorProto" }, { name: "service", number: 6, type: 11, label: 3, typeName: ".google.protobuf.ServiceDescriptorProto" }, { name: "extension", number: 7, type: 11, label: 3, typeName: ".google.protobuf.FieldDescriptorProto" }, { name: "options", number: 8, type: 11, label: 1, typeName: ".google.protobuf.FileOptions" }, { name: "source_code_info", number: 9, type: 11, label: 1, typeName: ".google.protobuf.SourceCodeInfo" }, { name: "syntax", number: 12, type: 9, label: 1 }, { name: "edition", number: 14, type: 14, label: 1, typeName: ".google.protobuf.Edition" }] }, { name: "DescriptorProto", field: [{ name: "name", number: 1, type: 9, label: 1 }, { name: "field", number: 2, type: 11, label: 3, typeName: ".google.protobuf.FieldDescriptorProto" }, { name: "extension", number: 6, type: 11, label: 3, typeName: ".google.protobuf.FieldDescriptorProto" }, { name: "nested_type", number: 3, type: 11, label: 3, typeName: ".google.protobuf.DescriptorProto" }, { name: "enum_type", number: 4, type: 11, label: 3, typeName: ".google.protobuf.EnumDescriptorProto" }, { name: "extension_range", number: 5, type: 11, label: 3, typeName: ".google.protobuf.DescriptorProto.ExtensionRange" }, { name: "oneof_decl", number: 8, type: 11, label: 3, typeName: ".google.protobuf.OneofDescriptorProto" }, { name: "options", number: 7, type: 11, label: 1, typeName: ".google.protobuf.MessageOptions" }, { name: "reserved_range", number: 9, type: 11, label: 3, typeName: ".google.protobuf.DescriptorProto.ReservedRange" }, { name: "reserved_name", number: 10, type: 9, label: 3 }, { name: "visibility", number: 11, type: 14, label: 1, typeName: ".google.protobuf.SymbolVisibility" }], nestedType: [{ name: "ExtensionRange", field: [{ name: "start", number: 1, type: 5, label: 1 }, { name: "end", number: 2, type: 5, label: 1 }, { name: "options", number: 3, type: 11, label: 1, typeName: ".google.protobuf.ExtensionRangeOptions" }] }, { name: "ReservedRange", field: [{ name: "start", number: 1, type: 5, label: 1 }, { name: "end", number: 2, type: 5, label: 1 }] }] }, { name: "ExtensionRangeOptions", field: [{ name: "uninterpreted_option", number: 999, type: 11, label: 3, typeName: ".google.protobuf.UninterpretedOption" }, { name: "declaration", number: 2, type: 11, label: 3, typeName: ".google.protobuf.ExtensionRangeOptions.Declaration", options: { retention: 2 } }, { name: "features", number: 50, type: 11, label: 1, typeName: ".google.protobuf.FeatureSet" }, { name: "verification", number: 3, type: 14, label: 1, typeName: ".google.protobuf.ExtensionRangeOptions.VerificationState", defaultValue: "UNVERIFIED", options: { retention: 2 } }], nestedType: [{ name: "Declaration", field: [{ name: "number", number: 1, type: 5, label: 1 }, { name: "full_name", number: 2, type: 9, label: 1 }, { name: "type", number: 3, type: 9, label: 1 }, { name: "reserved", number: 5, type: 8, label: 1 }, { name: "repeated", number: 6, type: 8, label: 1 }] }], enumType: [{ name: "VerificationState", value: [{ name: "DECLARATION", number: 0 }, { name: "UNVERIFIED", number: 1 }] }], extensionRange: [{ start: 1000, end: 536870912 }] }, { name: "FieldDescriptorProto", field: [{ name: "name", number: 1, type: 9, label: 1 }, { name: "number", number: 3, type: 5, label: 1 }, { name: "label", number: 4, type: 14, label: 1, typeName: ".google.protobuf.FieldDescriptorProto.Label" }, { name: "type", number: 5, type: 14, label: 1, typeName: ".google.protobuf.FieldDescriptorProto.Type" }, { name: "type_name", number: 6, type: 9, label: 1 }, { name: "extendee", number: 2, type: 9, label: 1 }, { name: "default_value", number: 7, type: 9, label: 1 }, { name: "oneof_index", number: 9, type: 5, label: 1 }, { name: "json_name", number: 10, type: 9, label: 1 }, { name: "options", number: 8, type: 11, label: 1, typeName: ".google.protobuf.FieldOptions" }, { name: "proto3_optional", number: 17, type: 8, label: 1 }], enumType: [{ name: "Type", value: [{ name: "TYPE_DOUBLE", number: 1 }, { name: "TYPE_FLOAT", number: 2 }, { name: "TYPE_INT64", number: 3 }, { name: "TYPE_UINT64", number: 4 }, { name: "TYPE_INT32", number: 5 }, { name: "TYPE_FIXED64", number: 6 }, { name: "TYPE_FIXED32", number: 7 }, { name: "TYPE_BOOL", number: 8 }, { name: "TYPE_STRING", number: 9 }, { name: "TYPE_GROUP", number: 10 }, { name: "TYPE_MESSAGE", number: 11 }, { name: "TYPE_BYTES", number: 12 }, { name: "TYPE_UINT32", number: 13 }, { name: "TYPE_ENUM", number: 14 }, { name: "TYPE_SFIXED32", number: 15 }, { name: "TYPE_SFIXED64", number: 16 }, { name: "TYPE_SINT32", number: 17 }, { name: "TYPE_SINT64", number: 18 }] }, { name: "Label", value: [{ name: "LABEL_OPTIONAL", number: 1 }, { name: "LABEL_REPEATED", number: 3 }, { name: "LABEL_REQUIRED", number: 2 }] }] }, { name: "OneofDescriptorProto", field: [{ name: "name", number: 1, type: 9, label: 1 }, { name: "options", number: 2, type: 11, label: 1, typeName: ".google.protobuf.OneofOptions" }] }, { name: "EnumDescriptorProto", field: [{ name: "name", number: 1, type: 9, label: 1 }, { name: "value", number: 2, type: 11, label: 3, typeName: ".google.protobuf.EnumValueDescriptorProto" }, { name: "options", number: 3, type: 11, label: 1, typeName: ".google.protobuf.EnumOptions" }, { name: "reserved_range", number: 4, type: 11, label: 3, typeName: ".google.protobuf.EnumDescriptorProto.EnumReservedRange" }, { name: "reserved_name", number: 5, type: 9, label: 3 }, { name: "visibility", number: 6, type: 14, label: 1, typeName: ".google.protobuf.SymbolVisibility" }], nestedType: [{ name: "EnumReservedRange", field: [{ name: "start", number: 1, type: 5, label: 1 }, { name: "end", number: 2, type: 5, label: 1 }] }] }, { name: "EnumValueDescriptorProto", field: [{ name: "name", number: 1, type: 9, label: 1 }, { name: "number", number: 2, type: 5, label: 1 }, { name: "options", number: 3, type: 11, label: 1, typeName: ".google.protobuf.EnumValueOptions" }] }, { name: "ServiceDescriptorProto", field: [{ name: "name", number: 1, type: 9, label: 1 }, { name: "method", number: 2, type: 11, label: 3, typeName: ".google.protobuf.MethodDescriptorProto" }, { name: "options", number: 3, type: 11, label: 1, typeName: ".google.protobuf.ServiceOptions" }] }, { name: "MethodDescriptorProto", field: [{ name: "name", number: 1, type: 9, label: 1 }, { name: "input_type", number: 2, type: 9, label: 1 }, { name: "output_type", number: 3, type: 9, label: 1 }, { name: "options", number: 4, type: 11, label: 1, typeName: ".google.protobuf.MethodOptions" }, { name: "client_streaming", number: 5, type: 8, label: 1, defaultValue: "false" }, { name: "server_streaming", number: 6, type: 8, label: 1, defaultValue: "false" }] }, { name: "FileOptions", field: [{ name: "java_package", number: 1, type: 9, label: 1 }, { name: "java_outer_classname", number: 8, type: 9, label: 1 }, { name: "java_multiple_files", number: 10, type: 8, label: 1, defaultValue: "false" }, { name: "java_generate_equals_and_hash", number: 20, type: 8, label: 1, options: { deprecated: true } }, { name: "java_string_check_utf8", number: 27, type: 8, label: 1, defaultValue: "false" }, { name: "optimize_for", number: 9, type: 14, label: 1, typeName: ".google.protobuf.FileOptions.OptimizeMode", defaultValue: "SPEED" }, { name: "go_package", number: 11, type: 9, label: 1 }, { name: "cc_generic_services", number: 16, type: 8, label: 1, defaultValue: "false" }, { name: "java_generic_services", number: 17, type: 8, label: 1, defaultValue: "false" }, { name: "py_generic_services", number: 18, type: 8, label: 1, defaultValue: "false" }, { name: "deprecated", number: 23, type: 8, label: 1, defaultValue: "false" }, { name: "cc_enable_arenas", number: 31, type: 8, label: 1, defaultValue: "true" }, { name: "objc_class_prefix", number: 36, type: 9, label: 1 }, { name: "csharp_namespace", number: 37, type: 9, label: 1 }, { name: "swift_prefix", number: 39, type: 9, label: 1 }, { name: "php_class_prefix", number: 40, type: 9, label: 1 }, { name: "php_namespace", number: 41, type: 9, label: 1 }, { name: "php_metadata_namespace", number: 44, type: 9, label: 1 }, { name: "ruby_package", number: 45, type: 9, label: 1 }, { name: "features", number: 50, type: 11, label: 1, typeName: ".google.protobuf.FeatureSet" }, { name: "uninterpreted_option", number: 999, type: 11, label: 3, typeName: ".google.protobuf.UninterpretedOption" }], enumType: [{ name: "OptimizeMode", value: [{ name: "SPEED", number: 1 }, { name: "CODE_SIZE", number: 2 }, { name: "LITE_RUNTIME", number: 3 }] }], extensionRange: [{ start: 1000, end: 536870912 }] }, { name: "MessageOptions", field: [{ name: "message_set_wire_format", number: 1, type: 8, label: 1, defaultValue: "false" }, { name: "no_standard_descriptor_accessor", number: 2, type: 8, label: 1, defaultValue: "false" }, { name: "deprecated", number: 3, type: 8, label: 1, defaultValue: "false" }, { name: "map_entry", number: 7, type: 8, label: 1 }, { name: "deprecated_legacy_json_field_conflicts", number: 11, type: 8, label: 1, options: { deprecated: true } }, { name: "features", number: 12, type: 11, label: 1, typeName: ".google.protobuf.FeatureSet" }, { name: "uninterpreted_option", number: 999, type: 11, label: 3, typeName: ".google.protobuf.UninterpretedOption" }], extensionRange: [{ start: 1000, end: 536870912 }] }, { name: "FieldOptions", field: [{ name: "ctype", number: 1, type: 14, label: 1, typeName: ".google.protobuf.FieldOptions.CType", defaultValue: "STRING" }, { name: "packed", number: 2, type: 8, label: 1 }, { name: "jstype", number: 6, type: 14, label: 1, typeName: ".google.protobuf.FieldOptions.JSType", defaultValue: "JS_NORMAL" }, { name: "lazy", number: 5, type: 8, label: 1, defaultValue: "false" }, { name: "unverified_lazy", number: 15, type: 8, label: 1, defaultValue: "false" }, { name: "deprecated", number: 3, type: 8, label: 1, defaultValue: "false" }, { name: "weak", number: 10, type: 8, label: 1, defaultValue: "false" }, { name: "debug_redact", number: 16, type: 8, label: 1, defaultValue: "false" }, { name: "retention", number: 17, type: 14, label: 1, typeName: ".google.protobuf.FieldOptions.OptionRetention" }, { name: "targets", number: 19, type: 14, label: 3, typeName: ".google.protobuf.FieldOptions.OptionTargetType" }, { name: "edition_defaults", number: 20, type: 11, label: 3, typeName: ".google.protobuf.FieldOptions.EditionDefault" }, { name: "features", number: 21, type: 11, label: 1, typeName: ".google.protobuf.FeatureSet" }, { name: "feature_support", number: 22, type: 11, label: 1, typeName: ".google.protobuf.FieldOptions.FeatureSupport" }, { name: "uninterpreted_option", number: 999, type: 11, label: 3, typeName: ".google.protobuf.UninterpretedOption" }], nestedType: [{ name: "EditionDefault", field: [{ name: "edition", number: 3, type: 14, label: 1, typeName: ".google.protobuf.Edition" }, { name: "value", number: 2, type: 9, label: 1 }] }, { name: "FeatureSupport", field: [{ name: "edition_introduced", number: 1, type: 14, label: 1, typeName: ".google.protobuf.Edition" }, { name: "edition_deprecated", number: 2, type: 14, label: 1, typeName: ".google.protobuf.Edition" }, { name: "deprecation_warning", number: 3, type: 9, label: 1 }, { name: "edition_removed", number: 4, type: 14, label: 1, typeName: ".google.protobuf.Edition" }] }], enumType: [{ name: "CType", value: [{ name: "STRING", number: 0 }, { name: "CORD", number: 1 }, { name: "STRING_PIECE", number: 2 }] }, { name: "JSType", value: [{ name: "JS_NORMAL", number: 0 }, { name: "JS_STRING", number: 1 }, { name: "JS_NUMBER", number: 2 }] }, { name: "OptionRetention", value: [{ name: "RETENTION_UNKNOWN", number: 0 }, { name: "RETENTION_RUNTIME", number: 1 }, { name: "RETENTION_SOURCE", number: 2 }] }, { name: "OptionTargetType", value: [{ name: "TARGET_TYPE_UNKNOWN", number: 0 }, { name: "TARGET_TYPE_FILE", number: 1 }, { name: "TARGET_TYPE_EXTENSION_RANGE", number: 2 }, { name: "TARGET_TYPE_MESSAGE", number: 3 }, { name: "TARGET_TYPE_FIELD", number: 4 }, { name: "TARGET_TYPE_ONEOF", number: 5 }, { name: "TARGET_TYPE_ENUM", number: 6 }, { name: "TARGET_TYPE_ENUM_ENTRY", number: 7 }, { name: "TARGET_TYPE_SERVICE", number: 8 }, { name: "TARGET_TYPE_METHOD", number: 9 }] }], extensionRange: [{ start: 1000, end: 536870912 }] }, { name: "OneofOptions", field: [{ name: "features", number: 1, type: 11, label: 1, typeName: ".google.protobuf.FeatureSet" }, { name: "uninterpreted_option", number: 999, type: 11, label: 3, typeName: ".google.protobuf.UninterpretedOption" }], extensionRange: [{ start: 1000, end: 536870912 }] }, { name: "EnumOptions", field: [{ name: "allow_alias", number: 2, type: 8, label: 1 }, { name: "deprecated", number: 3, type: 8, label: 1, defaultValue: "false" }, { name: "deprecated_legacy_json_field_conflicts", number: 6, type: 8, label: 1, options: { deprecated: true } }, { name: "features", number: 7, type: 11, label: 1, typeName: ".google.protobuf.FeatureSet" }, { name: "uninterpreted_option", number: 999, type: 11, label: 3, typeName: ".google.protobuf.UninterpretedOption" }], extensionRange: [{ start: 1000, end: 536870912 }] }, { name: "EnumValueOptions", field: [{ name: "deprecated", number: 1, type: 8, label: 1, defaultValue: "false" }, { name: "features", number: 2, type: 11, label: 1, typeName: ".google.protobuf.FeatureSet" }, { name: "debug_redact", number: 3, type: 8, label: 1, defaultValue: "false" }, { name: "feature_support", number: 4, type: 11, label: 1, typeName: ".google.protobuf.FieldOptions.FeatureSupport" }, { name: "uninterpreted_option", number: 999, type: 11, label: 3, typeName: ".google.protobuf.UninterpretedOption" }], extensionRange: [{ start: 1000, end: 536870912 }] }, { name: "ServiceOptions", field: [{ name: "features", number: 34, type: 11, label: 1, typeName: ".google.protobuf.FeatureSet" }, { name: "deprecated", number: 33, type: 8, label: 1, defaultValue: "false" }, { name: "uninterpreted_option", number: 999, type: 11, label: 3, typeName: ".google.protobuf.UninterpretedOption" }], extensionRange: [{ start: 1000, end: 536870912 }] }, { name: "MethodOptions", field: [{ name: "deprecated", number: 33, type: 8, label: 1, defaultValue: "false" }, { name: "idempotency_level", number: 34, type: 14, label: 1, typeName: ".google.protobuf.MethodOptions.IdempotencyLevel", defaultValue: "IDEMPOTENCY_UNKNOWN" }, { name: "features", number: 35, type: 11, label: 1, typeName: ".google.protobuf.FeatureSet" }, { name: "uninterpreted_option", number: 999, type: 11, label: 3, typeName: ".google.protobuf.UninterpretedOption" }], enumType: [{ name: "IdempotencyLevel", value: [{ name: "IDEMPOTENCY_UNKNOWN", number: 0 }, { name: "NO_SIDE_EFFECTS", number: 1 }, { name: "IDEMPOTENT", number: 2 }] }], extensionRange: [{ start: 1000, end: 536870912 }] }, { name: "UninterpretedOption", field: [{ name: "name", number: 2, type: 11, label: 3, typeName: ".google.protobuf.UninterpretedOption.NamePart" }, { name: "identifier_value", number: 3, type: 9, label: 1 }, { name: "positive_int_value", number: 4, type: 4, label: 1 }, { name: "negative_int_value", number: 5, type: 3, label: 1 }, { name: "double_value", number: 6, type: 1, label: 1 }, { name: "string_value", number: 7, type: 12, label: 1 }, { name: "aggregate_value", number: 8, type: 9, label: 1 }], nestedType: [{ name: "NamePart", field: [{ name: "name_part", number: 1, type: 9, label: 2 }, { name: "is_extension", number: 2, type: 8, label: 2 }] }] }, { name: "FeatureSet", field: [{ name: "field_presence", number: 1, type: 14, label: 1, typeName: ".google.protobuf.FeatureSet.FieldPresence", options: { retention: 1, targets: [4, 1], editionDefaults: [{ value: "EXPLICIT", edition: 900 }, { value: "IMPLICIT", edition: 999 }, { value: "EXPLICIT", edition: 1000 }] } }, { name: "enum_type", number: 2, type: 14, label: 1, typeName: ".google.protobuf.FeatureSet.EnumType", options: { retention: 1, targets: [6, 1], editionDefaults: [{ value: "CLOSED", edition: 900 }, { value: "OPEN", edition: 999 }] } }, { name: "repeated_field_encoding", number: 3, type: 14, label: 1, typeName: ".google.protobuf.FeatureSet.RepeatedFieldEncoding", options: { retention: 1, targets: [4, 1], editionDefaults: [{ value: "EXPANDED", edition: 900 }, { value: "PACKED", edition: 999 }] } }, { name: "utf8_validation", number: 4, type: 14, label: 1, typeName: ".google.protobuf.FeatureSet.Utf8Validation", options: { retention: 1, targets: [4, 1], editionDefaults: [{ value: "NONE", edition: 900 }, { value: "VERIFY", edition: 999 }] } }, { name: "message_encoding", number: 5, type: 14, label: 1, typeName: ".google.protobuf.FeatureSet.MessageEncoding", options: { retention: 1, targets: [4, 1], editionDefaults: [{ value: "LENGTH_PREFIXED", edition: 900 }] } }, { name: "json_format", number: 6, type: 14, label: 1, typeName: ".google.protobuf.FeatureSet.JsonFormat", options: { retention: 1, targets: [3, 6, 1], editionDefaults: [{ value: "LEGACY_BEST_EFFORT", edition: 900 }, { value: "ALLOW", edition: 999 }] } }, { name: "enforce_naming_style", number: 7, type: 14, label: 1, typeName: ".google.protobuf.FeatureSet.EnforceNamingStyle", options: { retention: 2, targets: [1, 2, 3, 4, 5, 6, 7, 8, 9], editionDefaults: [{ value: "STYLE_LEGACY", edition: 900 }, { value: "STYLE2024", edition: 1001 }] } }, { name: "default_symbol_visibility", number: 8, type: 14, label: 1, typeName: ".google.protobuf.FeatureSet.VisibilityFeature.DefaultSymbolVisibility", options: { retention: 2, targets: [1], editionDefaults: [{ value: "EXPORT_ALL", edition: 900 }, { value: "EXPORT_TOP_LEVEL", edition: 1001 }] } }], nestedType: [{ name: "VisibilityFeature", enumType: [{ name: "DefaultSymbolVisibility", value: [{ name: "DEFAULT_SYMBOL_VISIBILITY_UNKNOWN", number: 0 }, { name: "EXPORT_ALL", number: 1 }, { name: "EXPORT_TOP_LEVEL", number: 2 }, { name: "LOCAL_ALL", number: 3 }, { name: "STRICT", number: 4 }] }] }], enumType: [{ name: "FieldPresence", value: [{ name: "FIELD_PRESENCE_UNKNOWN", number: 0 }, { name: "EXPLICIT", number: 1 }, { name: "IMPLICIT", number: 2 }, { name: "LEGACY_REQUIRED", number: 3 }] }, { name: "EnumType", value: [{ name: "ENUM_TYPE_UNKNOWN", number: 0 }, { name: "OPEN", number: 1 }, { name: "CLOSED", number: 2 }] }, { name: "RepeatedFieldEncoding", value: [{ name: "REPEATED_FIELD_ENCODING_UNKNOWN", number: 0 }, { name: "PACKED", number: 1 }, { name: "EXPANDED", number: 2 }] }, { name: "Utf8Validation", value: [{ name: "UTF8_VALIDATION_UNKNOWN", number: 0 }, { name: "VERIFY", number: 2 }, { name: "NONE", number: 3 }] }, { name: "MessageEncoding", value: [{ name: "MESSAGE_ENCODING_UNKNOWN", number: 0 }, { name: "LENGTH_PREFIXED", number: 1 }, { name: "DELIMITED", number: 2 }] }, { name: "JsonFormat", value: [{ name: "JSON_FORMAT_UNKNOWN", number: 0 }, { name: "ALLOW", number: 1 }, { name: "LEGACY_BEST_EFFORT", number: 2 }] }, { name: "EnforceNamingStyle", value: [{ name: "ENFORCE_NAMING_STYLE_UNKNOWN", number: 0 }, { name: "STYLE2024", number: 1 }, { name: "STYLE_LEGACY", number: 2 }] }], extensionRange: [{ start: 1000, end: 9995 }, { start: 9995, end: 1e4 }, { start: 1e4, end: 10001 }] }, { name: "FeatureSetDefaults", field: [{ name: "defaults", number: 1, type: 11, label: 3, typeName: ".google.protobuf.FeatureSetDefaults.FeatureSetEditionDefault" }, { name: "minimum_edition", number: 4, type: 14, label: 1, typeName: ".google.protobuf.Edition" }, { name: "maximum_edition", number: 5, type: 14, label: 1, typeName: ".google.protobuf.Edition" }], nestedType: [{ name: "FeatureSetEditionDefault", field: [{ name: "edition", number: 3, type: 14, label: 1, typeName: ".google.protobuf.Edition" }, { name: "overridable_features", number: 4, type: 11, label: 1, typeName: ".google.protobuf.FeatureSet" }, { name: "fixed_features", number: 5, type: 11, label: 1, typeName: ".google.protobuf.FeatureSet" }] }] }, { name: "SourceCodeInfo", field: [{ name: "location", number: 1, type: 11, label: 3, typeName: ".google.protobuf.SourceCodeInfo.Location" }], nestedType: [{ name: "Location", field: [{ name: "path", number: 1, type: 5, label: 3, options: { packed: true } }, { name: "span", number: 2, type: 5, label: 3, options: { packed: true } }, { name: "leading_comments", number: 3, type: 9, label: 1 }, { name: "trailing_comments", number: 4, type: 9, label: 1 }, { name: "leading_detached_comments", number: 6, type: 9, label: 3 }] }], extensionRange: [{ start: 536000000, end: 536000001 }] }, { name: "GeneratedCodeInfo", field: [{ name: "annotation", number: 1, type: 11, label: 3, typeName: ".google.protobuf.GeneratedCodeInfo.Annotation" }], nestedType: [{ name: "Annotation", field: [{ name: "path", number: 1, type: 5, label: 3, options: { packed: true } }, { name: "source_file", number: 2, type: 9, label: 1 }, { name: "begin", number: 3, type: 5, label: 1 }, { name: "end", number: 4, type: 5, label: 1 }, { name: "semantic", number: 5, type: 14, label: 1, typeName: ".google.protobuf.GeneratedCodeInfo.Annotation.Semantic" }], enumType: [{ name: "Semantic", value: [{ name: "NONE", number: 0 }, { name: "SET", number: 1 }, { name: "ALIAS", number: 2 }] }] }] }], enumType: [{ name: "Edition", value: [{ name: "EDITION_UNKNOWN", number: 0 }, { name: "EDITION_LEGACY", number: 900 }, { name: "EDITION_PROTO2", number: 998 }, { name: "EDITION_PROTO3", number: 999 }, { name: "EDITION_2023", number: 1000 }, { name: "EDITION_2024", number: 1001 }, { name: "EDITION_1_TEST_ONLY", number: 1 }, { name: "EDITION_2_TEST_ONLY", number: 2 }, { name: "EDITION_99997_TEST_ONLY", number: 99997 }, { name: "EDITION_99998_TEST_ONLY", number: 99998 }, { name: "EDITION_99999_TEST_ONLY", number: 99999 }, { name: "EDITION_MAX", number: 2147483647 }] }, { name: "SymbolVisibility", value: [{ name: "VISIBILITY_UNSET", number: 0 }, { name: "VISIBILITY_LOCAL", number: 1 }, { name: "VISIBILITY_EXPORT", number: 2 }] }] });
var FileDescriptorProtoSchema = /* @__PURE__ */ messageDesc(file_google_protobuf_descriptor, 1);
var ExtensionRangeOptions_VerificationState;
(function(ExtensionRangeOptions_VerificationState2) {
  ExtensionRangeOptions_VerificationState2[ExtensionRangeOptions_VerificationState2["DECLARATION"] = 0] = "DECLARATION";
  ExtensionRangeOptions_VerificationState2[ExtensionRangeOptions_VerificationState2["UNVERIFIED"] = 1] = "UNVERIFIED";
})(ExtensionRangeOptions_VerificationState || (ExtensionRangeOptions_VerificationState = {}));
var FieldDescriptorProto_Type;
(function(FieldDescriptorProto_Type2) {
  FieldDescriptorProto_Type2[FieldDescriptorProto_Type2["DOUBLE"] = 1] = "DOUBLE";
  FieldDescriptorProto_Type2[FieldDescriptorProto_Type2["FLOAT"] = 2] = "FLOAT";
  FieldDescriptorProto_Type2[FieldDescriptorProto_Type2["INT64"] = 3] = "INT64";
  FieldDescriptorProto_Type2[FieldDescriptorProto_Type2["UINT64"] = 4] = "UINT64";
  FieldDescriptorProto_Type2[FieldDescriptorProto_Type2["INT32"] = 5] = "INT32";
  FieldDescriptorProto_Type2[FieldDescriptorProto_Type2["FIXED64"] = 6] = "FIXED64";
  FieldDescriptorProto_Type2[FieldDescriptorProto_Type2["FIXED32"] = 7] = "FIXED32";
  FieldDescriptorProto_Type2[FieldDescriptorProto_Type2["BOOL"] = 8] = "BOOL";
  FieldDescriptorProto_Type2[FieldDescriptorProto_Type2["STRING"] = 9] = "STRING";
  FieldDescriptorProto_Type2[FieldDescriptorProto_Type2["GROUP"] = 10] = "GROUP";
  FieldDescriptorProto_Type2[FieldDescriptorProto_Type2["MESSAGE"] = 11] = "MESSAGE";
  FieldDescriptorProto_Type2[FieldDescriptorProto_Type2["BYTES"] = 12] = "BYTES";
  FieldDescriptorProto_Type2[FieldDescriptorProto_Type2["UINT32"] = 13] = "UINT32";
  FieldDescriptorProto_Type2[FieldDescriptorProto_Type2["ENUM"] = 14] = "ENUM";
  FieldDescriptorProto_Type2[FieldDescriptorProto_Type2["SFIXED32"] = 15] = "SFIXED32";
  FieldDescriptorProto_Type2[FieldDescriptorProto_Type2["SFIXED64"] = 16] = "SFIXED64";
  FieldDescriptorProto_Type2[FieldDescriptorProto_Type2["SINT32"] = 17] = "SINT32";
  FieldDescriptorProto_Type2[FieldDescriptorProto_Type2["SINT64"] = 18] = "SINT64";
})(FieldDescriptorProto_Type || (FieldDescriptorProto_Type = {}));
var FieldDescriptorProto_Label;
(function(FieldDescriptorProto_Label2) {
  FieldDescriptorProto_Label2[FieldDescriptorProto_Label2["OPTIONAL"] = 1] = "OPTIONAL";
  FieldDescriptorProto_Label2[FieldDescriptorProto_Label2["REPEATED"] = 3] = "REPEATED";
  FieldDescriptorProto_Label2[FieldDescriptorProto_Label2["REQUIRED"] = 2] = "REQUIRED";
})(FieldDescriptorProto_Label || (FieldDescriptorProto_Label = {}));
var FileOptions_OptimizeMode;
(function(FileOptions_OptimizeMode2) {
  FileOptions_OptimizeMode2[FileOptions_OptimizeMode2["SPEED"] = 1] = "SPEED";
  FileOptions_OptimizeMode2[FileOptions_OptimizeMode2["CODE_SIZE"] = 2] = "CODE_SIZE";
  FileOptions_OptimizeMode2[FileOptions_OptimizeMode2["LITE_RUNTIME"] = 3] = "LITE_RUNTIME";
})(FileOptions_OptimizeMode || (FileOptions_OptimizeMode = {}));
var FieldOptions_CType;
(function(FieldOptions_CType2) {
  FieldOptions_CType2[FieldOptions_CType2["STRING"] = 0] = "STRING";
  FieldOptions_CType2[FieldOptions_CType2["CORD"] = 1] = "CORD";
  FieldOptions_CType2[FieldOptions_CType2["STRING_PIECE"] = 2] = "STRING_PIECE";
})(FieldOptions_CType || (FieldOptions_CType = {}));
var FieldOptions_JSType;
(function(FieldOptions_JSType2) {
  FieldOptions_JSType2[FieldOptions_JSType2["JS_NORMAL"] = 0] = "JS_NORMAL";
  FieldOptions_JSType2[FieldOptions_JSType2["JS_STRING"] = 1] = "JS_STRING";
  FieldOptions_JSType2[FieldOptions_JSType2["JS_NUMBER"] = 2] = "JS_NUMBER";
})(FieldOptions_JSType || (FieldOptions_JSType = {}));
var FieldOptions_OptionRetention;
(function(FieldOptions_OptionRetention2) {
  FieldOptions_OptionRetention2[FieldOptions_OptionRetention2["RETENTION_UNKNOWN"] = 0] = "RETENTION_UNKNOWN";
  FieldOptions_OptionRetention2[FieldOptions_OptionRetention2["RETENTION_RUNTIME"] = 1] = "RETENTION_RUNTIME";
  FieldOptions_OptionRetention2[FieldOptions_OptionRetention2["RETENTION_SOURCE"] = 2] = "RETENTION_SOURCE";
})(FieldOptions_OptionRetention || (FieldOptions_OptionRetention = {}));
var FieldOptions_OptionTargetType;
(function(FieldOptions_OptionTargetType2) {
  FieldOptions_OptionTargetType2[FieldOptions_OptionTargetType2["TARGET_TYPE_UNKNOWN"] = 0] = "TARGET_TYPE_UNKNOWN";
  FieldOptions_OptionTargetType2[FieldOptions_OptionTargetType2["TARGET_TYPE_FILE"] = 1] = "TARGET_TYPE_FILE";
  FieldOptions_OptionTargetType2[FieldOptions_OptionTargetType2["TARGET_TYPE_EXTENSION_RANGE"] = 2] = "TARGET_TYPE_EXTENSION_RANGE";
  FieldOptions_OptionTargetType2[FieldOptions_OptionTargetType2["TARGET_TYPE_MESSAGE"] = 3] = "TARGET_TYPE_MESSAGE";
  FieldOptions_OptionTargetType2[FieldOptions_OptionTargetType2["TARGET_TYPE_FIELD"] = 4] = "TARGET_TYPE_FIELD";
  FieldOptions_OptionTargetType2[FieldOptions_OptionTargetType2["TARGET_TYPE_ONEOF"] = 5] = "TARGET_TYPE_ONEOF";
  FieldOptions_OptionTargetType2[FieldOptions_OptionTargetType2["TARGET_TYPE_ENUM"] = 6] = "TARGET_TYPE_ENUM";
  FieldOptions_OptionTargetType2[FieldOptions_OptionTargetType2["TARGET_TYPE_ENUM_ENTRY"] = 7] = "TARGET_TYPE_ENUM_ENTRY";
  FieldOptions_OptionTargetType2[FieldOptions_OptionTargetType2["TARGET_TYPE_SERVICE"] = 8] = "TARGET_TYPE_SERVICE";
  FieldOptions_OptionTargetType2[FieldOptions_OptionTargetType2["TARGET_TYPE_METHOD"] = 9] = "TARGET_TYPE_METHOD";
})(FieldOptions_OptionTargetType || (FieldOptions_OptionTargetType = {}));
var MethodOptions_IdempotencyLevel;
(function(MethodOptions_IdempotencyLevel2) {
  MethodOptions_IdempotencyLevel2[MethodOptions_IdempotencyLevel2["IDEMPOTENCY_UNKNOWN"] = 0] = "IDEMPOTENCY_UNKNOWN";
  MethodOptions_IdempotencyLevel2[MethodOptions_IdempotencyLevel2["NO_SIDE_EFFECTS"] = 1] = "NO_SIDE_EFFECTS";
  MethodOptions_IdempotencyLevel2[MethodOptions_IdempotencyLevel2["IDEMPOTENT"] = 2] = "IDEMPOTENT";
})(MethodOptions_IdempotencyLevel || (MethodOptions_IdempotencyLevel = {}));
var FeatureSet_VisibilityFeature_DefaultSymbolVisibility;
(function(FeatureSet_VisibilityFeature_DefaultSymbolVisibility2) {
  FeatureSet_VisibilityFeature_DefaultSymbolVisibility2[FeatureSet_VisibilityFeature_DefaultSymbolVisibility2["DEFAULT_SYMBOL_VISIBILITY_UNKNOWN"] = 0] = "DEFAULT_SYMBOL_VISIBILITY_UNKNOWN";
  FeatureSet_VisibilityFeature_DefaultSymbolVisibility2[FeatureSet_VisibilityFeature_DefaultSymbolVisibility2["EXPORT_ALL"] = 1] = "EXPORT_ALL";
  FeatureSet_VisibilityFeature_DefaultSymbolVisibility2[FeatureSet_VisibilityFeature_DefaultSymbolVisibility2["EXPORT_TOP_LEVEL"] = 2] = "EXPORT_TOP_LEVEL";
  FeatureSet_VisibilityFeature_DefaultSymbolVisibility2[FeatureSet_VisibilityFeature_DefaultSymbolVisibility2["LOCAL_ALL"] = 3] = "LOCAL_ALL";
  FeatureSet_VisibilityFeature_DefaultSymbolVisibility2[FeatureSet_VisibilityFeature_DefaultSymbolVisibility2["STRICT"] = 4] = "STRICT";
})(FeatureSet_VisibilityFeature_DefaultSymbolVisibility || (FeatureSet_VisibilityFeature_DefaultSymbolVisibility = {}));
var FeatureSet_FieldPresence;
(function(FeatureSet_FieldPresence2) {
  FeatureSet_FieldPresence2[FeatureSet_FieldPresence2["FIELD_PRESENCE_UNKNOWN"] = 0] = "FIELD_PRESENCE_UNKNOWN";
  FeatureSet_FieldPresence2[FeatureSet_FieldPresence2["EXPLICIT"] = 1] = "EXPLICIT";
  FeatureSet_FieldPresence2[FeatureSet_FieldPresence2["IMPLICIT"] = 2] = "IMPLICIT";
  FeatureSet_FieldPresence2[FeatureSet_FieldPresence2["LEGACY_REQUIRED"] = 3] = "LEGACY_REQUIRED";
})(FeatureSet_FieldPresence || (FeatureSet_FieldPresence = {}));
var FeatureSet_EnumType;
(function(FeatureSet_EnumType2) {
  FeatureSet_EnumType2[FeatureSet_EnumType2["ENUM_TYPE_UNKNOWN"] = 0] = "ENUM_TYPE_UNKNOWN";
  FeatureSet_EnumType2[FeatureSet_EnumType2["OPEN"] = 1] = "OPEN";
  FeatureSet_EnumType2[FeatureSet_EnumType2["CLOSED"] = 2] = "CLOSED";
})(FeatureSet_EnumType || (FeatureSet_EnumType = {}));
var FeatureSet_RepeatedFieldEncoding;
(function(FeatureSet_RepeatedFieldEncoding2) {
  FeatureSet_RepeatedFieldEncoding2[FeatureSet_RepeatedFieldEncoding2["REPEATED_FIELD_ENCODING_UNKNOWN"] = 0] = "REPEATED_FIELD_ENCODING_UNKNOWN";
  FeatureSet_RepeatedFieldEncoding2[FeatureSet_RepeatedFieldEncoding2["PACKED"] = 1] = "PACKED";
  FeatureSet_RepeatedFieldEncoding2[FeatureSet_RepeatedFieldEncoding2["EXPANDED"] = 2] = "EXPANDED";
})(FeatureSet_RepeatedFieldEncoding || (FeatureSet_RepeatedFieldEncoding = {}));
var FeatureSet_Utf8Validation;
(function(FeatureSet_Utf8Validation2) {
  FeatureSet_Utf8Validation2[FeatureSet_Utf8Validation2["UTF8_VALIDATION_UNKNOWN"] = 0] = "UTF8_VALIDATION_UNKNOWN";
  FeatureSet_Utf8Validation2[FeatureSet_Utf8Validation2["VERIFY"] = 2] = "VERIFY";
  FeatureSet_Utf8Validation2[FeatureSet_Utf8Validation2["NONE"] = 3] = "NONE";
})(FeatureSet_Utf8Validation || (FeatureSet_Utf8Validation = {}));
var FeatureSet_MessageEncoding;
(function(FeatureSet_MessageEncoding2) {
  FeatureSet_MessageEncoding2[FeatureSet_MessageEncoding2["MESSAGE_ENCODING_UNKNOWN"] = 0] = "MESSAGE_ENCODING_UNKNOWN";
  FeatureSet_MessageEncoding2[FeatureSet_MessageEncoding2["LENGTH_PREFIXED"] = 1] = "LENGTH_PREFIXED";
  FeatureSet_MessageEncoding2[FeatureSet_MessageEncoding2["DELIMITED"] = 2] = "DELIMITED";
})(FeatureSet_MessageEncoding || (FeatureSet_MessageEncoding = {}));
var FeatureSet_JsonFormat;
(function(FeatureSet_JsonFormat2) {
  FeatureSet_JsonFormat2[FeatureSet_JsonFormat2["JSON_FORMAT_UNKNOWN"] = 0] = "JSON_FORMAT_UNKNOWN";
  FeatureSet_JsonFormat2[FeatureSet_JsonFormat2["ALLOW"] = 1] = "ALLOW";
  FeatureSet_JsonFormat2[FeatureSet_JsonFormat2["LEGACY_BEST_EFFORT"] = 2] = "LEGACY_BEST_EFFORT";
})(FeatureSet_JsonFormat || (FeatureSet_JsonFormat = {}));
var FeatureSet_EnforceNamingStyle;
(function(FeatureSet_EnforceNamingStyle2) {
  FeatureSet_EnforceNamingStyle2[FeatureSet_EnforceNamingStyle2["ENFORCE_NAMING_STYLE_UNKNOWN"] = 0] = "ENFORCE_NAMING_STYLE_UNKNOWN";
  FeatureSet_EnforceNamingStyle2[FeatureSet_EnforceNamingStyle2["STYLE2024"] = 1] = "STYLE2024";
  FeatureSet_EnforceNamingStyle2[FeatureSet_EnforceNamingStyle2["STYLE_LEGACY"] = 2] = "STYLE_LEGACY";
})(FeatureSet_EnforceNamingStyle || (FeatureSet_EnforceNamingStyle = {}));
var GeneratedCodeInfo_Annotation_Semantic;
(function(GeneratedCodeInfo_Annotation_Semantic2) {
  GeneratedCodeInfo_Annotation_Semantic2[GeneratedCodeInfo_Annotation_Semantic2["NONE"] = 0] = "NONE";
  GeneratedCodeInfo_Annotation_Semantic2[GeneratedCodeInfo_Annotation_Semantic2["SET"] = 1] = "SET";
  GeneratedCodeInfo_Annotation_Semantic2[GeneratedCodeInfo_Annotation_Semantic2["ALIAS"] = 2] = "ALIAS";
})(GeneratedCodeInfo_Annotation_Semantic || (GeneratedCodeInfo_Annotation_Semantic = {}));
var Edition;
(function(Edition2) {
  Edition2[Edition2["EDITION_UNKNOWN"] = 0] = "EDITION_UNKNOWN";
  Edition2[Edition2["EDITION_LEGACY"] = 900] = "EDITION_LEGACY";
  Edition2[Edition2["EDITION_PROTO2"] = 998] = "EDITION_PROTO2";
  Edition2[Edition2["EDITION_PROTO3"] = 999] = "EDITION_PROTO3";
  Edition2[Edition2["EDITION_2023"] = 1000] = "EDITION_2023";
  Edition2[Edition2["EDITION_2024"] = 1001] = "EDITION_2024";
  Edition2[Edition2["EDITION_1_TEST_ONLY"] = 1] = "EDITION_1_TEST_ONLY";
  Edition2[Edition2["EDITION_2_TEST_ONLY"] = 2] = "EDITION_2_TEST_ONLY";
  Edition2[Edition2["EDITION_99997_TEST_ONLY"] = 99997] = "EDITION_99997_TEST_ONLY";
  Edition2[Edition2["EDITION_99998_TEST_ONLY"] = 99998] = "EDITION_99998_TEST_ONLY";
  Edition2[Edition2["EDITION_99999_TEST_ONLY"] = 99999] = "EDITION_99999_TEST_ONLY";
  Edition2[Edition2["EDITION_MAX"] = 2147483647] = "EDITION_MAX";
})(Edition || (Edition = {}));
var SymbolVisibility;
(function(SymbolVisibility2) {
  SymbolVisibility2[SymbolVisibility2["VISIBILITY_UNSET"] = 0] = "VISIBILITY_UNSET";
  SymbolVisibility2[SymbolVisibility2["VISIBILITY_LOCAL"] = 1] = "VISIBILITY_LOCAL";
  SymbolVisibility2[SymbolVisibility2["VISIBILITY_EXPORT"] = 2] = "VISIBILITY_EXPORT";
})(SymbolVisibility || (SymbolVisibility = {}));
var readDefaults = {
  readUnknownFields: true
};
function makeReadOptions(options) {
  return options ? Object.assign(Object.assign({}, readDefaults), options) : readDefaults;
}
function fromBinary(schema, bytes, options) {
  const msg = reflect(schema, undefined, false);
  readMessage(msg, new BinaryReader(bytes), makeReadOptions(options), false, bytes.byteLength);
  return msg.message;
}
function readMessage(message, reader, options, delimited, lengthOrDelimitedFieldNo) {
  var _a;
  const end = delimited ? reader.len : reader.pos + lengthOrDelimitedFieldNo;
  let fieldNo;
  let wireType;
  const unknownFields = (_a = message.getUnknown()) !== null && _a !== undefined ? _a : [];
  while (reader.pos < end) {
    [fieldNo, wireType] = reader.tag();
    if (delimited && wireType == WireType.EndGroup) {
      break;
    }
    const field = message.findNumber(fieldNo);
    if (!field) {
      const data = reader.skip(wireType, fieldNo);
      if (options.readUnknownFields) {
        unknownFields.push({ no: fieldNo, wireType, data });
      }
      continue;
    }
    readField(message, reader, field, wireType, options);
  }
  if (delimited) {
    if (wireType != WireType.EndGroup || fieldNo !== lengthOrDelimitedFieldNo) {
      throw new Error("invalid end group tag");
    }
  }
  if (unknownFields.length > 0) {
    message.setUnknown(unknownFields);
  }
}
function readField(message, reader, field, wireType, options) {
  var _a;
  switch (field.fieldKind) {
    case "scalar":
      message.set(field, readScalar(reader, field.scalar));
      break;
    case "enum":
      const val = readScalar(reader, ScalarType.INT32);
      if (field.enum.open) {
        message.set(field, val);
      } else {
        const ok = field.enum.values.some((v) => v.number === val);
        if (ok) {
          message.set(field, val);
        } else if (options.readUnknownFields) {
          const data = new BinaryWriter().int32(val).finish();
          const unknownFields = (_a = message.getUnknown()) !== null && _a !== undefined ? _a : [];
          unknownFields.push({ no: field.number, wireType, data });
          message.setUnknown(unknownFields);
        }
      }
      break;
    case "message":
      message.set(field, readMessageField(reader, options, field, message.get(field)));
      break;
    case "list":
      readListField(reader, wireType, message.get(field), options);
      break;
    case "map":
      readMapEntry(reader, message.get(field), options);
      break;
  }
}
function readMapEntry(reader, map, options) {
  const field = map.field();
  let key;
  let val;
  const len = reader.uint32();
  const end = reader.pos + len;
  while (reader.pos < end) {
    const [fieldNo] = reader.tag();
    switch (fieldNo) {
      case 1:
        key = readScalar(reader, field.mapKey);
        break;
      case 2:
        switch (field.mapKind) {
          case "scalar":
            val = readScalar(reader, field.scalar);
            break;
          case "enum":
            val = reader.int32();
            break;
          case "message":
            val = readMessageField(reader, options, field);
            break;
        }
        break;
    }
  }
  if (key === undefined) {
    key = scalarZeroValue(field.mapKey, false);
  }
  if (val === undefined) {
    switch (field.mapKind) {
      case "scalar":
        val = scalarZeroValue(field.scalar, false);
        break;
      case "enum":
        val = field.enum.values[0].number;
        break;
      case "message":
        val = reflect(field.message, undefined, false);
        break;
    }
  }
  map.set(key, val);
}
function readListField(reader, wireType, list, options) {
  var _a;
  const field = list.field();
  if (field.listKind === "message") {
    list.add(readMessageField(reader, options, field));
    return;
  }
  const scalarType = (_a = field.scalar) !== null && _a !== undefined ? _a : ScalarType.INT32;
  const packed = wireType == WireType.LengthDelimited && scalarType != ScalarType.STRING && scalarType != ScalarType.BYTES;
  if (!packed) {
    list.add(readScalar(reader, scalarType));
    return;
  }
  const e = reader.uint32() + reader.pos;
  while (reader.pos < e) {
    list.add(readScalar(reader, scalarType));
  }
}
function readMessageField(reader, options, field, mergeMessage) {
  const delimited = field.delimitedEncoding;
  const message = mergeMessage !== null && mergeMessage !== undefined ? mergeMessage : reflect(field.message, undefined, false);
  readMessage(message, reader, options, delimited, delimited ? field.number : reader.uint32());
  return message;
}
function readScalar(reader, type) {
  switch (type) {
    case ScalarType.STRING:
      return reader.string();
    case ScalarType.BOOL:
      return reader.bool();
    case ScalarType.DOUBLE:
      return reader.double();
    case ScalarType.FLOAT:
      return reader.float();
    case ScalarType.INT32:
      return reader.int32();
    case ScalarType.INT64:
      return reader.int64();
    case ScalarType.UINT64:
      return reader.uint64();
    case ScalarType.FIXED64:
      return reader.fixed64();
    case ScalarType.BYTES:
      return reader.bytes();
    case ScalarType.FIXED32:
      return reader.fixed32();
    case ScalarType.SFIXED32:
      return reader.sfixed32();
    case ScalarType.SFIXED64:
      return reader.sfixed64();
    case ScalarType.SINT64:
      return reader.sint64();
    case ScalarType.UINT32:
      return reader.uint32();
    case ScalarType.SINT32:
      return reader.sint32();
  }
}
function fileDesc(b64, imports) {
  var _a;
  const root = fromBinary(FileDescriptorProtoSchema, base64Decode(b64));
  root.messageType.forEach(restoreJsonNames);
  root.dependency = (_a = imports === null || imports === undefined ? undefined : imports.map((f) => f.proto.name)) !== null && _a !== undefined ? _a : [];
  const reg = createFileRegistry(root, (protoFileName) => imports === null || imports === undefined ? undefined : imports.find((f) => f.proto.name === protoFileName));
  return reg.getFile(root.name);
}
var file_google_protobuf_timestamp = /* @__PURE__ */ fileDesc("Ch9nb29nbGUvcHJvdG9idWYvdGltZXN0YW1wLnByb3RvEg9nb29nbGUucHJvdG9idWYiKwoJVGltZXN0YW1wEg8KB3NlY29uZHMYASABKAMSDQoFbmFub3MYAiABKAVChQEKE2NvbS5nb29nbGUucHJvdG9idWZCDlRpbWVzdGFtcFByb3RvUAFaMmdvb2dsZS5nb2xhbmcub3JnL3Byb3RvYnVmL3R5cGVzL2tub3duL3RpbWVzdGFtcHBi+AEBogIDR1BCqgIeR29vZ2xlLlByb3RvYnVmLldlbGxLbm93blR5cGVzYgZwcm90bzM");
var TimestampSchema = /* @__PURE__ */ messageDesc(file_google_protobuf_timestamp, 0);
function timestampFromDate(date) {
  return timestampFromMs(date.getTime());
}
function timestampDate(timestamp) {
  return new Date(timestampMs(timestamp));
}
function timestampFromMs(timestampMs) {
  const seconds = Math.floor(timestampMs / 1000);
  return create(TimestampSchema, {
    seconds: protoInt64.parse(seconds),
    nanos: (timestampMs - seconds * 1000) * 1e6
  });
}
function timestampMs(timestamp) {
  return Number(timestamp.seconds) * 1000 + Math.round(timestamp.nanos / 1e6);
}
var file_google_protobuf_any = /* @__PURE__ */ fileDesc("Chlnb29nbGUvcHJvdG9idWYvYW55LnByb3RvEg9nb29nbGUucHJvdG9idWYiJgoDQW55EhAKCHR5cGVfdXJsGAEgASgJEg0KBXZhbHVlGAIgASgMQnYKE2NvbS5nb29nbGUucHJvdG9idWZCCEFueVByb3RvUAFaLGdvb2dsZS5nb2xhbmcub3JnL3Byb3RvYnVmL3R5cGVzL2tub3duL2FueXBiogIDR1BCqgIeR29vZ2xlLlByb3RvYnVmLldlbGxLbm93blR5cGVzYgZwcm90bzM");
var AnySchema = /* @__PURE__ */ messageDesc(file_google_protobuf_any, 0);
var LEGACY_REQUIRED2 = 3;
var writeDefaults = {
  writeUnknownFields: true
};
function makeWriteOptions(options) {
  return options ? Object.assign(Object.assign({}, writeDefaults), options) : writeDefaults;
}
function toBinary(schema, message, options) {
  return writeFields(new BinaryWriter, makeWriteOptions(options), reflect(schema, message)).finish();
}
function writeFields(writer, opts, msg) {
  var _a;
  for (const f of msg.sortedFields) {
    if (!msg.isSet(f)) {
      if (f.presence == LEGACY_REQUIRED2) {
        throw new Error(`cannot encode ${f} to binary: required field not set`);
      }
      continue;
    }
    writeField(writer, opts, msg, f);
  }
  if (opts.writeUnknownFields) {
    for (const { no, wireType, data } of (_a = msg.getUnknown()) !== null && _a !== undefined ? _a : []) {
      writer.tag(no, wireType).raw(data);
    }
  }
  return writer;
}
function writeField(writer, opts, msg, field) {
  var _a;
  switch (field.fieldKind) {
    case "scalar":
    case "enum":
      writeScalar(writer, msg.desc.typeName, field.name, (_a = field.scalar) !== null && _a !== undefined ? _a : ScalarType.INT32, field.number, msg.get(field));
      break;
    case "list":
      writeListField(writer, opts, field, msg.get(field));
      break;
    case "message":
      writeMessageField(writer, opts, field, msg.get(field));
      break;
    case "map":
      for (const [key, val] of msg.get(field)) {
        writeMapEntry(writer, opts, field, key, val);
      }
      break;
  }
}
function writeScalar(writer, msgName, fieldName, scalarType, fieldNo, value) {
  writeScalarValue(writer.tag(fieldNo, writeTypeOfScalar(scalarType)), msgName, fieldName, scalarType, value);
}
function writeMessageField(writer, opts, field, message) {
  if (field.delimitedEncoding) {
    writeFields(writer.tag(field.number, WireType.StartGroup), opts, message).tag(field.number, WireType.EndGroup);
  } else {
    writeFields(writer.tag(field.number, WireType.LengthDelimited).fork(), opts, message).join();
  }
}
function writeListField(writer, opts, field, list) {
  var _a;
  if (field.listKind == "message") {
    for (const item of list) {
      writeMessageField(writer, opts, field, item);
    }
    return;
  }
  const scalarType = (_a = field.scalar) !== null && _a !== undefined ? _a : ScalarType.INT32;
  if (field.packed) {
    if (!list.size) {
      return;
    }
    writer.tag(field.number, WireType.LengthDelimited).fork();
    for (const item of list) {
      writeScalarValue(writer, field.parent.typeName, field.name, scalarType, item);
    }
    writer.join();
    return;
  }
  for (const item of list) {
    writeScalar(writer, field.parent.typeName, field.name, scalarType, field.number, item);
  }
}
function writeMapEntry(writer, opts, field, key, value) {
  var _a;
  writer.tag(field.number, WireType.LengthDelimited).fork();
  writeScalar(writer, field.parent.typeName, field.name, field.mapKey, 1, key);
  switch (field.mapKind) {
    case "scalar":
    case "enum":
      writeScalar(writer, field.parent.typeName, field.name, (_a = field.scalar) !== null && _a !== undefined ? _a : ScalarType.INT32, 2, value);
      break;
    case "message":
      writeFields(writer.tag(2, WireType.LengthDelimited).fork(), opts, value).join();
      break;
  }
  writer.join();
}
function writeScalarValue(writer, msgName, fieldName, type, value) {
  try {
    switch (type) {
      case ScalarType.STRING:
        writer.string(value);
        break;
      case ScalarType.BOOL:
        writer.bool(value);
        break;
      case ScalarType.DOUBLE:
        writer.double(value);
        break;
      case ScalarType.FLOAT:
        writer.float(value);
        break;
      case ScalarType.INT32:
        writer.int32(value);
        break;
      case ScalarType.INT64:
        writer.int64(value);
        break;
      case ScalarType.UINT64:
        writer.uint64(value);
        break;
      case ScalarType.FIXED64:
        writer.fixed64(value);
        break;
      case ScalarType.BYTES:
        writer.bytes(value);
        break;
      case ScalarType.FIXED32:
        writer.fixed32(value);
        break;
      case ScalarType.SFIXED32:
        writer.sfixed32(value);
        break;
      case ScalarType.SFIXED64:
        writer.sfixed64(value);
        break;
      case ScalarType.SINT64:
        writer.sint64(value);
        break;
      case ScalarType.UINT32:
        writer.uint32(value);
        break;
      case ScalarType.SINT32:
        writer.sint32(value);
        break;
    }
  } catch (e) {
    if (e instanceof Error) {
      throw new Error(`cannot encode field ${msgName}.${fieldName} to binary: ${e.message}`);
    }
    throw e;
  }
}
function writeTypeOfScalar(type) {
  switch (type) {
    case ScalarType.BYTES:
    case ScalarType.STRING:
      return WireType.LengthDelimited;
    case ScalarType.DOUBLE:
    case ScalarType.FIXED64:
    case ScalarType.SFIXED64:
      return WireType.Bit64;
    case ScalarType.FIXED32:
    case ScalarType.SFIXED32:
    case ScalarType.FLOAT:
      return WireType.Bit32;
    default:
      return WireType.Varint;
  }
}
function anyPack(schema, message, into) {
  let ret = false;
  if (!into) {
    into = create(AnySchema);
    ret = true;
  }
  into.value = toBinary(schema, message);
  into.typeUrl = typeNameToUrl(message.$typeName);
  return ret ? into : undefined;
}
function anyIs(any, descOrTypeName) {
  if (any.typeUrl === "") {
    return false;
  }
  const want = typeof descOrTypeName == "string" ? descOrTypeName : descOrTypeName.typeName;
  const got = typeUrlToName(any.typeUrl);
  return want === got;
}
function anyUnpack(any, registryOrMessageDesc) {
  if (any.typeUrl === "") {
    return;
  }
  const desc = registryOrMessageDesc.kind == "message" ? registryOrMessageDesc : registryOrMessageDesc.getMessage(typeUrlToName(any.typeUrl));
  if (!desc || !anyIs(any, desc)) {
    return;
  }
  return fromBinary(desc, any.value);
}
function typeNameToUrl(name) {
  return `type.googleapis.com/${name}`;
}
function typeUrlToName(url) {
  const slash = url.lastIndexOf("/");
  const name = slash >= 0 ? url.substring(slash + 1) : url;
  if (!name.length) {
    throw new Error(`invalid type url: ${url}`);
  }
  return name;
}
var file_google_protobuf_empty = /* @__PURE__ */ fileDesc("Chtnb29nbGUvcHJvdG9idWYvZW1wdHkucHJvdG8SD2dvb2dsZS5wcm90b2J1ZiIHCgVFbXB0eUJ9ChNjb20uZ29vZ2xlLnByb3RvYnVmQgpFbXB0eVByb3RvUAFaLmdvb2dsZS5nb2xhbmcub3JnL3Byb3RvYnVmL3R5cGVzL2tub3duL2VtcHR5cGL4AQGiAgNHUEKqAh5Hb29nbGUuUHJvdG9idWYuV2VsbEtub3duVHlwZXNiBnByb3RvMw");
var EmptySchema = /* @__PURE__ */ messageDesc(file_google_protobuf_empty, 0);
var file_google_protobuf_struct = /* @__PURE__ */ fileDesc("Chxnb29nbGUvcHJvdG9idWYvc3RydWN0LnByb3RvEg9nb29nbGUucHJvdG9idWYihAEKBlN0cnVjdBIzCgZmaWVsZHMYASADKAsyIy5nb29nbGUucHJvdG9idWYuU3RydWN0LkZpZWxkc0VudHJ5GkUKC0ZpZWxkc0VudHJ5EgsKA2tleRgBIAEoCRIlCgV2YWx1ZRgCIAEoCzIWLmdvb2dsZS5wcm90b2J1Zi5WYWx1ZToCOAEi6gEKBVZhbHVlEjAKCm51bGxfdmFsdWUYASABKA4yGi5nb29nbGUucHJvdG9idWYuTnVsbFZhbHVlSAASFgoMbnVtYmVyX3ZhbHVlGAIgASgBSAASFgoMc3RyaW5nX3ZhbHVlGAMgASgJSAASFAoKYm9vbF92YWx1ZRgEIAEoCEgAEi8KDHN0cnVjdF92YWx1ZRgFIAEoCzIXLmdvb2dsZS5wcm90b2J1Zi5TdHJ1Y3RIABIwCgpsaXN0X3ZhbHVlGAYgASgLMhouZ29vZ2xlLnByb3RvYnVmLkxpc3RWYWx1ZUgAQgYKBGtpbmQiMwoJTGlzdFZhbHVlEiYKBnZhbHVlcxgBIAMoCzIWLmdvb2dsZS5wcm90b2J1Zi5WYWx1ZSobCglOdWxsVmFsdWUSDgoKTlVMTF9WQUxVRRAAQn8KE2NvbS5nb29nbGUucHJvdG9idWZCC1N0cnVjdFByb3RvUAFaL2dvb2dsZS5nb2xhbmcub3JnL3Byb3RvYnVmL3R5cGVzL2tub3duL3N0cnVjdHBi+AEBogIDR1BCqgIeR29vZ2xlLlByb3RvYnVmLldlbGxLbm93blR5cGVzYgZwcm90bzM");
var StructSchema = /* @__PURE__ */ messageDesc(file_google_protobuf_struct, 0);
var ValueSchema = /* @__PURE__ */ messageDesc(file_google_protobuf_struct, 1);
var ListValueSchema = /* @__PURE__ */ messageDesc(file_google_protobuf_struct, 2);
var NullValue;
(function(NullValue2) {
  NullValue2[NullValue2["NULL_VALUE"] = 0] = "NULL_VALUE";
})(NullValue || (NullValue = {}));
function setExtension(message, extension, value) {
  var _a;
  assertExtendee(extension, message);
  const ufs = ((_a = message.$unknown) !== null && _a !== undefined ? _a : []).filter((uf) => uf.no !== extension.number);
  const [container, field] = createExtensionContainer(extension, value);
  const writer = new BinaryWriter;
  writeField(writer, { writeUnknownFields: true }, container, field);
  const reader = new BinaryReader(writer.finish());
  while (reader.pos < reader.len) {
    const [no, wireType] = reader.tag();
    const data = reader.skip(wireType, no);
    ufs.push({ no, wireType, data });
  }
  message.$unknown = ufs;
}
function createExtensionContainer(extension, value) {
  const localName = extension.typeName;
  const field = Object.assign(Object.assign({}, extension), { kind: "field", parent: extension.extendee, localName });
  const desc = Object.assign(Object.assign({}, extension.extendee), { fields: [field], members: [field], oneofs: [] });
  const container = create(desc, value !== undefined ? { [localName]: value } : undefined);
  return [
    reflect(desc, container),
    field,
    () => {
      const value2 = container[localName];
      if (value2 === undefined) {
        const desc2 = extension.message;
        if (isWrapperDesc(desc2)) {
          return scalarZeroValue(desc2.fields[0].scalar, desc2.fields[0].longAsString);
        }
        return create(desc2);
      }
      return value2;
    }
  ];
}
function assertExtendee(extension, message) {
  if (extension.extendee.typeName != message.$typeName) {
    throw new Error(`extension ${extension.typeName} can only be applied to message ${extension.extendee.typeName}`);
  }
}
var jsonReadDefaults = {
  ignoreUnknownFields: false
};
function makeReadOptions2(options) {
  return options ? Object.assign(Object.assign({}, jsonReadDefaults), options) : jsonReadDefaults;
}
function fromJson(schema, json, options) {
  const msg = reflect(schema);
  try {
    readMessage2(msg, json, makeReadOptions2(options));
  } catch (e) {
    if (isFieldError(e)) {
      throw new Error(`cannot decode ${e.field()} from JSON: ${e.message}`, {
        cause: e
      });
    }
    throw e;
  }
  return msg.message;
}
function readMessage2(msg, json, opts) {
  var _a;
  if (tryWktFromJson(msg, json, opts)) {
    return;
  }
  if (json == null || Array.isArray(json) || typeof json != "object") {
    throw new Error(`cannot decode ${msg.desc} from JSON: ${formatVal(json)}`);
  }
  const oneofSeen = new Map;
  const jsonNames = new Map;
  for (const field of msg.desc.fields) {
    jsonNames.set(field.name, field).set(field.jsonName, field);
  }
  for (const [jsonKey, jsonValue] of Object.entries(json)) {
    const field = jsonNames.get(jsonKey);
    if (field) {
      if (field.oneof) {
        if (jsonValue === null && field.fieldKind == "scalar") {
          continue;
        }
        const seen = oneofSeen.get(field.oneof);
        if (seen !== undefined) {
          throw new FieldError(field.oneof, `oneof set multiple times by ${seen.name} and ${field.name}`);
        }
        oneofSeen.set(field.oneof, field);
      }
      readField2(msg, field, jsonValue, opts);
    } else {
      let extension = undefined;
      if (jsonKey.startsWith("[") && jsonKey.endsWith("]") && (extension = (_a = opts.registry) === null || _a === undefined ? undefined : _a.getExtension(jsonKey.substring(1, jsonKey.length - 1))) && extension.extendee.typeName === msg.desc.typeName) {
        const [container, field2, get] = createExtensionContainer(extension);
        readField2(container, field2, jsonValue, opts);
        setExtension(msg.message, extension, get());
      }
      if (!extension && !opts.ignoreUnknownFields) {
        throw new Error(`cannot decode ${msg.desc} from JSON: key "${jsonKey}" is unknown`);
      }
    }
  }
}
function readField2(msg, field, json, opts) {
  switch (field.fieldKind) {
    case "scalar":
      readScalarField(msg, field, json);
      break;
    case "enum":
      readEnumField(msg, field, json, opts);
      break;
    case "message":
      readMessageField2(msg, field, json, opts);
      break;
    case "list":
      readListField2(msg.get(field), json, opts);
      break;
    case "map":
      readMapField(msg.get(field), json, opts);
      break;
  }
}
function readMapField(map, json, opts) {
  if (json === null) {
    return;
  }
  const field = map.field();
  if (typeof json != "object" || Array.isArray(json)) {
    throw new FieldError(field, "expected object, got " + formatVal(json));
  }
  for (const [jsonMapKey, jsonMapValue] of Object.entries(json)) {
    if (jsonMapValue === null) {
      throw new FieldError(field, "map value must not be null");
    }
    let value;
    switch (field.mapKind) {
      case "message":
        const msgValue = reflect(field.message);
        readMessage2(msgValue, jsonMapValue, opts);
        value = msgValue;
        break;
      case "enum":
        value = readEnum(field.enum, jsonMapValue, opts.ignoreUnknownFields, true);
        if (value === tokenIgnoredUnknownEnum) {
          return;
        }
        break;
      case "scalar":
        value = scalarFromJson(field, jsonMapValue, true);
        break;
    }
    const key = mapKeyFromJson(field.mapKey, jsonMapKey);
    map.set(key, value);
  }
}
function readListField2(list, json, opts) {
  if (json === null) {
    return;
  }
  const field = list.field();
  if (!Array.isArray(json)) {
    throw new FieldError(field, "expected Array, got " + formatVal(json));
  }
  for (const jsonItem of json) {
    if (jsonItem === null) {
      throw new FieldError(field, "list item must not be null");
    }
    switch (field.listKind) {
      case "message":
        const msgValue = reflect(field.message);
        readMessage2(msgValue, jsonItem, opts);
        list.add(msgValue);
        break;
      case "enum":
        const enumValue = readEnum(field.enum, jsonItem, opts.ignoreUnknownFields, true);
        if (enumValue !== tokenIgnoredUnknownEnum) {
          list.add(enumValue);
        }
        break;
      case "scalar":
        list.add(scalarFromJson(field, jsonItem, true));
        break;
    }
  }
}
function readMessageField2(msg, field, json, opts) {
  if (json === null && field.message.typeName != "google.protobuf.Value") {
    msg.clear(field);
    return;
  }
  const msgValue = msg.isSet(field) ? msg.get(field) : reflect(field.message);
  readMessage2(msgValue, json, opts);
  msg.set(field, msgValue);
}
function readEnumField(msg, field, json, opts) {
  const enumValue = readEnum(field.enum, json, opts.ignoreUnknownFields, false);
  if (enumValue === tokenNull) {
    msg.clear(field);
  } else if (enumValue !== tokenIgnoredUnknownEnum) {
    msg.set(field, enumValue);
  }
}
function readScalarField(msg, field, json) {
  const scalarValue = scalarFromJson(field, json, false);
  if (scalarValue === tokenNull) {
    msg.clear(field);
  } else {
    msg.set(field, scalarValue);
  }
}
var tokenIgnoredUnknownEnum = Symbol();
function readEnum(desc, json, ignoreUnknownFields, nullAsZeroValue) {
  if (json === null) {
    if (desc.typeName == "google.protobuf.NullValue") {
      return 0;
    }
    return nullAsZeroValue ? desc.values[0].number : tokenNull;
  }
  switch (typeof json) {
    case "number":
      if (Number.isInteger(json)) {
        return json;
      }
      break;
    case "string":
      const value = desc.values.find((ev) => ev.name === json);
      if (value !== undefined) {
        return value.number;
      }
      if (ignoreUnknownFields) {
        return tokenIgnoredUnknownEnum;
      }
      break;
  }
  throw new Error(`cannot decode ${desc} from JSON: ${formatVal(json)}`);
}
var tokenNull = Symbol();
function scalarFromJson(field, json, nullAsZeroValue) {
  if (json === null) {
    if (nullAsZeroValue) {
      return scalarZeroValue(field.scalar, false);
    }
    return tokenNull;
  }
  switch (field.scalar) {
    case ScalarType.DOUBLE:
    case ScalarType.FLOAT:
      if (json === "NaN")
        return NaN;
      if (json === "Infinity")
        return Number.POSITIVE_INFINITY;
      if (json === "-Infinity")
        return Number.NEGATIVE_INFINITY;
      if (typeof json == "number") {
        if (Number.isNaN(json)) {
          throw new FieldError(field, "unexpected NaN number");
        }
        if (!Number.isFinite(json)) {
          throw new FieldError(field, "unexpected infinite number");
        }
        break;
      }
      if (typeof json == "string") {
        if (json === "") {
          break;
        }
        if (json.trim().length !== json.length) {
          break;
        }
        const float = Number(json);
        if (!Number.isFinite(float)) {
          break;
        }
        return float;
      }
      break;
    case ScalarType.INT32:
    case ScalarType.FIXED32:
    case ScalarType.SFIXED32:
    case ScalarType.SINT32:
    case ScalarType.UINT32:
      return int32FromJson(json);
    case ScalarType.BYTES:
      if (typeof json == "string") {
        if (json === "") {
          return new Uint8Array(0);
        }
        try {
          return base64Decode(json);
        } catch (e) {
          const message = e instanceof Error ? e.message : String(e);
          throw new FieldError(field, message);
        }
      }
      break;
  }
  return json;
}
function mapKeyFromJson(type, json) {
  switch (type) {
    case ScalarType.BOOL:
      switch (json) {
        case "true":
          return true;
        case "false":
          return false;
      }
      return json;
    case ScalarType.INT32:
    case ScalarType.FIXED32:
    case ScalarType.UINT32:
    case ScalarType.SFIXED32:
    case ScalarType.SINT32:
      return int32FromJson(json);
    default:
      return json;
  }
}
function int32FromJson(json) {
  if (typeof json == "string") {
    if (json === "") {
      return json;
    }
    if (json.trim().length !== json.length) {
      return json;
    }
    const num = Number(json);
    if (Number.isNaN(num)) {
      return json;
    }
    return num;
  }
  return json;
}
function tryWktFromJson(msg, jsonValue, opts) {
  if (!msg.desc.typeName.startsWith("google.protobuf.")) {
    return false;
  }
  switch (msg.desc.typeName) {
    case "google.protobuf.Any":
      anyFromJson(msg.message, jsonValue, opts);
      return true;
    case "google.protobuf.Timestamp":
      timestampFromJson(msg.message, jsonValue);
      return true;
    case "google.protobuf.Duration":
      durationFromJson(msg.message, jsonValue);
      return true;
    case "google.protobuf.FieldMask":
      fieldMaskFromJson(msg.message, jsonValue);
      return true;
    case "google.protobuf.Struct":
      structFromJson(msg.message, jsonValue);
      return true;
    case "google.protobuf.Value":
      valueFromJson(msg.message, jsonValue);
      return true;
    case "google.protobuf.ListValue":
      listValueFromJson(msg.message, jsonValue);
      return true;
    default:
      if (isWrapperDesc(msg.desc)) {
        const valueField = msg.desc.fields[0];
        if (jsonValue === null) {
          msg.clear(valueField);
        } else {
          msg.set(valueField, scalarFromJson(valueField, jsonValue, true));
        }
        return true;
      }
      return false;
  }
}
function anyFromJson(any, json, opts) {
  var _a;
  if (json === null || Array.isArray(json) || typeof json != "object") {
    throw new Error(`cannot decode message ${any.$typeName} from JSON: expected object but got ${formatVal(json)}`);
  }
  if (Object.keys(json).length == 0) {
    return;
  }
  const typeUrl = json["@type"];
  if (typeof typeUrl != "string" || typeUrl == "") {
    throw new Error(`cannot decode message ${any.$typeName} from JSON: "@type" is empty`);
  }
  const typeName = typeUrl.includes("/") ? typeUrl.substring(typeUrl.lastIndexOf("/") + 1) : typeUrl;
  if (!typeName.length) {
    throw new Error(`cannot decode message ${any.$typeName} from JSON: "@type" is invalid`);
  }
  const desc = (_a = opts.registry) === null || _a === undefined ? undefined : _a.getMessage(typeName);
  if (!desc) {
    throw new Error(`cannot decode message ${any.$typeName} from JSON: ${typeUrl} is not in the type registry`);
  }
  const msg = reflect(desc);
  if (typeName.startsWith("google.protobuf.") && Object.prototype.hasOwnProperty.call(json, "value")) {
    const value = json.value;
    readMessage2(msg, value, opts);
  } else {
    const copy = Object.assign({}, json);
    delete copy["@type"];
    readMessage2(msg, copy, opts);
  }
  anyPack(msg.desc, msg.message, any);
}
function timestampFromJson(timestamp, json) {
  if (typeof json !== "string") {
    throw new Error(`cannot decode message ${timestamp.$typeName} from JSON: ${formatVal(json)}`);
  }
  const matches = json.match(/^([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})(?:\.([0-9]{1,9}))?(?:Z|([+-][0-9][0-9]:[0-9][0-9]))$/);
  if (!matches) {
    throw new Error(`cannot decode message ${timestamp.$typeName} from JSON: invalid RFC 3339 string`);
  }
  const ms = Date.parse(matches[1] + "-" + matches[2] + "-" + matches[3] + "T" + matches[4] + ":" + matches[5] + ":" + matches[6] + (matches[8] ? matches[8] : "Z"));
  if (Number.isNaN(ms)) {
    throw new Error(`cannot decode message ${timestamp.$typeName} from JSON: invalid RFC 3339 string`);
  }
  if (ms < Date.parse("0001-01-01T00:00:00Z") || ms > Date.parse("9999-12-31T23:59:59Z")) {
    throw new Error(`cannot decode message ${timestamp.$typeName} from JSON: must be from 0001-01-01T00:00:00Z to 9999-12-31T23:59:59Z inclusive`);
  }
  timestamp.seconds = protoInt64.parse(ms / 1000);
  timestamp.nanos = 0;
  if (matches[7]) {
    timestamp.nanos = parseInt("1" + matches[7] + "0".repeat(9 - matches[7].length)) - 1e9;
  }
}
function durationFromJson(duration, json) {
  if (typeof json !== "string") {
    throw new Error(`cannot decode message ${duration.$typeName} from JSON: ${formatVal(json)}`);
  }
  const match = json.match(/^(-?[0-9]+)(?:\.([0-9]+))?s/);
  if (match === null) {
    throw new Error(`cannot decode message ${duration.$typeName} from JSON: ${formatVal(json)}`);
  }
  const longSeconds = Number(match[1]);
  if (longSeconds > 315576000000 || longSeconds < -315576000000) {
    throw new Error(`cannot decode message ${duration.$typeName} from JSON: ${formatVal(json)}`);
  }
  duration.seconds = protoInt64.parse(longSeconds);
  if (typeof match[2] !== "string") {
    return;
  }
  const nanosStr = match[2] + "0".repeat(9 - match[2].length);
  duration.nanos = parseInt(nanosStr);
  if (longSeconds < 0 || Object.is(longSeconds, -0)) {
    duration.nanos = -duration.nanos;
  }
}
function fieldMaskFromJson(fieldMask, json) {
  if (typeof json !== "string") {
    throw new Error(`cannot decode message ${fieldMask.$typeName} from JSON: ${formatVal(json)}`);
  }
  if (json === "") {
    return;
  }
  function camelToSnake(str) {
    if (str.includes("_")) {
      throw new Error(`cannot decode message ${fieldMask.$typeName} from JSON: path names must be lowerCamelCase`);
    }
    const sc = str.replace(/[A-Z]/g, (letter) => "_" + letter.toLowerCase());
    return sc[0] === "_" ? sc.substring(1) : sc;
  }
  fieldMask.paths = json.split(",").map(camelToSnake);
}
function structFromJson(struct, json) {
  if (typeof json != "object" || json == null || Array.isArray(json)) {
    throw new Error(`cannot decode message ${struct.$typeName} from JSON ${formatVal(json)}`);
  }
  for (const [k, v] of Object.entries(json)) {
    const parsedV = create(ValueSchema);
    valueFromJson(parsedV, v);
    struct.fields[k] = parsedV;
  }
}
function valueFromJson(value, json) {
  switch (typeof json) {
    case "number":
      value.kind = { case: "numberValue", value: json };
      break;
    case "string":
      value.kind = { case: "stringValue", value: json };
      break;
    case "boolean":
      value.kind = { case: "boolValue", value: json };
      break;
    case "object":
      if (json === null) {
        value.kind = { case: "nullValue", value: NullValue.NULL_VALUE };
      } else if (Array.isArray(json)) {
        const listValue = create(ListValueSchema);
        listValueFromJson(listValue, json);
        value.kind = { case: "listValue", value: listValue };
      } else {
        const struct = create(StructSchema);
        structFromJson(struct, json);
        value.kind = { case: "structValue", value: struct };
      }
      break;
    default:
      throw new Error(`cannot decode message ${value.$typeName} from JSON ${formatVal(json)}`);
  }
  return value;
}
function listValueFromJson(listValue, json) {
  if (!Array.isArray(json)) {
    throw new Error(`cannot decode message ${listValue.$typeName} from JSON ${formatVal(json)}`);
  }
  for (const e of json) {
    const value = create(ValueSchema);
    valueFromJson(value, e);
    listValue.values.push(value);
  }
}
var file_values_v1_values = /* @__PURE__ */ fileDesc("ChZ2YWx1ZXMvdjEvdmFsdWVzLnByb3RvEgl2YWx1ZXMudjEigQMKBVZhbHVlEhYKDHN0cmluZ192YWx1ZRgBIAEoCUgAEhQKCmJvb2xfdmFsdWUYAiABKAhIABIVCgtieXRlc192YWx1ZRgDIAEoDEgAEiMKCW1hcF92YWx1ZRgEIAEoCzIOLnZhbHVlcy52MS5NYXBIABIlCgpsaXN0X3ZhbHVlGAUgASgLMg8udmFsdWVzLnYxLkxpc3RIABIrCg1kZWNpbWFsX3ZhbHVlGAYgASgLMhIudmFsdWVzLnYxLkRlY2ltYWxIABIZCgtpbnQ2NF92YWx1ZRgHIAEoA0ICMABIABIpCgxiaWdpbnRfdmFsdWUYCSABKAsyES52YWx1ZXMudjEuQmlnSW50SAASMAoKdGltZV92YWx1ZRgKIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBIABIXCg1mbG9hdDY0X3ZhbHVlGAsgASgBSAASGgoMdWludDY0X3ZhbHVlGAwgASgEQgIwAEgAQgcKBXZhbHVlSgQICBAJIisKBkJpZ0ludBIPCgdhYnNfdmFsGAEgASgMEhAKBHNpZ24YAiABKANCAjAAInIKA01hcBIqCgZmaWVsZHMYASADKAsyGi52YWx1ZXMudjEuTWFwLkZpZWxkc0VudHJ5Gj8KC0ZpZWxkc0VudHJ5EgsKA2tleRgBIAEoCRIfCgV2YWx1ZRgCIAEoCzIQLnZhbHVlcy52MS5WYWx1ZToCOAEiKAoETGlzdBIgCgZmaWVsZHMYAiADKAsyEC52YWx1ZXMudjEuVmFsdWUiQwoHRGVjaW1hbBImCgtjb2VmZmljaWVudBgBIAEoCzIRLnZhbHVlcy52MS5CaWdJbnQSEAoIZXhwb25lbnQYAiABKAVCYQoNY29tLnZhbHVlcy52MUILVmFsdWVzUHJvdG9QAaICA1ZYWKoCCVZhbHVlcy5WMcoCCVZhbHVlc1xWMeICFVZhbHVlc1xWMVxHUEJNZXRhZGF0YeoCClZhbHVlczo6VjFiBnByb3RvMw", [file_google_protobuf_timestamp]);
var ValueSchema2 = /* @__PURE__ */ messageDesc(file_values_v1_values, 0);
var BigIntSchema = /* @__PURE__ */ messageDesc(file_values_v1_values, 1);
var MapSchema = /* @__PURE__ */ messageDesc(file_values_v1_values, 2);
var ListSchema = /* @__PURE__ */ messageDesc(file_values_v1_values, 3);
var DecimalSchema = /* @__PURE__ */ messageDesc(file_values_v1_values, 4);
var file_sdk_v1alpha_sdk = /* @__PURE__ */ fileDesc("ChVzZGsvdjFhbHBoYS9zZGsucHJvdG8SC3Nkay52MWFscGhhIrQBChVTaW1wbGVDb25zZW5zdXNJbnB1dHMSIQoFdmFsdWUYASABKAsyEC52YWx1ZXMudjEuVmFsdWVIABIPCgVlcnJvchgCIAEoCUgAEjUKC2Rlc2NyaXB0b3JzGAMgASgLMiAuc2RrLnYxYWxwaGEuQ29uc2Vuc3VzRGVzY3JpcHRvchIhCgdkZWZhdWx0GAQgASgLMhAudmFsdWVzLnYxLlZhbHVlQg0KC29ic2VydmF0aW9uIpABCglGaWVsZHNNYXASMgoGZmllbGRzGAEgAygLMiIuc2RrLnYxYWxwaGEuRmllbGRzTWFwLkZpZWxkc0VudHJ5Gk8KC0ZpZWxkc0VudHJ5EgsKA2tleRgBIAEoCRIvCgV2YWx1ZRgCIAEoCzIgLnNkay52MWFscGhhLkNvbnNlbnN1c0Rlc2NyaXB0b3I6AjgBIoYBChNDb25zZW5zdXNEZXNjcmlwdG9yEjMKC2FnZ3JlZ2F0aW9uGAEgASgOMhwuc2RrLnYxYWxwaGEuQWdncmVnYXRpb25UeXBlSAASLAoKZmllbGRzX21hcBgCIAEoCzIWLnNkay52MWFscGhhLkZpZWxkc01hcEgAQgwKCmRlc2NyaXB0b3IiagoNUmVwb3J0UmVxdWVzdBIXCg9lbmNvZGVkX3BheWxvYWQYASABKAwSFAoMZW5jb2Rlcl9uYW1lGAIgASgJEhQKDHNpZ25pbmdfYWxnbxgDIAEoCRIUCgxoYXNoaW5nX2FsZ28YBCABKAkilwEKDlJlcG9ydFJlc3BvbnNlEhUKDWNvbmZpZ19kaWdlc3QYASABKAwSEgoGc2VxX25yGAIgASgEQgIwABIWCg5yZXBvcnRfY29udGV4dBgDIAEoDBISCgpyYXdfcmVwb3J0GAQgASgMEi4KBHNpZ3MYBSADKAsyIC5zZGsudjFhbHBoYS5BdHRyaWJ1dGVkU2lnbmF0dXJlIjsKE0F0dHJpYnV0ZWRTaWduYXR1cmUSEQoJc2lnbmF0dXJlGAEgASgMEhEKCXNpZ25lcl9pZBgCIAEoDSJrChFDYXBhYmlsaXR5UmVxdWVzdBIKCgJpZBgBIAEoCRIlCgdwYXlsb2FkGAIgASgLMhQuZ29vZ2xlLnByb3RvYnVmLkFueRIOCgZtZXRob2QYAyABKAkSEwoLY2FsbGJhY2tfaWQYBCABKAUiWgoSQ2FwYWJpbGl0eVJlc3BvbnNlEicKB3BheWxvYWQYASABKAsyFC5nb29nbGUucHJvdG9idWYuQW55SAASDwoFZXJyb3IYAiABKAlIAEIKCghyZXNwb25zZSJYChNUcmlnZ2VyU3Vic2NyaXB0aW9uEgoKAmlkGAEgASgJEiUKB3BheWxvYWQYAiABKAsyFC5nb29nbGUucHJvdG9idWYuQW55Eg4KBm1ldGhvZBgDIAEoCSJVChpUcmlnZ2VyU3Vic2NyaXB0aW9uUmVxdWVzdBI3Cg1zdWJzY3JpcHRpb25zGAEgAygLMiAuc2RrLnYxYWxwaGEuVHJpZ2dlclN1YnNjcmlwdGlvbiJACgdUcmlnZ2VyEg4KAmlkGAEgASgEQgIwABIlCgdwYXlsb2FkGAIgASgLMhQuZ29vZ2xlLnByb3RvYnVmLkFueSInChhBd2FpdENhcGFiaWxpdGllc1JlcXVlc3QSCwoDaWRzGAEgAygFIrgBChlBd2FpdENhcGFiaWxpdGllc1Jlc3BvbnNlEkgKCXJlc3BvbnNlcxgBIAMoCzI1LnNkay52MWFscGhhLkF3YWl0Q2FwYWJpbGl0aWVzUmVzcG9uc2UuUmVzcG9uc2VzRW50cnkaUQoOUmVzcG9uc2VzRW50cnkSCwoDa2V5GAEgASgFEi4KBXZhbHVlGAIgASgLMh8uc2RrLnYxYWxwaGEuQ2FwYWJpbGl0eVJlc3BvbnNlOgI4ASKgAQoORXhlY3V0ZVJlcXVlc3QSDgoGY29uZmlnGAEgASgMEisKCXN1YnNjcmliZRgCIAEoCzIWLmdvb2dsZS5wcm90b2J1Zi5FbXB0eUgAEicKB3RyaWdnZXIYAyABKAsyFC5zZGsudjFhbHBoYS5UcmlnZ2VySAASHQoRbWF4X3Jlc3BvbnNlX3NpemUYBCABKARCAjAAQgkKB3JlcXVlc3QimQEKD0V4ZWN1dGlvblJlc3VsdBIhCgV2YWx1ZRgBIAEoCzIQLnZhbHVlcy52MS5WYWx1ZUgAEg8KBWVycm9yGAIgASgJSAASSAoVdHJpZ2dlcl9zdWJzY3JpcHRpb25zGAMgASgLMicuc2RrLnYxYWxwaGEuVHJpZ2dlclN1YnNjcmlwdGlvblJlcXVlc3RIAEIICgZyZXN1bHQiVgoRR2V0U2VjcmV0c1JlcXVlc3QSLAoIcmVxdWVzdHMYASADKAsyGi5zZGsudjFhbHBoYS5TZWNyZXRSZXF1ZXN0EhMKC2NhbGxiYWNrX2lkGAIgASgFIiIKE0F3YWl0U2VjcmV0c1JlcXVlc3QSCwoDaWRzGAEgAygFIqsBChRBd2FpdFNlY3JldHNSZXNwb25zZRJDCglyZXNwb25zZXMYASADKAsyMC5zZGsudjFhbHBoYS5Bd2FpdFNlY3JldHNSZXNwb25zZS5SZXNwb25zZXNFbnRyeRpOCg5SZXNwb25zZXNFbnRyeRILCgNrZXkYASABKAUSKwoFdmFsdWUYAiABKAsyHC5zZGsudjFhbHBoYS5TZWNyZXRSZXNwb25zZXM6AjgBIi4KDVNlY3JldFJlcXVlc3QSCgoCaWQYASABKAkSEQoJbmFtZXNwYWNlGAIgASgJIkUKBlNlY3JldBIKCgJpZBgBIAEoCRIRCgluYW1lc3BhY2UYAiABKAkSDQoFb3duZXIYAyABKAkSDQoFdmFsdWUYBCABKAkiSgoLU2VjcmV0RXJyb3ISCgoCaWQYASABKAkSEQoJbmFtZXNwYWNlGAIgASgJEg0KBW93bmVyGAMgASgJEg0KBWVycm9yGAQgASgJIm4KDlNlY3JldFJlc3BvbnNlEiUKBnNlY3JldBgBIAEoCzITLnNkay52MWFscGhhLlNlY3JldEgAEikKBWVycm9yGAIgASgLMhguc2RrLnYxYWxwaGEuU2VjcmV0RXJyb3JIAEIKCghyZXNwb25zZSJBCg9TZWNyZXRSZXNwb25zZXMSLgoJcmVzcG9uc2VzGAEgAygLMhsuc2RrLnYxYWxwaGEuU2VjcmV0UmVzcG9uc2UquAEKD0FnZ3JlZ2F0aW9uVHlwZRIgChxBR0dSRUdBVElPTl9UWVBFX1VOU1BFQ0lGSUVEEAASGwoXQUdHUkVHQVRJT05fVFlQRV9NRURJQU4QARIeChpBR0dSRUdBVElPTl9UWVBFX0lERU5USUNBTBACEiIKHkFHR1JFR0FUSU9OX1RZUEVfQ09NTU9OX1BSRUZJWBADEiIKHkFHR1JFR0FUSU9OX1RZUEVfQ09NTU9OX1NVRkZJWBAEKjkKBE1vZGUSFAoQTU9ERV9VTlNQRUNJRklFRBAAEgwKCE1PREVfRE9OEAESDQoJTU9ERV9OT0RFEAJCaAoPY29tLnNkay52MWFscGhhQghTZGtQcm90b1ABogIDU1hYqgILU2RrLlYxYWxwaGHKAgtTZGtcVjFhbHBoYeICF1Nka1xWMWFscGhhXEdQQk1ldGFkYXRh6gIMU2RrOjpWMWFscGhhYgZwcm90bzM", [file_google_protobuf_any, file_google_protobuf_empty, file_values_v1_values]);
var SimpleConsensusInputsSchema = /* @__PURE__ */ messageDesc(file_sdk_v1alpha_sdk, 0);
var ReportRequestSchema = /* @__PURE__ */ messageDesc(file_sdk_v1alpha_sdk, 3);
var ReportResponseSchema = /* @__PURE__ */ messageDesc(file_sdk_v1alpha_sdk, 4);
var CapabilityRequestSchema = /* @__PURE__ */ messageDesc(file_sdk_v1alpha_sdk, 6);
var TriggerSubscriptionRequestSchema = /* @__PURE__ */ messageDesc(file_sdk_v1alpha_sdk, 9);
var AwaitCapabilitiesRequestSchema = /* @__PURE__ */ messageDesc(file_sdk_v1alpha_sdk, 11);
var AwaitCapabilitiesResponseSchema = /* @__PURE__ */ messageDesc(file_sdk_v1alpha_sdk, 12);
var ExecuteRequestSchema = /* @__PURE__ */ messageDesc(file_sdk_v1alpha_sdk, 13);
var ExecutionResultSchema = /* @__PURE__ */ messageDesc(file_sdk_v1alpha_sdk, 14);
var GetSecretsRequestSchema = /* @__PURE__ */ messageDesc(file_sdk_v1alpha_sdk, 15);
var AwaitSecretsRequestSchema = /* @__PURE__ */ messageDesc(file_sdk_v1alpha_sdk, 16);
var AwaitSecretsResponseSchema = /* @__PURE__ */ messageDesc(file_sdk_v1alpha_sdk, 17);
var SecretRequestSchema = /* @__PURE__ */ messageDesc(file_sdk_v1alpha_sdk, 18);
var AggregationType;
(function(AggregationType2) {
  AggregationType2[AggregationType2["UNSPECIFIED"] = 0] = "UNSPECIFIED";
  AggregationType2[AggregationType2["MEDIAN"] = 1] = "MEDIAN";
  AggregationType2[AggregationType2["IDENTICAL"] = 2] = "IDENTICAL";
  AggregationType2[AggregationType2["COMMON_PREFIX"] = 3] = "COMMON_PREFIX";
  AggregationType2[AggregationType2["COMMON_SUFFIX"] = 4] = "COMMON_SUFFIX";
})(AggregationType || (AggregationType = {}));
var Mode;
(function(Mode2) {
  Mode2[Mode2["UNSPECIFIED"] = 0] = "UNSPECIFIED";
  Mode2[Mode2["DON"] = 1] = "DON";
  Mode2[Mode2["NODE"] = 2] = "NODE";
})(Mode || (Mode = {}));
var file_tools_generator_v1alpha_cre_metadata = /* @__PURE__ */ fileDesc("Cip0b29scy9nZW5lcmF0b3IvdjFhbHBoYS9jcmVfbWV0YWRhdGEucHJvdG8SF3Rvb2xzLmdlbmVyYXRvci52MWFscGhhIoQBCgtTdHJpbmdMYWJlbBJECghkZWZhdWx0cxgBIAMoCzIyLnRvb2xzLmdlbmVyYXRvci52MWFscGhhLlN0cmluZ0xhYmVsLkRlZmF1bHRzRW50cnkaLwoNRGVmYXVsdHNFbnRyeRILCgNrZXkYASABKAkSDQoFdmFsdWUYAiABKAk6AjgBIogBCgtVaW50NjRMYWJlbBJECghkZWZhdWx0cxgBIAMoCzIyLnRvb2xzLmdlbmVyYXRvci52MWFscGhhLlVpbnQ2NExhYmVsLkRlZmF1bHRzRW50cnkaMwoNRGVmYXVsdHNFbnRyeRILCgNrZXkYASABKAkSEQoFdmFsdWUYAiABKARCAjAAOgI4ASKEAQoLVWludDMyTGFiZWwSRAoIZGVmYXVsdHMYASADKAsyMi50b29scy5nZW5lcmF0b3IudjFhbHBoYS5VaW50MzJMYWJlbC5EZWZhdWx0c0VudHJ5Gi8KDURlZmF1bHRzRW50cnkSCwoDa2V5GAEgASgJEg0KBXZhbHVlGAIgASgNOgI4ASKGAQoKSW50NjRMYWJlbBJDCghkZWZhdWx0cxgBIAMoCzIxLnRvb2xzLmdlbmVyYXRvci52MWFscGhhLkludDY0TGFiZWwuRGVmYXVsdHNFbnRyeRozCg1EZWZhdWx0c0VudHJ5EgsKA2tleRgBIAEoCRIRCgV2YWx1ZRgCIAEoA0ICMAA6AjgBIoIBCgpJbnQzMkxhYmVsEkMKCGRlZmF1bHRzGAEgAygLMjEudG9vbHMuZ2VuZXJhdG9yLnYxYWxwaGEuSW50MzJMYWJlbC5EZWZhdWx0c0VudHJ5Gi8KDURlZmF1bHRzRW50cnkSCwoDa2V5GAEgASgJEg0KBXZhbHVlGAIgASgFOgI4ASLBAgoFTGFiZWwSPAoMc3RyaW5nX2xhYmVsGAEgASgLMiQudG9vbHMuZ2VuZXJhdG9yLnYxYWxwaGEuU3RyaW5nTGFiZWxIABI8Cgx1aW50NjRfbGFiZWwYAiABKAsyJC50b29scy5nZW5lcmF0b3IudjFhbHBoYS5VaW50NjRMYWJlbEgAEjoKC2ludDY0X2xhYmVsGAMgASgLMiMudG9vbHMuZ2VuZXJhdG9yLnYxYWxwaGEuSW50NjRMYWJlbEgAEjwKDHVpbnQzMl9sYWJlbBgEIAEoCzIkLnRvb2xzLmdlbmVyYXRvci52MWFscGhhLlVpbnQzMkxhYmVsSAASOgoLaW50MzJfbGFiZWwYBSABKAsyIy50b29scy5nZW5lcmF0b3IudjFhbHBoYS5JbnQzMkxhYmVsSABCBgoEa2luZCLkAQoSQ2FwYWJpbGl0eU1ldGFkYXRhEh8KBG1vZGUYASABKA4yES5zZGsudjFhbHBoYS5Nb2RlEhUKDWNhcGFiaWxpdHlfaWQYAiABKAkSRwoGbGFiZWxzGAMgAygLMjcudG9vbHMuZ2VuZXJhdG9yLnYxYWxwaGEuQ2FwYWJpbGl0eU1ldGFkYXRhLkxhYmVsc0VudHJ5Gk0KC0xhYmVsc0VudHJ5EgsKA2tleRgBIAEoCRItCgV2YWx1ZRgCIAEoCzIeLnRvb2xzLmdlbmVyYXRvci52MWFscGhhLkxhYmVsOgI4ASI2ChhDYXBhYmlsaXR5TWV0aG9kTWV0YWRhdGESGgoSbWFwX3RvX3VudHlwZWRfYXBpGAEgASgIOm4KCmNhcGFiaWxpdHkSHy5nb29nbGUucHJvdG9idWYuU2VydmljZU9wdGlvbnMY0IYDIAEoCzIrLnRvb2xzLmdlbmVyYXRvci52MWFscGhhLkNhcGFiaWxpdHlNZXRhZGF0YVIKY2FwYWJpbGl0eTprCgZtZXRob2QSHi5nb29nbGUucHJvdG9idWYuTWV0aG9kT3B0aW9ucxjRhgMgASgLMjEudG9vbHMuZ2VuZXJhdG9yLnYxYWxwaGEuQ2FwYWJpbGl0eU1ldGhvZE1ldGFkYXRhUgZtZXRob2RCrwEKG2NvbS50b29scy5nZW5lcmF0b3IudjFhbHBoYUIQQ3JlTWV0YWRhdGFQcm90b1ABogIDVEdYqgIXVG9vbHMuR2VuZXJhdG9yLlYxYWxwaGHKAhhUb29sc1xHZW5lcmF0b3JfXFYxYWxwaGHiAiRUb29sc1xHZW5lcmF0b3JfXFYxYWxwaGFcR1BCTWV0YWRhdGHqAhlUb29sczo6R2VuZXJhdG9yOjpWMWFscGhhYgZwcm90bzM", [file_google_protobuf_descriptor, file_sdk_v1alpha_sdk]);
var file_capabilities_blockchain_evm_v1alpha_client = /* @__PURE__ */ fileDesc("CjBjYXBhYmlsaXRpZXMvYmxvY2tjaGFpbi9ldm0vdjFhbHBoYS9jbGllbnQucHJvdG8SI2NhcGFiaWxpdGllcy5ibG9ja2NoYWluLmV2bS52MWFscGhhIh0KC1RvcGljVmFsdWVzEg4KBnZhbHVlcxgBIAMoDCK4AQoXRmlsdGVyTG9nVHJpZ2dlclJlcXVlc3QSEQoJYWRkcmVzc2VzGAEgAygMEkAKBnRvcGljcxgCIAMoCzIwLmNhcGFiaWxpdGllcy5ibG9ja2NoYWluLmV2bS52MWFscGhhLlRvcGljVmFsdWVzEkgKCmNvbmZpZGVuY2UYAyABKA4yNC5jYXBhYmlsaXRpZXMuYmxvY2tjaGFpbi5ldm0udjFhbHBoYS5Db25maWRlbmNlTGV2ZWwiegoTQ2FsbENvbnRyYWN0UmVxdWVzdBI6CgRjYWxsGAEgASgLMiwuY2FwYWJpbGl0aWVzLmJsb2NrY2hhaW4uZXZtLnYxYWxwaGEuQ2FsbE1zZxInCgxibG9ja19udW1iZXIYAiABKAsyES52YWx1ZXMudjEuQmlnSW50IiEKEUNhbGxDb250cmFjdFJlcGx5EgwKBGRhdGEYASABKAwiWwoRRmlsdGVyTG9nc1JlcXVlc3QSRgoMZmlsdGVyX3F1ZXJ5GAEgASgLMjAuY2FwYWJpbGl0aWVzLmJsb2NrY2hhaW4uZXZtLnYxYWxwaGEuRmlsdGVyUXVlcnkiSQoPRmlsdGVyTG9nc1JlcGx5EjYKBGxvZ3MYASADKAsyKC5jYXBhYmlsaXRpZXMuYmxvY2tjaGFpbi5ldm0udjFhbHBoYS5Mb2cixwEKA0xvZxIPCgdhZGRyZXNzGAEgASgMEg4KBnRvcGljcxgCIAMoDBIPCgd0eF9oYXNoGAMgASgMEhIKCmJsb2NrX2hhc2gYBCABKAwSDAoEZGF0YRgFIAEoDBIRCglldmVudF9zaWcYBiABKAwSJwoMYmxvY2tfbnVtYmVyGAcgASgLMhEudmFsdWVzLnYxLkJpZ0ludBIQCgh0eF9pbmRleBgIIAEoDRINCgVpbmRleBgJIAEoDRIPCgdyZW1vdmVkGAogASgIIjEKB0NhbGxNc2cSDAoEZnJvbRgBIAEoDBIKCgJ0bxgCIAEoDBIMCgRkYXRhGAMgASgMIr0BCgtGaWx0ZXJRdWVyeRISCgpibG9ja19oYXNoGAEgASgMEiUKCmZyb21fYmxvY2sYAiABKAsyES52YWx1ZXMudjEuQmlnSW50EiMKCHRvX2Jsb2NrGAMgASgLMhEudmFsdWVzLnYxLkJpZ0ludBIRCglhZGRyZXNzZXMYBCADKAwSOwoGdG9waWNzGAUgAygLMisuY2FwYWJpbGl0aWVzLmJsb2NrY2hhaW4uZXZtLnYxYWxwaGEuVG9waWNzIhcKBlRvcGljcxINCgV0b3BpYxgBIAMoDCJMChBCYWxhbmNlQXRSZXF1ZXN0Eg8KB2FjY291bnQYASABKAwSJwoMYmxvY2tfbnVtYmVyGAIgASgLMhEudmFsdWVzLnYxLkJpZ0ludCI0Cg5CYWxhbmNlQXRSZXBseRIiCgdiYWxhbmNlGAEgASgLMhEudmFsdWVzLnYxLkJpZ0ludCJPChJFc3RpbWF0ZUdhc1JlcXVlc3QSOQoDbXNnGAEgASgLMiwuY2FwYWJpbGl0aWVzLmJsb2NrY2hhaW4uZXZtLnYxYWxwaGEuQ2FsbE1zZyIjChBFc3RpbWF0ZUdhc1JlcGx5Eg8KA2dhcxgBIAEoBEICMAAiKwobR2V0VHJhbnNhY3Rpb25CeUhhc2hSZXF1ZXN0EgwKBGhhc2gYASABKAwiYgoZR2V0VHJhbnNhY3Rpb25CeUhhc2hSZXBseRJFCgt0cmFuc2FjdGlvbhgBIAEoCzIwLmNhcGFiaWxpdGllcy5ibG9ja2NoYWluLmV2bS52MWFscGhhLlRyYW5zYWN0aW9uIqEBCgtUcmFuc2FjdGlvbhIRCgVub25jZRgBIAEoBEICMAASDwoDZ2FzGAIgASgEQgIwABIKCgJ0bxgDIAEoDBIMCgRkYXRhGAQgASgMEgwKBGhhc2gYBSABKAwSIAoFdmFsdWUYBiABKAsyES52YWx1ZXMudjEuQmlnSW50EiQKCWdhc19wcmljZRgHIAEoCzIRLnZhbHVlcy52MS5CaWdJbnQiLAocR2V0VHJhbnNhY3Rpb25SZWNlaXB0UmVxdWVzdBIMCgRoYXNoGAEgASgMIlsKGkdldFRyYW5zYWN0aW9uUmVjZWlwdFJlcGx5Ej0KB3JlY2VpcHQYASABKAsyLC5jYXBhYmlsaXRpZXMuYmxvY2tjaGFpbi5ldm0udjFhbHBoYS5SZWNlaXB0IpkCCgdSZWNlaXB0EhIKBnN0YXR1cxgBIAEoBEICMAASFAoIZ2FzX3VzZWQYAiABKARCAjAAEhQKCHR4X2luZGV4GAMgASgEQgIwABISCgpibG9ja19oYXNoGAQgASgMEjYKBGxvZ3MYBiADKAsyKC5jYXBhYmlsaXRpZXMuYmxvY2tjaGFpbi5ldm0udjFhbHBoYS5Mb2cSDwoHdHhfaGFzaBgHIAEoDBIuChNlZmZlY3RpdmVfZ2FzX3ByaWNlGAggASgLMhEudmFsdWVzLnYxLkJpZ0ludBInCgxibG9ja19udW1iZXIYCSABKAsyES52YWx1ZXMudjEuQmlnSW50EhgKEGNvbnRyYWN0X2FkZHJlc3MYCiABKAwiQAoVSGVhZGVyQnlOdW1iZXJSZXF1ZXN0EicKDGJsb2NrX251bWJlchgBIAEoCzIRLnZhbHVlcy52MS5CaWdJbnQiUgoTSGVhZGVyQnlOdW1iZXJSZXBseRI7CgZoZWFkZXIYASABKAsyKy5jYXBhYmlsaXRpZXMuYmxvY2tjaGFpbi5ldm0udjFhbHBoYS5IZWFkZXIiawoGSGVhZGVyEhUKCXRpbWVzdGFtcBgBIAEoBEICMAASJwoMYmxvY2tfbnVtYmVyGAIgASgLMhEudmFsdWVzLnYxLkJpZ0ludBIMCgRoYXNoGAMgASgMEhMKC3BhcmVudF9oYXNoGAQgASgMIlsKGlJlZ2lzdGVyTG9nVHJhY2tpbmdSZXF1ZXN0Ej0KBmZpbHRlchgBIAEoCzItLmNhcGFiaWxpdGllcy5ibG9ja2NoYWluLmV2bS52MWFscGhhLkxQRmlsdGVyIsIBCghMUEZpbHRlchIZCg1tYXhfbG9nc19rZXB0GAEgASgEQgIwABIaCg5yZXRlbnRpb25fdGltZRgCIAEoA0ICMAASGgoObG9nc19wZXJfYmxvY2sYAyABKARCAjAAEgwKBG5hbWUYBCABKAkSEQoJYWRkcmVzc2VzGAUgAygMEhIKCmV2ZW50X3NpZ3MYBiADKAwSDgoGdG9waWMyGAcgAygMEg4KBnRvcGljMxgIIAMoDBIOCgZ0b3BpYzQYCSADKAwiMwocVW5yZWdpc3RlckxvZ1RyYWNraW5nUmVxdWVzdBITCgtmaWx0ZXJfbmFtZRgBIAEoCSKrAQoSV3JpdGVSZXBvcnRSZXF1ZXN0EhAKCHJlY2VpdmVyGAEgASgMEisKBnJlcG9ydBgCIAEoCzIbLnNkay52MWFscGhhLlJlcG9ydFJlc3BvbnNlEkcKCmdhc19jb25maWcYAyABKAsyLi5jYXBhYmlsaXRpZXMuYmxvY2tjaGFpbi5ldm0udjFhbHBoYS5HYXNDb25maWdIAIgBAUINCgtfZ2FzX2NvbmZpZyIiCglHYXNDb25maWcSFQoJZ2FzX2xpbWl0GAEgASgEQgIwACKHAwoQV3JpdGVSZXBvcnRSZXBseRJACgl0eF9zdGF0dXMYASABKA4yLS5jYXBhYmlsaXRpZXMuYmxvY2tjaGFpbi5ldm0udjFhbHBoYS5UeFN0YXR1cxJ1CiJyZWNlaXZlcl9jb250cmFjdF9leGVjdXRpb25fc3RhdHVzGAIgASgOMkQuY2FwYWJpbGl0aWVzLmJsb2NrY2hhaW4uZXZtLnYxYWxwaGEuUmVjZWl2ZXJDb250cmFjdEV4ZWN1dGlvblN0YXR1c0gAiAEBEhQKB3R4X2hhc2gYAyABKAxIAYgBARIvCg90cmFuc2FjdGlvbl9mZWUYBCABKAsyES52YWx1ZXMudjEuQmlnSW50SAKIAQESGgoNZXJyb3JfbWVzc2FnZRgFIAEoCUgDiAEBQiUKI19yZWNlaXZlcl9jb250cmFjdF9leGVjdXRpb25fc3RhdHVzQgoKCF90eF9oYXNoQhIKEF90cmFuc2FjdGlvbl9mZWVCEAoOX2Vycm9yX21lc3NhZ2UqaQoPQ29uZmlkZW5jZUxldmVsEhkKFUNPTkZJREVOQ0VfTEVWRUxfU0FGRRAAEhsKF0NPTkZJREVOQ0VfTEVWRUxfTEFURVNUEAESHgoaQ09ORklERU5DRV9MRVZFTF9GSU5BTElaRUQQAiqCAQofUmVjZWl2ZXJDb250cmFjdEV4ZWN1dGlvblN0YXR1cxIuCipSRUNFSVZFUl9DT05UUkFDVF9FWEVDVVRJT05fU1RBVFVTX1NVQ0NFU1MQABIvCitSRUNFSVZFUl9DT05UUkFDVF9FWEVDVVRJT05fU1RBVFVTX1JFVkVSVEVEEAEqTgoIVHhTdGF0dXMSEwoPVFhfU1RBVFVTX0ZBVEFMEAASFgoSVFhfU1RBVFVTX1JFVkVSVEVEEAESFQoRVFhfU1RBVFVTX1NVQ0NFU1MQAjLXDwoGQ2xpZW50EoABCgxDYWxsQ29udHJhY3QSOC5jYXBhYmlsaXRpZXMuYmxvY2tjaGFpbi5ldm0udjFhbHBoYS5DYWxsQ29udHJhY3RSZXF1ZXN0GjYuY2FwYWJpbGl0aWVzLmJsb2NrY2hhaW4uZXZtLnYxYWxwaGEuQ2FsbENvbnRyYWN0UmVwbHkSegoKRmlsdGVyTG9ncxI2LmNhcGFiaWxpdGllcy5ibG9ja2NoYWluLmV2bS52MWFscGhhLkZpbHRlckxvZ3NSZXF1ZXN0GjQuY2FwYWJpbGl0aWVzLmJsb2NrY2hhaW4uZXZtLnYxYWxwaGEuRmlsdGVyTG9nc1JlcGx5EncKCUJhbGFuY2VBdBI1LmNhcGFiaWxpdGllcy5ibG9ja2NoYWluLmV2bS52MWFscGhhLkJhbGFuY2VBdFJlcXVlc3QaMy5jYXBhYmlsaXRpZXMuYmxvY2tjaGFpbi5ldm0udjFhbHBoYS5CYWxhbmNlQXRSZXBseRJ9CgtFc3RpbWF0ZUdhcxI3LmNhcGFiaWxpdGllcy5ibG9ja2NoYWluLmV2bS52MWFscGhhLkVzdGltYXRlR2FzUmVxdWVzdBo1LmNhcGFiaWxpdGllcy5ibG9ja2NoYWluLmV2bS52MWFscGhhLkVzdGltYXRlR2FzUmVwbHkSmAEKFEdldFRyYW5zYWN0aW9uQnlIYXNoEkAuY2FwYWJpbGl0aWVzLmJsb2NrY2hhaW4uZXZtLnYxYWxwaGEuR2V0VHJhbnNhY3Rpb25CeUhhc2hSZXF1ZXN0Gj4uY2FwYWJpbGl0aWVzLmJsb2NrY2hhaW4uZXZtLnYxYWxwaGEuR2V0VHJhbnNhY3Rpb25CeUhhc2hSZXBseRKbAQoVR2V0VHJhbnNhY3Rpb25SZWNlaXB0EkEuY2FwYWJpbGl0aWVzLmJsb2NrY2hhaW4uZXZtLnYxYWxwaGEuR2V0VHJhbnNhY3Rpb25SZWNlaXB0UmVxdWVzdBo/LmNhcGFiaWxpdGllcy5ibG9ja2NoYWluLmV2bS52MWFscGhhLkdldFRyYW5zYWN0aW9uUmVjZWlwdFJlcGx5EoYBCg5IZWFkZXJCeU51bWJlchI6LmNhcGFiaWxpdGllcy5ibG9ja2NoYWluLmV2bS52MWFscGhhLkhlYWRlckJ5TnVtYmVyUmVxdWVzdBo4LmNhcGFiaWxpdGllcy5ibG9ja2NoYWluLmV2bS52MWFscGhhLkhlYWRlckJ5TnVtYmVyUmVwbHkSbgoTUmVnaXN0ZXJMb2dUcmFja2luZxI/LmNhcGFiaWxpdGllcy5ibG9ja2NoYWluLmV2bS52MWFscGhhLlJlZ2lzdGVyTG9nVHJhY2tpbmdSZXF1ZXN0GhYuZ29vZ2xlLnByb3RvYnVmLkVtcHR5EnIKFVVucmVnaXN0ZXJMb2dUcmFja2luZxJBLmNhcGFiaWxpdGllcy5ibG9ja2NoYWluLmV2bS52MWFscGhhLlVucmVnaXN0ZXJMb2dUcmFja2luZ1JlcXVlc3QaFi5nb29nbGUucHJvdG9idWYuRW1wdHkSdgoKTG9nVHJpZ2dlchI8LmNhcGFiaWxpdGllcy5ibG9ja2NoYWluLmV2bS52MWFscGhhLkZpbHRlckxvZ1RyaWdnZXJSZXF1ZXN0GiguY2FwYWJpbGl0aWVzLmJsb2NrY2hhaW4uZXZtLnYxYWxwaGEuTG9nMAESfQoLV3JpdGVSZXBvcnQSNy5jYXBhYmlsaXRpZXMuYmxvY2tjaGFpbi5ldm0udjFhbHBoYS5Xcml0ZVJlcG9ydFJlcXVlc3QaNS5jYXBhYmlsaXRpZXMuYmxvY2tjaGFpbi5ldm0udjFhbHBoYS5Xcml0ZVJlcG9ydFJlcGx5GrgEgrUYswQIARIJZXZtQDEuMC4wGqMECg1DaGFpblNlbGVjdG9yEpEEEo4ECh0KEWF2YWxhbmNoZS1tYWlubmV0ENXnisDh1ZikWQojChZhdmFsYW5jaGUtdGVzdG5ldC1mdWppEJv5/JCi46j4zAEKLwojYmluYW5jZV9zbWFydF9jaGFpbi1tYWlubmV0LW9wYm5iLTEQia2P75PG17sGCjAKI2JpbmFuY2Vfc21hcnRfY2hhaW4tdGVzdG5ldC1vcGJuYi0xEI71hZHBg4+cuAEKHAoQZXRoZXJldW0tbWFpbm5ldBCV9vHkz7KmwkUKJwobZXRoZXJldW0tbWFpbm5ldC1hcmJpdHJ1bS0xEMTojc2Om6HXRAonChtldGhlcmV1bS1tYWlubmV0LW9wdGltaXNtLTEQuJWPw/f+0OkzCiUKGGV0aGVyZXVtLXRlc3RuZXQtc2Vwb2xpYRDZteTO/MnuoN4BCi8KI2V0aGVyZXVtLXRlc3RuZXQtc2Vwb2xpYS1hcmJpdHJ1bS0xEOrO7v/qtoSjMAosCh9ldGhlcmV1bS10ZXN0bmV0LXNlcG9saWEtYmFzZS0xELjKue/2kK7IjwEKLwojZXRoZXJldW0tdGVzdG5ldC1zZXBvbGlhLW9wdGltaXNtLTEQn4bFob7Yw8BIChsKD3BvbHlnb24tbWFpbm5ldBCxq+TwmpKGnTgKIQoUcG9seWdvbi10ZXN0bmV0LWFtb3kQzY/W3/HHkPrhAULlAQonY29tLmNhcGFiaWxpdGllcy5ibG9ja2NoYWluLmV2bS52MWFscGhhQgtDbGllbnRQcm90b1ABogIDQ0JFqgIjQ2FwYWJpbGl0aWVzLkJsb2NrY2hhaW4uRXZtLlYxYWxwaGHKAiNDYXBhYmlsaXRpZXNcQmxvY2tjaGFpblxFdm1cVjFhbHBoYeICL0NhcGFiaWxpdGllc1xCbG9ja2NoYWluXEV2bVxWMWFscGhhXEdQQk1ldGFkYXRh6gImQ2FwYWJpbGl0aWVzOjpCbG9ja2NoYWluOjpFdm06OlYxYWxwaGFiBnByb3RvMw", [
  file_google_protobuf_empty,
  file_sdk_v1alpha_sdk,
  file_tools_generator_v1alpha_cre_metadata,
  file_values_v1_values
]);
var FilterLogTriggerRequestSchema = /* @__PURE__ */ messageDesc(file_capabilities_blockchain_evm_v1alpha_client, 1);
var CallContractRequestSchema = /* @__PURE__ */ messageDesc(file_capabilities_blockchain_evm_v1alpha_client, 2);
var CallContractReplySchema = /* @__PURE__ */ messageDesc(file_capabilities_blockchain_evm_v1alpha_client, 3);
var FilterLogsRequestSchema = /* @__PURE__ */ messageDesc(file_capabilities_blockchain_evm_v1alpha_client, 4);
var FilterLogsReplySchema = /* @__PURE__ */ messageDesc(file_capabilities_blockchain_evm_v1alpha_client, 5);
var LogSchema = /* @__PURE__ */ messageDesc(file_capabilities_blockchain_evm_v1alpha_client, 6);
var BalanceAtRequestSchema = /* @__PURE__ */ messageDesc(file_capabilities_blockchain_evm_v1alpha_client, 10);
var BalanceAtReplySchema = /* @__PURE__ */ messageDesc(file_capabilities_blockchain_evm_v1alpha_client, 11);
var EstimateGasRequestSchema = /* @__PURE__ */ messageDesc(file_capabilities_blockchain_evm_v1alpha_client, 12);
var EstimateGasReplySchema = /* @__PURE__ */ messageDesc(file_capabilities_blockchain_evm_v1alpha_client, 13);
var GetTransactionByHashRequestSchema = /* @__PURE__ */ messageDesc(file_capabilities_blockchain_evm_v1alpha_client, 14);
var GetTransactionByHashReplySchema = /* @__PURE__ */ messageDesc(file_capabilities_blockchain_evm_v1alpha_client, 15);
var GetTransactionReceiptRequestSchema = /* @__PURE__ */ messageDesc(file_capabilities_blockchain_evm_v1alpha_client, 17);
var GetTransactionReceiptReplySchema = /* @__PURE__ */ messageDesc(file_capabilities_blockchain_evm_v1alpha_client, 18);
var HeaderByNumberRequestSchema = /* @__PURE__ */ messageDesc(file_capabilities_blockchain_evm_v1alpha_client, 20);
var HeaderByNumberReplySchema = /* @__PURE__ */ messageDesc(file_capabilities_blockchain_evm_v1alpha_client, 21);
var RegisterLogTrackingRequestSchema = /* @__PURE__ */ messageDesc(file_capabilities_blockchain_evm_v1alpha_client, 23);
var UnregisterLogTrackingRequestSchema = /* @__PURE__ */ messageDesc(file_capabilities_blockchain_evm_v1alpha_client, 25);
var WriteReportRequestSchema = /* @__PURE__ */ messageDesc(file_capabilities_blockchain_evm_v1alpha_client, 26);
var GasConfigSchema = /* @__PURE__ */ messageDesc(file_capabilities_blockchain_evm_v1alpha_client, 27);
var WriteReportReplySchema = /* @__PURE__ */ messageDesc(file_capabilities_blockchain_evm_v1alpha_client, 28);
var ConfidenceLevel;
(function(ConfidenceLevel2) {
  ConfidenceLevel2[ConfidenceLevel2["SAFE"] = 0] = "SAFE";
  ConfidenceLevel2[ConfidenceLevel2["LATEST"] = 1] = "LATEST";
  ConfidenceLevel2[ConfidenceLevel2["FINALIZED"] = 2] = "FINALIZED";
})(ConfidenceLevel || (ConfidenceLevel = {}));
var ReceiverContractExecutionStatus;
(function(ReceiverContractExecutionStatus2) {
  ReceiverContractExecutionStatus2[ReceiverContractExecutionStatus2["SUCCESS"] = 0] = "SUCCESS";
  ReceiverContractExecutionStatus2[ReceiverContractExecutionStatus2["REVERTED"] = 1] = "REVERTED";
})(ReceiverContractExecutionStatus || (ReceiverContractExecutionStatus = {}));
var TxStatus;
(function(TxStatus2) {
  TxStatus2[TxStatus2["FATAL"] = 0] = "FATAL";
  TxStatus2[TxStatus2["REVERTED"] = 1] = "REVERTED";
  TxStatus2[TxStatus2["SUCCESS"] = 2] = "SUCCESS";
})(TxStatus || (TxStatus = {}));

class Report {
  report;
  constructor(report) {
    this.report = report.$typeName ? report : fromJson(ReportResponseSchema, report);
  }
  x_generatedCodeOnly_unwrap() {
    return this.report;
  }
}
var hexToBytes = (hexStr) => {
  if (!hexStr.startsWith("0x")) {
    throw new Error(`Invalid hex string: ${hexStr}`);
  }
  if (!/^0x[0-9a-fA-F]*$/.test(hexStr)) {
    throw new Error(`Invalid hex string: ${hexStr}`);
  }
  if ((hexStr.length - 2) % 2 !== 0) {
    throw new Error(`Hex string must have an even number of characters: ${hexStr}`);
  }
  const hex = hexStr.slice(2);
  const bytes = new Uint8Array(hex.length / 2);
  for (let i = 0;i < hex.length; i += 2) {
    bytes[i / 2] = Number.parseInt(hex.slice(i, i + 2), 16);
  }
  return bytes;
};
function createWriteCreReportRequest(input) {
  return {
    receiver: hexToBytes(input.receiver),
    report: input.report,
    gasConfig: input.gasConfig !== undefined ? fromJson(GasConfigSchema, input.gasConfig) : undefined,
    $report: true
  };
}
function x_generatedCodeOnly_unwrap_WriteCreReportRequest(input) {
  return create(WriteReportRequestSchema, {
    receiver: input.receiver,
    report: input.report !== undefined ? input.report.x_generatedCodeOnly_unwrap() : undefined,
    gasConfig: input.gasConfig
  });
}

class ClientCapability {
  chainSelector;
  static CAPABILITY_ID = "evm@1.0.0";
  static CAPABILITY_NAME = "evm";
  static CAPABILITY_VERSION = "1.0.0";
  static SUPPORTED_CHAINS = {
    "avalanche-mainnet": 6433500567565415381n,
    "avalanche-testnet-fuji": 14767482510784806043n,
    "binance_smart_chain-mainnet-opbnb-1": 465944652040885897n,
    "binance_smart_chain-testnet-opbnb-1": 13274425992935471758n,
    "ethereum-mainnet": 5009297550715157269n,
    "ethereum-mainnet-arbitrum-1": 4949039107694359620n,
    "ethereum-mainnet-optimism-1": 3734403246176062136n,
    "ethereum-testnet-sepolia": 16015286601757825753n,
    "ethereum-testnet-sepolia-arbitrum-1": 3478487238524512106n,
    "ethereum-testnet-sepolia-base-1": 10344971235874465080n,
    "ethereum-testnet-sepolia-optimism-1": 5224473277236331295n,
    "polygon-mainnet": 4051577828743386545n,
    "polygon-testnet-amoy": 16281711391670634445n
  };
  constructor(chainSelector) {
    this.chainSelector = chainSelector;
  }
  callContract(runtime, input) {
    let payload;
    if (input.$typeName) {
      payload = input;
    } else {
      payload = fromJson(CallContractRequestSchema, input);
    }
    const capabilityId = this.chainSelector ? `${ClientCapability.CAPABILITY_NAME}:ChainSelector:${this.chainSelector}@${ClientCapability.CAPABILITY_VERSION}` : ClientCapability.CAPABILITY_ID;
    const capabilityResponse = runtime.callCapability({
      capabilityId,
      method: "CallContract",
      payload,
      inputSchema: CallContractRequestSchema,
      outputSchema: CallContractReplySchema
    });
    return {
      result: () => {
        const result = capabilityResponse.result();
        return result;
      }
    };
  }
  filterLogs(runtime, input) {
    let payload;
    if (input.$typeName) {
      payload = input;
    } else {
      payload = fromJson(FilterLogsRequestSchema, input);
    }
    const capabilityId = this.chainSelector ? `${ClientCapability.CAPABILITY_NAME}:ChainSelector:${this.chainSelector}@${ClientCapability.CAPABILITY_VERSION}` : ClientCapability.CAPABILITY_ID;
    const capabilityResponse = runtime.callCapability({
      capabilityId,
      method: "FilterLogs",
      payload,
      inputSchema: FilterLogsRequestSchema,
      outputSchema: FilterLogsReplySchema
    });
    return {
      result: () => {
        const result = capabilityResponse.result();
        return result;
      }
    };
  }
  balanceAt(runtime, input) {
    let payload;
    if (input.$typeName) {
      payload = input;
    } else {
      payload = fromJson(BalanceAtRequestSchema, input);
    }
    const capabilityId = this.chainSelector ? `${ClientCapability.CAPABILITY_NAME}:ChainSelector:${this.chainSelector}@${ClientCapability.CAPABILITY_VERSION}` : ClientCapability.CAPABILITY_ID;
    const capabilityResponse = runtime.callCapability({
      capabilityId,
      method: "BalanceAt",
      payload,
      inputSchema: BalanceAtRequestSchema,
      outputSchema: BalanceAtReplySchema
    });
    return {
      result: () => {
        const result = capabilityResponse.result();
        return result;
      }
    };
  }
  estimateGas(runtime, input) {
    let payload;
    if (input.$typeName) {
      payload = input;
    } else {
      payload = fromJson(EstimateGasRequestSchema, input);
    }
    const capabilityId = this.chainSelector ? `${ClientCapability.CAPABILITY_NAME}:ChainSelector:${this.chainSelector}@${ClientCapability.CAPABILITY_VERSION}` : ClientCapability.CAPABILITY_ID;
    const capabilityResponse = runtime.callCapability({
      capabilityId,
      method: "EstimateGas",
      payload,
      inputSchema: EstimateGasRequestSchema,
      outputSchema: EstimateGasReplySchema
    });
    return {
      result: () => {
        const result = capabilityResponse.result();
        return result;
      }
    };
  }
  getTransactionByHash(runtime, input) {
    let payload;
    if (input.$typeName) {
      payload = input;
    } else {
      payload = fromJson(GetTransactionByHashRequestSchema, input);
    }
    const capabilityId = this.chainSelector ? `${ClientCapability.CAPABILITY_NAME}:ChainSelector:${this.chainSelector}@${ClientCapability.CAPABILITY_VERSION}` : ClientCapability.CAPABILITY_ID;
    const capabilityResponse = runtime.callCapability({
      capabilityId,
      method: "GetTransactionByHash",
      payload,
      inputSchema: GetTransactionByHashRequestSchema,
      outputSchema: GetTransactionByHashReplySchema
    });
    return {
      result: () => {
        const result = capabilityResponse.result();
        return result;
      }
    };
  }
  getTransactionReceipt(runtime, input) {
    let payload;
    if (input.$typeName) {
      payload = input;
    } else {
      payload = fromJson(GetTransactionReceiptRequestSchema, input);
    }
    const capabilityId = this.chainSelector ? `${ClientCapability.CAPABILITY_NAME}:ChainSelector:${this.chainSelector}@${ClientCapability.CAPABILITY_VERSION}` : ClientCapability.CAPABILITY_ID;
    const capabilityResponse = runtime.callCapability({
      capabilityId,
      method: "GetTransactionReceipt",
      payload,
      inputSchema: GetTransactionReceiptRequestSchema,
      outputSchema: GetTransactionReceiptReplySchema
    });
    return {
      result: () => {
        const result = capabilityResponse.result();
        return result;
      }
    };
  }
  headerByNumber(runtime, input) {
    let payload;
    if (input.$typeName) {
      payload = input;
    } else {
      payload = fromJson(HeaderByNumberRequestSchema, input);
    }
    const capabilityId = this.chainSelector ? `${ClientCapability.CAPABILITY_NAME}:ChainSelector:${this.chainSelector}@${ClientCapability.CAPABILITY_VERSION}` : ClientCapability.CAPABILITY_ID;
    const capabilityResponse = runtime.callCapability({
      capabilityId,
      method: "HeaderByNumber",
      payload,
      inputSchema: HeaderByNumberRequestSchema,
      outputSchema: HeaderByNumberReplySchema
    });
    return {
      result: () => {
        const result = capabilityResponse.result();
        return result;
      }
    };
  }
  registerLogTracking(runtime, input) {
    let payload;
    if (input.$typeName) {
      payload = input;
    } else {
      payload = fromJson(RegisterLogTrackingRequestSchema, input);
    }
    const capabilityId = this.chainSelector ? `${ClientCapability.CAPABILITY_NAME}:ChainSelector:${this.chainSelector}@${ClientCapability.CAPABILITY_VERSION}` : ClientCapability.CAPABILITY_ID;
    const capabilityResponse = runtime.callCapability({
      capabilityId,
      method: "RegisterLogTracking",
      payload,
      inputSchema: RegisterLogTrackingRequestSchema,
      outputSchema: EmptySchema
    });
    return {
      result: () => {
        const result = capabilityResponse.result();
        return result;
      }
    };
  }
  unregisterLogTracking(runtime, input) {
    let payload;
    if (input.$typeName) {
      payload = input;
    } else {
      payload = fromJson(UnregisterLogTrackingRequestSchema, input);
    }
    const capabilityId = this.chainSelector ? `${ClientCapability.CAPABILITY_NAME}:ChainSelector:${this.chainSelector}@${ClientCapability.CAPABILITY_VERSION}` : ClientCapability.CAPABILITY_ID;
    const capabilityResponse = runtime.callCapability({
      capabilityId,
      method: "UnregisterLogTracking",
      payload,
      inputSchema: UnregisterLogTrackingRequestSchema,
      outputSchema: EmptySchema
    });
    return {
      result: () => {
        const result = capabilityResponse.result();
        return result;
      }
    };
  }
  logTrigger(config) {
    const capabilityId = this.chainSelector ? `${ClientCapability.CAPABILITY_NAME}:ChainSelector:${this.chainSelector}@${ClientCapability.CAPABILITY_VERSION}` : ClientCapability.CAPABILITY_ID;
    return new ClientLogTrigger(config, capabilityId, "LogTrigger");
  }
  writeReport(runtime, input) {
    let payload;
    if (input.$report) {
      payload = x_generatedCodeOnly_unwrap_WriteCreReportRequest(input);
    } else {
      payload = x_generatedCodeOnly_unwrap_WriteCreReportRequest(createWriteCreReportRequest(input));
    }
    const capabilityId = this.chainSelector ? `${ClientCapability.CAPABILITY_NAME}:ChainSelector:${this.chainSelector}@${ClientCapability.CAPABILITY_VERSION}` : ClientCapability.CAPABILITY_ID;
    const capabilityResponse = runtime.callCapability({
      capabilityId,
      method: "WriteReport",
      payload,
      inputSchema: WriteReportRequestSchema,
      outputSchema: WriteReportReplySchema
    });
    return {
      result: () => {
        const result = capabilityResponse.result();
        return result;
      }
    };
  }
}

class ClientLogTrigger {
  _capabilityId;
  _method;
  config;
  constructor(config, _capabilityId, _method) {
    this._capabilityId = _capabilityId;
    this._method = _method;
    this.config = config.$typeName ? config : fromJson(FilterLogTriggerRequestSchema, config);
  }
  capabilityId() {
    return this._capabilityId;
  }
  method() {
    return this._method;
  }
  outputSchema() {
    return LogSchema;
  }
  configAsAny() {
    return anyPack(FilterLogTriggerRequestSchema, this.config);
  }
  adapt(rawOutput) {
    return rawOutput;
  }
}
var file_capabilities_networking_http_v1alpha_client = /* @__PURE__ */ fileDesc("CjFjYXBhYmlsaXRpZXMvbmV0d29ya2luZy9odHRwL3YxYWxwaGEvY2xpZW50LnByb3RvEiRjYXBhYmlsaXRpZXMubmV0d29ya2luZy5odHRwLnYxYWxwaGEiPAoNQ2FjaGVTZXR0aW5ncxIXCg9yZWFkX2Zyb21fY2FjaGUYASABKAgSEgoKbWF4X2FnZV9tcxgCIAEoBSKSAgoHUmVxdWVzdBILCgN1cmwYASABKAkSDgoGbWV0aG9kGAIgASgJEksKB2hlYWRlcnMYAyADKAsyOi5jYXBhYmlsaXRpZXMubmV0d29ya2luZy5odHRwLnYxYWxwaGEuUmVxdWVzdC5IZWFkZXJzRW50cnkSDAoEYm9keRgEIAEoDBISCgp0aW1lb3V0X21zGAUgASgFEksKDmNhY2hlX3NldHRpbmdzGAYgASgLMjMuY2FwYWJpbGl0aWVzLm5ldHdvcmtpbmcuaHR0cC52MWFscGhhLkNhY2hlU2V0dGluZ3MaLgoMSGVhZGVyc0VudHJ5EgsKA2tleRgBIAEoCRINCgV2YWx1ZRgCIAEoCToCOAEiqwEKCFJlc3BvbnNlEhMKC3N0YXR1c19jb2RlGAEgASgNEkwKB2hlYWRlcnMYAiADKAsyOy5jYXBhYmlsaXRpZXMubmV0d29ya2luZy5odHRwLnYxYWxwaGEuUmVzcG9uc2UuSGVhZGVyc0VudHJ5EgwKBGJvZHkYAyABKAwaLgoMSGVhZGVyc0VudHJ5EgsKA2tleRgBIAEoCRINCgV2YWx1ZRgCIAEoCToCOAEymAEKBkNsaWVudBJsCgtTZW5kUmVxdWVzdBItLmNhcGFiaWxpdGllcy5uZXR3b3JraW5nLmh0dHAudjFhbHBoYS5SZXF1ZXN0Gi4uY2FwYWJpbGl0aWVzLm5ldHdvcmtpbmcuaHR0cC52MWFscGhhLlJlc3BvbnNlGiCCtRgcCAISGGh0dHAtYWN0aW9uc0AxLjAuMC1hbHBoYULqAQooY29tLmNhcGFiaWxpdGllcy5uZXR3b3JraW5nLmh0dHAudjFhbHBoYUILQ2xpZW50UHJvdG9QAaICA0NOSKoCJENhcGFiaWxpdGllcy5OZXR3b3JraW5nLkh0dHAuVjFhbHBoYcoCJENhcGFiaWxpdGllc1xOZXR3b3JraW5nXEh0dHBcVjFhbHBoYeICMENhcGFiaWxpdGllc1xOZXR3b3JraW5nXEh0dHBcVjFhbHBoYVxHUEJNZXRhZGF0YeoCJ0NhcGFiaWxpdGllczo6TmV0d29ya2luZzo6SHR0cDo6VjFhbHBoYWIGcHJvdG8z", [file_tools_generator_v1alpha_cre_metadata]);
var RequestSchema = /* @__PURE__ */ messageDesc(file_capabilities_networking_http_v1alpha_client, 1);
var ResponseSchema = /* @__PURE__ */ messageDesc(file_capabilities_networking_http_v1alpha_client, 2);

class SendRequester {
  runtime;
  client;
  constructor(runtime, client) {
    this.runtime = runtime;
    this.client = client;
  }
  sendRequest(input) {
    return this.client.sendRequest(this.runtime, input);
  }
}

class ClientCapability2 {
  static CAPABILITY_ID = "http-actions@1.0.0-alpha";
  static CAPABILITY_NAME = "http-actions";
  static CAPABILITY_VERSION = "1.0.0-alpha";
  sendRequest(...args) {
    if (typeof args[1] === "function") {
      const [runtime2, fn, consensusAggregation, unwrapOptions] = args;
      return this.sendRequestSugarHelper(runtime2, fn, consensusAggregation, unwrapOptions);
    }
    const [runtime, input] = args;
    return this.sendRequestCallHelper(runtime, input);
  }
  sendRequestCallHelper(runtime, input) {
    let payload;
    if (input.$typeName) {
      payload = input;
    } else {
      payload = fromJson(RequestSchema, input);
    }
    const capabilityId = ClientCapability2.CAPABILITY_ID;
    const capabilityResponse = runtime.callCapability({
      capabilityId,
      method: "SendRequest",
      payload,
      inputSchema: RequestSchema,
      outputSchema: ResponseSchema
    });
    return {
      result: () => {
        const result = capabilityResponse.result();
        return result;
      }
    };
  }
  sendRequestSugarHelper(runtime, fn, consensusAggregation, unwrapOptions) {
    const wrappedFn = (runtime2, ...args) => {
      const sendRequester = new SendRequester(runtime2, this);
      return fn(sendRequester, ...args);
    };
    return runtime.runInNodeMode(wrappedFn, consensusAggregation, unwrapOptions);
  }
}
var file_capabilities_networking_http_v1alpha_trigger = /* @__PURE__ */ fileDesc("CjJjYXBhYmlsaXRpZXMvbmV0d29ya2luZy9odHRwL3YxYWxwaGEvdHJpZ2dlci5wcm90bxIkY2FwYWJpbGl0aWVzLm5ldHdvcmtpbmcuaHR0cC52MWFscGhhIlYKBkNvbmZpZxJMCg9hdXRob3JpemVkX2tleXMYASADKAsyMy5jYXBhYmlsaXRpZXMubmV0d29ya2luZy5odHRwLnYxYWxwaGEuQXV0aG9yaXplZEtleSJaCgdQYXlsb2FkEg0KBWlucHV0GAEgASgMEkAKA2tleRgCIAEoCzIzLmNhcGFiaWxpdGllcy5uZXR3b3JraW5nLmh0dHAudjFhbHBoYS5BdXRob3JpemVkS2V5ImAKDUF1dGhvcml6ZWRLZXkSOwoEdHlwZRgBIAEoDjItLmNhcGFiaWxpdGllcy5uZXR3b3JraW5nLmh0dHAudjFhbHBoYS5LZXlUeXBlEhIKCnB1YmxpY19rZXkYAiABKAkqOwoHS2V5VHlwZRIYChRLRVlfVFlQRV9VTlNQRUNJRklFRBAAEhYKEktFWV9UWVBFX0VDRFNBX0VWTRABMpIBCgRIVFRQEmgKB1RyaWdnZXISLC5jYXBhYmlsaXRpZXMubmV0d29ya2luZy5odHRwLnYxYWxwaGEuQ29uZmlnGi0uY2FwYWJpbGl0aWVzLm5ldHdvcmtpbmcuaHR0cC52MWFscGhhLlBheWxvYWQwARoggrUYHAgBEhhodHRwLXRyaWdnZXJAMS4wLjAtYWxwaGFC6wEKKGNvbS5jYXBhYmlsaXRpZXMubmV0d29ya2luZy5odHRwLnYxYWxwaGFCDFRyaWdnZXJQcm90b1ABogIDQ05IqgIkQ2FwYWJpbGl0aWVzLk5ldHdvcmtpbmcuSHR0cC5WMWFscGhhygIkQ2FwYWJpbGl0aWVzXE5ldHdvcmtpbmdcSHR0cFxWMWFscGhh4gIwQ2FwYWJpbGl0aWVzXE5ldHdvcmtpbmdcSHR0cFxWMWFscGhhXEdQQk1ldGFkYXRh6gInQ2FwYWJpbGl0aWVzOjpOZXR3b3JraW5nOjpIdHRwOjpWMWFscGhhYgZwcm90bzM", [file_tools_generator_v1alpha_cre_metadata]);
var ConfigSchema = /* @__PURE__ */ messageDesc(file_capabilities_networking_http_v1alpha_trigger, 0);
var PayloadSchema = /* @__PURE__ */ messageDesc(file_capabilities_networking_http_v1alpha_trigger, 1);
var KeyType;
(function(KeyType2) {
  KeyType2[KeyType2["UNSPECIFIED"] = 0] = "UNSPECIFIED";
  KeyType2[KeyType2["ECDSA_EVM"] = 1] = "ECDSA_EVM";
})(KeyType || (KeyType = {}));

class HTTPCapability {
  static CAPABILITY_ID = "http-trigger@1.0.0-alpha";
  static CAPABILITY_NAME = "http-trigger";
  static CAPABILITY_VERSION = "1.0.0-alpha";
  trigger(config) {
    const capabilityId = HTTPCapability.CAPABILITY_ID;
    return new HTTPTrigger(config, capabilityId, "Trigger");
  }
}

class HTTPTrigger {
  _capabilityId;
  _method;
  config;
  constructor(config, _capabilityId, _method) {
    this._capabilityId = _capabilityId;
    this._method = _method;
    this.config = config.$typeName ? config : fromJson(ConfigSchema, config);
  }
  capabilityId() {
    return this._capabilityId;
  }
  method() {
    return this._method;
  }
  outputSchema() {
    return PayloadSchema;
  }
  configAsAny() {
    return anyPack(ConfigSchema, this.config);
  }
  adapt(rawOutput) {
    return rawOutput;
  }
}
var file_capabilities_scheduler_cron_v1_trigger = /* @__PURE__ */ fileDesc("CixjYXBhYmlsaXRpZXMvc2NoZWR1bGVyL2Nyb24vdjEvdHJpZ2dlci5wcm90bxIeY2FwYWJpbGl0aWVzLnNjaGVkdWxlci5jcm9uLnYxIhoKBkNvbmZpZxIQCghzY2hlZHVsZRgBIAEoCSJHCgdQYXlsb2FkEjwKGHNjaGVkdWxlZF9leGVjdXRpb25fdGltZRgBIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXAiNQoNTGVnYWN5UGF5bG9hZBIgChhzY2hlZHVsZWRfZXhlY3V0aW9uX3RpbWUYASABKAk6AhgBMvUBCgRDcm9uElwKB1RyaWdnZXISJi5jYXBhYmlsaXRpZXMuc2NoZWR1bGVyLmNyb24udjEuQ29uZmlnGicuY2FwYWJpbGl0aWVzLnNjaGVkdWxlci5jcm9uLnYxLlBheWxvYWQwARJzCg1MZWdhY3lUcmlnZ2VyEiYuY2FwYWJpbGl0aWVzLnNjaGVkdWxlci5jcm9uLnYxLkNvbmZpZxotLmNhcGFiaWxpdGllcy5zY2hlZHVsZXIuY3Jvbi52MS5MZWdhY3lQYXlsb2FkIgmIAgGKtRgCCAEwARoagrUYFggBEhJjcm9uLXRyaWdnZXJAMS4wLjBCzQEKImNvbS5jYXBhYmlsaXRpZXMuc2NoZWR1bGVyLmNyb24udjFCDFRyaWdnZXJQcm90b1ABogIDQ1NDqgIeQ2FwYWJpbGl0aWVzLlNjaGVkdWxlci5Dcm9uLlYxygIeQ2FwYWJpbGl0aWVzXFNjaGVkdWxlclxDcm9uXFYx4gIqQ2FwYWJpbGl0aWVzXFNjaGVkdWxlclxDcm9uXFYxXEdQQk1ldGFkYXRh6gIhQ2FwYWJpbGl0aWVzOjpTY2hlZHVsZXI6OkNyb246OlYxYgZwcm90bzM", [file_google_protobuf_timestamp, file_tools_generator_v1alpha_cre_metadata]);
var ConfigSchema2 = /* @__PURE__ */ messageDesc(file_capabilities_scheduler_cron_v1_trigger, 0);
var PayloadSchema2 = /* @__PURE__ */ messageDesc(file_capabilities_scheduler_cron_v1_trigger, 1);

class CronCapability {
  static CAPABILITY_ID = "cron-trigger@1.0.0";
  static CAPABILITY_NAME = "cron-trigger";
  static CAPABILITY_VERSION = "1.0.0";
  trigger(config) {
    const capabilityId = CronCapability.CAPABILITY_ID;
    return new CronTrigger(config, capabilityId, "Trigger");
  }
}

class CronTrigger {
  _capabilityId;
  _method;
  config;
  constructor(config, _capabilityId, _method) {
    this._capabilityId = _capabilityId;
    this._method = _method;
    this.config = config.$typeName ? config : fromJson(ConfigSchema2, config);
  }
  capabilityId() {
    return this._capabilityId;
  }
  method() {
    return this._method;
  }
  outputSchema() {
    return PayloadSchema2;
  }
  configAsAny() {
    return anyPack(ConfigSchema2, this.config);
  }
  adapt(rawOutput) {
    return rawOutput;
  }
}
var lookup = [];
var revLookup = [];
var code = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
for (i = 0, len = code.length;i < len; ++i)
  lookup[i] = code[i], revLookup[code.charCodeAt(i)] = i;
var i;
var len;
revLookup[45] = 62;
revLookup[95] = 63;
function getLens(b64) {
  var len2 = b64.length;
  if (len2 % 4 > 0)
    throw new Error("Invalid string. Length must be a multiple of 4");
  var validLen = b64.indexOf("=");
  if (validLen === -1)
    validLen = len2;
  var placeHoldersLen = validLen === len2 ? 0 : 4 - validLen % 4;
  return [validLen, placeHoldersLen];
}
function _byteLength(validLen, placeHoldersLen) {
  return (validLen + placeHoldersLen) * 3 / 4 - placeHoldersLen;
}
function toByteArray(b64) {
  var tmp, lens = getLens(b64), validLen = lens[0], placeHoldersLen = lens[1], arr = new Uint8Array(_byteLength(validLen, placeHoldersLen)), curByte = 0, len2 = placeHoldersLen > 0 ? validLen - 4 : validLen, i2;
  for (i2 = 0;i2 < len2; i2 += 4)
    tmp = revLookup[b64.charCodeAt(i2)] << 18 | revLookup[b64.charCodeAt(i2 + 1)] << 12 | revLookup[b64.charCodeAt(i2 + 2)] << 6 | revLookup[b64.charCodeAt(i2 + 3)], arr[curByte++] = tmp >> 16 & 255, arr[curByte++] = tmp >> 8 & 255, arr[curByte++] = tmp & 255;
  if (placeHoldersLen === 2)
    tmp = revLookup[b64.charCodeAt(i2)] << 2 | revLookup[b64.charCodeAt(i2 + 1)] >> 4, arr[curByte++] = tmp & 255;
  if (placeHoldersLen === 1)
    tmp = revLookup[b64.charCodeAt(i2)] << 10 | revLookup[b64.charCodeAt(i2 + 1)] << 4 | revLookup[b64.charCodeAt(i2 + 2)] >> 2, arr[curByte++] = tmp >> 8 & 255, arr[curByte++] = tmp & 255;
  return arr;
}
function tripletToBase64(num) {
  return lookup[num >> 18 & 63] + lookup[num >> 12 & 63] + lookup[num >> 6 & 63] + lookup[num & 63];
}
function encodeChunk(uint8, start, end) {
  var tmp, output = [];
  for (var i2 = start;i2 < end; i2 += 3)
    tmp = (uint8[i2] << 16 & 16711680) + (uint8[i2 + 1] << 8 & 65280) + (uint8[i2 + 2] & 255), output.push(tripletToBase64(tmp));
  return output.join("");
}
function fromByteArray(uint8) {
  var tmp, len2 = uint8.length, extraBytes = len2 % 3, parts = [], maxChunkLength = 16383;
  for (var i2 = 0, len22 = len2 - extraBytes;i2 < len22; i2 += maxChunkLength)
    parts.push(encodeChunk(uint8, i2, i2 + maxChunkLength > len22 ? len22 : i2 + maxChunkLength));
  if (extraBytes === 1)
    tmp = uint8[len2 - 1], parts.push(lookup[tmp >> 2] + lookup[tmp << 4 & 63] + "==");
  else if (extraBytes === 2)
    tmp = (uint8[len2 - 2] << 8) + uint8[len2 - 1], parts.push(lookup[tmp >> 10] + lookup[tmp >> 4 & 63] + lookup[tmp << 2 & 63] + "=");
  return parts.join("");
}
function read(buffer, offset, isLE, mLen, nBytes) {
  var e, m, eLen = nBytes * 8 - mLen - 1, eMax = (1 << eLen) - 1, eBias = eMax >> 1, nBits = -7, i2 = isLE ? nBytes - 1 : 0, d = isLE ? -1 : 1, s = buffer[offset + i2];
  i2 += d, e = s & (1 << -nBits) - 1, s >>= -nBits, nBits += eLen;
  for (;nBits > 0; e = e * 256 + buffer[offset + i2], i2 += d, nBits -= 8)
    ;
  m = e & (1 << -nBits) - 1, e >>= -nBits, nBits += mLen;
  for (;nBits > 0; m = m * 256 + buffer[offset + i2], i2 += d, nBits -= 8)
    ;
  if (e === 0)
    e = 1 - eBias;
  else if (e === eMax)
    return m ? NaN : (s ? -1 : 1) * (1 / 0);
  else
    m = m + Math.pow(2, mLen), e = e - eBias;
  return (s ? -1 : 1) * m * Math.pow(2, e - mLen);
}
function write(buffer, value, offset, isLE, mLen, nBytes) {
  var e, m, c, eLen = nBytes * 8 - mLen - 1, eMax = (1 << eLen) - 1, eBias = eMax >> 1, rt = mLen === 23 ? Math.pow(2, -24) - Math.pow(2, -77) : 0, i2 = isLE ? 0 : nBytes - 1, d = isLE ? 1 : -1, s = value < 0 || value === 0 && 1 / value < 0 ? 1 : 0;
  if (value = Math.abs(value), isNaN(value) || value === 1 / 0)
    m = isNaN(value) ? 1 : 0, e = eMax;
  else {
    if (e = Math.floor(Math.log(value) / Math.LN2), value * (c = Math.pow(2, -e)) < 1)
      e--, c *= 2;
    if (e + eBias >= 1)
      value += rt / c;
    else
      value += rt * Math.pow(2, 1 - eBias);
    if (value * c >= 2)
      e++, c /= 2;
    if (e + eBias >= eMax)
      m = 0, e = eMax;
    else if (e + eBias >= 1)
      m = (value * c - 1) * Math.pow(2, mLen), e = e + eBias;
    else
      m = value * Math.pow(2, eBias - 1) * Math.pow(2, mLen), e = 0;
  }
  for (;mLen >= 8; buffer[offset + i2] = m & 255, i2 += d, m /= 256, mLen -= 8)
    ;
  e = e << mLen | m, eLen += mLen;
  for (;eLen > 0; buffer[offset + i2] = e & 255, i2 += d, e /= 256, eLen -= 8)
    ;
  buffer[offset + i2 - d] |= s * 128;
}
var customInspectSymbol = typeof Symbol === "function" && typeof Symbol.for === "function" ? Symbol.for("nodejs.util.inspect.custom") : null;
var INSPECT_MAX_BYTES = 50;
var kMaxLength = 2147483647;
var btoa = globalThis.btoa;
var atob2 = globalThis.atob;
var File = globalThis.File;
var Blob = globalThis.Blob;
function createBuffer(length) {
  if (length > kMaxLength)
    throw new RangeError('The value "' + length + '" is invalid for option "size"');
  let buf = new Uint8Array(length);
  return Object.setPrototypeOf(buf, Buffer2.prototype), buf;
}
function E(sym, getMessage, Base) {
  return class NodeError extends Base {
    constructor() {
      super();
      Object.defineProperty(this, "message", { value: getMessage.apply(this, arguments), writable: true, configurable: true }), this.name = `${this.name} [${sym}]`, this.stack, delete this.name;
    }
    get code() {
      return sym;
    }
    set code(value) {
      Object.defineProperty(this, "code", { configurable: true, enumerable: true, value, writable: true });
    }
    toString() {
      return `${this.name} [${sym}]: ${this.message}`;
    }
  };
}
var ERR_BUFFER_OUT_OF_BOUNDS = E("ERR_BUFFER_OUT_OF_BOUNDS", function(name) {
  if (name)
    return `${name} is outside of buffer bounds`;
  return "Attempt to access memory outside buffer bounds";
}, RangeError);
var ERR_INVALID_ARG_TYPE = E("ERR_INVALID_ARG_TYPE", function(name, actual) {
  return `The "${name}" argument must be of type number. Received type ${typeof actual}`;
}, TypeError);
var ERR_OUT_OF_RANGE = E("ERR_OUT_OF_RANGE", function(str, range, input) {
  let msg = `The value of "${str}" is out of range.`, received = input;
  if (Number.isInteger(input) && Math.abs(input) > 4294967296)
    received = addNumericalSeparator(String(input));
  else if (typeof input === "bigint") {
    if (received = String(input), input > BigInt(2) ** BigInt(32) || input < -(BigInt(2) ** BigInt(32)))
      received = addNumericalSeparator(received);
    received += "n";
  }
  return msg += ` It must be ${range}. Received ${received}`, msg;
}, RangeError);
function Buffer2(arg, encodingOrOffset, length) {
  if (typeof arg === "number") {
    if (typeof encodingOrOffset === "string")
      throw new TypeError('The "string" argument must be of type string. Received type number');
    return allocUnsafe(arg);
  }
  return from(arg, encodingOrOffset, length);
}
Object.defineProperty(Buffer2.prototype, "parent", { enumerable: true, get: function() {
  if (!Buffer2.isBuffer(this))
    return;
  return this.buffer;
} });
Object.defineProperty(Buffer2.prototype, "offset", { enumerable: true, get: function() {
  if (!Buffer2.isBuffer(this))
    return;
  return this.byteOffset;
} });
Buffer2.poolSize = 8192;
function from(value, encodingOrOffset, length) {
  if (typeof value === "string")
    return fromString(value, encodingOrOffset);
  if (ArrayBuffer.isView(value))
    return fromArrayView(value);
  if (value == null)
    throw new TypeError("The first argument must be one of type string, Buffer, ArrayBuffer, Array, or Array-like Object. Received type " + typeof value);
  if (isInstance(value, ArrayBuffer) || value && isInstance(value.buffer, ArrayBuffer))
    return fromArrayBuffer(value, encodingOrOffset, length);
  if (typeof SharedArrayBuffer !== "undefined" && (isInstance(value, SharedArrayBuffer) || value && isInstance(value.buffer, SharedArrayBuffer)))
    return fromArrayBuffer(value, encodingOrOffset, length);
  if (typeof value === "number")
    throw new TypeError('The "value" argument must not be of type number. Received type number');
  let valueOf = value.valueOf && value.valueOf();
  if (valueOf != null && valueOf !== value)
    return Buffer2.from(valueOf, encodingOrOffset, length);
  let b = fromObject(value);
  if (b)
    return b;
  if (typeof Symbol !== "undefined" && Symbol.toPrimitive != null && typeof value[Symbol.toPrimitive] === "function")
    return Buffer2.from(value[Symbol.toPrimitive]("string"), encodingOrOffset, length);
  throw new TypeError("The first argument must be one of type string, Buffer, ArrayBuffer, Array, or Array-like Object. Received type " + typeof value);
}
Buffer2.from = function(value, encodingOrOffset, length) {
  return from(value, encodingOrOffset, length);
};
Object.setPrototypeOf(Buffer2.prototype, Uint8Array.prototype);
Object.setPrototypeOf(Buffer2, Uint8Array);
function assertSize(size) {
  if (typeof size !== "number")
    throw new TypeError('"size" argument must be of type number');
  else if (size < 0)
    throw new RangeError('The value "' + size + '" is invalid for option "size"');
}
function alloc(size, fill, encoding) {
  if (assertSize(size), size <= 0)
    return createBuffer(size);
  if (fill !== undefined)
    return typeof encoding === "string" ? createBuffer(size).fill(fill, encoding) : createBuffer(size).fill(fill);
  return createBuffer(size);
}
Buffer2.alloc = function(size, fill, encoding) {
  return alloc(size, fill, encoding);
};
function allocUnsafe(size) {
  return assertSize(size), createBuffer(size < 0 ? 0 : checked(size) | 0);
}
Buffer2.allocUnsafe = function(size) {
  return allocUnsafe(size);
};
Buffer2.allocUnsafeSlow = function(size) {
  return allocUnsafe(size);
};
function fromString(string, encoding) {
  if (typeof encoding !== "string" || encoding === "")
    encoding = "utf8";
  if (!Buffer2.isEncoding(encoding))
    throw new TypeError("Unknown encoding: " + encoding);
  let length = byteLength(string, encoding) | 0, buf = createBuffer(length), actual = buf.write(string, encoding);
  if (actual !== length)
    buf = buf.slice(0, actual);
  return buf;
}
function fromArrayLike(array) {
  let length = array.length < 0 ? 0 : checked(array.length) | 0, buf = createBuffer(length);
  for (let i2 = 0;i2 < length; i2 += 1)
    buf[i2] = array[i2] & 255;
  return buf;
}
function fromArrayView(arrayView) {
  if (isInstance(arrayView, Uint8Array)) {
    let copy = new Uint8Array(arrayView);
    return fromArrayBuffer(copy.buffer, copy.byteOffset, copy.byteLength);
  }
  return fromArrayLike(arrayView);
}
function fromArrayBuffer(array, byteOffset, length) {
  if (byteOffset < 0 || array.byteLength < byteOffset)
    throw new RangeError('"offset" is outside of buffer bounds');
  if (array.byteLength < byteOffset + (length || 0))
    throw new RangeError('"length" is outside of buffer bounds');
  let buf;
  if (byteOffset === undefined && length === undefined)
    buf = new Uint8Array(array);
  else if (length === undefined)
    buf = new Uint8Array(array, byteOffset);
  else
    buf = new Uint8Array(array, byteOffset, length);
  return Object.setPrototypeOf(buf, Buffer2.prototype), buf;
}
function fromObject(obj) {
  if (Buffer2.isBuffer(obj)) {
    let len2 = checked(obj.length) | 0, buf = createBuffer(len2);
    if (buf.length === 0)
      return buf;
    return obj.copy(buf, 0, 0, len2), buf;
  }
  if (obj.length !== undefined) {
    if (typeof obj.length !== "number" || Number.isNaN(obj.length))
      return createBuffer(0);
    return fromArrayLike(obj);
  }
  if (obj.type === "Buffer" && Array.isArray(obj.data))
    return fromArrayLike(obj.data);
}
function checked(length) {
  if (length >= kMaxLength)
    throw new RangeError("Attempt to allocate Buffer larger than maximum size: 0x" + kMaxLength.toString(16) + " bytes");
  return length | 0;
}
Buffer2.isBuffer = function isBuffer(b) {
  return b != null && b._isBuffer === true && b !== Buffer2.prototype;
};
Buffer2.compare = function compare(a, b) {
  if (isInstance(a, Uint8Array))
    a = Buffer2.from(a, a.offset, a.byteLength);
  if (isInstance(b, Uint8Array))
    b = Buffer2.from(b, b.offset, b.byteLength);
  if (!Buffer2.isBuffer(a) || !Buffer2.isBuffer(b))
    throw new TypeError('The "buf1", "buf2" arguments must be one of type Buffer or Uint8Array');
  if (a === b)
    return 0;
  let x = a.length, y = b.length;
  for (let i2 = 0, len2 = Math.min(x, y);i2 < len2; ++i2)
    if (a[i2] !== b[i2]) {
      x = a[i2], y = b[i2];
      break;
    }
  if (x < y)
    return -1;
  if (y < x)
    return 1;
  return 0;
};
Buffer2.isEncoding = function isEncoding(encoding) {
  switch (String(encoding).toLowerCase()) {
    case "hex":
    case "utf8":
    case "utf-8":
    case "ascii":
    case "latin1":
    case "binary":
    case "base64":
    case "ucs2":
    case "ucs-2":
    case "utf16le":
    case "utf-16le":
      return true;
    default:
      return false;
  }
};
Buffer2.concat = function concat(list, length) {
  if (!Array.isArray(list))
    throw new TypeError('"list" argument must be an Array of Buffers');
  if (list.length === 0)
    return Buffer2.alloc(0);
  let i2;
  if (length === undefined) {
    length = 0;
    for (i2 = 0;i2 < list.length; ++i2)
      length += list[i2].length;
  }
  let buffer = Buffer2.allocUnsafe(length), pos = 0;
  for (i2 = 0;i2 < list.length; ++i2) {
    let buf = list[i2];
    if (isInstance(buf, Uint8Array))
      if (pos + buf.length > buffer.length) {
        if (!Buffer2.isBuffer(buf))
          buf = Buffer2.from(buf);
        buf.copy(buffer, pos);
      } else
        Uint8Array.prototype.set.call(buffer, buf, pos);
    else if (!Buffer2.isBuffer(buf))
      throw new TypeError('"list" argument must be an Array of Buffers');
    else
      buf.copy(buffer, pos);
    pos += buf.length;
  }
  return buffer;
};
function byteLength(string, encoding) {
  if (Buffer2.isBuffer(string))
    return string.length;
  if (ArrayBuffer.isView(string) || isInstance(string, ArrayBuffer))
    return string.byteLength;
  if (typeof string !== "string")
    throw new TypeError('The "string" argument must be one of type string, Buffer, or ArrayBuffer. Received type ' + typeof string);
  let len2 = string.length, mustMatch = arguments.length > 2 && arguments[2] === true;
  if (!mustMatch && len2 === 0)
    return 0;
  let loweredCase = false;
  for (;; )
    switch (encoding) {
      case "ascii":
      case "latin1":
      case "binary":
        return len2;
      case "utf8":
      case "utf-8":
        return utf8ToBytes(string).length;
      case "ucs2":
      case "ucs-2":
      case "utf16le":
      case "utf-16le":
        return len2 * 2;
      case "hex":
        return len2 >>> 1;
      case "base64":
        return base64ToBytes(string).length;
      default:
        if (loweredCase)
          return mustMatch ? -1 : utf8ToBytes(string).length;
        encoding = ("" + encoding).toLowerCase(), loweredCase = true;
    }
}
Buffer2.byteLength = byteLength;
function slowToString(encoding, start, end) {
  let loweredCase = false;
  if (start === undefined || start < 0)
    start = 0;
  if (start > this.length)
    return "";
  if (end === undefined || end > this.length)
    end = this.length;
  if (end <= 0)
    return "";
  if (end >>>= 0, start >>>= 0, end <= start)
    return "";
  if (!encoding)
    encoding = "utf8";
  while (true)
    switch (encoding) {
      case "hex":
        return hexSlice(this, start, end);
      case "utf8":
      case "utf-8":
        return utf8Slice(this, start, end);
      case "ascii":
        return asciiSlice(this, start, end);
      case "latin1":
      case "binary":
        return latin1Slice(this, start, end);
      case "base64":
        return base64Slice(this, start, end);
      case "ucs2":
      case "ucs-2":
      case "utf16le":
      case "utf-16le":
        return utf16leSlice(this, start, end);
      default:
        if (loweredCase)
          throw new TypeError("Unknown encoding: " + encoding);
        encoding = (encoding + "").toLowerCase(), loweredCase = true;
    }
}
Buffer2.prototype._isBuffer = true;
function swap(b, n, m) {
  let i2 = b[n];
  b[n] = b[m], b[m] = i2;
}
Buffer2.prototype.swap16 = function swap16() {
  let len2 = this.length;
  if (len2 % 2 !== 0)
    throw new RangeError("Buffer size must be a multiple of 16-bits");
  for (let i2 = 0;i2 < len2; i2 += 2)
    swap(this, i2, i2 + 1);
  return this;
};
Buffer2.prototype.swap32 = function swap32() {
  let len2 = this.length;
  if (len2 % 4 !== 0)
    throw new RangeError("Buffer size must be a multiple of 32-bits");
  for (let i2 = 0;i2 < len2; i2 += 4)
    swap(this, i2, i2 + 3), swap(this, i2 + 1, i2 + 2);
  return this;
};
Buffer2.prototype.swap64 = function swap64() {
  let len2 = this.length;
  if (len2 % 8 !== 0)
    throw new RangeError("Buffer size must be a multiple of 64-bits");
  for (let i2 = 0;i2 < len2; i2 += 8)
    swap(this, i2, i2 + 7), swap(this, i2 + 1, i2 + 6), swap(this, i2 + 2, i2 + 5), swap(this, i2 + 3, i2 + 4);
  return this;
};
Buffer2.prototype.toString = function toString() {
  let length = this.length;
  if (length === 0)
    return "";
  if (arguments.length === 0)
    return utf8Slice(this, 0, length);
  return slowToString.apply(this, arguments);
};
Buffer2.prototype.toLocaleString = Buffer2.prototype.toString;
Buffer2.prototype.equals = function equals2(b) {
  if (!Buffer2.isBuffer(b))
    throw new TypeError("Argument must be a Buffer");
  if (this === b)
    return true;
  return Buffer2.compare(this, b) === 0;
};
Buffer2.prototype.inspect = function inspect() {
  let str = "", max = INSPECT_MAX_BYTES;
  if (str = this.toString("hex", 0, max).replace(/(.{2})/g, "$1 ").trim(), this.length > max)
    str += " ... ";
  return "<Buffer " + str + ">";
};
if (customInspectSymbol)
  Buffer2.prototype[customInspectSymbol] = Buffer2.prototype.inspect;
Buffer2.prototype.compare = function compare2(target, start, end, thisStart, thisEnd) {
  if (isInstance(target, Uint8Array))
    target = Buffer2.from(target, target.offset, target.byteLength);
  if (!Buffer2.isBuffer(target))
    throw new TypeError('The "target" argument must be one of type Buffer or Uint8Array. Received type ' + typeof target);
  if (start === undefined)
    start = 0;
  if (end === undefined)
    end = target ? target.length : 0;
  if (thisStart === undefined)
    thisStart = 0;
  if (thisEnd === undefined)
    thisEnd = this.length;
  if (start < 0 || end > target.length || thisStart < 0 || thisEnd > this.length)
    throw new RangeError("out of range index");
  if (thisStart >= thisEnd && start >= end)
    return 0;
  if (thisStart >= thisEnd)
    return -1;
  if (start >= end)
    return 1;
  if (start >>>= 0, end >>>= 0, thisStart >>>= 0, thisEnd >>>= 0, this === target)
    return 0;
  let x = thisEnd - thisStart, y = end - start, len2 = Math.min(x, y), thisCopy = this.slice(thisStart, thisEnd), targetCopy = target.slice(start, end);
  for (let i2 = 0;i2 < len2; ++i2)
    if (thisCopy[i2] !== targetCopy[i2]) {
      x = thisCopy[i2], y = targetCopy[i2];
      break;
    }
  if (x < y)
    return -1;
  if (y < x)
    return 1;
  return 0;
};
function bidirectionalIndexOf(buffer, val, byteOffset, encoding, dir) {
  if (buffer.length === 0)
    return -1;
  if (typeof byteOffset === "string")
    encoding = byteOffset, byteOffset = 0;
  else if (byteOffset > 2147483647)
    byteOffset = 2147483647;
  else if (byteOffset < -2147483648)
    byteOffset = -2147483648;
  if (byteOffset = +byteOffset, Number.isNaN(byteOffset))
    byteOffset = dir ? 0 : buffer.length - 1;
  if (byteOffset < 0)
    byteOffset = buffer.length + byteOffset;
  if (byteOffset >= buffer.length)
    if (dir)
      return -1;
    else
      byteOffset = buffer.length - 1;
  else if (byteOffset < 0)
    if (dir)
      byteOffset = 0;
    else
      return -1;
  if (typeof val === "string")
    val = Buffer2.from(val, encoding);
  if (Buffer2.isBuffer(val)) {
    if (val.length === 0)
      return -1;
    return arrayIndexOf(buffer, val, byteOffset, encoding, dir);
  } else if (typeof val === "number") {
    if (val = val & 255, typeof Uint8Array.prototype.indexOf === "function")
      if (dir)
        return Uint8Array.prototype.indexOf.call(buffer, val, byteOffset);
      else
        return Uint8Array.prototype.lastIndexOf.call(buffer, val, byteOffset);
    return arrayIndexOf(buffer, [val], byteOffset, encoding, dir);
  }
  throw new TypeError("val must be string, number or Buffer");
}
function arrayIndexOf(arr, val, byteOffset, encoding, dir) {
  let indexSize = 1, arrLength = arr.length, valLength = val.length;
  if (encoding !== undefined) {
    if (encoding = String(encoding).toLowerCase(), encoding === "ucs2" || encoding === "ucs-2" || encoding === "utf16le" || encoding === "utf-16le") {
      if (arr.length < 2 || val.length < 2)
        return -1;
      indexSize = 2, arrLength /= 2, valLength /= 2, byteOffset /= 2;
    }
  }
  function read2(buf, i3) {
    if (indexSize === 1)
      return buf[i3];
    else
      return buf.readUInt16BE(i3 * indexSize);
  }
  let i2;
  if (dir) {
    let foundIndex = -1;
    for (i2 = byteOffset;i2 < arrLength; i2++)
      if (read2(arr, i2) === read2(val, foundIndex === -1 ? 0 : i2 - foundIndex)) {
        if (foundIndex === -1)
          foundIndex = i2;
        if (i2 - foundIndex + 1 === valLength)
          return foundIndex * indexSize;
      } else {
        if (foundIndex !== -1)
          i2 -= i2 - foundIndex;
        foundIndex = -1;
      }
  } else {
    if (byteOffset + valLength > arrLength)
      byteOffset = arrLength - valLength;
    for (i2 = byteOffset;i2 >= 0; i2--) {
      let found = true;
      for (let j = 0;j < valLength; j++)
        if (read2(arr, i2 + j) !== read2(val, j)) {
          found = false;
          break;
        }
      if (found)
        return i2;
    }
  }
  return -1;
}
Buffer2.prototype.includes = function includes(val, byteOffset, encoding) {
  return this.indexOf(val, byteOffset, encoding) !== -1;
};
Buffer2.prototype.indexOf = function indexOf(val, byteOffset, encoding) {
  return bidirectionalIndexOf(this, val, byteOffset, encoding, true);
};
Buffer2.prototype.lastIndexOf = function lastIndexOf(val, byteOffset, encoding) {
  return bidirectionalIndexOf(this, val, byteOffset, encoding, false);
};
function hexWrite(buf, string, offset, length) {
  offset = Number(offset) || 0;
  let remaining = buf.length - offset;
  if (!length)
    length = remaining;
  else if (length = Number(length), length > remaining)
    length = remaining;
  let strLen = string.length;
  if (length > strLen / 2)
    length = strLen / 2;
  let i2;
  for (i2 = 0;i2 < length; ++i2) {
    let parsed = parseInt(string.substr(i2 * 2, 2), 16);
    if (Number.isNaN(parsed))
      return i2;
    buf[offset + i2] = parsed;
  }
  return i2;
}
function utf8Write(buf, string, offset, length) {
  return blitBuffer(utf8ToBytes(string, buf.length - offset), buf, offset, length);
}
function asciiWrite(buf, string, offset, length) {
  return blitBuffer(asciiToBytes(string), buf, offset, length);
}
function base64Write(buf, string, offset, length) {
  return blitBuffer(base64ToBytes(string), buf, offset, length);
}
function ucs2Write(buf, string, offset, length) {
  return blitBuffer(utf16leToBytes(string, buf.length - offset), buf, offset, length);
}
Buffer2.prototype.write = function write2(string, offset, length, encoding) {
  if (offset === undefined)
    encoding = "utf8", length = this.length, offset = 0;
  else if (length === undefined && typeof offset === "string")
    encoding = offset, length = this.length, offset = 0;
  else if (isFinite(offset))
    if (offset = offset >>> 0, isFinite(length)) {
      if (length = length >>> 0, encoding === undefined)
        encoding = "utf8";
    } else
      encoding = length, length = undefined;
  else
    throw new Error("Buffer.write(string, encoding, offset[, length]) is no longer supported");
  let remaining = this.length - offset;
  if (length === undefined || length > remaining)
    length = remaining;
  if (string.length > 0 && (length < 0 || offset < 0) || offset > this.length)
    throw new RangeError("Attempt to write outside buffer bounds");
  if (!encoding)
    encoding = "utf8";
  let loweredCase = false;
  for (;; )
    switch (encoding) {
      case "hex":
        return hexWrite(this, string, offset, length);
      case "utf8":
      case "utf-8":
        return utf8Write(this, string, offset, length);
      case "ascii":
      case "latin1":
      case "binary":
        return asciiWrite(this, string, offset, length);
      case "base64":
        return base64Write(this, string, offset, length);
      case "ucs2":
      case "ucs-2":
      case "utf16le":
      case "utf-16le":
        return ucs2Write(this, string, offset, length);
      default:
        if (loweredCase)
          throw new TypeError("Unknown encoding: " + encoding);
        encoding = ("" + encoding).toLowerCase(), loweredCase = true;
    }
};
Buffer2.prototype.toJSON = function toJSON() {
  return { type: "Buffer", data: Array.prototype.slice.call(this._arr || this, 0) };
};
function base64Slice(buf, start, end) {
  if (start === 0 && end === buf.length)
    return fromByteArray(buf);
  else
    return fromByteArray(buf.slice(start, end));
}
function utf8Slice(buf, start, end) {
  end = Math.min(buf.length, end);
  let res = [], i2 = start;
  while (i2 < end) {
    let firstByte = buf[i2], codePoint = null, bytesPerSequence = firstByte > 239 ? 4 : firstByte > 223 ? 3 : firstByte > 191 ? 2 : 1;
    if (i2 + bytesPerSequence <= end) {
      let secondByte, thirdByte, fourthByte, tempCodePoint;
      switch (bytesPerSequence) {
        case 1:
          if (firstByte < 128)
            codePoint = firstByte;
          break;
        case 2:
          if (secondByte = buf[i2 + 1], (secondByte & 192) === 128) {
            if (tempCodePoint = (firstByte & 31) << 6 | secondByte & 63, tempCodePoint > 127)
              codePoint = tempCodePoint;
          }
          break;
        case 3:
          if (secondByte = buf[i2 + 1], thirdByte = buf[i2 + 2], (secondByte & 192) === 128 && (thirdByte & 192) === 128) {
            if (tempCodePoint = (firstByte & 15) << 12 | (secondByte & 63) << 6 | thirdByte & 63, tempCodePoint > 2047 && (tempCodePoint < 55296 || tempCodePoint > 57343))
              codePoint = tempCodePoint;
          }
          break;
        case 4:
          if (secondByte = buf[i2 + 1], thirdByte = buf[i2 + 2], fourthByte = buf[i2 + 3], (secondByte & 192) === 128 && (thirdByte & 192) === 128 && (fourthByte & 192) === 128) {
            if (tempCodePoint = (firstByte & 15) << 18 | (secondByte & 63) << 12 | (thirdByte & 63) << 6 | fourthByte & 63, tempCodePoint > 65535 && tempCodePoint < 1114112)
              codePoint = tempCodePoint;
          }
      }
    }
    if (codePoint === null)
      codePoint = 65533, bytesPerSequence = 1;
    else if (codePoint > 65535)
      codePoint -= 65536, res.push(codePoint >>> 10 & 1023 | 55296), codePoint = 56320 | codePoint & 1023;
    res.push(codePoint), i2 += bytesPerSequence;
  }
  return decodeCodePointsArray(res);
}
var MAX_ARGUMENTS_LENGTH = 4096;
function decodeCodePointsArray(codePoints) {
  let len2 = codePoints.length;
  if (len2 <= MAX_ARGUMENTS_LENGTH)
    return String.fromCharCode.apply(String, codePoints);
  let res = "", i2 = 0;
  while (i2 < len2)
    res += String.fromCharCode.apply(String, codePoints.slice(i2, i2 += MAX_ARGUMENTS_LENGTH));
  return res;
}
function asciiSlice(buf, start, end) {
  let ret = "";
  end = Math.min(buf.length, end);
  for (let i2 = start;i2 < end; ++i2)
    ret += String.fromCharCode(buf[i2] & 127);
  return ret;
}
function latin1Slice(buf, start, end) {
  let ret = "";
  end = Math.min(buf.length, end);
  for (let i2 = start;i2 < end; ++i2)
    ret += String.fromCharCode(buf[i2]);
  return ret;
}
function hexSlice(buf, start, end) {
  let len2 = buf.length;
  if (!start || start < 0)
    start = 0;
  if (!end || end < 0 || end > len2)
    end = len2;
  let out = "";
  for (let i2 = start;i2 < end; ++i2)
    out += hexSliceLookupTable[buf[i2]];
  return out;
}
function utf16leSlice(buf, start, end) {
  let bytes = buf.slice(start, end), res = "";
  for (let i2 = 0;i2 < bytes.length - 1; i2 += 2)
    res += String.fromCharCode(bytes[i2] + bytes[i2 + 1] * 256);
  return res;
}
Buffer2.prototype.slice = function slice(start, end) {
  let len2 = this.length;
  if (start = ~~start, end = end === undefined ? len2 : ~~end, start < 0) {
    if (start += len2, start < 0)
      start = 0;
  } else if (start > len2)
    start = len2;
  if (end < 0) {
    if (end += len2, end < 0)
      end = 0;
  } else if (end > len2)
    end = len2;
  if (end < start)
    end = start;
  let newBuf = this.subarray(start, end);
  return Object.setPrototypeOf(newBuf, Buffer2.prototype), newBuf;
};
function checkOffset(offset, ext, length) {
  if (offset % 1 !== 0 || offset < 0)
    throw new RangeError("offset is not uint");
  if (offset + ext > length)
    throw new RangeError("Trying to access beyond buffer length");
}
Buffer2.prototype.readUintLE = Buffer2.prototype.readUIntLE = function readUIntLE(offset, byteLength2, noAssert) {
  if (offset = offset >>> 0, byteLength2 = byteLength2 >>> 0, !noAssert)
    checkOffset(offset, byteLength2, this.length);
  let val = this[offset], mul = 1, i2 = 0;
  while (++i2 < byteLength2 && (mul *= 256))
    val += this[offset + i2] * mul;
  return val;
};
Buffer2.prototype.readUintBE = Buffer2.prototype.readUIntBE = function readUIntBE(offset, byteLength2, noAssert) {
  if (offset = offset >>> 0, byteLength2 = byteLength2 >>> 0, !noAssert)
    checkOffset(offset, byteLength2, this.length);
  let val = this[offset + --byteLength2], mul = 1;
  while (byteLength2 > 0 && (mul *= 256))
    val += this[offset + --byteLength2] * mul;
  return val;
};
Buffer2.prototype.readUint8 = Buffer2.prototype.readUInt8 = function readUInt8(offset, noAssert) {
  if (offset = offset >>> 0, !noAssert)
    checkOffset(offset, 1, this.length);
  return this[offset];
};
Buffer2.prototype.readUint16LE = Buffer2.prototype.readUInt16LE = function readUInt16LE(offset, noAssert) {
  if (offset = offset >>> 0, !noAssert)
    checkOffset(offset, 2, this.length);
  return this[offset] | this[offset + 1] << 8;
};
Buffer2.prototype.readUint16BE = Buffer2.prototype.readUInt16BE = function readUInt16BE(offset, noAssert) {
  if (offset = offset >>> 0, !noAssert)
    checkOffset(offset, 2, this.length);
  return this[offset] << 8 | this[offset + 1];
};
Buffer2.prototype.readUint32LE = Buffer2.prototype.readUInt32LE = function readUInt32LE(offset, noAssert) {
  if (offset = offset >>> 0, !noAssert)
    checkOffset(offset, 4, this.length);
  return (this[offset] | this[offset + 1] << 8 | this[offset + 2] << 16) + this[offset + 3] * 16777216;
};
Buffer2.prototype.readUint32BE = Buffer2.prototype.readUInt32BE = function readUInt32BE(offset, noAssert) {
  if (offset = offset >>> 0, !noAssert)
    checkOffset(offset, 4, this.length);
  return this[offset] * 16777216 + (this[offset + 1] << 16 | this[offset + 2] << 8 | this[offset + 3]);
};
Buffer2.prototype.readBigUInt64LE = defineBigIntMethod(function readBigUInt64LE(offset) {
  offset = offset >>> 0, validateNumber(offset, "offset");
  let first = this[offset], last = this[offset + 7];
  if (first === undefined || last === undefined)
    boundsError(offset, this.length - 8);
  let lo = first + this[++offset] * 256 + this[++offset] * 65536 + this[++offset] * 16777216, hi = this[++offset] + this[++offset] * 256 + this[++offset] * 65536 + last * 16777216;
  return BigInt(lo) + (BigInt(hi) << BigInt(32));
});
Buffer2.prototype.readBigUInt64BE = defineBigIntMethod(function readBigUInt64BE(offset) {
  offset = offset >>> 0, validateNumber(offset, "offset");
  let first = this[offset], last = this[offset + 7];
  if (first === undefined || last === undefined)
    boundsError(offset, this.length - 8);
  let hi = first * 16777216 + this[++offset] * 65536 + this[++offset] * 256 + this[++offset], lo = this[++offset] * 16777216 + this[++offset] * 65536 + this[++offset] * 256 + last;
  return (BigInt(hi) << BigInt(32)) + BigInt(lo);
});
Buffer2.prototype.readIntLE = function readIntLE(offset, byteLength2, noAssert) {
  if (offset = offset >>> 0, byteLength2 = byteLength2 >>> 0, !noAssert)
    checkOffset(offset, byteLength2, this.length);
  let val = this[offset], mul = 1, i2 = 0;
  while (++i2 < byteLength2 && (mul *= 256))
    val += this[offset + i2] * mul;
  if (mul *= 128, val >= mul)
    val -= Math.pow(2, 8 * byteLength2);
  return val;
};
Buffer2.prototype.readIntBE = function readIntBE(offset, byteLength2, noAssert) {
  if (offset = offset >>> 0, byteLength2 = byteLength2 >>> 0, !noAssert)
    checkOffset(offset, byteLength2, this.length);
  let i2 = byteLength2, mul = 1, val = this[offset + --i2];
  while (i2 > 0 && (mul *= 256))
    val += this[offset + --i2] * mul;
  if (mul *= 128, val >= mul)
    val -= Math.pow(2, 8 * byteLength2);
  return val;
};
Buffer2.prototype.readInt8 = function readInt8(offset, noAssert) {
  if (offset = offset >>> 0, !noAssert)
    checkOffset(offset, 1, this.length);
  if (!(this[offset] & 128))
    return this[offset];
  return (255 - this[offset] + 1) * -1;
};
Buffer2.prototype.readInt16LE = function readInt16LE(offset, noAssert) {
  if (offset = offset >>> 0, !noAssert)
    checkOffset(offset, 2, this.length);
  let val = this[offset] | this[offset + 1] << 8;
  return val & 32768 ? val | 4294901760 : val;
};
Buffer2.prototype.readInt16BE = function readInt16BE(offset, noAssert) {
  if (offset = offset >>> 0, !noAssert)
    checkOffset(offset, 2, this.length);
  let val = this[offset + 1] | this[offset] << 8;
  return val & 32768 ? val | 4294901760 : val;
};
Buffer2.prototype.readInt32LE = function readInt32LE(offset, noAssert) {
  if (offset = offset >>> 0, !noAssert)
    checkOffset(offset, 4, this.length);
  return this[offset] | this[offset + 1] << 8 | this[offset + 2] << 16 | this[offset + 3] << 24;
};
Buffer2.prototype.readInt32BE = function readInt32BE(offset, noAssert) {
  if (offset = offset >>> 0, !noAssert)
    checkOffset(offset, 4, this.length);
  return this[offset] << 24 | this[offset + 1] << 16 | this[offset + 2] << 8 | this[offset + 3];
};
Buffer2.prototype.readBigInt64LE = defineBigIntMethod(function readBigInt64LE(offset) {
  offset = offset >>> 0, validateNumber(offset, "offset");
  let first = this[offset], last = this[offset + 7];
  if (first === undefined || last === undefined)
    boundsError(offset, this.length - 8);
  let val = this[offset + 4] + this[offset + 5] * 256 + this[offset + 6] * 65536 + (last << 24);
  return (BigInt(val) << BigInt(32)) + BigInt(first + this[++offset] * 256 + this[++offset] * 65536 + this[++offset] * 16777216);
});
Buffer2.prototype.readBigInt64BE = defineBigIntMethod(function readBigInt64BE(offset) {
  offset = offset >>> 0, validateNumber(offset, "offset");
  let first = this[offset], last = this[offset + 7];
  if (first === undefined || last === undefined)
    boundsError(offset, this.length - 8);
  let val = (first << 24) + this[++offset] * 65536 + this[++offset] * 256 + this[++offset];
  return (BigInt(val) << BigInt(32)) + BigInt(this[++offset] * 16777216 + this[++offset] * 65536 + this[++offset] * 256 + last);
});
Buffer2.prototype.readFloatLE = function readFloatLE(offset, noAssert) {
  if (offset = offset >>> 0, !noAssert)
    checkOffset(offset, 4, this.length);
  return read(this, offset, true, 23, 4);
};
Buffer2.prototype.readFloatBE = function readFloatBE(offset, noAssert) {
  if (offset = offset >>> 0, !noAssert)
    checkOffset(offset, 4, this.length);
  return read(this, offset, false, 23, 4);
};
Buffer2.prototype.readDoubleLE = function readDoubleLE(offset, noAssert) {
  if (offset = offset >>> 0, !noAssert)
    checkOffset(offset, 8, this.length);
  return read(this, offset, true, 52, 8);
};
Buffer2.prototype.readDoubleBE = function readDoubleBE(offset, noAssert) {
  if (offset = offset >>> 0, !noAssert)
    checkOffset(offset, 8, this.length);
  return read(this, offset, false, 52, 8);
};
function checkInt(buf, value, offset, ext, max, min) {
  if (!Buffer2.isBuffer(buf))
    throw new TypeError('"buffer" argument must be a Buffer instance');
  if (value > max || value < min)
    throw new RangeError('"value" argument is out of bounds');
  if (offset + ext > buf.length)
    throw new RangeError("Index out of range");
}
Buffer2.prototype.writeUintLE = Buffer2.prototype.writeUIntLE = function writeUIntLE(value, offset, byteLength2, noAssert) {
  if (value = +value, offset = offset >>> 0, byteLength2 = byteLength2 >>> 0, !noAssert) {
    let maxBytes = Math.pow(2, 8 * byteLength2) - 1;
    checkInt(this, value, offset, byteLength2, maxBytes, 0);
  }
  let mul = 1, i2 = 0;
  this[offset] = value & 255;
  while (++i2 < byteLength2 && (mul *= 256))
    this[offset + i2] = value / mul & 255;
  return offset + byteLength2;
};
Buffer2.prototype.writeUintBE = Buffer2.prototype.writeUIntBE = function writeUIntBE(value, offset, byteLength2, noAssert) {
  if (value = +value, offset = offset >>> 0, byteLength2 = byteLength2 >>> 0, !noAssert) {
    let maxBytes = Math.pow(2, 8 * byteLength2) - 1;
    checkInt(this, value, offset, byteLength2, maxBytes, 0);
  }
  let i2 = byteLength2 - 1, mul = 1;
  this[offset + i2] = value & 255;
  while (--i2 >= 0 && (mul *= 256))
    this[offset + i2] = value / mul & 255;
  return offset + byteLength2;
};
Buffer2.prototype.writeUint8 = Buffer2.prototype.writeUInt8 = function writeUInt8(value, offset, noAssert) {
  if (value = +value, offset = offset >>> 0, !noAssert)
    checkInt(this, value, offset, 1, 255, 0);
  return this[offset] = value & 255, offset + 1;
};
Buffer2.prototype.writeUint16LE = Buffer2.prototype.writeUInt16LE = function writeUInt16LE(value, offset, noAssert) {
  if (value = +value, offset = offset >>> 0, !noAssert)
    checkInt(this, value, offset, 2, 65535, 0);
  return this[offset] = value & 255, this[offset + 1] = value >>> 8, offset + 2;
};
Buffer2.prototype.writeUint16BE = Buffer2.prototype.writeUInt16BE = function writeUInt16BE(value, offset, noAssert) {
  if (value = +value, offset = offset >>> 0, !noAssert)
    checkInt(this, value, offset, 2, 65535, 0);
  return this[offset] = value >>> 8, this[offset + 1] = value & 255, offset + 2;
};
Buffer2.prototype.writeUint32LE = Buffer2.prototype.writeUInt32LE = function writeUInt32LE(value, offset, noAssert) {
  if (value = +value, offset = offset >>> 0, !noAssert)
    checkInt(this, value, offset, 4, 4294967295, 0);
  return this[offset + 3] = value >>> 24, this[offset + 2] = value >>> 16, this[offset + 1] = value >>> 8, this[offset] = value & 255, offset + 4;
};
Buffer2.prototype.writeUint32BE = Buffer2.prototype.writeUInt32BE = function writeUInt32BE(value, offset, noAssert) {
  if (value = +value, offset = offset >>> 0, !noAssert)
    checkInt(this, value, offset, 4, 4294967295, 0);
  return this[offset] = value >>> 24, this[offset + 1] = value >>> 16, this[offset + 2] = value >>> 8, this[offset + 3] = value & 255, offset + 4;
};
function wrtBigUInt64LE(buf, value, offset, min, max) {
  checkIntBI(value, min, max, buf, offset, 7);
  let lo = Number(value & BigInt(4294967295));
  buf[offset++] = lo, lo = lo >> 8, buf[offset++] = lo, lo = lo >> 8, buf[offset++] = lo, lo = lo >> 8, buf[offset++] = lo;
  let hi = Number(value >> BigInt(32) & BigInt(4294967295));
  return buf[offset++] = hi, hi = hi >> 8, buf[offset++] = hi, hi = hi >> 8, buf[offset++] = hi, hi = hi >> 8, buf[offset++] = hi, offset;
}
function wrtBigUInt64BE(buf, value, offset, min, max) {
  checkIntBI(value, min, max, buf, offset, 7);
  let lo = Number(value & BigInt(4294967295));
  buf[offset + 7] = lo, lo = lo >> 8, buf[offset + 6] = lo, lo = lo >> 8, buf[offset + 5] = lo, lo = lo >> 8, buf[offset + 4] = lo;
  let hi = Number(value >> BigInt(32) & BigInt(4294967295));
  return buf[offset + 3] = hi, hi = hi >> 8, buf[offset + 2] = hi, hi = hi >> 8, buf[offset + 1] = hi, hi = hi >> 8, buf[offset] = hi, offset + 8;
}
Buffer2.prototype.writeBigUInt64LE = defineBigIntMethod(function writeBigUInt64LE(value, offset = 0) {
  return wrtBigUInt64LE(this, value, offset, BigInt(0), BigInt("0xffffffffffffffff"));
});
Buffer2.prototype.writeBigUInt64BE = defineBigIntMethod(function writeBigUInt64BE(value, offset = 0) {
  return wrtBigUInt64BE(this, value, offset, BigInt(0), BigInt("0xffffffffffffffff"));
});
Buffer2.prototype.writeIntLE = function writeIntLE(value, offset, byteLength2, noAssert) {
  if (value = +value, offset = offset >>> 0, !noAssert) {
    let limit = Math.pow(2, 8 * byteLength2 - 1);
    checkInt(this, value, offset, byteLength2, limit - 1, -limit);
  }
  let i2 = 0, mul = 1, sub = 0;
  this[offset] = value & 255;
  while (++i2 < byteLength2 && (mul *= 256)) {
    if (value < 0 && sub === 0 && this[offset + i2 - 1] !== 0)
      sub = 1;
    this[offset + i2] = (value / mul >> 0) - sub & 255;
  }
  return offset + byteLength2;
};
Buffer2.prototype.writeIntBE = function writeIntBE(value, offset, byteLength2, noAssert) {
  if (value = +value, offset = offset >>> 0, !noAssert) {
    let limit = Math.pow(2, 8 * byteLength2 - 1);
    checkInt(this, value, offset, byteLength2, limit - 1, -limit);
  }
  let i2 = byteLength2 - 1, mul = 1, sub = 0;
  this[offset + i2] = value & 255;
  while (--i2 >= 0 && (mul *= 256)) {
    if (value < 0 && sub === 0 && this[offset + i2 + 1] !== 0)
      sub = 1;
    this[offset + i2] = (value / mul >> 0) - sub & 255;
  }
  return offset + byteLength2;
};
Buffer2.prototype.writeInt8 = function writeInt8(value, offset, noAssert) {
  if (value = +value, offset = offset >>> 0, !noAssert)
    checkInt(this, value, offset, 1, 127, -128);
  if (value < 0)
    value = 255 + value + 1;
  return this[offset] = value & 255, offset + 1;
};
Buffer2.prototype.writeInt16LE = function writeInt16LE(value, offset, noAssert) {
  if (value = +value, offset = offset >>> 0, !noAssert)
    checkInt(this, value, offset, 2, 32767, -32768);
  return this[offset] = value & 255, this[offset + 1] = value >>> 8, offset + 2;
};
Buffer2.prototype.writeInt16BE = function writeInt16BE(value, offset, noAssert) {
  if (value = +value, offset = offset >>> 0, !noAssert)
    checkInt(this, value, offset, 2, 32767, -32768);
  return this[offset] = value >>> 8, this[offset + 1] = value & 255, offset + 2;
};
Buffer2.prototype.writeInt32LE = function writeInt32LE(value, offset, noAssert) {
  if (value = +value, offset = offset >>> 0, !noAssert)
    checkInt(this, value, offset, 4, 2147483647, -2147483648);
  return this[offset] = value & 255, this[offset + 1] = value >>> 8, this[offset + 2] = value >>> 16, this[offset + 3] = value >>> 24, offset + 4;
};
Buffer2.prototype.writeInt32BE = function writeInt32BE(value, offset, noAssert) {
  if (value = +value, offset = offset >>> 0, !noAssert)
    checkInt(this, value, offset, 4, 2147483647, -2147483648);
  if (value < 0)
    value = 4294967295 + value + 1;
  return this[offset] = value >>> 24, this[offset + 1] = value >>> 16, this[offset + 2] = value >>> 8, this[offset + 3] = value & 255, offset + 4;
};
Buffer2.prototype.writeBigInt64LE = defineBigIntMethod(function writeBigInt64LE(value, offset = 0) {
  return wrtBigUInt64LE(this, value, offset, -BigInt("0x8000000000000000"), BigInt("0x7fffffffffffffff"));
});
Buffer2.prototype.writeBigInt64BE = defineBigIntMethod(function writeBigInt64BE(value, offset = 0) {
  return wrtBigUInt64BE(this, value, offset, -BigInt("0x8000000000000000"), BigInt("0x7fffffffffffffff"));
});
function checkIEEE754(buf, value, offset, ext, max, min) {
  if (offset + ext > buf.length)
    throw new RangeError("Index out of range");
  if (offset < 0)
    throw new RangeError("Index out of range");
}
function writeFloat(buf, value, offset, littleEndian, noAssert) {
  if (value = +value, offset = offset >>> 0, !noAssert)
    checkIEEE754(buf, value, offset, 4, 340282346638528860000000000000000000000, -340282346638528860000000000000000000000);
  return write(buf, value, offset, littleEndian, 23, 4), offset + 4;
}
Buffer2.prototype.writeFloatLE = function writeFloatLE(value, offset, noAssert) {
  return writeFloat(this, value, offset, true, noAssert);
};
Buffer2.prototype.writeFloatBE = function writeFloatBE(value, offset, noAssert) {
  return writeFloat(this, value, offset, false, noAssert);
};
function writeDouble(buf, value, offset, littleEndian, noAssert) {
  if (value = +value, offset = offset >>> 0, !noAssert)
    checkIEEE754(buf, value, offset, 8, 179769313486231570000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000, -179769313486231570000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000);
  return write(buf, value, offset, littleEndian, 52, 8), offset + 8;
}
Buffer2.prototype.writeDoubleLE = function writeDoubleLE(value, offset, noAssert) {
  return writeDouble(this, value, offset, true, noAssert);
};
Buffer2.prototype.writeDoubleBE = function writeDoubleBE(value, offset, noAssert) {
  return writeDouble(this, value, offset, false, noAssert);
};
Buffer2.prototype.copy = function copy(target, targetStart, start, end) {
  if (!Buffer2.isBuffer(target))
    throw new TypeError("argument should be a Buffer");
  if (!start)
    start = 0;
  if (!end && end !== 0)
    end = this.length;
  if (targetStart >= target.length)
    targetStart = target.length;
  if (!targetStart)
    targetStart = 0;
  if (end > 0 && end < start)
    end = start;
  if (end === start)
    return 0;
  if (target.length === 0 || this.length === 0)
    return 0;
  if (targetStart < 0)
    throw new RangeError("targetStart out of bounds");
  if (start < 0 || start >= this.length)
    throw new RangeError("Index out of range");
  if (end < 0)
    throw new RangeError("sourceEnd out of bounds");
  if (end > this.length)
    end = this.length;
  if (target.length - targetStart < end - start)
    end = target.length - targetStart + start;
  let len2 = end - start;
  if (this === target && typeof Uint8Array.prototype.copyWithin === "function")
    this.copyWithin(targetStart, start, end);
  else
    Uint8Array.prototype.set.call(target, this.subarray(start, end), targetStart);
  return len2;
};
Buffer2.prototype.fill = function fill(val, start, end, encoding) {
  if (typeof val === "string") {
    if (typeof start === "string")
      encoding = start, start = 0, end = this.length;
    else if (typeof end === "string")
      encoding = end, end = this.length;
    if (encoding !== undefined && typeof encoding !== "string")
      throw new TypeError("encoding must be a string");
    if (typeof encoding === "string" && !Buffer2.isEncoding(encoding))
      throw new TypeError("Unknown encoding: " + encoding);
    if (val.length === 1) {
      let code2 = val.charCodeAt(0);
      if (encoding === "utf8" && code2 < 128 || encoding === "latin1")
        val = code2;
    }
  } else if (typeof val === "number")
    val = val & 255;
  else if (typeof val === "boolean")
    val = Number(val);
  if (start < 0 || this.length < start || this.length < end)
    throw new RangeError("Out of range index");
  if (end <= start)
    return this;
  if (start = start >>> 0, end = end === undefined ? this.length : end >>> 0, !val)
    val = 0;
  let i2;
  if (typeof val === "number")
    for (i2 = start;i2 < end; ++i2)
      this[i2] = val;
  else {
    let bytes = Buffer2.isBuffer(val) ? val : Buffer2.from(val, encoding), len2 = bytes.length;
    if (len2 === 0)
      throw new TypeError('The value "' + val + '" is invalid for argument "value"');
    for (i2 = 0;i2 < end - start; ++i2)
      this[i2 + start] = bytes[i2 % len2];
  }
  return this;
};
function addNumericalSeparator(val) {
  let res = "", i2 = val.length, start = val[0] === "-" ? 1 : 0;
  for (;i2 >= start + 4; i2 -= 3)
    res = `_${val.slice(i2 - 3, i2)}${res}`;
  return `${val.slice(0, i2)}${res}`;
}
function checkBounds(buf, offset, byteLength2) {
  if (validateNumber(offset, "offset"), buf[offset] === undefined || buf[offset + byteLength2] === undefined)
    boundsError(offset, buf.length - (byteLength2 + 1));
}
function checkIntBI(value, min, max, buf, offset, byteLength2) {
  if (value > max || value < min) {
    let n = typeof min === "bigint" ? "n" : "", range;
    if (byteLength2 > 3)
      if (min === 0 || min === BigInt(0))
        range = `>= 0${n} and < 2${n} ** ${(byteLength2 + 1) * 8}${n}`;
      else
        range = `>= -(2${n} ** ${(byteLength2 + 1) * 8 - 1}${n}) and < 2 ** ${(byteLength2 + 1) * 8 - 1}${n}`;
    else
      range = `>= ${min}${n} and <= ${max}${n}`;
    throw new ERR_OUT_OF_RANGE("value", range, value);
  }
  checkBounds(buf, offset, byteLength2);
}
function validateNumber(value, name) {
  if (typeof value !== "number")
    throw new ERR_INVALID_ARG_TYPE(name, "number", value);
}
function boundsError(value, length, type) {
  if (Math.floor(value) !== value)
    throw validateNumber(value, type), new ERR_OUT_OF_RANGE(type || "offset", "an integer", value);
  if (length < 0)
    throw new ERR_BUFFER_OUT_OF_BOUNDS;
  throw new ERR_OUT_OF_RANGE(type || "offset", `>= ${type ? 1 : 0} and <= ${length}`, value);
}
var INVALID_BASE64_RE = /[^+/0-9A-Za-z-_]/g;
function base64clean(str) {
  if (str = str.split("=")[0], str = str.trim().replace(INVALID_BASE64_RE, ""), str.length < 2)
    return "";
  while (str.length % 4 !== 0)
    str = str + "=";
  return str;
}
function utf8ToBytes(string, units) {
  units = units || 1 / 0;
  let codePoint, length = string.length, leadSurrogate = null, bytes = [];
  for (let i2 = 0;i2 < length; ++i2) {
    if (codePoint = string.charCodeAt(i2), codePoint > 55295 && codePoint < 57344) {
      if (!leadSurrogate) {
        if (codePoint > 56319) {
          if ((units -= 3) > -1)
            bytes.push(239, 191, 189);
          continue;
        } else if (i2 + 1 === length) {
          if ((units -= 3) > -1)
            bytes.push(239, 191, 189);
          continue;
        }
        leadSurrogate = codePoint;
        continue;
      }
      if (codePoint < 56320) {
        if ((units -= 3) > -1)
          bytes.push(239, 191, 189);
        leadSurrogate = codePoint;
        continue;
      }
      codePoint = (leadSurrogate - 55296 << 10 | codePoint - 56320) + 65536;
    } else if (leadSurrogate) {
      if ((units -= 3) > -1)
        bytes.push(239, 191, 189);
    }
    if (leadSurrogate = null, codePoint < 128) {
      if ((units -= 1) < 0)
        break;
      bytes.push(codePoint);
    } else if (codePoint < 2048) {
      if ((units -= 2) < 0)
        break;
      bytes.push(codePoint >> 6 | 192, codePoint & 63 | 128);
    } else if (codePoint < 65536) {
      if ((units -= 3) < 0)
        break;
      bytes.push(codePoint >> 12 | 224, codePoint >> 6 & 63 | 128, codePoint & 63 | 128);
    } else if (codePoint < 1114112) {
      if ((units -= 4) < 0)
        break;
      bytes.push(codePoint >> 18 | 240, codePoint >> 12 & 63 | 128, codePoint >> 6 & 63 | 128, codePoint & 63 | 128);
    } else
      throw new Error("Invalid code point");
  }
  return bytes;
}
function asciiToBytes(str) {
  let byteArray = [];
  for (let i2 = 0;i2 < str.length; ++i2)
    byteArray.push(str.charCodeAt(i2) & 255);
  return byteArray;
}
function utf16leToBytes(str, units) {
  let c, hi, lo, byteArray = [];
  for (let i2 = 0;i2 < str.length; ++i2) {
    if ((units -= 2) < 0)
      break;
    c = str.charCodeAt(i2), hi = c >> 8, lo = c % 256, byteArray.push(lo), byteArray.push(hi);
  }
  return byteArray;
}
function base64ToBytes(str) {
  return toByteArray(base64clean(str));
}
function blitBuffer(src, dst, offset, length) {
  let i2;
  for (i2 = 0;i2 < length; ++i2) {
    if (i2 + offset >= dst.length || i2 >= src.length)
      break;
    dst[i2 + offset] = src[i2];
  }
  return i2;
}
function isInstance(obj, type) {
  return obj instanceof type || obj != null && obj.constructor != null && obj.constructor.name != null && obj.constructor.name === type.name;
}
var hexSliceLookupTable = function() {
  let table = new Array(256);
  for (let i2 = 0;i2 < 16; ++i2) {
    let i16 = i2 * 16;
    for (let j = 0;j < 16; ++j)
      table[i16 + j] = "0123456789abcdef"[i2] + "0123456789abcdef"[j];
  }
  return table;
}();
function defineBigIntMethod(fn) {
  return typeof BigInt === "undefined" ? BufferBigIntNotDefined : fn;
}
function BufferBigIntNotDefined() {
  throw new Error("BigInt not supported");
}
function notimpl(name) {
  return () => {
    throw new Error(name + " is not implemented for node:buffer browser polyfill");
  };
}
var resolveObjectURL = notimpl("resolveObjectURL");
var isUtf8 = notimpl("isUtf8");
var transcode = notimpl("transcode");
var prepareRuntime = () => {
  globalThis.Buffer = Buffer2;
};
var handler = (trigger, fn) => ({
  trigger,
  fn
});
prepareRuntime();
var cre = {
  capabilities: {
    CronCapability,
    HTTPCapability,
    HTTPClient: ClientCapability2,
    EVMClient: ClientCapability
  },
  handler
};
var LAST_FINALIZED_BLOCK_NUMBER = {
  absVal: Buffer.from([3]).toString("base64"),
  sign: "-1"
};
var LATEST_BLOCK_NUMBER = {
  absVal: Buffer.from([2]).toString("base64"),
  sign: "-1"
};
function sendReport(runtime, report, fn) {
  const rawReport = report.x_generatedCodeOnly_unwrap();
  const request = fn(rawReport);
  return this.sendRequest(runtime, request);
}
function sendRequesterSendReport(report, fn) {
  const rawReport = report.x_generatedCodeOnly_unwrap();
  const request = fn(rawReport);
  return this.sendRequest(request);
}
ClientCapability2.prototype.sendReport = sendReport;
SendRequester.prototype.sendReport = sendRequesterSendReport;
var network = {
  chainId: "1",
  chainSelector: {
    name: "aptos-mainnet",
    selector: 4741433654826277614n
  },
  chainFamily: "aptos",
  networkType: "mainnet"
};
var aptos_mainnet_default = network;
var network2 = {
  chainId: "16661",
  chainSelector: {
    name: "0g-mainnet",
    selector: 4426351306075016396n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var _0g_mainnet_default = network2;
var network3 = {
  chainId: "2741",
  chainSelector: {
    name: "abstract-mainnet",
    selector: 3577778157919314504n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var abstract_mainnet_default = network3;
var network4 = {
  chainId: "33139",
  chainSelector: {
    name: "apechain-mainnet",
    selector: 14894068710063348487n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var apechain_mainnet_default = network4;
var network5 = {
  chainId: "463",
  chainSelector: {
    name: "areon-mainnet",
    selector: 1939936305787790600n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var areon_mainnet_default = network5;
var network6 = {
  chainId: "43114",
  chainSelector: {
    name: "avalanche-mainnet",
    selector: 6433500567565415381n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var avalanche_mainnet_default = network6;
var network7 = {
  chainId: "432204",
  chainSelector: {
    name: "avalanche-subnet-dexalot-mainnet",
    selector: 5463201557265485081n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var avalanche_subnet_dexalot_mainnet_default = network7;
var network8 = {
  chainId: "80094",
  chainSelector: {
    name: "berachain-mainnet",
    selector: 1294465214383781161n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var berachain_mainnet_default = network8;
var network9 = {
  chainId: "56",
  chainSelector: {
    name: "binance_smart_chain-mainnet",
    selector: 11344663589394136015n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var binance_smart_chain_mainnet_default = network9;
var network10 = {
  chainId: "204",
  chainSelector: {
    name: "binance_smart_chain-mainnet-opbnb-1",
    selector: 465944652040885897n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var binance_smart_chain_mainnet_opbnb_1_default = network10;
var network11 = {
  chainId: "1907",
  chainSelector: {
    name: "bitcichain-mainnet",
    selector: 4874388048629246000n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var bitcichain_mainnet_default = network11;
var network12 = {
  chainId: "200901",
  chainSelector: {
    name: "bitcoin-mainnet-bitlayer-1",
    selector: 7937294810946806131n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var bitcoin_mainnet_bitlayer_1_default = network12;
var network13 = {
  chainId: "60808",
  chainSelector: {
    name: "bitcoin-mainnet-bob-1",
    selector: 3849287863852499584n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var bitcoin_mainnet_bob_1_default = network13;
var network14 = {
  chainId: "3637",
  chainSelector: {
    name: "bitcoin-mainnet-botanix",
    selector: 4560701533377838164n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var bitcoin_mainnet_botanix_default = network14;
var network15 = {
  chainId: "223",
  chainSelector: {
    name: "bitcoin-mainnet-bsquared-1",
    selector: 5406759801798337480n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var bitcoin_mainnet_bsquared_1_default = network15;
var network16 = {
  chainId: "4200",
  chainSelector: {
    name: "bitcoin-merlin-mainnet",
    selector: 241851231317828981n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var bitcoin_merlin_mainnet_default = network16;
var network17 = {
  chainId: "964",
  chainSelector: {
    name: "bittensor-mainnet",
    selector: 2135107236357186872n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var bittensor_mainnet_default = network17;
var network18 = {
  chainId: "199",
  chainSelector: {
    name: "bittorrent_chain-mainnet",
    selector: 3776006016387883143n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var bittorrent_chain_mainnet_default = network18;
var network19 = {
  chainId: "42220",
  chainSelector: {
    name: "celo-mainnet",
    selector: 1346049177634351622n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var celo_mainnet_default = network19;
var network20 = {
  chainId: "81224",
  chainSelector: {
    name: "codex-mainnet",
    selector: 9478124434908827753n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var codex_mainnet_default = network20;
var network21 = {
  chainId: "52",
  chainSelector: {
    name: "coinex_smart_chain-mainnet",
    selector: 1761333065194157300n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var coinex_smart_chain_mainnet_default = network21;
var network22 = {
  chainId: "1030",
  chainSelector: {
    name: "conflux-mainnet",
    selector: 3358365939762719202n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var conflux_mainnet_default = network22;
var network23 = {
  chainId: "1116",
  chainSelector: {
    name: "core-mainnet",
    selector: 1224752112135636129n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var core_mainnet_default = network23;
var network24 = {
  chainId: "21000000",
  chainSelector: {
    name: "corn-mainnet",
    selector: 9043146809313071210n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var corn_mainnet_default = network24;
var network25 = {
  chainId: "25",
  chainSelector: {
    name: "cronos-mainnet",
    selector: 1456215246176062136n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var cronos_mainnet_default = network25;
var network26 = {
  chainId: "388",
  chainSelector: {
    name: "cronos-zkevm-mainnet",
    selector: 8788096068760390840n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var cronos_zkevm_mainnet_default = network26;
var network27 = {
  chainId: "1",
  chainSelector: {
    name: "ethereum-mainnet",
    selector: 5009297550715157269n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var ethereum_mainnet_default = network27;
var network28 = {
  chainId: "42161",
  chainSelector: {
    name: "ethereum-mainnet-arbitrum-1",
    selector: 4949039107694359620n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var ethereum_mainnet_arbitrum_1_default = network28;
var network29 = {
  chainId: "12324",
  chainSelector: {
    name: "ethereum-mainnet-arbitrum-1-l3x-1",
    selector: 3162193654116181371n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var ethereum_mainnet_arbitrum_1_l3x_1_default = network29;
var network30 = {
  chainId: "978670",
  chainSelector: {
    name: "ethereum-mainnet-arbitrum-1-treasure-1",
    selector: 1010349088906777999n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var ethereum_mainnet_arbitrum_1_treasure_1_default = network30;
var network31 = {
  chainId: "3776",
  chainSelector: {
    name: "ethereum-mainnet-astar-zkevm-1",
    selector: 1540201334317828111n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var ethereum_mainnet_astar_zkevm_1_default = network31;
var network32 = {
  chainId: "8453",
  chainSelector: {
    name: "ethereum-mainnet-base-1",
    selector: 15971525489660198786n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var ethereum_mainnet_base_1_default = network32;
var network33 = {
  chainId: "81457",
  chainSelector: {
    name: "ethereum-mainnet-blast-1",
    selector: 4411394078118774322n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var ethereum_mainnet_blast_1_default = network33;
var network34 = {
  chainId: "177",
  chainSelector: {
    name: "ethereum-mainnet-hashkey-1",
    selector: 7613811247471741961n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var ethereum_mainnet_hashkey_1_default = network34;
var network35 = {
  chainId: "13371",
  chainSelector: {
    name: "ethereum-mainnet-immutable-zkevm-1",
    selector: 1237925231416731909n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var ethereum_mainnet_immutable_zkevm_1_default = network35;
var network36 = {
  chainId: "57073",
  chainSelector: {
    name: "ethereum-mainnet-ink-1",
    selector: 3461204551265785888n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var ethereum_mainnet_ink_1_default = network36;
var network37 = {
  chainId: "255",
  chainSelector: {
    name: "ethereum-mainnet-kroma-1",
    selector: 3719320017875267166n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var ethereum_mainnet_kroma_1_default = network37;
var network38 = {
  chainId: "59144",
  chainSelector: {
    name: "ethereum-mainnet-linea-1",
    selector: 4627098889531055414n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var ethereum_mainnet_linea_1_default = network38;
var network39 = {
  chainId: "5000",
  chainSelector: {
    name: "ethereum-mainnet-mantle-1",
    selector: 1556008542357238666n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var ethereum_mainnet_mantle_1_default = network39;
var network40 = {
  chainId: "1088",
  chainSelector: {
    name: "ethereum-mainnet-metis-1",
    selector: 8805746078405598895n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var ethereum_mainnet_metis_1_default = network40;
var network41 = {
  chainId: "34443",
  chainSelector: {
    name: "ethereum-mainnet-mode-1",
    selector: 7264351850409363825n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var ethereum_mainnet_mode_1_default = network41;
var network42 = {
  chainId: "10",
  chainSelector: {
    name: "ethereum-mainnet-optimism-1",
    selector: 3734403246176062136n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var ethereum_mainnet_optimism_1_default = network42;
var network43 = {
  chainId: "1101",
  chainSelector: {
    name: "ethereum-mainnet-polygon-zkevm-1",
    selector: 4348158687435793198n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var ethereum_mainnet_polygon_zkevm_1_default = network43;
var network44 = {
  chainId: "534352",
  chainSelector: {
    name: "ethereum-mainnet-scroll-1",
    selector: 13204309965629103672n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var ethereum_mainnet_scroll_1_default = network44;
var network45 = {
  chainId: "167000",
  chainSelector: {
    name: "ethereum-mainnet-taiko-1",
    selector: 16468599424800719238n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var ethereum_mainnet_taiko_1_default = network45;
var network46 = {
  chainId: "130",
  chainSelector: {
    name: "ethereum-mainnet-unichain-1",
    selector: 1923510103922296319n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var ethereum_mainnet_unichain_1_default = network46;
var network47 = {
  chainId: "480",
  chainSelector: {
    name: "ethereum-mainnet-worldchain-1",
    selector: 2049429975587534727n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var ethereum_mainnet_worldchain_1_default = network47;
var network48 = {
  chainId: "196",
  chainSelector: {
    name: "ethereum-mainnet-xlayer-1",
    selector: 3016212468291539606n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var ethereum_mainnet_xlayer_1_default = network48;
var network49 = {
  chainId: "48900",
  chainSelector: {
    name: "ethereum-mainnet-zircuit-1",
    selector: 17198166215261833993n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var ethereum_mainnet_zircuit_1_default = network49;
var network50 = {
  chainId: "324",
  chainSelector: {
    name: "ethereum-mainnet-zksync-1",
    selector: 1562403441176082196n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var ethereum_mainnet_zksync_1_default = network50;
var network51 = {
  chainId: "42793",
  chainSelector: {
    name: "etherlink-mainnet",
    selector: 13624601974233774587n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var etherlink_mainnet_default = network51;
var network52 = {
  chainId: "250",
  chainSelector: {
    name: "fantom-mainnet",
    selector: 3768048213127883732n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var fantom_mainnet_default = network52;
var network53 = {
  chainId: "314",
  chainSelector: {
    name: "filecoin-mainnet",
    selector: 4561443241176882990n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var filecoin_mainnet_default = network53;
var network54 = {
  chainId: "252",
  chainSelector: {
    name: "fraxtal-mainnet",
    selector: 1462016016387883143n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var fraxtal_mainnet_default = network54;
var network55 = {
  chainId: "100",
  chainSelector: {
    name: "gnosis_chain-mainnet",
    selector: 465200170687744372n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var gnosis_chain_mainnet_default = network55;
var network56 = {
  chainId: "295",
  chainSelector: {
    name: "hedera-mainnet",
    selector: 3229138320728879060n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var hedera_mainnet_default = network56;
var network57 = {
  chainId: "43111",
  chainSelector: {
    name: "hemi-mainnet",
    selector: 1804312132722180201n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var hemi_mainnet_default = network57;
var network58 = {
  chainId: "999",
  chainSelector: {
    name: "hyperliquid-mainnet",
    selector: 2442541497099098535n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var hyperliquid_mainnet_default = network58;
var network59 = {
  chainId: "678",
  chainSelector: {
    name: "janction-mainnet",
    selector: 9107126442626377432n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var janction_mainnet_default = network59;
var network60 = {
  chainId: "8217",
  chainSelector: {
    name: "kaia-mainnet",
    selector: 9813823125703490621n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var kaia_mainnet_default = network60;
var network61 = {
  chainId: "2222",
  chainSelector: {
    name: "kava-mainnet",
    selector: 7550000543357438061n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var kava_mainnet_default = network61;
var network62 = {
  chainId: "1285",
  chainSelector: {
    name: "kusama-mainnet-moonriver",
    selector: 1355020143337428062n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var kusama_mainnet_moonriver_default = network62;
var network63 = {
  chainId: "232",
  chainSelector: {
    name: "lens-mainnet",
    selector: 5608378062013572713n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var lens_mainnet_default = network63;
var network64 = {
  chainId: "1135",
  chainSelector: {
    name: "lisk-mainnet",
    selector: 15293031020466096408n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var lisk_mainnet_default = network64;
var network65 = {
  chainId: "51888",
  chainSelector: {
    name: "memento-mainnet",
    selector: 6473245816409426016n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var memento_mainnet_default = network65;
var network66 = {
  chainId: "1750",
  chainSelector: {
    name: "metal-mainnet",
    selector: 13447077090413146373n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var metal_mainnet_default = network66;
var network67 = {
  chainId: "228",
  chainSelector: {
    name: "mind-mainnet",
    selector: 11690709103138290329n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var mind_mainnet_default = network67;
var network68 = {
  chainId: "185",
  chainSelector: {
    name: "mint-mainnet",
    selector: 17164792800244661392n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var mint_mainnet_default = network68;
var network69 = {
  chainId: "143",
  chainSelector: {
    name: "monad-mainnet",
    selector: 8481857512324358265n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var monad_mainnet_default = network69;
var network70 = {
  chainId: "2818",
  chainSelector: {
    name: "morph-mainnet",
    selector: 18164309074156128038n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var morph_mainnet_default = network70;
var network71 = {
  chainId: "397",
  chainSelector: {
    name: "near-mainnet",
    selector: 2039744413822257700n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var near_mainnet_default = network71;
var network72 = {
  chainId: "259",
  chainSelector: {
    name: "neonlink-mainnet",
    selector: 8239338020728974000n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var neonlink_mainnet_default = network72;
var network73 = {
  chainId: "47763",
  chainSelector: {
    name: "neox-mainnet",
    selector: 7222032299962346917n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var neox_mainnet_default = network73;
var network74 = {
  chainId: "68414",
  chainSelector: {
    name: "nexon-mainnet-henesys",
    selector: 12657445206920369324n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var nexon_mainnet_henesys_default = network74;
var network75 = {
  chainId: "60118",
  chainSelector: {
    name: "nexon-mainnet-lith",
    selector: 15758750456714168963n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var nexon_mainnet_lith_default = network75;
var network76 = {
  chainId: "807424",
  chainSelector: {
    name: "nexon-qa",
    selector: 14632960069656270105n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var nexon_qa_default = network76;
var network77 = {
  chainId: "847799",
  chainSelector: {
    name: "nexon-stage",
    selector: 5556806327594153475n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var nexon_stage_default = network77;
var network78 = {
  chainId: "6900",
  chainSelector: {
    name: "nibiru-mainnet",
    selector: 17349189558768828726n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var nibiru_mainnet_default = network78;
var network79 = {
  chainId: "9745",
  chainSelector: {
    name: "plasma-mainnet",
    selector: 9335212494177455608n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var plasma_mainnet_default = network79;
var network80 = {
  chainId: "98866",
  chainSelector: {
    name: "plume-mainnet",
    selector: 17912061998839310979n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var plume_mainnet_default = network80;
var network81 = {
  chainId: "592",
  chainSelector: {
    name: "polkadot-mainnet-astar",
    selector: 6422105447186081193n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var polkadot_mainnet_astar_default = network81;
var network82 = {
  chainId: "2031",
  chainSelector: {
    name: "polkadot-mainnet-centrifuge",
    selector: 8175830712062617656n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var polkadot_mainnet_centrifuge_default = network82;
var network83 = {
  chainId: "46",
  chainSelector: {
    name: "polkadot-mainnet-darwinia",
    selector: 8866418665544333000n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var polkadot_mainnet_darwinia_default = network83;
var network84 = {
  chainId: "1284",
  chainSelector: {
    name: "polkadot-mainnet-moonbeam",
    selector: 1252863800116739621n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var polkadot_mainnet_moonbeam_default = network84;
var network85 = {
  chainId: "137",
  chainSelector: {
    name: "polygon-mainnet",
    selector: 4051577828743386545n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var polygon_mainnet_default = network85;
var network86 = {
  chainId: "747474",
  chainSelector: {
    name: "polygon-mainnet-katana",
    selector: 2459028469735686113n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var polygon_mainnet_katana_default = network86;
var network87 = {
  chainId: "2020",
  chainSelector: {
    name: "ronin-mainnet",
    selector: 6916147374840168594n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var ronin_mainnet_default = network87;
var network88 = {
  chainId: "30",
  chainSelector: {
    name: "rootstock-mainnet",
    selector: 11964252391146578476n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var rootstock_mainnet_default = network88;
var network89 = {
  chainId: "1329",
  chainSelector: {
    name: "sei-mainnet",
    selector: 9027416829622342829n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var sei_mainnet_default = network89;
var network90 = {
  chainId: "109",
  chainSelector: {
    name: "shibarium-mainnet",
    selector: 3993510008929295315n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var shibarium_mainnet_default = network90;
var network91 = {
  chainId: "1868",
  chainSelector: {
    name: "soneium-mainnet",
    selector: 12505351618335765396n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var soneium_mainnet_default = network91;
var network92 = {
  chainId: "146",
  chainSelector: {
    name: "sonic-mainnet",
    selector: 1673871237479749969n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var sonic_mainnet_default = network92;
var network93 = {
  chainId: "5330",
  chainSelector: {
    name: "superseed-mainnet",
    selector: 470401360549526817n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var superseed_mainnet_default = network93;
var network94 = {
  chainId: "239",
  chainSelector: {
    name: "tac-mainnet",
    selector: 5936861837188149645n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var tac_mainnet_default = network94;
var network95 = {
  chainId: "40",
  chainSelector: {
    name: "telos-evm-mainnet",
    selector: 1477345371608778000n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var telos_evm_mainnet_default = network95;
var network96 = {
  chainId: "61166",
  chainSelector: {
    name: "treasure-mainnet",
    selector: 5214452172935136222n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var treasure_mainnet_default = network96;
var network97 = {
  chainId: "728126428",
  chainSelector: {
    name: "tron-mainnet-evm",
    selector: 1546563616611573946n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var tron_mainnet_evm_default = network97;
var network98 = {
  chainId: "106",
  chainSelector: {
    name: "velas-mainnet",
    selector: 374210358663784372n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var velas_mainnet_default = network98;
var network99 = {
  chainId: "1111",
  chainSelector: {
    name: "wemix-mainnet",
    selector: 5142893604156789321n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var wemix_mainnet_default = network99;
var network100 = {
  chainId: "50",
  chainSelector: {
    name: "xdc-mainnet",
    selector: 17673274061779414707n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var xdc_mainnet_default = network100;
var network101 = {
  chainId: "7000",
  chainSelector: {
    name: "zetachain-mainnet",
    selector: 10817664450262215148n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var zetachain_mainnet_default = network101;
var network102 = {
  chainId: "810180",
  chainSelector: {
    name: "zklink_nova-mainnet",
    selector: 4350319965322101699n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var zklink_nova_mainnet_default = network102;
var network103 = {
  chainId: "7777777",
  chainSelector: {
    name: "zora-mainnet",
    selector: 3555797439612589184n
  },
  chainFamily: "evm",
  networkType: "mainnet"
};
var zora_mainnet_default = network103;
var network104 = {
  chainId: "5eykt4UsFv8P8NJdTREpY1vzqKqZKvdpKuc147dw2N9d",
  chainSelector: {
    name: "solana-mainnet",
    selector: 124615329519749607n
  },
  chainFamily: "solana",
  networkType: "mainnet"
};
var solana_mainnet_default = network104;
var network105 = {
  chainId: "1",
  chainSelector: {
    name: "sui-mainnet",
    selector: 17529533435026248318n
  },
  chainFamily: "sui",
  networkType: "mainnet"
};
var sui_mainnet_default = network105;
var network106 = {
  chainId: "-239",
  chainSelector: {
    name: "ton-mainnet",
    selector: 16448340667252469081n
  },
  chainFamily: "ton",
  networkType: "mainnet"
};
var ton_mainnet_default = network106;
var network107 = {
  chainId: "728126428",
  chainSelector: {
    name: "tron-mainnet",
    selector: 1546563616611573945n
  },
  chainFamily: "tron",
  networkType: "mainnet"
};
var tron_mainnet_default = network107;
var network108 = {
  chainId: "4",
  chainSelector: {
    name: "aptos-localnet",
    selector: 4457093679053095497n
  },
  chainFamily: "aptos",
  networkType: "testnet"
};
var aptos_localnet_default = network108;
var network109 = {
  chainId: "2",
  chainSelector: {
    name: "aptos-testnet",
    selector: 743186221051783445n
  },
  chainFamily: "aptos",
  networkType: "testnet"
};
var aptos_testnet_default = network109;
var network110 = {
  chainId: "16601",
  chainSelector: {
    name: "0g-testnet-galileo",
    selector: 2131427466778448014n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var _0g_testnet_galileo_default = network110;
var network111 = {
  chainId: "16600",
  chainSelector: {
    name: "0g-testnet-newton",
    selector: 16088006396410204581n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var _0g_testnet_newton_default = network111;
var network112 = {
  chainId: "11124",
  chainSelector: {
    name: "abstract-testnet",
    selector: 16235373811196386733n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var abstract_testnet_default = network112;
var network113 = {
  chainId: "31337",
  chainSelector: {
    name: "anvil-devnet",
    selector: 7759470850252068959n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var anvil_devnet_default = network113;
var network114 = {
  chainId: "33111",
  chainSelector: {
    name: "apechain-testnet-curtis",
    selector: 9900119385908781505n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var apechain_testnet_curtis_default = network114;
var network115 = {
  chainId: "462",
  chainSelector: {
    name: "areon-testnet",
    selector: 7317911323415911000n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var areon_testnet_default = network115;
var network116 = {
  chainId: "432201",
  chainSelector: {
    name: "avalanche-subnet-dexalot-testnet",
    selector: 1458281248224512906n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var avalanche_subnet_dexalot_testnet_default = network116;
var network117 = {
  chainId: "43113",
  chainSelector: {
    name: "avalanche-testnet-fuji",
    selector: 14767482510784806043n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var avalanche_testnet_fuji_default = network117;
var network118 = {
  chainId: "595581",
  chainSelector: {
    name: "avalanche-testnet-nexon",
    selector: 7837562506228496256n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var avalanche_testnet_nexon_default = network118;
var network119 = {
  chainId: "80085",
  chainSelector: {
    name: "berachain-testnet-artio",
    selector: 12336603543561911511n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var berachain_testnet_artio_default = network119;
var network120 = {
  chainId: "80084",
  chainSelector: {
    name: "berachain-testnet-bartio",
    selector: 8999465244383784164n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var berachain_testnet_bartio_default = network120;
var network121 = {
  chainId: "80069",
  chainSelector: {
    name: "berachain-testnet-bepolia",
    selector: 7728255861635209484n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var berachain_testnet_bepolia_default = network121;
var network122 = {
  chainId: "97",
  chainSelector: {
    name: "binance_smart_chain-testnet",
    selector: 13264668187771770619n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var binance_smart_chain_testnet_default = network122;
var network123 = {
  chainId: "5611",
  chainSelector: {
    name: "binance_smart_chain-testnet-opbnb-1",
    selector: 13274425992935471758n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var binance_smart_chain_testnet_opbnb_1_default = network123;
var network124 = {
  chainId: "1908",
  chainSelector: {
    name: "bitcichain-testnet",
    selector: 4888058894222120000n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var bitcichain_testnet_default = network124;
var network125 = {
  chainId: "200810",
  chainSelector: {
    name: "bitcoin-testnet-bitlayer-1",
    selector: 3789623672476206327n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var bitcoin_testnet_bitlayer_1_default = network125;
var network126 = {
  chainId: "3636",
  chainSelector: {
    name: "bitcoin-testnet-botanix",
    selector: 1467223411771711614n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var bitcoin_testnet_botanix_default = network126;
var network127 = {
  chainId: "1123",
  chainSelector: {
    name: "bitcoin-testnet-bsquared-1",
    selector: 1948510578179542068n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var bitcoin_testnet_bsquared_1_default = network127;
var network128 = {
  chainId: "686868",
  chainSelector: {
    name: "bitcoin-testnet-merlin",
    selector: 5269261765892944301n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var bitcoin_testnet_merlin_default = network128;
var network129 = {
  chainId: "31",
  chainSelector: {
    name: "bitcoin-testnet-rootstock",
    selector: 8953668971247136127n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var bitcoin_testnet_rootstock_default = network129;
var network130 = {
  chainId: "808813",
  chainSelector: {
    name: "bitcoin-testnet-sepolia-bob-1",
    selector: 5535534526963509396n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var bitcoin_testnet_sepolia_bob_1_default = network130;
var network131 = {
  chainId: "945",
  chainSelector: {
    name: "bittensor-testnet",
    selector: 2177900824115119161n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var bittensor_testnet_default = network131;
var network132 = {
  chainId: "1029",
  chainSelector: {
    name: "bittorrent_chain-testnet",
    selector: 4459371029167934217n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var bittorrent_chain_testnet_default = network132;
var network133 = {
  chainId: "44787",
  chainSelector: {
    name: "celo-testnet-alfajores",
    selector: 3552045678561919002n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var celo_testnet_alfajores_default = network133;
var network134 = {
  chainId: "812242",
  chainSelector: {
    name: "codex-testnet",
    selector: 7225665875429174318n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var codex_testnet_default = network134;
var network135 = {
  chainId: "53",
  chainSelector: {
    name: "coinex_smart_chain-testnet",
    selector: 8955032871639343000n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var coinex_smart_chain_testnet_default = network135;
var network136 = {
  chainId: "1114",
  chainSelector: {
    name: "core-testnet",
    selector: 4264732132125536123n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var core_testnet_default = network136;
var network137 = {
  chainId: "338",
  chainSelector: {
    name: "cronos-testnet",
    selector: 2995292832068775165n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var cronos_testnet_default = network137;
var network138 = {
  chainId: "282",
  chainSelector: {
    name: "cronos-testnet-zkevm-1",
    selector: 3842103497652714138n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var cronos_testnet_zkevm_1_default = network138;
var network139 = {
  chainId: "240",
  chainSelector: {
    name: "cronos-zkevm-testnet-sepolia",
    selector: 16487132492576884721n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var cronos_zkevm_testnet_sepolia_default = network139;
var network140 = {
  chainId: "2025",
  chainSelector: {
    name: "dtcc-testnet-andesite",
    selector: 15513093881969820114n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var dtcc_testnet_andesite_default = network140;
var network141 = {
  chainId: "421613",
  chainSelector: {
    name: "ethereum-testnet-goerli-arbitrum-1",
    selector: 6101244977088475029n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_goerli_arbitrum_1_default = network141;
var network142 = {
  chainId: "84531",
  chainSelector: {
    name: "ethereum-testnet-goerli-base-1",
    selector: 5790810961207155433n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_goerli_base_1_default = network142;
var network143 = {
  chainId: "59140",
  chainSelector: {
    name: "ethereum-testnet-goerli-linea-1",
    selector: 1355246678561316402n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_goerli_linea_1_default = network143;
var network144 = {
  chainId: "5001",
  chainSelector: {
    name: "ethereum-testnet-goerli-mantle-1",
    selector: 4168263376276232250n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_goerli_mantle_1_default = network144;
var network145 = {
  chainId: "420",
  chainSelector: {
    name: "ethereum-testnet-goerli-optimism-1",
    selector: 2664363617261496610n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_goerli_optimism_1_default = network145;
var network146 = {
  chainId: "1442",
  chainSelector: {
    name: "ethereum-testnet-goerli-polygon-zkevm-1",
    selector: 11059667695644972511n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_goerli_polygon_zkevm_1_default = network146;
var network147 = {
  chainId: "280",
  chainSelector: {
    name: "ethereum-testnet-goerli-zksync-1",
    selector: 6802309497652714138n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_goerli_zksync_1_default = network147;
var network148 = {
  chainId: "17000",
  chainSelector: {
    name: "ethereum-testnet-holesky",
    selector: 7717148896336251131n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_holesky_default = network148;
var network149 = {
  chainId: "2522",
  chainSelector: {
    name: "ethereum-testnet-holesky-fraxtal-1",
    selector: 8901520481741771655n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_holesky_fraxtal_1_default = network149;
var network150 = {
  chainId: "2810",
  chainSelector: {
    name: "ethereum-testnet-holesky-morph-1",
    selector: 8304510386741731151n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_holesky_morph_1_default = network150;
var network151 = {
  chainId: "167009",
  chainSelector: {
    name: "ethereum-testnet-holesky-taiko-1",
    selector: 7248756420937879088n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_holesky_taiko_1_default = network151;
var network152 = {
  chainId: "11155111",
  chainSelector: {
    name: "ethereum-testnet-sepolia",
    selector: 16015286601757825753n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_default = network152;
var network153 = {
  chainId: "421614",
  chainSelector: {
    name: "ethereum-testnet-sepolia-arbitrum-1",
    selector: 3478487238524512106n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_arbitrum_1_default = network153;
var network154 = {
  chainId: "12325",
  chainSelector: {
    name: "ethereum-testnet-sepolia-arbitrum-1-l3x-1",
    selector: 3486622437121596122n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_arbitrum_1_l3x_1_default = network154;
var network155 = {
  chainId: "978657",
  chainSelector: {
    name: "ethereum-testnet-sepolia-arbitrum-1-treasure-1",
    selector: 10443705513486043421n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_arbitrum_1_treasure_1_default = network155;
var network156 = {
  chainId: "84532",
  chainSelector: {
    name: "ethereum-testnet-sepolia-base-1",
    selector: 10344971235874465080n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_base_1_default = network156;
var network157 = {
  chainId: "168587773",
  chainSelector: {
    name: "ethereum-testnet-sepolia-blast-1",
    selector: 2027362563942762617n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_blast_1_default = network157;
var network158 = {
  chainId: "21000001",
  chainSelector: {
    name: "ethereum-testnet-sepolia-corn-1",
    selector: 1467427327723633929n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_corn_1_default = network158;
var network159 = {
  chainId: "133",
  chainSelector: {
    name: "ethereum-testnet-sepolia-hashkey-1",
    selector: 4356164186791070119n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_hashkey_1_default = network159;
var network160 = {
  chainId: "13473",
  chainSelector: {
    name: "ethereum-testnet-sepolia-immutable-zkevm-1",
    selector: 4526165231216331901n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_immutable_zkevm_1_default = network160;
var network161 = {
  chainId: "2358",
  chainSelector: {
    name: "ethereum-testnet-sepolia-kroma-1",
    selector: 5990477251245693094n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_kroma_1_default = network161;
var network162 = {
  chainId: "37111",
  chainSelector: {
    name: "ethereum-testnet-sepolia-lens-1",
    selector: 6827576821754315911n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_lens_1_default = network162;
var network163 = {
  chainId: "59141",
  chainSelector: {
    name: "ethereum-testnet-sepolia-linea-1",
    selector: 5719461335882077547n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_linea_1_default = network163;
var network164 = {
  chainId: "4202",
  chainSelector: {
    name: "ethereum-testnet-sepolia-lisk-1",
    selector: 5298399861320400553n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_lisk_1_default = network164;
var network165 = {
  chainId: "5003",
  chainSelector: {
    name: "ethereum-testnet-sepolia-mantle-1",
    selector: 8236463271206331221n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_mantle_1_default = network165;
var network166 = {
  chainId: "59902",
  chainSelector: {
    name: "ethereum-testnet-sepolia-metis-1",
    selector: 3777822886988675105n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_metis_1_default = network166;
var network167 = {
  chainId: "919",
  chainSelector: {
    name: "ethereum-testnet-sepolia-mode-1",
    selector: 829525985033418733n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_mode_1_default = network167;
var network168 = {
  chainId: "11155420",
  chainSelector: {
    name: "ethereum-testnet-sepolia-optimism-1",
    selector: 5224473277236331295n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_optimism_1_default = network168;
var network169 = {
  chainId: "717160",
  chainSelector: {
    name: "ethereum-testnet-sepolia-polygon-validium-1",
    selector: 4418231248214522936n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_polygon_validium_1_default = network169;
var network170 = {
  chainId: "2442",
  chainSelector: {
    name: "ethereum-testnet-sepolia-polygon-zkevm-1",
    selector: 1654667687261492630n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_polygon_zkevm_1_default = network170;
var network171 = {
  chainId: "534351",
  chainSelector: {
    name: "ethereum-testnet-sepolia-scroll-1",
    selector: 2279865765895943307n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_scroll_1_default = network171;
var network172 = {
  chainId: "1946",
  chainSelector: {
    name: "ethereum-testnet-sepolia-soneium-1",
    selector: 686603546605904534n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_soneium_1_default = network172;
var network173 = {
  chainId: "1301",
  chainSelector: {
    name: "ethereum-testnet-sepolia-unichain-1",
    selector: 14135854469784514356n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_unichain_1_default = network173;
var network174 = {
  chainId: "4801",
  chainSelector: {
    name: "ethereum-testnet-sepolia-worldchain-1",
    selector: 5299555114858065850n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_worldchain_1_default = network174;
var network175 = {
  chainId: "195",
  chainSelector: {
    name: "ethereum-testnet-sepolia-xlayer-1",
    selector: 2066098519157881736n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_xlayer_1_default = network175;
var network176 = {
  chainId: "48899",
  chainSelector: {
    name: "ethereum-testnet-sepolia-zircuit-1",
    selector: 4562743618362911021n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_zircuit_1_default = network176;
var network177 = {
  chainId: "300",
  chainSelector: {
    name: "ethereum-testnet-sepolia-zksync-1",
    selector: 6898391096552792247n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ethereum_testnet_sepolia_zksync_1_default = network177;
var network178 = {
  chainId: "128123",
  chainSelector: {
    name: "etherlink-testnet",
    selector: 1910019406958449359n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var etherlink_testnet_default = network178;
var network179 = {
  chainId: "4002",
  chainSelector: {
    name: "fantom-testnet",
    selector: 4905564228793744293n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var fantom_testnet_default = network179;
var network180 = {
  chainId: "31415926",
  chainSelector: {
    name: "filecoin-testnet",
    selector: 7060342227814389000n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var filecoin_testnet_default = network180;
var network181 = {
  chainId: "1337",
  chainSelector: {
    name: "geth-testnet",
    selector: 3379446385462418246n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var geth_testnet_default = network181;
var network182 = {
  chainId: "10200",
  chainSelector: {
    name: "gnosis_chain-testnet-chiado",
    selector: 8871595565390010547n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var gnosis_chain_testnet_chiado_default = network182;
var network183 = {
  chainId: "296",
  chainSelector: {
    name: "hedera-testnet",
    selector: 222782988166878823n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var hedera_testnet_default = network183;
var network184 = {
  chainId: "743111",
  chainSelector: {
    name: "hemi-testnet-sepolia",
    selector: 16126893759944359622n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var hemi_testnet_sepolia_default = network184;
var network185 = {
  chainId: "998",
  chainSelector: {
    name: "hyperliquid-testnet",
    selector: 4286062357653186312n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var hyperliquid_testnet_default = network185;
var network186 = {
  chainId: "763373",
  chainSelector: {
    name: "ink-testnet-sepolia",
    selector: 9763904284804119144n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ink_testnet_sepolia_default = network186;
var network187 = {
  chainId: "679",
  chainSelector: {
    name: "janction-testnet-sepolia",
    selector: 5059197667603797935n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var janction_testnet_sepolia_default = network187;
var network188 = {
  chainId: "2019775",
  chainSelector: {
    name: "jovay-testnet",
    selector: 945045181441419236n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var jovay_testnet_default = network188;
var network189 = {
  chainId: "1001",
  chainSelector: {
    name: "kaia-testnet-kairos",
    selector: 2624132734533621656n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var kaia_testnet_kairos_default = network189;
var network190 = {
  chainId: "2221",
  chainSelector: {
    name: "kava-testnet",
    selector: 2110537777356199208n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var kava_testnet_default = network190;
var network191 = {
  chainId: "6342",
  chainSelector: {
    name: "megaeth-testnet",
    selector: 2443239559770384419n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var megaeth_testnet_default = network191;
var network192 = {
  chainId: "2129",
  chainSelector: {
    name: "memento-testnet",
    selector: 12168171414969487009n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var memento_testnet_default = network192;
var network193 = {
  chainId: "1740",
  chainSelector: {
    name: "metal-testnet",
    selector: 6286293440461807648n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var metal_testnet_default = network193;
var network194 = {
  chainId: "192940",
  chainSelector: {
    name: "mind-testnet",
    selector: 7189150270347329685n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var mind_testnet_default = network194;
var network195 = {
  chainId: "1687",
  chainSelector: {
    name: "mint-testnet",
    selector: 10749384167430721561n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var mint_testnet_default = network195;
var network196 = {
  chainId: "10143",
  chainSelector: {
    name: "monad-testnet",
    selector: 2183018362218727504n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var monad_testnet_default = network196;
var network197 = {
  chainId: "398",
  chainSelector: {
    name: "near-testnet",
    selector: 5061593697262339000n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var near_testnet_default = network197;
var network198 = {
  chainId: "9559",
  chainSelector: {
    name: "neonlink-testnet",
    selector: 1113014352258747600n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var neonlink_testnet_default = network198;
var network199 = {
  chainId: "12227332",
  chainSelector: {
    name: "neox-testnet-t4",
    selector: 2217764097022649312n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var neox_testnet_t4_default = network199;
var network200 = {
  chainId: "5668",
  chainSelector: {
    name: "nexon-dev",
    selector: 8911150974185440581n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var nexon_dev_default = network200;
var network201 = {
  chainId: "6930",
  chainSelector: {
    name: "nibiru-testnet",
    selector: 305104239123120457n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var nibiru_testnet_default = network201;
var network202 = {
  chainId: "9000",
  chainSelector: {
    name: "ondo-testnet",
    selector: 344208382356656551n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ondo_testnet_default = network202;
var network203 = {
  chainId: "688688",
  chainSelector: {
    name: "pharos-testnet",
    selector: 4012524741200567430n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var pharos_testnet_default = network203;
var network204 = {
  chainId: "9746",
  chainSelector: {
    name: "plasma-testnet",
    selector: 3967220077692964309n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var plasma_testnet_default = network204;
var network205 = {
  chainId: "98864",
  chainSelector: {
    name: "plume-devnet",
    selector: 3743020999916460931n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var plume_devnet_default = network205;
var network206 = {
  chainId: "161221135",
  chainSelector: {
    name: "plume-testnet",
    selector: 14684575664602284776n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var plume_testnet_default = network206;
var network207 = {
  chainId: "98867",
  chainSelector: {
    name: "plume-testnet-sepolia",
    selector: 13874588925447303949n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var plume_testnet_sepolia_default = network207;
var network208 = {
  chainId: "81",
  chainSelector: {
    name: "polkadot-testnet-astar-shibuya",
    selector: 6955638871347136141n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var polkadot_testnet_astar_shibuya_default = network208;
var network209 = {
  chainId: "2088",
  chainSelector: {
    name: "polkadot-testnet-centrifuge-altair",
    selector: 2333097300889804761n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var polkadot_testnet_centrifuge_altair_default = network209;
var network210 = {
  chainId: "45",
  chainSelector: {
    name: "polkadot-testnet-darwinia-pangoro",
    selector: 4340886533089894000n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var polkadot_testnet_darwinia_pangoro_default = network210;
var network211 = {
  chainId: "1287",
  chainSelector: {
    name: "polkadot-testnet-moonbeam-moonbase",
    selector: 5361632739113536121n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var polkadot_testnet_moonbeam_moonbase_default = network211;
var network212 = {
  chainId: "80002",
  chainSelector: {
    name: "polygon-testnet-amoy",
    selector: 16281711391670634445n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var polygon_testnet_amoy_default = network212;
var network213 = {
  chainId: "80001",
  chainSelector: {
    name: "polygon-testnet-mumbai",
    selector: 12532609583862916517n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var polygon_testnet_mumbai_default = network213;
var network214 = {
  chainId: "129399",
  chainSelector: {
    name: "polygon-testnet-tatara",
    selector: 9090863410735740267n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var polygon_testnet_tatara_default = network214;
var network215 = {
  chainId: "2024",
  chainSelector: {
    name: "private-testnet-andesite",
    selector: 6915682381028791124n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var private_testnet_andesite_default = network215;
var network216 = {
  chainId: "2023",
  chainSelector: {
    name: "private-testnet-granite",
    selector: 3260900564719373474n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var private_testnet_granite_default = network216;
var network217 = {
  chainId: "424242",
  chainSelector: {
    name: "private-testnet-mica",
    selector: 4489326297382772450n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var private_testnet_mica_default = network217;
var network218 = {
  chainId: "682",
  chainSelector: {
    name: "private-testnet-obsidian",
    selector: 6260932437388305511n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var private_testnet_obsidian_default = network218;
var network219 = {
  chainId: "45439",
  chainSelector: {
    name: "private-testnet-opala",
    selector: 8446413392851542429n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var private_testnet_opala_default = network219;
var network220 = {
  chainId: "2021",
  chainSelector: {
    name: "ronin-testnet-saigon",
    selector: 13116810400804392105n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var ronin_testnet_saigon_default = network220;
var network221 = {
  chainId: "1328",
  chainSelector: {
    name: "sei-testnet-atlantic",
    selector: 1216300075444106652n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var sei_testnet_atlantic_default = network221;
var network222 = {
  chainId: "157",
  chainSelector: {
    name: "shibarium-testnet-puppynet",
    selector: 17833296867764334567n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var shibarium_testnet_puppynet_default = network222;
var network223 = {
  chainId: "57054",
  chainSelector: {
    name: "sonic-testnet-blaze",
    selector: 3676871237479449268n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var sonic_testnet_blaze_default = network223;
var network224 = {
  chainId: "1513",
  chainSelector: {
    name: "story-testnet",
    selector: 4237030917318060427n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var story_testnet_default = network224;
var network225 = {
  chainId: "53302",
  chainSelector: {
    name: "superseed-testnet",
    selector: 13694007683517087973n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var superseed_testnet_default = network225;
var network226 = {
  chainId: "2391",
  chainSelector: {
    name: "tac-testnet",
    selector: 9488606126177218005n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var tac_testnet_default = network226;
var network227 = {
  chainId: "41",
  chainSelector: {
    name: "telos-evm-testnet",
    selector: 729797994450396300n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var telos_evm_testnet_default = network227;
var network228 = {
  chainId: "978658",
  chainSelector: {
    name: "treasure-testnet-topaz",
    selector: 3676916124122457866n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var treasure_testnet_topaz_default = network228;
var network229 = {
  chainId: "3360022319",
  chainSelector: {
    name: "tron-devnet-evm",
    selector: 13231703482326770600n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var tron_devnet_evm_default = network229;
var network230 = {
  chainId: "3448148188",
  chainSelector: {
    name: "tron-testnet-nile-evm",
    selector: 2052925811360307749n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var tron_testnet_nile_evm_default = network230;
var network231 = {
  chainId: "2494104990",
  chainSelector: {
    name: "tron-testnet-shasta-evm",
    selector: 13231703482326770598n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var tron_testnet_shasta_evm_default = network231;
var network232 = {
  chainId: "111",
  chainSelector: {
    name: "velas-testnet",
    selector: 572210378683744374n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var velas_testnet_default = network232;
var network233 = {
  chainId: "1112",
  chainSelector: {
    name: "wemix-testnet",
    selector: 9284632837123596123n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var wemix_testnet_default = network233;
var network234 = {
  chainId: "51",
  chainSelector: {
    name: "xdc-testnet",
    selector: 3017758115101368649n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var xdc_testnet_default = network234;
var network235 = {
  chainId: "80087",
  chainSelector: {
    name: "zero-g-testnet-galileo",
    selector: 2285225387454015855n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var zero_g_testnet_galileo_default = network235;
var network236 = {
  chainId: "48898",
  chainSelector: {
    name: "zircuit-testnet-garfield",
    selector: 13781831279385219069n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var zircuit_testnet_garfield_default = network236;
var network237 = {
  chainId: "810181",
  chainSelector: {
    name: "zklink_nova-testnet",
    selector: 5837261596322416298n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var zklink_nova_testnet_default = network237;
var network238 = {
  chainId: "999999999",
  chainSelector: {
    name: "zora-testnet",
    selector: 16244020411108056671n
  },
  chainFamily: "evm",
  networkType: "testnet"
};
var zora_testnet_default = network238;
var network239 = {
  chainId: "EtWTRABZaYq6iMfeYKouRu166VU2xqa1wcaWoxPkrZBG",
  chainSelector: {
    name: "solana-devnet",
    selector: 16423721717087811551n
  },
  chainFamily: "solana",
  networkType: "testnet"
};
var solana_devnet_default = network239;
var network240 = {
  chainId: "4uhcVJyU9pJkvQyS88uRDiswHXSCkY3zQawwpjk2NsNY",
  chainSelector: {
    name: "solana-testnet",
    selector: 6302590918974934319n
  },
  chainFamily: "solana",
  networkType: "testnet"
};
var solana_testnet_default = network240;
var network241 = {
  chainId: "4",
  chainSelector: {
    name: "sui-localnet",
    selector: 18395503381733958356n
  },
  chainFamily: "sui",
  networkType: "testnet"
};
var sui_localnet_default = network241;
var network242 = {
  chainId: "2",
  chainSelector: {
    name: "sui-testnet",
    selector: 9762610643973837292n
  },
  chainFamily: "sui",
  networkType: "testnet"
};
var sui_testnet_default = network242;
var network243 = {
  chainId: "-217",
  chainSelector: {
    name: "ton-localnet",
    selector: 13879075125137744094n
  },
  chainFamily: "ton",
  networkType: "testnet"
};
var ton_localnet_default = network243;
var network244 = {
  chainId: "-3",
  chainSelector: {
    name: "ton-testnet",
    selector: 1399300952838017768n
  },
  chainFamily: "ton",
  networkType: "testnet"
};
var ton_testnet_default = network244;
var network245 = {
  chainId: "3360022319",
  chainSelector: {
    name: "tron-devnet",
    selector: 13231703482326770599n
  },
  chainFamily: "tron",
  networkType: "testnet"
};
var tron_devnet_default = network245;
var network246 = {
  chainId: "3448148188",
  chainSelector: {
    name: "tron-testnet-nile",
    selector: 2052925811360307740n
  },
  chainFamily: "tron",
  networkType: "testnet"
};
var tron_testnet_nile_default = network246;
var network247 = {
  chainId: "2494104990",
  chainSelector: {
    name: "tron-testnet-shasta",
    selector: 13231703482326770597n
  },
  chainFamily: "tron",
  networkType: "testnet"
};
var tron_testnet_shasta_default = network247;
var mainnetBySelector = new Map([
  [5009297550715157269n, ethereum_mainnet_default],
  [3734403246176062136n, ethereum_mainnet_optimism_1_default],
  [1456215246176062136n, cronos_mainnet_default],
  [11964252391146578476n, rootstock_mainnet_default],
  [1477345371608778000n, telos_evm_mainnet_default],
  [8866418665544333000n, polkadot_mainnet_darwinia_default],
  [17673274061779414707n, xdc_mainnet_default],
  [1761333065194157300n, coinex_smart_chain_mainnet_default],
  [11344663589394136015n, binance_smart_chain_mainnet_default],
  [465200170687744372n, gnosis_chain_mainnet_default],
  [374210358663784372n, velas_mainnet_default],
  [3993510008929295315n, shibarium_mainnet_default],
  [1923510103922296319n, ethereum_mainnet_unichain_1_default],
  [4051577828743386545n, polygon_mainnet_default],
  [8481857512324358265n, monad_mainnet_default],
  [1673871237479749969n, sonic_mainnet_default],
  [7613811247471741961n, ethereum_mainnet_hashkey_1_default],
  [17164792800244661392n, mint_mainnet_default],
  [3016212468291539606n, ethereum_mainnet_xlayer_1_default],
  [3776006016387883143n, bittorrent_chain_mainnet_default],
  [465944652040885897n, binance_smart_chain_mainnet_opbnb_1_default],
  [5406759801798337480n, bitcoin_mainnet_bsquared_1_default],
  [11690709103138290329n, mind_mainnet_default],
  [5608378062013572713n, lens_mainnet_default],
  [5936861837188149645n, tac_mainnet_default],
  [3768048213127883732n, fantom_mainnet_default],
  [1462016016387883143n, fraxtal_mainnet_default],
  [3719320017875267166n, ethereum_mainnet_kroma_1_default],
  [8239338020728974000n, neonlink_mainnet_default],
  [3229138320728879060n, hedera_mainnet_default],
  [4561443241176882990n, filecoin_mainnet_default],
  [1562403441176082196n, ethereum_mainnet_zksync_1_default],
  [8788096068760390840n, cronos_zkevm_mainnet_default],
  [2039744413822257700n, near_mainnet_default],
  [1939936305787790600n, areon_mainnet_default],
  [2049429975587534727n, ethereum_mainnet_worldchain_1_default],
  [6422105447186081193n, polkadot_mainnet_astar_default],
  [9107126442626377432n, janction_mainnet_default],
  [2135107236357186872n, bittensor_mainnet_default],
  [2442541497099098535n, hyperliquid_mainnet_default],
  [3358365939762719202n, conflux_mainnet_default],
  [8805746078405598895n, ethereum_mainnet_metis_1_default],
  [4348158687435793198n, ethereum_mainnet_polygon_zkevm_1_default],
  [5142893604156789321n, wemix_mainnet_default],
  [1224752112135636129n, core_mainnet_default],
  [15293031020466096408n, lisk_mainnet_default],
  [1252863800116739621n, polkadot_mainnet_moonbeam_default],
  [1355020143337428062n, kusama_mainnet_moonriver_default],
  [9027416829622342829n, sei_mainnet_default],
  [13447077090413146373n, metal_mainnet_default],
  [12505351618335765396n, soneium_mainnet_default],
  [4874388048629246000n, bitcichain_mainnet_default],
  [6916147374840168594n, ronin_mainnet_default],
  [8175830712062617656n, polkadot_mainnet_centrifuge_default],
  [7550000543357438061n, kava_mainnet_default],
  [3577778157919314504n, abstract_mainnet_default],
  [18164309074156128038n, morph_mainnet_default],
  [4560701533377838164n, bitcoin_mainnet_botanix_default],
  [1540201334317828111n, ethereum_mainnet_astar_zkevm_1_default],
  [241851231317828981n, bitcoin_merlin_mainnet_default],
  [1556008542357238666n, ethereum_mainnet_mantle_1_default],
  [470401360549526817n, superseed_mainnet_default],
  [17349189558768828726n, nibiru_mainnet_default],
  [10817664450262215148n, zetachain_mainnet_default],
  [9813823125703490621n, kaia_mainnet_default],
  [15971525489660198786n, ethereum_mainnet_base_1_default],
  [9335212494177455608n, plasma_mainnet_default],
  [3162193654116181371n, ethereum_mainnet_arbitrum_1_l3x_1_default],
  [1237925231416731909n, ethereum_mainnet_immutable_zkevm_1_default],
  [4426351306075016396n, _0g_mainnet_default],
  [14894068710063348487n, apechain_mainnet_default],
  [7264351850409363825n, ethereum_mainnet_mode_1_default],
  [4949039107694359620n, ethereum_mainnet_arbitrum_1_default],
  [1346049177634351622n, celo_mainnet_default],
  [13624601974233774587n, etherlink_mainnet_default],
  [1804312132722180201n, hemi_mainnet_default],
  [6433500567565415381n, avalanche_mainnet_default],
  [7222032299962346917n, neox_mainnet_default],
  [17198166215261833993n, ethereum_mainnet_zircuit_1_default],
  [6473245816409426016n, memento_mainnet_default],
  [3461204551265785888n, ethereum_mainnet_ink_1_default],
  [4627098889531055414n, ethereum_mainnet_linea_1_default],
  [15758750456714168963n, nexon_mainnet_lith_default],
  [3849287863852499584n, bitcoin_mainnet_bob_1_default],
  [5214452172935136222n, treasure_mainnet_default],
  [12657445206920369324n, nexon_mainnet_henesys_default],
  [1294465214383781161n, berachain_mainnet_default],
  [9478124434908827753n, codex_mainnet_default],
  [4411394078118774322n, ethereum_mainnet_blast_1_default],
  [17912061998839310979n, plume_mainnet_default],
  [16468599424800719238n, ethereum_mainnet_taiko_1_default],
  [7937294810946806131n, bitcoin_mainnet_bitlayer_1_default],
  [5463201557265485081n, avalanche_subnet_dexalot_mainnet_default],
  [13204309965629103672n, ethereum_mainnet_scroll_1_default],
  [2459028469735686113n, polygon_mainnet_katana_default],
  [14632960069656270105n, nexon_qa_default],
  [4350319965322101699n, zklink_nova_mainnet_default],
  [5556806327594153475n, nexon_stage_default],
  [1010349088906777999n, ethereum_mainnet_arbitrum_1_treasure_1_default],
  [3555797439612589184n, zora_mainnet_default],
  [9043146809313071210n, corn_mainnet_default],
  [1546563616611573946n, tron_mainnet_evm_default],
  [124615329519749607n, solana_mainnet_default],
  [4741433654826277614n, aptos_mainnet_default],
  [17529533435026248318n, sui_mainnet_default],
  [16448340667252469081n, ton_mainnet_default],
  [1546563616611573945n, tron_mainnet_default]
]);
var testnetBySelector = new Map([
  [8953668971247136127n, bitcoin_testnet_rootstock_default],
  [729797994450396300n, telos_evm_testnet_default],
  [4340886533089894000n, polkadot_testnet_darwinia_pangoro_default],
  [3017758115101368649n, xdc_testnet_default],
  [8955032871639343000n, coinex_smart_chain_testnet_default],
  [6955638871347136141n, polkadot_testnet_astar_shibuya_default],
  [13264668187771770619n, binance_smart_chain_testnet_default],
  [572210378683744374n, velas_testnet_default],
  [4356164186791070119n, ethereum_testnet_sepolia_hashkey_1_default],
  [17833296867764334567n, shibarium_testnet_puppynet_default],
  [2066098519157881736n, ethereum_testnet_sepolia_xlayer_1_default],
  [16487132492576884721n, cronos_zkevm_testnet_sepolia_default],
  [6802309497652714138n, ethereum_testnet_goerli_zksync_1_default],
  [3842103497652714138n, cronos_testnet_zkevm_1_default],
  [222782988166878823n, hedera_testnet_default],
  [6898391096552792247n, ethereum_testnet_sepolia_zksync_1_default],
  [2995292832068775165n, cronos_testnet_default],
  [5061593697262339000n, near_testnet_default],
  [2664363617261496610n, ethereum_testnet_goerli_optimism_1_default],
  [7317911323415911000n, areon_testnet_default],
  [5059197667603797935n, janction_testnet_sepolia_default],
  [6260932437388305511n, private_testnet_obsidian_default],
  [829525985033418733n, ethereum_testnet_sepolia_mode_1_default],
  [2177900824115119161n, bittensor_testnet_default],
  [4286062357653186312n, hyperliquid_testnet_default],
  [2624132734533621656n, kaia_testnet_kairos_default],
  [4459371029167934217n, bittorrent_chain_testnet_default],
  [9284632837123596123n, wemix_testnet_default],
  [4264732132125536123n, core_testnet_default],
  [1948510578179542068n, bitcoin_testnet_bsquared_1_default],
  [5361632739113536121n, polkadot_testnet_moonbeam_moonbase_default],
  [14135854469784514356n, ethereum_testnet_sepolia_unichain_1_default],
  [1216300075444106652n, sei_testnet_atlantic_default],
  [3379446385462418246n, geth_testnet_default],
  [11059667695644972511n, ethereum_testnet_goerli_polygon_zkevm_1_default],
  [4237030917318060427n, story_testnet_default],
  [10749384167430721561n, mint_testnet_default],
  [6286293440461807648n, metal_testnet_default],
  [4888058894222120000n, bitcichain_testnet_default],
  [686603546605904534n, ethereum_testnet_sepolia_soneium_1_default],
  [13116810400804392105n, ronin_testnet_saigon_default],
  [3260900564719373474n, private_testnet_granite_default],
  [6915682381028791124n, private_testnet_andesite_default],
  [15513093881969820114n, dtcc_testnet_andesite_default],
  [2333097300889804761n, polkadot_testnet_centrifuge_altair_default],
  [12168171414969487009n, memento_testnet_default],
  [2110537777356199208n, kava_testnet_default],
  [5990477251245693094n, ethereum_testnet_sepolia_kroma_1_default],
  [9488606126177218005n, tac_testnet_default],
  [1654667687261492630n, ethereum_testnet_sepolia_polygon_zkevm_1_default],
  [8901520481741771655n, ethereum_testnet_holesky_fraxtal_1_default],
  [8304510386741731151n, ethereum_testnet_holesky_morph_1_default],
  [1467223411771711614n, bitcoin_testnet_botanix_default],
  [4905564228793744293n, fantom_testnet_default],
  [5298399861320400553n, ethereum_testnet_sepolia_lisk_1_default],
  [5299555114858065850n, ethereum_testnet_sepolia_worldchain_1_default],
  [4168263376276232250n, ethereum_testnet_goerli_mantle_1_default],
  [8236463271206331221n, ethereum_testnet_sepolia_mantle_1_default],
  [13274425992935471758n, binance_smart_chain_testnet_opbnb_1_default],
  [8911150974185440581n, nexon_dev_default],
  [2443239559770384419n, megaeth_testnet_default],
  [305104239123120457n, nibiru_testnet_default],
  [344208382356656551n, ondo_testnet_default],
  [1113014352258747600n, neonlink_testnet_default],
  [3967220077692964309n, plasma_testnet_default],
  [2183018362218727504n, monad_testnet_default],
  [8871595565390010547n, gnosis_chain_testnet_chiado_default],
  [16235373811196386733n, abstract_testnet_default],
  [3486622437121596122n, ethereum_testnet_sepolia_arbitrum_1_l3x_1_default],
  [4526165231216331901n, ethereum_testnet_sepolia_immutable_zkevm_1_default],
  [16088006396410204581n, _0g_testnet_newton_default],
  [2131427466778448014n, _0g_testnet_galileo_default],
  [7717148896336251131n, ethereum_testnet_holesky_default],
  [7759470850252068959n, anvil_devnet_default],
  [9900119385908781505n, apechain_testnet_curtis_default],
  [6827576821754315911n, ethereum_testnet_sepolia_lens_1_default],
  [14767482510784806043n, avalanche_testnet_fuji_default],
  [3552045678561919002n, celo_testnet_alfajores_default],
  [8446413392851542429n, private_testnet_opala_default],
  [13781831279385219069n, zircuit_testnet_garfield_default],
  [4562743618362911021n, ethereum_testnet_sepolia_zircuit_1_default],
  [13694007683517087973n, superseed_testnet_default],
  [3676871237479449268n, sonic_testnet_blaze_default],
  [1355246678561316402n, ethereum_testnet_goerli_linea_1_default],
  [5719461335882077547n, ethereum_testnet_sepolia_linea_1_default],
  [3777822886988675105n, ethereum_testnet_sepolia_metis_1_default],
  [12532609583862916517n, polygon_testnet_mumbai_default],
  [16281711391670634445n, polygon_testnet_amoy_default],
  [7728255861635209484n, berachain_testnet_bepolia_default],
  [8999465244383784164n, berachain_testnet_bartio_default],
  [12336603543561911511n, berachain_testnet_artio_default],
  [2285225387454015855n, zero_g_testnet_galileo_default],
  [5790810961207155433n, ethereum_testnet_goerli_base_1_default],
  [10344971235874465080n, ethereum_testnet_sepolia_base_1_default],
  [3743020999916460931n, plume_devnet_default],
  [13874588925447303949n, plume_testnet_sepolia_default],
  [1910019406958449359n, etherlink_testnet_default],
  [9090863410735740267n, polygon_testnet_tatara_default],
  [7248756420937879088n, ethereum_testnet_holesky_taiko_1_default],
  [7189150270347329685n, mind_testnet_default],
  [3789623672476206327n, bitcoin_testnet_bitlayer_1_default],
  [6101244977088475029n, ethereum_testnet_goerli_arbitrum_1_default],
  [3478487238524512106n, ethereum_testnet_sepolia_arbitrum_1_default],
  [4489326297382772450n, private_testnet_mica_default],
  [1458281248224512906n, avalanche_subnet_dexalot_testnet_default],
  [2279865765895943307n, ethereum_testnet_sepolia_scroll_1_default],
  [7837562506228496256n, avalanche_testnet_nexon_default],
  [5269261765892944301n, bitcoin_testnet_merlin_default],
  [4012524741200567430n, pharos_testnet_default],
  [4418231248214522936n, ethereum_testnet_sepolia_polygon_validium_1_default],
  [16126893759944359622n, hemi_testnet_sepolia_default],
  [9763904284804119144n, ink_testnet_sepolia_default],
  [5535534526963509396n, bitcoin_testnet_sepolia_bob_1_default],
  [5837261596322416298n, zklink_nova_testnet_default],
  [7225665875429174318n, codex_testnet_default],
  [10443705513486043421n, ethereum_testnet_sepolia_arbitrum_1_treasure_1_default],
  [3676916124122457866n, treasure_testnet_topaz_default],
  [945045181441419236n, jovay_testnet_default],
  [16015286601757825753n, ethereum_testnet_sepolia_default],
  [5224473277236331295n, ethereum_testnet_sepolia_optimism_1_default],
  [2217764097022649312n, neox_testnet_t4_default],
  [1467427327723633929n, ethereum_testnet_sepolia_corn_1_default],
  [7060342227814389000n, filecoin_testnet_default],
  [14684575664602284776n, plume_testnet_default],
  [2027362563942762617n, ethereum_testnet_sepolia_blast_1_default],
  [16244020411108056671n, zora_testnet_default],
  [13231703482326770598n, tron_testnet_shasta_evm_default],
  [13231703482326770600n, tron_devnet_evm_default],
  [2052925811360307749n, tron_testnet_nile_evm_default],
  [6302590918974934319n, solana_testnet_default],
  [16423721717087811551n, solana_devnet_default],
  [743186221051783445n, aptos_testnet_default],
  [4457093679053095497n, aptos_localnet_default],
  [9762610643973837292n, sui_testnet_default],
  [18395503381733958356n, sui_localnet_default],
  [1399300952838017768n, ton_testnet_default],
  [13879075125137744094n, ton_localnet_default],
  [13231703482326770597n, tron_testnet_shasta_default],
  [13231703482326770599n, tron_devnet_default],
  [2052925811360307740n, tron_testnet_nile_default]
]);
var mainnetByName = new Map([
  ["ethereum-mainnet", ethereum_mainnet_default],
  ["ethereum-mainnet-optimism-1", ethereum_mainnet_optimism_1_default],
  ["cronos-mainnet", cronos_mainnet_default],
  ["rootstock-mainnet", rootstock_mainnet_default],
  ["telos-evm-mainnet", telos_evm_mainnet_default],
  ["polkadot-mainnet-darwinia", polkadot_mainnet_darwinia_default],
  ["xdc-mainnet", xdc_mainnet_default],
  ["coinex_smart_chain-mainnet", coinex_smart_chain_mainnet_default],
  ["binance_smart_chain-mainnet", binance_smart_chain_mainnet_default],
  ["gnosis_chain-mainnet", gnosis_chain_mainnet_default],
  ["velas-mainnet", velas_mainnet_default],
  ["shibarium-mainnet", shibarium_mainnet_default],
  ["ethereum-mainnet-unichain-1", ethereum_mainnet_unichain_1_default],
  ["polygon-mainnet", polygon_mainnet_default],
  ["monad-mainnet", monad_mainnet_default],
  ["sonic-mainnet", sonic_mainnet_default],
  ["ethereum-mainnet-hashkey-1", ethereum_mainnet_hashkey_1_default],
  ["mint-mainnet", mint_mainnet_default],
  ["ethereum-mainnet-xlayer-1", ethereum_mainnet_xlayer_1_default],
  ["bittorrent_chain-mainnet", bittorrent_chain_mainnet_default],
  ["binance_smart_chain-mainnet-opbnb-1", binance_smart_chain_mainnet_opbnb_1_default],
  ["bitcoin-mainnet-bsquared-1", bitcoin_mainnet_bsquared_1_default],
  ["mind-mainnet", mind_mainnet_default],
  ["lens-mainnet", lens_mainnet_default],
  ["tac-mainnet", tac_mainnet_default],
  ["fantom-mainnet", fantom_mainnet_default],
  ["fraxtal-mainnet", fraxtal_mainnet_default],
  ["ethereum-mainnet-kroma-1", ethereum_mainnet_kroma_1_default],
  ["neonlink-mainnet", neonlink_mainnet_default],
  ["hedera-mainnet", hedera_mainnet_default],
  ["filecoin-mainnet", filecoin_mainnet_default],
  ["ethereum-mainnet-zksync-1", ethereum_mainnet_zksync_1_default],
  ["cronos-zkevm-mainnet", cronos_zkevm_mainnet_default],
  ["near-mainnet", near_mainnet_default],
  ["areon-mainnet", areon_mainnet_default],
  ["ethereum-mainnet-worldchain-1", ethereum_mainnet_worldchain_1_default],
  ["polkadot-mainnet-astar", polkadot_mainnet_astar_default],
  ["janction-mainnet", janction_mainnet_default],
  ["bittensor-mainnet", bittensor_mainnet_default],
  ["hyperliquid-mainnet", hyperliquid_mainnet_default],
  ["conflux-mainnet", conflux_mainnet_default],
  ["ethereum-mainnet-metis-1", ethereum_mainnet_metis_1_default],
  ["ethereum-mainnet-polygon-zkevm-1", ethereum_mainnet_polygon_zkevm_1_default],
  ["wemix-mainnet", wemix_mainnet_default],
  ["core-mainnet", core_mainnet_default],
  ["lisk-mainnet", lisk_mainnet_default],
  ["polkadot-mainnet-moonbeam", polkadot_mainnet_moonbeam_default],
  ["kusama-mainnet-moonriver", kusama_mainnet_moonriver_default],
  ["sei-mainnet", sei_mainnet_default],
  ["metal-mainnet", metal_mainnet_default],
  ["soneium-mainnet", soneium_mainnet_default],
  ["bitcichain-mainnet", bitcichain_mainnet_default],
  ["ronin-mainnet", ronin_mainnet_default],
  ["polkadot-mainnet-centrifuge", polkadot_mainnet_centrifuge_default],
  ["kava-mainnet", kava_mainnet_default],
  ["abstract-mainnet", abstract_mainnet_default],
  ["morph-mainnet", morph_mainnet_default],
  ["bitcoin-mainnet-botanix", bitcoin_mainnet_botanix_default],
  ["ethereum-mainnet-astar-zkevm-1", ethereum_mainnet_astar_zkevm_1_default],
  ["bitcoin-merlin-mainnet", bitcoin_merlin_mainnet_default],
  ["ethereum-mainnet-mantle-1", ethereum_mainnet_mantle_1_default],
  ["superseed-mainnet", superseed_mainnet_default],
  ["nibiru-mainnet", nibiru_mainnet_default],
  ["zetachain-mainnet", zetachain_mainnet_default],
  ["kaia-mainnet", kaia_mainnet_default],
  ["ethereum-mainnet-base-1", ethereum_mainnet_base_1_default],
  ["plasma-mainnet", plasma_mainnet_default],
  ["ethereum-mainnet-arbitrum-1-l3x-1", ethereum_mainnet_arbitrum_1_l3x_1_default],
  ["ethereum-mainnet-immutable-zkevm-1", ethereum_mainnet_immutable_zkevm_1_default],
  ["0g-mainnet", _0g_mainnet_default],
  ["apechain-mainnet", apechain_mainnet_default],
  ["ethereum-mainnet-mode-1", ethereum_mainnet_mode_1_default],
  ["ethereum-mainnet-arbitrum-1", ethereum_mainnet_arbitrum_1_default],
  ["celo-mainnet", celo_mainnet_default],
  ["etherlink-mainnet", etherlink_mainnet_default],
  ["hemi-mainnet", hemi_mainnet_default],
  ["avalanche-mainnet", avalanche_mainnet_default],
  ["neox-mainnet", neox_mainnet_default],
  ["ethereum-mainnet-zircuit-1", ethereum_mainnet_zircuit_1_default],
  ["memento-mainnet", memento_mainnet_default],
  ["ethereum-mainnet-ink-1", ethereum_mainnet_ink_1_default],
  ["ethereum-mainnet-linea-1", ethereum_mainnet_linea_1_default],
  ["nexon-mainnet-lith", nexon_mainnet_lith_default],
  ["bitcoin-mainnet-bob-1", bitcoin_mainnet_bob_1_default],
  ["treasure-mainnet", treasure_mainnet_default],
  ["nexon-mainnet-henesys", nexon_mainnet_henesys_default],
  ["berachain-mainnet", berachain_mainnet_default],
  ["codex-mainnet", codex_mainnet_default],
  ["ethereum-mainnet-blast-1", ethereum_mainnet_blast_1_default],
  ["plume-mainnet", plume_mainnet_default],
  ["ethereum-mainnet-taiko-1", ethereum_mainnet_taiko_1_default],
  ["bitcoin-mainnet-bitlayer-1", bitcoin_mainnet_bitlayer_1_default],
  ["avalanche-subnet-dexalot-mainnet", avalanche_subnet_dexalot_mainnet_default],
  ["ethereum-mainnet-scroll-1", ethereum_mainnet_scroll_1_default],
  ["polygon-mainnet-katana", polygon_mainnet_katana_default],
  ["nexon-qa", nexon_qa_default],
  ["zklink_nova-mainnet", zklink_nova_mainnet_default],
  ["nexon-stage", nexon_stage_default],
  ["ethereum-mainnet-arbitrum-1-treasure-1", ethereum_mainnet_arbitrum_1_treasure_1_default],
  ["zora-mainnet", zora_mainnet_default],
  ["corn-mainnet", corn_mainnet_default],
  ["tron-mainnet-evm", tron_mainnet_evm_default],
  ["solana-mainnet", solana_mainnet_default],
  ["aptos-mainnet", aptos_mainnet_default],
  ["sui-mainnet", sui_mainnet_default],
  ["ton-mainnet", ton_mainnet_default],
  ["tron-mainnet", tron_mainnet_default]
]);
var testnetByName = new Map([
  ["bitcoin-testnet-rootstock", bitcoin_testnet_rootstock_default],
  ["telos-evm-testnet", telos_evm_testnet_default],
  ["polkadot-testnet-darwinia-pangoro", polkadot_testnet_darwinia_pangoro_default],
  ["xdc-testnet", xdc_testnet_default],
  ["coinex_smart_chain-testnet", coinex_smart_chain_testnet_default],
  ["polkadot-testnet-astar-shibuya", polkadot_testnet_astar_shibuya_default],
  ["binance_smart_chain-testnet", binance_smart_chain_testnet_default],
  ["velas-testnet", velas_testnet_default],
  ["ethereum-testnet-sepolia-hashkey-1", ethereum_testnet_sepolia_hashkey_1_default],
  ["shibarium-testnet-puppynet", shibarium_testnet_puppynet_default],
  ["ethereum-testnet-sepolia-xlayer-1", ethereum_testnet_sepolia_xlayer_1_default],
  ["cronos-zkevm-testnet-sepolia", cronos_zkevm_testnet_sepolia_default],
  ["ethereum-testnet-goerli-zksync-1", ethereum_testnet_goerli_zksync_1_default],
  ["cronos-testnet-zkevm-1", cronos_testnet_zkevm_1_default],
  ["hedera-testnet", hedera_testnet_default],
  ["ethereum-testnet-sepolia-zksync-1", ethereum_testnet_sepolia_zksync_1_default],
  ["cronos-testnet", cronos_testnet_default],
  ["near-testnet", near_testnet_default],
  ["ethereum-testnet-goerli-optimism-1", ethereum_testnet_goerli_optimism_1_default],
  ["areon-testnet", areon_testnet_default],
  ["janction-testnet-sepolia", janction_testnet_sepolia_default],
  ["private-testnet-obsidian", private_testnet_obsidian_default],
  ["ethereum-testnet-sepolia-mode-1", ethereum_testnet_sepolia_mode_1_default],
  ["bittensor-testnet", bittensor_testnet_default],
  ["hyperliquid-testnet", hyperliquid_testnet_default],
  ["kaia-testnet-kairos", kaia_testnet_kairos_default],
  ["bittorrent_chain-testnet", bittorrent_chain_testnet_default],
  ["wemix-testnet", wemix_testnet_default],
  ["core-testnet", core_testnet_default],
  ["bitcoin-testnet-bsquared-1", bitcoin_testnet_bsquared_1_default],
  ["polkadot-testnet-moonbeam-moonbase", polkadot_testnet_moonbeam_moonbase_default],
  ["ethereum-testnet-sepolia-unichain-1", ethereum_testnet_sepolia_unichain_1_default],
  ["sei-testnet-atlantic", sei_testnet_atlantic_default],
  ["geth-testnet", geth_testnet_default],
  ["ethereum-testnet-goerli-polygon-zkevm-1", ethereum_testnet_goerli_polygon_zkevm_1_default],
  ["story-testnet", story_testnet_default],
  ["mint-testnet", mint_testnet_default],
  ["metal-testnet", metal_testnet_default],
  ["bitcichain-testnet", bitcichain_testnet_default],
  ["ethereum-testnet-sepolia-soneium-1", ethereum_testnet_sepolia_soneium_1_default],
  ["ronin-testnet-saigon", ronin_testnet_saigon_default],
  ["private-testnet-granite", private_testnet_granite_default],
  ["private-testnet-andesite", private_testnet_andesite_default],
  ["dtcc-testnet-andesite", dtcc_testnet_andesite_default],
  ["polkadot-testnet-centrifuge-altair", polkadot_testnet_centrifuge_altair_default],
  ["memento-testnet", memento_testnet_default],
  ["kava-testnet", kava_testnet_default],
  ["ethereum-testnet-sepolia-kroma-1", ethereum_testnet_sepolia_kroma_1_default],
  ["tac-testnet", tac_testnet_default],
  [
    "ethereum-testnet-sepolia-polygon-zkevm-1",
    ethereum_testnet_sepolia_polygon_zkevm_1_default
  ],
  ["ethereum-testnet-holesky-fraxtal-1", ethereum_testnet_holesky_fraxtal_1_default],
  ["ethereum-testnet-holesky-morph-1", ethereum_testnet_holesky_morph_1_default],
  ["bitcoin-testnet-botanix", bitcoin_testnet_botanix_default],
  ["fantom-testnet", fantom_testnet_default],
  ["ethereum-testnet-sepolia-lisk-1", ethereum_testnet_sepolia_lisk_1_default],
  ["ethereum-testnet-sepolia-worldchain-1", ethereum_testnet_sepolia_worldchain_1_default],
  ["ethereum-testnet-goerli-mantle-1", ethereum_testnet_goerli_mantle_1_default],
  ["ethereum-testnet-sepolia-mantle-1", ethereum_testnet_sepolia_mantle_1_default],
  ["binance_smart_chain-testnet-opbnb-1", binance_smart_chain_testnet_opbnb_1_default],
  ["nexon-dev", nexon_dev_default],
  ["megaeth-testnet", megaeth_testnet_default],
  ["nibiru-testnet", nibiru_testnet_default],
  ["ondo-testnet", ondo_testnet_default],
  ["neonlink-testnet", neonlink_testnet_default],
  ["plasma-testnet", plasma_testnet_default],
  ["monad-testnet", monad_testnet_default],
  ["gnosis_chain-testnet-chiado", gnosis_chain_testnet_chiado_default],
  ["abstract-testnet", abstract_testnet_default],
  [
    "ethereum-testnet-sepolia-arbitrum-1-l3x-1",
    ethereum_testnet_sepolia_arbitrum_1_l3x_1_default
  ],
  [
    "ethereum-testnet-sepolia-immutable-zkevm-1",
    ethereum_testnet_sepolia_immutable_zkevm_1_default
  ],
  ["0g-testnet-newton", _0g_testnet_newton_default],
  ["0g-testnet-galileo", _0g_testnet_galileo_default],
  ["ethereum-testnet-holesky", ethereum_testnet_holesky_default],
  ["anvil-devnet", anvil_devnet_default],
  ["apechain-testnet-curtis", apechain_testnet_curtis_default],
  ["ethereum-testnet-sepolia-lens-1", ethereum_testnet_sepolia_lens_1_default],
  ["avalanche-testnet-fuji", avalanche_testnet_fuji_default],
  ["celo-testnet-alfajores", celo_testnet_alfajores_default],
  ["private-testnet-opala", private_testnet_opala_default],
  ["zircuit-testnet-garfield", zircuit_testnet_garfield_default],
  ["ethereum-testnet-sepolia-zircuit-1", ethereum_testnet_sepolia_zircuit_1_default],
  ["superseed-testnet", superseed_testnet_default],
  ["sonic-testnet-blaze", sonic_testnet_blaze_default],
  ["ethereum-testnet-goerli-linea-1", ethereum_testnet_goerli_linea_1_default],
  ["ethereum-testnet-sepolia-linea-1", ethereum_testnet_sepolia_linea_1_default],
  ["ethereum-testnet-sepolia-metis-1", ethereum_testnet_sepolia_metis_1_default],
  ["polygon-testnet-mumbai", polygon_testnet_mumbai_default],
  ["polygon-testnet-amoy", polygon_testnet_amoy_default],
  ["berachain-testnet-bepolia", berachain_testnet_bepolia_default],
  ["berachain-testnet-bartio", berachain_testnet_bartio_default],
  ["berachain-testnet-artio", berachain_testnet_artio_default],
  ["zero-g-testnet-galileo", zero_g_testnet_galileo_default],
  ["ethereum-testnet-goerli-base-1", ethereum_testnet_goerli_base_1_default],
  ["ethereum-testnet-sepolia-base-1", ethereum_testnet_sepolia_base_1_default],
  ["plume-devnet", plume_devnet_default],
  ["plume-testnet-sepolia", plume_testnet_sepolia_default],
  ["etherlink-testnet", etherlink_testnet_default],
  ["polygon-testnet-tatara", polygon_testnet_tatara_default],
  ["ethereum-testnet-holesky-taiko-1", ethereum_testnet_holesky_taiko_1_default],
  ["mind-testnet", mind_testnet_default],
  ["bitcoin-testnet-bitlayer-1", bitcoin_testnet_bitlayer_1_default],
  ["ethereum-testnet-goerli-arbitrum-1", ethereum_testnet_goerli_arbitrum_1_default],
  ["ethereum-testnet-sepolia-arbitrum-1", ethereum_testnet_sepolia_arbitrum_1_default],
  ["private-testnet-mica", private_testnet_mica_default],
  ["avalanche-subnet-dexalot-testnet", avalanche_subnet_dexalot_testnet_default],
  ["ethereum-testnet-sepolia-scroll-1", ethereum_testnet_sepolia_scroll_1_default],
  ["avalanche-testnet-nexon", avalanche_testnet_nexon_default],
  ["bitcoin-testnet-merlin", bitcoin_testnet_merlin_default],
  ["pharos-testnet", pharos_testnet_default],
  [
    "ethereum-testnet-sepolia-polygon-validium-1",
    ethereum_testnet_sepolia_polygon_validium_1_default
  ],
  ["hemi-testnet-sepolia", hemi_testnet_sepolia_default],
  ["ink-testnet-sepolia", ink_testnet_sepolia_default],
  ["bitcoin-testnet-sepolia-bob-1", bitcoin_testnet_sepolia_bob_1_default],
  ["zklink_nova-testnet", zklink_nova_testnet_default],
  ["codex-testnet", codex_testnet_default],
  [
    "ethereum-testnet-sepolia-arbitrum-1-treasure-1",
    ethereum_testnet_sepolia_arbitrum_1_treasure_1_default
  ],
  ["treasure-testnet-topaz", treasure_testnet_topaz_default],
  ["jovay-testnet", jovay_testnet_default],
  ["ethereum-testnet-sepolia", ethereum_testnet_sepolia_default],
  ["ethereum-testnet-sepolia-optimism-1", ethereum_testnet_sepolia_optimism_1_default],
  ["neox-testnet-t4", neox_testnet_t4_default],
  ["ethereum-testnet-sepolia-corn-1", ethereum_testnet_sepolia_corn_1_default],
  ["filecoin-testnet", filecoin_testnet_default],
  ["plume-testnet", plume_testnet_default],
  ["ethereum-testnet-sepolia-blast-1", ethereum_testnet_sepolia_blast_1_default],
  ["zora-testnet", zora_testnet_default],
  ["tron-testnet-shasta-evm", tron_testnet_shasta_evm_default],
  ["tron-devnet-evm", tron_devnet_evm_default],
  ["tron-testnet-nile-evm", tron_testnet_nile_evm_default],
  ["solana-testnet", solana_testnet_default],
  ["solana-devnet", solana_devnet_default],
  ["aptos-testnet", aptos_testnet_default],
  ["aptos-localnet", aptos_localnet_default],
  ["sui-testnet", sui_testnet_default],
  ["sui-localnet", sui_localnet_default],
  ["ton-testnet", ton_testnet_default],
  ["ton-localnet", ton_localnet_default],
  ["tron-testnet-shasta", tron_testnet_shasta_default],
  ["tron-devnet", tron_devnet_default],
  ["tron-testnet-nile", tron_testnet_nile_default]
]);
var mainnetBySelectorByFamily = {
  evm: new Map([
    [5009297550715157269n, ethereum_mainnet_default],
    [3734403246176062136n, ethereum_mainnet_optimism_1_default],
    [1456215246176062136n, cronos_mainnet_default],
    [11964252391146578476n, rootstock_mainnet_default],
    [1477345371608778000n, telos_evm_mainnet_default],
    [8866418665544333000n, polkadot_mainnet_darwinia_default],
    [17673274061779414707n, xdc_mainnet_default],
    [1761333065194157300n, coinex_smart_chain_mainnet_default],
    [11344663589394136015n, binance_smart_chain_mainnet_default],
    [465200170687744372n, gnosis_chain_mainnet_default],
    [374210358663784372n, velas_mainnet_default],
    [3993510008929295315n, shibarium_mainnet_default],
    [1923510103922296319n, ethereum_mainnet_unichain_1_default],
    [4051577828743386545n, polygon_mainnet_default],
    [8481857512324358265n, monad_mainnet_default],
    [1673871237479749969n, sonic_mainnet_default],
    [7613811247471741961n, ethereum_mainnet_hashkey_1_default],
    [17164792800244661392n, mint_mainnet_default],
    [3016212468291539606n, ethereum_mainnet_xlayer_1_default],
    [3776006016387883143n, bittorrent_chain_mainnet_default],
    [465944652040885897n, binance_smart_chain_mainnet_opbnb_1_default],
    [5406759801798337480n, bitcoin_mainnet_bsquared_1_default],
    [11690709103138290329n, mind_mainnet_default],
    [5608378062013572713n, lens_mainnet_default],
    [5936861837188149645n, tac_mainnet_default],
    [3768048213127883732n, fantom_mainnet_default],
    [1462016016387883143n, fraxtal_mainnet_default],
    [3719320017875267166n, ethereum_mainnet_kroma_1_default],
    [8239338020728974000n, neonlink_mainnet_default],
    [3229138320728879060n, hedera_mainnet_default],
    [4561443241176882990n, filecoin_mainnet_default],
    [1562403441176082196n, ethereum_mainnet_zksync_1_default],
    [8788096068760390840n, cronos_zkevm_mainnet_default],
    [2039744413822257700n, near_mainnet_default],
    [1939936305787790600n, areon_mainnet_default],
    [2049429975587534727n, ethereum_mainnet_worldchain_1_default],
    [6422105447186081193n, polkadot_mainnet_astar_default],
    [9107126442626377432n, janction_mainnet_default],
    [2135107236357186872n, bittensor_mainnet_default],
    [2442541497099098535n, hyperliquid_mainnet_default],
    [3358365939762719202n, conflux_mainnet_default],
    [8805746078405598895n, ethereum_mainnet_metis_1_default],
    [4348158687435793198n, ethereum_mainnet_polygon_zkevm_1_default],
    [5142893604156789321n, wemix_mainnet_default],
    [1224752112135636129n, core_mainnet_default],
    [15293031020466096408n, lisk_mainnet_default],
    [1252863800116739621n, polkadot_mainnet_moonbeam_default],
    [1355020143337428062n, kusama_mainnet_moonriver_default],
    [9027416829622342829n, sei_mainnet_default],
    [13447077090413146373n, metal_mainnet_default],
    [12505351618335765396n, soneium_mainnet_default],
    [4874388048629246000n, bitcichain_mainnet_default],
    [6916147374840168594n, ronin_mainnet_default],
    [8175830712062617656n, polkadot_mainnet_centrifuge_default],
    [7550000543357438061n, kava_mainnet_default],
    [3577778157919314504n, abstract_mainnet_default],
    [18164309074156128038n, morph_mainnet_default],
    [4560701533377838164n, bitcoin_mainnet_botanix_default],
    [1540201334317828111n, ethereum_mainnet_astar_zkevm_1_default],
    [241851231317828981n, bitcoin_merlin_mainnet_default],
    [1556008542357238666n, ethereum_mainnet_mantle_1_default],
    [470401360549526817n, superseed_mainnet_default],
    [17349189558768828726n, nibiru_mainnet_default],
    [10817664450262215148n, zetachain_mainnet_default],
    [9813823125703490621n, kaia_mainnet_default],
    [15971525489660198786n, ethereum_mainnet_base_1_default],
    [9335212494177455608n, plasma_mainnet_default],
    [3162193654116181371n, ethereum_mainnet_arbitrum_1_l3x_1_default],
    [1237925231416731909n, ethereum_mainnet_immutable_zkevm_1_default],
    [4426351306075016396n, _0g_mainnet_default],
    [14894068710063348487n, apechain_mainnet_default],
    [7264351850409363825n, ethereum_mainnet_mode_1_default],
    [4949039107694359620n, ethereum_mainnet_arbitrum_1_default],
    [1346049177634351622n, celo_mainnet_default],
    [13624601974233774587n, etherlink_mainnet_default],
    [1804312132722180201n, hemi_mainnet_default],
    [6433500567565415381n, avalanche_mainnet_default],
    [7222032299962346917n, neox_mainnet_default],
    [17198166215261833993n, ethereum_mainnet_zircuit_1_default],
    [6473245816409426016n, memento_mainnet_default],
    [3461204551265785888n, ethereum_mainnet_ink_1_default],
    [4627098889531055414n, ethereum_mainnet_linea_1_default],
    [15758750456714168963n, nexon_mainnet_lith_default],
    [3849287863852499584n, bitcoin_mainnet_bob_1_default],
    [5214452172935136222n, treasure_mainnet_default],
    [12657445206920369324n, nexon_mainnet_henesys_default],
    [1294465214383781161n, berachain_mainnet_default],
    [9478124434908827753n, codex_mainnet_default],
    [4411394078118774322n, ethereum_mainnet_blast_1_default],
    [17912061998839310979n, plume_mainnet_default],
    [16468599424800719238n, ethereum_mainnet_taiko_1_default],
    [7937294810946806131n, bitcoin_mainnet_bitlayer_1_default],
    [5463201557265485081n, avalanche_subnet_dexalot_mainnet_default],
    [13204309965629103672n, ethereum_mainnet_scroll_1_default],
    [2459028469735686113n, polygon_mainnet_katana_default],
    [14632960069656270105n, nexon_qa_default],
    [4350319965322101699n, zklink_nova_mainnet_default],
    [5556806327594153475n, nexon_stage_default],
    [1010349088906777999n, ethereum_mainnet_arbitrum_1_treasure_1_default],
    [3555797439612589184n, zora_mainnet_default],
    [9043146809313071210n, corn_mainnet_default],
    [1546563616611573946n, tron_mainnet_evm_default]
  ]),
  solana: new Map([[124615329519749607n, solana_mainnet_default]]),
  aptos: new Map([[4741433654826277614n, aptos_mainnet_default]]),
  sui: new Map([[17529533435026248318n, sui_mainnet_default]]),
  ton: new Map([[16448340667252469081n, ton_mainnet_default]]),
  tron: new Map([[1546563616611573945n, tron_mainnet_default]])
};
var testnetBySelectorByFamily = {
  evm: new Map([
    [8953668971247136127n, bitcoin_testnet_rootstock_default],
    [729797994450396300n, telos_evm_testnet_default],
    [4340886533089894000n, polkadot_testnet_darwinia_pangoro_default],
    [3017758115101368649n, xdc_testnet_default],
    [8955032871639343000n, coinex_smart_chain_testnet_default],
    [6955638871347136141n, polkadot_testnet_astar_shibuya_default],
    [13264668187771770619n, binance_smart_chain_testnet_default],
    [572210378683744374n, velas_testnet_default],
    [4356164186791070119n, ethereum_testnet_sepolia_hashkey_1_default],
    [17833296867764334567n, shibarium_testnet_puppynet_default],
    [2066098519157881736n, ethereum_testnet_sepolia_xlayer_1_default],
    [16487132492576884721n, cronos_zkevm_testnet_sepolia_default],
    [6802309497652714138n, ethereum_testnet_goerli_zksync_1_default],
    [3842103497652714138n, cronos_testnet_zkevm_1_default],
    [222782988166878823n, hedera_testnet_default],
    [6898391096552792247n, ethereum_testnet_sepolia_zksync_1_default],
    [2995292832068775165n, cronos_testnet_default],
    [5061593697262339000n, near_testnet_default],
    [2664363617261496610n, ethereum_testnet_goerli_optimism_1_default],
    [7317911323415911000n, areon_testnet_default],
    [5059197667603797935n, janction_testnet_sepolia_default],
    [6260932437388305511n, private_testnet_obsidian_default],
    [829525985033418733n, ethereum_testnet_sepolia_mode_1_default],
    [2177900824115119161n, bittensor_testnet_default],
    [4286062357653186312n, hyperliquid_testnet_default],
    [2624132734533621656n, kaia_testnet_kairos_default],
    [4459371029167934217n, bittorrent_chain_testnet_default],
    [9284632837123596123n, wemix_testnet_default],
    [4264732132125536123n, core_testnet_default],
    [1948510578179542068n, bitcoin_testnet_bsquared_1_default],
    [5361632739113536121n, polkadot_testnet_moonbeam_moonbase_default],
    [14135854469784514356n, ethereum_testnet_sepolia_unichain_1_default],
    [1216300075444106652n, sei_testnet_atlantic_default],
    [3379446385462418246n, geth_testnet_default],
    [11059667695644972511n, ethereum_testnet_goerli_polygon_zkevm_1_default],
    [4237030917318060427n, story_testnet_default],
    [10749384167430721561n, mint_testnet_default],
    [6286293440461807648n, metal_testnet_default],
    [4888058894222120000n, bitcichain_testnet_default],
    [686603546605904534n, ethereum_testnet_sepolia_soneium_1_default],
    [13116810400804392105n, ronin_testnet_saigon_default],
    [3260900564719373474n, private_testnet_granite_default],
    [6915682381028791124n, private_testnet_andesite_default],
    [15513093881969820114n, dtcc_testnet_andesite_default],
    [2333097300889804761n, polkadot_testnet_centrifuge_altair_default],
    [12168171414969487009n, memento_testnet_default],
    [2110537777356199208n, kava_testnet_default],
    [5990477251245693094n, ethereum_testnet_sepolia_kroma_1_default],
    [9488606126177218005n, tac_testnet_default],
    [1654667687261492630n, ethereum_testnet_sepolia_polygon_zkevm_1_default],
    [8901520481741771655n, ethereum_testnet_holesky_fraxtal_1_default],
    [8304510386741731151n, ethereum_testnet_holesky_morph_1_default],
    [1467223411771711614n, bitcoin_testnet_botanix_default],
    [4905564228793744293n, fantom_testnet_default],
    [5298399861320400553n, ethereum_testnet_sepolia_lisk_1_default],
    [5299555114858065850n, ethereum_testnet_sepolia_worldchain_1_default],
    [4168263376276232250n, ethereum_testnet_goerli_mantle_1_default],
    [8236463271206331221n, ethereum_testnet_sepolia_mantle_1_default],
    [13274425992935471758n, binance_smart_chain_testnet_opbnb_1_default],
    [8911150974185440581n, nexon_dev_default],
    [2443239559770384419n, megaeth_testnet_default],
    [305104239123120457n, nibiru_testnet_default],
    [344208382356656551n, ondo_testnet_default],
    [1113014352258747600n, neonlink_testnet_default],
    [3967220077692964309n, plasma_testnet_default],
    [2183018362218727504n, monad_testnet_default],
    [8871595565390010547n, gnosis_chain_testnet_chiado_default],
    [16235373811196386733n, abstract_testnet_default],
    [3486622437121596122n, ethereum_testnet_sepolia_arbitrum_1_l3x_1_default],
    [4526165231216331901n, ethereum_testnet_sepolia_immutable_zkevm_1_default],
    [16088006396410204581n, _0g_testnet_newton_default],
    [2131427466778448014n, _0g_testnet_galileo_default],
    [7717148896336251131n, ethereum_testnet_holesky_default],
    [7759470850252068959n, anvil_devnet_default],
    [9900119385908781505n, apechain_testnet_curtis_default],
    [6827576821754315911n, ethereum_testnet_sepolia_lens_1_default],
    [14767482510784806043n, avalanche_testnet_fuji_default],
    [3552045678561919002n, celo_testnet_alfajores_default],
    [8446413392851542429n, private_testnet_opala_default],
    [13781831279385219069n, zircuit_testnet_garfield_default],
    [4562743618362911021n, ethereum_testnet_sepolia_zircuit_1_default],
    [13694007683517087973n, superseed_testnet_default],
    [3676871237479449268n, sonic_testnet_blaze_default],
    [1355246678561316402n, ethereum_testnet_goerli_linea_1_default],
    [5719461335882077547n, ethereum_testnet_sepolia_linea_1_default],
    [3777822886988675105n, ethereum_testnet_sepolia_metis_1_default],
    [12532609583862916517n, polygon_testnet_mumbai_default],
    [16281711391670634445n, polygon_testnet_amoy_default],
    [7728255861635209484n, berachain_testnet_bepolia_default],
    [8999465244383784164n, berachain_testnet_bartio_default],
    [12336603543561911511n, berachain_testnet_artio_default],
    [2285225387454015855n, zero_g_testnet_galileo_default],
    [5790810961207155433n, ethereum_testnet_goerli_base_1_default],
    [10344971235874465080n, ethereum_testnet_sepolia_base_1_default],
    [3743020999916460931n, plume_devnet_default],
    [13874588925447303949n, plume_testnet_sepolia_default],
    [1910019406958449359n, etherlink_testnet_default],
    [9090863410735740267n, polygon_testnet_tatara_default],
    [7248756420937879088n, ethereum_testnet_holesky_taiko_1_default],
    [7189150270347329685n, mind_testnet_default],
    [3789623672476206327n, bitcoin_testnet_bitlayer_1_default],
    [6101244977088475029n, ethereum_testnet_goerli_arbitrum_1_default],
    [3478487238524512106n, ethereum_testnet_sepolia_arbitrum_1_default],
    [4489326297382772450n, private_testnet_mica_default],
    [1458281248224512906n, avalanche_subnet_dexalot_testnet_default],
    [2279865765895943307n, ethereum_testnet_sepolia_scroll_1_default],
    [7837562506228496256n, avalanche_testnet_nexon_default],
    [5269261765892944301n, bitcoin_testnet_merlin_default],
    [4012524741200567430n, pharos_testnet_default],
    [4418231248214522936n, ethereum_testnet_sepolia_polygon_validium_1_default],
    [16126893759944359622n, hemi_testnet_sepolia_default],
    [9763904284804119144n, ink_testnet_sepolia_default],
    [5535534526963509396n, bitcoin_testnet_sepolia_bob_1_default],
    [5837261596322416298n, zklink_nova_testnet_default],
    [7225665875429174318n, codex_testnet_default],
    [10443705513486043421n, ethereum_testnet_sepolia_arbitrum_1_treasure_1_default],
    [3676916124122457866n, treasure_testnet_topaz_default],
    [945045181441419236n, jovay_testnet_default],
    [16015286601757825753n, ethereum_testnet_sepolia_default],
    [5224473277236331295n, ethereum_testnet_sepolia_optimism_1_default],
    [2217764097022649312n, neox_testnet_t4_default],
    [1467427327723633929n, ethereum_testnet_sepolia_corn_1_default],
    [7060342227814389000n, filecoin_testnet_default],
    [14684575664602284776n, plume_testnet_default],
    [2027362563942762617n, ethereum_testnet_sepolia_blast_1_default],
    [16244020411108056671n, zora_testnet_default],
    [13231703482326770598n, tron_testnet_shasta_evm_default],
    [13231703482326770600n, tron_devnet_evm_default],
    [2052925811360307749n, tron_testnet_nile_evm_default]
  ]),
  solana: new Map([
    [6302590918974934319n, solana_testnet_default],
    [16423721717087811551n, solana_devnet_default]
  ]),
  aptos: new Map([
    [743186221051783445n, aptos_testnet_default],
    [4457093679053095497n, aptos_localnet_default]
  ]),
  sui: new Map([
    [9762610643973837292n, sui_testnet_default],
    [18395503381733958356n, sui_localnet_default]
  ]),
  ton: new Map([
    [1399300952838017768n, ton_testnet_default],
    [13879075125137744094n, ton_localnet_default]
  ]),
  tron: new Map([
    [13231703482326770597n, tron_testnet_shasta_default],
    [13231703482326770599n, tron_devnet_default],
    [2052925811360307740n, tron_testnet_nile_default]
  ])
};
var mainnetByNameByFamily = {
  evm: new Map([
    ["ethereum-mainnet", ethereum_mainnet_default],
    ["ethereum-mainnet-optimism-1", ethereum_mainnet_optimism_1_default],
    ["cronos-mainnet", cronos_mainnet_default],
    ["rootstock-mainnet", rootstock_mainnet_default],
    ["telos-evm-mainnet", telos_evm_mainnet_default],
    ["polkadot-mainnet-darwinia", polkadot_mainnet_darwinia_default],
    ["xdc-mainnet", xdc_mainnet_default],
    ["coinex_smart_chain-mainnet", coinex_smart_chain_mainnet_default],
    ["binance_smart_chain-mainnet", binance_smart_chain_mainnet_default],
    ["gnosis_chain-mainnet", gnosis_chain_mainnet_default],
    ["velas-mainnet", velas_mainnet_default],
    ["shibarium-mainnet", shibarium_mainnet_default],
    ["ethereum-mainnet-unichain-1", ethereum_mainnet_unichain_1_default],
    ["polygon-mainnet", polygon_mainnet_default],
    ["monad-mainnet", monad_mainnet_default],
    ["sonic-mainnet", sonic_mainnet_default],
    ["ethereum-mainnet-hashkey-1", ethereum_mainnet_hashkey_1_default],
    ["mint-mainnet", mint_mainnet_default],
    ["ethereum-mainnet-xlayer-1", ethereum_mainnet_xlayer_1_default],
    ["bittorrent_chain-mainnet", bittorrent_chain_mainnet_default],
    ["binance_smart_chain-mainnet-opbnb-1", binance_smart_chain_mainnet_opbnb_1_default],
    ["bitcoin-mainnet-bsquared-1", bitcoin_mainnet_bsquared_1_default],
    ["mind-mainnet", mind_mainnet_default],
    ["lens-mainnet", lens_mainnet_default],
    ["tac-mainnet", tac_mainnet_default],
    ["fantom-mainnet", fantom_mainnet_default],
    ["fraxtal-mainnet", fraxtal_mainnet_default],
    ["ethereum-mainnet-kroma-1", ethereum_mainnet_kroma_1_default],
    ["neonlink-mainnet", neonlink_mainnet_default],
    ["hedera-mainnet", hedera_mainnet_default],
    ["filecoin-mainnet", filecoin_mainnet_default],
    ["ethereum-mainnet-zksync-1", ethereum_mainnet_zksync_1_default],
    ["cronos-zkevm-mainnet", cronos_zkevm_mainnet_default],
    ["near-mainnet", near_mainnet_default],
    ["areon-mainnet", areon_mainnet_default],
    ["ethereum-mainnet-worldchain-1", ethereum_mainnet_worldchain_1_default],
    ["polkadot-mainnet-astar", polkadot_mainnet_astar_default],
    ["janction-mainnet", janction_mainnet_default],
    ["bittensor-mainnet", bittensor_mainnet_default],
    ["hyperliquid-mainnet", hyperliquid_mainnet_default],
    ["conflux-mainnet", conflux_mainnet_default],
    ["ethereum-mainnet-metis-1", ethereum_mainnet_metis_1_default],
    ["ethereum-mainnet-polygon-zkevm-1", ethereum_mainnet_polygon_zkevm_1_default],
    ["wemix-mainnet", wemix_mainnet_default],
    ["core-mainnet", core_mainnet_default],
    ["lisk-mainnet", lisk_mainnet_default],
    ["polkadot-mainnet-moonbeam", polkadot_mainnet_moonbeam_default],
    ["kusama-mainnet-moonriver", kusama_mainnet_moonriver_default],
    ["sei-mainnet", sei_mainnet_default],
    ["metal-mainnet", metal_mainnet_default],
    ["soneium-mainnet", soneium_mainnet_default],
    ["bitcichain-mainnet", bitcichain_mainnet_default],
    ["ronin-mainnet", ronin_mainnet_default],
    ["polkadot-mainnet-centrifuge", polkadot_mainnet_centrifuge_default],
    ["kava-mainnet", kava_mainnet_default],
    ["abstract-mainnet", abstract_mainnet_default],
    ["morph-mainnet", morph_mainnet_default],
    ["bitcoin-mainnet-botanix", bitcoin_mainnet_botanix_default],
    ["ethereum-mainnet-astar-zkevm-1", ethereum_mainnet_astar_zkevm_1_default],
    ["bitcoin-merlin-mainnet", bitcoin_merlin_mainnet_default],
    ["ethereum-mainnet-mantle-1", ethereum_mainnet_mantle_1_default],
    ["superseed-mainnet", superseed_mainnet_default],
    ["nibiru-mainnet", nibiru_mainnet_default],
    ["zetachain-mainnet", zetachain_mainnet_default],
    ["kaia-mainnet", kaia_mainnet_default],
    ["ethereum-mainnet-base-1", ethereum_mainnet_base_1_default],
    ["plasma-mainnet", plasma_mainnet_default],
    ["ethereum-mainnet-arbitrum-1-l3x-1", ethereum_mainnet_arbitrum_1_l3x_1_default],
    ["ethereum-mainnet-immutable-zkevm-1", ethereum_mainnet_immutable_zkevm_1_default],
    ["0g-mainnet", _0g_mainnet_default],
    ["apechain-mainnet", apechain_mainnet_default],
    ["ethereum-mainnet-mode-1", ethereum_mainnet_mode_1_default],
    ["ethereum-mainnet-arbitrum-1", ethereum_mainnet_arbitrum_1_default],
    ["celo-mainnet", celo_mainnet_default],
    ["etherlink-mainnet", etherlink_mainnet_default],
    ["hemi-mainnet", hemi_mainnet_default],
    ["avalanche-mainnet", avalanche_mainnet_default],
    ["neox-mainnet", neox_mainnet_default],
    ["ethereum-mainnet-zircuit-1", ethereum_mainnet_zircuit_1_default],
    ["memento-mainnet", memento_mainnet_default],
    ["ethereum-mainnet-ink-1", ethereum_mainnet_ink_1_default],
    ["ethereum-mainnet-linea-1", ethereum_mainnet_linea_1_default],
    ["nexon-mainnet-lith", nexon_mainnet_lith_default],
    ["bitcoin-mainnet-bob-1", bitcoin_mainnet_bob_1_default],
    ["treasure-mainnet", treasure_mainnet_default],
    ["nexon-mainnet-henesys", nexon_mainnet_henesys_default],
    ["berachain-mainnet", berachain_mainnet_default],
    ["codex-mainnet", codex_mainnet_default],
    ["ethereum-mainnet-blast-1", ethereum_mainnet_blast_1_default],
    ["plume-mainnet", plume_mainnet_default],
    ["ethereum-mainnet-taiko-1", ethereum_mainnet_taiko_1_default],
    ["bitcoin-mainnet-bitlayer-1", bitcoin_mainnet_bitlayer_1_default],
    ["avalanche-subnet-dexalot-mainnet", avalanche_subnet_dexalot_mainnet_default],
    ["ethereum-mainnet-scroll-1", ethereum_mainnet_scroll_1_default],
    ["polygon-mainnet-katana", polygon_mainnet_katana_default],
    ["nexon-qa", nexon_qa_default],
    ["zklink_nova-mainnet", zklink_nova_mainnet_default],
    ["nexon-stage", nexon_stage_default],
    ["ethereum-mainnet-arbitrum-1-treasure-1", ethereum_mainnet_arbitrum_1_treasure_1_default],
    ["zora-mainnet", zora_mainnet_default],
    ["corn-mainnet", corn_mainnet_default],
    ["tron-mainnet-evm", tron_mainnet_evm_default]
  ]),
  solana: new Map([["solana-mainnet", solana_mainnet_default]]),
  aptos: new Map([["aptos-mainnet", aptos_mainnet_default]]),
  sui: new Map([["sui-mainnet", sui_mainnet_default]]),
  ton: new Map([["ton-mainnet", ton_mainnet_default]]),
  tron: new Map([["tron-mainnet", tron_mainnet_default]])
};
var testnetByNameByFamily = {
  evm: new Map([
    ["bitcoin-testnet-rootstock", bitcoin_testnet_rootstock_default],
    ["telos-evm-testnet", telos_evm_testnet_default],
    ["polkadot-testnet-darwinia-pangoro", polkadot_testnet_darwinia_pangoro_default],
    ["xdc-testnet", xdc_testnet_default],
    ["coinex_smart_chain-testnet", coinex_smart_chain_testnet_default],
    ["polkadot-testnet-astar-shibuya", polkadot_testnet_astar_shibuya_default],
    ["binance_smart_chain-testnet", binance_smart_chain_testnet_default],
    ["velas-testnet", velas_testnet_default],
    ["ethereum-testnet-sepolia-hashkey-1", ethereum_testnet_sepolia_hashkey_1_default],
    ["shibarium-testnet-puppynet", shibarium_testnet_puppynet_default],
    ["ethereum-testnet-sepolia-xlayer-1", ethereum_testnet_sepolia_xlayer_1_default],
    ["cronos-zkevm-testnet-sepolia", cronos_zkevm_testnet_sepolia_default],
    ["ethereum-testnet-goerli-zksync-1", ethereum_testnet_goerli_zksync_1_default],
    ["cronos-testnet-zkevm-1", cronos_testnet_zkevm_1_default],
    ["hedera-testnet", hedera_testnet_default],
    ["ethereum-testnet-sepolia-zksync-1", ethereum_testnet_sepolia_zksync_1_default],
    ["cronos-testnet", cronos_testnet_default],
    ["near-testnet", near_testnet_default],
    ["ethereum-testnet-goerli-optimism-1", ethereum_testnet_goerli_optimism_1_default],
    ["areon-testnet", areon_testnet_default],
    ["janction-testnet-sepolia", janction_testnet_sepolia_default],
    ["private-testnet-obsidian", private_testnet_obsidian_default],
    ["ethereum-testnet-sepolia-mode-1", ethereum_testnet_sepolia_mode_1_default],
    ["bittensor-testnet", bittensor_testnet_default],
    ["hyperliquid-testnet", hyperliquid_testnet_default],
    ["kaia-testnet-kairos", kaia_testnet_kairos_default],
    ["bittorrent_chain-testnet", bittorrent_chain_testnet_default],
    ["wemix-testnet", wemix_testnet_default],
    ["core-testnet", core_testnet_default],
    ["bitcoin-testnet-bsquared-1", bitcoin_testnet_bsquared_1_default],
    ["polkadot-testnet-moonbeam-moonbase", polkadot_testnet_moonbeam_moonbase_default],
    ["ethereum-testnet-sepolia-unichain-1", ethereum_testnet_sepolia_unichain_1_default],
    ["sei-testnet-atlantic", sei_testnet_atlantic_default],
    ["geth-testnet", geth_testnet_default],
    [
      "ethereum-testnet-goerli-polygon-zkevm-1",
      ethereum_testnet_goerli_polygon_zkevm_1_default
    ],
    ["story-testnet", story_testnet_default],
    ["mint-testnet", mint_testnet_default],
    ["metal-testnet", metal_testnet_default],
    ["bitcichain-testnet", bitcichain_testnet_default],
    ["ethereum-testnet-sepolia-soneium-1", ethereum_testnet_sepolia_soneium_1_default],
    ["ronin-testnet-saigon", ronin_testnet_saigon_default],
    ["private-testnet-granite", private_testnet_granite_default],
    ["private-testnet-andesite", private_testnet_andesite_default],
    ["dtcc-testnet-andesite", dtcc_testnet_andesite_default],
    ["polkadot-testnet-centrifuge-altair", polkadot_testnet_centrifuge_altair_default],
    ["memento-testnet", memento_testnet_default],
    ["kava-testnet", kava_testnet_default],
    ["ethereum-testnet-sepolia-kroma-1", ethereum_testnet_sepolia_kroma_1_default],
    ["tac-testnet", tac_testnet_default],
    [
      "ethereum-testnet-sepolia-polygon-zkevm-1",
      ethereum_testnet_sepolia_polygon_zkevm_1_default
    ],
    ["ethereum-testnet-holesky-fraxtal-1", ethereum_testnet_holesky_fraxtal_1_default],
    ["ethereum-testnet-holesky-morph-1", ethereum_testnet_holesky_morph_1_default],
    ["bitcoin-testnet-botanix", bitcoin_testnet_botanix_default],
    ["fantom-testnet", fantom_testnet_default],
    ["ethereum-testnet-sepolia-lisk-1", ethereum_testnet_sepolia_lisk_1_default],
    ["ethereum-testnet-sepolia-worldchain-1", ethereum_testnet_sepolia_worldchain_1_default],
    ["ethereum-testnet-goerli-mantle-1", ethereum_testnet_goerli_mantle_1_default],
    ["ethereum-testnet-sepolia-mantle-1", ethereum_testnet_sepolia_mantle_1_default],
    ["binance_smart_chain-testnet-opbnb-1", binance_smart_chain_testnet_opbnb_1_default],
    ["nexon-dev", nexon_dev_default],
    ["megaeth-testnet", megaeth_testnet_default],
    ["nibiru-testnet", nibiru_testnet_default],
    ["ondo-testnet", ondo_testnet_default],
    ["neonlink-testnet", neonlink_testnet_default],
    ["plasma-testnet", plasma_testnet_default],
    ["monad-testnet", monad_testnet_default],
    ["gnosis_chain-testnet-chiado", gnosis_chain_testnet_chiado_default],
    ["abstract-testnet", abstract_testnet_default],
    [
      "ethereum-testnet-sepolia-arbitrum-1-l3x-1",
      ethereum_testnet_sepolia_arbitrum_1_l3x_1_default
    ],
    [
      "ethereum-testnet-sepolia-immutable-zkevm-1",
      ethereum_testnet_sepolia_immutable_zkevm_1_default
    ],
    ["0g-testnet-newton", _0g_testnet_newton_default],
    ["0g-testnet-galileo", _0g_testnet_galileo_default],
    ["ethereum-testnet-holesky", ethereum_testnet_holesky_default],
    ["anvil-devnet", anvil_devnet_default],
    ["apechain-testnet-curtis", apechain_testnet_curtis_default],
    ["ethereum-testnet-sepolia-lens-1", ethereum_testnet_sepolia_lens_1_default],
    ["avalanche-testnet-fuji", avalanche_testnet_fuji_default],
    ["celo-testnet-alfajores", celo_testnet_alfajores_default],
    ["private-testnet-opala", private_testnet_opala_default],
    ["zircuit-testnet-garfield", zircuit_testnet_garfield_default],
    ["ethereum-testnet-sepolia-zircuit-1", ethereum_testnet_sepolia_zircuit_1_default],
    ["superseed-testnet", superseed_testnet_default],
    ["sonic-testnet-blaze", sonic_testnet_blaze_default],
    ["ethereum-testnet-goerli-linea-1", ethereum_testnet_goerli_linea_1_default],
    ["ethereum-testnet-sepolia-linea-1", ethereum_testnet_sepolia_linea_1_default],
    ["ethereum-testnet-sepolia-metis-1", ethereum_testnet_sepolia_metis_1_default],
    ["polygon-testnet-mumbai", polygon_testnet_mumbai_default],
    ["polygon-testnet-amoy", polygon_testnet_amoy_default],
    ["berachain-testnet-bepolia", berachain_testnet_bepolia_default],
    ["berachain-testnet-bartio", berachain_testnet_bartio_default],
    ["berachain-testnet-artio", berachain_testnet_artio_default],
    ["zero-g-testnet-galileo", zero_g_testnet_galileo_default],
    ["ethereum-testnet-goerli-base-1", ethereum_testnet_goerli_base_1_default],
    ["ethereum-testnet-sepolia-base-1", ethereum_testnet_sepolia_base_1_default],
    ["plume-devnet", plume_devnet_default],
    ["plume-testnet-sepolia", plume_testnet_sepolia_default],
    ["etherlink-testnet", etherlink_testnet_default],
    ["polygon-testnet-tatara", polygon_testnet_tatara_default],
    ["ethereum-testnet-holesky-taiko-1", ethereum_testnet_holesky_taiko_1_default],
    ["mind-testnet", mind_testnet_default],
    ["bitcoin-testnet-bitlayer-1", bitcoin_testnet_bitlayer_1_default],
    ["ethereum-testnet-goerli-arbitrum-1", ethereum_testnet_goerli_arbitrum_1_default],
    ["ethereum-testnet-sepolia-arbitrum-1", ethereum_testnet_sepolia_arbitrum_1_default],
    ["private-testnet-mica", private_testnet_mica_default],
    ["avalanche-subnet-dexalot-testnet", avalanche_subnet_dexalot_testnet_default],
    ["ethereum-testnet-sepolia-scroll-1", ethereum_testnet_sepolia_scroll_1_default],
    ["avalanche-testnet-nexon", avalanche_testnet_nexon_default],
    ["bitcoin-testnet-merlin", bitcoin_testnet_merlin_default],
    ["pharos-testnet", pharos_testnet_default],
    [
      "ethereum-testnet-sepolia-polygon-validium-1",
      ethereum_testnet_sepolia_polygon_validium_1_default
    ],
    ["hemi-testnet-sepolia", hemi_testnet_sepolia_default],
    ["ink-testnet-sepolia", ink_testnet_sepolia_default],
    ["bitcoin-testnet-sepolia-bob-1", bitcoin_testnet_sepolia_bob_1_default],
    ["zklink_nova-testnet", zklink_nova_testnet_default],
    ["codex-testnet", codex_testnet_default],
    [
      "ethereum-testnet-sepolia-arbitrum-1-treasure-1",
      ethereum_testnet_sepolia_arbitrum_1_treasure_1_default
    ],
    ["treasure-testnet-topaz", treasure_testnet_topaz_default],
    ["jovay-testnet", jovay_testnet_default],
    ["ethereum-testnet-sepolia", ethereum_testnet_sepolia_default],
    ["ethereum-testnet-sepolia-optimism-1", ethereum_testnet_sepolia_optimism_1_default],
    ["neox-testnet-t4", neox_testnet_t4_default],
    ["ethereum-testnet-sepolia-corn-1", ethereum_testnet_sepolia_corn_1_default],
    ["filecoin-testnet", filecoin_testnet_default],
    ["plume-testnet", plume_testnet_default],
    ["ethereum-testnet-sepolia-blast-1", ethereum_testnet_sepolia_blast_1_default],
    ["zora-testnet", zora_testnet_default],
    ["tron-testnet-shasta-evm", tron_testnet_shasta_evm_default],
    ["tron-devnet-evm", tron_devnet_evm_default],
    ["tron-testnet-nile-evm", tron_testnet_nile_evm_default]
  ]),
  solana: new Map([
    ["solana-testnet", solana_testnet_default],
    ["solana-devnet", solana_devnet_default]
  ]),
  aptos: new Map([
    ["aptos-testnet", aptos_testnet_default],
    ["aptos-localnet", aptos_localnet_default]
  ]),
  sui: new Map([
    ["sui-testnet", sui_testnet_default],
    ["sui-localnet", sui_localnet_default]
  ]),
  ton: new Map([
    ["ton-testnet", ton_testnet_default],
    ["ton-localnet", ton_localnet_default]
  ]),
  tron: new Map([
    ["tron-testnet-shasta", tron_testnet_shasta_default],
    ["tron-devnet", tron_devnet_default],
    ["tron-testnet-nile", tron_testnet_nile_default]
  ])
};

class Int64 {
  static INT64_MIN = -(2n ** 63n);
  static INT64_MAX = 2n ** 63n - 1n;
  value;
  static toInt64Bigint(v) {
    if (typeof v === "string") {
      const bi2 = BigInt(v);
      return Int64.toInt64Bigint(bi2);
    }
    if (typeof v === "bigint") {
      if (v > Int64.INT64_MAX)
        throw new Error("int64 overflow");
      else if (v < Int64.INT64_MIN)
        throw new Error("int64 underflow");
      return v;
    }
    if (!Number.isFinite(v) || !Number.isInteger(v))
      throw new Error("int64 requires an integer number");
    const bi = BigInt(v);
    if (bi > Int64.INT64_MAX)
      throw new Error("int64 overflow");
    else if (bi < Int64.INT64_MIN)
      throw new Error("int64 underflow");
    return bi;
  }
  constructor(v) {
    this.value = Int64.toInt64Bigint(v);
  }
  add(i2, safe = true) {
    return safe ? new Int64(this.value + i2.value) : new Int64(BigInt.asIntN(64, this.value + i2.value));
  }
  sub(i2, safe = true) {
    return safe ? new Int64(this.value - i2.value) : new Int64(BigInt.asIntN(64, this.value - i2.value));
  }
  mul(i2, safe = true) {
    return safe ? new Int64(this.value * i2.value) : new Int64(BigInt.asIntN(64, this.value * i2.value));
  }
  div(i2, safe = true) {
    return new Int64(this.value / i2.value);
  }
}

class UInt64 {
  static UINT64_MAX = 2n ** 64n - 1n;
  value;
  static toUint64Bigint(v) {
    if (typeof v === "string") {
      const bi2 = BigInt(v);
      return UInt64.toUint64Bigint(bi2);
    }
    if (typeof v === "bigint") {
      if (v > UInt64.UINT64_MAX)
        throw new Error("uint64 overflow");
      else if (v < 0n)
        throw new Error("uint64 underflow");
      return v;
    }
    if (!Number.isFinite(v) || !Number.isInteger(v))
      throw new Error("int64 requires an integer number");
    const bi = BigInt(v);
    if (bi > UInt64.UINT64_MAX)
      throw new Error("uint64 overflow");
    else if (bi < 0n)
      throw new Error("uint64 underflow");
    return bi;
  }
  constructor(v) {
    this.value = UInt64.toUint64Bigint(v);
  }
  add(i2, safe = true) {
    return safe ? new UInt64(this.value + i2.value) : new UInt64(BigInt.asUintN(64, this.value + i2.value));
  }
  sub(i2, safe = true) {
    return safe ? new UInt64(this.value - i2.value) : new UInt64(BigInt.asUintN(64, this.value - i2.value));
  }
  mul(i2, safe = true) {
    return safe ? new UInt64(this.value * i2.value) : new UInt64(BigInt.asUintN(64, this.value * i2.value));
  }
  div(i2, safe = true) {
    return new UInt64(this.value / i2.value);
  }
}

class Decimal {
  coeffecient;
  exponent;
  static parse(s) {
    const m = /^([+-])?(\d+)(?:\.(\d+))?$/.exec(s.trim());
    if (!m)
      throw new Error("invalid decimal string");
    const signStr = m[1] ?? "+";
    const intPart = m[2] ?? "0";
    let fracPart = m[3] ?? "";
    while (fracPart.length > 0 && fracPart[fracPart.length - 1] === "0") {
      fracPart = fracPart.slice(0, -1);
    }
    const exponent = fracPart.length === 0 ? 0 : -fracPart.length;
    const digits = intPart + fracPart || "0";
    const coeffecient = BigInt((signStr === "-" ? "-" : "") + digits);
    return new Decimal(coeffecient, exponent);
  }
  constructor(coeffecient, exponent) {
    this.coeffecient = coeffecient;
    this.exponent = exponent;
  }
}

class Value {
  value;
  static from(value) {
    return new Value(value);
  }
  static wrap(value) {
    return new Value(value);
  }
  constructor(value) {
    if (value instanceof Value) {
      this.value = value.value;
    } else if (isValueProto(value)) {
      this.value = value;
    } else {
      this.value = Value.wrapInternal(value);
    }
  }
  proto() {
    return this.value;
  }
  static toUint8Array(input) {
    return input instanceof Uint8Array ? input : new Uint8Array(input);
  }
  static bigintToBytesBE(abs) {
    if (abs === 0n)
      return new Uint8Array;
    let hex = abs.toString(16);
    if (hex.length % 2 === 1)
      hex = "0" + hex;
    const len2 = hex.length / 2;
    const out = new Uint8Array(len2);
    for (let i2 = 0;i2 < len2; i2++) {
      out[i2] = parseInt(hex.slice(i2 * 2, i2 * 2 + 2), 16);
    }
    return out;
  }
  static bigIntToProtoBigInt(v) {
    const sign = v === 0n ? 0n : v < 0n ? -1n : 1n;
    const abs = v < 0n ? -v : v;
    return create(BigIntSchema, {
      absVal: Value.bigintToBytesBE(abs),
      sign
    });
  }
  static toTimestamp(d) {
    const date = d instanceof Date ? d : new Date(d);
    return timestampFromDate(date);
  }
  static isPlainObject(v) {
    return typeof v === "object" && v !== null && v.constructor === Object;
  }
  static isObject(v) {
    return typeof v === "object" && v !== null;
  }
  static wrapInternal(v) {
    if (v === null || v === undefined)
      throw new Error("cannot wrap null/undefined into Value");
    if (v instanceof Value) {
      return v.proto();
    }
    if (v instanceof Uint8Array)
      return create(ValueSchema2, { value: { case: "bytesValue", value: v } });
    if (v instanceof ArrayBuffer)
      return create(ValueSchema2, {
        value: { case: "bytesValue", value: Value.toUint8Array(v) }
      });
    if (v instanceof Date)
      return create(ValueSchema2, {
        value: { case: "timeValue", value: Value.toTimestamp(v) }
      });
    if (v instanceof Int64) {
      return create(ValueSchema2, {
        value: { case: "int64Value", value: v.value }
      });
    }
    if (v instanceof UInt64) {
      return create(ValueSchema2, {
        value: { case: "uint64Value", value: v.value }
      });
    }
    if (v instanceof Decimal) {
      const decimalProto = create(DecimalSchema, {
        coefficient: Value.bigIntToProtoBigInt(v.coeffecient),
        exponent: v.exponent
      });
      return create(ValueSchema2, {
        value: { case: "decimalValue", value: decimalProto }
      });
    }
    switch (typeof v) {
      case "string":
        return create(ValueSchema2, {
          value: { case: "stringValue", value: v }
        });
      case "boolean":
        return create(ValueSchema2, { value: { case: "boolValue", value: v } });
      case "bigint": {
        return create(ValueSchema2, {
          value: { case: "bigintValue", value: Value.bigIntToProtoBigInt(v) }
        });
      }
      case "number": {
        return create(ValueSchema2, {
          value: { case: "float64Value", value: v }
        });
      }
      case "object":
        break;
      default:
        throw new Error(`unsupported type: ${typeof v}`);
    }
    if (Array.isArray(v)) {
      const fields2 = v.map(Value.wrapInternal);
      const list = create(ListSchema, { fields: fields2 });
      return create(ValueSchema2, { value: { case: "listValue", value: list } });
    }
    if (Value.isPlainObject(v)) {
      const fields2 = {};
      for (const [k, vv] of Object.entries(v)) {
        fields2[k] = Value.wrapInternal(vv);
      }
      const map = create(MapSchema, { fields: fields2 });
      return create(ValueSchema2, { value: { case: "mapValue", value: map } });
    }
    if (Value.isObject(v) && v.constructor !== Object) {
      const fields2 = {};
      for (const [k, vv] of Object.entries(v)) {
        fields2[k] = Value.wrapInternal(vv);
      }
      const map = create(MapSchema, { fields: fields2 });
      return create(ValueSchema2, { value: { case: "mapValue", value: map } });
    }
    throw new Error("unsupported object instance");
  }
  unwrap() {
    return unwrap(this.value);
  }
  unwrapToType(options) {
    const unwrapped = this.unwrap();
    if ("instance" in options) {
      if (typeof unwrapped !== typeof options.instance) {
        throw new Error(`Cannot unwrap to type ${typeof options.instance}`);
      }
      return unwrapped;
    }
    if (options.schema) {
      return options.schema.parse(unwrapped);
    }
    const obj = options.factory();
    if (typeof unwrapped === "object" && unwrapped !== null) {
      Object.assign(obj, unwrapped);
    } else {
      throw new Error(`Cannot copy properties from primitive value to object instance. Use a schema instead.`);
    }
    return obj;
  }
}
function unwrap(value) {
  switch (value.value.case) {
    case "stringValue":
      return value.value.value;
    case "boolValue":
      return value.value.value;
    case "bytesValue":
      return value.value.value;
    case "int64Value":
      return new Int64(value.value.value);
    case "uint64Value":
      return new UInt64(value.value.value);
    case "float64Value":
      return value.value.value;
    case "bigintValue": {
      const bigIntValue = value.value.value;
      const absVal = bigIntValue.absVal;
      const sign = bigIntValue.sign;
      let result = 0n;
      for (const byte of absVal) {
        result = result << 8n | BigInt(byte);
      }
      return sign < 0n ? -result : result;
    }
    case "timeValue": {
      return timestampDate(value.value.value);
    }
    case "listValue": {
      const list = value.value.value;
      return list.fields.map(unwrap);
    }
    case "mapValue": {
      const map = value.value.value;
      const result = {};
      for (const [key, val] of Object.entries(map.fields)) {
        result[key] = unwrap(val);
      }
      return result;
    }
    case "decimalValue": {
      const decimal = value.value.value;
      const coefficient = decimal.coefficient;
      const exponent = decimal.exponent;
      if (!coefficient) {
        return new Decimal(0n, 0);
      }
      let coeffBigInt;
      const absVal = coefficient.absVal;
      const sign = coefficient.sign;
      let result = 0n;
      for (const byte of absVal) {
        result = result << 8n | BigInt(byte);
      }
      coeffBigInt = sign < 0n ? -result : result;
      return new Decimal(coeffBigInt, exponent);
    }
    default:
      throw new Error(`Unsupported value type: ${value.value.case}`);
  }
}
function isValueProto(value) {
  return value.$typeName && typeof value.$typeName === "string" && value.$typeName === "values.v1.Value";
}
async function standardValidate(schema, input) {
  let result = schema["~standard"].validate(input);
  if (result instanceof Promise)
    result = await result;
  if (result.issues) {
    throw new Error(JSON.stringify(result.issues, null, 2));
  }
  return result.value;
}
var defaultJsonParser = (config) => JSON.parse(Buffer.from(config).toString());
var configHandler = async (request, { configParser, configSchema } = {}) => {
  const config = request.config;
  const parser = configParser || defaultJsonParser;
  let intermediateConfig;
  try {
    intermediateConfig = parser(config);
  } catch (error) {
    if (error instanceof Error) {
      throw new Error(`Failed to parse configuration: ${error.message}`);
    } else {
      throw new Error(`Failed to parse configuration: unknown error`);
    }
  }
  return configSchema ? standardValidate(configSchema, intermediateConfig) : intermediateConfig;
};
var exports_external = {};
__export(exports_external, {
  void: () => voidType,
  util: () => util,
  unknown: () => unknownType,
  union: () => unionType,
  undefined: () => undefinedType,
  tuple: () => tupleType,
  transformer: () => effectsType,
  symbol: () => symbolType,
  string: () => stringType,
  strictObject: () => strictObjectType,
  setErrorMap: () => setErrorMap,
  set: () => setType,
  record: () => recordType,
  quotelessJson: () => quotelessJson,
  promise: () => promiseType,
  preprocess: () => preprocessType,
  pipeline: () => pipelineType,
  ostring: () => ostring,
  optional: () => optionalType,
  onumber: () => onumber,
  oboolean: () => oboolean,
  objectUtil: () => objectUtil,
  object: () => objectType,
  number: () => numberType,
  nullable: () => nullableType,
  null: () => nullType,
  never: () => neverType,
  nativeEnum: () => nativeEnumType,
  nan: () => nanType,
  map: () => mapType,
  makeIssue: () => makeIssue,
  literal: () => literalType,
  lazy: () => lazyType,
  late: () => late,
  isValid: () => isValid,
  isDirty: () => isDirty,
  isAsync: () => isAsync,
  isAborted: () => isAborted,
  intersection: () => intersectionType,
  instanceof: () => instanceOfType,
  getParsedType: () => getParsedType,
  getErrorMap: () => getErrorMap,
  function: () => functionType,
  enum: () => enumType,
  effect: () => effectsType,
  discriminatedUnion: () => discriminatedUnionType,
  defaultErrorMap: () => en_default,
  datetimeRegex: () => datetimeRegex,
  date: () => dateType,
  custom: () => custom,
  coerce: () => coerce,
  boolean: () => booleanType,
  bigint: () => bigIntType,
  array: () => arrayType,
  any: () => anyType,
  addIssueToContext: () => addIssueToContext,
  ZodVoid: () => ZodVoid,
  ZodUnknown: () => ZodUnknown,
  ZodUnion: () => ZodUnion,
  ZodUndefined: () => ZodUndefined,
  ZodType: () => ZodType,
  ZodTuple: () => ZodTuple,
  ZodTransformer: () => ZodEffects,
  ZodSymbol: () => ZodSymbol,
  ZodString: () => ZodString,
  ZodSet: () => ZodSet,
  ZodSchema: () => ZodType,
  ZodRecord: () => ZodRecord,
  ZodReadonly: () => ZodReadonly,
  ZodPromise: () => ZodPromise,
  ZodPipeline: () => ZodPipeline,
  ZodParsedType: () => ZodParsedType,
  ZodOptional: () => ZodOptional,
  ZodObject: () => ZodObject,
  ZodNumber: () => ZodNumber,
  ZodNullable: () => ZodNullable,
  ZodNull: () => ZodNull,
  ZodNever: () => ZodNever,
  ZodNativeEnum: () => ZodNativeEnum,
  ZodNaN: () => ZodNaN,
  ZodMap: () => ZodMap,
  ZodLiteral: () => ZodLiteral,
  ZodLazy: () => ZodLazy,
  ZodIssueCode: () => ZodIssueCode,
  ZodIntersection: () => ZodIntersection,
  ZodFunction: () => ZodFunction,
  ZodFirstPartyTypeKind: () => ZodFirstPartyTypeKind,
  ZodError: () => ZodError,
  ZodEnum: () => ZodEnum,
  ZodEffects: () => ZodEffects,
  ZodDiscriminatedUnion: () => ZodDiscriminatedUnion,
  ZodDefault: () => ZodDefault,
  ZodDate: () => ZodDate,
  ZodCatch: () => ZodCatch,
  ZodBranded: () => ZodBranded,
  ZodBoolean: () => ZodBoolean,
  ZodBigInt: () => ZodBigInt,
  ZodArray: () => ZodArray,
  ZodAny: () => ZodAny,
  Schema: () => ZodType,
  ParseStatus: () => ParseStatus,
  OK: () => OK,
  NEVER: () => NEVER,
  INVALID: () => INVALID,
  EMPTY_PATH: () => EMPTY_PATH,
  DIRTY: () => DIRTY,
  BRAND: () => BRAND
});
var util;
(function(util2) {
  util2.assertEqual = (_) => {};
  function assertIs(_arg) {}
  util2.assertIs = assertIs;
  function assertNever(_x) {
    throw new Error;
  }
  util2.assertNever = assertNever;
  util2.arrayToEnum = (items) => {
    const obj = {};
    for (const item of items) {
      obj[item] = item;
    }
    return obj;
  };
  util2.getValidEnumValues = (obj) => {
    const validKeys = util2.objectKeys(obj).filter((k) => typeof obj[obj[k]] !== "number");
    const filtered = {};
    for (const k of validKeys) {
      filtered[k] = obj[k];
    }
    return util2.objectValues(filtered);
  };
  util2.objectValues = (obj) => {
    return util2.objectKeys(obj).map(function(e) {
      return obj[e];
    });
  };
  util2.objectKeys = typeof Object.keys === "function" ? (obj) => Object.keys(obj) : (object) => {
    const keys = [];
    for (const key in object) {
      if (Object.prototype.hasOwnProperty.call(object, key)) {
        keys.push(key);
      }
    }
    return keys;
  };
  util2.find = (arr, checker) => {
    for (const item of arr) {
      if (checker(item))
        return item;
    }
    return;
  };
  util2.isInteger = typeof Number.isInteger === "function" ? (val) => Number.isInteger(val) : (val) => typeof val === "number" && Number.isFinite(val) && Math.floor(val) === val;
  function joinValues(array, separator = " | ") {
    return array.map((val) => typeof val === "string" ? `'${val}'` : val).join(separator);
  }
  util2.joinValues = joinValues;
  util2.jsonStringifyReplacer = (_, value2) => {
    if (typeof value2 === "bigint") {
      return value2.toString();
    }
    return value2;
  };
})(util || (util = {}));
var objectUtil;
(function(objectUtil2) {
  objectUtil2.mergeShapes = (first, second) => {
    return {
      ...first,
      ...second
    };
  };
})(objectUtil || (objectUtil = {}));
var ZodParsedType = util.arrayToEnum([
  "string",
  "nan",
  "number",
  "integer",
  "float",
  "boolean",
  "date",
  "bigint",
  "symbol",
  "function",
  "undefined",
  "null",
  "array",
  "object",
  "unknown",
  "promise",
  "void",
  "never",
  "map",
  "set"
]);
var getParsedType = (data) => {
  const t = typeof data;
  switch (t) {
    case "undefined":
      return ZodParsedType.undefined;
    case "string":
      return ZodParsedType.string;
    case "number":
      return Number.isNaN(data) ? ZodParsedType.nan : ZodParsedType.number;
    case "boolean":
      return ZodParsedType.boolean;
    case "function":
      return ZodParsedType.function;
    case "bigint":
      return ZodParsedType.bigint;
    case "symbol":
      return ZodParsedType.symbol;
    case "object":
      if (Array.isArray(data)) {
        return ZodParsedType.array;
      }
      if (data === null) {
        return ZodParsedType.null;
      }
      if (data.then && typeof data.then === "function" && data.catch && typeof data.catch === "function") {
        return ZodParsedType.promise;
      }
      if (typeof Map !== "undefined" && data instanceof Map) {
        return ZodParsedType.map;
      }
      if (typeof Set !== "undefined" && data instanceof Set) {
        return ZodParsedType.set;
      }
      if (typeof Date !== "undefined" && data instanceof Date) {
        return ZodParsedType.date;
      }
      return ZodParsedType.object;
    default:
      return ZodParsedType.unknown;
  }
};
var ZodIssueCode = util.arrayToEnum([
  "invalid_type",
  "invalid_literal",
  "custom",
  "invalid_union",
  "invalid_union_discriminator",
  "invalid_enum_value",
  "unrecognized_keys",
  "invalid_arguments",
  "invalid_return_type",
  "invalid_date",
  "invalid_string",
  "too_small",
  "too_big",
  "invalid_intersection_types",
  "not_multiple_of",
  "not_finite"
]);
var quotelessJson = (obj) => {
  const json = JSON.stringify(obj, null, 2);
  return json.replace(/"([^"]+)":/g, "$1:");
};

class ZodError extends Error {
  get errors() {
    return this.issues;
  }
  constructor(issues) {
    super();
    this.issues = [];
    this.addIssue = (sub) => {
      this.issues = [...this.issues, sub];
    };
    this.addIssues = (subs = []) => {
      this.issues = [...this.issues, ...subs];
    };
    const actualProto = new.target.prototype;
    if (Object.setPrototypeOf) {
      Object.setPrototypeOf(this, actualProto);
    } else {
      this.__proto__ = actualProto;
    }
    this.name = "ZodError";
    this.issues = issues;
  }
  format(_mapper) {
    const mapper = _mapper || function(issue) {
      return issue.message;
    };
    const fieldErrors = { _errors: [] };
    const processError = (error) => {
      for (const issue of error.issues) {
        if (issue.code === "invalid_union") {
          issue.unionErrors.map(processError);
        } else if (issue.code === "invalid_return_type") {
          processError(issue.returnTypeError);
        } else if (issue.code === "invalid_arguments") {
          processError(issue.argumentsError);
        } else if (issue.path.length === 0) {
          fieldErrors._errors.push(mapper(issue));
        } else {
          let curr = fieldErrors;
          let i2 = 0;
          while (i2 < issue.path.length) {
            const el = issue.path[i2];
            const terminal = i2 === issue.path.length - 1;
            if (!terminal) {
              curr[el] = curr[el] || { _errors: [] };
            } else {
              curr[el] = curr[el] || { _errors: [] };
              curr[el]._errors.push(mapper(issue));
            }
            curr = curr[el];
            i2++;
          }
        }
      }
    };
    processError(this);
    return fieldErrors;
  }
  static assert(value2) {
    if (!(value2 instanceof ZodError)) {
      throw new Error(`Not a ZodError: ${value2}`);
    }
  }
  toString() {
    return this.message;
  }
  get message() {
    return JSON.stringify(this.issues, util.jsonStringifyReplacer, 2);
  }
  get isEmpty() {
    return this.issues.length === 0;
  }
  flatten(mapper = (issue) => issue.message) {
    const fieldErrors = {};
    const formErrors = [];
    for (const sub of this.issues) {
      if (sub.path.length > 0) {
        const firstEl = sub.path[0];
        fieldErrors[firstEl] = fieldErrors[firstEl] || [];
        fieldErrors[firstEl].push(mapper(sub));
      } else {
        formErrors.push(mapper(sub));
      }
    }
    return { formErrors, fieldErrors };
  }
  get formErrors() {
    return this.flatten();
  }
}
ZodError.create = (issues) => {
  const error = new ZodError(issues);
  return error;
};
var errorMap = (issue, _ctx) => {
  let message;
  switch (issue.code) {
    case ZodIssueCode.invalid_type:
      if (issue.received === ZodParsedType.undefined) {
        message = "Required";
      } else {
        message = `Expected ${issue.expected}, received ${issue.received}`;
      }
      break;
    case ZodIssueCode.invalid_literal:
      message = `Invalid literal value, expected ${JSON.stringify(issue.expected, util.jsonStringifyReplacer)}`;
      break;
    case ZodIssueCode.unrecognized_keys:
      message = `Unrecognized key(s) in object: ${util.joinValues(issue.keys, ", ")}`;
      break;
    case ZodIssueCode.invalid_union:
      message = `Invalid input`;
      break;
    case ZodIssueCode.invalid_union_discriminator:
      message = `Invalid discriminator value. Expected ${util.joinValues(issue.options)}`;
      break;
    case ZodIssueCode.invalid_enum_value:
      message = `Invalid enum value. Expected ${util.joinValues(issue.options)}, received '${issue.received}'`;
      break;
    case ZodIssueCode.invalid_arguments:
      message = `Invalid function arguments`;
      break;
    case ZodIssueCode.invalid_return_type:
      message = `Invalid function return type`;
      break;
    case ZodIssueCode.invalid_date:
      message = `Invalid date`;
      break;
    case ZodIssueCode.invalid_string:
      if (typeof issue.validation === "object") {
        if ("includes" in issue.validation) {
          message = `Invalid input: must include "${issue.validation.includes}"`;
          if (typeof issue.validation.position === "number") {
            message = `${message} at one or more positions greater than or equal to ${issue.validation.position}`;
          }
        } else if ("startsWith" in issue.validation) {
          message = `Invalid input: must start with "${issue.validation.startsWith}"`;
        } else if ("endsWith" in issue.validation) {
          message = `Invalid input: must end with "${issue.validation.endsWith}"`;
        } else {
          util.assertNever(issue.validation);
        }
      } else if (issue.validation !== "regex") {
        message = `Invalid ${issue.validation}`;
      } else {
        message = "Invalid";
      }
      break;
    case ZodIssueCode.too_small:
      if (issue.type === "array")
        message = `Array must contain ${issue.exact ? "exactly" : issue.inclusive ? `at least` : `more than`} ${issue.minimum} element(s)`;
      else if (issue.type === "string")
        message = `String must contain ${issue.exact ? "exactly" : issue.inclusive ? `at least` : `over`} ${issue.minimum} character(s)`;
      else if (issue.type === "number")
        message = `Number must be ${issue.exact ? `exactly equal to ` : issue.inclusive ? `greater than or equal to ` : `greater than `}${issue.minimum}`;
      else if (issue.type === "bigint")
        message = `Number must be ${issue.exact ? `exactly equal to ` : issue.inclusive ? `greater than or equal to ` : `greater than `}${issue.minimum}`;
      else if (issue.type === "date")
        message = `Date must be ${issue.exact ? `exactly equal to ` : issue.inclusive ? `greater than or equal to ` : `greater than `}${new Date(Number(issue.minimum))}`;
      else
        message = "Invalid input";
      break;
    case ZodIssueCode.too_big:
      if (issue.type === "array")
        message = `Array must contain ${issue.exact ? `exactly` : issue.inclusive ? `at most` : `less than`} ${issue.maximum} element(s)`;
      else if (issue.type === "string")
        message = `String must contain ${issue.exact ? `exactly` : issue.inclusive ? `at most` : `under`} ${issue.maximum} character(s)`;
      else if (issue.type === "number")
        message = `Number must be ${issue.exact ? `exactly` : issue.inclusive ? `less than or equal to` : `less than`} ${issue.maximum}`;
      else if (issue.type === "bigint")
        message = `BigInt must be ${issue.exact ? `exactly` : issue.inclusive ? `less than or equal to` : `less than`} ${issue.maximum}`;
      else if (issue.type === "date")
        message = `Date must be ${issue.exact ? `exactly` : issue.inclusive ? `smaller than or equal to` : `smaller than`} ${new Date(Number(issue.maximum))}`;
      else
        message = "Invalid input";
      break;
    case ZodIssueCode.custom:
      message = `Invalid input`;
      break;
    case ZodIssueCode.invalid_intersection_types:
      message = `Intersection results could not be merged`;
      break;
    case ZodIssueCode.not_multiple_of:
      message = `Number must be a multiple of ${issue.multipleOf}`;
      break;
    case ZodIssueCode.not_finite:
      message = "Number must be finite";
      break;
    default:
      message = _ctx.defaultError;
      util.assertNever(issue);
  }
  return { message };
};
var en_default = errorMap;
var overrideErrorMap = en_default;
function setErrorMap(map) {
  overrideErrorMap = map;
}
function getErrorMap() {
  return overrideErrorMap;
}
var makeIssue = (params) => {
  const { data, path, errorMaps, issueData } = params;
  const fullPath = [...path, ...issueData.path || []];
  const fullIssue = {
    ...issueData,
    path: fullPath
  };
  if (issueData.message !== undefined) {
    return {
      ...issueData,
      path: fullPath,
      message: issueData.message
    };
  }
  let errorMessage = "";
  const maps = errorMaps.filter((m) => !!m).slice().reverse();
  for (const map of maps) {
    errorMessage = map(fullIssue, { data, defaultError: errorMessage }).message;
  }
  return {
    ...issueData,
    path: fullPath,
    message: errorMessage
  };
};
var EMPTY_PATH = [];
function addIssueToContext(ctx, issueData) {
  const overrideMap = getErrorMap();
  const issue = makeIssue({
    issueData,
    data: ctx.data,
    path: ctx.path,
    errorMaps: [
      ctx.common.contextualErrorMap,
      ctx.schemaErrorMap,
      overrideMap,
      overrideMap === en_default ? undefined : en_default
    ].filter((x) => !!x)
  });
  ctx.common.issues.push(issue);
}

class ParseStatus {
  constructor() {
    this.value = "valid";
  }
  dirty() {
    if (this.value === "valid")
      this.value = "dirty";
  }
  abort() {
    if (this.value !== "aborted")
      this.value = "aborted";
  }
  static mergeArray(status, results) {
    const arrayValue = [];
    for (const s of results) {
      if (s.status === "aborted")
        return INVALID;
      if (s.status === "dirty")
        status.dirty();
      arrayValue.push(s.value);
    }
    return { status: status.value, value: arrayValue };
  }
  static async mergeObjectAsync(status, pairs) {
    const syncPairs = [];
    for (const pair of pairs) {
      const key = await pair.key;
      const value2 = await pair.value;
      syncPairs.push({
        key,
        value: value2
      });
    }
    return ParseStatus.mergeObjectSync(status, syncPairs);
  }
  static mergeObjectSync(status, pairs) {
    const finalObject = {};
    for (const pair of pairs) {
      const { key, value: value2 } = pair;
      if (key.status === "aborted")
        return INVALID;
      if (value2.status === "aborted")
        return INVALID;
      if (key.status === "dirty")
        status.dirty();
      if (value2.status === "dirty")
        status.dirty();
      if (key.value !== "__proto__" && (typeof value2.value !== "undefined" || pair.alwaysSet)) {
        finalObject[key.value] = value2.value;
      }
    }
    return { status: status.value, value: finalObject };
  }
}
var INVALID = Object.freeze({
  status: "aborted"
});
var DIRTY = (value2) => ({ status: "dirty", value: value2 });
var OK = (value2) => ({ status: "valid", value: value2 });
var isAborted = (x) => x.status === "aborted";
var isDirty = (x) => x.status === "dirty";
var isValid = (x) => x.status === "valid";
var isAsync = (x) => typeof Promise !== "undefined" && x instanceof Promise;
var errorUtil;
(function(errorUtil2) {
  errorUtil2.errToObj = (message) => typeof message === "string" ? { message } : message || {};
  errorUtil2.toString = (message) => typeof message === "string" ? message : message?.message;
})(errorUtil || (errorUtil = {}));

class ParseInputLazyPath {
  constructor(parent, value2, path, key) {
    this._cachedPath = [];
    this.parent = parent;
    this.data = value2;
    this._path = path;
    this._key = key;
  }
  get path() {
    if (!this._cachedPath.length) {
      if (Array.isArray(this._key)) {
        this._cachedPath.push(...this._path, ...this._key);
      } else {
        this._cachedPath.push(...this._path, this._key);
      }
    }
    return this._cachedPath;
  }
}
var handleResult = (ctx, result) => {
  if (isValid(result)) {
    return { success: true, data: result.value };
  } else {
    if (!ctx.common.issues.length) {
      throw new Error("Validation failed but no issues detected.");
    }
    return {
      success: false,
      get error() {
        if (this._error)
          return this._error;
        const error = new ZodError(ctx.common.issues);
        this._error = error;
        return this._error;
      }
    };
  }
};
function processCreateParams(params) {
  if (!params)
    return {};
  const { errorMap: errorMap2, invalid_type_error, required_error, description } = params;
  if (errorMap2 && (invalid_type_error || required_error)) {
    throw new Error(`Can't use "invalid_type_error" or "required_error" in conjunction with custom error map.`);
  }
  if (errorMap2)
    return { errorMap: errorMap2, description };
  const customMap = (iss, ctx) => {
    const { message } = params;
    if (iss.code === "invalid_enum_value") {
      return { message: message ?? ctx.defaultError };
    }
    if (typeof ctx.data === "undefined") {
      return { message: message ?? required_error ?? ctx.defaultError };
    }
    if (iss.code !== "invalid_type")
      return { message: ctx.defaultError };
    return { message: message ?? invalid_type_error ?? ctx.defaultError };
  };
  return { errorMap: customMap, description };
}

class ZodType {
  get description() {
    return this._def.description;
  }
  _getType(input) {
    return getParsedType(input.data);
  }
  _getOrReturnCtx(input, ctx) {
    return ctx || {
      common: input.parent.common,
      data: input.data,
      parsedType: getParsedType(input.data),
      schemaErrorMap: this._def.errorMap,
      path: input.path,
      parent: input.parent
    };
  }
  _processInputParams(input) {
    return {
      status: new ParseStatus,
      ctx: {
        common: input.parent.common,
        data: input.data,
        parsedType: getParsedType(input.data),
        schemaErrorMap: this._def.errorMap,
        path: input.path,
        parent: input.parent
      }
    };
  }
  _parseSync(input) {
    const result = this._parse(input);
    if (isAsync(result)) {
      throw new Error("Synchronous parse encountered promise.");
    }
    return result;
  }
  _parseAsync(input) {
    const result = this._parse(input);
    return Promise.resolve(result);
  }
  parse(data, params) {
    const result = this.safeParse(data, params);
    if (result.success)
      return result.data;
    throw result.error;
  }
  safeParse(data, params) {
    const ctx = {
      common: {
        issues: [],
        async: params?.async ?? false,
        contextualErrorMap: params?.errorMap
      },
      path: params?.path || [],
      schemaErrorMap: this._def.errorMap,
      parent: null,
      data,
      parsedType: getParsedType(data)
    };
    const result = this._parseSync({ data, path: ctx.path, parent: ctx });
    return handleResult(ctx, result);
  }
  "~validate"(data) {
    const ctx = {
      common: {
        issues: [],
        async: !!this["~standard"].async
      },
      path: [],
      schemaErrorMap: this._def.errorMap,
      parent: null,
      data,
      parsedType: getParsedType(data)
    };
    if (!this["~standard"].async) {
      try {
        const result = this._parseSync({ data, path: [], parent: ctx });
        return isValid(result) ? {
          value: result.value
        } : {
          issues: ctx.common.issues
        };
      } catch (err) {
        if (err?.message?.toLowerCase()?.includes("encountered")) {
          this["~standard"].async = true;
        }
        ctx.common = {
          issues: [],
          async: true
        };
      }
    }
    return this._parseAsync({ data, path: [], parent: ctx }).then((result) => isValid(result) ? {
      value: result.value
    } : {
      issues: ctx.common.issues
    });
  }
  async parseAsync(data, params) {
    const result = await this.safeParseAsync(data, params);
    if (result.success)
      return result.data;
    throw result.error;
  }
  async safeParseAsync(data, params) {
    const ctx = {
      common: {
        issues: [],
        contextualErrorMap: params?.errorMap,
        async: true
      },
      path: params?.path || [],
      schemaErrorMap: this._def.errorMap,
      parent: null,
      data,
      parsedType: getParsedType(data)
    };
    const maybeAsyncResult = this._parse({ data, path: ctx.path, parent: ctx });
    const result = await (isAsync(maybeAsyncResult) ? maybeAsyncResult : Promise.resolve(maybeAsyncResult));
    return handleResult(ctx, result);
  }
  refine(check, message) {
    const getIssueProperties = (val) => {
      if (typeof message === "string" || typeof message === "undefined") {
        return { message };
      } else if (typeof message === "function") {
        return message(val);
      } else {
        return message;
      }
    };
    return this._refinement((val, ctx) => {
      const result = check(val);
      const setError = () => ctx.addIssue({
        code: ZodIssueCode.custom,
        ...getIssueProperties(val)
      });
      if (typeof Promise !== "undefined" && result instanceof Promise) {
        return result.then((data) => {
          if (!data) {
            setError();
            return false;
          } else {
            return true;
          }
        });
      }
      if (!result) {
        setError();
        return false;
      } else {
        return true;
      }
    });
  }
  refinement(check, refinementData) {
    return this._refinement((val, ctx) => {
      if (!check(val)) {
        ctx.addIssue(typeof refinementData === "function" ? refinementData(val, ctx) : refinementData);
        return false;
      } else {
        return true;
      }
    });
  }
  _refinement(refinement) {
    return new ZodEffects({
      schema: this,
      typeName: ZodFirstPartyTypeKind.ZodEffects,
      effect: { type: "refinement", refinement }
    });
  }
  superRefine(refinement) {
    return this._refinement(refinement);
  }
  constructor(def) {
    this.spa = this.safeParseAsync;
    this._def = def;
    this.parse = this.parse.bind(this);
    this.safeParse = this.safeParse.bind(this);
    this.parseAsync = this.parseAsync.bind(this);
    this.safeParseAsync = this.safeParseAsync.bind(this);
    this.spa = this.spa.bind(this);
    this.refine = this.refine.bind(this);
    this.refinement = this.refinement.bind(this);
    this.superRefine = this.superRefine.bind(this);
    this.optional = this.optional.bind(this);
    this.nullable = this.nullable.bind(this);
    this.nullish = this.nullish.bind(this);
    this.array = this.array.bind(this);
    this.promise = this.promise.bind(this);
    this.or = this.or.bind(this);
    this.and = this.and.bind(this);
    this.transform = this.transform.bind(this);
    this.brand = this.brand.bind(this);
    this.default = this.default.bind(this);
    this.catch = this.catch.bind(this);
    this.describe = this.describe.bind(this);
    this.pipe = this.pipe.bind(this);
    this.readonly = this.readonly.bind(this);
    this.isNullable = this.isNullable.bind(this);
    this.isOptional = this.isOptional.bind(this);
    this["~standard"] = {
      version: 1,
      vendor: "zod",
      validate: (data) => this["~validate"](data)
    };
  }
  optional() {
    return ZodOptional.create(this, this._def);
  }
  nullable() {
    return ZodNullable.create(this, this._def);
  }
  nullish() {
    return this.nullable().optional();
  }
  array() {
    return ZodArray.create(this);
  }
  promise() {
    return ZodPromise.create(this, this._def);
  }
  or(option) {
    return ZodUnion.create([this, option], this._def);
  }
  and(incoming) {
    return ZodIntersection.create(this, incoming, this._def);
  }
  transform(transform) {
    return new ZodEffects({
      ...processCreateParams(this._def),
      schema: this,
      typeName: ZodFirstPartyTypeKind.ZodEffects,
      effect: { type: "transform", transform }
    });
  }
  default(def) {
    const defaultValueFunc = typeof def === "function" ? def : () => def;
    return new ZodDefault({
      ...processCreateParams(this._def),
      innerType: this,
      defaultValue: defaultValueFunc,
      typeName: ZodFirstPartyTypeKind.ZodDefault
    });
  }
  brand() {
    return new ZodBranded({
      typeName: ZodFirstPartyTypeKind.ZodBranded,
      type: this,
      ...processCreateParams(this._def)
    });
  }
  catch(def) {
    const catchValueFunc = typeof def === "function" ? def : () => def;
    return new ZodCatch({
      ...processCreateParams(this._def),
      innerType: this,
      catchValue: catchValueFunc,
      typeName: ZodFirstPartyTypeKind.ZodCatch
    });
  }
  describe(description) {
    const This = this.constructor;
    return new This({
      ...this._def,
      description
    });
  }
  pipe(target) {
    return ZodPipeline.create(this, target);
  }
  readonly() {
    return ZodReadonly.create(this);
  }
  isOptional() {
    return this.safeParse(undefined).success;
  }
  isNullable() {
    return this.safeParse(null).success;
  }
}
var cuidRegex = /^c[^\s-]{8,}$/i;
var cuid2Regex = /^[0-9a-z]+$/;
var ulidRegex = /^[0-9A-HJKMNP-TV-Z]{26}$/i;
var uuidRegex = /^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$/i;
var nanoidRegex = /^[a-z0-9_-]{21}$/i;
var jwtRegex = /^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]*$/;
var durationRegex = /^[-+]?P(?!$)(?:(?:[-+]?\d+Y)|(?:[-+]?\d+[.,]\d+Y$))?(?:(?:[-+]?\d+M)|(?:[-+]?\d+[.,]\d+M$))?(?:(?:[-+]?\d+W)|(?:[-+]?\d+[.,]\d+W$))?(?:(?:[-+]?\d+D)|(?:[-+]?\d+[.,]\d+D$))?(?:T(?=[\d+-])(?:(?:[-+]?\d+H)|(?:[-+]?\d+[.,]\d+H$))?(?:(?:[-+]?\d+M)|(?:[-+]?\d+[.,]\d+M$))?(?:[-+]?\d+(?:[.,]\d+)?S)?)??$/;
var emailRegex = /^(?!\.)(?!.*\.\.)([A-Z0-9_'+\-\.]*)[A-Z0-9_+-]@([A-Z0-9][A-Z0-9\-]*\.)+[A-Z]{2,}$/i;
var _emojiRegex = `^(\\p{Extended_Pictographic}|\\p{Emoji_Component})+$`;
var emojiRegex;
var ipv4Regex = /^(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])$/;
var ipv4CidrRegex = /^(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])\/(3[0-2]|[12]?[0-9])$/;
var ipv6Regex = /^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$/;
var ipv6CidrRegex = /^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))\/(12[0-8]|1[01][0-9]|[1-9]?[0-9])$/;
var base64Regex = /^([0-9a-zA-Z+/]{4})*(([0-9a-zA-Z+/]{2}==)|([0-9a-zA-Z+/]{3}=))?$/;
var base64urlRegex = /^([0-9a-zA-Z-_]{4})*(([0-9a-zA-Z-_]{2}(==)?)|([0-9a-zA-Z-_]{3}(=)?))?$/;
var dateRegexSource = `((\\d\\d[2468][048]|\\d\\d[13579][26]|\\d\\d0[48]|[02468][048]00|[13579][26]00)-02-29|\\d{4}-((0[13578]|1[02])-(0[1-9]|[12]\\d|3[01])|(0[469]|11)-(0[1-9]|[12]\\d|30)|(02)-(0[1-9]|1\\d|2[0-8])))`;
var dateRegex = new RegExp(`^${dateRegexSource}$`);
function timeRegexSource(args) {
  let secondsRegexSource = `[0-5]\\d`;
  if (args.precision) {
    secondsRegexSource = `${secondsRegexSource}\\.\\d{${args.precision}}`;
  } else if (args.precision == null) {
    secondsRegexSource = `${secondsRegexSource}(\\.\\d+)?`;
  }
  const secondsQuantifier = args.precision ? "+" : "?";
  return `([01]\\d|2[0-3]):[0-5]\\d(:${secondsRegexSource})${secondsQuantifier}`;
}
function timeRegex(args) {
  return new RegExp(`^${timeRegexSource(args)}$`);
}
function datetimeRegex(args) {
  let regex = `${dateRegexSource}T${timeRegexSource(args)}`;
  const opts = [];
  opts.push(args.local ? `Z?` : `Z`);
  if (args.offset)
    opts.push(`([+-]\\d{2}:?\\d{2})`);
  regex = `${regex}(${opts.join("|")})`;
  return new RegExp(`^${regex}$`);
}
function isValidIP(ip, version) {
  if ((version === "v4" || !version) && ipv4Regex.test(ip)) {
    return true;
  }
  if ((version === "v6" || !version) && ipv6Regex.test(ip)) {
    return true;
  }
  return false;
}
function isValidJWT(jwt, alg) {
  if (!jwtRegex.test(jwt))
    return false;
  try {
    const [header] = jwt.split(".");
    if (!header)
      return false;
    const base64 = header.replace(/-/g, "+").replace(/_/g, "/").padEnd(header.length + (4 - header.length % 4) % 4, "=");
    const decoded = JSON.parse(atob(base64));
    if (typeof decoded !== "object" || decoded === null)
      return false;
    if ("typ" in decoded && decoded?.typ !== "JWT")
      return false;
    if (!decoded.alg)
      return false;
    if (alg && decoded.alg !== alg)
      return false;
    return true;
  } catch {
    return false;
  }
}
function isValidCidr(ip, version) {
  if ((version === "v4" || !version) && ipv4CidrRegex.test(ip)) {
    return true;
  }
  if ((version === "v6" || !version) && ipv6CidrRegex.test(ip)) {
    return true;
  }
  return false;
}

class ZodString extends ZodType {
  _parse(input) {
    if (this._def.coerce) {
      input.data = String(input.data);
    }
    const parsedType = this._getType(input);
    if (parsedType !== ZodParsedType.string) {
      const ctx2 = this._getOrReturnCtx(input);
      addIssueToContext(ctx2, {
        code: ZodIssueCode.invalid_type,
        expected: ZodParsedType.string,
        received: ctx2.parsedType
      });
      return INVALID;
    }
    const status = new ParseStatus;
    let ctx = undefined;
    for (const check of this._def.checks) {
      if (check.kind === "min") {
        if (input.data.length < check.value) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            code: ZodIssueCode.too_small,
            minimum: check.value,
            type: "string",
            inclusive: true,
            exact: false,
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "max") {
        if (input.data.length > check.value) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            code: ZodIssueCode.too_big,
            maximum: check.value,
            type: "string",
            inclusive: true,
            exact: false,
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "length") {
        const tooBig = input.data.length > check.value;
        const tooSmall = input.data.length < check.value;
        if (tooBig || tooSmall) {
          ctx = this._getOrReturnCtx(input, ctx);
          if (tooBig) {
            addIssueToContext(ctx, {
              code: ZodIssueCode.too_big,
              maximum: check.value,
              type: "string",
              inclusive: true,
              exact: true,
              message: check.message
            });
          } else if (tooSmall) {
            addIssueToContext(ctx, {
              code: ZodIssueCode.too_small,
              minimum: check.value,
              type: "string",
              inclusive: true,
              exact: true,
              message: check.message
            });
          }
          status.dirty();
        }
      } else if (check.kind === "email") {
        if (!emailRegex.test(input.data)) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            validation: "email",
            code: ZodIssueCode.invalid_string,
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "emoji") {
        if (!emojiRegex) {
          emojiRegex = new RegExp(_emojiRegex, "u");
        }
        if (!emojiRegex.test(input.data)) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            validation: "emoji",
            code: ZodIssueCode.invalid_string,
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "uuid") {
        if (!uuidRegex.test(input.data)) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            validation: "uuid",
            code: ZodIssueCode.invalid_string,
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "nanoid") {
        if (!nanoidRegex.test(input.data)) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            validation: "nanoid",
            code: ZodIssueCode.invalid_string,
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "cuid") {
        if (!cuidRegex.test(input.data)) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            validation: "cuid",
            code: ZodIssueCode.invalid_string,
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "cuid2") {
        if (!cuid2Regex.test(input.data)) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            validation: "cuid2",
            code: ZodIssueCode.invalid_string,
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "ulid") {
        if (!ulidRegex.test(input.data)) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            validation: "ulid",
            code: ZodIssueCode.invalid_string,
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "url") {
        try {
          new URL(input.data);
        } catch {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            validation: "url",
            code: ZodIssueCode.invalid_string,
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "regex") {
        check.regex.lastIndex = 0;
        const testResult = check.regex.test(input.data);
        if (!testResult) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            validation: "regex",
            code: ZodIssueCode.invalid_string,
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "trim") {
        input.data = input.data.trim();
      } else if (check.kind === "includes") {
        if (!input.data.includes(check.value, check.position)) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            code: ZodIssueCode.invalid_string,
            validation: { includes: check.value, position: check.position },
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "toLowerCase") {
        input.data = input.data.toLowerCase();
      } else if (check.kind === "toUpperCase") {
        input.data = input.data.toUpperCase();
      } else if (check.kind === "startsWith") {
        if (!input.data.startsWith(check.value)) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            code: ZodIssueCode.invalid_string,
            validation: { startsWith: check.value },
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "endsWith") {
        if (!input.data.endsWith(check.value)) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            code: ZodIssueCode.invalid_string,
            validation: { endsWith: check.value },
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "datetime") {
        const regex = datetimeRegex(check);
        if (!regex.test(input.data)) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            code: ZodIssueCode.invalid_string,
            validation: "datetime",
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "date") {
        const regex = dateRegex;
        if (!regex.test(input.data)) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            code: ZodIssueCode.invalid_string,
            validation: "date",
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "time") {
        const regex = timeRegex(check);
        if (!regex.test(input.data)) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            code: ZodIssueCode.invalid_string,
            validation: "time",
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "duration") {
        if (!durationRegex.test(input.data)) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            validation: "duration",
            code: ZodIssueCode.invalid_string,
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "ip") {
        if (!isValidIP(input.data, check.version)) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            validation: "ip",
            code: ZodIssueCode.invalid_string,
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "jwt") {
        if (!isValidJWT(input.data, check.alg)) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            validation: "jwt",
            code: ZodIssueCode.invalid_string,
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "cidr") {
        if (!isValidCidr(input.data, check.version)) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            validation: "cidr",
            code: ZodIssueCode.invalid_string,
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "base64") {
        if (!base64Regex.test(input.data)) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            validation: "base64",
            code: ZodIssueCode.invalid_string,
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "base64url") {
        if (!base64urlRegex.test(input.data)) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            validation: "base64url",
            code: ZodIssueCode.invalid_string,
            message: check.message
          });
          status.dirty();
        }
      } else {
        util.assertNever(check);
      }
    }
    return { status: status.value, value: input.data };
  }
  _regex(regex, validation, message) {
    return this.refinement((data) => regex.test(data), {
      validation,
      code: ZodIssueCode.invalid_string,
      ...errorUtil.errToObj(message)
    });
  }
  _addCheck(check) {
    return new ZodString({
      ...this._def,
      checks: [...this._def.checks, check]
    });
  }
  email(message) {
    return this._addCheck({ kind: "email", ...errorUtil.errToObj(message) });
  }
  url(message) {
    return this._addCheck({ kind: "url", ...errorUtil.errToObj(message) });
  }
  emoji(message) {
    return this._addCheck({ kind: "emoji", ...errorUtil.errToObj(message) });
  }
  uuid(message) {
    return this._addCheck({ kind: "uuid", ...errorUtil.errToObj(message) });
  }
  nanoid(message) {
    return this._addCheck({ kind: "nanoid", ...errorUtil.errToObj(message) });
  }
  cuid(message) {
    return this._addCheck({ kind: "cuid", ...errorUtil.errToObj(message) });
  }
  cuid2(message) {
    return this._addCheck({ kind: "cuid2", ...errorUtil.errToObj(message) });
  }
  ulid(message) {
    return this._addCheck({ kind: "ulid", ...errorUtil.errToObj(message) });
  }
  base64(message) {
    return this._addCheck({ kind: "base64", ...errorUtil.errToObj(message) });
  }
  base64url(message) {
    return this._addCheck({
      kind: "base64url",
      ...errorUtil.errToObj(message)
    });
  }
  jwt(options) {
    return this._addCheck({ kind: "jwt", ...errorUtil.errToObj(options) });
  }
  ip(options) {
    return this._addCheck({ kind: "ip", ...errorUtil.errToObj(options) });
  }
  cidr(options) {
    return this._addCheck({ kind: "cidr", ...errorUtil.errToObj(options) });
  }
  datetime(options) {
    if (typeof options === "string") {
      return this._addCheck({
        kind: "datetime",
        precision: null,
        offset: false,
        local: false,
        message: options
      });
    }
    return this._addCheck({
      kind: "datetime",
      precision: typeof options?.precision === "undefined" ? null : options?.precision,
      offset: options?.offset ?? false,
      local: options?.local ?? false,
      ...errorUtil.errToObj(options?.message)
    });
  }
  date(message) {
    return this._addCheck({ kind: "date", message });
  }
  time(options) {
    if (typeof options === "string") {
      return this._addCheck({
        kind: "time",
        precision: null,
        message: options
      });
    }
    return this._addCheck({
      kind: "time",
      precision: typeof options?.precision === "undefined" ? null : options?.precision,
      ...errorUtil.errToObj(options?.message)
    });
  }
  duration(message) {
    return this._addCheck({ kind: "duration", ...errorUtil.errToObj(message) });
  }
  regex(regex, message) {
    return this._addCheck({
      kind: "regex",
      regex,
      ...errorUtil.errToObj(message)
    });
  }
  includes(value2, options) {
    return this._addCheck({
      kind: "includes",
      value: value2,
      position: options?.position,
      ...errorUtil.errToObj(options?.message)
    });
  }
  startsWith(value2, message) {
    return this._addCheck({
      kind: "startsWith",
      value: value2,
      ...errorUtil.errToObj(message)
    });
  }
  endsWith(value2, message) {
    return this._addCheck({
      kind: "endsWith",
      value: value2,
      ...errorUtil.errToObj(message)
    });
  }
  min(minLength, message) {
    return this._addCheck({
      kind: "min",
      value: minLength,
      ...errorUtil.errToObj(message)
    });
  }
  max(maxLength, message) {
    return this._addCheck({
      kind: "max",
      value: maxLength,
      ...errorUtil.errToObj(message)
    });
  }
  length(len2, message) {
    return this._addCheck({
      kind: "length",
      value: len2,
      ...errorUtil.errToObj(message)
    });
  }
  nonempty(message) {
    return this.min(1, errorUtil.errToObj(message));
  }
  trim() {
    return new ZodString({
      ...this._def,
      checks: [...this._def.checks, { kind: "trim" }]
    });
  }
  toLowerCase() {
    return new ZodString({
      ...this._def,
      checks: [...this._def.checks, { kind: "toLowerCase" }]
    });
  }
  toUpperCase() {
    return new ZodString({
      ...this._def,
      checks: [...this._def.checks, { kind: "toUpperCase" }]
    });
  }
  get isDatetime() {
    return !!this._def.checks.find((ch) => ch.kind === "datetime");
  }
  get isDate() {
    return !!this._def.checks.find((ch) => ch.kind === "date");
  }
  get isTime() {
    return !!this._def.checks.find((ch) => ch.kind === "time");
  }
  get isDuration() {
    return !!this._def.checks.find((ch) => ch.kind === "duration");
  }
  get isEmail() {
    return !!this._def.checks.find((ch) => ch.kind === "email");
  }
  get isURL() {
    return !!this._def.checks.find((ch) => ch.kind === "url");
  }
  get isEmoji() {
    return !!this._def.checks.find((ch) => ch.kind === "emoji");
  }
  get isUUID() {
    return !!this._def.checks.find((ch) => ch.kind === "uuid");
  }
  get isNANOID() {
    return !!this._def.checks.find((ch) => ch.kind === "nanoid");
  }
  get isCUID() {
    return !!this._def.checks.find((ch) => ch.kind === "cuid");
  }
  get isCUID2() {
    return !!this._def.checks.find((ch) => ch.kind === "cuid2");
  }
  get isULID() {
    return !!this._def.checks.find((ch) => ch.kind === "ulid");
  }
  get isIP() {
    return !!this._def.checks.find((ch) => ch.kind === "ip");
  }
  get isCIDR() {
    return !!this._def.checks.find((ch) => ch.kind === "cidr");
  }
  get isBase64() {
    return !!this._def.checks.find((ch) => ch.kind === "base64");
  }
  get isBase64url() {
    return !!this._def.checks.find((ch) => ch.kind === "base64url");
  }
  get minLength() {
    let min = null;
    for (const ch of this._def.checks) {
      if (ch.kind === "min") {
        if (min === null || ch.value > min)
          min = ch.value;
      }
    }
    return min;
  }
  get maxLength() {
    let max = null;
    for (const ch of this._def.checks) {
      if (ch.kind === "max") {
        if (max === null || ch.value < max)
          max = ch.value;
      }
    }
    return max;
  }
}
ZodString.create = (params) => {
  return new ZodString({
    checks: [],
    typeName: ZodFirstPartyTypeKind.ZodString,
    coerce: params?.coerce ?? false,
    ...processCreateParams(params)
  });
};
function floatSafeRemainder(val, step) {
  const valDecCount = (val.toString().split(".")[1] || "").length;
  const stepDecCount = (step.toString().split(".")[1] || "").length;
  const decCount = valDecCount > stepDecCount ? valDecCount : stepDecCount;
  const valInt = Number.parseInt(val.toFixed(decCount).replace(".", ""));
  const stepInt = Number.parseInt(step.toFixed(decCount).replace(".", ""));
  return valInt % stepInt / 10 ** decCount;
}

class ZodNumber extends ZodType {
  constructor() {
    super(...arguments);
    this.min = this.gte;
    this.max = this.lte;
    this.step = this.multipleOf;
  }
  _parse(input) {
    if (this._def.coerce) {
      input.data = Number(input.data);
    }
    const parsedType = this._getType(input);
    if (parsedType !== ZodParsedType.number) {
      const ctx2 = this._getOrReturnCtx(input);
      addIssueToContext(ctx2, {
        code: ZodIssueCode.invalid_type,
        expected: ZodParsedType.number,
        received: ctx2.parsedType
      });
      return INVALID;
    }
    let ctx = undefined;
    const status = new ParseStatus;
    for (const check of this._def.checks) {
      if (check.kind === "int") {
        if (!util.isInteger(input.data)) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            code: ZodIssueCode.invalid_type,
            expected: "integer",
            received: "float",
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "min") {
        const tooSmall = check.inclusive ? input.data < check.value : input.data <= check.value;
        if (tooSmall) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            code: ZodIssueCode.too_small,
            minimum: check.value,
            type: "number",
            inclusive: check.inclusive,
            exact: false,
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "max") {
        const tooBig = check.inclusive ? input.data > check.value : input.data >= check.value;
        if (tooBig) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            code: ZodIssueCode.too_big,
            maximum: check.value,
            type: "number",
            inclusive: check.inclusive,
            exact: false,
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "multipleOf") {
        if (floatSafeRemainder(input.data, check.value) !== 0) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            code: ZodIssueCode.not_multiple_of,
            multipleOf: check.value,
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "finite") {
        if (!Number.isFinite(input.data)) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            code: ZodIssueCode.not_finite,
            message: check.message
          });
          status.dirty();
        }
      } else {
        util.assertNever(check);
      }
    }
    return { status: status.value, value: input.data };
  }
  gte(value2, message) {
    return this.setLimit("min", value2, true, errorUtil.toString(message));
  }
  gt(value2, message) {
    return this.setLimit("min", value2, false, errorUtil.toString(message));
  }
  lte(value2, message) {
    return this.setLimit("max", value2, true, errorUtil.toString(message));
  }
  lt(value2, message) {
    return this.setLimit("max", value2, false, errorUtil.toString(message));
  }
  setLimit(kind, value2, inclusive, message) {
    return new ZodNumber({
      ...this._def,
      checks: [
        ...this._def.checks,
        {
          kind,
          value: value2,
          inclusive,
          message: errorUtil.toString(message)
        }
      ]
    });
  }
  _addCheck(check) {
    return new ZodNumber({
      ...this._def,
      checks: [...this._def.checks, check]
    });
  }
  int(message) {
    return this._addCheck({
      kind: "int",
      message: errorUtil.toString(message)
    });
  }
  positive(message) {
    return this._addCheck({
      kind: "min",
      value: 0,
      inclusive: false,
      message: errorUtil.toString(message)
    });
  }
  negative(message) {
    return this._addCheck({
      kind: "max",
      value: 0,
      inclusive: false,
      message: errorUtil.toString(message)
    });
  }
  nonpositive(message) {
    return this._addCheck({
      kind: "max",
      value: 0,
      inclusive: true,
      message: errorUtil.toString(message)
    });
  }
  nonnegative(message) {
    return this._addCheck({
      kind: "min",
      value: 0,
      inclusive: true,
      message: errorUtil.toString(message)
    });
  }
  multipleOf(value2, message) {
    return this._addCheck({
      kind: "multipleOf",
      value: value2,
      message: errorUtil.toString(message)
    });
  }
  finite(message) {
    return this._addCheck({
      kind: "finite",
      message: errorUtil.toString(message)
    });
  }
  safe(message) {
    return this._addCheck({
      kind: "min",
      inclusive: true,
      value: Number.MIN_SAFE_INTEGER,
      message: errorUtil.toString(message)
    })._addCheck({
      kind: "max",
      inclusive: true,
      value: Number.MAX_SAFE_INTEGER,
      message: errorUtil.toString(message)
    });
  }
  get minValue() {
    let min = null;
    for (const ch of this._def.checks) {
      if (ch.kind === "min") {
        if (min === null || ch.value > min)
          min = ch.value;
      }
    }
    return min;
  }
  get maxValue() {
    let max = null;
    for (const ch of this._def.checks) {
      if (ch.kind === "max") {
        if (max === null || ch.value < max)
          max = ch.value;
      }
    }
    return max;
  }
  get isInt() {
    return !!this._def.checks.find((ch) => ch.kind === "int" || ch.kind === "multipleOf" && util.isInteger(ch.value));
  }
  get isFinite() {
    let max = null;
    let min = null;
    for (const ch of this._def.checks) {
      if (ch.kind === "finite" || ch.kind === "int" || ch.kind === "multipleOf") {
        return true;
      } else if (ch.kind === "min") {
        if (min === null || ch.value > min)
          min = ch.value;
      } else if (ch.kind === "max") {
        if (max === null || ch.value < max)
          max = ch.value;
      }
    }
    return Number.isFinite(min) && Number.isFinite(max);
  }
}
ZodNumber.create = (params) => {
  return new ZodNumber({
    checks: [],
    typeName: ZodFirstPartyTypeKind.ZodNumber,
    coerce: params?.coerce || false,
    ...processCreateParams(params)
  });
};

class ZodBigInt extends ZodType {
  constructor() {
    super(...arguments);
    this.min = this.gte;
    this.max = this.lte;
  }
  _parse(input) {
    if (this._def.coerce) {
      try {
        input.data = BigInt(input.data);
      } catch {
        return this._getInvalidInput(input);
      }
    }
    const parsedType = this._getType(input);
    if (parsedType !== ZodParsedType.bigint) {
      return this._getInvalidInput(input);
    }
    let ctx = undefined;
    const status = new ParseStatus;
    for (const check of this._def.checks) {
      if (check.kind === "min") {
        const tooSmall = check.inclusive ? input.data < check.value : input.data <= check.value;
        if (tooSmall) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            code: ZodIssueCode.too_small,
            type: "bigint",
            minimum: check.value,
            inclusive: check.inclusive,
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "max") {
        const tooBig = check.inclusive ? input.data > check.value : input.data >= check.value;
        if (tooBig) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            code: ZodIssueCode.too_big,
            type: "bigint",
            maximum: check.value,
            inclusive: check.inclusive,
            message: check.message
          });
          status.dirty();
        }
      } else if (check.kind === "multipleOf") {
        if (input.data % check.value !== BigInt(0)) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            code: ZodIssueCode.not_multiple_of,
            multipleOf: check.value,
            message: check.message
          });
          status.dirty();
        }
      } else {
        util.assertNever(check);
      }
    }
    return { status: status.value, value: input.data };
  }
  _getInvalidInput(input) {
    const ctx = this._getOrReturnCtx(input);
    addIssueToContext(ctx, {
      code: ZodIssueCode.invalid_type,
      expected: ZodParsedType.bigint,
      received: ctx.parsedType
    });
    return INVALID;
  }
  gte(value2, message) {
    return this.setLimit("min", value2, true, errorUtil.toString(message));
  }
  gt(value2, message) {
    return this.setLimit("min", value2, false, errorUtil.toString(message));
  }
  lte(value2, message) {
    return this.setLimit("max", value2, true, errorUtil.toString(message));
  }
  lt(value2, message) {
    return this.setLimit("max", value2, false, errorUtil.toString(message));
  }
  setLimit(kind, value2, inclusive, message) {
    return new ZodBigInt({
      ...this._def,
      checks: [
        ...this._def.checks,
        {
          kind,
          value: value2,
          inclusive,
          message: errorUtil.toString(message)
        }
      ]
    });
  }
  _addCheck(check) {
    return new ZodBigInt({
      ...this._def,
      checks: [...this._def.checks, check]
    });
  }
  positive(message) {
    return this._addCheck({
      kind: "min",
      value: BigInt(0),
      inclusive: false,
      message: errorUtil.toString(message)
    });
  }
  negative(message) {
    return this._addCheck({
      kind: "max",
      value: BigInt(0),
      inclusive: false,
      message: errorUtil.toString(message)
    });
  }
  nonpositive(message) {
    return this._addCheck({
      kind: "max",
      value: BigInt(0),
      inclusive: true,
      message: errorUtil.toString(message)
    });
  }
  nonnegative(message) {
    return this._addCheck({
      kind: "min",
      value: BigInt(0),
      inclusive: true,
      message: errorUtil.toString(message)
    });
  }
  multipleOf(value2, message) {
    return this._addCheck({
      kind: "multipleOf",
      value: value2,
      message: errorUtil.toString(message)
    });
  }
  get minValue() {
    let min = null;
    for (const ch of this._def.checks) {
      if (ch.kind === "min") {
        if (min === null || ch.value > min)
          min = ch.value;
      }
    }
    return min;
  }
  get maxValue() {
    let max = null;
    for (const ch of this._def.checks) {
      if (ch.kind === "max") {
        if (max === null || ch.value < max)
          max = ch.value;
      }
    }
    return max;
  }
}
ZodBigInt.create = (params) => {
  return new ZodBigInt({
    checks: [],
    typeName: ZodFirstPartyTypeKind.ZodBigInt,
    coerce: params?.coerce ?? false,
    ...processCreateParams(params)
  });
};

class ZodBoolean extends ZodType {
  _parse(input) {
    if (this._def.coerce) {
      input.data = Boolean(input.data);
    }
    const parsedType = this._getType(input);
    if (parsedType !== ZodParsedType.boolean) {
      const ctx = this._getOrReturnCtx(input);
      addIssueToContext(ctx, {
        code: ZodIssueCode.invalid_type,
        expected: ZodParsedType.boolean,
        received: ctx.parsedType
      });
      return INVALID;
    }
    return OK(input.data);
  }
}
ZodBoolean.create = (params) => {
  return new ZodBoolean({
    typeName: ZodFirstPartyTypeKind.ZodBoolean,
    coerce: params?.coerce || false,
    ...processCreateParams(params)
  });
};

class ZodDate extends ZodType {
  _parse(input) {
    if (this._def.coerce) {
      input.data = new Date(input.data);
    }
    const parsedType = this._getType(input);
    if (parsedType !== ZodParsedType.date) {
      const ctx2 = this._getOrReturnCtx(input);
      addIssueToContext(ctx2, {
        code: ZodIssueCode.invalid_type,
        expected: ZodParsedType.date,
        received: ctx2.parsedType
      });
      return INVALID;
    }
    if (Number.isNaN(input.data.getTime())) {
      const ctx2 = this._getOrReturnCtx(input);
      addIssueToContext(ctx2, {
        code: ZodIssueCode.invalid_date
      });
      return INVALID;
    }
    const status = new ParseStatus;
    let ctx = undefined;
    for (const check of this._def.checks) {
      if (check.kind === "min") {
        if (input.data.getTime() < check.value) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            code: ZodIssueCode.too_small,
            message: check.message,
            inclusive: true,
            exact: false,
            minimum: check.value,
            type: "date"
          });
          status.dirty();
        }
      } else if (check.kind === "max") {
        if (input.data.getTime() > check.value) {
          ctx = this._getOrReturnCtx(input, ctx);
          addIssueToContext(ctx, {
            code: ZodIssueCode.too_big,
            message: check.message,
            inclusive: true,
            exact: false,
            maximum: check.value,
            type: "date"
          });
          status.dirty();
        }
      } else {
        util.assertNever(check);
      }
    }
    return {
      status: status.value,
      value: new Date(input.data.getTime())
    };
  }
  _addCheck(check) {
    return new ZodDate({
      ...this._def,
      checks: [...this._def.checks, check]
    });
  }
  min(minDate, message) {
    return this._addCheck({
      kind: "min",
      value: minDate.getTime(),
      message: errorUtil.toString(message)
    });
  }
  max(maxDate, message) {
    return this._addCheck({
      kind: "max",
      value: maxDate.getTime(),
      message: errorUtil.toString(message)
    });
  }
  get minDate() {
    let min = null;
    for (const ch of this._def.checks) {
      if (ch.kind === "min") {
        if (min === null || ch.value > min)
          min = ch.value;
      }
    }
    return min != null ? new Date(min) : null;
  }
  get maxDate() {
    let max = null;
    for (const ch of this._def.checks) {
      if (ch.kind === "max") {
        if (max === null || ch.value < max)
          max = ch.value;
      }
    }
    return max != null ? new Date(max) : null;
  }
}
ZodDate.create = (params) => {
  return new ZodDate({
    checks: [],
    coerce: params?.coerce || false,
    typeName: ZodFirstPartyTypeKind.ZodDate,
    ...processCreateParams(params)
  });
};

class ZodSymbol extends ZodType {
  _parse(input) {
    const parsedType = this._getType(input);
    if (parsedType !== ZodParsedType.symbol) {
      const ctx = this._getOrReturnCtx(input);
      addIssueToContext(ctx, {
        code: ZodIssueCode.invalid_type,
        expected: ZodParsedType.symbol,
        received: ctx.parsedType
      });
      return INVALID;
    }
    return OK(input.data);
  }
}
ZodSymbol.create = (params) => {
  return new ZodSymbol({
    typeName: ZodFirstPartyTypeKind.ZodSymbol,
    ...processCreateParams(params)
  });
};

class ZodUndefined extends ZodType {
  _parse(input) {
    const parsedType = this._getType(input);
    if (parsedType !== ZodParsedType.undefined) {
      const ctx = this._getOrReturnCtx(input);
      addIssueToContext(ctx, {
        code: ZodIssueCode.invalid_type,
        expected: ZodParsedType.undefined,
        received: ctx.parsedType
      });
      return INVALID;
    }
    return OK(input.data);
  }
}
ZodUndefined.create = (params) => {
  return new ZodUndefined({
    typeName: ZodFirstPartyTypeKind.ZodUndefined,
    ...processCreateParams(params)
  });
};

class ZodNull extends ZodType {
  _parse(input) {
    const parsedType = this._getType(input);
    if (parsedType !== ZodParsedType.null) {
      const ctx = this._getOrReturnCtx(input);
      addIssueToContext(ctx, {
        code: ZodIssueCode.invalid_type,
        expected: ZodParsedType.null,
        received: ctx.parsedType
      });
      return INVALID;
    }
    return OK(input.data);
  }
}
ZodNull.create = (params) => {
  return new ZodNull({
    typeName: ZodFirstPartyTypeKind.ZodNull,
    ...processCreateParams(params)
  });
};

class ZodAny extends ZodType {
  constructor() {
    super(...arguments);
    this._any = true;
  }
  _parse(input) {
    return OK(input.data);
  }
}
ZodAny.create = (params) => {
  return new ZodAny({
    typeName: ZodFirstPartyTypeKind.ZodAny,
    ...processCreateParams(params)
  });
};

class ZodUnknown extends ZodType {
  constructor() {
    super(...arguments);
    this._unknown = true;
  }
  _parse(input) {
    return OK(input.data);
  }
}
ZodUnknown.create = (params) => {
  return new ZodUnknown({
    typeName: ZodFirstPartyTypeKind.ZodUnknown,
    ...processCreateParams(params)
  });
};

class ZodNever extends ZodType {
  _parse(input) {
    const ctx = this._getOrReturnCtx(input);
    addIssueToContext(ctx, {
      code: ZodIssueCode.invalid_type,
      expected: ZodParsedType.never,
      received: ctx.parsedType
    });
    return INVALID;
  }
}
ZodNever.create = (params) => {
  return new ZodNever({
    typeName: ZodFirstPartyTypeKind.ZodNever,
    ...processCreateParams(params)
  });
};

class ZodVoid extends ZodType {
  _parse(input) {
    const parsedType = this._getType(input);
    if (parsedType !== ZodParsedType.undefined) {
      const ctx = this._getOrReturnCtx(input);
      addIssueToContext(ctx, {
        code: ZodIssueCode.invalid_type,
        expected: ZodParsedType.void,
        received: ctx.parsedType
      });
      return INVALID;
    }
    return OK(input.data);
  }
}
ZodVoid.create = (params) => {
  return new ZodVoid({
    typeName: ZodFirstPartyTypeKind.ZodVoid,
    ...processCreateParams(params)
  });
};

class ZodArray extends ZodType {
  _parse(input) {
    const { ctx, status } = this._processInputParams(input);
    const def = this._def;
    if (ctx.parsedType !== ZodParsedType.array) {
      addIssueToContext(ctx, {
        code: ZodIssueCode.invalid_type,
        expected: ZodParsedType.array,
        received: ctx.parsedType
      });
      return INVALID;
    }
    if (def.exactLength !== null) {
      const tooBig = ctx.data.length > def.exactLength.value;
      const tooSmall = ctx.data.length < def.exactLength.value;
      if (tooBig || tooSmall) {
        addIssueToContext(ctx, {
          code: tooBig ? ZodIssueCode.too_big : ZodIssueCode.too_small,
          minimum: tooSmall ? def.exactLength.value : undefined,
          maximum: tooBig ? def.exactLength.value : undefined,
          type: "array",
          inclusive: true,
          exact: true,
          message: def.exactLength.message
        });
        status.dirty();
      }
    }
    if (def.minLength !== null) {
      if (ctx.data.length < def.minLength.value) {
        addIssueToContext(ctx, {
          code: ZodIssueCode.too_small,
          minimum: def.minLength.value,
          type: "array",
          inclusive: true,
          exact: false,
          message: def.minLength.message
        });
        status.dirty();
      }
    }
    if (def.maxLength !== null) {
      if (ctx.data.length > def.maxLength.value) {
        addIssueToContext(ctx, {
          code: ZodIssueCode.too_big,
          maximum: def.maxLength.value,
          type: "array",
          inclusive: true,
          exact: false,
          message: def.maxLength.message
        });
        status.dirty();
      }
    }
    if (ctx.common.async) {
      return Promise.all([...ctx.data].map((item, i2) => {
        return def.type._parseAsync(new ParseInputLazyPath(ctx, item, ctx.path, i2));
      })).then((result2) => {
        return ParseStatus.mergeArray(status, result2);
      });
    }
    const result = [...ctx.data].map((item, i2) => {
      return def.type._parseSync(new ParseInputLazyPath(ctx, item, ctx.path, i2));
    });
    return ParseStatus.mergeArray(status, result);
  }
  get element() {
    return this._def.type;
  }
  min(minLength, message) {
    return new ZodArray({
      ...this._def,
      minLength: { value: minLength, message: errorUtil.toString(message) }
    });
  }
  max(maxLength, message) {
    return new ZodArray({
      ...this._def,
      maxLength: { value: maxLength, message: errorUtil.toString(message) }
    });
  }
  length(len2, message) {
    return new ZodArray({
      ...this._def,
      exactLength: { value: len2, message: errorUtil.toString(message) }
    });
  }
  nonempty(message) {
    return this.min(1, message);
  }
}
ZodArray.create = (schema, params) => {
  return new ZodArray({
    type: schema,
    minLength: null,
    maxLength: null,
    exactLength: null,
    typeName: ZodFirstPartyTypeKind.ZodArray,
    ...processCreateParams(params)
  });
};
function deepPartialify(schema) {
  if (schema instanceof ZodObject) {
    const newShape = {};
    for (const key in schema.shape) {
      const fieldSchema = schema.shape[key];
      newShape[key] = ZodOptional.create(deepPartialify(fieldSchema));
    }
    return new ZodObject({
      ...schema._def,
      shape: () => newShape
    });
  } else if (schema instanceof ZodArray) {
    return new ZodArray({
      ...schema._def,
      type: deepPartialify(schema.element)
    });
  } else if (schema instanceof ZodOptional) {
    return ZodOptional.create(deepPartialify(schema.unwrap()));
  } else if (schema instanceof ZodNullable) {
    return ZodNullable.create(deepPartialify(schema.unwrap()));
  } else if (schema instanceof ZodTuple) {
    return ZodTuple.create(schema.items.map((item) => deepPartialify(item)));
  } else {
    return schema;
  }
}

class ZodObject extends ZodType {
  constructor() {
    super(...arguments);
    this._cached = null;
    this.nonstrict = this.passthrough;
    this.augment = this.extend;
  }
  _getCached() {
    if (this._cached !== null)
      return this._cached;
    const shape = this._def.shape();
    const keys = util.objectKeys(shape);
    this._cached = { shape, keys };
    return this._cached;
  }
  _parse(input) {
    const parsedType = this._getType(input);
    if (parsedType !== ZodParsedType.object) {
      const ctx2 = this._getOrReturnCtx(input);
      addIssueToContext(ctx2, {
        code: ZodIssueCode.invalid_type,
        expected: ZodParsedType.object,
        received: ctx2.parsedType
      });
      return INVALID;
    }
    const { status, ctx } = this._processInputParams(input);
    const { shape, keys: shapeKeys } = this._getCached();
    const extraKeys = [];
    if (!(this._def.catchall instanceof ZodNever && this._def.unknownKeys === "strip")) {
      for (const key in ctx.data) {
        if (!shapeKeys.includes(key)) {
          extraKeys.push(key);
        }
      }
    }
    const pairs = [];
    for (const key of shapeKeys) {
      const keyValidator = shape[key];
      const value2 = ctx.data[key];
      pairs.push({
        key: { status: "valid", value: key },
        value: keyValidator._parse(new ParseInputLazyPath(ctx, value2, ctx.path, key)),
        alwaysSet: key in ctx.data
      });
    }
    if (this._def.catchall instanceof ZodNever) {
      const unknownKeys = this._def.unknownKeys;
      if (unknownKeys === "passthrough") {
        for (const key of extraKeys) {
          pairs.push({
            key: { status: "valid", value: key },
            value: { status: "valid", value: ctx.data[key] }
          });
        }
      } else if (unknownKeys === "strict") {
        if (extraKeys.length > 0) {
          addIssueToContext(ctx, {
            code: ZodIssueCode.unrecognized_keys,
            keys: extraKeys
          });
          status.dirty();
        }
      } else if (unknownKeys === "strip") {} else {
        throw new Error(`Internal ZodObject error: invalid unknownKeys value.`);
      }
    } else {
      const catchall = this._def.catchall;
      for (const key of extraKeys) {
        const value2 = ctx.data[key];
        pairs.push({
          key: { status: "valid", value: key },
          value: catchall._parse(new ParseInputLazyPath(ctx, value2, ctx.path, key)),
          alwaysSet: key in ctx.data
        });
      }
    }
    if (ctx.common.async) {
      return Promise.resolve().then(async () => {
        const syncPairs = [];
        for (const pair of pairs) {
          const key = await pair.key;
          const value2 = await pair.value;
          syncPairs.push({
            key,
            value: value2,
            alwaysSet: pair.alwaysSet
          });
        }
        return syncPairs;
      }).then((syncPairs) => {
        return ParseStatus.mergeObjectSync(status, syncPairs);
      });
    } else {
      return ParseStatus.mergeObjectSync(status, pairs);
    }
  }
  get shape() {
    return this._def.shape();
  }
  strict(message) {
    errorUtil.errToObj;
    return new ZodObject({
      ...this._def,
      unknownKeys: "strict",
      ...message !== undefined ? {
        errorMap: (issue, ctx) => {
          const defaultError = this._def.errorMap?.(issue, ctx).message ?? ctx.defaultError;
          if (issue.code === "unrecognized_keys")
            return {
              message: errorUtil.errToObj(message).message ?? defaultError
            };
          return {
            message: defaultError
          };
        }
      } : {}
    });
  }
  strip() {
    return new ZodObject({
      ...this._def,
      unknownKeys: "strip"
    });
  }
  passthrough() {
    return new ZodObject({
      ...this._def,
      unknownKeys: "passthrough"
    });
  }
  extend(augmentation) {
    return new ZodObject({
      ...this._def,
      shape: () => ({
        ...this._def.shape(),
        ...augmentation
      })
    });
  }
  merge(merging) {
    const merged = new ZodObject({
      unknownKeys: merging._def.unknownKeys,
      catchall: merging._def.catchall,
      shape: () => ({
        ...this._def.shape(),
        ...merging._def.shape()
      }),
      typeName: ZodFirstPartyTypeKind.ZodObject
    });
    return merged;
  }
  setKey(key, schema) {
    return this.augment({ [key]: schema });
  }
  catchall(index) {
    return new ZodObject({
      ...this._def,
      catchall: index
    });
  }
  pick(mask) {
    const shape = {};
    for (const key of util.objectKeys(mask)) {
      if (mask[key] && this.shape[key]) {
        shape[key] = this.shape[key];
      }
    }
    return new ZodObject({
      ...this._def,
      shape: () => shape
    });
  }
  omit(mask) {
    const shape = {};
    for (const key of util.objectKeys(this.shape)) {
      if (!mask[key]) {
        shape[key] = this.shape[key];
      }
    }
    return new ZodObject({
      ...this._def,
      shape: () => shape
    });
  }
  deepPartial() {
    return deepPartialify(this);
  }
  partial(mask) {
    const newShape = {};
    for (const key of util.objectKeys(this.shape)) {
      const fieldSchema = this.shape[key];
      if (mask && !mask[key]) {
        newShape[key] = fieldSchema;
      } else {
        newShape[key] = fieldSchema.optional();
      }
    }
    return new ZodObject({
      ...this._def,
      shape: () => newShape
    });
  }
  required(mask) {
    const newShape = {};
    for (const key of util.objectKeys(this.shape)) {
      if (mask && !mask[key]) {
        newShape[key] = this.shape[key];
      } else {
        const fieldSchema = this.shape[key];
        let newField2 = fieldSchema;
        while (newField2 instanceof ZodOptional) {
          newField2 = newField2._def.innerType;
        }
        newShape[key] = newField2;
      }
    }
    return new ZodObject({
      ...this._def,
      shape: () => newShape
    });
  }
  keyof() {
    return createZodEnum(util.objectKeys(this.shape));
  }
}
ZodObject.create = (shape, params) => {
  return new ZodObject({
    shape: () => shape,
    unknownKeys: "strip",
    catchall: ZodNever.create(),
    typeName: ZodFirstPartyTypeKind.ZodObject,
    ...processCreateParams(params)
  });
};
ZodObject.strictCreate = (shape, params) => {
  return new ZodObject({
    shape: () => shape,
    unknownKeys: "strict",
    catchall: ZodNever.create(),
    typeName: ZodFirstPartyTypeKind.ZodObject,
    ...processCreateParams(params)
  });
};
ZodObject.lazycreate = (shape, params) => {
  return new ZodObject({
    shape,
    unknownKeys: "strip",
    catchall: ZodNever.create(),
    typeName: ZodFirstPartyTypeKind.ZodObject,
    ...processCreateParams(params)
  });
};

class ZodUnion extends ZodType {
  _parse(input) {
    const { ctx } = this._processInputParams(input);
    const options = this._def.options;
    function handleResults(results) {
      for (const result of results) {
        if (result.result.status === "valid") {
          return result.result;
        }
      }
      for (const result of results) {
        if (result.result.status === "dirty") {
          ctx.common.issues.push(...result.ctx.common.issues);
          return result.result;
        }
      }
      const unionErrors = results.map((result) => new ZodError(result.ctx.common.issues));
      addIssueToContext(ctx, {
        code: ZodIssueCode.invalid_union,
        unionErrors
      });
      return INVALID;
    }
    if (ctx.common.async) {
      return Promise.all(options.map(async (option) => {
        const childCtx = {
          ...ctx,
          common: {
            ...ctx.common,
            issues: []
          },
          parent: null
        };
        return {
          result: await option._parseAsync({
            data: ctx.data,
            path: ctx.path,
            parent: childCtx
          }),
          ctx: childCtx
        };
      })).then(handleResults);
    } else {
      let dirty = undefined;
      const issues = [];
      for (const option of options) {
        const childCtx = {
          ...ctx,
          common: {
            ...ctx.common,
            issues: []
          },
          parent: null
        };
        const result = option._parseSync({
          data: ctx.data,
          path: ctx.path,
          parent: childCtx
        });
        if (result.status === "valid") {
          return result;
        } else if (result.status === "dirty" && !dirty) {
          dirty = { result, ctx: childCtx };
        }
        if (childCtx.common.issues.length) {
          issues.push(childCtx.common.issues);
        }
      }
      if (dirty) {
        ctx.common.issues.push(...dirty.ctx.common.issues);
        return dirty.result;
      }
      const unionErrors = issues.map((issues2) => new ZodError(issues2));
      addIssueToContext(ctx, {
        code: ZodIssueCode.invalid_union,
        unionErrors
      });
      return INVALID;
    }
  }
  get options() {
    return this._def.options;
  }
}
ZodUnion.create = (types3, params) => {
  return new ZodUnion({
    options: types3,
    typeName: ZodFirstPartyTypeKind.ZodUnion,
    ...processCreateParams(params)
  });
};
var getDiscriminator = (type) => {
  if (type instanceof ZodLazy) {
    return getDiscriminator(type.schema);
  } else if (type instanceof ZodEffects) {
    return getDiscriminator(type.innerType());
  } else if (type instanceof ZodLiteral) {
    return [type.value];
  } else if (type instanceof ZodEnum) {
    return type.options;
  } else if (type instanceof ZodNativeEnum) {
    return util.objectValues(type.enum);
  } else if (type instanceof ZodDefault) {
    return getDiscriminator(type._def.innerType);
  } else if (type instanceof ZodUndefined) {
    return [undefined];
  } else if (type instanceof ZodNull) {
    return [null];
  } else if (type instanceof ZodOptional) {
    return [undefined, ...getDiscriminator(type.unwrap())];
  } else if (type instanceof ZodNullable) {
    return [null, ...getDiscriminator(type.unwrap())];
  } else if (type instanceof ZodBranded) {
    return getDiscriminator(type.unwrap());
  } else if (type instanceof ZodReadonly) {
    return getDiscriminator(type.unwrap());
  } else if (type instanceof ZodCatch) {
    return getDiscriminator(type._def.innerType);
  } else {
    return [];
  }
};

class ZodDiscriminatedUnion extends ZodType {
  _parse(input) {
    const { ctx } = this._processInputParams(input);
    if (ctx.parsedType !== ZodParsedType.object) {
      addIssueToContext(ctx, {
        code: ZodIssueCode.invalid_type,
        expected: ZodParsedType.object,
        received: ctx.parsedType
      });
      return INVALID;
    }
    const discriminator = this.discriminator;
    const discriminatorValue = ctx.data[discriminator];
    const option = this.optionsMap.get(discriminatorValue);
    if (!option) {
      addIssueToContext(ctx, {
        code: ZodIssueCode.invalid_union_discriminator,
        options: Array.from(this.optionsMap.keys()),
        path: [discriminator]
      });
      return INVALID;
    }
    if (ctx.common.async) {
      return option._parseAsync({
        data: ctx.data,
        path: ctx.path,
        parent: ctx
      });
    } else {
      return option._parseSync({
        data: ctx.data,
        path: ctx.path,
        parent: ctx
      });
    }
  }
  get discriminator() {
    return this._def.discriminator;
  }
  get options() {
    return this._def.options;
  }
  get optionsMap() {
    return this._def.optionsMap;
  }
  static create(discriminator, options, params) {
    const optionsMap = new Map;
    for (const type of options) {
      const discriminatorValues = getDiscriminator(type.shape[discriminator]);
      if (!discriminatorValues.length) {
        throw new Error(`A discriminator value for key \`${discriminator}\` could not be extracted from all schema options`);
      }
      for (const value2 of discriminatorValues) {
        if (optionsMap.has(value2)) {
          throw new Error(`Discriminator property ${String(discriminator)} has duplicate value ${String(value2)}`);
        }
        optionsMap.set(value2, type);
      }
    }
    return new ZodDiscriminatedUnion({
      typeName: ZodFirstPartyTypeKind.ZodDiscriminatedUnion,
      discriminator,
      options,
      optionsMap,
      ...processCreateParams(params)
    });
  }
}
function mergeValues(a, b) {
  const aType = getParsedType(a);
  const bType = getParsedType(b);
  if (a === b) {
    return { valid: true, data: a };
  } else if (aType === ZodParsedType.object && bType === ZodParsedType.object) {
    const bKeys = util.objectKeys(b);
    const sharedKeys = util.objectKeys(a).filter((key) => bKeys.indexOf(key) !== -1);
    const newObj = { ...a, ...b };
    for (const key of sharedKeys) {
      const sharedValue = mergeValues(a[key], b[key]);
      if (!sharedValue.valid) {
        return { valid: false };
      }
      newObj[key] = sharedValue.data;
    }
    return { valid: true, data: newObj };
  } else if (aType === ZodParsedType.array && bType === ZodParsedType.array) {
    if (a.length !== b.length) {
      return { valid: false };
    }
    const newArray = [];
    for (let index = 0;index < a.length; index++) {
      const itemA = a[index];
      const itemB = b[index];
      const sharedValue = mergeValues(itemA, itemB);
      if (!sharedValue.valid) {
        return { valid: false };
      }
      newArray.push(sharedValue.data);
    }
    return { valid: true, data: newArray };
  } else if (aType === ZodParsedType.date && bType === ZodParsedType.date && +a === +b) {
    return { valid: true, data: a };
  } else {
    return { valid: false };
  }
}

class ZodIntersection extends ZodType {
  _parse(input) {
    const { status, ctx } = this._processInputParams(input);
    const handleParsed = (parsedLeft, parsedRight) => {
      if (isAborted(parsedLeft) || isAborted(parsedRight)) {
        return INVALID;
      }
      const merged = mergeValues(parsedLeft.value, parsedRight.value);
      if (!merged.valid) {
        addIssueToContext(ctx, {
          code: ZodIssueCode.invalid_intersection_types
        });
        return INVALID;
      }
      if (isDirty(parsedLeft) || isDirty(parsedRight)) {
        status.dirty();
      }
      return { status: status.value, value: merged.data };
    };
    if (ctx.common.async) {
      return Promise.all([
        this._def.left._parseAsync({
          data: ctx.data,
          path: ctx.path,
          parent: ctx
        }),
        this._def.right._parseAsync({
          data: ctx.data,
          path: ctx.path,
          parent: ctx
        })
      ]).then(([left, right]) => handleParsed(left, right));
    } else {
      return handleParsed(this._def.left._parseSync({
        data: ctx.data,
        path: ctx.path,
        parent: ctx
      }), this._def.right._parseSync({
        data: ctx.data,
        path: ctx.path,
        parent: ctx
      }));
    }
  }
}
ZodIntersection.create = (left, right, params) => {
  return new ZodIntersection({
    left,
    right,
    typeName: ZodFirstPartyTypeKind.ZodIntersection,
    ...processCreateParams(params)
  });
};

class ZodTuple extends ZodType {
  _parse(input) {
    const { status, ctx } = this._processInputParams(input);
    if (ctx.parsedType !== ZodParsedType.array) {
      addIssueToContext(ctx, {
        code: ZodIssueCode.invalid_type,
        expected: ZodParsedType.array,
        received: ctx.parsedType
      });
      return INVALID;
    }
    if (ctx.data.length < this._def.items.length) {
      addIssueToContext(ctx, {
        code: ZodIssueCode.too_small,
        minimum: this._def.items.length,
        inclusive: true,
        exact: false,
        type: "array"
      });
      return INVALID;
    }
    const rest = this._def.rest;
    if (!rest && ctx.data.length > this._def.items.length) {
      addIssueToContext(ctx, {
        code: ZodIssueCode.too_big,
        maximum: this._def.items.length,
        inclusive: true,
        exact: false,
        type: "array"
      });
      status.dirty();
    }
    const items = [...ctx.data].map((item, itemIndex) => {
      const schema = this._def.items[itemIndex] || this._def.rest;
      if (!schema)
        return null;
      return schema._parse(new ParseInputLazyPath(ctx, item, ctx.path, itemIndex));
    }).filter((x) => !!x);
    if (ctx.common.async) {
      return Promise.all(items).then((results) => {
        return ParseStatus.mergeArray(status, results);
      });
    } else {
      return ParseStatus.mergeArray(status, items);
    }
  }
  get items() {
    return this._def.items;
  }
  rest(rest) {
    return new ZodTuple({
      ...this._def,
      rest
    });
  }
}
ZodTuple.create = (schemas, params) => {
  if (!Array.isArray(schemas)) {
    throw new Error("You must pass an array of schemas to z.tuple([ ... ])");
  }
  return new ZodTuple({
    items: schemas,
    typeName: ZodFirstPartyTypeKind.ZodTuple,
    rest: null,
    ...processCreateParams(params)
  });
};

class ZodRecord extends ZodType {
  get keySchema() {
    return this._def.keyType;
  }
  get valueSchema() {
    return this._def.valueType;
  }
  _parse(input) {
    const { status, ctx } = this._processInputParams(input);
    if (ctx.parsedType !== ZodParsedType.object) {
      addIssueToContext(ctx, {
        code: ZodIssueCode.invalid_type,
        expected: ZodParsedType.object,
        received: ctx.parsedType
      });
      return INVALID;
    }
    const pairs = [];
    const keyType = this._def.keyType;
    const valueType = this._def.valueType;
    for (const key in ctx.data) {
      pairs.push({
        key: keyType._parse(new ParseInputLazyPath(ctx, key, ctx.path, key)),
        value: valueType._parse(new ParseInputLazyPath(ctx, ctx.data[key], ctx.path, key)),
        alwaysSet: key in ctx.data
      });
    }
    if (ctx.common.async) {
      return ParseStatus.mergeObjectAsync(status, pairs);
    } else {
      return ParseStatus.mergeObjectSync(status, pairs);
    }
  }
  get element() {
    return this._def.valueType;
  }
  static create(first, second, third) {
    if (second instanceof ZodType) {
      return new ZodRecord({
        keyType: first,
        valueType: second,
        typeName: ZodFirstPartyTypeKind.ZodRecord,
        ...processCreateParams(third)
      });
    }
    return new ZodRecord({
      keyType: ZodString.create(),
      valueType: first,
      typeName: ZodFirstPartyTypeKind.ZodRecord,
      ...processCreateParams(second)
    });
  }
}

class ZodMap extends ZodType {
  get keySchema() {
    return this._def.keyType;
  }
  get valueSchema() {
    return this._def.valueType;
  }
  _parse(input) {
    const { status, ctx } = this._processInputParams(input);
    if (ctx.parsedType !== ZodParsedType.map) {
      addIssueToContext(ctx, {
        code: ZodIssueCode.invalid_type,
        expected: ZodParsedType.map,
        received: ctx.parsedType
      });
      return INVALID;
    }
    const keyType = this._def.keyType;
    const valueType = this._def.valueType;
    const pairs = [...ctx.data.entries()].map(([key, value2], index) => {
      return {
        key: keyType._parse(new ParseInputLazyPath(ctx, key, ctx.path, [index, "key"])),
        value: valueType._parse(new ParseInputLazyPath(ctx, value2, ctx.path, [index, "value"]))
      };
    });
    if (ctx.common.async) {
      const finalMap = new Map;
      return Promise.resolve().then(async () => {
        for (const pair of pairs) {
          const key = await pair.key;
          const value2 = await pair.value;
          if (key.status === "aborted" || value2.status === "aborted") {
            return INVALID;
          }
          if (key.status === "dirty" || value2.status === "dirty") {
            status.dirty();
          }
          finalMap.set(key.value, value2.value);
        }
        return { status: status.value, value: finalMap };
      });
    } else {
      const finalMap = new Map;
      for (const pair of pairs) {
        const key = pair.key;
        const value2 = pair.value;
        if (key.status === "aborted" || value2.status === "aborted") {
          return INVALID;
        }
        if (key.status === "dirty" || value2.status === "dirty") {
          status.dirty();
        }
        finalMap.set(key.value, value2.value);
      }
      return { status: status.value, value: finalMap };
    }
  }
}
ZodMap.create = (keyType, valueType, params) => {
  return new ZodMap({
    valueType,
    keyType,
    typeName: ZodFirstPartyTypeKind.ZodMap,
    ...processCreateParams(params)
  });
};

class ZodSet extends ZodType {
  _parse(input) {
    const { status, ctx } = this._processInputParams(input);
    if (ctx.parsedType !== ZodParsedType.set) {
      addIssueToContext(ctx, {
        code: ZodIssueCode.invalid_type,
        expected: ZodParsedType.set,
        received: ctx.parsedType
      });
      return INVALID;
    }
    const def = this._def;
    if (def.minSize !== null) {
      if (ctx.data.size < def.minSize.value) {
        addIssueToContext(ctx, {
          code: ZodIssueCode.too_small,
          minimum: def.minSize.value,
          type: "set",
          inclusive: true,
          exact: false,
          message: def.minSize.message
        });
        status.dirty();
      }
    }
    if (def.maxSize !== null) {
      if (ctx.data.size > def.maxSize.value) {
        addIssueToContext(ctx, {
          code: ZodIssueCode.too_big,
          maximum: def.maxSize.value,
          type: "set",
          inclusive: true,
          exact: false,
          message: def.maxSize.message
        });
        status.dirty();
      }
    }
    const valueType = this._def.valueType;
    function finalizeSet(elements2) {
      const parsedSet = new Set;
      for (const element of elements2) {
        if (element.status === "aborted")
          return INVALID;
        if (element.status === "dirty")
          status.dirty();
        parsedSet.add(element.value);
      }
      return { status: status.value, value: parsedSet };
    }
    const elements = [...ctx.data.values()].map((item, i2) => valueType._parse(new ParseInputLazyPath(ctx, item, ctx.path, i2)));
    if (ctx.common.async) {
      return Promise.all(elements).then((elements2) => finalizeSet(elements2));
    } else {
      return finalizeSet(elements);
    }
  }
  min(minSize, message) {
    return new ZodSet({
      ...this._def,
      minSize: { value: minSize, message: errorUtil.toString(message) }
    });
  }
  max(maxSize, message) {
    return new ZodSet({
      ...this._def,
      maxSize: { value: maxSize, message: errorUtil.toString(message) }
    });
  }
  size(size, message) {
    return this.min(size, message).max(size, message);
  }
  nonempty(message) {
    return this.min(1, message);
  }
}
ZodSet.create = (valueType, params) => {
  return new ZodSet({
    valueType,
    minSize: null,
    maxSize: null,
    typeName: ZodFirstPartyTypeKind.ZodSet,
    ...processCreateParams(params)
  });
};

class ZodFunction extends ZodType {
  constructor() {
    super(...arguments);
    this.validate = this.implement;
  }
  _parse(input) {
    const { ctx } = this._processInputParams(input);
    if (ctx.parsedType !== ZodParsedType.function) {
      addIssueToContext(ctx, {
        code: ZodIssueCode.invalid_type,
        expected: ZodParsedType.function,
        received: ctx.parsedType
      });
      return INVALID;
    }
    function makeArgsIssue(args, error) {
      return makeIssue({
        data: args,
        path: ctx.path,
        errorMaps: [ctx.common.contextualErrorMap, ctx.schemaErrorMap, getErrorMap(), en_default].filter((x) => !!x),
        issueData: {
          code: ZodIssueCode.invalid_arguments,
          argumentsError: error
        }
      });
    }
    function makeReturnsIssue(returns, error) {
      return makeIssue({
        data: returns,
        path: ctx.path,
        errorMaps: [ctx.common.contextualErrorMap, ctx.schemaErrorMap, getErrorMap(), en_default].filter((x) => !!x),
        issueData: {
          code: ZodIssueCode.invalid_return_type,
          returnTypeError: error
        }
      });
    }
    const params = { errorMap: ctx.common.contextualErrorMap };
    const fn = ctx.data;
    if (this._def.returns instanceof ZodPromise) {
      const me = this;
      return OK(async function(...args) {
        const error = new ZodError([]);
        const parsedArgs = await me._def.args.parseAsync(args, params).catch((e) => {
          error.addIssue(makeArgsIssue(args, e));
          throw error;
        });
        const result = await Reflect.apply(fn, this, parsedArgs);
        const parsedReturns = await me._def.returns._def.type.parseAsync(result, params).catch((e) => {
          error.addIssue(makeReturnsIssue(result, e));
          throw error;
        });
        return parsedReturns;
      });
    } else {
      const me = this;
      return OK(function(...args) {
        const parsedArgs = me._def.args.safeParse(args, params);
        if (!parsedArgs.success) {
          throw new ZodError([makeArgsIssue(args, parsedArgs.error)]);
        }
        const result = Reflect.apply(fn, this, parsedArgs.data);
        const parsedReturns = me._def.returns.safeParse(result, params);
        if (!parsedReturns.success) {
          throw new ZodError([makeReturnsIssue(result, parsedReturns.error)]);
        }
        return parsedReturns.data;
      });
    }
  }
  parameters() {
    return this._def.args;
  }
  returnType() {
    return this._def.returns;
  }
  args(...items) {
    return new ZodFunction({
      ...this._def,
      args: ZodTuple.create(items).rest(ZodUnknown.create())
    });
  }
  returns(returnType) {
    return new ZodFunction({
      ...this._def,
      returns: returnType
    });
  }
  implement(func) {
    const validatedFunc = this.parse(func);
    return validatedFunc;
  }
  strictImplement(func) {
    const validatedFunc = this.parse(func);
    return validatedFunc;
  }
  static create(args, returns, params) {
    return new ZodFunction({
      args: args ? args : ZodTuple.create([]).rest(ZodUnknown.create()),
      returns: returns || ZodUnknown.create(),
      typeName: ZodFirstPartyTypeKind.ZodFunction,
      ...processCreateParams(params)
    });
  }
}

class ZodLazy extends ZodType {
  get schema() {
    return this._def.getter();
  }
  _parse(input) {
    const { ctx } = this._processInputParams(input);
    const lazySchema = this._def.getter();
    return lazySchema._parse({ data: ctx.data, path: ctx.path, parent: ctx });
  }
}
ZodLazy.create = (getter, params) => {
  return new ZodLazy({
    getter,
    typeName: ZodFirstPartyTypeKind.ZodLazy,
    ...processCreateParams(params)
  });
};

class ZodLiteral extends ZodType {
  _parse(input) {
    if (input.data !== this._def.value) {
      const ctx = this._getOrReturnCtx(input);
      addIssueToContext(ctx, {
        received: ctx.data,
        code: ZodIssueCode.invalid_literal,
        expected: this._def.value
      });
      return INVALID;
    }
    return { status: "valid", value: input.data };
  }
  get value() {
    return this._def.value;
  }
}
ZodLiteral.create = (value2, params) => {
  return new ZodLiteral({
    value: value2,
    typeName: ZodFirstPartyTypeKind.ZodLiteral,
    ...processCreateParams(params)
  });
};
function createZodEnum(values, params) {
  return new ZodEnum({
    values,
    typeName: ZodFirstPartyTypeKind.ZodEnum,
    ...processCreateParams(params)
  });
}

class ZodEnum extends ZodType {
  _parse(input) {
    if (typeof input.data !== "string") {
      const ctx = this._getOrReturnCtx(input);
      const expectedValues = this._def.values;
      addIssueToContext(ctx, {
        expected: util.joinValues(expectedValues),
        received: ctx.parsedType,
        code: ZodIssueCode.invalid_type
      });
      return INVALID;
    }
    if (!this._cache) {
      this._cache = new Set(this._def.values);
    }
    if (!this._cache.has(input.data)) {
      const ctx = this._getOrReturnCtx(input);
      const expectedValues = this._def.values;
      addIssueToContext(ctx, {
        received: ctx.data,
        code: ZodIssueCode.invalid_enum_value,
        options: expectedValues
      });
      return INVALID;
    }
    return OK(input.data);
  }
  get options() {
    return this._def.values;
  }
  get enum() {
    const enumValues = {};
    for (const val of this._def.values) {
      enumValues[val] = val;
    }
    return enumValues;
  }
  get Values() {
    const enumValues = {};
    for (const val of this._def.values) {
      enumValues[val] = val;
    }
    return enumValues;
  }
  get Enum() {
    const enumValues = {};
    for (const val of this._def.values) {
      enumValues[val] = val;
    }
    return enumValues;
  }
  extract(values, newDef = this._def) {
    return ZodEnum.create(values, {
      ...this._def,
      ...newDef
    });
  }
  exclude(values, newDef = this._def) {
    return ZodEnum.create(this.options.filter((opt) => !values.includes(opt)), {
      ...this._def,
      ...newDef
    });
  }
}
ZodEnum.create = createZodEnum;

class ZodNativeEnum extends ZodType {
  _parse(input) {
    const nativeEnumValues = util.getValidEnumValues(this._def.values);
    const ctx = this._getOrReturnCtx(input);
    if (ctx.parsedType !== ZodParsedType.string && ctx.parsedType !== ZodParsedType.number) {
      const expectedValues = util.objectValues(nativeEnumValues);
      addIssueToContext(ctx, {
        expected: util.joinValues(expectedValues),
        received: ctx.parsedType,
        code: ZodIssueCode.invalid_type
      });
      return INVALID;
    }
    if (!this._cache) {
      this._cache = new Set(util.getValidEnumValues(this._def.values));
    }
    if (!this._cache.has(input.data)) {
      const expectedValues = util.objectValues(nativeEnumValues);
      addIssueToContext(ctx, {
        received: ctx.data,
        code: ZodIssueCode.invalid_enum_value,
        options: expectedValues
      });
      return INVALID;
    }
    return OK(input.data);
  }
  get enum() {
    return this._def.values;
  }
}
ZodNativeEnum.create = (values, params) => {
  return new ZodNativeEnum({
    values,
    typeName: ZodFirstPartyTypeKind.ZodNativeEnum,
    ...processCreateParams(params)
  });
};

class ZodPromise extends ZodType {
  unwrap() {
    return this._def.type;
  }
  _parse(input) {
    const { ctx } = this._processInputParams(input);
    if (ctx.parsedType !== ZodParsedType.promise && ctx.common.async === false) {
      addIssueToContext(ctx, {
        code: ZodIssueCode.invalid_type,
        expected: ZodParsedType.promise,
        received: ctx.parsedType
      });
      return INVALID;
    }
    const promisified = ctx.parsedType === ZodParsedType.promise ? ctx.data : Promise.resolve(ctx.data);
    return OK(promisified.then((data) => {
      return this._def.type.parseAsync(data, {
        path: ctx.path,
        errorMap: ctx.common.contextualErrorMap
      });
    }));
  }
}
ZodPromise.create = (schema, params) => {
  return new ZodPromise({
    type: schema,
    typeName: ZodFirstPartyTypeKind.ZodPromise,
    ...processCreateParams(params)
  });
};

class ZodEffects extends ZodType {
  innerType() {
    return this._def.schema;
  }
  sourceType() {
    return this._def.schema._def.typeName === ZodFirstPartyTypeKind.ZodEffects ? this._def.schema.sourceType() : this._def.schema;
  }
  _parse(input) {
    const { status, ctx } = this._processInputParams(input);
    const effect = this._def.effect || null;
    const checkCtx = {
      addIssue: (arg) => {
        addIssueToContext(ctx, arg);
        if (arg.fatal) {
          status.abort();
        } else {
          status.dirty();
        }
      },
      get path() {
        return ctx.path;
      }
    };
    checkCtx.addIssue = checkCtx.addIssue.bind(checkCtx);
    if (effect.type === "preprocess") {
      const processed = effect.transform(ctx.data, checkCtx);
      if (ctx.common.async) {
        return Promise.resolve(processed).then(async (processed2) => {
          if (status.value === "aborted")
            return INVALID;
          const result = await this._def.schema._parseAsync({
            data: processed2,
            path: ctx.path,
            parent: ctx
          });
          if (result.status === "aborted")
            return INVALID;
          if (result.status === "dirty")
            return DIRTY(result.value);
          if (status.value === "dirty")
            return DIRTY(result.value);
          return result;
        });
      } else {
        if (status.value === "aborted")
          return INVALID;
        const result = this._def.schema._parseSync({
          data: processed,
          path: ctx.path,
          parent: ctx
        });
        if (result.status === "aborted")
          return INVALID;
        if (result.status === "dirty")
          return DIRTY(result.value);
        if (status.value === "dirty")
          return DIRTY(result.value);
        return result;
      }
    }
    if (effect.type === "refinement") {
      const executeRefinement = (acc) => {
        const result = effect.refinement(acc, checkCtx);
        if (ctx.common.async) {
          return Promise.resolve(result);
        }
        if (result instanceof Promise) {
          throw new Error("Async refinement encountered during synchronous parse operation. Use .parseAsync instead.");
        }
        return acc;
      };
      if (ctx.common.async === false) {
        const inner = this._def.schema._parseSync({
          data: ctx.data,
          path: ctx.path,
          parent: ctx
        });
        if (inner.status === "aborted")
          return INVALID;
        if (inner.status === "dirty")
          status.dirty();
        executeRefinement(inner.value);
        return { status: status.value, value: inner.value };
      } else {
        return this._def.schema._parseAsync({ data: ctx.data, path: ctx.path, parent: ctx }).then((inner) => {
          if (inner.status === "aborted")
            return INVALID;
          if (inner.status === "dirty")
            status.dirty();
          return executeRefinement(inner.value).then(() => {
            return { status: status.value, value: inner.value };
          });
        });
      }
    }
    if (effect.type === "transform") {
      if (ctx.common.async === false) {
        const base = this._def.schema._parseSync({
          data: ctx.data,
          path: ctx.path,
          parent: ctx
        });
        if (!isValid(base))
          return INVALID;
        const result = effect.transform(base.value, checkCtx);
        if (result instanceof Promise) {
          throw new Error(`Asynchronous transform encountered during synchronous parse operation. Use .parseAsync instead.`);
        }
        return { status: status.value, value: result };
      } else {
        return this._def.schema._parseAsync({ data: ctx.data, path: ctx.path, parent: ctx }).then((base) => {
          if (!isValid(base))
            return INVALID;
          return Promise.resolve(effect.transform(base.value, checkCtx)).then((result) => ({
            status: status.value,
            value: result
          }));
        });
      }
    }
    util.assertNever(effect);
  }
}
ZodEffects.create = (schema, effect, params) => {
  return new ZodEffects({
    schema,
    typeName: ZodFirstPartyTypeKind.ZodEffects,
    effect,
    ...processCreateParams(params)
  });
};
ZodEffects.createWithPreprocess = (preprocess, schema, params) => {
  return new ZodEffects({
    schema,
    effect: { type: "preprocess", transform: preprocess },
    typeName: ZodFirstPartyTypeKind.ZodEffects,
    ...processCreateParams(params)
  });
};

class ZodOptional extends ZodType {
  _parse(input) {
    const parsedType = this._getType(input);
    if (parsedType === ZodParsedType.undefined) {
      return OK(undefined);
    }
    return this._def.innerType._parse(input);
  }
  unwrap() {
    return this._def.innerType;
  }
}
ZodOptional.create = (type, params) => {
  return new ZodOptional({
    innerType: type,
    typeName: ZodFirstPartyTypeKind.ZodOptional,
    ...processCreateParams(params)
  });
};

class ZodNullable extends ZodType {
  _parse(input) {
    const parsedType = this._getType(input);
    if (parsedType === ZodParsedType.null) {
      return OK(null);
    }
    return this._def.innerType._parse(input);
  }
  unwrap() {
    return this._def.innerType;
  }
}
ZodNullable.create = (type, params) => {
  return new ZodNullable({
    innerType: type,
    typeName: ZodFirstPartyTypeKind.ZodNullable,
    ...processCreateParams(params)
  });
};

class ZodDefault extends ZodType {
  _parse(input) {
    const { ctx } = this._processInputParams(input);
    let data = ctx.data;
    if (ctx.parsedType === ZodParsedType.undefined) {
      data = this._def.defaultValue();
    }
    return this._def.innerType._parse({
      data,
      path: ctx.path,
      parent: ctx
    });
  }
  removeDefault() {
    return this._def.innerType;
  }
}
ZodDefault.create = (type, params) => {
  return new ZodDefault({
    innerType: type,
    typeName: ZodFirstPartyTypeKind.ZodDefault,
    defaultValue: typeof params.default === "function" ? params.default : () => params.default,
    ...processCreateParams(params)
  });
};

class ZodCatch extends ZodType {
  _parse(input) {
    const { ctx } = this._processInputParams(input);
    const newCtx = {
      ...ctx,
      common: {
        ...ctx.common,
        issues: []
      }
    };
    const result = this._def.innerType._parse({
      data: newCtx.data,
      path: newCtx.path,
      parent: {
        ...newCtx
      }
    });
    if (isAsync(result)) {
      return result.then((result2) => {
        return {
          status: "valid",
          value: result2.status === "valid" ? result2.value : this._def.catchValue({
            get error() {
              return new ZodError(newCtx.common.issues);
            },
            input: newCtx.data
          })
        };
      });
    } else {
      return {
        status: "valid",
        value: result.status === "valid" ? result.value : this._def.catchValue({
          get error() {
            return new ZodError(newCtx.common.issues);
          },
          input: newCtx.data
        })
      };
    }
  }
  removeCatch() {
    return this._def.innerType;
  }
}
ZodCatch.create = (type, params) => {
  return new ZodCatch({
    innerType: type,
    typeName: ZodFirstPartyTypeKind.ZodCatch,
    catchValue: typeof params.catch === "function" ? params.catch : () => params.catch,
    ...processCreateParams(params)
  });
};

class ZodNaN extends ZodType {
  _parse(input) {
    const parsedType = this._getType(input);
    if (parsedType !== ZodParsedType.nan) {
      const ctx = this._getOrReturnCtx(input);
      addIssueToContext(ctx, {
        code: ZodIssueCode.invalid_type,
        expected: ZodParsedType.nan,
        received: ctx.parsedType
      });
      return INVALID;
    }
    return { status: "valid", value: input.data };
  }
}
ZodNaN.create = (params) => {
  return new ZodNaN({
    typeName: ZodFirstPartyTypeKind.ZodNaN,
    ...processCreateParams(params)
  });
};
var BRAND = Symbol("zod_brand");

class ZodBranded extends ZodType {
  _parse(input) {
    const { ctx } = this._processInputParams(input);
    const data = ctx.data;
    return this._def.type._parse({
      data,
      path: ctx.path,
      parent: ctx
    });
  }
  unwrap() {
    return this._def.type;
  }
}

class ZodPipeline extends ZodType {
  _parse(input) {
    const { status, ctx } = this._processInputParams(input);
    if (ctx.common.async) {
      const handleAsync = async () => {
        const inResult = await this._def.in._parseAsync({
          data: ctx.data,
          path: ctx.path,
          parent: ctx
        });
        if (inResult.status === "aborted")
          return INVALID;
        if (inResult.status === "dirty") {
          status.dirty();
          return DIRTY(inResult.value);
        } else {
          return this._def.out._parseAsync({
            data: inResult.value,
            path: ctx.path,
            parent: ctx
          });
        }
      };
      return handleAsync();
    } else {
      const inResult = this._def.in._parseSync({
        data: ctx.data,
        path: ctx.path,
        parent: ctx
      });
      if (inResult.status === "aborted")
        return INVALID;
      if (inResult.status === "dirty") {
        status.dirty();
        return {
          status: "dirty",
          value: inResult.value
        };
      } else {
        return this._def.out._parseSync({
          data: inResult.value,
          path: ctx.path,
          parent: ctx
        });
      }
    }
  }
  static create(a, b) {
    return new ZodPipeline({
      in: a,
      out: b,
      typeName: ZodFirstPartyTypeKind.ZodPipeline
    });
  }
}

class ZodReadonly extends ZodType {
  _parse(input) {
    const result = this._def.innerType._parse(input);
    const freeze = (data) => {
      if (isValid(data)) {
        data.value = Object.freeze(data.value);
      }
      return data;
    };
    return isAsync(result) ? result.then((data) => freeze(data)) : freeze(result);
  }
  unwrap() {
    return this._def.innerType;
  }
}
ZodReadonly.create = (type, params) => {
  return new ZodReadonly({
    innerType: type,
    typeName: ZodFirstPartyTypeKind.ZodReadonly,
    ...processCreateParams(params)
  });
};
function cleanParams(params, data) {
  const p = typeof params === "function" ? params(data) : typeof params === "string" ? { message: params } : params;
  const p2 = typeof p === "string" ? { message: p } : p;
  return p2;
}
function custom(check, _params = {}, fatal) {
  if (check)
    return ZodAny.create().superRefine((data, ctx) => {
      const r = check(data);
      if (r instanceof Promise) {
        return r.then((r2) => {
          if (!r2) {
            const params = cleanParams(_params, data);
            const _fatal = params.fatal ?? fatal ?? true;
            ctx.addIssue({ code: "custom", ...params, fatal: _fatal });
          }
        });
      }
      if (!r) {
        const params = cleanParams(_params, data);
        const _fatal = params.fatal ?? fatal ?? true;
        ctx.addIssue({ code: "custom", ...params, fatal: _fatal });
      }
      return;
    });
  return ZodAny.create();
}
var late = {
  object: ZodObject.lazycreate
};
var ZodFirstPartyTypeKind;
(function(ZodFirstPartyTypeKind2) {
  ZodFirstPartyTypeKind2["ZodString"] = "ZodString";
  ZodFirstPartyTypeKind2["ZodNumber"] = "ZodNumber";
  ZodFirstPartyTypeKind2["ZodNaN"] = "ZodNaN";
  ZodFirstPartyTypeKind2["ZodBigInt"] = "ZodBigInt";
  ZodFirstPartyTypeKind2["ZodBoolean"] = "ZodBoolean";
  ZodFirstPartyTypeKind2["ZodDate"] = "ZodDate";
  ZodFirstPartyTypeKind2["ZodSymbol"] = "ZodSymbol";
  ZodFirstPartyTypeKind2["ZodUndefined"] = "ZodUndefined";
  ZodFirstPartyTypeKind2["ZodNull"] = "ZodNull";
  ZodFirstPartyTypeKind2["ZodAny"] = "ZodAny";
  ZodFirstPartyTypeKind2["ZodUnknown"] = "ZodUnknown";
  ZodFirstPartyTypeKind2["ZodNever"] = "ZodNever";
  ZodFirstPartyTypeKind2["ZodVoid"] = "ZodVoid";
  ZodFirstPartyTypeKind2["ZodArray"] = "ZodArray";
  ZodFirstPartyTypeKind2["ZodObject"] = "ZodObject";
  ZodFirstPartyTypeKind2["ZodUnion"] = "ZodUnion";
  ZodFirstPartyTypeKind2["ZodDiscriminatedUnion"] = "ZodDiscriminatedUnion";
  ZodFirstPartyTypeKind2["ZodIntersection"] = "ZodIntersection";
  ZodFirstPartyTypeKind2["ZodTuple"] = "ZodTuple";
  ZodFirstPartyTypeKind2["ZodRecord"] = "ZodRecord";
  ZodFirstPartyTypeKind2["ZodMap"] = "ZodMap";
  ZodFirstPartyTypeKind2["ZodSet"] = "ZodSet";
  ZodFirstPartyTypeKind2["ZodFunction"] = "ZodFunction";
  ZodFirstPartyTypeKind2["ZodLazy"] = "ZodLazy";
  ZodFirstPartyTypeKind2["ZodLiteral"] = "ZodLiteral";
  ZodFirstPartyTypeKind2["ZodEnum"] = "ZodEnum";
  ZodFirstPartyTypeKind2["ZodEffects"] = "ZodEffects";
  ZodFirstPartyTypeKind2["ZodNativeEnum"] = "ZodNativeEnum";
  ZodFirstPartyTypeKind2["ZodOptional"] = "ZodOptional";
  ZodFirstPartyTypeKind2["ZodNullable"] = "ZodNullable";
  ZodFirstPartyTypeKind2["ZodDefault"] = "ZodDefault";
  ZodFirstPartyTypeKind2["ZodCatch"] = "ZodCatch";
  ZodFirstPartyTypeKind2["ZodPromise"] = "ZodPromise";
  ZodFirstPartyTypeKind2["ZodBranded"] = "ZodBranded";
  ZodFirstPartyTypeKind2["ZodPipeline"] = "ZodPipeline";
  ZodFirstPartyTypeKind2["ZodReadonly"] = "ZodReadonly";
})(ZodFirstPartyTypeKind || (ZodFirstPartyTypeKind = {}));
var instanceOfType = (cls, params = {
  message: `Input not instance of ${cls.name}`
}) => custom((data) => data instanceof cls, params);
var stringType = ZodString.create;
var numberType = ZodNumber.create;
var nanType = ZodNaN.create;
var bigIntType = ZodBigInt.create;
var booleanType = ZodBoolean.create;
var dateType = ZodDate.create;
var symbolType = ZodSymbol.create;
var undefinedType = ZodUndefined.create;
var nullType = ZodNull.create;
var anyType = ZodAny.create;
var unknownType = ZodUnknown.create;
var neverType = ZodNever.create;
var voidType = ZodVoid.create;
var arrayType = ZodArray.create;
var objectType = ZodObject.create;
var strictObjectType = ZodObject.strictCreate;
var unionType = ZodUnion.create;
var discriminatedUnionType = ZodDiscriminatedUnion.create;
var intersectionType = ZodIntersection.create;
var tupleType = ZodTuple.create;
var recordType = ZodRecord.create;
var mapType = ZodMap.create;
var setType = ZodSet.create;
var functionType = ZodFunction.create;
var lazyType = ZodLazy.create;
var literalType = ZodLiteral.create;
var enumType = ZodEnum.create;
var nativeEnumType = ZodNativeEnum.create;
var promiseType = ZodPromise.create;
var effectsType = ZodEffects.create;
var optionalType = ZodOptional.create;
var nullableType = ZodNullable.create;
var preprocessType = ZodEffects.createWithPreprocess;
var pipelineType = ZodPipeline.create;
var ostring = () => stringType().optional();
var onumber = () => numberType().optional();
var oboolean = () => booleanType().optional();
var coerce = {
  string: (arg) => ZodString.create({ ...arg, coerce: true }),
  number: (arg) => ZodNumber.create({ ...arg, coerce: true }),
  boolean: (arg) => ZodBoolean.create({
    ...arg,
    coerce: true
  }),
  bigint: (arg) => ZodBigInt.create({ ...arg, coerce: true }),
  date: (arg) => ZodDate.create({ ...arg, coerce: true })
};
var NEVER = INVALID;
var globalHostBindingsSchema = exports_external.object({
  switchModes: exports_external.function().args(exports_external.nativeEnum(Mode)).returns(exports_external.void()),
  log: exports_external.function().args(exports_external.string()).returns(exports_external.void()),
  sendResponse: exports_external.function().args(exports_external.union([exports_external.instanceof(Uint8Array), exports_external.custom()])).returns(exports_external.number()),
  versionV2: exports_external.function().args().returns(exports_external.void()),
  callCapability: exports_external.function().args(exports_external.union([exports_external.instanceof(Uint8Array), exports_external.custom()])).returns(exports_external.number()),
  awaitCapabilities: exports_external.function().args(exports_external.union([exports_external.instanceof(Uint8Array), exports_external.custom()]), exports_external.number()).returns(exports_external.union([exports_external.instanceof(Uint8Array), exports_external.custom()])),
  getSecrets: exports_external.function().args(exports_external.union([exports_external.instanceof(Uint8Array), exports_external.custom()]), exports_external.number()).returns(exports_external.any()),
  awaitSecrets: exports_external.function().args(exports_external.union([exports_external.instanceof(Uint8Array), exports_external.custom()]), exports_external.number()).returns(exports_external.union([exports_external.instanceof(Uint8Array), exports_external.custom()])),
  getWasiArgs: exports_external.function().args().returns(exports_external.string()),
  now: exports_external.function().args().returns(exports_external.number())
});
var validateGlobalHostBindings = () => {
  const globalFunctions = globalThis;
  try {
    return globalHostBindingsSchema.parse(globalFunctions);
  } catch (error) {
    const missingFunctions = Object.keys(globalHostBindingsSchema.shape).filter((key) => !(key in globalFunctions));
    throw new Error(`Missing required global host functions: ${missingFunctions.join(", ")}. ` + `This indicates the runtime environment is not properly configured.`);
  }
};
var _hostBindings = null;
var hostBindings = new Proxy({}, {
  get(target, prop) {
    if (!_hostBindings) {
      _hostBindings = validateGlobalHostBindings();
    }
    return _hostBindings[prop];
  }
});

class ConsensusCapability {
  static CAPABILITY_ID = "consensus@1.0.0-alpha";
  static CAPABILITY_NAME = "consensus";
  static CAPABILITY_VERSION = "1.0.0-alpha";
  simple(runtime, input) {
    let payload;
    if (input.$typeName) {
      payload = input;
    } else {
      payload = fromJson(SimpleConsensusInputsSchema, input);
    }
    const capabilityId = ConsensusCapability.CAPABILITY_ID;
    const capabilityResponse = runtime.callCapability({
      capabilityId,
      method: "Simple",
      payload,
      inputSchema: SimpleConsensusInputsSchema,
      outputSchema: ValueSchema2
    });
    return {
      result: () => {
        const result = capabilityResponse.result();
        return result;
      }
    };
  }
  report(runtime, input) {
    let payload;
    if (input.$typeName) {
      payload = input;
    } else {
      payload = fromJson(ReportRequestSchema, input);
    }
    const capabilityId = ConsensusCapability.CAPABILITY_ID;
    const capabilityResponse = runtime.callCapability({
      capabilityId,
      method: "Report",
      payload,
      inputSchema: ReportRequestSchema,
      outputSchema: ReportResponseSchema
    });
    return {
      result: () => {
        const result = capabilityResponse.result();
        return new Report(result);
      }
    };
  }
}

class CapabilityError extends Error {
  name;
  capabilityId;
  method;
  callbackId;
  constructor(message, options) {
    super(message);
    this.name = "CapabilityError";
    if (options) {
      this.capabilityId = options.capabilityId;
      this.method = options.method;
      this.callbackId = options.callbackId;
    }
  }
}

class DonModeError extends Error {
  constructor() {
    super("cannot use Runtime inside RunInNodeMode");
  }
}

class NodeModeError extends Error {
  constructor() {
    super("cannot use NodeRuntime outside RunInNodeMode");
    this.name = "NodeModeError";
  }
}

class SecretsError extends Error {
  secretRequest;
  error;
  constructor(secretRequest, error) {
    super(`error fetching ${secretRequest}: ${error}`);
    this.secretRequest = secretRequest;
    this.error = error;
  }
}

class BaseRuntimeImpl {
  config;
  nextCallId;
  helpers;
  maxResponseSize;
  mode;
  modeError;
  constructor(config, nextCallId, helpers, maxResponseSize, mode) {
    this.config = config;
    this.nextCallId = nextCallId;
    this.helpers = helpers;
    this.maxResponseSize = maxResponseSize;
    this.mode = mode;
  }
  callCapability({ capabilityId, method, payload, inputSchema, outputSchema }) {
    if (this.modeError) {
      return {
        result: () => {
          throw this.modeError;
        }
      };
    }
    const anyPayload = anyPack(inputSchema, payload);
    const callbackId = this.nextCallId;
    if (this.mode === Mode.DON) {
      this.nextCallId++;
    } else {
      this.nextCallId--;
    }
    const req = create(CapabilityRequestSchema, {
      id: capabilityId,
      method,
      payload: anyPayload,
      callbackId
    });
    if (!this.helpers.call(req)) {
      return {
        result: () => {
          throw new CapabilityError(`Capability not found ${capabilityId}`, {
            callbackId,
            method,
            capabilityId
          });
        }
      };
    }
    return {
      result: () => {
        const awaitRequest = create(AwaitCapabilitiesRequestSchema, {
          ids: [callbackId]
        });
        const awaitResponse = this.helpers.await(awaitRequest, this.maxResponseSize);
        const capabilityResponse = awaitResponse.responses[callbackId];
        if (!capabilityResponse) {
          throw new CapabilityError(`No response found for callback ID ${callbackId}`, {
            capabilityId,
            method,
            callbackId
          });
        }
        const response = capabilityResponse.response;
        switch (response.case) {
          case "payload": {
            try {
              return anyUnpack(response.value, outputSchema);
            } catch {
              throw new CapabilityError(`Error cannot unwrap payload`, {
                capabilityId,
                method,
                callbackId
              });
            }
          }
          case "error":
            throw new CapabilityError(`Error ${response.value}`, {
              capabilityId,
              method,
              callbackId
            });
          default:
            throw new CapabilityError(`Error cannot unwrap ${response.case}`, {
              capabilityId,
              method,
              callbackId
            });
        }
      }
    };
  }
  getNextCallId() {
    return this.nextCallId;
  }
  now() {
    return new Date(this.helpers.now() / 1e6);
  }
  log(message) {
    this.helpers.log(message);
  }
}

class NodeRuntimeImpl extends BaseRuntimeImpl {
  _isNodeRuntime = true;
  constructor(config, nextCallId, helpers, maxResponseSize) {
    helpers.switchModes(Mode.NODE);
    super(config, nextCallId, helpers, maxResponseSize, Mode.NODE);
  }
}

class RuntimeImpl extends BaseRuntimeImpl {
  nextNodeCallId = -1;
  constructor(config, nextCallId, helpers, maxResponseSize) {
    helpers.switchModes(Mode.DON);
    super(config, nextCallId, helpers, maxResponseSize, Mode.DON);
  }
  runInNodeMode(fn, consensusAggretation, unwrapOptions) {
    return (...args) => {
      this.modeError = new DonModeError;
      const nodeRuntime = new NodeRuntimeImpl(this.config, this.nextNodeCallId, this.helpers, this.maxResponseSize);
      const consensusInput = create(SimpleConsensusInputsSchema, {
        descriptors: consensusAggretation.descriptor
      });
      if (consensusAggretation.defaultValue) {
        consensusInput.default = Value.from(consensusAggretation.defaultValue).proto();
      }
      try {
        const observation = fn(nodeRuntime, ...args);
        consensusInput.observation = {
          case: "value",
          value: Value.from(observation).proto()
        };
      } catch (e) {
        consensusInput.observation = {
          case: "error",
          value: e instanceof Error && e.message || String(e)
        };
      } finally {
        this.modeError = undefined;
        this.nextNodeCallId = nodeRuntime.nextCallId;
        nodeRuntime.modeError = new NodeModeError;
        this.helpers.switchModes(Mode.DON);
      }
      const consensus = new ConsensusCapability;
      const call = consensus.simple(this, consensusInput);
      return {
        result: () => {
          const result = call.result();
          const wrappedValue = Value.wrap(result);
          return unwrapOptions ? wrappedValue.unwrapToType(unwrapOptions) : wrappedValue.unwrap();
        }
      };
    };
  }
  getSecret(request) {
    if (this.modeError) {
      return {
        result: () => {
          throw this.modeError;
        }
      };
    }
    const secretRequest = request.$typeName ? request : create(SecretRequestSchema, request);
    const id = this.nextCallId;
    this.nextCallId++;
    const secretsReq = create(GetSecretsRequestSchema, {
      callbackId: id,
      requests: [secretRequest]
    });
    if (!this.helpers.getSecrets(secretsReq, this.maxResponseSize)) {
      return {
        result: () => {
          throw new SecretsError(secretRequest, "host is not making the secrets request");
        }
      };
    }
    return {
      result: () => {
        const awaitRequest = create(AwaitSecretsRequestSchema, { ids: [id] });
        const awaitResponse = this.helpers.awaitSecrets(awaitRequest, this.maxResponseSize);
        const secretsResponse = awaitResponse.responses[id];
        if (!secretsResponse) {
          throw new SecretsError(secretRequest, "no response");
        }
        const responses = secretsResponse.responses;
        if (responses.length !== 1) {
          throw new SecretsError(secretRequest, "invalid value returned from host");
        }
        const response = responses[0].response;
        switch (response.case) {
          case "secret":
            return response.value;
          case "error":
            throw new SecretsError(secretRequest, response.value.error);
          default:
            throw new SecretsError(secretRequest, "cannot unmarshal returned value from host");
        }
      }
    };
  }
  report(input) {
    const consensus = new ConsensusCapability;
    const call = consensus.report(this, input);
    return {
      result: () => {
        return call.result();
      }
    };
  }
}

class Runtime extends RuntimeImpl {
  constructor(config, nextCallId, maxResponseSize) {
    super(config, nextCallId, WasmRuntimeHelpers.getInstance(), maxResponseSize);
  }
}

class WasmRuntimeHelpers {
  static instance;
  constructor() {}
  now() {
    return hostBindings.now();
  }
  static getInstance() {
    if (!WasmRuntimeHelpers.instance) {
      WasmRuntimeHelpers.instance = new WasmRuntimeHelpers;
    }
    return WasmRuntimeHelpers.instance;
  }
  call(request) {
    return hostBindings.callCapability(toBinary(CapabilityRequestSchema, request)) >= 0;
  }
  await(request, maxResponseSize) {
    const responseSize = Math.trunc(Number(maxResponseSize));
    const response = hostBindings.awaitCapabilities(toBinary(AwaitCapabilitiesRequestSchema, request), responseSize);
    const responseBytes = Array.isArray(response) ? new Uint8Array(response) : response;
    return fromBinary(AwaitCapabilitiesResponseSchema, responseBytes);
  }
  getSecrets(request, maxResponseSize) {
    const responseSize = Math.trunc(Number(maxResponseSize));
    return hostBindings.getSecrets(toBinary(GetSecretsRequestSchema, request), responseSize) >= 0;
  }
  awaitSecrets(request, maxResponseSize) {
    const responseSize = Math.trunc(Number(maxResponseSize));
    const response = hostBindings.awaitSecrets(toBinary(AwaitSecretsRequestSchema, request), responseSize);
    const responseBytes = Array.isArray(response) ? new Uint8Array(response) : response;
    return fromBinary(AwaitSecretsResponseSchema, responseBytes);
  }
  switchModes(mode) {
    hostBindings.switchModes(mode);
  }
  log(message) {
    hostBindings.log(message);
  }
}

class Runner {
  config;
  request;
  constructor(config, request) {
    this.config = config;
    this.request = request;
  }
  static async newRunner(configHandlerParams) {
    hostBindings.versionV2();
    const request = Runner.getRequest();
    const config = await configHandler(request, configHandlerParams);
    return new Runner(config, request);
  }
  static getRequest() {
    const argsString = hostBindings.getWasiArgs();
    let args;
    try {
      args = JSON.parse(argsString);
    } catch (e) {
      throw new Error("Invalid request: could not parse arguments");
    }
    if (args.length !== 2) {
      throw new Error("Invalid request: must contain payload");
    }
    const base64Request = args[1];
    const bytes = Buffer.from(base64Request, "base64");
    return fromBinary(ExecuteRequestSchema, bytes);
  }
  async run(initFn) {
    const runtime = new Runtime(this.config, 0, this.request.maxResponseSize);
    var result;
    try {
      const workflow = await initFn(this.config, {
        getSecret: runtime.getSecret.bind(runtime)
      });
      switch (this.request.request.case) {
        case "subscribe":
          result = this.handleSubscribePhase(this.request, workflow);
          break;
        case "trigger":
          result = this.handleExecutionPhase(this.request, workflow, runtime);
          break;
        default:
          throw new Error("Unknown request type");
      }
    } catch (e) {
      const err = e instanceof Error ? e.message : String(e);
      result = create(ExecutionResultSchema, {
        result: { case: "error", value: err }
      });
    }
    const awaitedResult = await result;
    hostBindings.sendResponse(toBinary(ExecutionResultSchema, awaitedResult));
  }
  async handleExecutionPhase(req, workflow, runtime) {
    if (req.request.case !== "trigger") {
      throw new Error("cannot handle non-trigger request as a trigger");
    }
    const triggerMsg = req.request.value;
    const id = BigInt(triggerMsg.id);
    if (id > BigInt(Number.MAX_SAFE_INTEGER)) {
      throw new Error(`Trigger ID ${id} exceeds safe integer range`);
    }
    const index = Number(triggerMsg.id);
    if (Number.isFinite(index) && index >= 0 && index < workflow.length) {
      const entry = workflow[index];
      const schema = entry.trigger.outputSchema();
      const payloadAny = triggerMsg.payload;
      const decoded = fromBinary(schema, payloadAny.value);
      const adapted = entry.trigger.adapt(decoded);
      try {
        const result = await entry.fn(runtime, adapted);
        const wrapped = Value.wrap(result);
        return create(ExecutionResultSchema, {
          result: { case: "value", value: wrapped.proto() }
        });
      } catch (e) {
        const err = e instanceof Error ? e.message : String(e);
        return create(ExecutionResultSchema, {
          result: { case: "error", value: err }
        });
      }
    }
    return create(ExecutionResultSchema, {
      result: { case: "error", value: "trigger not found" }
    });
  }
  handleSubscribePhase(req, workflow) {
    if (req.request.case !== "subscribe") {
      return create(ExecutionResultSchema, {
        result: { case: "error", value: "subscribe request expected" }
      });
    }
    const subscriptions = workflow.map((entry) => ({
      id: entry.trigger.capabilityId(),
      method: entry.trigger.method(),
      payload: entry.trigger.configAsAny()
    }));
    const subscriptionRequest = create(TriggerSubscriptionRequestSchema, {
      subscriptions
    });
    return create(ExecutionResultSchema, {
      result: { case: "triggerSubscriptions", value: subscriptionRequest }
    });
  }
}
var onCronTrigger = (runtime2) => {
  runtime2.log("Hello world! Workflow triggered.");
  return "Hello world!";
};
var initWorkflow = (config) => {
  const cron = new cre.capabilities.CronCapability;
  return [cre.handler(cron.trigger({ schedule: config.schedule }), onCronTrigger)];
};
async function main() {
  const runner = await Runner.newRunner();
  await runner.run(initWorkflow);
}
main();
export {
  main
};
