package pterm_test

import (
	"testing"

	"github.com/MarvinJWendt/testza"

	"github.com/pterm/pterm"
)

func TestInteractiveSecretInputPrinter_WithMessage(t *testing.T) {
	p := pterm.InteractiveSecretInputPrinterWithDefaultMessage.WithMessage("Enter Secret: ")
	testza.AssertEqual(t, p.Message, "Enter Secret: ")
}
