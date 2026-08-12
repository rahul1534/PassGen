//go:build js && wasm

package main

import (
	"fmt"
	"syscall/js"

	"github.com/rahul1534/PassGen/internal/web"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			showFatal(fmt.Sprintf("Application error: %v", r))
		}
	}()

	web.SetGitHubLink()
	_, err := web.NewApp()
	if err != nil {
		showFatal(err.Error())
		return
	}
	select {}
}

func showFatal(message string) {
	doc := js.Global().Get("document")
	el := doc.Call("getElementById", "validation-error")
	if !el.IsNull() {
		el.Set("textContent", message)
		el.Get("classList").Call("remove", "hidden")
	}
}
