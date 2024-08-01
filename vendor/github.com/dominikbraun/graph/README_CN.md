[中文版](README_CN.md) | [English Version](README.md)

# <img src="img/logo.svg" width="300">

这是一款用于创建通用图数据结构、对其进行修改、分析和可视化的库。

# 特性

* 支持任意类型的通用顶点，例如 `int` 或 `City`。
* 图的特征和相应的验证，例如在无环图中进行循环检查。
* 寻找路径或连通图的算法，例如最短路径或强连通图。
* 转换和表示的算法，例如传递闭包或拓扑排序。
* 非递归图遍历的算法，例如 DFS 或 BFS。
* 顶点和边可以包含可选的元数据，例如权重或自定义属性。
* 使用 DOT 语言和 Graphviz 进行图形可视化。
* 通过使用自己的 `Store` 实现，可以集成任何存储后端。
* 包含广泛的测试，覆盖率约为 90%，且没有任何依赖项。

> 状态：由于 graph 版本处于 0.x 阶段，公共 API 不应被视为稳定的。
 
> README 可能包含未发布的更改。请查看 [latest documentation](https://pkg.go.dev/github.com/dominikbraun/graph).

# 入门指南

```
go get github.com/dominikbraun/graph
```

# 快速示例

## 创建整数类型节点ID图

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

## 创建整数类型节点ID有向无环图

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

## 创建自定义类型节点ID图

要详细了解此示例，请参见 [concept of hashes](https://pkg.go.dev/github.com/dominikbraun/graph@v0.17.0-rc4#hdr-Hashes).

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

## 创建边带权重的图

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

## 执行深度优先搜索

这个示例按 DFS 顺序遍历并打印图中的所有顶点。

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

## 查找强联通分量

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

## 查找最短路径

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

## 查找生成树

![minimum spanning tree](img/mst.svg)

```go
g := graph.New(graph.StringHash, graph.Weighted())

// Add vertices and edges ...

mst, _ := graph.MinimumSpanningTree(g)
```

## 执行拓扑排序

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

## 执行传递闭包削减

![transitive reduction](img/transitive-reduction-before.svg)

```go
g := graph.New(graph.StringHash, graph.Directed(), graph.PreventCycles())

// Add vertices and edges ...

transitiveReduction, _ := graph.TransitiveReduction(g)
```

![transitive reduction](img/transitive-reduction-after.svg)

## 禁止创建环路

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
panic: 在 2 和 3 之间创建的边将会引入一个环
```

## 使用 Graphviz 图可视化

以下示例将为 `g` 生成一个 DOT 描述，并将其写入给定的文件中。

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

要使用 Graphviz 从创建的文件生成 SVG，请使用如下命令：

```
dot -Tsvg -O mygraph.gv
```

`DOT` 函数还支持渲染图属性：

```go
_ = draw.DOT(g, file, draw.GraphAttribute("label", "my-graph"))
```

### 按照此文档绘制图

![simple graph](img/simple.svg)

图使用以下程序进行渲染：

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

使用 neato 引擎进行可视化：

```
dot -Tsvg -Kneato -O simple.gv
```

这个例子使用Graphviz支持的 [Brewer color scheme](https://graphviz.org/doc/info/colors.html#brewer)。

## 存储边属性

边可以具有一个或多个属性，用于存储元数据。在[visualizing a graph](#visualize-a-graph-using-graphviz) 时将考虑这些属性。
例如，此边将呈现为红色：

```go
_ = g.AddEdge(1, 2, graph.EdgeAttribute("color", "red"))
```

要获取所有支持的属性的概述，请查看
[DOT documentation](https://graphviz.org/doc/info/attrs.html).

The stored attributes can be retrieved by getting the edge and accessing the `Properties.Attributes`
field.
可以通过获取边并访问 `Properties.Attributes` 字段来检索存储的属性。

```go
edge, _ := g.Edge(1, 2)
color := edge.Properties.Attributes["color"] 
```

## 存储边数据

还可以在边上存储任意类型属性数据，而不仅仅是键值字符串对。此数据类型为 `any`。

```go
_  = g.AddEdge(1, 2, graph.EdgeData(myData))
```

可以通过获取边并访问 `Properties.Data` 字段来检索存储的数据。

```go
edge, _ := g.Edge(1, 2)
myData := edge.Properties.Data 
```

### 更新边数据

可以使用 `Graph.UpdateEdge` 更新边属性。以下示例向边 (A,B) 添加了一个新的 `color` 属性，并将边权重设置为 10。

```go
_ = g.UpdateEdge("A", "B", graph.EdgeAttribute("color", "red"), graph.EdgeWeight(10))
```

`Graph.UpdateEdge` 的方法签名和接受的函数选项与 `Graph.AddEdge` 完全相同。

## 存储点属性

顶点可能具有一个或多个属性，可用于存储元数据。在 [visualizing a graph](#visualize-a-graph-using-graphviz) 时将考虑这些属性。
例如，此顶点将以红色渲染：

```go
_ = g.AddVertex(1, graph.VertexAttribute("style", "filled"))
```

存储在顶点中的数据可以通过使用 `VertexWithProperties` 获取顶点，并访问 `Attributes` 字段来检索。

```go
vertex, properties, _ := g.VertexWithProperties(1)
style := properties.Attributes["style"]
```

要获取所有支持的属性的概述，请查看
[DOT documentation](https://graphviz.org/doc/info/attrs.html).

## 将图存储在自定义存储中

可以通过实现 `Store` 接口并使用它初始化一个新的图，来集成任何存储后端：

```go
g := graph.NewWithStore(graph.IntHash, myStore)
```

恰当实现 `Store` 接口，参考 [documentation](https://pkg.go.dev/github.com/dominikbraun/graph#Store)。
[`graph-sql`](https://github.com/dominikbraun/graph-sql) 是一个可直接使用的 SQL 存储实现。
# 文档

完整文档可在以下位置找到： [pkg.go.dev](https://pkg.go.dev/github.com/dominikbraun/graph).
