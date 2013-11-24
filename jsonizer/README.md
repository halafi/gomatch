CONFIGURATION
--------------------
	o PATTERNS.TXT
	Each line in patterns.txt corresponds to one match, that will be searched for.
	You have three options when defining patterns, and you need to separate them with spaces:
		TOKEN (regular expression defined in tokens.txt) surrounded by <> (e.g. <IP>)
		SPECIFIC WORD surrounded by {} (e.g. {cr020r01-3.sac.overture.com})
		ANYTHING for that you can type _ and search for that will match anything
	Example line: <IP> _ _ <DATE> {"GET}

	o TOKENS.TXT
	Tokens in tokens.txt needs to be defined on separate lines like this:
		NAME regex
	The syntax of the regular expressions accepted is the same general syntax used by
	Perl, Python, and other languages. 
	More precisely, it is the syntax accepted by RE2 and described at
	http://code.google.com/p/re2/wiki/Syntax, except for \C.

	o Make sure that there are no extra spaces or endlines in these files and that they are ANSI encoded.
