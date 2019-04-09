(module
  ;; Allocate a page of linear memory (64kb). Export it as "memory"
  (memory (export "memory") 1)

  ;; Write the string at the start of the linear memory.
  (data (i32.const 0) "Hello, world!") ;; write string at location 0

  ;; Export the position and length of the string.
  (global (export "length") i32 (i32.const 12))
  (global (export "position") i32 (i32.const 0)))
