package main

var PodInfoTemplate string = `
Podname: {{.Name}}
========================
Master: {{.MasterIP}}:{{.MasterPort}}
Quorum: {{.Quorum}}
Auth Token: {{.Authpass}}
Known Sentinels: {{ range .KnownSentinels }}
	{{.}}
{{ end }}
Known Slaves: {{ range .KnownSlaves }}
	{{.}}
{{ end }}
Settings: {{ range $k,$v := .Settings }} 
	{{printf "%-30s" $k}}     {{printf "%10s" $v}} {{ end }}
`
