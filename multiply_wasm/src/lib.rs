#[cfg(no_std)]

#[no_mangle]
fn perform() {
}

#[cfg(test)]
mod tests {
    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4);
    }
}
