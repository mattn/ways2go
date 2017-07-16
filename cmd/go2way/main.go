package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/mattn/go2way"
)

type env map[string]interface{}

func (e *env) String() string {
	return fmt.Sprint(*e)
}

func (e *env) Set(v string) error {
	if *e == nil {
		*e = map[string]interface{}{}
	}
	token := strings.SplitN(v, "=", 2)
	if len(token) == 1 {
		(*e)[token[0]] = nil
	} else {
		var v interface{}
		err := json.Unmarshal([]byte(token[1]), &v)
		if err != nil {
			return err
		}
		(*e)[token[0]] = v
	}
	return nil
}

type namedSign go2way.NamedSign

func (s *namedSign) String() string {
	return go2way.NamedSign(*s).String()
}

func (s *namedSign) Set(v string) error {
	switch v {
	case "?":
		*s = namedSign(go2way.Question)
	case ":":
		*s = namedSign(go2way.Colon)
	case "$":
		*s = namedSign(go2way.Dollar)
	default:
		return errors.New("invalid named sign")
	}
	return nil
}

func main() {
	var file string
	var e env
	var ns namedSign = namedSign(go2way.Question)
	flag.StringVar(&file, "f", "", "SQL file")
	flag.Var(&e, "e", "envs")
	flag.Var(&ns, "s", "named variable sign")
	flag.Parse()

	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	sql, err := go2way.Eval(string(b), e, go2way.NamedSign(ns))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(sql)
}
