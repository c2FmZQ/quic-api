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

type argument struct {
	varName      string
	localName    string
	origTypeName string
	typeName     string
	needWrapper  bool
}

type arguments []argument

func (aa arguments) sig(includeVar bool) string {
	parts := make([]string, 0, len(aa))
	for _, a := range aa {
		if includeVar {
			parts = append(parts, fmt.Sprintf("%s %s", a.varName, a.typeName))
		} else {
			parts = append(parts, a.typeName)
		}
	}
	return strings.Join(parts, ", ")
}

func (aa arguments) callParams() string {
	parts := make([]string, 0, len(aa))
	for _, a := range aa {
		if a.needWrapper {
			parts = append(parts, fmt.Sprintf("%s.(*%sWrapper).Base", a.varName, a.typeName))
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
		for _, v := range varNames {
			if v.re.MatchString(a.typeName) {
				a.varName = v.name
				break
			}
		}
		if a.varName == "" {
			if a.typeName == "[]byte" {
				a.varName = "b"
			} else {
				short := a.typeName
				if _, n, ok := strings.Cut(short, "."); ok {
					short = n
				}
				a.varName = strings.ToLower(string(short[0]))
			}
		}
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
		fmt.Printf("func (w *%s) %s(%s)", receiver, name, inArgs.sig(true))
	} else {
		fmt.Printf("func %s(%s)", name, inArgs.sig(true))
	}
	fmt.Printf(" (%s)", outArgs.sig(hasWrappers))
	fmt.Printf(" {\n")
	for _, oa := range outArgs {
		if oa.needWrapper {
			fmt.Printf("\tvar %s %s\n", oa.localName, oa.origTypeName)
		}
	}
	prefix := "w.Base"
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
			fmt.Printf("\tif %s != nil {\n\t\t%s = &%sWrapper{Base:%s}\n\t}\n", oa.localName, oa.varName, oa.typeName, oa.localName)
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
		fmt.Printf("// %s is an auto-generated interface for [quic.%s]\n", t.name, t.name)
		fmt.Printf("type %s interface {\n", t.name)

		vt := v.Type()

		for i := 0; i < vt.NumMethod(); i++ {
			m := vt.Method(i)
			inArgs := b.args(m.Type.In, m.Type.NumIn(), true)
			outArgs := b.args(m.Type.Out, m.Type.NumOut(), false)
			fmt.Printf("\t%s(%s) (%s)\n", m.Name, inArgs.sig(false), outArgs.sig(false))
		}
		fmt.Printf("}\n\n")
	}

	for _, t := range types {
		v := reflect.ValueOf(t.typ)
		if v.Kind() != reflect.Pointer {
			continue
		}
		structName := t.name + "Wrapper"
		fmt.Printf("var _ %s = (*%s)(nil)\n\n", t.name, structName)
		fmt.Printf("// %s is an auto-generated wrapper for [quic.%s]\n", structName, t.name)
		fmt.Printf("type %s struct {\n", structName)

		vt := v.Type()

		fmt.Printf("\tBase %s\n", vt)
		fmt.Printf("}\n\n")

		for i := 0; i < vt.NumMethod(); i++ {
			m := vt.Method(i)
			b.writeFunc(structName, m.Name, m.Type)
		}
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
