module github.com/himidori/pm

go 1.14

replace github.com/himidori/pm/db => ./db

replace github.com/himidori/pm/utils => ./utils

require (
	github.com/Difrex/gpg v0.0.0-20190524122925-075df532c02f // indirect
	github.com/atotto/clipboard v0.1.2
	github.com/fatih/color v1.9.0
	github.com/himidori/pm/db v0.0.0-00010101000000-000000000000
	github.com/himidori/pm/utils v0.0.0-00010101000000-000000000000
	github.com/ogier/pflag v0.0.1
	github.com/sirupsen/logrus v1.6.0 // indirect
)
