//go:build js && wasm

package web

import (
	"fmt"
	"strconv"
	"syscall/js"

	"github.com/rahul1534/PassGen/internal/config"
	"github.com/rahul1534/PassGen/internal/generator"
	"github.com/rahul1534/PassGen/internal/random"
)

type mode int

const (
	modeRandom mode = iota
	modeStrong
	modePassphrase
	modePIN
)

type theme int

const (
	themeSystem theme = iota
	themeLight
	themeDark
)

// App manages UI state and browser integration.
type App struct {
	doc        js.Value
	rng        random.Source
	mode       mode
	theme      theme
	password   generator.PasswordOptions
	passphrase generator.PassphraseOptions
	pin        generator.PINOptions
	output     string
	errorMsg   string
	copyReset  js.Func
	callbacks  []js.Func
}

// NewApp creates the application and binds DOM events.
func NewApp() (*App, error) {
	src, err := random.NewCryptoSource()
	if err != nil {
		return nil, err
	}

	app := &App{
		doc:        js.Global().Get("document"),
		rng:        src,
		mode:       modeRandom,
		password:   generator.DefaultPasswordOptions(),
		passphrase: generator.DefaultPassphraseOptions(),
		pin:        generator.DefaultPINOptions(),
		theme:      themeSystem,
	}

	app.bindEvents()
	app.applyTheme()
	app.syncControlsFromState()
	app.generate()
	return app, nil
}

func (a *App) bindEvents() {
	a.onClick("btn-generate", a.generate)
	a.onClick("btn-copy", a.copyPassword)
	a.onClick("btn-reset", a.resetDefaults)
	a.onClick("btn-advanced-toggle", a.toggleAdvanced)

	a.onChange("mode-random", func() { a.setMode(modeRandom) })
	a.onChange("mode-strong", func() { a.setMode(modeStrong) })
	a.onChange("mode-passphrase", func() { a.setMode(modePassphrase) })
	a.onChange("mode-pin", func() { a.setMode(modePIN) })

	a.onInput("input-length", a.readPasswordControls)
	a.onInput("input-min-upper", a.readPasswordControls)
	a.onInput("input-min-lower", a.readPasswordControls)
	a.onInput("input-min-numbers", a.readPasswordControls)
	a.onInput("input-min-symbols", a.readPasswordControls)
	a.onInput("input-excluded", a.readPasswordControls)
	a.onInput("input-words", a.readPassphraseControls)
	a.onInput("input-separator", a.readPassphraseControls)
	a.onInput("input-pin-length", a.readPINControls)

	for _, id := range []string{
		"chk-upper", "chk-lower", "chk-numbers", "chk-symbols",
		"chk-exclude-ambiguous", "chk-prevent-repeated",
		"chk-capitalize", "chk-add-number", "chk-add-symbol",
		"chk-pin-repeated", "chk-pin-avoid-patterns",
	} {
		a.onChange(id, a.readAllControls)
	}

	a.onChange("theme-select", a.readTheme)

	a.retain(js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		e := args[0]
		key := e.Get("key").String()
		if key == "Enter" && (e.Get("metaKey").Bool() || e.Get("ctrlKey").Bool()) {
			a.generate()
		}
		return nil
	}))
	a.doc.Call("addEventListener", "keydown", a.callbacks[len(a.callbacks)-1])
}

func (a *App) retain(fn js.Func) js.Func {
	a.callbacks = append(a.callbacks, fn)
	return fn
}

func (a *App) onClick(id string, fn func()) {
	el := a.doc.Call("getElementById", id)
	if el.IsNull() {
		return
	}
	handler := a.retain(js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		fn()
		return nil
	}))
	el.Call("addEventListener", "click", handler)
}

func (a *App) onChange(id string, fn func()) {
	el := a.doc.Call("getElementById", id)
	if el.IsNull() {
		return
	}
	handler := a.retain(js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		fn()
		return nil
	}))
	el.Call("addEventListener", "change", handler)
}

func (a *App) onInput(id string, fn func()) {
	el := a.doc.Call("getElementById", id)
	if el.IsNull() {
		return
	}
	handler := a.retain(js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		fn()
		return nil
	}))
	el.Call("addEventListener", "input", handler)
}

func (a *App) setMode(m mode) {
	a.mode = m
	a.showPanel()
	a.generate()
}

func (a *App) showPanel() {
	panels := map[mode]string{
		modeRandom:     "panel-random",
		modeStrong:     "panel-strong",
		modePassphrase: "panel-passphrase",
		modePIN:        "panel-pin",
	}
	for m, id := range panels {
		el := a.doc.Call("getElementById", id)
		if el.IsNull() {
			continue
		}
		if m == a.mode {
			el.Get("classList").Call("remove", "hidden")
		} else {
			el.Get("classList").Call("add", "hidden")
		}
	}

	modeIDs := map[mode]string{
		modeRandom:     "mode-random",
		modeStrong:     "panel-strong",
		modePassphrase: "mode-passphrase",
		modePIN:        "mode-pin",
	}
	for m, id := range modeIDs {
		el := a.doc.Call("getElementById", id)
		if !el.IsNull() {
			el.Set("checked", m == a.mode)
		}
	}
}

func (a *App) readAllControls() {
	a.readPasswordControls()
	a.readPassphraseControls()
	a.readPINControls()
}

func (a *App) readPasswordControls() {
	a.password.Length = a.intValue("input-length", a.password.Length)
	a.password.Uppercase = a.checked("chk-upper")
	a.password.Lowercase = a.checked("chk-lower")
	a.password.Numbers = a.checked("chk-numbers")
	a.password.Symbols = a.checked("chk-symbols")
	a.password.MinUppercase = a.intValue("input-min-upper", a.password.MinUppercase)
	a.password.MinLowercase = a.intValue("input-min-lower", a.password.MinLowercase)
	a.password.MinNumbers = a.intValue("input-min-numbers", a.password.MinNumbers)
	a.password.MinSymbols = a.intValue("input-min-symbols", a.password.MinSymbols)
	a.password.ExcludeAmbiguous = a.checked("chk-exclude-ambiguous")
	a.password.PreventRepeated = a.checked("chk-prevent-repeated")
	a.password.ExcludedCharacters = generator.NormalizeExcludedChars(a.stringValue("input-excluded"))
}

func (a *App) readPassphraseControls() {
	a.passphrase.Words = a.intValue("input-words", a.passphrase.Words)
	a.passphrase.Separator = a.stringValue("input-separator")
	if a.passphrase.Separator == "" {
		a.passphrase.Separator = "-"
	}
	a.passphrase.Capitalize = a.checked("chk-capitalize")
	a.passphrase.AddNumber = a.checked("chk-add-number")
	a.passphrase.AddSymbol = a.checked("chk-add-symbol")
}

func (a *App) readPINControls() {
	a.pin.Length = a.intValue("input-pin-length", a.pin.Length)
	a.pin.AllowRepeatedDigits = a.checked("chk-pin-repeated")
	a.pin.AvoidAmbiguousPatterns = a.checked("chk-pin-avoid-patterns")
}

func (a *App) readTheme() {
	val := a.stringValue("theme-select")
	switch val {
	case "light":
		a.theme = themeLight
	case "dark":
		a.theme = themeDark
	default:
		a.theme = themeSystem
	}
	a.applyTheme()
}

func (a *App) syncControlsFromState() {
	a.setIntValue("input-length", a.password.Length)
	a.setChecked("chk-upper", a.password.Uppercase)
	a.setChecked("chk-lower", a.password.Lowercase)
	a.setChecked("chk-numbers", a.password.Numbers)
	a.setChecked("chk-symbols", a.password.Symbols)
	a.setIntValue("input-min-upper", a.password.MinUppercase)
	a.setIntValue("input-min-lower", a.password.MinLowercase)
	a.setIntValue("input-min-numbers", a.password.MinNumbers)
	a.setIntValue("input-min-symbols", a.password.MinSymbols)
	a.setChecked("chk-exclude-ambiguous", a.password.ExcludeAmbiguous)
	a.setChecked("chk-prevent-repeated", a.password.PreventRepeated)
	a.setValue("input-excluded", a.password.ExcludedCharacters)

	a.setIntValue("input-words", a.passphrase.Words)
	a.setValue("input-separator", a.passphrase.Separator)
	a.setChecked("chk-capitalize", a.passphrase.Capitalize)
	a.setChecked("chk-add-number", a.passphrase.AddNumber)
	a.setChecked("chk-add-symbol", a.passphrase.AddSymbol)

	a.setIntValue("input-pin-length", a.pin.Length)
	a.setChecked("chk-pin-repeated", a.pin.AllowRepeatedDigits)
	a.setChecked("chk-pin-avoid-patterns", a.pin.AvoidAmbiguousPatterns)

	a.showPanel()
}

func (a *App) generate() {
	a.readAllControls()
	a.errorMsg = ""

	var (
		result   string
		err      error
		strength generator.StrengthResult
	)

	switch a.mode {
	case modeRandom:
		result, err = generator.GeneratePassword(a.rng, a.strongPassword)
		if err == nil {
			strength = generator.EstimatePasswordStrength(result, a.strongPassword)
		}
	case modeStrong:
		result, err = generator.GeneratePassword(a.rng, a.password)
		if err == nil {
			strength = generator.EstimatePasswordStrength(result, a.password)
		}
	case modePassphrase:
		result, err = generator.GeneratePassphrase(a.rng, a.passphrase)
		if err == nil {
			strength = generator.EstimatePassphraseStrength(a.passphrase, generator.WordListSize())
		}
	case modePIN:
		result, err = generator.GeneratePIN(a.rng, a.pin)
		if err == nil {
			strength = generator.EstimatePINStrength(a.pin.Length, !a.pin.AllowRepeatedDigits)
		}
	}

	if err != nil {
		a.errorMsg = generator.UserMessage(err)
		a.output = ""
	} else {
		a.output = result
	}
	a.render(strength)
}

func (a *App) render(strength generator.StrengthResult) {
	output := a.doc.Call("getElementById", "password-output")
	if !output.IsNull() {
		output.Set("textContent", a.output)
	}

	errEl := a.doc.Call("getElementById", "validation-error")
	if !errEl.IsNull() {
		errEl.Set("textContent", a.errorMsg)
		if a.errorMsg != "" {
			errEl.Get("classList").Call("remove", "hidden")
		} else {
			errEl.Get("classList").Call("add", "hidden")
		}
	}

	strengthLabel := a.doc.Call("getElementById", "strength-label")
	if !strengthLabel.IsNull() {
		strengthLabel.Set("textContent", strength.Level.String())
	}

	entropyEl := a.doc.Call("getElementById", "entropy-bits")
	if !entropyEl.IsNull() {
		if a.output == "" {
			entropyEl.Set("textContent", "")
		} else {
			entropyEl.Set("textContent", fmt.Sprintf("~%.0f bits estimated entropy", strength.Entropy))
		}
	}

	strengthBar := a.doc.Call("getElementById", "strength-bar")
	if !strengthBar.IsNull() {
		strengthBar.Set("style", fmt.Sprintf("width: %d%%", strength.Level.BarWidth()))
		strengthBar.Set("data-level", strength.Level.String())
	}

	genBtn := a.doc.Call("getElementById", "btn-generate")
	if !genBtn.IsNull() {
		genBtn.Set("disabled", a.errorMsg != "" && a.output == "")
	}
}

func (a *App) copyPassword() {
	if a.output == "" {
		return
	}
	btn := a.doc.Call("getElementById", "btn-copy")
	clipboard := js.Global().Get("navigator").Get("clipboard")
	if clipboard.IsUndefined() || clipboard.IsNull() {
		a.setStatus("Unable to access clipboard. Select and copy the password manually.")
		return
	}

	promise := clipboard.Call("writeText", a.output)
	thenFn := a.retain(js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		if !btn.IsNull() {
			btn.Set("textContent", "Copied!")
		}
		if a.copyReset.Type() == js.TypeUndefined {
			a.copyReset = a.retain(js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
				if !btn.IsNull() {
					btn.Set("textContent", "Copy")
				}
				return nil
			}))
		}
		js.Global().Call("setTimeout", a.copyReset, 2000)
		return nil
	}))
	catchFn := a.retain(js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		a.setStatus("Unable to access clipboard. Select and copy the password manually.")
		return nil
	}))
	promise.Call("then", thenFn).Call("catch", catchFn)
}

func (a *App) setStatus(msg string) {
	el := a.doc.Call("getElementById", "status-message")
	if !el.IsNull() {
		el.Set("textContent", msg)
	}
}

func (a *App) resetDefaults() {
	a.mode = modeRandom
	a.password = generator.DefaultPasswordOptions()
	a.strongPassword = generator.StrongPasswordOptions()
	a.passphrase = generator.DefaultPassphraseOptions()
	a.pin = generator.DefaultPINOptions()
	a.theme = themeSystem
	a.setValue("theme-select", "system")
	a.applyTheme()
	a.syncControlsFromState()
	a.generate()
}

func (a *App) toggleAdvanced() {
	el := a.doc.Call("getElementById", "advanced-panel")
	btn := a.doc.Call("getElementById", "btn-advanced-toggle")
	if el.IsNull() {
		return
	}
	hidden := el.Get("classList").Call("contains", "hidden").Bool()
	if hidden {
		el.Get("classList").Call("remove", "hidden")
		if !btn.IsNull() {
			btn.Set("aria-expanded", "true")
			btn.Set("textContent", "Advanced Options ▲")
		}
	} else {
		el.Get("classList").Call("add", "hidden")
		if !btn.IsNull() {
			btn.Set("aria-expanded", "false")
			btn.Set("textContent", "Advanced Options ▼")
		}
	}
}

func (a *App) applyTheme() {
	root := a.doc.Get("documentElement")
	switch a.theme {
	case themeLight:
		root.Set("data-theme", "light")
	case themeDark:
		root.Set("data-theme", "dark")
	default:
		root.Set("data-theme", "system")
	}
}

func (a *App) checked(id string) bool {
	el := a.doc.Call("getElementById", id)
	if el.IsNull() {
		return false
	}
	return el.Get("checked").Bool()
}

func (a *App) setChecked(id string, value bool) {
	el := a.doc.Call("getElementById", id)
	if !el.IsNull() {
		el.Set("checked", value)
	}
}

func (a *App) stringValue(id string) string {
	el := a.doc.Call("getElementById", id)
	if el.IsNull() {
		return ""
	}
	return el.Get("value").String()
}

func (a *App) setValue(id, value string) {
	el := a.doc.Call("getElementById", id)
	if !el.IsNull() {
		el.Set("value", value)
	}
}

func (a *App) intValue(id string, fallback int) int {
	s := a.stringValue(id)
	if s == "" {
		return fallback
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return n
}

func (a *App) setIntValue(id string, value int) {
	a.setValue(id, strconv.Itoa(value))
}

// SetGitHubLink configures repository links from build constants.
func SetGitHubLink() {
	doc := js.Global().Get("document")
	for _, id := range []string{"github-link", "footer-github-link"} {
		link := doc.Call("getElementById", id)
		if !link.IsNull() {
			link.Set("href", config.GitHubURL)
		}
	}
	title := doc.Call("getElementById", "app-title")
	if !title.IsNull() {
		title.Set("textContent", config.AppName)
	}
}
