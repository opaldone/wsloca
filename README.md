<h1 align="center">
  <img src="./lo.svg" alt="WSloca">
  <br />
  WSloca
  <br />
</h1>
<h4 align="center">
  WebSocket server for broadcasting location information between senders and receivers
</h4>
<h3>
Built with these excellent libraries
<img src="https://go.dev/blog/go-brand/Go-Logo/SVG/Go-Logo_Blue.svg" height="45px" vertical-align="middle" />
</h3>

* [gorilla-websocket](https://github.com/gorilla/websocket)
* [julienschmidt-httprouter](https://github.com/julienschmidt/httprouter)
* [acme-autocert](https://pkg.go.dev/golang.org/x/crypto/acme/autocert)
<h1></h1>

### How to install and compile
##### Clonning
```bash
git clone https://github.com/opaldone/wsloca.git
```
##### Go to the root "wsloca" directory
```bash
cd wsloca
```
##### Set your GOPATH to the "wsloca" directory to keep your global GOPATH clean
```bash
export GOPATH=$(pwd)
```
##### Go to the source folder
```bash
cd src/wsloca
```
##### Installing the required Golang packages
```bash
go mod init
```
```bash
go mod tidy
```
##### Return to the "wsloca" root directory, You can see the "wsloca/pkg" folder that contains the required Golang packages
```bash
cd ../..
```
##### Compiling by the "r" bash script
> r - means "run", b - means "build"
```bash
./r b
```
##### Creating the required folders structure and copying the frontend part by the "u" bash script
> The "u" script is a watching script then for stopping press Ctrl+C \
> u - means "update"
```bash
./u
```
> The "u" script reads sub file "watch_files" \
> E_FOLDERS - the array of creating empty folders \
> C_FOLDERS - the array of folders to simple copy \
> W_FILES - the array of files whose changes are tracked
```bash
./watch_files
```
##### You can check the "wsloca/bin" folder. It should contain the necessary structure of folders and files
```bash
ls -lash --group-directories-first bin
```
##### Start the server
```bash
./r
```
### About config
The config file is located here __wsloca/bin/config.json__
```JavaScript
{
  // Just a name of application
  "appname": "wsloca",

  // IP address of the server, zeros mean current host
  "address": "0.0.0.0",

  // Port, don't forget to open for firewall
  "port": 8080,

  // The folder that stores the frontend part of the site
  "static": "static",

  // Set "acme": true if You need to use acme/autocert
  // false - if You use self-signed certificates
  "acme": false,

  // The array of domain names, set "acme": true
  "acmehost": [
    "opaldone.click",
    "206.189.101.23",
    "www.opaldone.click"
  ],

  // The folder where acme/autocert will store the keys, set "acme": true
  "dirCache": "./certs",

  // The paths to your self-signed HTTPS keys, set "acme": false
  "crt": "/server.crt",
  "key": "/server.key",
}
```

### License
MIT License - see [LICENSE](LICENSE) for full text
