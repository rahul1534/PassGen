//go:build js && wasm

package web

import (
	"fmt"
	"strconv"
	"syscall/js"
	"unicode/utf8"

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
	doc            js.Value
	rng            random.Source
	mode           mode
	theme          theme
	password       generator.PasswordOptions
	strongPassword generator.PasswordOptions
	passphrase     generator.PassphraseOptions
	pin            generator.PINOptions
	output         string
	errorMsg       string
	copyReset      js.Func
	callbacks      []js.Func
}

// NewApp creates the application and binds DOM events.
func NewApp() (*App, error) {
	src, err := random.NewCryptoSource()
	if err != nil {
		return nil, err
	}

	app := &App{
		doc:            js.Global().Get("document"),
		rng:            src,
		mode:           modeRandom,
		password:       generator.DefaultPasswordOptions(),
		strongPassword: generator.StrongPasswordOptions(generator.DefaultPasswordOptions().Length),
		passphrase:     generator.DefaultPassphraseOptions(),
		pin:            generator.DefaultPINOptions(),
		theme:          themeSystem,
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

	// Options regenerate the result immediately, so the readout and strength
	// meter always describe the settings on screen. Number fields commit on
	// change (blur/Enter/steppers) so half-typed values don't flash errors;
	// sliders and text fields update as they move.
	a.bindLength("input-length", "input-length-range")
	a.bindLength("input-strong-length", "input-strong-length-range")
	a.bindLength("input-words", "input-words-range")
	a.bindLength("input-pin-length", "input-pin-length-range")

	for _, id := range []string{
		"input-min-upper", "input-min-lower", "input-min-numbers", "input-min-symbols",
		"select-wordlist",
	} {
		a.onChange(id, a.generate)
	}
	a.onInput("input-excluded", a.generate)
	a.onInput("input-separator", a.generate)

	for _, id := range []string{
		"chk-upper", "chk-lower", "chk-numbers", "chk-symbols",
		"chk-exclude-ambiguous", "chk-prevent-repeated",
		"chk-capitalize", "chk-add-number", "chk-add-symbol",
		"chk-pin-repeated", "chk-pin-avoid-patterns",
	} {
		a.onChange(id, a.generate)
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

// bindLength keeps a number field and its slider in step and regenerates.
func (a *App) bindLength(numID, rangeID string) {
	a.onChange(numID, func() {
		if n, err := strconv.Atoi(a.stringValue(numID)); err == nil {
			a.setValue(rangeID, strconv.Itoa(n)) // the browser clamps to the slider's range
		}
		a.generate()
	})
	a.onInput(rangeID, func() {
		a.setValue(numID, a.stringValue(rangeID))
		a.generate()
	})
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
		modeStrong:     "mode-strong",
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
	a.readStrongPasswordControls()
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

func (a *App) readStrongPasswordControls() {
	a.strongPassword.Length = a.intValue("input-strong-length", a.strongPassword.Length)
}

func (a *App) readPassphraseControls() {
	a.passphrase.Words = a.intValue("input-words", a.passphrase.Words)
	a.passphrase.WordList = a.stringValue("select-wordlist")
	if a.passphrase.WordList == "" {
		a.passphrase.WordList = generator.DefaultPassphraseWordList
	}
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
	a.setIntValue("input-strong-length", a.strongPassword.Length)
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
	a.setValue("select-wordlist", a.passphrase.WordList)
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
		result, err = generator.GeneratePassword(a.rng, a.password)
		if err == nil {
			strength = generator.EstimatePasswordStrength(result, a.password)
		}
	case modeStrong:
		result, err = generator.GeneratePassword(a.rng, a.strongPassword)
		if err == nil {
			strength = generator.EstimatePasswordStrength(result, a.strongPassword)
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
		a.renderSecret(output, a.output)
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

	// With no result (validation error) show an empty meter rather than a
	// misleading "Very Weak".
	level, levelWidth, bits := "—", 0, ""
	if a.output != "" {
		level = strength.Level.String()
		levelWidth = strength.Level.BarWidth()
		bits = fmt.Sprintf("~%.0f bits estimated entropy", strength.Entropy)
	}

	if el := a.doc.Call("getElementById", "strength-label"); !el.IsNull() {
		el.Set("textContent", level)
	}
	if el := a.doc.Call("getElementById", "entropy-bits"); !el.IsNull() {
		el.Set("textContent", bits)
	}
	if bar := a.doc.Call("getElementById", "strength-bar"); !bar.IsNull() {
		bar.Set("style", fmt.Sprintf("width: %d%%", levelWidth))
		dataLevel := level
		if a.output == "" {
			dataLevel = ""
		}
		bar.Call("setAttribute", "data-level", dataLevel)
	}
	if track := a.doc.Call("getElementById", "strength-track"); !track.IsNull() {
		track.Call("setAttribute", "aria-valuenow", strconv.Itoa(levelWidth))
		track.Call("setAttribute", "aria-valuetext", level)
	}

	genBtn := a.doc.Call("getElementById", "btn-generate")
	if !genBtn.IsNull() {
		genBtn.Set("disabled", false)
	}
	if copyBtn := a.doc.Call("getElementById", "btn-copy"); !copyBtn.IsNull() {
		copyBtn.Set("disabled", a.output == "")
	}
}

// renderSecret shows s in el as runs of letters, digits and symbols so digits
// and symbols can be coloured. It only ever sets textContent on freshly
// created spans, never HTML, and el.textContent stays exactly s.
func (a *App) renderSecret(el js.Value, s string) {
	el.Call("replaceChildren")
	density := ""
	if utf8.RuneCountInString(s) > 48 {
		density = "long"
	}
	if density == "" {
		el.Call("removeAttribute", "data-density")
	} else {
		el.Call("setAttribute", "data-density", density)
	}
	for _, r := range splitRuns(s) {
		span := a.doc.Call("createElement", "span")
		if cls := r.Kind.class(); cls != "" {
			span.Set("className", cls)
		}
		span.Set("textContent", r.Text)
		el.Call("appendChild", span)
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
			btn.Call("setAttribute", "data-state", "copied")
		}
		a.announce("Copied to clipboard.")
		if a.copyReset.Type() == js.TypeUndefined {
			a.copyReset = a.retain(js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
				if !btn.IsNull() {
					btn.Set("textContent", "Copy")
					btn.Call("removeAttribute", "data-state")
				}
				a.announce("")
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

// announce updates a visually hidden live region for screen-reader users.
func (a *App) announce(msg string) {
	if el := a.doc.Call("getElementById", "copy-announcer"); !el.IsNull() {
		el.Set("textContent", msg)
	}
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
	a.strongPassword = generator.StrongPasswordOptions(generator.DefaultPasswordOptions().Length)
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
			btn.Call("setAttribute", "aria-expanded", "true")
		}
	} else {
		el.Get("classList").Call("add", "hidden")
		if !btn.IsNull() {
			btn.Call("setAttribute", "aria-expanded", "false")
		}
	}
}

func (a *App) applyTheme() {
	root := a.doc.Get("documentElement")
	switch a.theme {
	case themeLight:
		root.Call("setAttribute", "data-theme", "light")
	case themeDark:
		root.Call("setAttribute", "data-theme", "dark")
	default:
		root.Call("setAttribute", "data-theme", "system")
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
	a.setValue(id+"-range", strconv.Itoa(value)) // no-op for fields without a slider
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
