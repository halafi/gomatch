Each line in patterns.txt corresponds to one match, that will be searched for.

You have three options when defining patterns, and you need to separate them with spaces:

#1 TOKEN (regular expression defined in tokens.txt) surrounded by <> (for example <IP>
#2 SPECIFIC WORD surrounded by {} (for example {cr020r01-3.sac.overture.com})
#3 ANY WORD for that you can type anything without {} or [] and the search for that word will be skipped (for example _ or skip)

Some examples:
<IP> _ _ <DATE> {"GET}
{4.37.97.186} _ _ {[11/Mar/2004:13:12:54 -0800]}