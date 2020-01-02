package main

import (
	"fmt"

	"github.com/joshkrueger/cmdmux"
)

func main() {
	router := cmdmux.NewMux()

	router.Handle("ping", func(a cmdmux.Args) error {
		fmt.Printf("pong\n")
		return nil
	})

	router.Handle("cmd :arg one", func(a cmdmux.Args) error {
		fmt.Printf("Command One: Arg = %s\n", a.Get("arg"))
		return nil
	})

	router.Handle("cmd :arg two", func(a cmdmux.Args) error {
		fmt.Printf("Command Two: Arg = %s\n", a.Get("arg"))
		return nil
	})

	router.Handle("cmd :arg two :more", func(a cmdmux.Args) error {
		fmt.Printf("Command Two: Arg = %s | More = %s\n", a.Get("arg"), a.Get("more"))
		return nil
	})

	router.Handle("cmd :arg", func(a cmdmux.Args) error {
		fmt.Printf("Command: Arg = %s\n", a.Get("arg"))
		return nil
	})

	router.Handle("md :arg", func(a cmdmux.Args) error {
		fmt.Printf("Metadata Command: Arg = %s | Metadata[foo] = %s \n", a.Get("arg"), a.GetMeta("foo"))
		return nil
	})

	inputs := []string{
		"ping",
		"cmd foo one",
		"cmd foo two",
		"cmd bar two baz",
		"cmd qux",
		"cmd does not exist",
	}

	for _, i := range inputs {
		fmt.Printf("Input: %s\n * ", i)
		err := router.Execute(i)
		if err != nil {
			fmt.Println("ERROR:", err)
		}
	}

	router.ExecuteWithMetadata(cmdmux.Metadata{
		"foo": "baz",
	}, "md bar")
}
