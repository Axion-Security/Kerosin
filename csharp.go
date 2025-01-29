package main

import (
	"fmt"
	"strings"
)

func ConvertGET_CSharp(input string) string {
	lines := strings.Split(input, "\n")

	headers := []string{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			if strings.HasPrefix(line, "REQUEST GET") {
				// Extract URL
				urlStr := strings.TrimSpace(strings.TrimPrefix(line, "REQUEST GET"))
				headers = append(headers, fmt.Sprintf("string url = %s;", urlStr))
			} else if strings.HasPrefix(line, "HEADER ") {
				// Extract headers
				header := strings.TrimSpace(strings.TrimPrefix(line, "HEADER "))
				headerParts := strings.SplitN(header, ": ", 2)
				if len(headerParts) == 2 {
					key := headerParts[0]
					value := headerParts[1]
					headers = append(headers, fmt.Sprintf("request.Headers.Add(%s\", \"%s);", key, value))
				}
			}
		}
	}

	functionName := generateFunctionName(strings.ReplaceAll(headers[0], "\"", ""))
	code := fmt.Sprintf("public async Task<string> %sAsync()\n", functionName)
	code += "{\n"
	code += "    string method = \"GET\";\n"
	code += "    " + headers[0] + "\n"
	code += "    HttpRequestMessage request = new HttpRequestMessage(HttpMethod.Get, url);\n\n"

	for _, header := range headers[1:] {
		code += "    " + header + "\n"
	}

	code += "\n"
	code += "    HttpResponseMessage response = await _httpClient.SendAsync(request);\n"
	code += "\n"
	code += "    if (!response.IsSuccessStatusCode)\n"
	code += "    {\n"
	code += "        return \"\"; // Handle error here\n"
	code += "    }\n"
	code += "\n"
	code += "    string responseBody = await response.Content.ReadAsStringAsync();\n"
	code += "\n"
	code += "    return responseBody;\n"
	code += "}\n"

	return code
}
func ConvertPOST_CSharp(input string) string {
	lines := strings.Split(input, "\n")

	headers := []string{}
	var content, contentType string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			if strings.HasPrefix(line, "REQUEST POST") {
				// Extract URL
				urlStr := strings.TrimSpace(strings.TrimPrefix(line, "REQUEST POST"))
				headers = append(headers, fmt.Sprintf("string url = \"%s\";", urlStr))
			} else if strings.HasPrefix(line, "HEADER ") {
				// Extract headers
				header := strings.TrimSpace(strings.TrimPrefix(line, "HEADER "))
				headerParts := strings.SplitN(header, ": ", 2)
				if len(headerParts) == 2 {
					key := headerParts[0]
					value := headerParts[1]
					headers = append(headers, fmt.Sprintf("request.Headers.Add(\"%s\", \"%s\");", key, value))
				}
			} else if strings.HasPrefix(line, "CONTENT ") {
				// Extract content
				content = line[len("CONTENT "):]
			} else if strings.HasPrefix(line, "CONTENTTYPE ") {
				// Extract content type
				contentType = line[len("CONTENTTYPE "):]
			}
		}
	}

	functionName := generateFunctionName(strings.ReplaceAll(headers[0], "\"", ""))
	code := fmt.Sprintf("public async Task<string> %sAsync()\n", functionName)
	code += "{\n"
	code += "    string method = \"POST\";\n"
	code += fmt.Sprintf("    string url = %s;\n", headers[0]) // Extracted URL
	code += fmt.Sprintf("    string requestBody = \"%s\";\n", content)
	code += fmt.Sprintf("    string contentType = \"%s\";\n", contentType)
	code += "    HttpRequestMessage request = new HttpRequestMessage(HttpMethod.Post, url);\n"
	code += "    request.Content = new StringContent(requestBody, Encoding.UTF8, contentType);\n\n"

	for _, header := range headers[1:] { // Skip the URL header
		code += "    " + header + "\n"
	}

	code += "\n"
	code += "    HttpResponseMessage response = await _httpClient.SendAsync(request);\n"
	code += "\n"
	code += "    if (!response.IsSuccessStatusCode)\n"
	code += "    {\n"
	code += "        return \"\"; // Handle error here\n"
	code += "    }\n"
	code += "\n"
	code += "    string responseBody = await response.Content.ReadAsStringAsync();\n"
	code += "\n"
	code += "    return responseBody;\n"
	code += "}\n"

	return code
}
func ConvertKeycheck_CSharp(input string) string {
	functionName := "KeyCheck"
	code := fmt.Sprintf("public string %s(string response)\n", functionName)
	code += "{\n"

	conditions := strings.Split(input, "KEYCHAIN")

	for _, condition := range conditions {
		condition = strings.TrimSpace(condition)
		keywords := strings.Split(condition, "KEY")
		keychainKeyword := strings.TrimSpace(keywords[0])

		if keychainKeyword != "" {
			keywords = keywords[1:]
			code += fmt.Sprintf("    if (%s)\n", generateKeychainCondition(keychainKeyword, keywords))
			code += "    {\n"
			code += "        return \"success\";\n" // or "failure" for failure cases
			code += "    }\n"
		}
	}

	code += "    return \"Unknown\";\n"
	code += "}\n\n"

	return code
}
func ConvertJSParser_CSharp(input string) string {
	code := ""

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line != "" && strings.HasPrefix(line, "PARSE") {
			// Split the line into tokens
			tokens := strings.Fields(line)

			// Check if the line has enough tokens
			if len(tokens) >= 6 {
				// Extract relevant information from the tokens
				//openingTag := tokens[1]
				parserType := tokens[2]
				jsonPath := tokens[3]
				varName := tokens[len(tokens)-1]

				// Check if the parser type is "JSON"
				if parserType == "JSON" {
					functionName := fmt.Sprintf("%sParse", strings.Title(strings.ReplaceAll(jsonPath, "\"", "")))
					code += fmt.Sprintf("public static (string, Exception) %s(string response)\n", functionName)
					code += "{\n"
					code += "    try\n"
					code += "    {\n"
					code += "        Dictionary<string, object> data = JsonConvert.DeserializeObject<Dictionary<string, object>>(response);\n"
					code += fmt.Sprintf("        if (data.TryGetValue(%s, out object value) && value is string stringValue)\n", jsonPath)
					code += "        {\n"
					code += fmt.Sprintf("            var %s = stringValue;\n", varName)
					code += "        }\n"
					code += "        else\n"
					code += "        {\n"
					code += fmt.Sprintf("            return (null, new Exception(\"Key '%s' not found or not a string in the JSON response\"));\n", jsonPath)
					code += "        }\n"
					code += "    }\n"
					code += "    catch (Exception ex)\n"
					code += "    {\n"
					code += "        return (null, ex);\n"
					code += "    }\n"
					code += "}\n\n"
				}
			}
		}
	}

	return code
}
func ConvertLRParser_CSharp(input string) string {
	code := ""

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line != "" && strings.HasPrefix(line, "PARSE") {
			words := strings.Fields(line)
			if len(words) >= 7 {
				before := strings.Trim(words[3], " ")
				after := strings.Trim(words[4], " ")
				functionName := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(before, "\"", ""), "\\", ""), ":", "")
				code += fmt.Sprintf("public static (string, Exception) %s(string response)\n", functionName)
				code += "{\n"
				code += fmt.Sprintf("    int idx = response.IndexOf(%s);\n", before)
				code += "    if (idx == -1)\n"
				code += "    {\n"
				code += "        return (\"\", new Exception(\"substring not found in input\"));\n"
				code += "    }\n"
				code += fmt.Sprintf("    int start = idx + %s.Length;\n", before)
				code += fmt.Sprintf("    int end = response.IndexOf(%s, start);\n", after)
				code += "    if (end == -1)\n"
				code += "    {\n"
				code += "        return (\"\", new Exception(\"substring not found in after\"));\n"
				code += "    }\n"
				code += "    return (response.Substring(start, end - start), null);\n"
				code += "}\n\n"
			}
		}
	}

	return code
}
func generateKeychainCondition(keychainKeyword string, keywords []string) string {
	conditions := ""

	for _, keyword := range keywords {
		keyword = strings.TrimSpace(keyword)
		if keyword != "" {
			conditions += fmt.Sprintf("response.Contains(%s) || ", keyword)
		}
	}

	if len(conditions) > 4 {
		conditions = conditions[:len(conditions)-4] // Remove the trailing " || "
	}

	return conditions
}
