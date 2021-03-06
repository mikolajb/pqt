package pqtgo

import (
	"fmt"
	"go/types"
	"io"
	"reflect"
	"sort"
	"strings"

	"github.com/huandu/xstrings"
	"github.com/piotrkowalczuk/pqt"
)

func generateTypeBuiltin(t BuiltinType, m int32) (r string) {
	switch types.BasicKind(t) {
	case types.Bool:
		r = chooseType("bool", "*bool", "*bool", m)
	case types.Int:
		r = chooseType("int", "*int", "*int", m)
	case types.Int8:
		r = chooseType("int8", "*int8", "*int8", m)
	case types.Int16:
		r = chooseType("int16", "*int16", "*int16", m)
	case types.Int32:
		r = chooseType("int32", "*int32", "*int32", m)
	case types.Int64:
		r = chooseType("int64", "*int64", "*int64", m)
	case types.Uint:
		r = chooseType("uint", "*uint", "*uint", m)
	case types.Uint8:
		r = chooseType("uint8", "*uint8", "*uint8", m)
	case types.Uint16:
		r = chooseType("uint16", "*uint16", "*uint16", m)
	case types.Uint32:
		r = chooseType("uint32", "*uint32", "*uint32", m)
	case types.Uint64:
		r = chooseType("uint64", "*uint64", "*uint64", m)
	case types.Float32:
		r = chooseType("float32", "*float32", "*float32", m)
	case types.Float64:
		r = chooseType("float64", "*float64", "*float64", m)
	case types.Complex64:
		r = chooseType("complex64", "*complex64", "*complex64", m)
	case types.Complex128:
		r = chooseType("complex128", "*complex128", "*complex128", m)
	case types.String:
		r = chooseType("string", "*string", "*string", m)
	default:
		r = "invalid"
	}

	return
}

func generateTypeBase(t pqt.Type, m int32) string {
	switch t {
	case pqt.TypeText():
		return chooseType("string", "sql.NullString", "sql.NullString", m)
	case pqt.TypeBool():
		return chooseType("bool", "sql.NullBool", "sql.NullBool", m)
	case pqt.TypeIntegerSmall():
		return chooseType("int16", "*int16", "*int16", m)
	case pqt.TypeInteger():
		return chooseType("int32", "*int32", "*int32", m)
	case pqt.TypeIntegerBig():
		return chooseType("int64", "sql.NullInt64", "sql.NullInt64", m)
	case pqt.TypeSerial():
		return chooseType("int32", "*int32", "*int32", m)
	case pqt.TypeSerialSmall():
		return chooseType("int16", "*int16", "*int16", m)
	case pqt.TypeSerialBig():
		return chooseType("int64", "sql.NullInt64", "sql.NullInt64", m)
	case pqt.TypeTimestamp(), pqt.TypeTimestampTZ():
		return chooseType("time.Time", "pq.NullTime", "pq.NullTime", m)
	case pqt.TypeReal():
		return chooseType("float32", "*float32", "*float32", m)
	case pqt.TypeDoublePrecision():
		return chooseType("float64", "sql.NullFloat64", "sql.NullFloat64", m)
	case pqt.TypeBytea(), pqt.TypeJSON(), pqt.TypeJSONB():
		return "[]byte"
	case pqt.TypeUUID():
		return chooseType("string", "sql.NullString", "sql.NullString", m)
	default:
		gt := t.String()
		switch {
		case strings.HasPrefix(gt, "SMALLINT["), strings.HasPrefix(gt, "INTEGER["), strings.HasPrefix(gt, "BIGINT["):
			return chooseType("pq.Int64Array", "NullInt64Array", "NullInt64Array", m)
		case strings.HasPrefix(gt, "DOUBLE PRECISION["):
			return chooseType("pq.Float64Array", "NullFloat64Array", "NullFloat64Array", m)
		case strings.HasPrefix(gt, "TEXT["):
			return chooseType("pq.StringArray", "NullStringArray", "NullStringArray", m)
		case strings.HasPrefix(gt, "DECIMAL"), strings.HasPrefix(gt, "NUMERIC"):
			return chooseType("float64", "sql.NullFloat64", "sql.NullFloat64", m)
		case strings.HasPrefix(gt, "VARCHAR"), strings.HasPrefix(gt, "CHARACTER"):
			return chooseType("string", "sql.NullString", "sql.NullString", m)
		default:
			return "interface{}"
		}
	}
}

func closeBrace(w io.Writer, n int) {
	for i := 0; i < n; i++ {
		fmt.Fprintln(w, `
		}`)
	}
}

var keywords = map[string]string{
	"break":       "brk",
	"default":     "def",
	"func":        "fn",
	"interface":   "intf",
	"select":      "selec",
	"case":        "cas",
	"defer":       "defe",
	"go":          "g",
	"map":         "ma",
	"struct":      "struc",
	"chan":        "cha",
	"else":        "els",
	"goto":        "got",
	"package":     "pkg",
	"switch":      "switc",
	"const":       "cons",
	"fallthrough": "fallthroug",
	"if":          "i",
	"range":       "rang",
	"type":        "typ",
	"continue":    "cont",
	"for":         "fo",
	"import":      "impor",
	"return":      "rtn",
	"var":         "va",
}

type structField struct {
	Name     string
	Type     string
	Tags     reflect.StructTag
	ReadOnly bool
}

func or(s1, s2 string) string {
	if s1 == "" {
		return s2
	}
	return s1
}

func columnForeignName(c *pqt.Column) string {
	return c.Table.Name + "_" + c.Name
}

func sortedColumns(columns []*pqt.Column) []string {
	tmp := make([]string, 0, len(columns))
	for _, c := range columns {
		tmp = append(tmp, c.Name)
	}
	sort.Strings(tmp)

	return tmp
}

func snake(s string, private bool, acronyms map[string]string) string {
	var parts []string
	parts1 := strings.Split(s, "_")
	for _, p1 := range parts1 {
		parts2 := strings.Split(p1, "/")
		for _, p2 := range parts2 {
			parts3 := strings.Split(p2, "-")
			parts = append(parts, parts3...)
		}
	}

	for i, part := range parts {
		if !private || i > 0 {
			if formatted, ok := acronyms[part]; ok {
				parts[i] = formatted

				continue
			}
		}

		parts[i] = xstrings.FirstRuneToUpper(part)
	}

	if private {
		parts[0] = xstrings.FirstRuneToLower(parts[0])
	}
	return strings.Join(parts, "")
}

func chooseType(tm, to, tc string, mode int32) string {
	switch mode {
	case ModeCriteria:
		return tc
	case ModeMandatory:
		return tm
	case ModeOptional:
		return to
	case ModeDefault:
		return to
	default:
		panic("unknown mode")
	}
}

func generateCustomType(t CustomType, m int32) string {
	goType := func(tp reflect.Type) string {
		if tp == nil {
			return "<nil>"
		}
		return tp.String()
	}
	return chooseType(
		goType(t.mandatoryTypeOf),
		goType(t.optionalTypeOf),
		goType(t.criteriaTypeOf),
		m,
	)
}

func columnMode(c *pqt.Column, m int32) int32 {
	switch m {
	case ModeCriteria:
	case ModeMandatory:
	case ModeOptional:
	default:
		if c.NotNull || c.PrimaryKey {
			m = ModeMandatory
		}
	}
	return m
}
