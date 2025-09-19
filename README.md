# Switch Parser Plugin for XJS

This plugin adds support for `switch` statements to the **XJS** language.

## Usage

```go
package main

import (
    "github.com/xjslang/xjs/lexer"
    "github.com/xjslang/xjs/parser"
    switchparser "path/to/switch-parser"
)

func main() {
    input := `
    switch (os) {
        case 'linux':
            console.log('Linux system')
            break
        case 'macos':
            console.log('macOS system')
            break
        default:
            console.log('Unknown OS')
    }`
    
    lb := lexer.NewBuilder()
    parser := parser.NewBuilder(lb).Install(switchparser.Plugin).Build(input)
    program, err := parser.ParseProgram()
    if err != nil {
        panic(fmt.Sprintf("ParseProgram() error: %q", err))
    }
    fmt.Println(program)
}
```
