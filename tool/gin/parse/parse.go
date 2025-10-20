package parse

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

func CheckPrefix(line string, object ...string) bool {
	for _, o := range object {
		if strings.HasPrefix(line, o) {
			return true
		}
	}
	return false
}

func StandardAddr(line string) string {
	if CheckPrefix(line, "/") {
		return line
	}
	return "/" + line
}

func ParesApi(addr string) (*Api, error) {
	apiContent, err := os.Open(addr)
	if err != nil {
		return nil, err
	}
	defer apiContent.Close()

	var api Api
	var ans strings.Builder
	scanner := bufio.NewScanner(apiContent)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 提取type
		if CheckPrefix(line, "type") {
			if CheckPrefix(line, "type (", "type(") {
				ans.WriteString(line + "\n")
				GetType(scanner, ")", &ans)
			} else {
				line = strings.Replace(line, "{", "struct {", 1)
				ans.WriteString(line + "\n")
				GetType(scanner, "}", &ans)
			}
		}

		// 提取server
		if CheckPrefix(line, "@server") {
			api.Server = GetServer(scanner)
		}

		// 提取service
		if CheckPrefix(line, "service") {
			line = strings.Replace(line, "{", "", 1)
			line = strings.Replace(line, "service", "", 1)
			line = strings.TrimSpace(line)
			api.ServiceName = line
			api.Service = GetService(scanner)
		}
	}

	api.T = ans.String()
	return &api, nil
}

func GetType(scanner *bufio.Scanner, symbol string, ans *strings.Builder) {
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasSuffix(line, "{") {
			line = strings.Replace(line, "{", "struct {", 1)
		}
		if CheckPrefix(line, symbol) {
			ans.WriteString(line + "\n" + "\n")
			return
		}
		ans.WriteString(line + "\n")
	}
}

func GetServer(scanner *bufio.Scanner) *Server {
	var prefix, group string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if CheckPrefix(line, "prefix") {
			prefix = StandardAddr(strings.TrimSpace(line[len("prefix:"):]))
		}
		if CheckPrefix(line, "group") {
			group = strings.TrimSpace(line[len("group:"):])
		}
		if CheckPrefix(line, ")") {
			return &Server{prefix, group}
		}
	}
	return nil
}

func GetService(scanner *bufio.Scanner) []*Service {
	var service []*Service
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if CheckPrefix(line, "@handler") {
			var s Service
			s.Handler = strings.TrimSpace(line[len("@handler:"):])
			scanner.Scan()
			line = strings.TrimSpace(scanner.Text())
			s.Method = parseMethod(line)
			service = append(service, &s)
		}
	}
	return service
}

func parseMethod(s string) *Method {
	re := regexp.MustCompile(`(?i)^\s*(\w+)\s+(\S+)(?:\s+([^\s]+))?(?:\s+returns\s+([^\s]+))?\s*$`)
	matches := re.FindStringSubmatch(s)
	if len(matches) < 3 {
		return nil
	}

	if matches[3] != "" {
		matches[3] = matches[3][1 : len(matches[3])-1]
	}

	if matches[4] != "" {
		matches[4] = matches[4][1 : len(matches[4])-1]
	}

	return &Method{
		Method: strings.ToUpper(matches[1]),
		Route:  StandardAddr(matches[2]),
		Req:    matches[3],
		Resp:   matches[4],
	}
}
