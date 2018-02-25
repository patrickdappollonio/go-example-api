package app

import (
	"html/template"
	"strings"
	"time"
)

var tmplCreate, tmplRead *template.Template

const _create_bin = `<!doctype html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
		<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css"
		integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
		<title>Create an inspector bin</title>
		<style type="text/css">
		#top {
			margin-top: 20px;
		}
		input[type="text"] {
			text-overflow: ellipsis;
		}
		</style>
	</head>
	<body>
		<div id="top" class="container">
			<div class="col">
				<div class="jumbotron">
					<h1 class="display-4">Create URLs and inspect its requests</h1>
					<p class="lead">This utility allows you to create URLs you can send requests to that can be read
					later through the browser.</p>
					<p>Send any <code>GET</code>, <code>POST</code>, <code>PUT</code>, <code>PATCH</code> and <code>DELETE</code>
					requests to the URL given to you and review it using the inspector URL below.</p>
				</div>
			</div>
		</div>
		<div class="container">
			<div class="col" id="app">
				<div class="form-group">
					<label class="lead">Use the following URL <strong>to send your requests to</strong>:</label>
					<input type="text" class="form-control-plaintext form-control-lg" readonly :value="requestURL" @click="select">
				</div>
				<div class="form-group">
					<label class="lead">Use the following URL to <strong>inspect your requests</strong>:</label>
					<input type="text" class="form-control-plaintext form-control-lg" readonly :value="inspectURL" @click="select">
				</div>
				<div class="form-group">
					<a class="btn btn-primary" :href="inspectURL">Go to inspect URL &rarr;</a>
				</div>
			</div>
		</div>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/vue/2.5.13/vue.js"
		integrity="sha256-pU9euBaEcVl8Gtg+FRYCtin2vKLN8sx5/4npZDmY2VA=" crossorigin="anonymous"></script>

		<script type="text/javascript">
		new Vue({
			el: '#app',
			data: {
				id: '{{ .ID }}'
			},
			computed: {
				inspectURL: function() {
					return window.location.origin + '/inspector/' + this.id
				},
				requestURL: function() {
					return window.location.origin + '/r/' + this.id
				}
			},
			methods: {
				select: function(evt) {
					evt.target.select()
				}
			}
		})
		</script>
	</body>
</html>
`
const _read_bins = `<!doctype html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
		<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css"
		integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
		<title>Reading requests for /r/{{ .ID }}</title>
		<style type="text/css">
		body {
			padding-top: 20px;
			padding-bottom: 40px;
		}
		.sep {
			display: block;
			margin-top: 35px;
		}
		.card {
			margin-top: 45px;
		}
		.card:hover {
			background-color: rgba(0, 0, 0, .05);
		}
		.card .section-title {
			border-bottom: 1px solid rgba(0, 0, 0, 0.125);
			padding-bottom: 5px;
			font-weight: bold;
		}
		</style>
	</head>
	<body>
		<div id="top" class="container">
			<div class="col">
				<h1 class="display-4">Reviewing requests to <code>/r/{{ .ID }}</code></h1>
				<p class="lead">Below you'll find all requests sent to the URL above, showing newest first. Requests
				<strong>will live for 24 hours before being deleted from the server</strong>.</p>
			</div>
		</div>
		<div class="container">
			<div class="col">
				{{ if .Records }}{{ range .Records }}
				<div class="card">
					<div class="card-body">
						<h5 class="card-title"><code>{{ .Method }} {{ .Path }}</code></h5>
						<h6 class="card-subtitle mb-2 text-muted">Saved on {{ .Time.Format date }}</h6>

						<div class="sep"></div>

						<h6 class="section-title">Headers:</h6>
						<p class="card-text">
							<dl class="row">
								{{ range $key, $value := .Headers }}
								<dt class="col-sm-3">{{ $key }}</dt>
								<dd class="col-sm-9">
									{{ range $value }}
										<code>{{ . }}</code><br>
									{{ end }}
								</dt>
								{{ end }}
							</dl>
						</p>

						{{ with .Query }}
						<h6 class="section-title">Querystring parameters:</h6>
						<p class="card-text">
							<dl class="row">
								{{ range $key, $value := . }}
								<dt class="col-sm-3"><code>{{ $key }}</code></dt>
								<dd class="col-sm-9"><code>{{ $value | join }}</code></dt>
								{{ end }}
							</dl>
						</p>
						{{ end }}

						<h6 class="section-title">Body:</h6>
						<p class="card-text">
							{{ if .Body }}
								<pre><code>{{ .Body }}</code></pre>
							{{ else }}
								<em>Request body had no content.</em>
							{{ end }}
						</p>
					</div>
				</div>
				{{ end }}{{ else }}
				<div class="alert alert-primary" role="alert">
					<strong>No requests captured yet!</strong> Start by making some. Need help on that? Visit
					<a href="javascript:return false" data-id="{{ .ID }}">{{ .ID }}</a> with your browser or
					use something like <a href="https://www.getpostman.com/" target="_blank">Postman</a>,
					<a href="https://insomnia.rest/" target="_blank">Insomnia</a> or
					<a href="https://paw.cloud/" target="_blank">Paw (for MacOS)</a>
				</div>
				{{ end }}
			</div>
		</div>

		<script src="https://code.jquery.com/jquery-3.2.1.slim.min.js"
		integrity="sha384-KJ3o2DKtIkvYIK3UENzmM7KCkRr/rE9/Qpg6aAZGJwFDMVNA/GpGFF93hXpG5KkN" crossorigin="anonymous"></script>
		<script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/js/bootstrap.min.js"
		integrity="sha384-JZR6Spejh4U02d8jOt6vLEHfe/JQGiRRSQQxSfFWpi1MquVdAyjUar5+76PVCmYl" crossorigin="anonymous"></script>
		<script type="text/javascript">
		$(window).on('load', function(){
			'use strict';
			$('[data-id]').each(function(){
				var e = $(this);
				e.text('/r/' + e.attr('data-id'));
				e.attr('href', window.location.origin + '/r/' + e.attr('data-id'));
				e.attr('target', '_blank');
			})
		})
		</script>
	</body>
</html>
`

func setupTemplate() {
	funcs := template.FuncMap{
		"date": func() string {
			return time.RFC1123
		},
		"join": func(a []string) string {
			return strings.Join(a, ", ")
		},
	}

	tmplCreate = template.Must(template.New("create").Funcs(funcs).Parse(_create_bin))
	tmplRead = template.Must(template.New("read").Funcs(funcs).Parse(_read_bins))
}
