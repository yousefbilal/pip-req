package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
)

func main() {
	var args struct {
		PackageName    string `arg:"positional,required"`
		PackageVersion string `arg:"-v,--version"`
	}
	arg.MustParse(&args)

	apiUrl := "https://pypi.org/pypi/" + args.PackageName
	if args.PackageVersion != "" {
		apiUrl += "/" + args.PackageVersion
	}
	apiUrl += "/json"

	resp, err := http.Get(apiUrl)

	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	var respBody map[string]any

	err = json.NewDecoder(resp.Body).Decode(&respBody)

	if err != nil {
		fmt.Printf("error decoding response body: %s\n", err)
		os.Exit(1)
	}

	data := respBody["info"].(map[string]interface{})["requires_dist"].([]interface{})
	var requiresDist []string
	for _, item := range data {
		requiresDist = append(requiresDist, item.(string))
	}

	fmt.Println(strings.Join(requiresDist, "\n"))
}
