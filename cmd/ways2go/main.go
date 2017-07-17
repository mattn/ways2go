package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/mattn/ways2go"
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

type namedSign ways2go.NamedSign

func (s *namedSign) String() string {
	return ways2go.NamedSign(*s).String()
}

func (s *namedSign) Set(v string) error {
	switch v {
	case "?":
		*s = namedSign(ways2go.Question)
	case ":":
		*s = namedSign(ways2go.Colon)
	case "$":
		*s = namedSign(ways2go.Dollar)
	default:
		return errors.New("invalid named sign")
	}
	return nil
}

func main() {
	var file string
	var e env
	var ns namedSign = namedSign(ways2go.Question)
	flag.StringVar(&file, "f", "", "SQL file")
	flag.Var(&e, "e", "envs")
	flag.Var(&ns, "s", "named variable sign")
	flag.Parse()

	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	sql, err := ways2go.Eval(string(b), e, ways2go.NamedSign(ns))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(sql)
}
