Multi-pattern string matching in Go
==================
Tool used for searching multiple patterns in a single Log file (single line Logs only) and printing results in JSON file.

By <a href="mailto:xgam33@gmail.com">Filip Halas</a>.

Setting it up on Linux
-----------------------
1. Download this repository.
2. Open a command-line and navigate to the root of this tool directory.
3. Review Makefile if you wish.
4. Use <code>make</code> to build the tool.
5. After that you can use: <code>./jsonizer</code> to run it.

Setting it up on Windows
-----------------------------
1. Download and install Go from <a href="https://code.google.com/p/go/downloads/list">here</a>.
2. Check that environmental variables are set correctly - try executing: <code>go</code> in your command line.
 * In case of failure (or when using ZIP archive) you might need to set them manually:
    * navigate to <code>Control Panel - System - Advanced (tab) - Environment Variables - System variables</code>
    * append to variable <code>Path</code>your Go installation location: <code>d:\Program Files\Go\bin;</code>
4. After that you can either create a standalone executable (compile) or run this tool.
 * Open a command-line and navigate to the root of this tool directory.
    * For <b>running</b> it once use: <code>go run jsonizer.go</code>.
    * For <b>compiling</b> it to Windows executable use: <code>go build jsonizer.go</code>.
