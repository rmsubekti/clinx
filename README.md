Just a tool for deploying local documentation server, it help me to configure virtual host server for static website using nginx

## How to use
- clone this repository
- run `go build`
- copy/move binary to /usr/local/bin `sudo mv clinx /usr/local/bin `

### Creating virtual host
- run command with sudo `sudo clinx -d=example.com.local -f=folderNameContainStaticFiles`

### Removing virtual host
- run command with sudo `sudo clinx -r=example.com.local`

## Changed Directory/Files
- /etc/hosts
- /usr/share/nginx
- /etc/nginx/sites-available
- /etc/nginx/sites-enabled