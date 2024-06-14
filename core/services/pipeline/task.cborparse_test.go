package pipeline_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

func TestCBORParseTask(t *testing.T) {
	tests := []struct {
		name                  string
		data                  string
		vars                  pipeline.Vars
		inputs                []pipeline.Result
		expected              map[string]interface{}
		expectedErrorCause    error
		expectedErrorContains string
	}{
		{
			"hello world",
			"$(foo)",
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": "0xbf6375726c781a68747470733a2f2f657468657270726963652e636f6d2f61706964706174689f66726563656e7463757364ffff",
			}),
			nil,
			map[string]interface{}{
				"path": []interface{}{"recent", "usd"},
				"url":  "https://etherprice.com/api",
			},
			nil,
			"",
		},
		{
			"trailing empty bytes",
			"$(foo)",
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": "0xbf6375726c781a68747470733a2f2f657468657270726963652e636f6d2f61706964706174689f66726563656e7463757364ffff000000",
			}),
			nil,
			map[string]interface{}{
				"path": []interface{}{"recent", "usd"},
				"url":  "https://etherprice.com/api",
			},
			nil,
			"",
		},
		{
			"nested maps",
			"$(foo)",
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": "0xbf657461736b739f6868747470706f7374ff66706172616d73bf636d73676f68656c6c6f5f636861696e6c696e6b6375726c75687474703a2f2f6c6f63616c686f73743a36363930ffff",
			}),
			nil,
			map[string]interface{}{
				"params": map[string]interface{}{
					"msg": "hello_chainlink",
					"url": "http://localhost:6690",
				},
				"tasks": []interface{}{"httppost"},
			},
			nil,
			"",
		},
		{
			"bignums",
			"$(foo)",
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": "0x" +
					"bf" + // map(*)
					"67" + // text(7)
					"6269676e756d73" + // "bignums"
					"9f" + // array(*)
					"c2" + // tag(2) == unsigned bignum
					"5820" + // bytes(32)
					"0000000000000000000000000000000000000000000000010000000000000000" +
					// int(18446744073709551616)
					"c2" + // tag(2) == unsigned bignum
					"5820" + // bytes(32)
					"4000000000000000000000000000000000000000000000000000000000000000" +
					// int(28948022309329048855892746252171976963317496166410141009864396001978282409984)
					"c3" + // tag(3) == signed bignum
					"5820" + // bytes(32)
					"0000000000000000000000000000000000000000000000010000000000000000" +
					// int(18446744073709551616)
					"c3" + // tag(3) == signed bignum
					"5820" + // bytes(32)
					"3fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
					// int(28948022309329048855892746252171976963317496166410141009864396001978282409983)
					"ff" + // primitive(*)
					"ff", // primitive(*)
			}),
			nil,
			map[string]interface{}{
				"bignums": []interface{}{
					testutils.MustParseBigInt(t, "18446744073709551616"),
					testutils.MustParseBigInt(t, "28948022309329048855892746252171976963317496166410141009864396001978282409984"),
					testutils.MustParseBigInt(t, "-18446744073709551617"),
					testutils.MustParseBigInt(t, "-28948022309329048855892746252171976963317496166410141009864396001978282409984"),
				},
			},
			nil,
			"",
		},

		{
			"empty data",
			"$(foo)",
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": nil,
			}),
			nil,
			map[string]interface{}{},
			nil,
			"data",
		},
		{
			"error input",
			"",
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Error: errors.New("foo")}},
			nil,
			pipeline.ErrTooManyErrors,
			"task inputs",
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			task := pipeline.CBORParseTask{
				BaseTask: pipeline.NewBaseTask(0, "cbor", nil, nil, 0),
				Data:     test.data,
			}

			result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), test.vars, test.inputs)
			assert.False(t, runInfo.IsPending)
			assert.False(t, runInfo.IsRetryable)

			if test.expectedErrorCause != nil {
				assert.Equal(t, test.expectedErrorCause, errors.Cause(result.Error))
				assert.Nil(t, result.Value)
				if test.expectedErrorContains != "" {
					assert.Contains(t, result.Error.Error(), test.expectedErrorContains)
				}
			} else {
				assert.NoError(t, result.Error)
				assert.Equal(t, test.expected, result.Value)
			}
		})
	}
}
