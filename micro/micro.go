package micro

import (
	"github.com/samber/lo"
	"golang.org/x/exp/constraints"
)

//三目表达式
func IF[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

func MAX[T constraints.Ordered](params ...T) T {
	return lo.Max(params)
}

func MIN[T constraints.Ordered](params ...T) T {
	return lo.Min(params)
}
