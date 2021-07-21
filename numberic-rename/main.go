package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"

	// "reflect"
	"regexp"
	"strings"
	// "unicode"
)

func test() {
	var length int
	var revert bool
	var excludeAlphabet bool
	var preview bool
	flag.IntVar(&length, "l", 2, "length of prefix number")
	flag.BoolVar(&revert, "r", false, "revert number prefix")
	flag.BoolVar(&excludeAlphabet, "e", false, "exclude leading with alphabet char file names")
	flag.BoolVar(&preview, "p", false, "do NOT rename, just preview the renaming result")
	flag.Parse()

	wd, _ := os.Getwd()
	fmt.Println("Current Directory:", wd)
	format := fmt.Sprintf("%%0%dd_%%s", length)
	prefixExp := regexp.MustCompile(`^\d+_`)
	alphabetExp := regexp.MustCompile(`^[a-zA-Z0-9]+`)
	n := 0

	// alphabet := "0123456789abcdefghigklmnopqrstuvwxyz"
	for _, name := range flag.Args() {
		if excludeAlphabet && alphabetExp.MatchString(name) {
			continue
		}
		if revert {
			if prefixExp.MatchString(name) {
				newName := strings.SplitN(name, "_", 2)[1]
				fmt.Printf("Renaming \"%s\" -> \"%s\"\n", name, newName)
				if !preview {
					os.Rename(name, newName)
				}
			}
		} else {
			newName := fmt.Sprintf(format, n, name)
			n++
			fmt.Printf("Renaming \"%s\" -> \"%s\"\n", name, newName)
			if !preview {
				os.Rename(name, newName)
			}
		}
	}
}

type Dog struct {
	Name    string
	Age     int
	Contact map[string]int //name:number
}

func (d *Dog) Say(m string) {
	fmt.Println("say:", m, d.Name, d.Age)
}
func inspect(i interface{}) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	fmt.Printf("v.NumMethod(): %v\n", t.NumMethod())
	v.MethodByName("Say").Call([]reflect.Value{reflect.ValueOf("ff")})
	t = t.Elem()
	v = v.Elem()
	for i := 0; i < t.NumField(); i++ {
		tf := t.Field(i)
		fmt.Printf("tf.Type.String(): %v\n", tf.Type.String())
		vf := v.Field(i)
		switch vf.Kind() {
		case reflect.Int:
			fmt.Printf("%s -> %d\n", tf.Name, vf.Int())
		case reflect.String:
			fmt.Printf("%s -> %s\n", tf.Name, vf.String())
		case reflect.Map:
			fmt.Printf("%s -> ", tf.Name)
			iter := vf.MapRange()
			for iter.Next() {
				k := iter.Key()
				u := iter.Value()
				fmt.Printf("%s:%d,", k.String(), u.Int())
			}
			fmt.Println()
		}

	}
}

type T struct{}

func (t *T) Geeks() {
	fmt.Println("GeekforGeeks")
}

func main() {
	var t T
	val := reflect.ValueOf(&t).MethodByName("Geeks").Call([]reflect.Value{})
	fmt.Println(val)
	d := &Dog{Name: "you", Age: 15, Contact: map[string]int{"abc": 123, "def": 456}}
	d.Say("go")
	inspect(d)
}
