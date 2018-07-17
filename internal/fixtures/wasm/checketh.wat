(module
  (func $perform (param $value i64) (result i32)
    (i64.lt_s (i64.const 450) (get_local $value))
  )
  (export "perform" (func $perform))
)
