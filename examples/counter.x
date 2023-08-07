fn main {
  var count = 0
  for i <- 0..2 {
    call(89)
  }
}

/*

b0: <global>

b1: <main>
  v0 = 0 <count>
  jmp L0 <for>

L0: <for>
  v1 = 0 <i>
  call b2
  v2 = v1 + 1 <i++>
  cmp v2 2
  jl call

b2: <call>
  v3 = 89
  call call
  
*/