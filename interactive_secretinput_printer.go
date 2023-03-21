package pterm

import (
	"atomicgo.dev/cursor"
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"

	"github.com/pterm/pterm/internal"
)

// InteractiveSecretInputPrinterWithDefaultMessage is the default
// InteractiveSecretInputPrinter with the default message.
var InteractiveSecretInputPrinterWithDefaultMessage = InteractiveSecretInputPrinter{
	Message: "Enter Secret: ",
}

// InteractiveSecretInputPrinter is a prompt for the user to enter a secret
// value (e.g. password or API token). The entered value will be masked with
// "***" and will not be shown on screen.
type InteractiveSecretInputPrinter struct {
	Message string // the message to display to the user

	input     string // the entered value
	cursorPos int    // the cursor position
	text      string // the text to display
}

// WithMessage sets the message to display to the user.
func (p InteractiveSecretInputPrinter) WithMessage(message string) *InteractiveSecretInputPrinter {
	p.Message = message
	return &p
}

// Show displays the prompt to the user, prints asterisks(*) for each character
// entered and returns the entered value.
func (p InteractiveSecretInputPrinter) Show() (string, error) {
	// Print message
	Println(p.Message)

	cancel, exit := internal.NewCancelationSignal()
	defer exit()

	areaText := ""

	p.text = areaText

	area, err := DefaultArea.Start(areaText)
	defer area.Stop()
	if err != nil {
		return "", err
	}

	cursor.Up(1)
	cursor.StartOfLine()

	err = keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.Enter:
			return true, nil
		case keys.RuneKey:
			p.input = string(append([]rune(p.input)[:len([]rune(p.input))+p.cursorPos], append([]rune(key.String()), []rune(p.input)[len([]rune(p.input))+p.cursorPos:]...)...))
			p.text = string(append([]rune(p.text)[:len([]rune(p.text))+p.cursorPos], append([]rune("*"), []rune(p.text)[len([]rune(p.text))+p.cursorPos:]...)...))
		case keys.Space:
			p.input = string(append([]rune(p.input)[:len([]rune(p.input))+p.cursorPos], append([]rune(" "), []rune(p.input)[len([]rune(p.input))+p.cursorPos:]...)...))
			p.text = string(append([]rune(p.text)[:len([]rune(p.text))+p.cursorPos], append([]rune("*"), []rune(p.text)[len([]rune(p.text))+p.cursorPos:]...)...))
		case keys.Backspace:
			if len([]rune(p.input))+p.cursorPos > 0 {
				p.input = string(append([]rune(p.input)[:len([]rune(p.input))+p.cursorPos-1], []rune(p.input)[len([]rune(p.input))+p.cursorPos:]...))
				p.text = string(append([]rune(p.text)[:len([]rune(p.text))+p.cursorPos-1], []rune(p.text)[len([]rune(p.text))+p.cursorPos:]...))
			}
		case keys.Delete:
			if len([]rune(p.input))+p.cursorPos < len([]rune(p.input)) {
				p.input = string(append([]rune(p.input)[:len([]rune(p.input))+p.cursorPos], []rune(p.input)[len([]rune(p.input))+p.cursorPos+1:]...))
				p.text = string(append([]rune(p.text)[:len([]rune(p.text))+p.cursorPos], []rune(p.text)[len([]rune(p.text))+p.cursorPos+1:]...))
			}
		case keys.CtrlC:
			cancel()
			return true, nil
		}

		if internal.GetStringMaxWidth(p.input) > 0 {
			switch key.Code {
			case keys.Right:
				if p.cursorPos < 0 {
					p.cursorPos++
				} else if p.cursorPos < len(p.input)-1 {
					p.cursorPos = -internal.GetStringMaxWidth(p.input)
				}
			case keys.Left:
				if p.cursorPos+internal.GetStringMaxWidth(p.input) > 0 {
					p.cursorPos--
				} else if p.cursorPos < 0 {
					p.cursorPos = 0
				}
			}
		}

		p.updateArea(area)

		return false, nil
	})

	if err != nil {
		return "", err
	}

	// Add new line
	Println()

	return p.input, nil
}

func (p InteractiveSecretInputPrinter) updateArea(area *AreaPrinter) string {
	areaText := p.text

	if p.cursorPos+internal.GetStringMaxWidth(p.input) < 1 {
		p.cursorPos = -internal.GetStringMaxWidth(p.input)
	}

	cursor.StartOfLine()
	area.Update(areaText)
	if len(p.input) > 0 {
		cursor.Up(1)
	}
	cursor.StartOfLine()

	cursor.Right(internal.GetStringMaxWidth(areaText) + p.cursorPos)

	return areaText
}
