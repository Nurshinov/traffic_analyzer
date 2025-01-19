package web

import (
	"html/template"
	"net/http"
	"sort"
	"traffic_analyzer/internal/db"
)

func getNetworkStat(w http.ResponseWriter, _ *http.Request) {
	dbc := db.New()
	defer dbc.Close()

	ipSl := dbc.GetAll()

	sort.Slice(ipSl, func(i, j int) bool {
		return ipSl[i].All > ipSl[j].All
	})

	if err := index.Execute(w, ipSl); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var index = template.Must(template.New("index").Parse(indexTmpl))

const indexTmpl = `<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Traffic analyzer</title>
	</head>
	<body>
		<table>
			<tr>
				<th>IP</th>
				<th>ALL</th>
				<th>Retransmitted</th>
			</tr>
		{{range .}}
			<tr>
				<td>{{.IP}}</td>
				<td>{{.All}}</td>
				<td>{{.Retransmitted}}</td>
			</tr>
		{{end}}
		</table>
	</body>
</html>
`
