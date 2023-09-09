package main

var config string = `
##
# Virtual Host configuration for example.com
#
# You can move that to a different file under sites-available/ and symlink that
# to sites-enabled/ to enable it.
#
server {
		listen 80;
		listen [::]:80;

		server_name example.com;

		root /usr/share/nginx/example.com/html;
		index index.html;

		location / {
				try_files $uri $uri/ =404;
		}
}`

var usage string = `
Usage:
Create and update virtual host
clinx -d=[local domain name] -f=[folder to serve]

Delete virtual host
clinx -r=[local domain name]
`
