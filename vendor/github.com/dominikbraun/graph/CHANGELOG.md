# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),  and this project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.23.0] - 2023-07-05

**Are you using graph? [Check out the graph user survey](https://forms.gle/MLKUZKMeCRxTfj4v9)**

### Added
* Added the `AllPathsBetween` function for computing all paths between two vertices.

## [0.22.3] - 2023-06-14

### Changed
* Changed `StableTopologicalSort` to invoke the `less` function as few as possible, reducing comparisons.
* Changed `CreatesCycle` to use an optimized path if the default in-memory store is being used.
* Changed map allocations to use pre-defined memory sizes.

## [0.22.2] - 2023-06-06

### Fixed
* Fixed the major performance issues of `StableTopologicalSort`.

## [0.22.1] - 2023-06-05

### Fixed
* Fixed `TopologicalSort` to retain its original performance.

## [0.22.0] - 2023-05-24

### Added
* Added the `StableTopologicalSort` function for deterministic topological orderings.
* Added the `VertexAttributes` functional option for setting an entire vertex attributes map.

## [0.21.0] - 2023-05-18

### Added
* Added the `BFSWithDepth` function for performing a BFS with depth information.

### Fixed
* Fixed false positives of `ErrVertexHasEdges` when removing a vertex.

## [0.20.0] - 2023-05-01

**Release post: [graph Version 0.20 Is Out](https://dominikbraun.io/blog/graph-version-0.20-is-out/)**

### Added
* Added the `Graph.AddVerticesFrom` method for adding all vertices from another graph.
* Added the `Graph.AddEdgesFrom` method for adding all edges from another graph.
* Added the `Graph.Edges` method for obtaining all edges as a slice.
* Added the `Graph.UpdateEdge` method for updating the properties of an edge.
* Added the `Store.UpdateEdge` method for updating the properties of an edge.
* Added the `NewLike` function for creating a new graph that is "like" the given graph.
* Added the `EdgeAttributes` functional option for setting an entire edge attributes map.

### Changed
* Changed `Graph.Clone` to use the built-in in-memory store for storing vertices and edges for cloned graphs.

## [0.19.0] - 2023-04-23

### Added
* Added the `MinimumSpanningTree` function for finding a minimum spanning tree.
* Added the `MaximumSpanningTree` function for finding a maximum spanning tree.

## [0.18.0] - 2023-04-16

### Added
* Added the `Graph.RemoveVertex` method for removing a vertex.
* Added the `Store.RemoveVertex` method for removing a vertex.
* Added the `ErrVertexHasEdges` error instance.
* Added the `Union` function for combining two graphs into one.

## [0.17.0] - 2023-04-12

### Added
* Added the `draw.GraphAttributes` functional option for `draw.DOT` for rendering graph attributes.

### Changed
* Changed the library's GoDoc documentation.

## [0.16.2] - 2023-03-27

### Fixed
* Fixed `ShortestPath` for an edge case.

## [0.16.1] - 2023-03-06

### Fixed
* Fixed `TransitiveReduction` not to incorrectly report cycles.

## [0.16.0] - 2023-03-01

**This release contains breaking changes of the public API (see "Changed").**

### Added
* Added the `Store` interface, introducing support for custom storage implementations.
* Added the `NewWithStore` function for explicitly initializing a graph with a `Store` instance.
* Added the `EdgeData` functional option that can be used with `AddEdge`, introducing support for arbitrary data.
* Added the `Data` field to `EdgeProperties` for retrieving data added using `EdgeData`.

### Changed
* Changed `Order` to additionally return an error instance (breaking change).
* Changed `Size` to additionally return an error instance (breaking change).

## [0.15.1] - 2023-01-18

### Changed
* Changed `ShortestPath` to return `ErrTargetNotReachable` if the target vertex is not reachable.

### Fixed
* Fixed `ShortestPath` to return correct results for large unweighted graphs.

## [0.15.0] - 2022-11-25

### Added
* Added the `ErrVertexAlreadyExists` error instance. Use `errors.Is` to check for this instance.
* Added the `ErrEdgeAlreadyExists` error instance. Use `errors.Is` to check for this instance.
* Added the `ErrEdgeCreatesCycle` error instance. Use `errors.Is` to check for this instance.

### Changed
* Changed `AddVertex` to return `ErrVertexAlreadyExists` if the vertex already exists.
* Changed `VertexWithProperties` to return `ErrVertexNotFound` if the vertex doesn't exist.
* Changed `AddEdge` to return `ErrVertexNotFound` if either vertex doesn't exist.
* Changed `AddEdge` to return `ErrEdgeAlreadyExists` if the edge already exists.
* Changed `AddEdge` to return `ErrEdgeCreatesCycle` if cycle prevention is active and the edge would create a cycle.
* Changed `Edge` to return `ErrEdgeNotFound` if the edge doesn't exist.
* Changed `RemoveEdge` to return the error instances returned by `Edge`.

## [0.14.0] - 2022-11-01

### Added
* Added the `ErrVertexNotFound` error instance.

### Changed
* Changed `TopologicalSort` to fail at runtime when a cycle is detected.
* Changed `TransitiveReduction` to return the transitive reduction as a new graph and fail at runtime when a cycle is detected.
* Changed `Vertex` to return `ErrVertexNotFound` if the desired vertex couldn't be found.

## [0.13.0] - 2022-10-15

### Added
* Added the `VertexProperties` type for storing vertex-related properties.
* Added the `VertexWithProperties` method for retrieving a vertex and its properties.
* Added the `VertexWeight` functional option that can be used for `AddVertex`.
* Added the `VertexAttribute` functional option that can be used for `AddVertex`.
* Added support for rendering vertices with attributes using `draw.DOT`.

### Changed
* Changed `AddVertex` to accept functional options.
* Renamed `PermitCycles` to `PreventCycles`. This seems to be the price to pay if English isn't a library author's native language.

### Fixed
* Fixed the behavior of `ShortestPath` when the target vertex is not reachable from one of the visited vertices.

## [0.12.0] - 2022-09-19

### Added
* Added the `PermitCycles` option to explicitly prevent the creation of cycles.

### Changed
* Changed the `Acyclic` option to not implicitly impose cycle checks for operations like `AddEdge`. To prevent the creation of cycles, use `PermitCycles`.
* Changed `TopologicalSort` to only work for graphs created with `PermitCycles`. This is temporary.
* Changed `TransitiveReduction` to only work for graphs created with `PermitCycles`. This is temporary.

## [0.11.0] - 2022-09-15

### Added
* Added the `Order` method for retrieving the number of vertices in the graph.
* Added the `Size` method for retrieving the number of edges in the graph.

### Changed
* Changed the `graph` logo.
* Changed an internal operation of `ShortestPath` from O(n) to O(log(n)) by implementing the priority queue as a binary heap. Note that the actual complexity might still be defined by `ShortestPath` itself.

### Fixed
* Fixed `draw.DOT` to work correctly with vertices that contain special characters and whitespaces.

## [0.10.0] - 2022-09-09

### Added
* Added the `PredecessorMap` method for obtaining a map with all predecessors of each vertex.
* Added the `RemoveEdge` method for removing the edge between two vertices.
* Added the `Clone` method for retrieving a deep copy of the graph.
* Added the `TopologicalSort` function for obtaining the topological order of the vertices in the graph.
* Added the `TransitiveReduction` function for transforming the graph into its transitive reduction.

### Changed
* Changed the `visit` function of `DFS` to accept a vertex hash instead of the vertex value (i.e. `K` instead of `T`).
* Changed the `visit` function of `BFS` to accept a vertex hash instead of the vertex value (i.e. `K` instead of `T`).

### Removed
* Removed the `Predecessors` function. Use `PredecessorMap` instead and look up the respective vertex.

## [0.9.0] - 2022-08-17

### Added
* Added the `Graph.AddVertex` method for adding a vertex. This replaces `Graph.Vertex`.
* Added the `Graph.AddEdge` method for creating an edge. This replaces `Graph.Edge`.
* Added the `Graph.Vertex` method for retrieving a vertex by its hash. This is not to be confused with the old `Graph.Vertex` function for adding vertices that got replaced with `Graph.AddVertex`.
* Added the `Graph.Edge` method for retrieving an edge. This is not to be confused with the old `Graph.Edge` function for creating an edge that got replaced with `Graph.AddEdge`.
* Added the `Graph.Predecessors` function for retrieving a vertex' predecessors.
* Added the `DFS` function.
* Added the `BFS` function.
* Added the `CreatesCycle` function.
* Added the `StronglyConnectedComponents` function.
* Added the `ShortestPath` function.
* Added the `ErrEdgeNotFound` error indicating that a desired edge could not be found.

### Removed
* Removed the `Graph.EdgeByHashes` method. Use `Graph.AddEdge` instead.
* Removed the `Graph.GetEdgeByHashes` method. Use `Graph.Edge` instead.
* Removed the `Graph.DegreeByHash` method. Use `Graph.Degree` instead.
* Removed the `Graph.Degree` method.
* Removed the `Graph.DFS` and `Graph.DFSByHash` methods. Use `DFS` instead.
* Removed the `Graph.BFS` and `Graph.BFSByHash` methods. Use `BFS` instead.
* Removed the `Graph.CreatesCycle` and `Graph.CreatesCycleByHashes` methods. Use `CreatesCycle` instead.
* Removed the `Graph.StronglyConnectedComponents` method. Use `StronglyConnectedComponents` instead.
* Removed the `Graph.ShortestPath` and `Graph.ShortestPathByHash` methods. Use `ShortestPath` instead.

## [0.8.0] - 2022-08-01

### Added
* Added the `EdgeWeight` and `EdgeAttribute` functional options.
* Added the `Properties` field to `Edge`.

### Changed
* Changed `Edge` to accept a variadic `options` parameter.
* Changed `EdgeByHashes` to accept a variadic `options` parameter.
* Renamed `draw.Graph` to `draw.DOT` for more clarity regarding the rendering format.

### Removed
* Removed the `WeightedEdge` function. Use `Edge` with the `EdgeWeight` functional option instead.
* Removed the `WeightedEdgeByHashes` function. Use `EdgeByHashes` with the `EdgeWeight` functional option instead.

### Fixed
* Fixed missing edge attributes when drawing a graph using `draw.DOT`.

## [0.7.0] - 2022-07-26

### Added
* Added `draw` package for graph visualization using DOT-compatible renderers.
* Added `Traits` function for retrieving the graph's traits.

## [0.6.0] - 2022-07-22

### Added
* Added `AdjacencyMap` function for retrieving an adjancency map for all vertices.

### Removed
* Removed the `AdjacencyList` function.

## [0.5.0] - 2022-07-21

### Added
* Added `AdjacencyList` function for retrieving an adjacency list for all vertices.
  
### Changed
* Updated the examples in the documentation.

## [0.4.0] - 2022-07-01

### Added
* Added `ShortestPath` function for computing shortest paths.

### Changed
* Changed the term "properties" to "traits" in the code and documentation.
* Don't traverse all vertices in disconnected graphs by design.

## [0.3.0] - 2022-06-27

### Added
* Added `StronglyConnectedComponents` function for detecting SCCs.
* Added various images to usage examples.

## [0.2.0] - 2022-06-20

### Added
* Added `Degree` and `DegreeByHash` functions for determining vertex degrees.
* Added cycle checks when adding an edge using the `Edge` functions.

## [0.1.0] - 2022-06-19

### Added
* Added `CreatesCycle` and `CreatesCycleByHashes` functions for predicting cycles.

## [0.1.0-beta] - 2022-06-17

### Changed
* Introduced dedicated types for directed and undirected graphs, making `Graph[K, T]` an interface.

## [0.1.0-alpha] - 2022-06-13

### Added
* Introduced core types and methods.
