package treeout

import (
	"bytes"
	"strings"
)

// special character groups used in composing the heriarchy layout
const (
	BranchDelimiterBox = `└─ `

	BranchChainerBox = `│ `

	BranchSplitterBox = `├─ `

	Indent = "   "
)

type Tree struct {
	Doc      string
	branches []Branches

	isRoot bool
	level  int
	parent Branches
	index  int
	prefix string
}

type Branches interface {
	Child(string) Branches
	Add(Branches)
	ParentFunc(fn func(Branches))
	String() string

	padding() string
	children() []Branches
	prnt() Branches
	setPrefix(string)
	getPrefix() string
	selfIndex() int

	setBranches([]Branches)
	setLevel(level int)
	setParent(parent Branches)
	setIndex(index int)
}

func New(doc string) *Tree {
	return &Tree{
		Doc:    doc,
		isRoot: true,
		level:  0,
	}
}

func (t *Tree) setPrefix(s string) {
	t.prefix = s
}
func (t *Tree) getPrefix() string {
	return t.prefix
}

func (t *Tree) selfIndex() int {
	return t.index
}

func (t Tree) String() string {
	if t.isRoot {
		return foreachLine(t.Doc, func(total int, i int, s string) string {
			base := t.padding() + s
			if i == total-1 {
				return base
			}
			return base + "\n"
		}) + "\n" + formatArr(t.branches)
	}
	return t.branchLn(t.Doc) + formatArr(t.branches)
}

type sf func(int, int, string) string

// Apply given transformation func for each line in string
func foreachLine(str string, transform sf) (out string) {
	parts := strings.Split(str, "\n")
	for idx, line := range parts {
		out += transform(len(parts), idx, line)
	}
	return
}

func (t *Tree) padding() string {
	var padding string
	for i := 0; i <= t.level; i++ {
		padding += Indent
	}
	return padding
}

func (t *Tree) branchLn(doc string) string {
	if t.selfIndex() < len(t.prnt().children())-1 {
		return foreachLine(doc, func(total int, i int, s string) string {
			var base string
			if i == 0 {
				base = strings.TrimSuffix(t.getPrefix(), BranchChainerBox) + BranchSplitterBox + s
			} else {
				base = strings.TrimSuffix(t.getPrefix(), BranchChainerBox) + BranchChainerBox + " " + s
			}
			if i == total-1 {
				return base
			}
			return base + "\n"
		}) + "\n"
		// return strings.TrimSuffix(t.getPrefix(), BranchChainerBox) + BranchSplitterBox +
		// 	doc + "\n"
	}
	if t.selfIndex() == len(t.prnt().children())-1 {
		return foreachLine(doc, func(total int, i int, s string) string {
			var base string
			if i == 0 {
				base = strings.TrimSuffix(t.getPrefix(), BranchChainerBox) + BranchDelimiterBox + s
			} else {
				base = strings.TrimSuffix(t.getPrefix(), BranchChainerBox) + "   " + s
			}
			if i == total-1 {
				return base
			}
			return base + "\n"
		}) + "\n"
		// return strings.TrimSuffix(t.getPrefix(), BranchChainerBox) + BranchDelimiterBox +
		// 	doc + "\n"
	}
	return foreachLine(doc, func(total int, i int, s string) string {
		var base string
		if i == 0 {
			base = strings.TrimSuffix(t.getPrefix(), BranchChainerBox) + BranchDelimiterBox + s
		} else {
			base = strings.TrimSuffix(t.getPrefix(), BranchChainerBox) + BranchChainerBox + " " + s
		}
		if i == total-1 {
			return base
		}
		return base + "\n"
	}) + "\n"
	// return strings.TrimSuffix(t.getPrefix(), BranchChainerBox) + BranchDelimiterBox +
	// 	doc + "\n"
}

func (t *Tree) children() []Branches {
	return t.branches
}

func (t *Tree) prnt() Branches {
	return t.parent
}

func (t *Tree) Child(doc string) Branches {
	newT := &Tree{
		Doc:    doc,
		level:  t.level + 1,
		parent: t,
		index:  len(t.children()),
	}

	t.branches = append(t.branches, newT)
	return newT
}

func (t *Tree) ParentFunc(fn func(Branches)) {
	fn(t)
}

func (t *Tree) Add(children Branches) {
	children.setParent(t)
	children.setIndex(len(t.children()))
	t.branches = append(t.branches, children)
}

func (t *Tree) setBranches(branches []Branches) {
	t.branches = branches
}
func (t *Tree) setLevel(level int) {
	t.level = level
}
func (t *Tree) setParent(parent Branches) {
	t.parent = parent
}
func (t *Tree) setIndex(index int) {
	t.index = index
}

func formatArr(arr []Branches) string {
	var accumulator bytes.Buffer
	for i, v := range arr {
		if len(v.prnt().children()) > i+1 {
			v.setPrefix(v.prnt().getPrefix() + Indent + BranchChainerBox)
		} else {
			if i == len(arr)-1 {
				v.setPrefix(v.prnt().getPrefix() + Indent)
			} else {
				v.setPrefix(v.prnt().getPrefix() + Indent)
			}
		}
		accumulator.WriteString(v.String())
	}
	return accumulator.String()
}
