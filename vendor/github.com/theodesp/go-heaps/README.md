go-heaps
[![All Contributors](https://img.shields.io/badge/all_contributors-9-orange.svg?style=flat-square)](#contributors)
---
<a href="https://godoc.org/github.com/theodesp/go-heaps">
<img src="https://godoc.org/github.com/theodesp/go-heaps?status.svg" alt="GoDoc">
</a>

<a href="https://opensource.org/licenses/MIT" rel="nofollow">
<img src="https://img.shields.io/github/license/mashape/apistatus.svg" alt="License"/>
</a>

<a href="https://travis-ci.org/theodesp/go-heaps" rel="nofollow">
<img src="https://travis-ci.org/theodesp/go-heaps.svg?branch=master" />
</a>

<a href="https://codecov.io/gh/theodesp/go-heaps">
  <img src="https://codecov.io/gh/theodesp/go-heaps/branch/master/graph/badge.svg" />
</a>

Reference implementations of heap data structures in Go

## Installation
```bash
$ go get -u github.com/theodesp/go-heaps
```

## Contents

**Heaps**

* [Pairing Heap](https://en.wikipedia.org/wiki/Pairing_heap): A pairing heap is a type of heap data structure with relatively simple implementation and excellent practical amortized performance.
* [Leftist Heap](https://www.geeksforgeeks.org/leftist-tree-leftist-heap/): a variant of a binary heap. Every node has an s-value which is the distance to the nearest leaf. In contrast to a binary heap, a leftist tree attempts to be very unbalanced.
* [Skew Heap](https://en.wikipedia.org/wiki/Skew_heap): A skew heap (or self-adjusting heap) is a heap data structure implemented as a binary tree. Skew heaps are advantageous because of their ability to merge more quickly than binary heaps.
* [Fibonacci Heap](https://en.wikipedia.org/wiki/Fibonacci_heap): a Fibonacci heap is a data structure for priority queue operations, consisting of a collection of heap-ordered trees. It has a better amortized running time than many other priority queue data structures including the binary heap and binomial heap.
* [Binomial Heap](https://www.geeksforgeeks.org/binomial-heap-2/): A Binomial Heap is a collection of Binomial Trees. A Binomial Heap is a set of Binomial Trees where each Binomial Tree follows Min Heap property. And there can be at most one Binomial Tree of any degree.
* [Treap Heap](https://en.wikipedia.org/wiki/Treap): A Treap and the randomized binary search tree are two closely related forms of binary search tree data structures that maintain a dynamic set of ordered keys and allow binary searches among the keys.
* [Rank Pairing Heap](http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.153.4644&rep=rep1&type=pdf): A heap (priority queue) implementation that combines the asymptotic efficiency of Fibonacci heaps with much of the simplicity of pairing heaps

## Usage

```go
package main

import (
	"github.com/theodesp/go-heaps"
	pairingHeap "github.com/theodesp/go-heaps/pairing"
	"fmt"
)

func main()  {
	heap := pairingHeap.New()
	heap.Insert(go_heaps.Integer(4))
	heap.Insert(go_heaps.Integer(3))
	heap.Insert(go_heaps.Integer(2))
	heap.Insert(go_heaps.Integer(5))

	fmt.Println(heap.DeleteMin()) // 2
	fmt.Println(heap.DeleteMin()) // 3
	fmt.Println(heap.DeleteMin()) // 4
	fmt.Println(heap.DeleteMin()) // 5
}

```

## Complexity
| Operation     | Pairing       | Leftist      | Skew          | Fibonacci     | Binomial      | Treap         |
| ------------- |:-------------:|:-------------:|:-------------:|:-------------:|:-------------:|:-------------:|
| FindMin       | Θ(1)          | Θ(1)          | Θ(1)          | Θ(1)			| Θ(log n)      | O(n)          |
| DeleteMin     | O(log n)      | O(log n)      | O(log n)      | O(log n)	    | Θ(log n)      | O(n)          |
| Insert        | Θ(1)          | O(log n)      | O(log n)      | Θ(1)			| Θ(1)          | O(n)          |
| Find          | O(n)          |               |               |				|               |               |    
| Delete        | O(n)          |               | O(log n)      | O(n)			| Θ(log n)      | O(n)          |
| Adjust        | O(n)          |               | O(log n)      | O(n) 			| Θ(log n)      | O(n)          |
| Meld          | Θ(1)          |               |               |               |               |               |

| Operation     | Rank Pairing  | 
| ------------- |:-------------:|
| FindMin       | Θ(1)          |
| DeleteMin     | O(log n)      | 
| Insert        | Θ(1)          | 
| Find          | O(n)          |     
| Delete        | O(n)          |             
| Adjust        | O(n)          |
| Meld          | Θ(1)          |



## Contributors


Thanks goes to these wonderful people ([emoji key](https://github.com/kentcdodds/all-contributors#emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore -->
<table><tr><td align="center"><a href="http://mb-14.github.io"><img src="https://avatars1.githubusercontent.com/u/1137632?v=4" width="100px;" alt="Miroojin Bakshi"/><br /><sub><b>Miroojin Bakshi</b></sub></a><br /><a href="https://github.com/theodesp/go-heaps/commits?author=mb-14" title="Code">💻</a></td><td align="center"><a href="https://syfaro.net"><img src="https://avatars2.githubusercontent.com/u/1369709?v=4" width="100px;" alt="Syfaro"/><br /><sub><b>Syfaro</b></sub></a><br /><a href="https://github.com/theodesp/go-heaps/commits?author=Syfaro" title="Code">💻</a></td><td align="center"><a href="https://github.com/theodesp"><img src="https://avatars0.githubusercontent.com/u/328805?v=4" width="100px;" alt="Theofanis Despoudis"/><br /><sub><b>Theofanis Despoudis</b></sub></a><br /><a href="https://github.com/theodesp/go-heaps/commits?author=theodesp" title="Code">💻</a></td><td align="center"><a href="https://www.linkedin.com/in/ignacy-radlinski"><img src="https://avatars0.githubusercontent.com/u/26116041?v=4" width="100px;" alt="Radliński Ignacy"/><br /><sub><b>Radliński Ignacy</b></sub></a><br /><a href="https://github.com/theodesp/go-heaps/commits?author=radlinskii" title="Code">💻</a></td><td align="center"><a href="https://github.com/DonMcNamara"><img src="https://avatars3.githubusercontent.com/u/1152286?v=4" width="100px;" alt="Don McNamara"/><br /><sub><b>Don McNamara</b></sub></a><br /><a href="#infra-DonMcNamara" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a></td><td align="center"><a href="https://koneko096.github.io/"><img src="https://avatars3.githubusercontent.com/u/9217338?v=4" width="100px;" alt="Afrizal Fikri"/><br /><sub><b>Afrizal Fikri</b></sub></a><br /><a href="https://github.com/theodesp/go-heaps/commits?author=koneko096" title="Code">💻</a></td><td align="center"><a href="https://github.com/lhauspie"><img src="https://avatars1.githubusercontent.com/u/25682509?v=4" width="100px;" alt="Logan HAUSPIE"/><br /><sub><b>Logan HAUSPIE</b></sub></a><br /><a href="https://github.com/theodesp/go-heaps/commits?author=lhauspie" title="Code">💻</a></td></tr><tr><td align="center"><a href="https://ejq.me/"><img src="https://avatars3.githubusercontent.com/u/5517838?v=4" width="100px;" alt="Song Guo"/><br /><sub><b>Song Guo</b></sub></a><br /><a href="https://github.com/theodesp/go-heaps/commits?author=erjiaqing" title="Code">💻</a></td><td align="center"><a href="https://github.com/safwan-moha"><img src="https://avatars0.githubusercontent.com/u/7314967?v=4" width="100px;" alt="Safwan Mohammed"/><br /><sub><b>Safwan Mohammed</b></sub></a><br /><a href="https://github.com/theodesp/go-heaps/commits?author=safwan-moha" title="Tests">⚠️</a> <a href="https://github.com/theodesp/go-heaps/commits?author=safwan-moha" title="Code">💻</a></td></tr></table>

<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/kentcdodds/all-contributors) specification. Contributions of any kind welcome!


## LICENCE
Copyright © 2017 Theo Despoudis MIT license