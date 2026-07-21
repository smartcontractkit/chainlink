# Agent Instructions: Gin Concurrency & Data Race Prevention in core/web

## Critical Rule: Do NOT call `gin.CreateTestContext()` or `gin.SetMode()` in tests

### Problem

`gin.CreateTestContext()` internally calls `gin.SetMode(gin.TestMode)` on every invocation.
`gin.SetMode()` mutates the global un-synchronized package variable `gin.mode`.
When tests run in parallel with `t.Parallel()`, calling `gin.CreateTestContext()` or `gin.SetMode()` causes a Go **DATA RACE** with concurrent router initialization or handler calls reading `gin.mode` (e.g., via `gin.IsDebugging()`).

### Solution

1. **Never call `gin.CreateTestContext()` or `gin.SetMode()` inside individual unit tests.**
2. **Use `gin.New()` and `engine.ServeHTTP(w, req)`** to test controller handlers:
   ```go
   w := httptest.NewRecorder()
   engine := gin.New()
   engine.POST("/endpoint", controller.Handler)

   req, err := http.NewRequestWithContext(t.Context(), "POST", "/endpoint", body)
   require.NoError(t, err)

   engine.ServeHTTP(w, req)
   ```
3. `TestMain()` in [main_test.go](file:///Users/adamhamrick/Projects/chainlink/core/web/main_test.go) sets `gin.SetMode(gin.TestMode)` once before parallel tests execute.
