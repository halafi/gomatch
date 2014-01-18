Multi-pattern string matching in Go
==================
Tool used to search for multiple patterns in log records (single line Logs only) and printing results in JSON. 

See the project <a href="https://github.com/halafi/String-matching-Go/wiki">wiki</a> for help with configuration and some examples of usage.

převod logovacích záznamů do JSON notace

Linux installation
-----------------------
1. Download this repository, <code>git clone https://github.com/halafi/String-matching-Go.git</code> in your command-line should do.
2. Open a command-line and navigate to the root of this tool directory.
3. Review Makefile if you wish. Choose whether you wish to install/uninstall dependecies.
4. Use <code>make</code> to build the tool.
5. After that you can use: <code>jsonizer</code> to run it.

Windows installation
-----------------------
1. Download Go: https://code.google.com/p/go/downloads/list.
2. Download this repository: https://github.com/halafi/String-matching-Go/archive/master.zip.
3. Execute <code>win_build.cmd</code> to build <code>jsonizer.exe</code>.

By <a href="mailto:xgam33@gmail.com">Filip Halas</a>.
