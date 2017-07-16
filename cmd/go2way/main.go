package main

import (
	"encoding/json"
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

func main() {
	var file string
	var e env
	flag.StringVar(&file, "f", "", "SQL file")
	flag.Var(&e, "e", "envs")
	flag.Parse()

	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	sql, err := go2way.Eval(string(b), e, go2way.Question)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(sql)
}
