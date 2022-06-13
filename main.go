package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
)

type Arguments map[string]string

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

var flagVal0 string
var flagVal1 string
var flagVal2 string
var flagVal3 string

func Perform(args Arguments, writer io.Writer) error {

	if args["operation"] == "" {
		err := errors.New("-operation flag has to be specified")
		fmt.Println("-operation flag has to be specified %w", err)
		//panic(err)
		return err
	}
	if args["fileName"] == "" {
		err := errors.New("-fileName flag has to be specified")
		fmt.Println("-fileName flag has to be specified %w", err)
		//panic(err)
		return err
	}
	if args["operation"] != "list" && args["operation"] != "add" && args["operation"] != "remove" && args["operation"] != "findById" {
		err := errors.New("Operation " + args["operation"] + " not allowed!")
		fmt.Println("Operation not allowed! %w", err)
		//panic(err)
		return err
	}

	if args["operation"] == "add" {

		if args["item"] == "" {
			err := errors.New("-item flag has to be specified")
			fmt.Println("-item flag has to be specified %w", err)
			//panic(err)
			return err

		}

		f, err := os.OpenFile(args["fileName"], os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			//panic(err)
			return err
		}

		f0, err := ioutil.ReadFile(args["fileName"])
		if err != nil {
			return err
		}

		var s []User
		var s0 []User

		err = json.Unmarshal(f0, &s)
		if err != nil {
			fmt.Println("Operation not allowed! %w", err)
			//panic(err)
			//	return err
		}

		err = json.Unmarshal([]byte("["+args["item"]+"]"), &s0)
		if err != nil {
			fmt.Println("Operation not allowed! %w", err)
			//panic(err)
			//	return err
		}
		fmt.Println(" typeof item = ", reflect.TypeOf(args["item"]))
		for i, v := range s {
			fmt.Println("v = ", v)

			if s[i].Id == s0[0].Id {
				err := errors.New("Item with id " + strconv.Itoa(i+1) + " already exists")
				writer.Write([]byte(err.Error()))
				return nil //err

			}
		}

		defer f.Close()

		n, err := f.Write([]byte(args["item"]))
		if err != nil {
			//panic(err)
			return err
		}
		fmt.Printf("wrote %d bytes", n)

	}

	if args["operation"] == "list" {

		f, err := ioutil.ReadFile(args["fileName"])
		if err != nil {
			return err
		}
		writer.Write(f)
	}

	if args["operation"] == "findById" {

		if args["id"] == "" {
			err := errors.New("-id flag has to be specified")
			fmt.Println("-item flag has to be specified %w", err)
			//panic(err)
			return err

		}

		f0, err := ioutil.ReadFile(args["fileName"])
		if err != nil {
			return err
		}

		var s []User

		err = json.Unmarshal(f0, &s)
		if err != nil {
			fmt.Println("Operation not allowed! %w", err)
			//panic(err)
			return err
		}

		for i, v := range s {
			fmt.Println("v = ", v)
			var s0 []byte
			if s[i].Id == args["id"] {
				s0, err = json.Marshal(s[i])
				if err != nil {
					fmt.Println("Operation not allowed! %w", err)
					return err
				}

				writer.Write(s0)
				return nil //err

			}
		}

	}

	if args["operation"] == "remove" {

		if args["id"] == "" {
			err := errors.New("-id flag has to be specified")
			fmt.Println("-id flag has to be specified %w", err)
			//panic(err)
			return err

		}

		f, err := os.OpenFile(args["fileName"], os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			//panic(err)
			return err
		}
		defer f.Close()

		f0, err := ioutil.ReadFile(args["fileName"])
		if err != nil {
			return err
		}

		var s []User
		//		var s0 []User

		err = json.Unmarshal(f0, &s)
		if err != nil {
			fmt.Println("Operation not allowed! %w", err)
		}

		c := 0
		for i := 0; i < len(s); i++ {

			if s[i].Id == args["id"] {

				c++
				copy(s[i:], s[i+1:])
				s = s[:len(s)-1]
				i--

			}
		}

		if c == 0 {
			err := errors.New("Item with id " + args["id"] + " not found")
			writer.Write([]byte(err.Error()))
			return nil

		} else {

			f.Seek(0, io.SeekStart)

			f.Truncate(0)

			var s0 []byte

			f.Write([]byte("["))
			for i, v := range s {
				fmt.Println("v = ", v)
				s0, err = json.Marshal(s[i])
				if err != nil {
					fmt.Println("Operation not allowed! %w", err)
					return err
				}
				_, err := f.Write([]byte(s0))
				if err != nil {
					//panic(err)
					return err
				}

			}
			f.Write([]byte("]"))

		}

	}

	return nil

}

func init() {

	flag.StringVar(&flagVal0, "id", "", "help message for flagname id")
	flag.StringVar(&flagVal1, "operation", "", "help message for flagname operation")
	flag.StringVar(&flagVal2, "item", "", "help message for flagname item")
	flag.StringVar(&flagVal3, "fileName", "", "help message for flagname item")

}

func parseArgs() Arguments {

	flag.Parse()
	args := Arguments{}

	args["id"] = flagVal0
	args["operation"] = flagVal1
	args["item"] = flagVal2
	args["fileName"] = flagVal3

	return args
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
