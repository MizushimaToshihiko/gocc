.file 1 "testdata/alignof.go"
	.local .L..12
	.align 1
	.data
.L..12:
	.byte 79
	.byte 75
	.byte 0
	.local .L..11
	.align 1
	.data
.L..11:
	.byte 65
	.byte 108
	.byte 105
	.byte 103
	.byte 110
	.byte 111
	.byte 102
	.byte 40
	.byte 120
	.byte 41
	.byte 60
	.byte 60
	.byte 54
	.byte 51
	.byte 62
	.byte 62
	.byte 54
	.byte 51
	.byte 0
	.local .L..10
	.align 1
	.data
.L..10:
	.byte 65
	.byte 108
	.byte 105
	.byte 103
	.byte 110
	.byte 111
	.byte 102
	.byte 40
	.byte 120
	.byte 41
	.byte 60
	.byte 60
	.byte 51
	.byte 49
	.byte 62
	.byte 62
	.byte 51
	.byte 49
	.byte 0
	.local .L..9
	.align 1
	.data
.L..9:
	.byte 65
	.byte 108
	.byte 105
	.byte 103
	.byte 110
	.byte 111
	.byte 102
	.byte 40
	.byte 103
	.byte 57
	.byte 41
	.byte 0
	.local .L..8
	.align 1
	.data
.L..8:
	.byte 65
	.byte 108
	.byte 105
	.byte 103
	.byte 110
	.byte 111
	.byte 102
	.byte 40
	.byte 103
	.byte 56
	.byte 41
	.byte 0
	.local .L..7
	.align 1
	.data
.L..7:
	.byte 65
	.byte 108
	.byte 105
	.byte 103
	.byte 110
	.byte 111
	.byte 102
	.byte 40
	.byte 103
	.byte 55
	.byte 41
	.byte 0
	.local .L..6
	.align 1
	.data
.L..6:
	.byte 65
	.byte 108
	.byte 105
	.byte 103
	.byte 110
	.byte 111
	.byte 102
	.byte 40
	.byte 103
	.byte 54
	.byte 41
	.byte 0
	.local .L..5
	.align 1
	.data
.L..5:
	.byte 65
	.byte 108
	.byte 105
	.byte 103
	.byte 110
	.byte 111
	.byte 102
	.byte 40
	.byte 103
	.byte 53
	.byte 41
	.byte 0
	.local .L..4
	.align 1
	.data
.L..4:
	.byte 65
	.byte 108
	.byte 105
	.byte 103
	.byte 110
	.byte 111
	.byte 102
	.byte 40
	.byte 103
	.byte 52
	.byte 41
	.byte 0
	.local .L..3
	.align 1
	.data
.L..3:
	.byte 65
	.byte 108
	.byte 105
	.byte 103
	.byte 110
	.byte 111
	.byte 102
	.byte 40
	.byte 103
	.byte 51
	.byte 41
	.byte 0
	.local g9
	.align 1
	.data
g9:
	.byte 1
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 2
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.local g8
	.align 1
	.data
g8:
	.byte 1
	.byte 2
	.byte 0
	.byte 0
	.local .L..0
	.align 1
	.data
.L..0:
	.byte 97
	.byte 98
	.byte 99
	.byte 100
	.byte 101
	.byte 102
	.byte 0
	.local g7
	.align 1
	.data
g7:
	.quad .L..0+0
	.local g6
	.align 8
	.data
g6:
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.local g5
	.align 4
	.data
g5:
	.byte 0
	.byte 0
	.byte 0
	.byte 0
	.local g4
	.align 2
	.data
g4:
	.byte 0
	.byte 0
	.local g3
	.align 1
	.data
g3:
	.byte 0
	.globl main
	.text
main:
	push %rbp
	mov %rsp, %rbp
	sub $144, %rsp
  movl $0, -137(%rbp)
  movl $48, -133(%rbp)
  movq %rbp, -121(%rbp)
  addq $-113, -121(%rbp)
  movq %rdi, -113(%rbp)
  movq %rsi, -105(%rbp)
  movq %rdx, -97(%rbp)
  movq %rcx, -89(%rbp)
  movq %r8, -81(%rbp)
  movq %r9, -73(%rbp)
  movsd %xmm0, -65(%rbp)
  movsd %xmm1, -57(%rbp)
  movsd %xmm2, -49(%rbp)
  movsd %xmm3, -41(%rbp)
  movsd %xmm4, -33(%rbp)
  movsd %xmm5, -25(%rbp)
  movsd %xmm6, -17(%rbp)
  movsd %xmm7, -9(%rbp)
	.loc 1 21
	.loc 1 21
	.loc 1 21
	.loc 1 21
	.loc 1 21
	lea .L..3(%rip), %rax
	push %rax
	.loc 1 21
	.loc 1 21
	mov $1, %rax
	push %rax
	.loc 1 21
	.loc 1 21
	mov $1, %rax
	push %rax
	.loc 1 21
	mov assert@GOTPCREL(%rip), %rax
	pop %rdi
	pop %rsi
	pop %rdx
	mov %rax, %r10
	mov $0, %rax
	call *%r10
	add $0, %rsp
	.loc 1 22
	.loc 1 22
	.loc 1 22
	.loc 1 22
	lea .L..4(%rip), %rax
	push %rax
	.loc 1 22
	.loc 1 22
	mov $2, %rax
	push %rax
	.loc 1 22
	.loc 1 22
	mov $2, %rax
	push %rax
	.loc 1 22
	mov assert@GOTPCREL(%rip), %rax
	pop %rdi
	pop %rsi
	pop %rdx
	mov %rax, %r10
	mov $0, %rax
	call *%r10
	add $0, %rsp
	.loc 1 23
	.loc 1 23
	.loc 1 23
	.loc 1 23
	lea .L..5(%rip), %rax
	push %rax
	.loc 1 23
	.loc 1 23
	mov $4, %rax
	push %rax
	.loc 1 23
	.loc 1 23
	mov $4, %rax
	push %rax
	.loc 1 23
	mov assert@GOTPCREL(%rip), %rax
	pop %rdi
	pop %rsi
	pop %rdx
	mov %rax, %r10
	mov $0, %rax
	call *%r10
	add $0, %rsp
	.loc 1 24
	.loc 1 24
	.loc 1 24
	.loc 1 24
	lea .L..6(%rip), %rax
	push %rax
	.loc 1 24
	.loc 1 24
	mov $8, %rax
	push %rax
	.loc 1 24
	.loc 1 24
	mov $8, %rax
	push %rax
	.loc 1 24
	mov assert@GOTPCREL(%rip), %rax
	pop %rdi
	pop %rsi
	pop %rdx
	mov %rax, %r10
	mov $0, %rax
	call *%r10
	add $0, %rsp
	.loc 1 25
	.loc 1 25
	.loc 1 25
	.loc 1 25
	lea .L..7(%rip), %rax
	push %rax
	.loc 1 25
	.loc 1 25
	mov $8, %rax
	push %rax
	.loc 1 25
	.loc 1 25
	mov $8, %rax
	push %rax
	.loc 1 25
	mov assert@GOTPCREL(%rip), %rax
	pop %rdi
	pop %rsi
	pop %rdx
	mov %rax, %r10
	mov $0, %rax
	call *%r10
	add $0, %rsp
	.loc 1 26
	.loc 1 26
	.loc 1 26
	.loc 1 26
	lea .L..8(%rip), %rax
	push %rax
	.loc 1 26
	.loc 1 26
	mov $1, %rax
	push %rax
	.loc 1 26
	.loc 1 26
	mov $1, %rax
	push %rax
	.loc 1 26
	mov assert@GOTPCREL(%rip), %rax
	pop %rdi
	pop %rsi
	pop %rdx
	mov %rax, %r10
	mov $0, %rax
	call *%r10
	add $0, %rsp
	.loc 1 27
	.loc 1 27
	.loc 1 27
	.loc 1 27
	lea .L..9(%rip), %rax
	push %rax
	.loc 1 27
	.loc 1 27
	mov $8, %rax
	push %rax
	.loc 1 27
	.loc 1 27
	mov $8, %rax
	push %rax
	.loc 1 27
	mov assert@GOTPCREL(%rip), %rax
	pop %rdi
	pop %rsi
	pop %rdx
	mov %rax, %r10
	mov $0, %rax
	call *%r10
	add $0, %rsp
	.loc 1 29
	.loc 1 29
	.loc 1 29
	.loc 1 29
	mov $1, %rcx
	lea -1(%rbp), %rdi
	mov $0, %al
	rep stosb
	.loc 1 29
	lea -1(%rbp), %rax
	push %rax
	.loc 1 29
	.loc 1 29
	mov $0, %rax
	pop %rdi
	mov %al, (%rdi)
	.loc 1 30
	.loc 1 30
	.loc 1 30
	.loc 1 30
	lea .L..10(%rip), %rax
	push %rax
	.loc 1 30
	.loc 1 30
	.loc 1 30
	mov $31, %rax
	push %rax
	.loc 1 30
	.loc 1 30
	mov $31, %rax
	push %rax
	.loc 1 30
	mov $1, %rax
	pop %rdi
	mov %rdi, %rcx
	shl %cl, %rax
	pop %rdi
	mov %rdi, %rcx
	shr %cl, %rax
	push %rax
	.loc 1 30
	.loc 1 30
	mov $1, %rax
	push %rax
	.loc 1 30
	mov assert@GOTPCREL(%rip), %rax
	pop %rdi
	pop %rsi
	pop %rdx
	mov %rax, %r10
	mov $0, %rax
	call *%r10
	add $0, %rsp
	.loc 1 31
	.loc 1 31
	.loc 1 31
	.loc 1 31
	lea .L..11(%rip), %rax
	push %rax
	.loc 1 31
	.loc 1 31
	.loc 1 31
	mov $63, %rax
	push %rax
	.loc 1 31
	.loc 1 31
	mov $63, %rax
	push %rax
	.loc 1 31
	mov $1, %rax
	pop %rdi
	mov %rdi, %rcx
	shl %cl, %rax
	pop %rdi
	mov %rdi, %rcx
	shr %cl, %rax
	push %rax
	.loc 1 31
	.loc 1 31
	mov $1, %rax
	push %rax
	.loc 1 31
	mov assert@GOTPCREL(%rip), %rax
	pop %rdi
	pop %rsi
	pop %rdx
	mov %rax, %r10
	mov $0, %rax
	call *%r10
	add $0, %rsp
	.loc 1 33
	.loc 1 33
	.loc 1 33
	.loc 1 33
	lea .L..12(%rip), %rax
	push %rax
	.loc 1 33
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
