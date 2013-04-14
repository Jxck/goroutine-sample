def generator(n):
  i = 0
  while True:
    if i > n: break
    yield i
    i += 1

for i in generator(10):
  print i
