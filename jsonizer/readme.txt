README.TXT for Jsonizer v.0.1
=============================
Thank you for choosing to read this readme.

TABLE OF CONTENTS
-----------------
	MINIMUM REQUIREMENTS
	INSTALLATION
	CONFIGURATION

MINIMUM REQUIREMENTS
--------------------
	o For GO language system requirements go to: http://golang.org/doc/install#requirements
	o If you are able to run GO, you should be fine.

INSTALLATION
--------------------
	o github readme.MD

CONFIGURATION
--------------------
	o PATTERNS.TXT
	Each line in patterns.txt corresponds to one match, that will be searched for.
	You have three options when defining patterns, and you need to separate them with spaces:
	#1 TOKEN (regular expression defined in tokens.txt) surrounded by <> (for example <IP>)
	#2 SPECIFIC WORD surrounded by {} (for example {cr020r01-3.sac.overture.com})
	#3 ANYTHING for that you can type _ without {} or [] and the search for that word will be skipped
	Examples:<IP> _ _ <DATE> {"GET}(new_line){4.37.97.186} _ _ {[11/Mar/2004:13:12:54 -0800]}

	o TOKENS.TXT
	Tokens in tokens.txt needs to be defined on separate lines like this:
	NAME 'regular expression without quotes'	
	Examples: WORD ^\w+$(new_line)NUMBER ^[0-9]+$
	The syntax of the regular expressions accepted is the same general syntax used by Perl, Python, and other languages. 
	More precisely, it is the syntax accepted by RE2 and described at http://code.google.com/p/re2/wiki/Syntax,
	except for \C.

	o Make sure that there are no extra spaces or endlines in these files and that they are ANSI encoded.
