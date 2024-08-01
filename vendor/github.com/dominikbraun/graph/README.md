[中文版](README_CN.md) | [English Version](README.md)

# <img src="img/banner.png">

A library for creating generic graph data structures and modifying, analyzing,
and visualizing them.

**Are you using graph? [Check out the graph user survey.](https://forms.gle/MLKUZKMeCRxTfj4v9)**

# Features

* Generic vertices of any type, such as `int` or `City`.
* Graph traits with corresponding validations, such as cycle checks in acyclic graphs.
* Algorithms for finding paths or components, such as shortest paths or strongly connected components.
* Algorithms for transformations and representations, such as transitive reduction or topological order.
* Algorithms for non-recursive graph traversal, such as DFS or BFS.
* Vertices and edges with optional metadata, such as weights or custom attributes.
* Visualization of graphs using the DOT language and Graphviz.
* Integrate any storage backend by using your own `Store` implementation.
* Extensive tests with ~90% coverage, and zero dependencies.

> Status: Because `graph` is in version 0, the public API shouldn't be considered stable.

> This README may contain unreleased changes. Check out the [latest documentation](https://pkg.go.dev/github.com/dominikbraun/graph).

# Getting started

```
go get github.com/dominikbraun/graph
```

# Quick examples

## Create a graph of integers

![graph of integers](img/simple.svg)

```go
g := graph.New(graph.IntHash)

_ = g.AddVertex(1)
_ = g.AddVertex(2)
_ = g.AddVertex(3)
_ = g.AddVertex(4)
_ = g.AddVertex(5)

_ = g.AddEdge(1, 2)
_ = g.AddEdge(1, 4)
_ = g.AddEdge(2, 3)
_ = g.AddEdge(2, 4)
_ = g.AddEdge(2, 5)
_ = g.AddEdge(3, 5)
```

## Create a directed acyclic graph of integers

![directed acyclic graph](img/dag.svg)

```go
g := graph.New(graph.IntHash, graph.Directed(), graph.Acyclic())

_ = g.AddVertex(1)
_ = g.AddVertex(2)
_ = g.AddVertex(3)
_ = g.AddVertex(4)

_ = g.AddEdge(1, 2)
_ = g.AddEdge(1, 3)
_ = g.AddEdge(2, 3)
_ = g.AddEdge(2, 4)
_ = g.AddEdge(3, 4)
```

## Create a graph of a custom type

To understand this example in detail, see the [concept of hashes](https://pkg.go.dev/github.com/dominikbraun/graph#hdr-Hashes).

```go
type City struct {
    Name string
}

cityHash := func(c City) string {
    return c.Name
}

g := graph.New(cityHash)

_ = g.AddVertex(london)
```

## Create a weighted graph

![weighted graph](img/cities.svg)

```go
g := graph.New(cityHash, graph.Weighted())

_ = g.AddVertex(london)
_ = g.AddVertex(munich)
_ = g.AddVertex(paris)
_ = g.AddVertex(madrid)

_ = g.AddEdge("london", "munich", graph.EdgeWeight(3))
_ = g.AddEdge("london", "paris", graph.EdgeWeight(2))
_ = g.AddEdge("london", "madrid", graph.EdgeWeight(5))
_ = g.AddEdge("munich", "madrid", graph.EdgeWeight(6))
_ = g.AddEdge("munich", "paris", graph.EdgeWeight(2))
_ = g.AddEdge("paris", "madrid", graph.EdgeWeight(4))
```

## Perform a Depth-First Search

This example traverses and prints all vertices in the graph in DFS order.

![depth-first search](img/dfs.svg)

```go
g := graph.New(graph.IntHash, graph.Directed())

_ = g.AddVertex(1)
_ = g.AddVertex(2)
_ = g.AddVertex(3)
_ = g.AddVertex(4)

_ = g.AddEdge(1, 2)
_ = g.AddEdge(1, 3)
_ = g.AddEdge(3, 4)

_ = graph.DFS(g, 1, func(value int) bool {
    fmt.Println(value)
    return false
})
```

```
1 3 4 2
```

## Find strongly connected components

![strongly connected components](img/scc.svg)

```go
g := graph.New(graph.IntHash)

// Add vertices and edges ...

scc, _ := graph.StronglyConnectedComponents(g)

fmt.Println(scc)
```

```
[[1 2 5] [3 4 8] [6 7]]
```

## Find the shortest path

![shortest path algorithm](img/dijkstra.svg)

```go
g := graph.New(graph.StringHash, graph.Weighted())

// Add vertices and weighted edges ...

path, _ := graph.ShortestPath(g, "A", "B")

fmt.Println(path)
```

```
[A C E B]
```

## Find spanning trees

![minimum spanning tree](img/mst.svg)

```go
g := graph.New(graph.StringHash, graph.Weighted())

// Add vertices and edges ...

mst, _ := graph.MinimumSpanningTree(g)
```

## Perform a topological sort

![topological sort](img/topological-sort.svg)

```go
g := graph.New(graph.IntHash, graph.Directed(), graph.PreventCycles())

// Add vertices and edges ...

// For a deterministic topological ordering, use StableTopologicalSort.
order, _ := graph.TopologicalSort(g)

fmt.Println(order)
```

```
[1 2 3 4 5]
```

## Perform a transitive reduction

![transitive reduction](img/transitive-reduction-before.svg)

```go
g := graph.New(graph.StringHash, graph.Directed(), graph.PreventCycles())

// Add vertices and edges ...

transitiveReduction, _ := graph.TransitiveReduction(g)
```

![transitive reduction](img/transitive-reduction-after.svg)

## Prevent the creation of cycles

![cycle checks](img/cycles.svg)

```go
g := graph.New(graph.IntHash, graph.PreventCycles())

_ = g.AddVertex(1)
_ = g.AddVertex(2)
_ = g.AddVertex(3)

_ = g.AddEdge(1, 2)
_ = g.AddEdge(1, 3)

if err := g.AddEdge(2, 3); err != nil {
    panic(err)
}
```

```
panic: an edge between 2 and 3 would introduce a cycle
```

## Visualize a graph using Graphviz

The following example will generate a DOT description for `g` and write it into the given file.

```go
g := graph.New(graph.IntHash, graph.Directed())

_ = g.AddVertex(1)
_ = g.AddVertex(2)
_ = g.AddVertex(3)

_ = g.AddEdge(1, 2)
_ = g.AddEdge(1, 3)

file, _ := os.Create("./mygraph.gv")
_ = draw.DOT(g, file)
```

To generate an SVG from the created file using Graphviz, use a command such as the following:

```
dot -Tsvg -O mygraph.gv
```

The `DOT` function also supports rendering graph attributes:

```go
_ = draw.DOT(g, file, draw.GraphAttribute("label", "my-graph"))
```

### Draw a graph as in this documentation

![simple graph](img/simple.svg)

This graph has been rendered using the following program:

```go
package main

import (
	"os"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
)

func main() {
	g := graph.New(graph.IntHash)

	_ = g.AddVertex(1, graph.VertexAttribute("colorscheme", "blues3"), graph.VertexAttribute("style", "filled"), graph.VertexAttribute("color", "2"), graph.VertexAttribute("fillcolor", "1"))
	_ = g.AddVertex(2, graph.VertexAttribute("colorscheme", "greens3"), graph.VertexAttribute("style", "filled"), graph.VertexAttribute("color", "2"), graph.VertexAttribute("fillcolor", "1"))
	_ = g.AddVertex(3, graph.VertexAttribute("colorscheme", "purples3"), graph.VertexAttribute("style", "filled"), graph.VertexAttribute("color", "2"), graph.VertexAttribute("fillcolor", "1"))
	_ = g.AddVertex(4, graph.VertexAttribute("colorscheme", "ylorbr3"), graph.VertexAttribute("style", "filled"), graph.VertexAttribute("color", "2"), graph.VertexAttribute("fillcolor", "1"))
	_ = g.AddVertex(5, graph.VertexAttribute("colorscheme", "reds3"), graph.VertexAttribute("style", "filled"), graph.VertexAttribute("color", "2"), graph.VertexAttribute("fillcolor", "1"))

	_ = g.AddEdge(1, 2)
	_ = g.AddEdge(1, 4)
	_ = g.AddEdge(2, 3)
	_ = g.AddEdge(2, 4)
	_ = g.AddEdge(2, 5)
	_ = g.AddEdge(3, 5)

	file, _ := os.Create("./simple.gv")
	_ = draw.DOT(g, file)
}
```

It has been rendered using the `neato` engine:

```
dot -Tsvg -Kneato -O simple.gv
```

The example uses the [Brewer color scheme](https://graphviz.org/doc/info/colors.html#brewer) supported by Graphviz.

## Storing edge attributes

Edges may have one or more attributes which can be used to store metadata. Attributes will be taken
into account when [visualizing a graph](#visualize-a-graph-using-graphviz). For example, this edge
will be rendered in red color:

```go
_ = g.AddEdge(1, 2, graph.EdgeAttribute("color", "red"))
```

To get an overview of all supported attributes, take a look at the
[DOT documentation](https://graphviz.org/doc/info/attrs.html).

The stored attributes can be retrieved by getting the edge and accessing the `Properties.Attributes`
field.

```go
edge, _ := g.Edge(1, 2)
color := edge.Properties.Attributes["color"] 
```

## Storing edge data

It is also possible to store arbitrary data inside edges, not just key-value string pairs. This data
is of type `any`.

```go
_  = g.AddEdge(1, 2, graph.EdgeData(myData))
```

The stored data can be retrieved by getting the edge and accessing the `Properties.Data` field.

```go
edge, _ := g.Edge(1, 2)
myData := edge.Properties.Data 
```

### Updating edge data

Edge properties can be updated using `Graph.UpdateEdge`. The following example adds a new `color`
attribute to the edge (A,B) and sets the edge weight to 10.

```go
_ = g.UpdateEdge("A", "B", graph.EdgeAttribute("color", "red"), graph.EdgeWeight(10))
```

The method signature and the accepted functional options are exactly the same as for `Graph.AddEdge`.

## Storing vertex attributes

Vertices may have one or more attributes which can be used to store metadata. Attributes will be
taken into account when [visualizing a graph](#visualize-a-graph-using-graphviz). For example, this
vertex will be rendered in red color:

```go
_ = g.AddVertex(1, graph.VertexAttribute("style", "filled"))
```

The stored data can be retrieved by getting the vertex using `VertexWithProperties` and accessing
the `Attributes` field.

```go
vertex, properties, _ := g.VertexWithProperties(1)
style := properties.Attributes["style"]
```

To get an overview of all supported attributes, take a look at the
[DOT documentation](https://graphviz.org/doc/info/attrs.html).

## Store the graph in a custom storage

You can integrate any storage backend by implementing the `Store` interface and initializing a new
graph with it:

```go
g := graph.NewWithStore(graph.IntHash, myStore)
```

To implement the `Store` interface appropriately, take a look at the [documentation](https://pkg.go.dev/github.com/dominikbraun/graph#Store).
[`graph-sql`](https://github.com/dominikbraun/graph-sql) is a ready-to-use SQL store implementation.

# Documentation

The full documentation is available at [pkg.go.dev](https://pkg.go.dev/github.com/dominikbraun/graph).

**Are you using graph? [Check out the graph user survey.](https://forms.gle/MLKUZKMeCRxTfj4v9)**
