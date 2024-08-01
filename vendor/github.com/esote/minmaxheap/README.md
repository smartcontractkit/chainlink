# minmaxheap

Min-max heap operations for any type that implements heap.Interface. A min-max
heap can be used to implement a double-ended priority queue.

Min-max heap implementation from the 1986 paper "Min-Max Heaps and Generalized
Priority Queues" by Atkinson et. al. https://doi.org/10.1145/6617.6621.

| Operation | min-max heap | heap |
| --- | --- | --- |
| Init | O(n) | O(n) |
| Push | O(log n) | O(log n) |
| Pop | O(log n) | O(log n) |
| PopMax | **O(log n)** | O(n) |
| Fix | O(log n) | O(log n) |
