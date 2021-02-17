package gethwrappers

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"

	gethParams "github.com/ethereum/go-ethereum/params"
	"github.com/smartcontractkit/chainlink/core/utils"
	"golang.org/x/tools/go/ast/astutil"
)

const headerComment = `// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

`

// AbigenArgs is the arguments to the abigen executable. E.g., Bin is the -bin
// arg.
type AbigenArgs struct {
	Bin, ABI, Out, Type, Pkg string
}

// Abigen calls Abigen  with the given arguments
//
// It might seem like a shame, to shell out to another golang program like
// this, but the abigen executable is the stable public interface to the
// geth contract-wrapper machinery.
//
// Check whether native abigen is installed, and has correct version
func Abigen(a AbigenArgs) {
	var versionResponse bytes.Buffer
	abigenExecutablePath := filepath.Join(GetProjectRoot(), "tools/bin/abigen")
	abigenVersionCheck := exec.Command(abigenExecutablePath, "--version")
	abigenVersionCheck.Stdout = &versionResponse
	if err := abigenVersionCheck.Run(); err != nil {
		Exit("no native abigen; you must install it (`make abigen` in the "+
			"chainlink root dir)", err)
	}
	version := string(regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+`).Find(
		versionResponse.Bytes()))
	// TODO: Re-enable once geth 1.10 is released with the abigen patch included.
	// We do this to avoid running un-released code in geth which is also present in the abigen bug fix patch.
	// This way we _only_ use the patched abigen for code generation.
	// if version != gethParams.Version {
	if version != "1.9.26" {
		Exit(fmt.Sprintf("wrong version (%s) of abigen; install the correct one "+
			"(%s) with `make abigen` in the chainlink root dir", version,
			gethParams.Version),
			nil)
	}
	buildCommand := exec.Command(
		abigenExecutablePath,
		"-bin", a.Bin,
		"-abi", a.ABI,
		"-out", a.Out,
		"-type", a.Type,
		"-pkg", a.Pkg,
	)
	var buildResponse bytes.Buffer
	buildCommand.Stderr = &buildResponse
	if err := buildCommand.Run(); err != nil {
		Exit("failure while building "+a.Pkg+" wrapper, stderr: "+
			buildResponse.String(), err)
	}

	ImproveAbigenOutput(a.Out)
}

func ImproveAbigenOutput(path string) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		Exit("Error while improving abigen output", err)
	}

	fset, fileNode := parseFile(bs)
	contractName := getContractName(fileNode)
	fileNode = addAddressField(contractName, fileNode)
	fileNode = replaceAnonymousStructs(contractName, fileNode)
	bs = generateCode(fset, fileNode)
	bs = writeAdditionalMethods(contractName, bs)
	fset, fileNode = parseFile(bs)
	fileNode = writeInterface(contractName, fileNode)
	bs = generateCode(fset, fileNode)
	bs = addHeader(bs)

	err = ioutil.WriteFile(path, bs, 0600)
	if err != nil {
		Exit("Error while writing improved abigen source", err)
	}
}

func parseFile(bs []byte) (*token.FileSet, *ast.File) {
	fset := token.NewFileSet()
	fileNode, err := parser.ParseFile(fset, "", string(bs), parser.AllErrors)
	if err != nil {
		Exit("Error while improving abigen output", err)
	}
	return fset, fileNode
}

func generateCode(fset *token.FileSet, fileNode *ast.File) []byte {
	var buf bytes.Buffer
	err := format.Node(&buf, fset, fileNode)
	if err != nil {
		Exit("Error while writing improved abigen source", err)
	}
	return buf.Bytes()
}

func getContractName(fileNode *ast.File) string {
	// Grab the contract name. It's always the first type in the file.
	var contractName string
	astutil.Apply(fileNode, func(cursor *astutil.Cursor) bool {
		x, is := cursor.Node().(*ast.TypeSpec)
		if !is {
			return true
		}
		if contractName == "" {
			contractName = x.Name.Name
		}
		return false
	}, nil)
	return contractName
}

func addAddressField(contractName string, fileNode *ast.File) *ast.File {
	// Add an exported `.Address` field to the contract struct
	fileNode = astutil.Apply(fileNode, func(cursor *astutil.Cursor) bool {
		x, is := cursor.Node().(*ast.StructType)
		if !is {
			return true
		}
		theType, is := cursor.Parent().(*ast.TypeSpec)
		if !is {
			return false
		} else if theType.Name.Name != contractName {
			return false
		}

		addrField := &ast.Field{
			Names: []*ast.Ident{ast.NewIdent("address")},
			Type: &ast.SelectorExpr{
				X:   ast.NewIdent("common"),
				Sel: ast.NewIdent("Address"),
			},
		}
		x.Fields.List = append([]*ast.Field{addrField}, x.Fields.List...)

		return false
	}, nil).(*ast.File)

	// Add the field to the return value of the constructor
	fileNode = astutil.Apply(fileNode, func(cursor *astutil.Cursor) bool {
		x, is := cursor.Node().(*ast.FuncDecl)
		if !is {
			return true
		} else if x.Name.Name != "New"+contractName {
			return false
		}

		for _, stmt := range x.Body.List {
			returnStmt, is := stmt.(*ast.ReturnStmt)
			if !is {
				continue
			}
			lit, is := returnStmt.Results[0].(*ast.UnaryExpr).X.(*ast.CompositeLit)
			if !is {
				continue
			}
			kvExpr := &ast.KeyValueExpr{
				Key:   ast.NewIdent("address"),
				Value: ast.NewIdent("address"),
			}
			lit.Elts = append([]ast.Expr{kvExpr}, lit.Elts...)
		}

		return false
	}, nil).(*ast.File)

	return fileNode
}

func replaceAnonymousStructs(contractName string, fileNode *ast.File) *ast.File {
	done := map[string]bool{}
	return astutil.Apply(fileNode, func(cursor *astutil.Cursor) bool {
		// Replace all anonymous structs with named structs
		x, is := cursor.Node().(*ast.FuncDecl)
		if !is {
			return true
		} else if len(x.Type.Results.List) == 0 {
			return false
		}
		theStruct, is := x.Type.Results.List[0].Type.(*ast.StructType)
		if !is {
			return false
		}

		methodName := x.Name.Name
		x.Type.Results.List[0].Type = ast.NewIdent(methodName)

		x.Body = astutil.Apply(x.Body, func(cursor *astutil.Cursor) bool {
			if _, is := cursor.Node().(*ast.StructType); !is {
				return true
			}
			if call, is := cursor.Parent().(*ast.CallExpr); is {
				for i, arg := range call.Args {
					if arg == cursor.Node() {
						call.Args[i] = ast.NewIdent(methodName)
						break
					}
				}
			}
			return true
		}, nil).(*ast.BlockStmt)

		if done[contractName+methodName] {
			return true
		}

		// Add the named structs to the bottom of the file
		fileNode.Decls = append(fileNode.Decls, &ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: ast.NewIdent(methodName),
					Type: theStruct,
				},
			},
		})

		done[contractName+methodName] = true
		return false
	}, nil).(*ast.File)
}

func writeAdditionalMethods(contractName string, bs []byte) []byte {
	// Write the the UnpackLog method to the bottom of the file
	bs = append(bs, []byte(fmt.Sprintf(`
func (_%v *%v) UnpackLog(out interface{}, event string, log types.Log) error {
    return _%v.%vFilterer.contract.UnpackLog(out, event, log)
}
`, contractName, contractName, contractName, contractName))...)

	// Write the the Address method to the bottom of the file
	bs = append(bs, []byte(fmt.Sprintf(`
func (_%v *%v) Address() common.Address {
    return _%v.address
}
`, contractName, contractName, contractName))...)

	return bs
}

func writeInterface(contractName string, fileNode *ast.File) *ast.File {
	// Generate an interface for the contract
	var methods []*ast.Field
	ast.Inspect(fileNode, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Recv == nil {
				return true
			}
			star, is := x.Recv.List[0].Type.(*ast.StarExpr)
			if !is {
				return false
			}

			typeName := star.X.(*ast.Ident).String()
			if typeName != contractName && typeName != contractName+"Caller" && typeName != contractName+"Transactor" && typeName != contractName+"Filterer" {
				return true
			}

			methods = append(methods, &ast.Field{
				Names: []*ast.Ident{x.Name},
				Type:  x.Type,
			})
		}
		return true
	})

	fileNode.Decls = append(fileNode.Decls, &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(contractName + "Interface"),
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: methods,
					},
				},
			},
		},
	})

	return fileNode
}

func addHeader(code []byte) []byte {
	return utils.ConcatBytes([]byte(headerComment), code)
}
