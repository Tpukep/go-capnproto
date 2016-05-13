package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	C "github.com/glycerine/go-capnproto"
	"github.com/tpukep/bambam/bam"
	"github.com/tpukep/caps"
)

var g_nodes = make(map[uint64]*node)
var g_imported map[string]bool
var g_segment *C.Segment

const GO_CAPNP_IMPORT = "github.com/glycerine/go-capnproto"

type node struct {
	caps.Node
	pkg    string
	imp    string
	nodes  []*node
	name   string
	codecs map[uint64]bool
}

func assert(chk bool, format string, a ...interface{}) {
	if !chk {
		panic(fmt.Sprintf(format, a...))
		os.Exit(1)
	}
}

func copyData(obj C.Object) int {
	r, off, err := g_segment.NewRoot()
	assert(err == nil, "%v\n", err)
	err = r.Set(0, obj)
	assert(err == nil, "%v\n", err)
	return off
}

func findNode(id uint64) *node {
	n := g_nodes[id]
	assert(n != nil, "could not find node 0x%x\n", id)
	return n
}

func (n *node) remoteScope(from *node) string {
	assert(n.pkg != "", "missing package declaration for %s", n.DisplayName())

	if n.imp == from.imp {
		return ""
	} else {
		assert(n.imp != "", "missing import declaration for %s", n.DisplayName())

		g_imported[n.imp] = true
		return n.pkg + "."
	}
}

func (n *node) remoteName(from *node) string {
	return n.remoteScope(from) + n.name
}

func (n *node) resolveName(base, name string, file *node) {
	if na := nameAnnotation(n.Annotations()); na != "" {
		name = na
	}
	if base != "" {
		n.name = base + strings.Title(name)
	} else {
		n.name = strings.Title(name)
	}

	n.pkg = file.pkg
	n.imp = file.imp

	if n.Which() != caps.NODE_STRUCT || !n.Struct().IsGroup() {
		file.nodes = append(file.nodes, n)
	}

	for _, nn := range n.NestedNodes().ToArray() {
		if ni := g_nodes[nn.Id()]; ni != nil {
			ni.resolveName(n.name, nn.Name(), file)
		}
	}

	if n.Which() == caps.NODE_STRUCT {
		for _, f := range n.Struct().Fields().ToArray() {
			if f.Which() == caps.FIELD_GROUP {
				gname := f.Name()
				if na := nameAnnotation(f.Annotations()); na != "" {
					gname = na
				}
				findNode(f.Group().TypeId()).resolveName(n.name, gname, file)
			}
		}
	}
}

func nameAnnotation(annotations caps.Annotation_List) string {
	for _, a := range annotations.ToArray() {
		if a.Id() == C.Name {
			if name := a.Value().Text(); name != "" {
				return name
			}
		}
	}
	return ""
}

type enumval struct {
	caps.Enumerant
	val    int
	name   string
	tag    string
	parent *node
}

func (e *enumval) fullName() string {
	return fmt.Sprintf("%s_%s", strings.ToUpper(e.parent.name), strings.ToUpper(e.name))
}

func (n *node) defineEnum(w io.Writer, x *bam.Extractor) {
	for _, a := range n.Annotations().ToArray() {
		if a.Id() == C.Doc {
			fmt.Fprintf(w, "// %s\n", a.Value().Text())
		}
	}
	fmt.Fprintf(w, "type %s uint16\n", n.name)
	x.NewEnum(n.name)

	if es := n.Enum().Enumerants(); es.Len() > 0 {
		fmt.Fprintf(w, "const (\n")

		ev := make([]enumval, es.Len())
		for i := 0; i < es.Len(); i++ {
			e := es.At(i)
			ename := e.Name()
			if an := nameAnnotation(e.Annotations()); an != "" {
				ename = an
			}

			t := ename
			for _, an := range e.Annotations().ToArray() {
				if an.Id() == C.Tag {
					t = an.Value().Text()
				} else if an.Id() == C.Notag {
					t = ""
				}
			}
			ev[e.CodeOrder()] = enumval{e, i, ename, t, n}
		}

		// not an iota, so type has to go on each line
		for _, e := range ev {
			fmt.Fprintf(w, "%s %s = %d\n", e.fullName(), n.name, e.val)
		}

		fmt.Fprintf(w, ")\n")

		fmt.Fprintf(w, "func (c %s) String() string {\n", n.name)
		fmt.Fprintf(w, "switch c {\n")
		for _, e := range ev {
			if e.tag != "" {
				fmt.Fprintf(w, "case %s: return \"%s\"\n", e.fullName(), e.tag)
			}
		}
		fmt.Fprintf(w, "default: return \"\"\n")
		fmt.Fprintf(w, "}\n}\n\n")

		fmt.Fprintf(w, "func %sFromString(c string) %s {\n", n.name, n.name)
		fmt.Fprintf(w, "switch c {\n")
		for _, e := range ev {
			if e.tag != "" {
				fmt.Fprintf(w, "case \"%s\": return %s\n", e.tag, e.fullName())
			}
		}
		fmt.Fprintf(w, "default: return 0\n")
		fmt.Fprintf(w, "}\n}\n")
	}
}

func (n *node) writeValue(w io.Writer, t caps.Type, v caps.Value) {
	switch t.Which() {
	case caps.TYPE_VOID, caps.TYPE_INTERFACE:
		fmt.Fprintf(w, "struct{}")

	case caps.TYPE_BOOL:
		assert(v.Which() == caps.VALUE_BOOL, "expected bool value")
		if v.Bool() {
			fmt.Fprintf(w, "true")
		} else {
			fmt.Fprintf(w, "false")
		}

	case caps.TYPE_INT8:
		assert(v.Which() == caps.VALUE_INT8, "expected int8 value")
		fmt.Fprintf(w, "int8(%d)", v.Int8())

	case caps.TYPE_UINT8:
		assert(v.Which() == caps.VALUE_UINT8, "expected uint8 value")
		fmt.Fprintf(w, "uint8(%d)", v.Uint8())

	case caps.TYPE_INT16:
		assert(v.Which() == caps.VALUE_INT16, "expected int16 value")
		fmt.Fprintf(w, "int16(%d)", v.Int16())

	case caps.TYPE_UINT16:
		assert(v.Which() == caps.VALUE_UINT16, "expected uint16 value")
		fmt.Fprintf(w, "uint16(%d)", v.Uint16())

	case caps.TYPE_INT32:
		assert(v.Which() == caps.VALUE_INT32, "expected int32 value")
		fmt.Fprintf(w, "int32(%d)", v.Int32())

	case caps.TYPE_UINT32:
		assert(v.Which() == caps.VALUE_UINT32, "expected uint32 value")
		fmt.Fprintf(w, "uint32(%d)", v.Uint32())

	case caps.TYPE_INT64:
		assert(v.Which() == caps.VALUE_INT64, "expected int64 value")
		fmt.Fprintf(w, "int64(%d)", v.Int64())

	case caps.TYPE_UINT64:
		assert(v.Which() == caps.VALUE_UINT64, "expected uint64 value")
		fmt.Fprintf(w, "uint64(%d)", v.Uint64())

	case caps.TYPE_FLOAT32:
		assert(v.Which() == caps.VALUE_FLOAT32, "expected float32 value")
		fmt.Fprintf(w, "float32(%f)", v.Float32())

	case caps.TYPE_FLOAT64:
		assert(v.Which() == caps.VALUE_FLOAT64, "expected float64 value")
		fmt.Fprintf(w, "float64(%f)", v.Float64())

	case caps.TYPE_TEXT:
		assert(v.Which() == caps.VALUE_TEXT, "expected text value"+" got "+strconv.Itoa(int(v.Which())))
		fmt.Fprintf(w, "%s", strconv.Quote(v.Text()))

	case caps.TYPE_DATA:
		assert(v.Which() == caps.VALUE_DATA, "expected data value")
		fmt.Fprintf(w, "[]byte{")
		for i, b := range v.Data() {
			if i > 0 {
				fmt.Fprintf(w, ", ")
			}
			fmt.Fprintf(w, "%d", b)
		}
		fmt.Fprintf(w, "}")

	case caps.TYPE_ENUM:
		assert(v.Which() == caps.VALUE_ENUM, "expected enum value")
		en := findNode(t.Enum().TypeId())
		assert(en.Which() == caps.NODE_ENUM, "expected enum type ID")
		ev := en.Enum().Enumerants()
		if val := int(v.Enum()); val >= ev.Len() {
			fmt.Fprintf(w, "%s(%d)", en.remoteName(n), val)
		} else {
			fmt.Fprintf(w, "%s%s", en.remoteScope(n), ev.At(val).Name())
		}

	case caps.TYPE_STRUCT:
		fmt.Fprintf(w, "%s{ ", findNode(t.Struct().TypeId()).name)

		for _, f := range findNode(t.Struct().TypeId()).codeOrderFields() {
			fmt.Fprintf(w, "%s: ", strings.Title(f.Name()))

			switch slot := f.Slot(); slot.Type().Which() {
			case caps.TYPE_UINT8:
				val := v.Struct().ToStruct().Get8(int(slot.Offset()))
				fmt.Fprintf(w, "uint8(%d),", uint8(val))

			case caps.TYPE_UINT16:
				val := v.Struct().ToStruct().Get16(int(slot.Offset()))
				fmt.Fprintf(w, "uint16(%d),", uint16(val))

			case caps.TYPE_UINT32:
				val := v.Struct().ToStruct().Get32(int(slot.Offset()))
				fmt.Fprintf(w, "uint32(%d),", uint32(val))

			case caps.TYPE_UINT64:
				val := v.Struct().ToStruct().Get64(int(slot.Offset()))
				fmt.Fprintf(w, "uint64(%d),", uint64(val))

			case caps.TYPE_INT8:
				val := v.Struct().ToStruct().Get8(int(slot.Offset()))
				fmt.Fprintf(w, "int8(%d),", int8(val))

			case caps.TYPE_INT16:
				val := v.Struct().ToStruct().Get16(int(slot.Offset()))
				fmt.Fprintf(w, "int16(%d),", int16(val))

			case caps.TYPE_INT32:
				val := v.Struct().ToStruct().Get32(int(slot.Offset()))
				fmt.Fprintf(w, "int32(%d),", int32(val))

			case caps.TYPE_INT64:
				val := v.Struct().ToStruct().Get64(int(slot.Offset()))
				fmt.Fprintf(w, "int64(%d),", int64(val))

			case caps.TYPE_FLOAT32:
				val := v.Struct().ToStruct().Get32(int(slot.Offset()))
				fmt.Fprintf(w, "float32(%f),", float32(val))

			case caps.TYPE_FLOAT64:
				val := v.Struct().ToStruct().Get64(int(slot.Offset()))
				fmt.Fprintf(w, "float64(%f),", float64(val))

			case caps.TYPE_TEXT:
				val := v.Struct().ToStruct().GetObject(int(slot.Offset()))
				fmt.Fprintf(w, "\"%s\",", val.ToText())

			default:
				panic("Unsupported value type")
			}
		}

		fmt.Fprintf(w, "}\n")

	case caps.TYPE_ANYPOINTER:
		fmt.Fprintf(w, "interface{}")

	case caps.TYPE_LIST:
		assert(v.Which() == caps.VALUE_LIST, "expected list value")

		switch lt := t.List().ElementType(); lt.Which() {
		case caps.TYPE_VOID:
			fmt.Fprintf(w, "make([]C.Void, %d)", v.List().ToVoidList().Len())

		case caps.TYPE_INTERFACE:
			fmt.Fprintf(w, "make([]C.Interface, %d)", v.List().ToVoidList().Len())

		case caps.TYPE_BOOL:
			fmt.Fprintf(w, "[]bool{")
			for i, b := range v.List().ToBitList().ToArray() {
				if i > 0 {
					fmt.Fprintf(w, ", ")
				}
				fmt.Fprintf(w, "%v", b)
			}
			fmt.Fprintf(w, "}")

		case caps.TYPE_INT8:
			fmt.Fprintf(w, "[]int8{")
			for i, b := range v.List().ToInt8List().ToArray() {
				if i > 0 {
					fmt.Fprintf(w, ", ")
				}
				fmt.Fprintf(w, "%d", b)
			}
			fmt.Fprintf(w, "}")

		case caps.TYPE_UINT8:
			fmt.Fprintf(w, "[]uint8{")
			for i, b := range v.List().ToUInt8List().ToArray() {
				if i > 0 {
					fmt.Fprintf(w, ", ")
				}
				fmt.Fprintf(w, "%d", b)
			}
			fmt.Fprintf(w, "}")

		case caps.TYPE_INT16:
			fmt.Fprintf(w, "[]int16{")
			for i, b := range v.List().ToInt16List().ToArray() {
				if i > 0 {
					fmt.Fprintf(w, ", ")
				}
				fmt.Fprintf(w, "%d", b)
			}
			fmt.Fprintf(w, "}")

		case caps.TYPE_UINT16:
			fmt.Fprintf(w, "[]uint16{")
			for i, b := range v.List().ToUInt16List().ToArray() {
				if i > 0 {
					fmt.Fprintf(w, ", ")
				}
				fmt.Fprintf(w, "%d", b)
			}
			fmt.Fprintf(w, "}")

		case caps.TYPE_INT32:
			fmt.Fprintf(w, "[]int32{")
			for i, b := range v.List().ToInt32List().ToArray() {
				if i > 0 {
					fmt.Fprintf(w, ", ")
				}
				fmt.Fprintf(w, "%d", b)
			}
			fmt.Fprintf(w, "}")

		case caps.TYPE_UINT32:
			fmt.Fprintf(w, "[]uint32{")
			for i, b := range v.List().ToUInt32List().ToArray() {
				if i > 0 {
					fmt.Fprintf(w, ", ")
				}
				fmt.Fprintf(w, "%d", b)
			}
			fmt.Fprintf(w, "}")

		case caps.TYPE_FLOAT32:
			fmt.Fprintf(w, "[]float32{")
			for i, b := range v.List().ToFloat32List().ToArray() {
				if i > 0 {
					fmt.Fprintf(w, ", ")
				}
				fmt.Fprintf(w, "%g", b)
			}
			fmt.Fprintf(w, "}")

		case caps.TYPE_INT64:
			fmt.Fprintf(w, "[]int64{")
			for i, b := range v.List().ToInt64List().ToArray() {
				if i > 0 {
					fmt.Fprintf(w, ", ")
				}
				fmt.Fprintf(w, "%d", b)
			}
			fmt.Fprintf(w, "}")

		case caps.TYPE_UINT64:
			fmt.Fprintf(w, "[]uint64{")
			for i, b := range v.List().ToUInt64List().ToArray() {
				if i > 0 {
					fmt.Fprintf(w, ", ")
				}
				fmt.Fprintf(w, "%d", b)
			}
			fmt.Fprintf(w, "}")

		case caps.TYPE_FLOAT64:
			fmt.Fprintf(w, "[]float64{")
			for i, b := range v.List().ToFloat64List().ToArray() {
				if i > 0 {
					fmt.Fprintf(w, ", ")
				}
				fmt.Fprintf(w, "%g", b)
			}
			fmt.Fprintf(w, "}")

		case caps.TYPE_TEXT:
			fmt.Fprintf(w, "[]string{")
			for i, b := range v.List().ToTextList().ToArray() {
				if i > 0 {
					fmt.Fprintf(w, ", ")
				}
				fmt.Fprintf(w, "\"%s\"", b)
			}
			fmt.Fprintf(w, "}")

		case caps.TYPE_DATA:
			fmt.Fprintf(w, "[]byte{")
			for i, b := range v.List().ToDataList().ToArray() {
				if i > 0 {
					fmt.Fprintf(w, ", ")
				}
				fmt.Fprintf(w, "[]byte{")
				for i, ib := range b {
					if i > 0 {
						fmt.Fprintf(w, ", ")
					}
					fmt.Fprintf(w, "%d", ib)
				}
				fmt.Fprintf(w, "}")
			}
			fmt.Fprintf(w, "}")

		case caps.TYPE_ENUM:
			en := findNode(lt.Enum().TypeId())
			fmt.Fprintf(w, "[]%s{", en.remoteName(n))

			ev := en.Enum().Enumerants()
			a := v.List().ToUInt16List().ToEnumArray()

			for i, b := range *a {
				if i > 0 {
					fmt.Fprintf(w, ", ")
				}

				fmt.Fprintf(w, "%s", ev.At(int(b)).Name())
			}
			fmt.Fprintf(w, "}")
		case caps.TYPE_STRUCT:
			stype := findNode(lt.Struct().TypeId())
			fmt.Fprintf(w, "[]%s{", stype.remoteName(n))
			fmt.Fprintf(w, "/* Not implemented */")
			fmt.Fprintf(w, "}")
		case caps.TYPE_LIST, caps.TYPE_ANYPOINTER:
			fmt.Fprintf(w, "[]interface{")
			fmt.Fprintf(w, "/* Not implemented */")
			fmt.Fprintf(w, "}")
		}
	}
}

func (n *node) defineAnnotation(w io.Writer) {
	fmt.Fprintf(w, "const %s = uint64(0x%x)\n", n.name, n.Id())
}

func constIsVar(n *node) bool {
	switch n.Const().Type().Which() {
	case caps.TYPE_BOOL, caps.TYPE_INT8, caps.TYPE_UINT8, caps.TYPE_INT16,
		caps.TYPE_UINT16, caps.TYPE_INT32, caps.TYPE_UINT32, caps.TYPE_INT64,
		caps.TYPE_UINT64, caps.TYPE_FLOAT32, caps.TYPE_FLOAT64, caps.TYPE_TEXT, caps.TYPE_ENUM:
		return false
	default:
		return true
	}
}

func defineConstNodes(w io.Writer, nodes []*node) {
	any := false

	for _, n := range nodes {
		if n.Which() == caps.NODE_CONST && !constIsVar(n) {
			if !any {
				fmt.Fprintf(w, "const (\n")
				any = true
			}
			fmt.Fprintf(w, "%s = ", n.name)
			n.writeValue(w, n.Const().Type(), n.Const().Value())
			fmt.Fprintf(w, "\n")
		}
	}

	if any {
		fmt.Fprintf(w, ")\n")
	}

	any = false

	for _, n := range nodes {
		if n.Which() == caps.NODE_CONST && constIsVar(n) {
			if !any {
				fmt.Fprintf(w, "var (\n")
				any = true
			}
			fmt.Fprintf(w, "%s = ", n.name)
			n.writeValue(w, n.Const().Type(), n.Const().Value())
			fmt.Fprintf(w, "\n")
		}
	}

	if any {
		fmt.Fprintf(w, ")\n")
	}
}

func (n *node) defineField(w io.Writer, f caps.Field, x *bam.Extractor) {
	t := f.Slot().Type()

	if t.Which() == caps.TYPE_INTERFACE {
		return
	}

	var fname string

	if an := nameAnnotation(f.Annotations()); an != "" {
		fname = an
	} else {
		fname = f.Name()
	}

	fname = strings.Title(fname)

	var g, s bytes.Buffer

	if f.DiscriminantValue() != 0xFFFF {
		if t.Which() == caps.TYPE_VOID {
			x.SetUnionStruct()
			w.Write(s.Bytes())
			return
		}
	} else if t.Which() == caps.TYPE_VOID {
		return
	}

	customtype := ""
	for _, a := range f.Annotations().ToArray() {
		if a.Id() == C.Doc {
			fmt.Fprintf(&g, "// %s\n", a.Value().Text())
		}
		if a.Id() == C.Customtype {
			customtype = a.Value().Text()
			if i := strings.LastIndex(customtype, "."); i != -1 {
				g_imported[customtype[:i]] = true
			}
		}
	}

	if len(customtype) != 0 {
		log.Println("CUSTOM TYPE:", customtype)
	}

	fmt.Fprintf(&s, "%s ", fname)

	typeName := GoTypeName(n, f.Slot(), customtype)
	fmt.Fprintf(&s, "%s", typeName)

	fld := &ast.Field{}
	goseq := strings.SplitAfter(typeName, "[]")
	typePrefix := ""
	if len(goseq) == 2 {
		typeName = goseq[1]
		typePrefix = goseq[0]
	}

	x.GenerateStructField(fname, typePrefix, typeName, fld, t.Which() == caps.TYPE_LIST, fld.Tag, false, goseq)

	ans := f.Annotations()
	n.processAnnotations(&s, f, t.Which(), ans)

	fmt.Fprintf(&s, "\n")

	w.Write(g.Bytes())
	w.Write(s.Bytes())
}

func GoTypeName(n *node, s caps.FieldSlot, customtype string) string {
	def := s.DefaultValue()
	t := s.Type()

	switch t.Which() {
	case caps.TYPE_BOOL:
		assert(def.Which() == caps.VALUE_VOID || def.Which() == caps.VALUE_BOOL, "expected bool default")
		return "bool"

	case caps.TYPE_INT8:
		assert(def.Which() == caps.VALUE_VOID || def.Which() == caps.VALUE_INT8, "expected int8 default")
		return "int8"

	case caps.TYPE_UINT8:
		assert(def.Which() == caps.VALUE_VOID || def.Which() == caps.VALUE_UINT8, "expected uint8 default")
		return "uint8"

	case caps.TYPE_INT16:
		assert(def.Which() == caps.VALUE_VOID || def.Which() == caps.VALUE_INT16, "expected int16 default")
		return "int16"

	case caps.TYPE_UINT16:
		assert(def.Which() == caps.VALUE_VOID || def.Which() == caps.VALUE_UINT16, "expected uint16 default")
		return "uint16"

	case caps.TYPE_INT32:
		assert(def.Which() == caps.VALUE_VOID || def.Which() == caps.VALUE_INT32, "expected int32 default")
		return "int32"

	case caps.TYPE_UINT32:
		assert(def.Which() == caps.VALUE_VOID || def.Which() == caps.VALUE_UINT32, "expected uint32 default")
		return "uint32"

	case caps.TYPE_INT64:
		assert(def.Which() == caps.VALUE_VOID || def.Which() == caps.VALUE_INT64, "expected int64 default")
		return "int64"

	case caps.TYPE_UINT64:
		assert(def.Which() == caps.VALUE_VOID || def.Which() == caps.VALUE_UINT64, "expected uint64 default")
		return "uint64"

	case caps.TYPE_FLOAT32:
		assert(def.Which() == caps.VALUE_VOID || def.Which() == caps.VALUE_FLOAT32, "expected float32 default")
		return "float32"

	case caps.TYPE_FLOAT64:
		assert(def.Which() == caps.VALUE_VOID || def.Which() == caps.VALUE_FLOAT64, "expected float64 default")
		return "float64"

	case caps.TYPE_TEXT:
		assert(def.Which() == caps.VALUE_VOID || def.Which() == caps.VALUE_TEXT, "expected text default")

		return "string"

	case caps.TYPE_DATA:
		assert(def.Which() == caps.VALUE_VOID || def.Which() == caps.VALUE_DATA, "expected data default")
		if def.Which() == caps.VALUE_DATA && len(def.Data()) > 0 {
			dstr := "[]byte{"
			for i, b := range def.Data() {
				if i > 0 {
					dstr += ", "
				}
				dstr += fmt.Sprintf("%d", b)
			}
			dstr += "}"
			if len(customtype) != 0 {
				return fmt.Sprintf("%s\n", dstr)
			}
		} else {
			return "[]byte"
		}
	case caps.TYPE_ENUM:
		ni := findNode(t.Enum().TypeId())
		assert(def.Which() == caps.VALUE_VOID || def.Which() == caps.VALUE_ENUM, "expected enum default")
		return ni.remoteName(n)

	case caps.TYPE_STRUCT:
		ni := findNode(t.Struct().TypeId())
		assert(def.Which() == caps.VALUE_VOID || def.Which() == caps.VALUE_STRUCT, "expected struct default")
		return ni.remoteName(n)

	case caps.TYPE_ANYPOINTER:
		assert(def.Which() == caps.VALUE_VOID || def.Which() == caps.VALUE_ANYPOINTER, "expected object default")
		return "interface{}"

	case caps.TYPE_LIST:
		assert(def.Which() == caps.VALUE_VOID || def.Which() == caps.VALUE_LIST, "expected list default")

		switch lt := t.List().ElementType(); lt.Which() {
		case caps.TYPE_VOID, caps.TYPE_INTERFACE:
			return "[]struct{}"
		case caps.TYPE_BOOL:
			return "[]bool"
		case caps.TYPE_INT8:
			return "[]int8"
		case caps.TYPE_UINT8:
			return "[]uint8"
		case caps.TYPE_INT16:
			return "[]int16"
		case caps.TYPE_UINT16:
			return "[]uint16"
		case caps.TYPE_INT32:
			return "[]uint32"
		case caps.TYPE_UINT32:
			return "[]uint32"
		case caps.TYPE_INT64:
			return "[]int64"
		case caps.TYPE_UINT64:
			return "[]uint64"
		case caps.TYPE_FLOAT32:
			return "[]float32"
		case caps.TYPE_FLOAT64:
			return "[]float64"
		case caps.TYPE_TEXT:
			return "[]string"
		case caps.TYPE_DATA:
			return "[]byte"
		case caps.TYPE_ENUM:
			ni := findNode(lt.Enum().TypeId())
			return fmt.Sprintf("[]%s", ni.remoteName(n))
		case caps.TYPE_STRUCT:
			ni := findNode(lt.Struct().TypeId())

			return fmt.Sprintf("[]%s", ni.name)
		case caps.TYPE_ANYPOINTER, caps.TYPE_LIST:
			return "[]interface{}"
		}
	}

	panic("Unsupported type. Type_Which=" + strconv.Itoa(int(t.Which())))
}

func (n *node) codeOrderFields() []caps.Field {
	fields := n.Struct().Fields().ToArray()
	mbrs := make([]caps.Field, len(fields))
	for _, f := range fields {
		mbrs[f.CodeOrder()] = f
	}
	return mbrs
}

func (n *node) defineStructTypes(w io.Writer, baseNode *node, x *bam.Extractor) {
	assert(n.Which() == caps.NODE_STRUCT, "invalid struct node")

	for _, a := range n.Annotations().ToArray() {
		if a.Id() == C.Doc {
			fmt.Fprintf(w, "// %s\n", a.Value().Text())
		}
	}
	if baseNode == nil {
		x.StartStruct(n.name)

		fmt.Fprintf(w, "type %s struct {\n", n.name)
		n.defineStructFields(w, x)
		fmt.Fprintf(w, "}\n\n")

		baseNode = n
		x.EndStruct()
	}

	for _, f := range n.codeOrderFields() {
		if f.Which() == caps.FIELD_GROUP {
			findNode(f.Group().TypeId()).defineStructTypes(w, baseNode, x)
		}
	}
}

func (n *node) defineStructEnums(w io.Writer) {
	assert(n.Which() == caps.NODE_STRUCT, "invalid struct node")

	if n.Struct().DiscriminantCount() > 0 {
		fmt.Fprintf(w, "type %s_Which uint16\n", n.name)
		fmt.Fprintf(w, "const (\n")

		for _, f := range n.codeOrderFields() {
			if f.DiscriminantValue() == 0xFFFF {
				// Non-union member
			} else {
				fmt.Fprintf(w, "%s_%s %s_Which = %d\n", strings.ToUpper(n.name), strings.ToUpper(f.Name()), n.name, f.DiscriminantValue())
			}
		}
		fmt.Fprintf(w, ")\n")
	}

	for _, f := range n.codeOrderFields() {
		if f.Which() == caps.FIELD_GROUP {
			findNode(f.Group().TypeId()).defineStructEnums(w)
		}
	}
}

func (n *node) defineStructFields(w io.Writer, x *bam.Extractor) {
	assert(n.Which() == caps.NODE_STRUCT, "invalid struct node")

	for _, f := range n.codeOrderFields() {
		switch f.Which() {
		case caps.FIELD_SLOT:
			n.defineField(w, f, x)
		case caps.FIELD_GROUP:
			g := findNode(f.Group().TypeId())
			fname := f.Name()
			if an := nameAnnotation(f.Annotations()); an != "" {
				fname = an
			}
			fname = strings.Title(fname)

			typeName := ""
			fld := &ast.Field{}
			x.GenerateStructField(fname, "", typeName, fld, false, fld.Tag, true, []string{typeName})

			fmt.Fprintf(w, "%s struct {\n", fname)
			g.defineStructFields(w, x)

			fmt.Fprintf(w, "}\n")
		}
	}
}

func (n *node) writeImports(file *os.File) {
	if n.imp != "" || len(g_imported) > 0 {
		fmt.Fprintf(file, "import (\n")
		if n.imp != "" {
			fmt.Fprintf(file, "    %q\n", n.imp)
		}

		for imp := range g_imported {
			fmt.Fprintf(file, "    %q\n", imp)
		}

		fmt.Fprintf(file, ")\n\n")
	}
}

func (n *node) processAnnotations(w io.Writer, f caps.Field, t caps.Type_Which, ans caps.Annotation_List) {
	annotations := make(map[uint64]caps.Annotation)

	for _, a := range ans.ToArray() {
		annotations[a.Id()] = a
	}

	req, required := annotations[caps.FieldRequired]
	opt, optional := annotations[caps.FieldOptional]
	_, ignored := annotations[caps.FieldIgnored]

	assert(!(required && ignored), "Field annnotations 'required' and 'ignored' are incompatible.")
	assert(!(required && optional), "Annnotations 'required' and 'optional' are incompatible")
	assert(!(optional && ignored), "Annnotations 'optional' and 'ignored' are incompatible")

	var tags []string
	var checkTags []string

	// Codecs Tags
	if ignored {
		if _, found := n.codecs[caps.CodecJson]; found {
			tags = append(tags, fmt.Sprintf("json:\"-\""))
		}
		if _, found := n.codecs[caps.CodecMsgp]; found {
			tags = append(tags, fmt.Sprintf("msg:\"-\""))
		}

		checkTags = append(checkTags, "-")
	} else if optional {
		if _, found := n.codecs[caps.CodecJson]; found {
			tags = append(tags, fmt.Sprintf("json:\"%s,omitempty\"", opt.Value().Text()))
		}
		if _, found := n.codecs[caps.CodecMsgp]; found {
			tags = append(tags, fmt.Sprintf("msg:\"%s\"", opt.Value().Text()))
		}

		checkTags = append(checkTags, "omitempty")
	} else if required {
		if _, found := n.codecs[caps.CodecJson]; found {
			tags = append(tags, fmt.Sprintf("json:\"%s\"", req.Value().Text()))
		}
		if _, found := n.codecs[caps.CodecMsgp]; found {
			tags = append(tags, fmt.Sprintf("msg:\"%s\"", req.Value().Text()))
		}

		checkTags = append(checkTags, "required")
	} else {
		if _, found := n.codecs[caps.CodecJson]; found {
			tags = append(tags, fmt.Sprintf("json:\"%s\"", f.Name()))
		}
		if _, found := n.codecs[caps.CodecMsgp]; found {
			tags = append(tags, fmt.Sprintf("msg:\"%s\"", f.Name()))
		}
	}

	// Check annotation
	if exp, found := annotations[caps.CheckValue]; found {
		checkTags = append(checkTags, exp.Value().Text())
	}

	if len(checkTags) != 0 {
		tags = append(tags, fmt.Sprintf("validate:\"%s\"", strings.Join(checkTags, ",")))
	}

	if len(tags) != 0 {
		fmt.Fprintf(w, "`%s`", strings.Join(tags, " "))
	}
}

func main() {
	s, err := C.ReadFromStream(os.Stdin, nil)
	assert(err == nil, "%v\n", err)

	req := caps.ReadRootCodeGeneratorRequest(s)
	allfiles := []*node{}

	for _, ni := range req.Nodes().ToArray() {
		n := &node{Node: ni, codecs: make(map[uint64]bool)}
		g_nodes[n.Id()] = n

		if n.Which() == caps.NODE_FILE {
			allfiles = append(allfiles, n)
		}
	}

	g_imported = make(map[string]bool)

	for _, f := range allfiles {
		for _, a := range f.Annotations().ToArray() {
			if v := a.Value(); v.Which() == caps.VALUE_TEXT {
				switch a.Id() {
				case C.Package:
					f.pkg = v.Text()
				case C.Import:
					f.imp = v.Text()
				}
			} else {
				switch a.Id() {
				case caps.CodecCapnp:
					enableCodec(f, caps.CodecCapnp)
					g_imported["io"] = true
					g_imported[GO_CAPNP_IMPORT] = true
				case caps.CodecJson:
					enableCodec(f, caps.CodecJson)
				case caps.CodecMsgp:
					enableCodec(f, caps.CodecMsgp)
				}
			}
		}

		for _, nn := range f.NestedNodes().ToArray() {
			if ni := g_nodes[nn.Id()]; ni != nil {
				ni.resolveName("", nn.Name(), f)
			}
		}
	}

	for _, reqf := range req.RequestedFiles().ToArray() {
		x := bam.NewExtractor()
		x.FieldPrefix = "   "
		x.FieldSuffix = "\n"

		f := findNode(reqf.Id())
		buf := bytes.Buffer{}
		g_segment = C.NewBuffer([]byte{})

		defineConstNodes(&buf, f.nodes)

		for _, n := range f.nodes {
			switch n.Which() {
			case caps.NODE_ANNOTATION:
				n.defineAnnotation(&buf)
			case caps.NODE_ENUM:
				n.defineEnum(&buf, x)
			case caps.NODE_STRUCT:
				if !n.Struct().IsGroup() {
					n.defineStructTypes(&buf, nil, x)
					// n.defineStructEnums(&buf)
				}
			}
		}

		// Write translation functions
		if _, found := f.codecs[caps.CodecCapnp]; found {
			_, err = x.WriteToTranslators(&buf)
			assert(err == nil, "%v\n", err)
		}

		assert(f.pkg != "", "missing package annotation for %s", reqf.Filename())
		x.PkgName = f.pkg

		if dirPath, _ := filepath.Split(reqf.Filename()); dirPath != "" {
			err := os.MkdirAll(dirPath, os.ModePerm)
			assert(err == nil, "%v\n", err)
			x.OutDir = dirPath
		}

		// Create output file
		filename := strings.TrimSuffix(reqf.Filename(), ".capnp")

		file, err := os.Create(filename + ".go")
		assert(err == nil, "%v\n", err)

		// Write package
		fmt.Fprintf(file, "package %s\n\n", f.pkg)
		fmt.Fprintf(file, "// AUTO GENERATED - DO NOT EDIT\n\n")

		// Write imports
		f.writeImports(file)

		// Format sources
		clean, err := format.Source(buf.Bytes())
		assert(err == nil, "%v\n", err)
		file.Write(clean)

		defer file.Close()
	}
}

func enableCodec(n *node, codec uint64) {
	n.codecs[codec] = true
	for _, nst := range n.NestedNodes().ToArray() {
		nn := findNode(nst.Id())
		enableCodec(nn, codec)
	}
}
