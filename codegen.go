//
// code generator
//
package main

import (
	"fmt"
	"io"
	"unsafe"
)

type codeWriter struct {
	w   io.Writer
	err error
}

func (c *codeWriter) println(frmt string, a ...interface{}) {
	if c.err != nil {
		return
	}
	_, c.err = fmt.Fprintf(c.w, frmt, a...)
	_, c.err = fmt.Fprintln(c.w)
}

func (c *codeWriter) unreachable(frmt string, a ...interface{}) {
	if c.err == nil {
		c.err = fmt.Errorf(frmt, a...)
	} else {
		c.err = fmt.Errorf(c.err.Error()+"\n"+frmt, a)
	}
}

var i int = 1

func count() int {
	i++
	return i
}

const GP_MAX = 6
const FP_MAX = 8

var depth int
var argreg8 = []string{"%dil", "%sil", "%dl", "%cl", "%r8b", "%r9b"}
var argreg16 = []string{"%di", "%si", "%dx", "%cx", "%r8w", "%r9w"}
var argreg32 = []string{"%edi", "%esi", "%edx", "%ecx", "%r8d", "%r9d"}
var argreg64 = []string{"%rdi", "%rsi", "%rdx", "%rcx", "%r8", "%r9"}

var curFnInGen *Obj

func (c *codeWriter) push() {
	if c.err != nil {
		return
	}

	c.println("	push %%rax")
	depth++
}

func (c *codeWriter) pop(arg string) {
	if c.err != nil {
		return
	}

	c.println("	pop %s", arg)
	depth--
}

func (c *codeWriter) pushf() {
	if c.err != nil {
		return
	}

	c.println("	sub $8, %%rsp")
	c.println("	movsd %%xmm0, (%%rsp)")
	depth++
}

func (c *codeWriter) popf(reg int) {
	if c.err != nil {
		return
	}

	c.println("	movsd (%%rsp), %%xmm%d", reg)
	c.println("	add $8, %%rsp")
	depth--
}

// Round up `n` to the nearest multiple of `align`. For instance,
// alignTo(5,8) returns 8 and alignTo(11,8) returns 16.
func alignTo(n int, align int) int {
	return (n + align - 1) / align * align
}

// Pushes the given node's address to the stack
func (c *codeWriter) genAddr(node *Node) {
	if c.err != nil {
		return
	}

	switch node.Kind {
	case ND_VAR:
		//  Local variable
		if node.Obj.IsLocal {
			c.println("	lea %d(%%rbp), %%rax", node.Obj.Offset)
			return
		}

		// Here, we generate an absolute address of a function or global
		// variable, Even though they exist at a certain address at runtime,
		// their addresses are not known at a link-time for the following
		// two reasons.
		//
		//  - Address randomization: Executables are loaded to memory as a
		//    whole but it is not know what assress they are loaded to.
		//    Therefore, at link-time, relative address in the same
		//    executable (i.e. the distance between two functions in the
		//    same executable) is known, but the absolute address is not
		//    known.
		//
		//  - Dynamic linking: Dynamic shared objects (DSOs) or .so files
		//    are loaded to mwmory alongside an aexecutable at runtime and
		//    linked by the runtime loader in memory. We know nothing
		//    about address of global stuff that may be defined by DSOs
		//    until the runtime relocation is complete.
		//
		// In order to deal with the former case, we use RIP-relative
		// addressing, denoted by `(%rip)`. For the latter, we obtain an
		// address of a stuff that may be in a shared object file from the
		// Global Offset Table using `@GOTPCREL(%rip)` notation.

		// Function
		if node.Ty.Kind == TY_FUNC {
			if node.Obj.IsDef {
				c.println("	lea %s(%%rip), %%rax", node.Obj.Name)
			} else {
				c.println("	mov %s@GOTPCREL(%%rip), %%rax", node.Obj.Name)
			}
			return
		}

		// Global variable
		c.println("	lea %s(%%rip), %%rax", node.Obj.Name)
		return
	case ND_DEREF:
		c.genExpr(node.Lhs)
		return
	case ND_COMMA:
		c.genExpr(node.Lhs)
		c.genAddr(node.Rhs)
		return
	case ND_MEMBER:
		c.genAddr(node.Lhs)
		c.println("	add $%d, %%rax", node.Mem.Offset)
		return
	case ND_FUNCALL:
		if node.RetBuf != nil {
			c.genExpr(node)
			return
		}
	}

	c.unreachable(errorTok(node.Tok, "not a lvalue"))
}

func (c *codeWriter) load(ty *Type) {
	if c.err != nil {
		return
	}

	switch ty.Kind {
	case TY_ARRAY, TY_STRUCT, TY_FUNC:
		// If it is an array, do not attempt to load a value to the
		// register because in general we can't load an entire array to a
		// register. As a result, the result of an evaluation of an array
		// become not the array itself but the address of the array.
		// This is where "array is automatically converted to a pointer to
		// the first element of the array in C" occurs.
		return
	case TY_FLOAT:
		c.println("	movss (%%rax), %%xmm0")
		return
	case TY_DOUBLE:
		c.println("	movsd (%%rax), %%xmm0")
		return
	}

	var insn string
	if ty.IsUnsigned {
		insn = "movz"
	} else {
		insn = "movs"
	}

	// When we load a char or a short value to a register, we always
	// extend them to the size of int, so we can assume the lower half of a
	// register for char, short and int may contain garbage. When we load
	// a long value to a register, it simply occupies the entire register.
	switch ty.Sz {
	case 1:
		c.println("	%sbl (%%rax), %%eax", insn)
	case 2:
		c.println("	%swl (%%rax), %%eax", insn)
	case 4:
		c.println("	movsxd (%%rax), %%rax")
	default:
		c.println("	mov (%%rax), %%rax")
		return
	}
}

func (c *codeWriter) store(ty *Type) {
	if c.err != nil {
		return
	}

	c.pop("%rdi")

	switch ty.Kind {
	case TY_STRUCT:
		for i := 0; i < ty.Sz; i++ {
			c.println("	mov %d(%%rax), %%r8b", i)
			c.println("	mov %%r8b, %d(%%rdi)", i)
		}
		return
	case TY_FLOAT:
		c.println("	movss %%xmm0, (%%rdi)")
		return
	case TY_DOUBLE:
		c.println("	movsd %%xmm0, (%%rdi)")
		return
	}

	switch ty.Sz {
	case 1:
		c.println("	mov %%al, (%%rdi)")
	case 2:
		c.println("	mov %%ax, (%%rdi)")
	case 4:
		c.println("	mov %%eax, (%%rdi)")
	default:
		c.println("	mov %%rax, (%%rdi)")
	}
}

func (c *codeWriter) cmpZero(ty *Type) {
	if c.err != nil {
		return
	}

	switch ty.Kind {
	case TY_FLOAT:
		c.println("	xorps %%xmm1, xmm1")
		c.println("	ucomiss %%xmm1, %%xmm0")
		return
	case TY_DOUBLE:
		c.println("	xorpd %%xmm1, %%xmm1")
		c.println("	ucomisd %%xmm1, %%xmm0")
		return
	}

	if isInteger(ty) && ty.Sz <= 4 {
		c.println("	cmp $0, %%eax")
	} else {
		c.println("	cmp $0, %%rax")
	}
}

const (
	I8 = iota
	I16
	I32
	I64
	U8
	U16
	U32
	U64
	F32
	F64
)

func (c *codeWriter) getTypeId(ty *Type) int {
	if c.err != nil {
		return -1
	}

	switch ty.Kind {
	case TY_BYTE:
		if ty.IsUnsigned {
			return U8
		}
		return I8
	case TY_SHORT:
		if ty.IsUnsigned {
			return U16
		}
		return I16
	case TY_INT:
		if ty.IsUnsigned {
			return U32
		}
		return I32
	case TY_LONG:
		if ty.IsUnsigned {
			return U64
		}
		return I64
	case TY_FLOAT:
		return F32
	case TY_DOUBLE:
		return F64
	default:
		return U64
	}
}

// The table for type casts
const (
	i32i8  string = "movsbl %al, %eax"
	i32u8  string = "movzbl %al, %eax"
	i32i16 string = "movswl %ax, %eax"
	i32u16 string = "movzwl %ax, %eax"
	i32f32 string = "cvtsi2ssl %eax, %xmm0"
	i32i64 string = "movsxd %eax, %rax"
	i32f64 string = "cvtsi2sdl %eax, %xmm0"

	u32f32 string = "mov %eax, %eax; cvtsi2ssq %rax, %xmm0"
	u32i64 string = "mov %eax, %eax"
	u32f64 string = "mov %eax, %eax; cvtsi2sdq %rax, %xmm0"

	i64f32 string = "cvtsi2ssq %rax, %xmm0"
	i64f64 string = "cvtsi2sdq %rax, %xmm0"

	u64f32 string = "cvtsi2ssq %rax, %xmm0"
	u64f64 string = "" +
		"test %rax,%rax; js 1f; pxor %xmm0,%xmm0; cvtsi2sd %rax,%xmm0; jmp 2f; " +
		"1: mov %rax,%rdi; and $1,%eax; pxor %xmm0,%xmm0; shr %rdi; " +
		"or %rax,%rdi; cvtsi2sd %rdi,%xmm0; addsd %xmm0,%xmm0; 2:"

	f32i8  string = "cvttss2sil %xmm0, %eax; movsbl %al, %eax"
	f32u8  string = "cvttss2sil %xmm0, %eax; movzbl %al, %eax"
	f32i16 string = "cvttss2sil %xmm0, %eax; movswl %ax, %eax"
	f32u16 string = "cvttss2sil %xmm0, %eax; movzwl %ax, %eax"
	f32i32 string = "cvttss2sil %xmm0, %eax"
	f32u32 string = "cvttss2siq %xmm0, %rax"
	f32i64 string = "cvttss2siq %xmm0, %rax"
	f32u64 string = "cvttss2siq %xmm0, %rax"
	f32f64 string = "cvtss2sd %xmm0, %xmm0"

	f64i8  string = "cvttsd2sil %xmm0, %eax; movsbl %al, %eax"
	f64u8  string = "cvttsd2sil %xmm0, %eax; movzbl %al, %eax"
	f64i16 string = "cvttsd2sil %xmm0, %eax; movswl %ax, %eax"
	f64u16 string = "cvttsd2sil %xmm0, %eax; movzwl %ax, %eax"
	f64i32 string = "cvttsd2sil %xmm0, %eax"
	f64u32 string = "cvttsd2siq %xmm0, %rax"
	f64f32 string = "cvtsd2ss %xmm0, %xmm0"
	f64i64 string = "cvttsd2siq %xmm0, %rax"
	f64u64 string = "cvttsd2siq %xmm0, %rax"
)

var castTable = [10][10]string{
	//i8    i16     i32     i64     u8     u16     u32     u64     f32     f64
	{"nil", "null", "null", i32i64, i32u8, i32u16, "null", i32i64, i32f32, i32f64}, // i8
	{i32i8, "null", "null", i32i64, i32u8, i32u16, "null", i32i64, i32f32, i32f64}, // i16
	{i32i8, i32i16, "null", i32i64, i32u8, i32u16, "null", i32i64, i32f32, i32f64}, // i32
	{i32i8, i32i16, "null", "null", i32u8, i32u16, "null", "null", i64f32, i64f64}, // i64

	{i32i8, "null", "null", i32i64, "nil", "null", "null", i32i64, i32f32, i32f64}, // u8
	{i32i8, i32i16, "null", i32i64, i32u8, "null", "null", i32i64, i32f32, i32f64}, // u16
	{i32i8, i32i16, "null", u32i64, i32u8, i32u16, "null", u32i64, u32f32, u32f64}, // u32
	{i32i8, i32i16, "null", "null", i32u8, i32u16, "null", "null", u64f32, u64f64}, // u64

	{f32i8, f32i16, f32i32, f32i64, f32u8, f32u16, f32u32, f32u64, "null", f32f64}, // f32
	{f64i8, f64i16, f64i32, f64i64, f64u8, f64u16, f64i32, f64u64, f64f32, "null"}, // f64
}

func (c *codeWriter) cast(from *Type, to *Type) {
	if c.err != nil {
		return

	}

	if to.Kind == TY_VOID {
		return
	}

	if to.Kind == TY_BOOL {
		c.cmpZero(from)
		c.println("	setne %%al")
		c.println("	movzx %%al, %%eax")
		return
	}

	t1 := c.getTypeId(from)
	t2 := c.getTypeId(to)
	if castTable[t1][t2] != "nil" && castTable[t1][t2] != "null" {
		c.println("	%s", castTable[t1][t2])
	}
}

// Struct equal or smaller than 16 bytes are passed
// using up to two registers.
//
// If the first 8 bytes contains only floating-point type members,
// they are passed in an XMM register. Otherwise, they are passed
// in a general-purpose register.
//
// If a struct is larger than 8 bytes, the same rule is
// applied to the next 8 bytes chunk.
//
// This function returns true if `ty` has only floating-point
// member in its byte tange [lo hi).
func hasFlonum(ty *Type, lo int, hi int, offset int) bool {
	if ty.Kind == TY_STRUCT {
		for mem := ty.Mems; mem != nil; mem = mem.Next {
			if !hasFlonum(mem.Ty, lo, hi, offset+mem.Offset) {
				return false
			}
		}
		return true
	}

	if ty.Kind == TY_ARRAY {
		for i := 0; i < ty.ArrSz; i++ {
			if !hasFlonum(ty.Base, lo, hi, offset+ty.Base.Sz*i) {
				return false
			}
		}
		return true
	}

	return offset < lo || hi <= offset || isFlonum(ty)
}

func hasFlonum1(ty *Type) bool {
	return hasFlonum(ty, 0, 8, 0)
}

func hasFlonum2(ty *Type) bool {
	return hasFlonum(ty, 8, 16, 0)
}

func (c *codeWriter) pushStruct(ty *Type) {
	sz := alignTo(ty.Sz, 8)
	c.println("	sub $%d, %%rsp", sz)
	depth += sz / 8

	for i := 0; i < ty.Sz; i++ {
		c.println("	mov %d(%%rax), %%r10b", i)
		c.println("	mov %%r10b, %d(%%rsp)", i)
	}
}

func (c *codeWriter) pushArgs2(args *Node, firstPass bool) {
	if c.err != nil {
		return
	}

	if args == nil {
		return
	}

	c.pushArgs2(args.Next, firstPass)

	if (firstPass && !args.PassByStack) ||
		(!firstPass && args.PassByStack) {
		return
	}

	c.genExpr(args)

	switch args.Ty.Kind {
	case TY_STRUCT:
		c.pushStruct(args.Ty)
	case TY_FLOAT, TY_DOUBLE:
		c.pushf()
	default:
		c.push()
	}
}

// pushArgs loads function call arguments. Arguments are already evaluated and
// stored to the stack as local variables. What we need to do in this
// function is to load them to registers or push them to the stack as
// specified by the x86-64 psABI. Here is what the spec says:
//
//  - Up to 6 arguments of integral type are passed using RDI, RSI,
//    RDX, RCX, R8 and R9.
//
//  - UP tp 8 arguments of floating-point type are passed using XMM0 to
//    XMM7.
//
//  - If all registers of an appropriate type are already used, push an
//    argument to the stack in the right-to-left order.
//
//  - Each argument passed on the stack takes 8 bytes, and the end of
//    the argument area must be aligned to a 16 bytes boundary.
//
//  - If a function is variadic, set the number of floating-point type
//    arguments to RAX.
func (c *codeWriter) pushArgs(node *Node) int {
	if c.err != nil {
		return -1
	}

	var stack, gp, fp int

	// If the return type is a large struct, the caller passes
	// a pointer to a buffer as if it were the first argument.
	if node.RetBuf != nil && node.Ty.Sz > 16 {
		gp++
	}

	// Load as many arguments to the registers as possible.
	for arg := node.Args; arg != nil; arg = arg.Next {
		ty := arg.Ty

		switch ty.Kind {
		case TY_STRUCT:
			if ty.Sz > 16 {
				arg.PassByStack = true
				stack += alignTo(ty.Sz, 8) / 8
			} else {
				var fp1, fp2 int
				var notfp1, notfp2 int = 1, 1
				if hasFlonum1(ty) {
					fp1 = 1
					notfp1 = 0
				}
				if hasFlonum2(ty) {
					fp2 = 1
					notfp2 = 0
				}

				if fp+fp1+fp2 < FP_MAX && gp+notfp1+notfp2 < GP_MAX {
					fp = fp + fp1 + fp2
					gp = gp + notfp1 + notfp2
				} else {
					arg.PassByStack = true
					stack += alignTo(ty.Sz, 8) / 8
				}
			}
		case TY_FLOAT, TY_DOUBLE:
			if fp >= FP_MAX {
				arg.PassByStack = true
				stack++
			}
			fp++
		default:
			if gp >= GP_MAX {
				arg.PassByStack = true
				stack++
			}
			gp++
		}
	}

	if (depth+stack)%2 == 1 {
		c.println("	sub $8, %%rsp")
		depth++
		stack++
	}

	c.pushArgs2(node.Args, true)
	c.pushArgs2(node.Args, false)

	// If the return type is a large struct, the caller passes
	// a pointer to a buffer as if it were the first argument.
	if node.RetBuf != nil && node.Ty.Sz > 16 {
		c.println("	lea %d(%%rbp), %%rax", node.RetBuf.Offset)
		c.push()
	}

	return stack
}

func (c *codeWriter) copyRetBuf(v *Obj) {
	ty := v.Ty
	var gp, fp int

	if hasFlonum1(ty) {
		if ty.Sz != 4 && 8 > ty.Sz {
			c.unreachable("internal error")
			return
		}

		if ty.Sz == 4 {
			c.println("	movss %%xmm0, %d(%%rbp)", v.Offset)
		} else {
			c.println("	movsd %%xmm0, %d(%%rbp)", v.Offset)
		}
		fp++
	} else {
		for i := 0; i < min(8, ty.Sz); i++ {
			c.println("	mov %%al, %d(%%rbp)", v.Offset+i)
			c.println("	shr $8, %%rax")
		}
		gp++
	}

	if ty.Sz > 8 {
		if hasFlonum2(ty) {
			if ty.Sz != 12 && ty.Sz != 16 {
				c.unreachable("internal error")
				return
			}
			if ty.Sz == 12 {
				c.println("	movss %%xmm%d, %d(%%rbp)", fp, v.Offset+8)
			} else {
				c.println("	movsd %%xmm%d, %d(%%rbp)", fp, v.Offset+8)
			}
		} else {
			var reg1 string = "%al"
			var reg2 string = "%rax"
			if gp != 0 {
				reg1 = "%dl"
				reg2 = "%rdx"
			}
			for i := 8; i < min(16, ty.Sz); i++ {
				c.println("	mov %s, %d(%%rbp)", reg1, v.Offset+i)
				c.println("	shr $8, %s", reg2)
			}
		}
	}
}

func (c *codeWriter) copyStructReg() {
	ty := curFnInGen.Ty.RetTy
	var gp, fp int

	c.println("	mov %%rax, %%rdi")

	if hasFlonum(ty, 0, 8, 0) {
		if ty.Sz != 4 && 8 > ty.Sz {
			c.unreachable("internal error")
			return
		}
		if ty.Sz == 4 {
			c.println("	movss (%%rdi), %%xmm0")
		} else {
			c.println("	movsd (%%rdi), %%xmm0")
		}
		fp++
	} else {
		c.println("	mov $0, %%rax")
		for i := min(8, ty.Sz) - 1; i >= 0; i-- {
			c.println("	shl $8, %%rax")
			c.println("	mov %d(%%rdi), %%al", i)
		}
		gp++
	}

	if ty.Sz > 8 {
		if hasFlonum(ty, 8, 16, 0) {
			if ty.Sz != 12 && ty.Sz != 16 {
				c.unreachable("internal error")
				return
			}
			if ty.Sz == 12 {
				c.println("	movss 8(%%rdi), %%xmm%d", fp)
			} else {
				c.println("	movsd 8(%%rdi), %%xmm%d", fp)
			}
		} else {
			var reg1, reg2 string = "%al", "%rax"
			if gp != 0 {
				reg1 = "%dl"
				reg2 = "%rdx"
			}
			c.println("	mov $0, %s", reg2)
			for i := min(16, ty.Sz) - 1; i >= 8; i-- {
				c.println("	shl $8, %s", reg2)
				c.println("	mov %d(%%rdi), %s", i, reg1)
			}
		}
	}
}

func (c *codeWriter) copyStructMem() {
	ty := curFnInGen.Ty.RetTy
	v := curFnInGen.Params

	c.println("	mov %d(%%rbp), %%rdi", v.Offset)

	for i := 0; i < ty.Sz; i++ {
		c.println("	mov %d(%%rax), %%dl", i)
		c.println("	mov %%dl, %d(%%rdi)", i)
	}
}

func (c *codeWriter) genExpr(node *Node) {
	if c.err != nil {
		return
	}

	c.println("	.loc 1 %d", node.Tok.LineNo)

	switch node.Kind {
	case ND_NULL_EXPR:
		return
	case ND_NUM:
		// c.println("# ND_NUM")
		switch node.Ty.Kind {
		case TY_FLOAT:
			f32 := node.FVal
			c.println("	mov $%d, %%eax  # float %f", *(*uint32)(unsafe.Pointer(&f32)), f32)
			c.println("	movq %%rax, %%xmm0")
			return
		case TY_DOUBLE:
			f64 := node.FVal
			c.println("	mov $%d, %%rax  # double %f", *(*uint64)(unsafe.Pointer(&f64)), f64)
			c.println("	movq %%rax, %%xmm0")
			return
		}

		c.println("	mov $%d, %%rax", node.Val)
		return
	case ND_NEG:
		// c.println("# ND_NEG")
		c.genExpr(node.Lhs)

		switch node.Ty.Kind {
		case TY_FLOAT:
			c.println("	mov $1, %%rax")
			c.println("	shl $31, %%rax")
			c.println("	movq %%rax, %%xmm1")
			c.println("	xorps %%xmm1, %%xmm0")
			return
		case TY_DOUBLE:
			c.println("	mov $1, %%rax")
			c.println("	shl $63, %%rax")
			c.println("	movq %%rax, %%xmm1")
			c.println("	xorpd %%xmm1, %%xmm0")
			return
		}

		c.println("	neg %%rax")
		return
	case ND_VAR, ND_MEMBER:
		// c.println("# ND_VAR or ND_MEMBER")
		c.genAddr(node)
		c.load(node.Ty)
		return
	case ND_DEREF:
		// c.println("# ND_DEREF")
		c.genExpr(node.Lhs)
		c.load(node.Ty)
		return
	case ND_ADDR:
		// c.println("# ND_ADDR")
		c.genAddr(node.Lhs)
		return
	case ND_ASSIGN:
		// c.println("# ND_ASSIGN")
		c.genAddr(node.Lhs)
		c.push()
		c.genExpr(node.Rhs)
		c.store(node.Ty)
		return
	case ND_STMT_EXPR:
		// c.println("# ND_STMT_EXPR")
		for n := node.Body; n != nil; n = n.Next {
			c.genStmt(n)
		}
		return
	case ND_COMMA:
		// c.println("# ND_COMMA")
		c.genExpr(node.Lhs)
		c.genExpr(node.Rhs)
		return
	case ND_CAST:
		// c.println("# ND_CAST")
		c.genExpr(node.Lhs)
		c.cast(node.Lhs.Ty, node.Ty)
		return
	case ND_MEMZERO:
		// c.println("# ND_MEMZERO")
		// `rep stosb` is equivalent to `memset(%rdi, %al, %rcx)`.
		c.println("	mov $%d, %%rcx", node.Obj.Ty.Sz)
		c.println("	lea %d(%%rbp), %%rdi", node.Obj.Offset)
		c.println("	mov $0, %%al")
		c.println("	rep stosb")
		return
	case ND_COND:
		// c.println("# ND_COND")
		cnt := count()
		c.genExpr(node.Cond)
		c.cmpZero(node.Cond.Ty)
		c.println("	je .L.else.%d", cnt)
		c.genExpr(node.Then)
		c.println("	jmp .L.end.%d", cnt)
		c.println(".L.else.%d", cnt)
		c.genExpr(node.Els)
		c.println(".L.end.%d", cnt)
		return
	case ND_NOT:
		// c.println("# ND_NOT")
		c.genExpr(node.Lhs)
		c.cmpZero(node.Lhs.Ty)
		c.println("	sete %%al")
		c.println("	movzx %%al, %%rax")
		return
	case ND_BITNOT:
		// c.println("# ND_BITNOT")
		c.genExpr(node.Lhs)
		c.println("	not %%rax")
		return
	case ND_LOGAND:
		// c.println("# ND_BITAND")
		cnt := count()
		c.genExpr(node.Lhs)
		c.cmpZero(node.Lhs.Ty)
		c.println("	je  .L.false.%d", cnt)
		c.genExpr(node.Rhs)
		c.cmpZero(node.Rhs.Ty)
		c.println("	je  .L.false.%d", cnt)
		c.println("	mov $1, %%rax")
		c.println("	jmp  .L.end.%d", cnt)
		c.println(".L.false.%d:", cnt)
		c.println("	mov $0, %%rax")
		c.println(".L.end.%d:", cnt)
		return
	case ND_LOGOR:
		// c.println("# ND_LOGOR")
		cnt := count()
		c.genExpr(node.Lhs)
		c.cmpZero(node.Lhs.Ty)
		c.println("	jne  .L.true.%d", cnt)
		c.genExpr(node.Rhs)
		c.cmpZero(node.Rhs.Ty)
		c.println("	jne  .L.true.%d", cnt)
		c.println("	mov $0, %%rax")
		c.println("	jmp  .L.end.%d", cnt)
		c.println(".L.true.%d:", cnt)
		c.println("	mov $1, %%rax")
		c.println(".L.end.%d:", cnt)
		return
	case ND_FUNCALL:
		// c.println("# ND_FUNCALL")
		stackArgs := c.pushArgs(node)
		c.genExpr(node.Lhs)

		gp := 0
		fp := 0

		// If the return type is a large struct, the caller passes
		// a pointer to a buffer as if it were the first argument.
		if node.RetBuf != nil && node.Ty.Sz > 16 {
			c.pop(argreg64[gp])
			gp++
		}

		for arg := node.Args; arg != nil; arg = arg.Next {
			ty := arg.Ty

			switch ty.Kind {
			case TY_STRUCT:
				if ty.Sz > 16 {
					continue
				}

				var fp1, fp2 int
				var notfp1, notfp2 int = 1, 1
				if hasFlonum1(ty) {
					fp1 = 1
					notfp1 = 0
				}
				if hasFlonum2(ty) {
					fp2 = 1
					notfp2 = 0
				}

				if fp+fp1+fp2 < FP_MAX && gp+notfp1+notfp2 < GP_MAX {
					if fp1 != 0 {
						c.popf(fp)
						fp++
					} else {
						c.pop(argreg64[gp])
						gp++
					}

					if ty.Sz > 8 {
						if fp2 != 0 {
							c.popf(fp)
							fp++
						} else {
							c.pop(argreg64[gp])
							gp++
						}
					}
				}
			case TY_FLOAT, TY_DOUBLE:
				if fp < FP_MAX {
					c.popf(fp)
					fp++
				}
			default:
				if gp < GP_MAX {
					c.pop(argreg64[gp])
					gp++
				}
			}
		}

		c.println("	mov %%rax, %%r10")
		c.println("	mov $%d, %%rax", fp)
		c.println("	call *%%r10")
		c.println("	add $%d, %%rsp", stackArgs*8)

		depth -= stackArgs

		// It looks like the most significant 48 or 56 bits int RAX may
		// contain garbage if a function return type is short or bool/char,
		// respectively. We clear the upper bits here.
		switch node.Ty.Kind {
		case TY_BOOL:
			c.println("	movzx %%al, %%eax")
			return
		case TY_BYTE:
			if node.Ty.IsUnsigned {
				c.println("	movzbl %%al, %%eax")
				return
			}
			c.println("	movsbl %%al, %%eax")
			return
		case TY_SHORT:
			if node.Ty.IsUnsigned {
				c.println("	movzwl %%ax, %%eax")
				return
			}
			c.println("	movswl %%ax, %%eax")
			return
		}

		// If the return type is a small struct, a value is returned
		// using up to two registers.
		if node.RetBuf != nil && node.Ty.Sz <= 16 {
			c.copyRetBuf(node.RetBuf)
			c.println("	lea %d(%%rbp), %%rax", node.RetBuf.Offset)
		}

		return
	}

	if isFlonum(node.Lhs.Ty) {
		c.genExpr(node.Rhs)
		c.pushf()
		c.genExpr(node.Lhs)
		c.popf(1)

		var sz string
		if node.Lhs.Ty.Kind == TY_FLOAT {
			sz = "ss"
		} else {
			sz = "sd"
		}

		switch node.Kind {
		case ND_ADD:
			// c.println("# ND_ADD is_flonum")
			c.println("	add%s %%xmm1, %%xmm0", sz)
			return
		case ND_SUB:
			// c.println("# ND_SUB is_flonum")
			c.println("	sub%s %%xmm1, %%xmm0", sz)
			return
		case ND_MUL:
			// c.println("# ND_MUL is_flonum")
			c.println("	mul%s %%xmm1, %%xmm0", sz)
			return
		case ND_DIV:
			// c.println("# ND_DIV is_flonum")
			c.println("	div%s %%xmm1, %%xmm0", sz)
			return
		case ND_EQ, ND_NE, ND_LT, ND_LE:
			// c.println("# ND_EQ or ND_NE or ND_LT or ND_LE")
			c.println("	ucomi%s %%xmm0, %%xmm1", sz)

			switch node.Kind {
			case ND_EQ:
				// c.println("# ND_EQ")
				c.println("	sete %%al")
				c.println("	setnp %%dl")
				c.println("	and %%dl, %%al")
			case ND_NE:
				// c.println("# ND_NE")
				c.println("	setne %%al")
				c.println("	setp %%dl")
				c.println("	or %%dl, %%al")
			case ND_LT:
				// c.println("# ND_LT")
				c.println("	seta %%al")
			default:
				// c.println("# exept for ND_EQ, ND_NE, ND_LT")
				c.println("	setae %%al")
			}

			c.println("	and $1, %%al")
			c.println("	movzb %%al, %%rax")
			return
		default:
			c.unreachable(errorTok(node.Tok, "invalid expression"))
			return
		}
	}

	c.genExpr(node.Rhs)
	c.push()
	c.genExpr(node.Lhs)
	c.pop("%rdi")

	var ax, di, dx string

	if node.Lhs.Ty.Kind == TY_LONG || node.Lhs.Ty.Base != nil {
		ax = "%rax"
		di = "%rdi"
		dx = "%rdx"
	} else {
		ax = "%eax"
		di = "%edi"
		dx = "%edx"
	}

	switch node.Kind {
	case ND_ADD:
		// c.println("# ND_ADD")
		c.println("	add %s, %s", di, ax)
		return
	case ND_SUB:
		// c.println("# ND_SUB")
		c.println("	sub %s, %s", di, ax)
		return
	case ND_MUL:
		// c.println("# ND_MUL")
		c.println("	imul %s, %s", di, ax)
		return
	case ND_DIV, ND_MOD:
		// c.println("# ND_DIV or ND_MOD")
		if node.Ty.IsUnsigned {
			c.println("	mov $0, %s", dx)
			c.println("	div %s", di)
		} else {
			if node.Lhs.Ty.Sz == 8 {
				c.println("	cqo")
			} else {
				c.println("	cdq")
			}
			c.println("	idiv %s", di)
		}

		if node.Kind == ND_MOD {
			// c.println("# ND_MOD")
			c.println("	mov %%rdx, %%rax")
		}
		return
	case ND_BITAND:
		// c.println("# ND_BITAND")
		c.println("	and %s, %s", di, ax)
		return
	case ND_BITOR:
		// c.println("# ND_BITOR")
		c.println("	or %s, %s", di, ax)
		return
	case ND_BITXOR:
		// c.println("# ND_BITXOR")
		c.println("	xor %s, %s", di, ax)
		return
	case ND_EQ, ND_NE, ND_LT, ND_LE:
		// c.println("# ND_EQ or ND_NE or ND_LT or ND_LE")
		c.println("	cmp %s, %s", di, ax)

		switch node.Kind {
		case ND_EQ:
			// c.println("# ND_EQ")
			c.println("	sete %%al")
		case ND_NE:
			// c.println("# ND_NE")
			c.println("	setne %%al")
		case ND_LT:
			// c.println("# ND_LT")
			if node.Lhs.Ty.IsUnsigned {
				c.println("	setb %%al")
			} else {
				c.println("	setl %%al")
			}
		case ND_LE:
			// c.println("# ND_LE")
			if node.Lhs.Ty.IsUnsigned {
				c.println("	setbe %%al")
			} else {
				c.println("	setle %%al")
			}
		}

		c.println("	movzb %%al, %%rax")
		return
	case ND_SHL:
		// c.println("# ND_SHL")
		c.println("	mov %%rdi, %%rcx")
		c.println("	shl %%cl, %s", ax)
		return
	case ND_SHR:
		// c.println("# ND_SHR")
		c.println("	mov %%rdi, %%rcx")
		if node.Lhs.Ty.IsUnsigned {
			c.println("	shr %%cl, %s", ax)
		} else {
			c.println("	sar %%cl, %s", ax)
		}
		return
	}

	c.unreachable("invalid expression")
}

func (c *codeWriter) genStmt(node *Node) {
	if c.err != nil {
		return
	}

	c.println("	.loc 1 %d", node.Tok.LineNo)

	switch node.Kind {
	case ND_IF:
		// c.println("# ND_IF")
		cnt := count()
		if node.Init != nil {
			c.genStmt(node.Init)
		}
		c.genExpr(node.Cond)
		c.cmpZero(node.Cond.Ty)
		c.println("	je .L.else.%d", cnt)
		c.genStmt(node.Then)
		c.println("	jmp .L.end.%d", cnt)
		c.println(".L.else.%d:", cnt)
		if node.Els != nil {
			c.genStmt(node.Els)
		}
		c.println(".L.end.%d:", cnt)
		return
	case ND_FOR:
		// c.println("# ND_FOR")
		cnt := count()
		if node.Init != nil {
			c.genStmt(node.Init)
		}
		c.println(".L.begin.%d:", cnt)
		if node.Cond != nil {
			c.genExpr(node.Cond)
			c.cmpZero(node.Cond.Ty)
			c.println("	je %s", node.BrkLabel)
		}
		c.genStmt(node.Then)
		c.println("%s:", node.ContLabel)
		if node.Inc != nil {
			c.genExpr(node.Inc)
		}
		c.println("	jmp .L.begin.%d", cnt)
		c.println("%s:", node.BrkLabel)
		return
	case ND_SWITCH:
		// c.println("# ND_SWITCH")
		if node.Init != nil {
			c.genStmt(node.Init)
		}
		c.genExpr(node.Cond)

		var reg1, reg2 string
		if node.Cond.Ty.Sz == 8 {
			reg1 = "%rdx"
			reg2 = "%rax"
		} else {
			reg1 = "%edx"
			reg2 = "%eax"
		}
		// Escape the value in reg2 to reg1.
		c.println("	mov %s, %s", reg2, reg1)

		for n := node.CaseNext; n != nil; n = n.CaseNext {

			c.genExpr(n.Expr)

			c.println("	cmp %s, %s", reg1, reg2)
			c.println("	je %s", n.Lbl)
			n.CaseEndLbl = node.BrkLabel
		}

		if node.DefCase != nil {
			node.DefCase.CaseEndLbl = node.BrkLabel
			c.println("	jmp %s", node.DefCase.Lbl)
		}

		c.println("	jmp %s", node.BrkLabel)
		c.genStmt(node.Then)
		c.println("%s:", node.BrkLabel)
		return
	case ND_CASE:
		// c.println("# ND_CASE")
		c.println("%s:", node.Lbl)
		if node.Lhs != nil {
			c.genStmt(node.Lhs)
			c.println("	jmp %s", node.CaseEndLbl)
		}
		return
	case ND_BLOCK:
		// c.println("# ND_BLOCK")
		for n := node.Body; n != nil; n = n.Next {
			c.genStmt(n)
		}
		return
	case ND_GOTO:
		// c.println("# ND_GOTO")
		c.println("	jmp %s", node.UniqueLbl)
		return
	case ND_LABEL:
		// c.println("# ND_LABEL")
		c.println("%s:", node.UniqueLbl)
		c.genStmt(node.Lhs)
		return
	case ND_MULTIVALASSIGN:

		for lhs := node.Lhses; lhs != nil; lhs = lhs.Next {
			if lhs.Kind != ND_NULL_EXPR {
				c.println("# lhs: %s", lhs.Obj.Name)
				c.genAddr(lhs)
				c.push()
			}
		}

		for rhs := node.Rhses; rhs != nil; rhs = rhs.Next {
			if rhs.Obj != nil {
				c.println("# rhs: %s", rhs.Obj.Name)
			}
			c.genExpr(rhs)
			c.store(rhs.Ty)
		}
		return
	case ND_MULTIRETASSIGN:
		c.genExpr(node.Lhs)

		i := 0
		n := node.Masg
		for ; n != nil; n = n.Next {
			c.genAddr(n)
			c.push()
			if isFlonum(n.Ty) {
				c.println("	movq %s, %%xmm0", argreg64[i])
			} else {
				c.println("	mov %s, %%rax", argreg64[i])
			}
			c.store(n.Ty)
			i++
		}
		return
	case ND_RETURN:
		// c.println("# ND_RETURN")

		// Save passed-by-register return values to the stack.
		i := 0
		for ret := node.RetVals; ret != nil; ret = ret.Next {
			c.genExpr(ret)

			ty := ret.Ty
			if ty.Kind == TY_STRUCT {
				if ty.Sz <= 16 {
					c.copyStructReg()
				} else {
					c.copyStructMem()
				}
			} else {
				c.println("	mov %%rax, %s", argreg64[i])
				i++
			}
		}

		c.println("	jmp .L.return.%s", curFnInGen.Name)
		return
	case ND_EXPR_STMT:
		// c.println("# ND_EXPR_STMT")
		c.genExpr(node.Lhs)
		return
	}
	c.unreachable(errorTok(node.Tok, "invalid statement"))
}

// Assign offsets to local variables
func (c *codeWriter) assignLvarOffsets(prog *Obj) {
	if c.err != nil {
		return
	}

	for fn := prog; fn != nil; fn = fn.Next {
		if !fn.IsFunc {
			continue
		}

		// If a function has many parameters, some parameters are
		// inevitably passed by stack rather than by register.
		// The first passed-by-stack parameter resides at RBP+16.
		top := 16
		bottom := 0

		gp := 0
		fp := 0

		// Assign offsets to pass-by-stack parameters.
		for v := fn.Params; v != nil; v = v.Next {
			ty := v.Ty

			switch ty.Kind {
			case TY_STRUCT:
				if ty.Sz <= 16 {
					var fp1, fp2 int
					var notfp1, notfp2 int = 1, 1
					if hasFlonum(ty, 0, 8, 0) {
						fp1 = 1
						notfp1 = 0
					}
					if hasFlonum(ty, 8, 16, 8) {
						fp2 = 1
						notfp2 = 0
					}

					if fp+fp1+fp2 < FP_MAX && gp+notfp1+notfp2 < GP_MAX {
						fp = fp + fp1 + fp2
						gp = gp + notfp1 + notfp2
						continue
					}
				}
			case TY_FLOAT, TY_DOUBLE:
				if fp < FP_MAX {
					fp++
					continue
				}
				fp++
			default:
				if gp < GP_MAX {
					gp++
					continue
				}
				gp++
			}

			top = alignTo(top, 8)
			v.Offset = top
			top += v.Ty.Sz
		}

		// Assign offset to pass-by-register parameters and local variables.
		for v := fn.Locals; v != nil; v = v.Next {
			if v.Offset != 0 {
				continue
			}

			bottom += v.Ty.Sz
			bottom = alignTo(bottom, v.Align)
			v.Offset = -bottom
		}

		fn.StackSz = alignTo(bottom, 16)
	}
}

func (c *codeWriter) emitData(prog *Obj) {
	if c.err != nil {
		return
	}

	for v := prog; v != nil; v = v.Next {
		if v.IsFunc || !v.IsDef {
			continue
		}

		if v.IsStatic {
			c.println("	.local %s", v.Name)
		} else {
			c.println("	.globl %s", v.Name)
		}
		c.println("	.align %d", v.Align)

		if v.InitData != nil {
			c.println("	.data")
			c.println("%s:", v.Name)

			rel := v.Rel
			pos := 0
			for pos < v.Ty.Sz {
				if rel != nil && rel.Offset == pos {
					c.println("	.quad %s%+d", rel.Lbl, rel.Addend)
					rel = rel.Next
					pos += 8
					continue
				} else {
					c.println("	.byte %d", v.InitData[pos])
					pos++
				}
			}
			continue
		}

		// Initialize with 0(?).
		c.println("	.bss")
		c.println("%s:", v.Name)
		c.println("	.zero %d", v.Ty.Sz)
	}
}

func (c *codeWriter) storeFp(r int, offset int, sz int) {
	if c.err != nil {
		return
	}

	switch sz {
	case 4:
		c.println("	movss %%xmm%d, %d(%%rbp)", r, offset)
		return
	case 8:
		c.println("	movsd %%xmm%d, %d(%%rbp)", r, offset)
		return
	default:
		c.unreachable("internal error")
	}
}

func (c *codeWriter) storeGp(r, offset, sz int) {
	if c.err != nil {
		return
	}

	switch sz {
	case 1:
		c.println("	mov %s, %d(%%rbp)", argreg8[r], offset)
		return
	case 2:
		c.println("	mov %s, %d(%%rbp)", argreg16[r], offset)
		return
	case 4:
		c.println("	mov %s, %d(%%rbp)", argreg32[r], offset)
		return
	case 8:
		c.println("	mov %s, %d(%%rbp)", argreg64[r], offset)
		return
	default:
		for i := 0; i < sz; i++ {
			c.println("	mov %s, %d(%%rbp)", argreg8[r], offset+i)
			c.println("	shr $8, %s", argreg64[r])
		}
	}
}

func (c *codeWriter) emitText(prog *Obj) {
	if c.err != nil {
		return
	}

	for fn := prog; fn != nil; fn = fn.Next {
		if !fn.IsFunc || !fn.IsDef {
			continue
		}

		if fn.IsStatic {
			c.println("	.local %s", fn.Name)
		} else {
			c.println("	.globl %s", fn.Name)
		}

		c.println("	.text")
		c.println("%s:", fn.Name)
		curFnInGen = fn

		// Prologue
		c.println("	push %%rbp")
		c.println("	mov %%rsp, %%rbp")
		c.println("	sub $%d, %%rsp", int(fn.StackSz))

		// Save arg registers if function is variadic.
		if fn.VaArea != nil {
			var gp, fp int
			for v := fn.Params; v != nil; v = v.Next {
				if isFlonum(v.Ty) {
					fp++
				} else {
					gp++
				}
			}

			off := fn.VaArea.Offset

			// va_elem
			c.println("  movl $%d, %d(%%rbp)", gp*8, off)
			c.println("  movl $%d, %d(%%rbp)", fp*8+48, off+4)
			c.println("  movq %%rbp, %d(%%rbp)", off+16)
			c.println("  addq $%d, %d(%%rbp)", off+24, off+16)

			// __reg_save_area__
			c.println("  movq %%rdi, %d(%%rbp)", off+24)
			c.println("  movq %%rsi, %d(%%rbp)", off+32)
			c.println("  movq %%rdx, %d(%%rbp)", off+40)
			c.println("  movq %%rcx, %d(%%rbp)", off+48)
			c.println("  movq %%r8, %d(%%rbp)", off+56)
			c.println("  movq %%r9, %d(%%rbp)", off+64)
			c.println("  movsd %%xmm0, %d(%%rbp)", off+72)
			c.println("  movsd %%xmm1, %d(%%rbp)", off+80)
			c.println("  movsd %%xmm2, %d(%%rbp)", off+88)
			c.println("  movsd %%xmm3, %d(%%rbp)", off+96)
			c.println("  movsd %%xmm4, %d(%%rbp)", off+104)
			c.println("  movsd %%xmm5, %d(%%rbp)", off+112)
			c.println("  movsd %%xmm6, %d(%%rbp)", off+120)
			c.println("  movsd %%xmm7, %d(%%rbp)", off+128)
		}

		// Push passed-by-register arguments to the stack
		gp := 0
		fp := 0
		for v := fn.Params; v != nil; v = v.Next {
			if v.Offset > 0 {
				continue
			}

			ty := v.Ty

			switch ty.Kind {
			case TY_STRUCT:
				if ty.Sz > 16 {
					c.unreachable("internal error")
					return
				}
				if hasFlonum(ty, 0, 8, 0) {
					c.storeFp(fp, v.Offset, min(8, ty.Sz))
					fp++
				} else {
					c.storeGp(gp, v.Offset, min(8, ty.Sz))
					gp++
				}

				if ty.Sz > 8 {
					if hasFlonum(ty, 8, 16, 0) {
						c.storeFp(fp, v.Offset+8, ty.Sz-8)
						fp++
					} else {
						c.storeGp(gp, v.Offset+8, ty.Sz-8)
						gp++
					}
				}
			case TY_FLOAT, TY_DOUBLE:
				c.storeFp(fp, v.Offset, ty.Sz)
				fp++
			default:
				c.storeGp(gp, v.Offset, ty.Sz)
				gp++
			}
		}

		// Emit code
		c.genStmt(fn.Body)
		if depth != 0 {
			c.unreachable(fmt.Sprintf("expected depth is 0, but %d", depth))
			return
		}

		// 'main' function returns implicitly 0.
		// [https://www.sigbus.info/n1570#5.1.2.2.3p1] The C spec defines
		// a special rule for the main function. Reaching the end of the
		// main function is equivalent to returning 0, evan though the
		// behavior is undifined for the other functions.
		if fn.Name == "main" {
			c.println("	mov $0, %%rax")
			c.println("	jmp .L.return.%s", fn.Name)
		}

		// Epilogue
		c.println(".L.return.%s:", fn.Name)
		c.println("	mov %%rbp, %%rsp")
		c.println("	pop %%rbp")
		c.println("	ret")
	}
}

func codegen(w io.Writer, prog *Obj) error {
	c := &codeWriter{w: w}

	c.assignLvarOffsets(prog)
	c.emitData(prog)
	c.emitText(prog)

	return c.err
}
