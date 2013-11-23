Each line in patterns.txt corresponds to one match, that will be searched for.

You have three options when defining patterns, and you need to separate them with spaces:

#1 TOKEN (regular expression defined in tokens.txt) surrounded by <> (for example <IP>
#2 SPECIFIC WORD surrounded by {} (for example {cr020r01-3.sac.overture.com})
#3 ANY WORD for that you can type _ without {} or [] and the search for that word will be skipped

Some examples:
<IP> _ _ <DATE> {"GET}
{4.37.97.186} _ _ {[11/Mar/2004:13:12:54 -0800]}

---------------------------------------------------------------------------------------------------------------------
Tokens in tokens.txt needs to be defined like this:
NAME 'regular expression without quotes'

Some examples: 
WORD ^\w+$
NUMBER ^[0-9]+$

The syntax of the regular expressions accepted is the same general syntax used by Perl, Python, and other languages. 
More precisely, it is the syntax accepted by RE2 and described at http://code.google.com/p/re2/wiki/Syntax,
except for \C.