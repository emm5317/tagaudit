package main

import (
	"github.com/emm5317/tagaudit"
	"github.com/emm5317/tagaudit/rules"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(tagaudit.NewAnalyzer(rules.DefaultConfig()))
}
