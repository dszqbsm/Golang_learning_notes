module testGoModules

go 1.23.0

require github.com/q1mi/hello v0.1.2-0.20210219092711-2ccfaddad6a3 // indirect
require overtime v0.0.0
replace overtime => ./overtime
