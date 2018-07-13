// Copyright 2013-2015 The Rust Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// http://rust-lang.org/COPYRIGHT.
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

extern crate num;
#[macro_use]
extern crate num_derive;

#[derive(Debug, PartialEq, FromPrimitive)]
enum Color {
    Red,
    Blue = 5,
    Green,
}

#[test]
fn test_from_primitive_for_enum_with_custom_value() {
    let v: [Option<Color>; 4] = [num::FromPrimitive::from_u64(0),
                                 num::FromPrimitive::from_u64(5),
                                 num::FromPrimitive::from_u64(6),
                                 num::FromPrimitive::from_u64(3)];

    assert_eq!(v,
               [Some(Color::Red), Some(Color::Blue), Some(Color::Green), None]);
}
