package domain

const (
	/*INTEGER = iota
	REAL
	CHAR
	STRING*/

	INTEGER         = "integer"
	REAL            = "real"
	CHAR            = "char"
	STRING          = "string"
	TEXTFILE        = "textfile"
	INTEGERINTERVAL = "integerInterval"
)

func IsValidDataType(tpe string) bool{
	return tpe == INTEGER || tpe == REAL || tpe == STRING || tpe == CHAR || tpe == TEXTFILE || tpe == INTEGERINTERVAL
}
