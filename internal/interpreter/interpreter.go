package interpreter

import (
	"fmt"

	"github.com/shaftoe44/danlisp/internal/expr"
)

type Interpreter struct {
	environment map[string]interface{}
}

func NewInterpreter() Interpreter{
	return Interpreter{NewEnvironment()}
}

func (interpreter *Interpreter) Interpret(exprs []expr.Expr) (interface{}, error) {	
	var retval interface{}
	var err error
	for _, ex := range exprs {
		retval, err = interpreter.eval(ex)
		if err != nil {
			return nil, err
		}
	}
	return retval, nil
}

func (interpreter *Interpreter) eval(ex expr.Expr) (interface{}, error) {

	switch v := ex.(type) {
	case expr.Atom:
		return evalAtom(v), nil
	case expr.Seq:
		return interpreter.evalSeq(v)
	case expr.Symbol:
		return interpreter.evalSymbol(v)
	case expr.Def:
		return interpreter.evalDef(v)
	case expr.If:
		return interpreter.evalIf(v)
	}

	panic("Don't know how to eval this thing")
}

func evalAtom(ex expr.Atom) interface{} {
	return ex.Value
}

func (interpreter *Interpreter) evalSymbol(ex expr.Symbol) (interface{}, error) {
	val, ok := interpreter.environment[ex.Name]
	if !ok {
		return nil, fmt.Errorf("runtime error. Could not find symbol '%v'", ex.Name)
	}
	return val, nil
}

func (interpreter *Interpreter) evalSeq(ex expr.Seq) (interface{}, error) {
	symbol, err := interpreter.eval(ex.Exprs[0])
	if err != nil {
		return symbol, err
	}
	args := []interface{}{}
	for _, argex := range ex.Exprs[1:] {
		arg, _ := interpreter.eval(argex)
		args = append(args, arg)
	}
	applyer := symbol.(func(argv []interface{}) interface{})
	return applyer(args), nil
}

func (interpreter *Interpreter)  evalDef(ex expr.Def) (interface{}, error) {
	val, err := interpreter.eval(ex.Value)
	interpreter.environment[ex.Var.Name] = val
	return val, err
}

func (interpreter *Interpreter) evalIf(iff expr.If) (interface{}, error) {
	cond, _ := interpreter.eval(iff.Cond)
	var expr expr.Expr
	if isTruthy(cond) {
		expr = iff.TrueBranch
	} else {
		expr = iff.FalseBranch
	}
	return interpreter.eval(expr)
}

func NewEnvironment() map[string]interface{} {
	env := make(map[string]interface{})

	// Basic operators
	env["+"] = func(argv []interface{}) interface{} { return argv[0].(float64) + argv[1].(float64) }
	env["-"] = func(argv []interface{}) interface{} { return argv[0].(float64) - argv[1].(float64) }
	env["*"] = func(argv []interface{}) interface{} { return argv[0].(float64) * argv[1].(float64) }
	env["/"] = func(argv []interface{}) interface{} { return argv[0].(float64) / argv[1].(float64) }
	env["mod"] = func(argv []interface{}) interface{} { return float64(int(argv[0].(float64)) % int(argv[1].(float64))) }

	// Bitwise ops
	env["&"] = func(argv []interface{}) interface{} { return float64(int(argv[0].(float64)) & int(argv[1].(float64))) }
	env["|"] = func(argv []interface{}) interface{} { return float64(int(argv[0].(float64)) | int(argv[1].(float64))) }
	env["^"] = func(argv []interface{}) interface{} { return float64(int(argv[0].(float64)) ^ int(argv[1].(float64))) }
	env["&^"] = func(argv []interface{}) interface{} { return float64(int(argv[0].(float64)) &^ int(argv[1].(float64))) }
	env[">>"] = func(argv []interface{}) interface{} { return float64(int(argv[0].(float64)) >> int(argv[1].(float64))) }
	env["<<"] = func(argv []interface{}) interface{} { return float64(int(argv[0].(float64)) << int(argv[1].(float64))) }

	// Boleans
	env["="] = func(argv []interface{}) interface{} { return argv[0] == argv[1] }
	env["and"] = func(argv []interface{}) interface{} { return isTruthy(argv[0]) && isTruthy(argv[1]) }
	env["or"] = func(argv []interface{}) interface{} { return isTruthy(argv[0]) || isTruthy(argv[1]) }

	return env
}

func isTruthy(v interface{}) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return v != nil
}
