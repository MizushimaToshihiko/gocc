package main

func evalFuncall(node *Node) int64 {
	addType(node)

	enterScope()
	defer leaveScope()

	fn := node.Lhs.Obj
	createParamLvars(fn.Ty.Params)
	lvars := fn.Locals
	for lv := lvars; lv != nil; lv = lv.Next {
		pushScope(lv.Name).Obj = lv
	}

	// evaluate arguments
	for arg := node.Args; node.Args != nil; arg = arg.Next {
		findVar(lvars.Ty.Name).Obj.Val = eval(arg)
		lvars = lvars.Next
	}

	node2 := fn.Body
	for n := node2.Body; n != nil; n = n.Next {
		if n.Kind == ND_RETURN {
			return eval(n.RetVals)
		}
	}
	return 0
}
