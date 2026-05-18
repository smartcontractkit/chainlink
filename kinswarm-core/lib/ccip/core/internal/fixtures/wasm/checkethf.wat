(module
  (func $perform (param $value f64) (result i32)
    (f64.lt (f64.const 450.0) (get_local $value))
  )
  (export "perform" (func $perform))
)
