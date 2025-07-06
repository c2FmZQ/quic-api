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

func main() {
	typeToName := make(map[string]string)
	for _, t := range types {
		v := reflect.ValueOf(t.typ)
		typeToName[v.Type().String()] = t.name
	}
	type arg struct {
		varName     string
		tmpName     string
		origName    string
		typeName    string
		needWrapper bool
	}
	getArgs := func(get func(int) reflect.Type, num int, skipFirst bool) []arg {
		out := make([]arg, 0, num)
		for i := 0; i < num; i++ {
			if i == 0 && skipFirst {
				continue
			}
			var a arg
			a.origName = get(i).String()
			if a.origName == "[]uint8" {
				a.origName = "[]byte"
			}
			a.typeName = a.origName
			if n, ok := typeToName[a.origName]; ok {
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
			a.tmpName = a.varName + "Internal"
			out = append(out, a)
		}
		return out
	}
	joinArgs := func(args []arg, includeVar bool) string {
		parts := make([]string, 0, len(args))
		for _, a := range args {
			if includeVar {
				parts = append(parts, fmt.Sprintf("%s %s", a.varName, a.typeName))
			} else {
				parts = append(parts, a.typeName)
			}
		}
		return strings.Join(parts, ", ")
	}
	joinVars := func(args []arg) string {
		parts := make([]string, 0, len(args))
		for _, a := range args {
			if a.needWrapper {
				parts = append(parts, fmt.Sprintf("%s.(*%sWrapper).Base", a.varName, a.typeName))
			} else {
				parts = append(parts, a.varName)
			}
		}
		return strings.Join(parts, ", ")
	}
	joinTmp := func(args []arg) string {
		parts := make([]string, 0, len(args))
		for _, a := range args {
			if a.needWrapper {
				parts = append(parts, a.tmpName)
			} else {
				parts = append(parts, a.varName)
			}
		}
		return strings.Join(parts, ", ")
	}
	writeFunc := func(receiver, funcName string, rt reflect.Type) {
		inArgs := getArgs(rt.In, rt.NumIn(), receiver != "")
		outArgs := getArgs(rt.Out, rt.NumOut(), false)
		var hasWrappers bool
		for _, oa := range outArgs {
			if oa.needWrapper {
				hasWrappers = true
			}
		}
		if receiver != "" {
			fmt.Printf("func (w *%s) %s(%s)", receiver, funcName, joinArgs(inArgs, true))
		} else {
			fmt.Printf("func %s(%s)", funcName, joinArgs(inArgs, true))
		}
		fmt.Printf(" (%s)", joinArgs(outArgs, hasWrappers))
		fmt.Printf(" {\n")
		for _, oa := range outArgs {
			if oa.needWrapper {
				fmt.Printf("\tvar %s %s\n", oa.tmpName, oa.origName)
			}
		}
		prefix := "w.Base"
		if receiver == "" {
			prefix = "quic"
		}
		if len(outArgs) == 0 {
			fmt.Printf("\t%s.%s(%s)\n", prefix, funcName, joinVars(inArgs))
		} else if hasWrappers {
			fmt.Printf("\t%s = %s.%s(%s)\n", joinTmp(outArgs), prefix, funcName, joinVars(inArgs))
		} else {
			fmt.Printf("\treturn %s.%s(%s)\n", prefix, funcName, joinVars(inArgs))
		}
		for _, oa := range outArgs {
			if oa.needWrapper {
				fmt.Printf("\tif %s != nil {\n\t\t%s = &%sWrapper{Base:%s}\n\t}\n", oa.tmpName, oa.varName, oa.typeName, oa.tmpName)
			}
		}
		if len(outArgs) > 0 && hasWrappers {
			fmt.Printf("\treturn\n")
		}
		fmt.Printf("}\n\n")
	}

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
			inArgs := getArgs(m.Type.In, m.Type.NumIn(), false)
			outArgs := getArgs(m.Type.Out, m.Type.NumOut(), false)
			fmt.Printf("\t%s(%s) (%s)\n", m.Name, joinArgs(inArgs[1:], false), joinArgs(outArgs, false))
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
			writeFunc(structName, m.Name, m.Type)
		}
	}

	for _, t := range types {
		v := reflect.ValueOf(t.typ)
		if v.Kind() != reflect.Func {
			continue
		}
		fmt.Printf("// %s is an auto-generated wrapper for [quic.%s]\n", t.name, t.name)
		writeFunc("", t.name, v.Type())
	}
}
