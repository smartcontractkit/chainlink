# Gin Concurrency & Data Race Prevention

<intent>
Prevent race conditions in core/web unit tests caused by global Gin mode mutations.
</intent>

<rules>
<rule name="no-test-context-or-mode">
Never call gin.CreateTestContext() or gin.SetMode() in unit tests. TestMain() in main_test.go sets test mode once globally.
</rule>

<rule name="handler-test-pattern">
Test handlers with gin.New() and engine.ServeHTTP():

```go
w := httptest.NewRecorder()
engine := gin.New()
engine.POST("/endpoint", controller.Handler)

req, err := http.NewRequestWithContext(t.Context(), "POST", "/endpoint", body)
require.NoError(t, err)

engine.ServeHTTP(w, req)
```
</rule>
</rules>
