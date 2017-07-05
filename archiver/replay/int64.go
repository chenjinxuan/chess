package main

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	"strconv"
)

type Int64 int64

func (i Int64) register(L *lua.LState) {
	mt := L.NewTypeMetatable("int64")
	L.SetGlobal("int64", mt)
	// static attributes
	L.SetField(mt, "new", L.NewFunction(i.newInt64))
	// meta-methods
	L.SetFuncs(mt, map[string]lua.LGFunction{
		"__add":      i.add,
		"__sub":      i.sub,
		"__mul":      i.mul,
		"__div":      i.div,
		"__mod":      i.mod,
		"__unm":      i.unm,
		"__eq":       i.eq,
		"__lt":       i.lt,
		"__le":       i.le,
		"__tostring": i.tostring,
	})

	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"bxor": i.xor,
		"band": i.and,
		"bor":  i.or,
		"bnot": i.not,
		"shl":  i.lshift,
		"shr":  i.rshift,
	}))
}

func (i Int64) newInt64(L *lua.LState) int {
	v := L.CheckAny(1)
	switch v.(type) {
	case lua.LString:
		x, err := strconv.ParseInt(L.CheckString(1), 0, 64)
		if err == nil {
			ud := L.NewUserData()
			ud.Value = Int64(x)
			L.SetMetatable(ud, L.GetTypeMetatable("int64"))
			L.Push(ud)
			return 1
		} else {
			L.ArgError(1, err.Error())
			return 0
		}
	case lua.LNumber:
		ud := L.NewUserData()
		ud.Value = Int64(L.CheckNumber(1))
		L.SetMetatable(ud, L.GetTypeMetatable("int64"))
		L.Push(ud)
		return 1
	default:
		L.ArgError(1, "invalid datatype")
		return 0
	}
}

func (i Int64) add(L *lua.LState) int {
	return i.binop(L, func(x, y int64) int64 { return x + y })
}

func (i Int64) sub(L *lua.LState) int {
	return i.binop(L, func(x, y int64) int64 { return x - y })
}

func (i Int64) mul(L *lua.LState) int {
	return i.binop(L, func(x, y int64) int64 { return x * y })
}

func (i Int64) div(L *lua.LState) int {
	return i.binop(L, func(x, y int64) int64 { return x / y })
}

func (i Int64) mod(L *lua.LState) int {
	return i.binop(L, func(x, y int64) int64 { return x % y })
}

func (i Int64) unm(L *lua.LState) int {
	return i.unaryop(L, func(x int64) int64 { return -x })
}

func (i Int64) eq(L *lua.LState) int {
	return i.boolop(L, func(x, y int64) bool { return x == y })
}

func (i Int64) lt(L *lua.LState) int {
	return i.boolop(L, func(x, y int64) bool { return x < y })
}

func (i Int64) le(L *lua.LState) int {
	return i.boolop(L, func(x, y int64) bool { return x <= y })
}

func (i Int64) xor(L *lua.LState) int {
	return i.binop(L, func(x, y int64) int64 { return x ^ y })
}

func (i Int64) and(L *lua.LState) int {
	return i.binop(L, func(x, y int64) int64 { return x & y })
}

func (i Int64) or(L *lua.LState) int {
	return i.binop(L, func(x, y int64) int64 { return x | y })
}

func (i Int64) not(L *lua.LState) int {
	a := L.CheckUserData(1).Value.(Int64)
	ud := L.NewUserData()
	ud.Value = Int64(int64(a) ^ (1<<63 - 1))
	L.SetMetatable(ud, L.GetTypeMetatable("int64"))
	L.Push(ud)
	return 1
}

func (i Int64) lshift(L *lua.LState) int {
	a := L.CheckUserData(1).Value.(Int64)
	shift := uint(L.CheckInt(2))
	ud := L.NewUserData()
	ud.Value = Int64(int64(a) << shift)
	L.SetMetatable(ud, L.GetTypeMetatable("int64"))
	L.Push(ud)
	return 1
}

func (i Int64) rshift(L *lua.LState) int {
	a := L.CheckUserData(1).Value.(Int64)
	shift := uint(L.CheckInt(2))
	ud := L.NewUserData()
	ud.Value = Int64(int64(a) >> shift)
	L.SetMetatable(ud, L.GetTypeMetatable("int64"))
	L.Push(ud)
	return 1
}

func (i Int64) tostring(L *lua.LState) int {
	x := L.CheckUserData(1).Value.(Int64)
	L.Push(lua.LString(fmt.Sprint(x)))
	return 1
}

func (i Int64) binop(L *lua.LState, f func(x, y int64) int64) int {
	a := L.CheckUserData(1).Value.(Int64)
	b := L.CheckUserData(2).Value.(Int64)
	ud := L.NewUserData()
	ud.Value = Int64(f(int64(a), int64(b)))
	L.SetMetatable(ud, L.GetTypeMetatable("int64"))
	L.Push(ud)
	return 1
}

func (i Int64) boolop(L *lua.LState, f func(x, y int64) bool) int {
	a := L.CheckUserData(1).Value.(Int64)
	b := L.CheckUserData(2).Value.(Int64)
	L.Push(lua.LBool(f(int64(a), int64(b))))
	return 1
}

func (i Int64) unaryop(L *lua.LState, f func(x int64) int64) int {
	a := L.CheckUserData(1).Value.(Int64)
	ud := L.NewUserData()
	ud.Value = Int64(f(int64(a)))
	L.SetMetatable(ud, L.GetTypeMetatable("int64"))
	L.Push(ud)
	return 1
}
