package main

var eip1559 = map[int]bool{
	1:  true,
	4:  true,
	5:  true,
	42: false,

	137:   false,
	80001: false,

	56: false,
	97: false,
}
