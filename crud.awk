{
	gsub("TYPE", type)
	gsub("LOWER", tolower(type))
	gsub("PLURAL", tolower(type) "s")
	print
}
