String matching in Go
==================
Implementations of string matching alghoritms in Go language.

By <a href="mailto:xgam33@gmail.com">Filip Halas</a>.

setting up the Go environment
-----------------------------
1. download and install Go from <a href="https://code.google.com/p/go/downloads/list">here</a>
2. check that environmental variables are set correctly - try executing <code>go</code> in your command line 
3. in case of failure (or when using ZIP archive) you will need to set them manually
   * <b>Windows</b>: go to <code>Control Panel - System - Advanced (tab) - Environment variables - system variables</code>
   * <b>Unix/Linux</b>: try executing in your command line: <code>export GOROOT=$HOME/golang/go export PATH=$PATH:$GOROOT/bin</code>
   * variable <b>GOROOT</b> should be set to something like <code>C:\go</code>  on Windows or <code>$HOME/golang/go</code> on Unix/Linux
   * variable <b>GOPATH</b> to something like <code>$GOROOT\bin</code> on Windows or <code>$GOROOT/bin</code> on Unix/Linux

running the source code
-----------------------
* For <b>running</b> go file in your command line use: <code>go run filename.go</code>
* For <b>compiling</b> go file to Windows executable use: <code>go build filename.go</code>
