package bin

import (
	"strings"
	"unicode"
)

// Ported from https://github.com/withoutboats/heck
// https://github.com/withoutboats/heck/blob/master/LICENSE-APACHE
// https://github.com/withoutboats/heck/blob/master/LICENSE-MIT

// ToPascalCase converts a string to upper camel case.
func ToPascalCase(s string) string {
	return transform(
		s,
		capitalize,
		func(*strings.Builder) {},
	)
}

func transform(
	s string,
	with_word func(string, *strings.Builder),
	boundary func(*strings.Builder),
) string {
	builder := new(strings.Builder)

	first_word := true
	words := splitIntoWords(s)
	for _, word := range words {
		char_indices := newReader(word)
		init := 0
		mode := _Boundary

		for char_indices.Move() {
			i, c := char_indices.This()

			// Skip underscore characters
			if c == '_' {
				if init == i {
					init += 1
				}
				continue
			}

			if next_i, next := char_indices.Peek(); next_i != -1 {

				// The mode including the current character, assuming the
				// current character does not result in a word boundary.
				next_mode := func() _WordMode {
					if unicode.IsLower(c) {
						return _Lowercase
					} else if unicode.IsUpper(c) {
						return _Uppercase
					} else {
						return mode
					}
				}()

				// Word boundary after if next is underscore or current is
				// not uppercase and next is uppercase
				if next == '_' || (next_mode == _Lowercase && unicode.IsUpper(next)) {
					if !first_word {
						// boundary(f)?;
						boundary(builder)
					}
					{
						// with_word(&word[init..next_i], f)?;
						with_word(word[init:next_i], builder)
					}

					first_word = false
					init = next_i
					mode = _Boundary

					// Otherwise if current and previous are uppercase and next
					// is lowercase, word boundary before
				} else if mode == _Uppercase && unicode.IsUpper(c) && unicode.IsLower(next) {
					if !first_word {
						// boundary(f)?;
						boundary(builder)
					} else {
						first_word = false
					}
					{
						// with_word(&word[init..i], f)?;
						with_word(word[init:i], builder)
					}
					init = i
					mode = _Boundary

					// Otherwise no word boundary, just update the mode
				} else {
					mode = next_mode
				}

			} else {
				// Collect trailing characters as a word
				if !first_word {
					// boundary(f)?;
					boundary(builder)
				} else {
					first_word = false
				}
				{
					// with_word(&word[init..], f)?;
					with_word(word[init:], builder)
				}
				break
			}
		}
	}

	return builder.String()
}

// fn capitalize(s: &str, f: &mut fmt::Formatter) -> fmt::Result {
//     let mut char_indices = s.char_indices();
//     if let Some((_, c)) = char_indices.next() {
//         write!(f, "{}", c.to_uppercase())?;
//         if let Some((i, _)) = char_indices.next() {
//             lowercase(&s[i..], f)?;
//         }
//     }

//     Ok(())
// }

func capitalize(s string, f *strings.Builder) {
	char_indices := newReader(s)
	if i, c := char_indices.Next(); i != -1 {
		f.WriteString(strings.ToUpper(string(c)))
		if i, _ := char_indices.Next(); i != -1 {
			lowercase(s[i:], f)
		}
	}
}

// fn lowercase(s: &str, f: &mut fmt::Formatter) -> fmt::Result {
//     let mut chars = s.chars().peekable();
//     while let Some(c) = chars.next() {
//         if c == 'Σ' && chars.peek().is_none() {
//             write!(f, "ς")?;
//         } else {
//             write!(f, "{}", c.to_lowercase())?;
//         }
//     }

//     Ok(())
// }

func lowercase(s string, f *strings.Builder) {
	chars := newReader(s)
	for chars.Move() {
		_, c := chars.This()
		if c == 'Σ' && chars.PeekNext() == 0 {
			f.WriteString("ς")
		} else {
			f.WriteString(strings.ToLower(string(c)))
		}
	}
}

func ToSnakeForSighash(s string) string {
	return ToRustSnakeCase(s)
}

type reader struct {
	runes []rune
	index int
}

func splitStringByRune(s string) []rune {
	var res []rune
	iterateStringAsRunes(s, func(r rune) bool {
		res = append(res, r)
		return true
	})
	return res
}

func iterateStringAsRunes(s string, callback func(r rune) bool) {
	for _, rn := range s {
		doContinue := callback(rn)
		if !doContinue {
			return
		}
	}
}

func newReader(s string) *reader {
	return &reader{
		runes: splitStringByRune(s),
		index: -1,
	}
}

func (r reader) This() (int, rune) {
	return r.index, r.runes[r.index]
}

func (r reader) HasNext() bool {
	return r.index < len(r.runes)-1
}

func (r *reader) Next() (int, rune) {
	if r.HasNext() {
		r.index++
		return r.index, r.runes[r.index]
	}
	return -1, rune(0)
}

func (r reader) Peek() (int, rune) {
	if r.HasNext() {
		return r.index + 1, r.runes[r.index+1]
	}
	return -1, rune(0)
}

func (r reader) PeekNext() rune {
	if r.HasNext() {
		return r.runes[r.index+1]
	}
	return rune(0)
}

func (r *reader) Move() bool {
	if r.HasNext() {
		r.index++
		return true
	}
	return false
}

// #[cfg(not(feature = "unicode"))]
func splitIntoWords(s string) []string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return !(unicode.IsLetter(r) || unicode.IsDigit(r))
	})
	return parts
}

type _WordMode int

const (
	/// There have been no lowercase or uppercase characters in the current
	/// word.
	_Boundary _WordMode = iota
	/// The previous cased character in the current word is lowercase.
	_Lowercase
	/// The previous cased character in the current word is uppercase.
	_Uppercase
)

// ToRustSnakeCase converts the given string to a snake_case string.
// Ported from https://github.com/withoutboats/heck/blob/c501fc95db91ce20eaef248a511caec7142208b4/src/lib.rs#L75 as used by Anchor.
func ToRustSnakeCase(s string) string {
	builder := new(strings.Builder)

	first_word := true
	words := splitIntoWords(s)
	for _, word := range words {
		char_indices := newReader(word)
		init := 0
		mode := _Boundary

		for char_indices.Move() {
			i, c := char_indices.This()

			// Skip underscore characters
			if c == '_' {
				if init == i {
					init += 1
				}
				continue
			}

			if next_i, next := char_indices.Peek(); next_i != -1 {

				// The mode including the current character, assuming the
				// current character does not result in a word boundary.
				next_mode := func() _WordMode {
					if unicode.IsLower(c) {
						return _Lowercase
					} else if unicode.IsUpper(c) {
						return _Uppercase
					} else {
						return mode
					}
				}()

				// Word boundary after if next is underscore or current is
				// not uppercase and next is uppercase
				if next == '_' || (next_mode == _Lowercase && unicode.IsUpper(next)) {
					if !first_word {
						// boundary(f)?;
						builder.WriteRune('_')
					}
					{
						// with_word(&word[init..next_i], f)?;
						builder.WriteString(strings.ToLower(word[init:next_i]))
					}

					first_word = false
					init = next_i
					mode = _Boundary

					// Otherwise if current and previous are uppercase and next
					// is lowercase, word boundary before
				} else if mode == _Uppercase && unicode.IsUpper(c) && unicode.IsLower(next) {
					if !first_word {
						// boundary(f)?;
						builder.WriteRune('_')
					} else {
						first_word = false
					}
					{
						// with_word(&word[init..i], f)?;
						builder.WriteString(strings.ToLower(word[init:i]))
					}
					init = i
					mode = _Boundary

					// Otherwise no word boundary, just update the mode
				} else {
					mode = next_mode
				}

			} else {
				// Collect trailing characters as a word
				if !first_word {
					// boundary(f)?;
					builder.WriteRune('_')
				} else {
					first_word = false
				}
				{
					// with_word(&word[init..], f)?;
					builder.WriteString(strings.ToLower(word[init:]))
				}
				break
			}
		}
	}

	return builder.String()
}
