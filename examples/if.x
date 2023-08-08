fn main {
  var count = 0
  for _ <- 0..10 {
    if count > 5 {
      break
    } else {
      count = count + 1
    }
  }
}