package go2way

import (
	"bytes"
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/mattn/go2way/internal/scanner"
	"github.com/mattn/kinako/vm"
)

// NamedSign is behavior of named value for SQL.
type NamedSign int

const (
	// Question replace variable name into ?.
	Question NamedSign = iota
	// Dollar append $ for variable name.
	Dollar
	// Colon also append.
	Colon
)

var (
	reVar = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*$`)
)

func toInt64(v reflect.Value) int64 {
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		return int64(v.Float())
	case reflect.Int, reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.String:
		s := v.String()
		var i int64
		var err error
		if strings.HasPrefix(s, "0x") {
			i, err = strconv.ParseInt(s, 16, 64)
		} else {
			i, err = strconv.ParseInt(s, 10, 64)
		}
		if err == nil {
			return int64(i)
		}
	}
	return 0
}
func toBool(v reflect.Value) bool {
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		return v.Float() != 0.0
	case reflect.Int, reflect.Int32, reflect.Int64:
		return v.Int() != 0
	case reflect.Bool:
		return v.Bool()
	case reflect.String:
		if v.String() == "true" {
			return true
		}
		if toInt64(v) != 0 {
			return true
		}
	}
	return false
}

// Eval evaluate SQL with condition contained env. And replace variable names with sign. This return replaced SQL.
func Eval(sql string, env map[string]interface{}, sign NamedSign) (string, error) {
	venv := vm.NewEnv()
	for k, v := range env {
		venv.Define(k, v)
	}

	scan := scanner.NewScanner(strings.NewReader(sql))

	output := true
	evar := false

	buf := new(bytes.Buffer)
	for scan.Scan() {
		text := scan.Text()
		if scan.Token() != scanner.COMMENT {
			if output && evar == false {
				buf.WriteString(text)
			}
			evar = false
			continue
		}
		expr := strings.TrimSpace(text[2 : len(text)-2])
		if strings.HasPrefix(expr, "IF ") {
			expr = expr[3:]
			r, err := venv.Execute(expr)
			if err != nil {
				return "", err
			}
			if !toBool(r) {
				output = false
			}
		} else if expr == "ELSE" {
			output = !output
		} else if expr == "END" {
			output = true
		} else if reVar.MatchString(expr) {
			switch sign {
			case Question:
				buf.WriteString("?")
			case Dollar:
				buf.WriteString("$" + expr)
			case Colon:
				buf.WriteString(":" + expr)
			}
			evar = true
		} else {
			return "", errors.New("Invalid token: " + text)
		}
	}
	return buf.String(), scan.Err()
}
