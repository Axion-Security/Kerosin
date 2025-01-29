package main

import (
	"fmt"
	"strings"
)

func ConvertGET_Go(input string) string {
	lines := strings.Split(input, "\n")
	requestLine := strings.Fields(lines[0])

	method := requestLine[1]
	urlStr := requestLine[2]

	headers := []string{}

	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line != "" {
			headers = append(headers, line)
		}
	}

	functionName := generateFunctionName(strings.ReplaceAll(urlStr, "\"", ""))
	code := fmt.Sprintf("func %s() (string, error) {\n", functionName)
	code += fmt.Sprintf("    method := \"%s\"\n", method)
	code += fmt.Sprintf("    url := %s\n", urlStr)
	code += "    requestBody := \"\"\n\n"
	code += "    req, err := http.NewRequest(method, url, strings.NewReader(requestBody))\n"
	code += "    if err != nil {\n"
	code += "        return \"\", err\n"
	code += "    }\n\n"

	for _, header := range headers {
		headerCleaned := strings.ReplaceAll(header, "HEADER ", "")
		headerParts := strings.SplitN(headerCleaned, ": ", 2)
		if len(headerParts) == 2 {
			key := headerParts[0]
			value := headerParts[1]
			code += fmt.Sprintf("    req.Header.Set(%s\", \"%s)\n", key, value)
		}
	}

	code += "\n"
	code += "    client := &http.Client{}\n"
	code += "    resp, err := client.Do(req)\n"
	code += "    if err != nil {\n"
	code += "        return \"\", err\n"
	code += "    }\n"
	code += "    defer resp.Body.Close()\n\n"
	code += "    body, err := ioutil.ReadAll(resp.Body)\n"
	code += "    if err != nil {\n"
	code += "        return \"\", err\n"
	code += "    }\n\n"
	code += "    responseBody := string(body)\n\n"
	code += "    return responseBody, nil\n"
	code += "}\n\n"

	return code
}
func ConvertPOST_Go(input string) string {
	lines := strings.Split(input, "\n")
	requestLine := strings.Fields(lines[0])
	urlStr := requestLine[2]

	headers := []string{}
	var content, contentType string

	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line != "" {
			if strings.HasPrefix(line, "HEADER ") {
				headers = append(headers, line)
			} else if strings.HasPrefix(line, "CONTENT ") {
				content = strings.TrimSpace(strings.TrimPrefix(line, "CONTENT "))
			} else if strings.HasPrefix(line, "CONTENTTYPE ") {
				contentType = strings.TrimSpace(strings.TrimPrefix(line, "CONTENTTYPE "))
			}
		}
	}

	functionName := generateFunctionName(strings.ReplaceAll(urlStr, "\"", ""))
	code := fmt.Sprintf("func %s() (string, error) {\n", functionName)
	code += fmt.Sprintf("    url := %s\n", urlStr)
	code += fmt.Sprintf("    requestBody := %s\n", content)
	code += fmt.Sprintf("    contentType := %s\n\n", contentType)
	code += "    req, err := http.NewRequest(\"POST\", url, strings.NewReader(requestBody))\n"
	code += "    if err != nil {\n"
	code += "        return \"\", err\n"
	code += "    }\n\n"

	for _, header := range headers {
		headerCleaned := strings.ReplaceAll(header, "HEADER ", "")
		headerParts := strings.SplitN(headerCleaned, ": ", 2)
		if len(headerParts) == 2 {
			key := headerParts[0]
			value := headerParts[1]
			code += fmt.Sprintf("    req.Header.Set(%s\", \"%s)\n", key, value)
		}
	}

	code += "\n"
	code += "    client := &http.Client{}\n"
	code += "    resp, err := client.Do(req)\n"
	code += "    if err != nil {\n"
	code += "        return \"\", err\n"
	code += "    }\n"
	code += "    defer resp.Body.Close()\n\n"
	code += "    body, err := ioutil.ReadAll(resp.Body)\n"
	code += "    if err != nil {\n"
	code += "        return \"\", err\n"
	code += "    }\n\n"
	code += "    responseBody := string(body)\n\n"
	code += "    return responseBody, nil\n"
	code += "}\n\n"

	return code
}
func ConvertKeycheck_Go(input string) string {
	functionName := "KeyCheck"
	code := fmt.Sprintf("func %s(response string) string {\n", functionName)

	conditions := strings.Split(input, "KEYCHAIN")

	for _, condition := range conditions {
		condition = strings.TrimSpace(condition)
		keywords := strings.Split(condition, "KEY")
		keychainKeyword := strings.TrimSpace(keywords[0])

		if keychainKeyword != "" {
			keywords = keywords[1:]
			code += fmt.Sprintf("    if %s {\n", generateKeychainConditionGo(keychainKeyword, keywords))
			code += "        return \"success\"\n" // or "failure" for failure cases
			code += "    }\n"
		}
	}

	code += "    return \"Unknown\"\n"
	code += "}\n\n"

	return code
}
func ConvertJSParser_Go(input string) string {
	code := ""

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line != "" && strings.HasPrefix(line, "PARSE") {
			words := strings.Fields(line)
			if len(words) >= 4 {
				jsonPath := words[3]
				functionName := fmt.Sprintf("%sParse", strings.Title(strings.ReplaceAll(jsonPath, "\"", "")))
				code += fmt.Sprintf("func %s(response string) (string, error) {\n", functionName)
				code += "    var data map[string]interface{}\n"
				code += "    if err := json.Unmarshal([]byte(response), &data); err != nil {\n"
				code += "        return \"\", err\n"
				code += "    }\n"
				code += fmt.Sprintf("    value, ok := data[%s].(string)\n", jsonPath)
				code += "    if !ok {\n"
				code += fmt.Sprintf("        return \"\", fmt.Errorf(\"Key '%s' not found or not a string in the JSON response\")\n", strings.ReplaceAll(jsonPath, "\"", ""))
				code += "    }\n"
				code += "    return value, nil\n"
				code += "}\n\n"
			}
		}
	}

	return code
}
func ConvertLRParser_Go(input string) string {
	code := ""

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line != "" && strings.HasPrefix(line, "PARSE") {
			words := strings.Fields(line)
			if len(words) >= 7 {
				before := strings.Trim(words[3], " ")
				after := strings.Trim(words[4], " ")
				functionName := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(before, "\"", ""), "\\", ""), ":", "")
				code += fmt.Sprintf("func %s(response string) (string, error) {\n", functionName)
				code += fmt.Sprintf("    idx := strings.Index(response, %s)\n", before)
				code += "    if idx == -1 {\n"
				code += "        return \"\", fmt.Errorf(\"substring not found in input\")\n"
				code += "    }\n"
				code += fmt.Sprintf("    start := idx + len(%s)\n", before)
				code += fmt.Sprintf("    end := strings.Index(response[start:], %s)\n", after)
				code += "    if end == -1 {\n"
				code += "        return \"\", fmt.Errorf(\"substring not found in after\")\n"
				code += "    }\n"
				code += "    return response[start : start+end], nil\n"
				code += "}\n\n"
			}
		}
	}

	return code
}
func generateKeychainConditionGo(keychainKeyword string, keywords []string) string {
	conditions := ""

	for _, keyword := range keywords {
		keyword = strings.TrimSpace(keyword)
		if keyword != "" {
			conditions += fmt.Sprintf("strings.Contains(response, %s) || ", keyword)
		}
	}

	if len(conditions) > 4 {
		conditions = conditions[:len(conditions)-4] // Remove the trailing " || "
	}

	return conditions
}
