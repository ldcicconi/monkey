package object

import (
	"bytes"
	"fmt"
	"github.com/ldcicconi/monkey-interpreter/ast"
	"hash/fnv"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN_OBJ"
	ARRAY_OBJ        = "ARRAY_OBJ"
	HASH_OBJ         = "HASH_OBJ"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Integer struct {
	Value int64
}

func (i Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (i Integer) Type() ObjectType { return INTEGER_OBJ }
func (i Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}

func (b Boolean) HashKey() HashKey {
	val := 0
	if b.Value {
		val = 1
	}

	return HashKey{Type: b.Type(), Value: uint64(val)}
}

func (b Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (n Null) Type() ObjectType { return NULL_OBJ }
func (n Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message string
}

func (e Error) Type() ObjectType { return ERROR_OBJ }
func (e Error) Inspect() string  { return "ERROR: " + e.Message }

type Function struct {
	Parameters  []*ast.Identifier
	Body        *ast.BlockStatement
	Environment *Environment
}

func (f Function) Type() ObjectType { return FUNCTION_OBJ }
func (f Function) Inspect() string {
	var (
		out    bytes.Buffer
		params []string
	)

	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

type String struct {
	Value string
}

func (s String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

func (s String) Type() ObjectType { return STRING_OBJ }
func (s String) Inspect() string  { return s.Value }

type BuiltInFunction func(args ...Object) Object

type BuiltIn struct {
	Fn BuiltInFunction
}

func (b *BuiltIn) Type() ObjectType { return BUILTIN_OBJ }
func (b *BuiltIn) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var builder strings.Builder

	elements := make([]string, 0, len(a.Elements))
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	builder.WriteString("[")
	builder.WriteString(strings.Join(elements, ", "))
	builder.WriteString("]")

	return builder.String()
}

type HashPair struct {
	Key, Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h Hash) Type() ObjectType { return HASH_OBJ }

func (h Hash) Inspect() string {
	var builder strings.Builder
	pairs := make([]string, 0, len(h.Pairs))
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}
	builder.WriteString("{")
	builder.WriteString(strings.Join(pairs, ", "))
	builder.WriteString("}")
	return builder.String()
}

type Hashable interface {
	HashKey() HashKey
}
