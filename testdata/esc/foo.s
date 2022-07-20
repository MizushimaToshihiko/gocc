.file 1 "testdata/esc/foo.go"
	.globl main
	.text
main:
	push %rbp
	mov %rsp, %rbp
	sub $144, %rsp
	.loc 1 3
	.loc 1 3
	.loc 1 3
	.loc 1 3
	.loc 1 3
# ND_VAR: bar
	mov bar@GOTPCREL(%rip), %rax
	mov %rax, %r10
	mov $0, %rax
	call *%r10
	add $0, %rsp
	mov %rax, %r10
	jmp .L.return.main
	.loc 1 3
	.loc 1 3
	.loc 1 3
# ND_VAR: bar
	mov bar@GOTPCREL(%rip), %rax
	mov %rax, %r10
	mov $0, %rax
	call *%r10
	add $0, %rsp
	mov $0, %rax
	jmp .L.return.main
.L.return.main:
	mov %rbp, %rsp
	pop %rbp
	ret
