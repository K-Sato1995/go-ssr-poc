package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	esbuild "github.com/evanw/esbuild/pkg/api"
	v8 "rogchap.com/v8go"
)

// [Yaffle/TextEncoderTextDecoder.js](https://gist.github.com/Yaffle/5458286)
var textEncoderPolyfill = `function TextEncoder(){} TextEncoder.prototype.encode=function(string){var octets=[],length=string.length,i=0;while(i<length){var codePoint=string.codePointAt(i),c=0,bits=0;codePoint<=0x7F?(c=0,bits=0x00):codePoint<=0x7FF?(c=6,bits=0xC0):codePoint<=0xFFFF?(c=12,bits=0xE0):codePoint<=0x1FFFFF&&(c=18,bits=0xF0),octets.push(bits|(codePoint>>c)),c-=6;while(c>=0){octets.push(0x80|((codePoint>>c)&0x3F)),c-=6}i+=codePoint>=0x10000?2:1}return octets};function TextDecoder(){} TextDecoder.prototype.decode=function(octets){var string="",i=0;while(i<octets.length){var octet=octets[i],bytesNeeded=0,codePoint=0;octet<=0x7F?(bytesNeeded=0,codePoint=octet&0xFF):octet<=0xDF?(bytesNeeded=1,codePoint=octet&0x1F):octet<=0xEF?(bytesNeeded=2,codePoint=octet&0x0F):octet<=0xF4&&(bytesNeeded=3,codePoint=octet&0x07),octets.length-i-bytesNeeded>0?function(){for(var k=0;k<bytesNeeded;){octet=octets[i+k+1],codePoint=(codePoint<<6)|(octet&0x3F),k+=1}}():codePoint=0xFFFD,bytesNeeded=octets.length-i,string+=String.fromCodePoint(codePoint),i+=bytesNeeded+1}return string};`
var processPolyfill = `var process = {env: {NODE_ENV: "production"}};`
var consolePolyfill = `var console = {log: function(){}};`

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>React App</title>
</head>
<body>
    <div id="app">{{.RenderedContent}}</div>
	<script type="module">
	try{
	  {{ .JS }}
	} catch (e) {
	  showError(e.stack)
	}
  </script>
	<script>window.APP_PROPS = {{.InitialProps}};</script>
</body>
</html>
`

type PageData struct {
	RenderedContent template.HTML
	InitialProps    template.JS
	JS              template.JS
}

type InitialProps struct {
	Name          string
	InitialNumber int
}

func buildBackend() string {
	result := esbuild.Build(esbuild.BuildOptions{
		EntryPoints: []string{"./frontend/serverEntry.jsx"},
		Bundle:      true,
		Write:       false,
		Outdir:      "/",
		Format:      esbuild.FormatIIFE,
		Platform:    esbuild.PlatformNode,
		Target:      esbuild.ES2015,
		Banner: map[string]string{
			"js": textEncoderPolyfill + processPolyfill + consolePolyfill,
		},
		Loader: map[string]esbuild.Loader{
			".jsx": esbuild.LoaderJSX,
		},
	})
	script := fmt.Sprintf("%s", result.OutputFiles[0].Contents)
	return script
}

func buildClient() string {
	clientResult := esbuild.Build(esbuild.BuildOptions{
		EntryPoints: []string{"./frontend/clientEntry.jsx"},
		Bundle:      true,
		Write:       true,
	})
	clientBundleString := string(clientResult.OutputFiles[0].Contents)
	return clientBundleString
}

func main() {
	// reactをSSRするためのバンドルを作成
	backendBundle := buildBackend()
	// reactをhydrateするためのバンドルを作成
	clientBundle := buildClient()

	ctx := v8.NewContext(nil)
	_, err := ctx.RunScript(backendBundle, "bundle.js")
	if err != nil {
		log.Fatalf("Failed to evaluate bundled script: %v", err)
	}
	val, err := ctx.RunScript("renderApp()", "render.js")
	if err != nil {
		log.Fatalf("Failed to render React component: %v", err)
	}
	renderedHTML := val.String()

	tmpl, err := template.New("webpage").Parse(htmlTemplate)
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}

	// backendから渡すprops
	initialProps := InitialProps{
		Name:          "GoでReactをSSRする",
		InitialNumber: 100,
	}
	jsonProps, err := json.Marshal(initialProps)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		data := PageData{
			RenderedContent: template.HTML(renderedHTML),
			InitialProps:    template.JS(jsonProps),
			JS:              template.JS(clientBundle),
		}
		err := tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	fmt.Println("Server is running at http://localhost:3002")
	log.Fatal(http.ListenAndServe(":3002", nil))
}
