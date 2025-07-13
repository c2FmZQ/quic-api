package main

import (
	"fmt"
	"github.com/quic-go/quic-go"
	"reflect"
	"regexp"
	"strings"
)

var types = []struct {
	name string
	typ  any
}{
	{"Transport", &quic.Transport{}},
	{"Listener", &quic.Listener{}},
	{"EarlyListener", &quic.EarlyListener{}},
	{"Conn", &quic.Conn{}},
	{"SendStream", &quic.SendStream{}},
	{"ReceiveStream", &quic.ReceiveStream{}},
	{"Stream", &quic.Stream{}},
	{"Path", &quic.Path{}},
	{"Dial", quic.Dial},
	{"DialEarly", quic.DialEarly},
	{"DialAddr", quic.DialAddr},
	{"DialAddrEarly", quic.DialAddrEarly},
	{"Listen", quic.Listen},
	{"ListenEarly", quic.ListenEarly},
	{"ListenAddr", quic.ListenAddr},
	{"ListenAddrEarly", quic.ListenAddrEarly},
}

var varNames = []struct {
	re   *regexp.Regexp
	name string
}{
	{regexp.MustCompile(`^context.Context$`), "ctx"},
	{regexp.MustCompile(`^net.Addr$`), "addr"},
	{regexp.MustCompile(`chan `), "ch"},
	{regexp.MustCompile(`^error$`), "err"},
	{regexp.MustCompile(`^int$`), "n"},
	{regexp.MustCompile(`^Conn$`), "conn"},
	{regexp.MustCompile(`Stream$`), "stream"},
	{regexp.MustCompile(`Code$`), "code"},
	{regexp.MustCompile(`^[*]tls.Config$`), "tc"},
	{regexp.MustCompile(`^[*]quic.Config$`), "qc"},
	{regexp.MustCompile(`Listener$`), "ln"},
	{regexp.MustCompile(`ID$`), "id"},
}

func varName(t string) string {
	var out string
	for _, v := range varNames {
		if v.re.MatchString(t) {
			out = v.name
			break
		}
	}
	if out == "" {
		if t == "[]byte" {
			out = "b"
		} else {
			short := t
			if _, n, ok := strings.Cut(short, "."); ok {
				short = n
			}
			out = strings.ToLower(string(short[0]))
		}
	}
	return out
}

type argument struct {
	varName      string
	localName    string
	origTypeName string
	typeName     string
	needWrapper  bool
}

type arguments []argument

func (aa arguments) sig(includeVar, in bool) string {
	parts := make([]string, 0, len(aa))
	for _, a := range aa {
		typeName := a.typeName
		if in && a.needWrapper {
			typeName += "Unwrapper"
		}
		if includeVar {
			parts = append(parts, fmt.Sprintf("%s %s", a.varName, typeName))
		} else {
			parts = append(parts, typeName)
		}
	}
	return strings.Join(parts, ", ")
}

func (aa arguments) callParams() string {
	parts := make([]string, 0, len(aa))
	for _, a := range aa {
		if a.needWrapper {
			parts = append(parts, fmt.Sprintf("%s.Unwrap()", a.varName))
		} else {
			parts = append(parts, a.varName)
		}
	}
	return strings.Join(parts, ", ")
}

func (aa arguments) localVars() string {
	parts := make([]string, 0, len(aa))
	for _, a := range aa {
		if a.needWrapper {
			parts = append(parts, a.localName)
		} else {
			parts = append(parts, a.varName)
		}
	}
	return strings.Join(parts, ", ")
}

type builder struct {
	typeToName map[string]string
}

func (b *builder) init() {
	b.typeToName = make(map[string]string)
	for _, t := range types {
		v := reflect.ValueOf(t.typ)
		b.typeToName[v.Type().String()] = t.name
	}
}

func (b *builder) args(get func(int) reflect.Type, num int, skipFirst bool) arguments {
	out := make(arguments, 0, num)
	for i := 0; i < num; i++ {
		if i == 0 && skipFirst {
			continue
		}
		var a argument
		a.origTypeName = get(i).String()
		if a.origTypeName == "[]uint8" {
			a.origTypeName = "[]byte"
		}
		a.typeName = a.origTypeName
		if n, ok := b.typeToName[a.origTypeName]; ok {
			a.typeName = n
			a.needWrapper = true
		}
		a.varName = varName(a.typeName)
		a.localName = a.varName + "Internal"
		out = append(out, a)
	}
	names := make(map[string]bool)
	dups := make([]string, 0, len(out))
	for _, v := range out {
		if names[v.varName] {
			dups = append(dups, v.varName)
			continue
		}
		names[v.varName] = true
	}
	for _, d := range dups {
		count := 1
		for i, v := range out {
			if d == v.varName {
				v.varName = fmt.Sprintf("%s%d", d, count)
				out[i] = v
				count++
			}
		}
	}
	return out
}

func (b *builder) writeFunc(receiver, name string, rt reflect.Type) {
	inArgs := b.args(rt.In, rt.NumIn(), receiver != "")
	outArgs := b.args(rt.Out, rt.NumOut(), false)
	var hasWrappers bool
	for _, oa := range outArgs {
		if oa.needWrapper {
			hasWrappers = true
		}
	}
	if receiver != "" {
		fmt.Printf("func (w *%s) %s(%s)", receiver, name, inArgs.sig(true, true))
	} else {
		fmt.Printf("func %s(%s)", name, inArgs.sig(true, true))
	}
	fmt.Printf(" (%s)", outArgs.sig(hasWrappers, false))
	fmt.Printf(" {\n")
	for _, oa := range outArgs {
		if oa.needWrapper {
			fmt.Printf("\tvar %s %s\n", oa.localName, oa.origTypeName)
		}
	}
	prefix := "w.base"
	if receiver == "" {
		prefix = "quic"
	}
	if len(outArgs) == 0 {
		fmt.Printf("\t%s.%s(%s)\n", prefix, name, inArgs.callParams())
	} else if hasWrappers {
		fmt.Printf("\t%s = %s.%s(%s)\n", outArgs.localVars(), prefix, name, inArgs.callParams())
	} else {
		fmt.Printf("\treturn %s.%s(%s)\n", prefix, name, inArgs.callParams())
	}
	for _, oa := range outArgs {
		if oa.needWrapper {
			fmt.Printf("\t%s = Wrap%s(%s)\n", oa.varName, oa.typeName, oa.localName)
		}
	}
	if hasWrappers {
		fmt.Printf("\treturn\n")
	}
	fmt.Printf("}\n\n")
}

func main() {
	b := &builder{}
	b.init()

	for _, t := range types {
		v := reflect.ValueOf(t.typ)
		if v.Kind() != reflect.Pointer {
			continue
		}
		fmt.Printf("// %s is an auto-generated interface for [quic.%s].\n", t.name, t.name)
		fmt.Printf("// Use [Wrap%s] to convert a [*quic.%s] to a [%s].\n", t.name, t.name, t.name)
		fmt.Printf("type %s interface {\n", t.name)

		vt := v.Type()

		for i := 0; i < vt.NumMethod(); i++ {
			m := vt.Method(i)
			inArgs := b.args(m.Type.In, m.Type.NumIn(), true)
			outArgs := b.args(m.Type.Out, m.Type.NumOut(), false)
			fmt.Printf("\t%s(%s) (%s)\n", m.Name, inArgs.sig(false, true), outArgs.sig(false, false))
		}
		fmt.Printf("}\n\n")
	}

	for _, t := range types {
		v := reflect.ValueOf(t.typ)
		if v.Kind() != reflect.Pointer {
			continue
		}
		vt := v.Type()
		structName := "Wrapped" + t.name

		fmt.Printf("// %sUnwrapper is an auto-generated interface to unwrap a [*quic.%s].\n", t.name, t.name)
		fmt.Printf("// The value returned by [Wrap%s] implements this interface.\n", t.name)
		fmt.Printf("type %sUnwrapper interface {\n", t.name)
		fmt.Printf("\tUnwrap() *quic.%s\n", t.name)
		fmt.Printf("}\n\n")

		fmt.Printf("var _ %s = (*%s)(nil)\n", t.name, structName)
		fmt.Printf("var _ %sUnwrapper = (*%s)(nil)\n\n", t.name, structName)

		fmt.Printf("// Wrap%s converts a [quic.%s] to a [%s].\n", t.name, t.name, t.name)
		fmt.Printf("func Wrap%s (%s %s) *%s {\n", t.name, varName(t.name), vt, structName)
		fmt.Printf("\tif %s == nil {\n", varName(t.name))
		fmt.Printf("\t\treturn nil\n")
		fmt.Printf("\t}\n")
		fmt.Printf("\treturn &%s{base: %s}\n", structName, varName(t.name))
		fmt.Printf("}\n\n")

		fmt.Printf("// %s is an auto-generated wrapper for [quic.%s]. It implements the [%s] interface.\n", structName, t.name, t.name)
		fmt.Printf("type %s struct {\n", structName)

		fmt.Printf("\tbase %s\n", vt)
		fmt.Printf("}\n\n")

		for i := 0; i < vt.NumMethod(); i++ {
			m := vt.Method(i)
			b.writeFunc(structName, m.Name, m.Type)
		}
		fmt.Printf("// Unwrap returns the underlying [*quic.%s].\n", t.name)
		fmt.Printf("func (w *%s) Unwrap() *quic.%s {\n", structName, t.name)
		fmt.Printf("\tif w == nil {\n")
		fmt.Printf("\t\treturn nil\n")
		fmt.Printf("\t}\n")
		fmt.Printf("\treturn w.base\n")
		fmt.Printf("}\n\n")
	}

	for _, t := range types {
		v := reflect.ValueOf(t.typ)
		if v.Kind() != reflect.Func {
			continue
		}
		fmt.Printf("// %s is an auto-generated wrapper for [quic.%s]\n", t.name, t.name)
		b.writeFunc("", t.name, v.Type())
	}
}
