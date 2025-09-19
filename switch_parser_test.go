package switchparser

import (
	"fmt"
	"testing"

	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func TestPlugin(t *testing.T) {
	input := `
	switch (os) {
		case 'linux':
			console.log('Linux system')
			console.log('!')
			break
		case 'macos':
			console.log('macOS system')
			break
		case 'windows':
			console.log('Windows system')
			break
		default:
			console.log('not sure')
	}`
	lb := lexer.NewBuilder()
	p := parser.NewBuilder(lb).Install(Plugin).Build(input)
	program, err := p.ParseProgram()
	if err != nil {
		t.Errorf("ParseProgram() error: %q", err)
	}
	fmt.Println(program)
}
