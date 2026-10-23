package docs

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/config/configtest"
)

func TestCoreDefaults_notNil(t *testing.T) {
	configtest.AssertFieldsNotNil(t, CoreDefaults())
}

func TestCoreDefaults_distinctInstances(t *testing.T) {
	t.Parallel()
	a := CoreDefaults()
	b := CoreDefaults()
	*a.Insecure.DevWebServer = true
	*a.AuditLogger.Enabled = true
	assert.False(t, *b.Insecure.DevWebServer)
	assert.False(t, *b.AuditLogger.Enabled)
}

func TestCoreDefaults_noSharedReferences(t *testing.T) {
	t.Parallel()
	a, b := CoreDefaults(), CoreDefaults()
	assertNoSharedRefs(t, "Core", reflect.ValueOf(a), reflect.ValueOf(b))
}

// assertNoSharedRefs walks a and b in lockstep and fails if any non-nil
// pointer, slice, or map is the same reference in both.
func assertNoSharedRefs(t *testing.T, path string, a, b reflect.Value) {
	t.Helper()
	switch a.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice:
		if a.IsNil() || b.IsNil() {
			return
		}
		if a.Kind() != reflect.Slice || a.Len() > 0 {
			assert.NotEqual(t, a.Pointer(), b.Pointer(), "%s is shared between instances", path)
		}
		switch a.Kind() {
		case reflect.Pointer:
			assertNoSharedRefs(t, path, a.Elem(), b.Elem())
		case reflect.Slice:
			for i := 0; i < a.Len() && i < b.Len(); i++ {
				assertNoSharedRefs(t, fmt.Sprintf("%s[%d]", path, i), a.Index(i), b.Index(i))
			}
		}
	case reflect.Struct:
		for i := 0; i < a.NumField(); i++ {
			if a.Type().Field(i).IsExported() {
				assertNoSharedRefs(t, path+"."+a.Type().Field(i).Name, a.Field(i), b.Field(i))
			}
		}
	}
}
