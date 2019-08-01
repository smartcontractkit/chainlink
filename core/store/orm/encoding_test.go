package orm

//func TestStore_addressParser(t *testing.T) {
//zero := &common.Address{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
//fifteen := &common.Address{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 15}

//val, err := parseAddress("")
//assert.NoError(t, err)
//assert.Equal(t, nil, val)

//val, err = parseAddress("0x000000000000000000000000000000000000000F")
//assert.NoError(t, err)
//assert.Equal(t, fifteen, val)

//val, err = parseAddress("0X000000000000000000000000000000000000000F")
//assert.NoError(t, err)
//assert.Equal(t, fifteen, val)

//val, err = parseAddress("0")
//assert.NoError(t, err)
//assert.Equal(t, zero, val)

//val, err = parseAddress("15")
//assert.NoError(t, err)
//assert.Equal(t, fifteen, val)

//val, err = parseAddress("0x0")
//assert.Error(t, err)

//val, err = parseAddress("x")
//assert.Error(t, err)
//}

//func TestStore_bigIntParser(t *testing.T) {
//val, err := parseBigInt("0")
//assert.NoError(t, err)
//assert.Equal(t, new(big.Int).SetInt64(0), val)

//val, err = parseBigInt("15")
//assert.NoError(t, err)
//assert.Equal(t, new(big.Int).SetInt64(15), val)

//val, err = parseBigInt("x")
//assert.Error(t, err)

//val, err = parseBigInt("")
//assert.Error(t, err)
//}

//func TestStore_levelParser(t *testing.T) {
//val, err := parseLogLevel("ERROR")
//assert.NoError(t, err)
//assert.Equal(t, LogLevel{zapcore.ErrorLevel}, val)

//val, err = parseLogLevel("")
//assert.NoError(t, err)
//assert.Equal(t, LogLevel{zapcore.InfoLevel}, val)

//val, err = parseLogLevel("primus sucks")
//assert.Error(t, err)
//}

//func TestStore_urlParser(t *testing.T) {
//tests := []struct {
//name      string
//input     string
//wantError bool
//}{
//{"valid URL", "http://localhost:3000", false},
//{"invalid URL", ":", true},
//{"empty URL", "", false},
//}

//for _, test := range tests {
//t.Run(test.name, func(t *testing.T) {
//i, err := parseURL(test.input)

//if test.wantError {
//assert.Error(t, err)
//} else {
//require.NoError(t, err)
//w, ok := i.(*url.URL)
//require.True(t, ok)
//assert.Equal(t, test.input, w.String())
//}
//})
//}
//}
