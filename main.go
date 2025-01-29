package main

import (
	"net/url"
	"strings"
)

func generateFunctionName(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "createRequest"
	}
	hostParts := strings.Split(u.Hostname(), ".")
	if len(hostParts) < 2 {
		return "createRequest"
	}
	domainName := hostParts[len(hostParts)-2]
	return "createRequestFor" + strings.Title(domainName)
}

func ConvertToGo(module, input string) string {
	switch module {
	case "Request":
		if strings.Contains(input, "POST") {
			return ConvertPOST_Go(input)
		} else if strings.Contains(input, "GET") {
			return ConvertGET_Go(input)
		} else {
			return "Couldn't identify Block"
		}
	case "KeyCheck":
		return ConvertKeycheck_Go(input)
	case "Parser":
		if strings.Contains(input, "LR") {
			return ConvertLRParser_Go(input)
		} else if strings.Contains(input, "JSON") {
			return ConvertJSParser_Go(input)
		} else {
			return "Couldn't identify Block"
		}
	default:
		return "Invalid module"
	}
}

func ConvertToCSharp(module, input string) string {
	switch module {
	case "Request":
		if strings.Contains(input, "POST") {
			return ConvertPOST_CSharp(input)
		} else if strings.Contains(input, "GET") {
			return ConvertGET_CSharp(input)
		} else {
			return "Couldn't identify Block"
		}
	case "KeyCheck":
		return ConvertKeycheck_CSharp(input)
	case "Parser":
		if strings.Contains(input, "LR") {
			return ConvertLRParser_CSharp(input)
		} else if strings.Contains(input, "JSON") {
			return ConvertJSParser_CSharp(input)
		} else {
			return "Couldn't identify Block"
		}
	default:
		return "Invalid module"
	}
}

func main() {
	ConvertToGo("MODULE NAME", "MODULE")
	ConvertToCSharp("MODULE NAME", "MODULE")
}
