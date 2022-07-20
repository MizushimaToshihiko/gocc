.file 1 "testdata/esc/bar.go"
	.local bar
	.text
bar:
	push %rbp
	mov %rsp, %rbp
	sub $144, %rsp
	.loc 1 2
	.loc 1 2
	.loc 1 2
	.loc 1 2
	mov $42, %rax
	mov %rax, %r10
	jmp .L.return.bar
	.loc 1 2
	.loc 1 2
	mov $42, %rax
.L.return.bar:
	mov %rbp, %rsp
	pop %rbp
	ret
