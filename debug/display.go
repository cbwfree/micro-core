package debug

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"runtime"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
)

type pointerInfo struct {
	prev *pointerInfo
	n    int
	addr uintptr
	pos  int
	used []int
}

// Display print the data in console
func Display(data ...interface{}) {
	display(true, data...)
}

// GetDisplayString return data print string
func GetDisplayString(data ...interface{}) string {
	return display(false, data...)
}

func display(displayed bool, data ...interface{}) string {
	var pc, file, line, ok = runtime.Caller(2)

	if !ok {
		return ""
	}

	var buf = new(bytes.Buffer)

	_, _ = fmt.Fprintf(buf, "[Debug] at %s() [%s:%d]\n", function(pc), file, line)

	_, _ = fmt.Fprintf(buf, "\n[Variables]\n")

	for i := 0; i < len(data); i += 2 {
		var output = fomateinfo(len(data[i].(string))+3, data[i+1])
		_, _ = fmt.Fprintf(buf, "%s = %s", data[i], output)
	}

	if displayed {
		log.Print(buf)
	}
	return buf.String()
}

// return data dump and format bytes
func fomateinfo(headlen int, data ...interface{}) []byte {
	var buf = new(bytes.Buffer)

	if len(data) > 1 {
		_, _ = fmt.Fprint(buf, "    ")

		_, _ = fmt.Fprint(buf, "[")

		_, _ = fmt.Fprintln(buf)
	}

	for k, v := range data {
		var buf2 = new(bytes.Buffer)
		var pointers *pointerInfo
		var interfaces = make([]reflect.Value, 0, 10)

		printKeyValue(buf2, reflect.ValueOf(v), &pointers, &interfaces, nil, true, "    ", 1)

		if k < len(data)-1 {
			_, _ = fmt.Fprint(buf2, ", ")
		}

		_, _ = fmt.Fprintln(buf2)

		buf.Write(buf2.Bytes())
	}

	if len(data) > 1 {
		_, _ = fmt.Fprintln(buf)

		_, _ = fmt.Fprint(buf, "    ")

		_, _ = fmt.Fprint(buf, "]")
	}

	return buf.Bytes()
}

// check data is golang basic type
func isSimpleType(val reflect.Value, kind reflect.Kind, pointers **pointerInfo, interfaces *[]reflect.Value) bool {
	switch kind {
	case reflect.Bool:
		return true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.Complex64, reflect.Complex128:
		return true
	case reflect.String:
		return true
	case reflect.Chan:
		return true
	case reflect.Invalid:
		return true
	case reflect.Interface:
		for _, in := range *interfaces {
			if reflect.DeepEqual(in, val) {
				return true
			}
		}
		return false
	case reflect.UnsafePointer:
		if val.IsNil() {
			return true
		}

		var elem = val.Elem()

		if isSimpleType(elem, elem.Kind(), pointers, interfaces) {
			return true
		}

		var addr = val.Elem().UnsafeAddr()

		for p := *pointers; p != nil; p = p.prev {
			if addr == p.addr {
				return true
			}
		}

		return false
	}

	return false
}

// dump value
func printKeyValue(buf *bytes.Buffer, val reflect.Value, pointers **pointerInfo, interfaces *[]reflect.Value, structFilter func(string, string) bool, formatOutput bool, indent string, level int) {
	var t = val.Kind()

	switch t {
	case reflect.Bool:
		_, _ = fmt.Fprint(buf, val.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		_, _ = fmt.Fprint(buf, val.Int())
	case reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64:
		_, _ = fmt.Fprint(buf, val.Uint())
	case reflect.Float32, reflect.Float64:
		_, _ = fmt.Fprint(buf, val.Float())
	case reflect.Complex64, reflect.Complex128:
		_, _ = fmt.Fprint(buf, val.Complex())
	case reflect.UnsafePointer:
		_, _ = fmt.Fprintf(buf, "unsafe.Pointer(0x%X)", val.Pointer())
	case reflect.Ptr:
		if val.IsNil() {
			_, _ = fmt.Fprint(buf, "nil")
			return
		}

		var addr = val.Elem().UnsafeAddr()

		for p := *pointers; p != nil; p = p.prev {
			if addr == p.addr {
				p.used = append(p.used, buf.Len())
				_, _ = fmt.Fprintf(buf, "0x%X", addr)
				return
			}
		}

		*pointers = &pointerInfo{
			prev: *pointers,
			addr: addr,
			pos:  buf.Len(),
			used: make([]int, 0),
		}

		_, _ = fmt.Fprint(buf, "&")

		printKeyValue(buf, val.Elem(), pointers, interfaces, structFilter, formatOutput, indent, level)
	case reflect.String:
		_, _ = fmt.Fprint(buf, "\"", val.String(), "\"")
	case reflect.Interface:
		var value = val.Elem()

		if !value.IsValid() {
			_, _ = fmt.Fprint(buf, "nil")
		} else {
			for _, in := range *interfaces {
				if reflect.DeepEqual(in, val) {
					_, _ = fmt.Fprint(buf, "repeat")
					return
				}
			}

			*interfaces = append(*interfaces, val)

			printKeyValue(buf, value, pointers, interfaces, structFilter, formatOutput, indent, level+1)
		}
	case reflect.Struct:
		var t = val.Type()

		_, _ = fmt.Fprint(buf, t)
		_, _ = fmt.Fprint(buf, "{")

		for i := 0; i < val.NumField(); i++ {
			if formatOutput {
				_, _ = fmt.Fprintln(buf)
			} else {
				_, _ = fmt.Fprint(buf, " ")
			}

			var name = t.Field(i).Name

			if formatOutput {
				for ind := 0; ind < level; ind++ {
					_, _ = fmt.Fprint(buf, indent)
				}
			}

			_, _ = fmt.Fprint(buf, name)
			_, _ = fmt.Fprint(buf, ": ")

			if structFilter != nil && structFilter(t.String(), name) {
				_, _ = fmt.Fprint(buf, "ignore")
			} else {
				printKeyValue(buf, val.Field(i), pointers, interfaces, structFilter, formatOutput, indent, level+1)
			}

			_, _ = fmt.Fprint(buf, ",")
		}

		if formatOutput {
			_, _ = fmt.Fprintln(buf)

			for ind := 0; ind < level-1; ind++ {
				_, _ = fmt.Fprint(buf, indent)
			}
		} else {
			_, _ = fmt.Fprint(buf, " ")
		}

		_, _ = fmt.Fprint(buf, "}")
	case reflect.Array, reflect.Slice:
		_, _ = fmt.Fprint(buf, val.Type())
		_, _ = fmt.Fprint(buf, "{")

		var allSimple = true

		for i := 0; i < val.Len(); i++ {
			var elem = val.Index(i)

			var isSimple = isSimpleType(elem, elem.Kind(), pointers, interfaces)

			if !isSimple {
				allSimple = false
			}

			if formatOutput && !isSimple {
				_, _ = fmt.Fprintln(buf)
			} else {
				fmt.Fprint(buf, " ")
			}

			if formatOutput && !isSimple {
				for ind := 0; ind < level; ind++ {
					_, _ = fmt.Fprint(buf, indent)
				}
			}

			printKeyValue(buf, elem, pointers, interfaces, structFilter, formatOutput, indent, level+1)

			if i != val.Len()-1 || !allSimple {
				_, _ = fmt.Fprint(buf, ",")
			}
		}

		if formatOutput && !allSimple {
			_, _ = fmt.Fprintln(buf)

			for ind := 0; ind < level-1; ind++ {
				_, _ = fmt.Fprint(buf, indent)
			}
		} else {
			_, _ = fmt.Fprint(buf, " ")
		}

		_, _ = fmt.Fprint(buf, "}")
	case reflect.Map:
		var t = val.Type()
		var keys = val.MapKeys()

		_, _ = fmt.Fprint(buf, t)
		_, _ = fmt.Fprint(buf, "{")

		var allSimple = true

		for i := 0; i < len(keys); i++ {
			var elem = val.MapIndex(keys[i])

			var isSimple = isSimpleType(elem, elem.Kind(), pointers, interfaces)

			if !isSimple {
				allSimple = false
			}

			if formatOutput && !isSimple {
				_, _ = fmt.Fprintln(buf)
			} else {
				_, _ = fmt.Fprint(buf, " ")
			}

			if formatOutput && !isSimple {
				for ind := 0; ind <= level; ind++ {
					_, _ = fmt.Fprint(buf, indent)
				}
			}

			printKeyValue(buf, keys[i], pointers, interfaces, structFilter, formatOutput, indent, level+1)
			_, _ = fmt.Fprint(buf, ": ")
			printKeyValue(buf, elem, pointers, interfaces, structFilter, formatOutput, indent, level+1)

			if i != val.Len()-1 || !allSimple {
				_, _ = fmt.Fprint(buf, ",")
			}
		}

		if formatOutput && !allSimple {
			_, _ = fmt.Fprintln(buf)

			for ind := 0; ind < level-1; ind++ {
				_, _ = fmt.Fprint(buf, indent)
			}
		} else {
			_, _ = fmt.Fprint(buf, " ")
		}

		_, _ = fmt.Fprint(buf, "}")
	case reflect.Chan:
		_, _ = fmt.Fprint(buf, val.Type())
	case reflect.Invalid:
		_, _ = fmt.Fprint(buf, "invalid")
	default:
		_, _ = fmt.Fprint(buf, "unknow")
	}
}

// PrintPointerInfo dump pointer value
func PrintPointerInfo(buf *bytes.Buffer, headlen int, pointers *pointerInfo) {
	var anyused = false
	var pointerNum = 0

	for p := pointers; p != nil; p = p.prev {
		if len(p.used) > 0 {
			anyused = true
		}
		pointerNum++
		p.n = pointerNum
	}

	if anyused {
		var pointerBufs = make([][]rune, pointerNum+1)

		for i := 0; i < len(pointerBufs); i++ {
			var pointerBuf = make([]rune, buf.Len()+headlen)

			for j := 0; j < len(pointerBuf); j++ {
				pointerBuf[j] = ' '
			}

			pointerBufs[i] = pointerBuf
		}

		for pn := 0; pn <= pointerNum; pn++ {
			for p := pointers; p != nil; p = p.prev {
				if len(p.used) > 0 && p.n >= pn {
					if pn == p.n {
						pointerBufs[pn][p.pos+headlen] = '└'

						var maxpos = 0

						for i, pos := range p.used {
							if i < len(p.used)-1 {
								pointerBufs[pn][pos+headlen] = '┴'
							} else {
								pointerBufs[pn][pos+headlen] = '┘'
							}

							maxpos = pos
						}

						for i := 0; i < maxpos-p.pos-1; i++ {
							if pointerBufs[pn][i+p.pos+headlen+1] == ' ' {
								pointerBufs[pn][i+p.pos+headlen+1] = '─'
							}
						}
					} else {
						pointerBufs[pn][p.pos+headlen] = '│'

						for _, pos := range p.used {
							if pointerBufs[pn][pos+headlen] == ' ' {
								pointerBufs[pn][pos+headlen] = '│'
							} else {
								pointerBufs[pn][pos+headlen] = '┼'
							}
						}
					}
				}
			}

			buf.WriteString(string(pointerBufs[pn]) + "\n")
		}
	}
}

// Stack get stack bytes
func Stack(skip int, indent string) []byte {
	var buf = new(bytes.Buffer)

	for i := skip; ; i++ {
		var pc, file, line, ok = runtime.Caller(i)

		if !ok {
			break
		}

		buf.WriteString(indent)

		_, _ = fmt.Fprintf(buf, "at %s() [%s:%d]\n", function(pc), file, line)
	}

	return buf.Bytes()
}

// return the name of the function containing the PC if possible,
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
