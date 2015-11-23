package demo

// AUTO GENERATED - DO NOT EDIT

const (
En = foo
In8 = int8(3)
Uin8 = uint8(49)
In16 = int16(159)
Uin16 = uint16(59)
In32 = int32(59)
Uin32 = uint32(159)
In64 = int64(159)
Uin64 = uint64(4159)
F32 = float32(3.141590)
F64 = float64(3.141590)
Em = "lisa@example.com"
Nm = "Lisa"
)
var (
Eseq = []Rfc3092Variable{foo, bar, baz}
Seq8 = []int8{3, 5, 7, 9}
Sequ8 = []uint8{3, 5, 7, 9}
Seq16 = []int16{3, 5, 7, 9}
Sequ16 = []uint16{3, 5, 7, 9}
Seq32 = []int32{3, 5, 7, 9}
Sequ32 = []uint32{3, 5, 7, 9}
Seq64 = []int64{3, 5, 7, 9}
Sequ64 = []uint64{3, 5, 7, 9}
Seqf32 = []float32{3.5, 5.34, 7, 9}
Seqf64 = []float64{3.5, 5.34, 7, 9}
Ans = []bool{true, false, true}
Let = []string{"a", "b", "c"}
Letd = []byte{[]byte{159, 152, 115, 156, 43, 83, 131, 94, 103, 32, 160, 9, 7, 171, 212, 47}, []byte{159, 152, 115, 156, 43, 83, 131, 94, 103, 32, 160, 9, 7, 171, 212, 47}, []byte{159, 152, 115, 156, 43, 83, 131, 94, 103, 32, 160, 9, 7, 171, 212, 47}}
Secret = []byte{159, 152, 115, 156, 43, 83, 131, 94, 103, 32, 160, 9, 7, 171, 212, 47}
Bob = Person{ Name: "Bob",Email: "bob@example.com",}

Liz = Person{ Name: "Lisa",Email: "lisa@example.com",}

Seqp = []Person{/* Not implemented */}
)
type Person struct {
Name string
Email string
}
type Rfc3092Variable uint16
const (
RFC3092VARIABLE_FOO Rfc3092Variable = 0
RFC3092VARIABLE_BAR Rfc3092Variable = 1
RFC3092VARIABLE_BAZ Rfc3092Variable = 2
RFC3092VARIABLE_QUX Rfc3092Variable = 3
)
func (c Rfc3092Variable) String() string {
switch c {
case RFC3092VARIABLE_FOO: return "foo"
case RFC3092VARIABLE_BAR: return "bar"
case RFC3092VARIABLE_BAZ: return "baz"
case RFC3092VARIABLE_QUX: return "qux"
default: return ""
}
}

func Rfc3092VariableFromString(c string) Rfc3092Variable {
switch c {
case "foo": return RFC3092VARIABLE_FOO
case "bar": return RFC3092VARIABLE_BAR
case "baz": return RFC3092VARIABLE_BAZ
case "qux": return RFC3092VARIABLE_QUX
default: return 0
}
}
