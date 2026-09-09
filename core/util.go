package core

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
)

func DebugUnit(gUnit Unit) {
	debugUnit(gUnit, 0)
}

func PrintUnit(unit Unit) {
	fmt.Println(SprintUnit(unit))
}
func SprintUnit(unit Unit) string {
	switch unit.Type {
	case SchemeType:
		asScheme := unit.Value.(Scheme)
		var ret string = fmt.Sprintf("(%s", asScheme.Name)
		for i := range asScheme.Args {
			ret += fmt.Sprintf(" %s", SprintUnit(asScheme.Args[i]))
		}
		return ret + ")"
	case Symbol:
		return fmt.Sprintf("%s", unit.Value.(string))
	case Float:
		return fmt.Sprintf("%f", unit.Value.(float64))
	case Integer:
		return fmt.Sprintf("%d", unit.Value.(int64))
	case String:
		return fmt.Sprintf("%s", unit.Value.(string))
	case Bool:
		return fmt.Sprintf("%v", unit.Value.(bool))
	case Char:
		return fmt.Sprintf("%c", unit.Value.(rune))
	case Pair:
		asPair := unit.Value.(PairVal)
		return fmt.Sprintf("(%v . %v)", SprintUnit(asPair[0]), SprintUnit(asPair[1]))
	case Vector:
		asVec := unit.Value.(VectorVal)
		var ret string
		for i := range asVec {
			if i != 0 {
				ret = fmt.Sprintf("%s %s", ret, SprintUnit(asVec[i]))
			} else {
				ret = fmt.Sprintf("#(%s", SprintUnit(asVec[i]))
			}
		}
		return fmt.Sprintf("%s)", ret)
	default:
		panic(fmt.Sprintf("Unknown type %v", unit))
	}
}

func debugUnit(unit Unit, depth int) {
	switch unit.Type {
	case SchemeType:
		asGroup := unit.Value.(Scheme)
		log.Debugf("%sScheme '%s'", repeat(depth), asGroup.Name)
		for _, member := range asGroup.Args {
			debugUnit(member, depth+1)
		}
	case Symbol:
		log.Debugf("%sSymbol %s", repeat(depth), unit.Value.(string))
	case Float:
		log.Debugf("%sFloat %.1f", repeat(depth), unit.Value.(float64))
	case Integer:
		log.Debugf("%sInt %v", repeat(depth), unit.Value.(int64))
	case String:
		log.Debugf("%sString \"%s\"", repeat(depth), unit.Value.(string))
	case Char:
		log.Debugf("%sChar %c", repeat(depth), unit.Value.(rune))
	case Pair:
		asPair := unit.Value.(PairVal)
		log.Debugf("%sPair (%v . %v)", repeat(depth), SprintUnit(asPair[0]), SprintUnit(asPair[1]))
	case Vector:
		asVec := unit.Value.(VectorVal)
		var ret string
		for i := range asVec {
			if i != 0 {
				ret = fmt.Sprintf("%s %s", ret, SprintUnit(asVec[i]))
			} else {
				ret = fmt.Sprintf("#(%s", SprintUnit(asVec[i]))
			}
		}
		ret = fmt.Sprintf("%s)", ret)

		log.Debugf("%sVector %s", repeat(depth), ret)
	default:
		log.Warnf("%sUnknown %v", repeat(depth), unit.Value)
	}
}

func repeat(depth int) string {
	return strings.Repeat(". ", depth)
}