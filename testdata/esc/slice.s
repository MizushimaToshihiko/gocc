.file 1 "testdata/slice.go"
	.local .L..7
	.align 1
	.data
.L..7:
	.byte 79
	.byte 75
	.byte 0
	.local .L..6
	.align 1
	.data
.L..6:
	.byte 37
	.byte 115
	.byte 0
	.local .L..5
	.align 1
	.data
.L..5:
	.byte 0
	.local .L..4
	.align 1
	.data
.L..4:
	.byte 0
	.local .L..3
	.align 1
	.data
.L..3:
	.byte 0
	.local g06
	.align 8
	.data
g06:
	.quad .L..3+0
	.quad .L..4+0
	.quad .L..5+0
	.local .L..2
	.align 1
	.data
.L..2:
	.byte 103
	.byte 104
	.byte 105
	.byte 0
	.local .L..1
	.align 1
	.data
.L..1:
	.byte 100
	.byte 101
	.byte 102
	.byte 0
	.local .L..0
	.align 1
	.data
.L..0:
	.byte 97
	.byte 98
	.byte 99
	.byte 0
	.local g05
	.align 1
	.data
g05:
	.quad .L..0+0
	.quad .L..1+0
	.quad .L..2+0
	.globl main
	.text
main:
	push %rbp
	mov %rsp, %rbp
	sub $144, %rsp
  movl $0, -136(%rbp)
  movl $48, -132(%rbp)
  movq %rbp, -120(%rbp)
  addq $-112, -120(%rbp)
  movq %rdi, -112(%rbp)
  movq %rsi, -104(%rbp)
  movq %rdx, -96(%rbp)
  movq %rcx, -88(%rbp)
  movq %r8, -80(%rbp)
  movq %r9, -72(%rbp)
  movsd %xmm0, -64(%rbp)
  movsd %xmm1, -56(%rbp)
  movsd %xmm2, -48(%rbp)
  movsd %xmm3, -40(%rbp)
  movsd %xmm4, -32(%rbp)
  movsd %xmm5, -24(%rbp)
  movsd %xmm6, -16(%rbp)
  movsd %xmm7, -8(%rbp)
	.loc 1 32
	.loc 1 349
	.loc 1 349
# ND_ASSIGN
	.loc 1 349
	.loc 1 349
	.loc 1 349
	.loc 1 349
	.loc 1 349
	mov $8, %rax
	push %rax
	.loc 1 349
	.loc 1 349
	mov $1, %rax
	movsxd %eax, %rax
	pop %rdi
	imul %rdi, %rax
	push %rax
	.loc 1 349
	.loc 1 349
	lea g06(%rip), %rax
	pop %rdi
	add %rdi, %rax
	push %rax
	.loc 1 349
	.loc 1 349
	.loc 1 349
	.loc 1 349
	.loc 1 349
	.loc 1 349
	.loc 1 349
	mov $8, %rax
	push %rax
	.loc 1 349
	.loc 1 349
	mov $1, %rax
	movsxd %eax, %rax
	pop %rdi
	imul %rdi, %rax
	push %rax
	.loc 1 349
	.loc 1 349
	lea g05(%rip), %rax
	pop %rdi
	add %rdi, %rax
	mov (%rax), %rax
	pop %rdi
	mov %rax, (%rdi)
	.loc 1 350
	.loc 1 350
	.loc 1 350
	.loc 1 350
	.loc 1 350
	.loc 1 350
	.loc 1 350
	.loc 1 350
	mov $8, %rax
	push %rax
	.loc 1 350
	.loc 1 350
	mov $1, %rax
	movsxd %eax, %rax
	pop %rdi
	imul %rdi, %rax
	push %rax
	.loc 1 350
	.loc 1 350
	lea g06(%rip), %rax
	pop %rdi
	add %rdi, %rax
	mov (%rax), %rax
	push %rax
	.loc 1 350
	.loc 1 350
	lea .L..6(%rip), %rax
	push %rax
	.loc 1 350
	mov println@GOTPCREL(%rip), %rax
	pop %rdi
	pop %rsi
	mov %rax, %r10
	mov $0, %rax
	call *%r10
	add $0, %rsp
	.loc 1 352
	.loc 1 352
	.loc 1 352
	.loc 1 352
	lea .L..7(%rip), %rax
	push %rax
	.loc 1 352
	mov println@GOTPCREL(%rip), %rax
	pop %rdi
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
