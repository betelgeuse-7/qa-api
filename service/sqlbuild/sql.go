package sqlbuild

import "github.com/Masterminds/squirrel"

type Builder struct {
	B squirrel.StatementBuilderType
}

func New() *Builder {
	b := &Builder{}
	b.B = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	return b
}
