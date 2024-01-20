package main

import (
	"fmt"
	"log"

	esbuild "github.com/evanw/esbuild/pkg/api"

	"github.com/robertkrimen/otto"
)

var textEncoderPolyfill = `function TextEncoder(){}TextEncoder.prototype.encode=function(string){var octets=[];var length=string.length;var i=0;while(i<length){var codePoint=string.codePointAt(i);var c=0;var bits=0;if(codePoint<=0x0000007F){c=0;bits=0x00}else if(codePoint<=0x000007FF){c=6;bits=0xC0}else if(codePoint<=0x0000FFFF){c=12;bits=0xE0}else if(codePoint<=0x001FFFFF){c=18;bits=0xF0}octets.push(bits|(codePoint>>c));c-=6;while(c>=0){octets.push(0x80|((codePoint>>c)&0x3F));c-=6}i+=codePoint>=0x10000?2:1}return octets};function TextDecoder(){}TextDecoder.prototype.decode=function(octets){var string="";var i=0;while(i<octets.length){var octet=octets[i];var bytesNeeded=0;var codePoint=0;if(octet<=0x7F){bytesNeeded=0;codePoint=octet&0xFF}else if(octet<=0xDF){bytesNeeded=1;codePoint=octet&0x1F}else if(octet<=0xEF){bytesNeeded=2;codePoint=octet&0x0F}else if(octet<=0xF4){bytesNeeded=3;codePoint=octet&0x07}if(octets.length-i-bytesNeeded>0){var k=0;while(k<bytesNeeded){octet=octets[i+k+1];codePoint=(codePoint<<6)|(octet&0x3F);k+=1}}else{codePoint=0xFFFD;bytesNeeded=octets.length-i}string+=String.fromCodePoint(codePoint);i+=bytesNeeded+1}return string};`
var processPolyfill = `var process = {env: {NODE_ENV: "production"}};`
var consolePolyfill = `var console = {log: function(){}};`
var symbolPolyfill = `// Polyfill for Symbol
(function() {
  var root = this;
  var SymbolPolyfill = function(name) {
    return '__$' + name + '$' + Math.random().toString(36).substr(2) + '$__';
  };
  root.Symbol = SymbolPolyfill;
}());`

func main() {
	result := esbuild.Build(esbuild.BuildOptions{
		EntryPoints: []string{"./frontend/serverEntry.jsx"},
		Bundle:      true,
		Write:       false,
		Outdir:      "/",
		Format:      esbuild.FormatIIFE,
		Platform:    esbuild.PlatformNode,
		Target:      2,
		Banner: map[string]string{
			"js": textEncoderPolyfill + processPolyfill + consolePolyfill + symbolPolyfill,
		},
		Loader: map[string]esbuild.Loader{
			".jsx": esbuild.LoaderJSX,
		},
		// External:    []string{},
	})
	s := fmt.Sprintf("%s", result.OutputFiles[0].Contents)
	fmt.Println(s)
	if len(result.Errors) > 0 {
		log.Fatalf("Failed to build: %v", result.Errors)
	}
	vm := otto.New()
	value, err := vm.Run(string(result.OutputFiles[0].Contents))
	if err != nil {
		log.Fatal("Failed to evaluate app bundle:", err)
	}
	resultString, _ := value.ToString()
	fmt.Println("JavaScript Result:", resultString)
}
