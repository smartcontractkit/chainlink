package gethwrappers

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	gethParams "github.com/ethereum/go-ethereum/params"
)

// AbigenArgs is the arguments to the abigen executable. E.g., Bin is the -bin
// arg.

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
	if version != gethParams.Version {
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

	// Replace all anonymous structs in method return values
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", string(bs), parser.AllErrors)
	if err != nil {
		Exit("Error while improving abigen output", err)
	}

	type StructDef struct {
		ID  string
		Def string
	}

	structDefs := map[string]StructDef{}
	offset := 1
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Recv == nil { // Skip non-method functions
				return true
			} else if x.Recv.List[0].Names[0].Name == "it" { // We don't need to handle iterators
				return true
			}

			switch s := x.Type.Results.List[0].Type.(type) {
			case *ast.StructType:
				var (
					start      = int(s.Pos()) - offset
					end        = int(s.End()) - offset
					typeName   = x.Recv.List[0].Names[0].Name[1:]
					methodName = x.Name.Name
				)

				structDef := string(bs[start:end])

				trailing := make([]byte, len(bs[end:]))
				copy(trailing, bs[end:])

				replacementName := []byte(typeName + methodName + "Return")
				bs = append(bs[:start], replacementName...)
				bs = append(bs, trailing...)

				structDefs[string(replacementName)] = StructDef{
					ID:  getID(structDef),
					Def: structDef,
				}
				offset += (int(s.End()) - int(s.Pos())) - len(replacementName)
			}
		default:
			return true
		}

		return true
	})

	// Replace all other instances of the structs in question
	fset = token.NewFileSet()
	f, err = parser.ParseFile(fset, "", string(bs), parser.AllErrors)
	if err != nil {
		Exit("Error while improving abigen output", err)
	}
	offset = 1
	ast.Inspect(f, func(n ast.Node) bool {
		switch s := n.(type) {
		case *ast.StructType:
			var (
				start      = int(s.Pos()) - offset
				end        = int(s.End()) - offset
				bracketIdx = strings.Index(string(bs[start:]), "}")
				structSrc  = string(bs[start : start+bracketIdx+1])

				// We have to make sure we choose the right struct in case there are two methods with
				// identical return signatures. We do this by seeking the most recent `(` twice and
				// then extracting the struct name from the return signature.
				fnCallIdx    = strings.LastIndex(string(bs[:start-2]), "(")
				retvalIdx    = strings.LastIndex(string(bs[:fnCallIdx]), "(") + 1
				commaIdx     = strings.Index(string(bs[retvalIdx:]), ",")
				endRetvalIdx = strings.Index(string(bs[retvalIdx:]), ")")
			)
			var structNameEndIdx int
			if commaIdx == -1 {
				structNameEndIdx = endRetvalIdx
			} else {
				structNameEndIdx = commaIdx
			}
			structName := string(bs[retvalIdx : retvalIdx+structNameEndIdx])

			for _, def := range structDefs {
				if def.ID == getID(structSrc) {
					trailing := make([]byte, len(bs[end:]))
					copy(trailing, bs[end:])

					bs = append(bs[:start], structName...)
					bs = append(bs, trailing...)

					offset += (int(s.End()) - int(s.Pos())) - len(structName)
					break
				}
			}
		}
		return true
	})

	// Grab the contract name. It's always the first type in the file.
	fset = token.NewFileSet()
	f, err = parser.ParseFile(fset, "", string(bs), parser.AllErrors)
	if err != nil {
		Exit("Error while improving abigen output", err)
	}
	var contractName string
	ast.Inspect(f, func(n ast.Node) bool {
		switch s := n.(type) {
		case *ast.TypeSpec:
			if contractName == "" {
				contractName = s.Name.String()
			}
			return false
		}
		return true
	})

	// Generate an interface for the contract
	fset = token.NewFileSet()
	f, err = parser.ParseFile(fset, "", string(bs), parser.AllErrors)
	if err != nil {
		Exit("Error while improving abigen output", err)
	}

	var methods []string
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Recv == nil {
				return true
			}
			typeName := string(bs[int(x.Recv.List[0].Type.Pos()) : int(x.Recv.List[0].Type.End())-1])
			if typeName == contractName+"Caller" || typeName == contractName+"Transactor" || typeName == contractName+"Filterer" {
				nameStart := int(x.Name.Pos()) - 1
				bracketIdx := strings.Index(string(bs[nameStart:]), "{")
				methods = append(methods, string(bs[nameStart:nameStart+bracketIdx-1]))
			}
		}
		return true
	})

	// Write the named structs to the bottom of the file
	for name, def := range structDefs {
		src := strings.Replace(def.Def, "struct", "type "+name+" struct", 1)
		bs = append(bs, []byte("\n"+src+"\n")...)
	}

	// Write the the UnpackLog method to the bottom of the file
	bs = append(bs, []byte(fmt.Sprintf(`
func (_%v *%v) UnpackLog(out interface{}, event string, log types.Log) error {
    return _%v.%vFilterer.contract.UnpackLog(out, event, log)
}
`, contractName, contractName, contractName, contractName))...)

	// Write the interface to the bottom of the file
	var methodSrc string
	for _, method := range methods {
		methodSrc = methodSrc + "\t" + method + "\n"
	}
	methodSrc = fmt.Sprintf(`
type %vInterface interface {
%v
}
`, contractName, methodSrc)
	bs = append(bs, []byte(methodSrc)...)

	err = ioutil.WriteFile(path, bs, 0644)
	if err != nil {
		Exit("Error while writing improved abigen source", err)
	}
}

var re = regexp.MustCompile("\\s")

// We uniquely identify structs by removing all whitespace in their definitions.
func getID(structDef string) string {
	return re.ReplaceAllString(structDef, "")
}
