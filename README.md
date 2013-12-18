Multi-pattern string matching in Go
==================
Tool used for searching multiple patterns in a single Log file (single line Logs only) and printing results in JSON file.

By <a href="mailto:xgam33@gmail.com">Filip Halas</a>.

Setting it up on Linux
-----------------------
1. Download this repository.
2. Open a command-line and navigate to the root of this tool directory.
3. Review Makefile if you wish.
4. Use Make to build the tool:
 * <code>make</code>
 * this should take care of everything.
5. After that you can use:
 * <code>./jsonizer</code>
 * to run it.

Setting it up on Windows
-----------------------------
1. Download and install Go from <a href="https://code.google.com/p/go/downloads/list">here</a>.
2. Check that environmental variables are set correctly - try executing:
 * <code>go</code>
 * in your command line.
 * In case of failure (or when using ZIP archive) you might need to set them manually:
    * navigate to <code>Control Panel - System - Advanced (tab) - Environment Variables - System variables</code>
    * <code>GOROOT</code> should be set to something like <code>C:\Program Files\Go</code>,
    * <code>GOPATH</code> should be set to <code>$GOROOT\bin</code>.
4. After that you can either create a standalone executable (compile) or run this tool.
 * Open a command-line and navigate to the root of this tool directory.
    * For <b>running</b> go file in your command line use: <code>go run filename.go</code>.
    * For <b>compiling</b> go file to Windows executable use: <code>go build filename.go</code>.

Configuration
==================
Make sure that there are no extra spaces or endlines in these files and that they are ANSI encoded.

Patterns.txt
-----------------------------
* Each line in <b>patterns.txt</b> corresponds to one match, that will be searched for.
* Each pattern line starts with a name separated by <code>##</code>
* Pattern might consist of:
 * <code>&lt;TOKEN&gt;</code> - regular expression defined in <b>tokens.txt</b>.
 * <code>&lt;TOKEN:name&gt;</code> - regular expression defined in <b>tokens.txt</b> and a name that will be in output.
 * <code>specific_word</code> - simple word that will need to match.
* Words on each line needs to be separated by spaces.
* Sample pattern: <code>match 1##&lt;IP:sourceIPs&gt; &lt;DATE:date&gt; user &lt;USER:username&gt;</code>

Tokens.txt
-----------------------------
* One token definition per line like this: <code>token_name(space)regular_expression</code>.
* The syntax of the regular expressions accepted is the same general syntax used by
Perl, Python, and other languages. 
More precisely, it is the syntax accepted by RE2 and described at http://code.google.com/p/re2/wiki/Syntax, except for \C.
* Comments are allowed on different lines than the ones containing token definitions <code>#comment</code>.
