function f()
  coroutine.yield "one"
  coroutine.yield "two"
  coroutine.yield "three"
  return
end

local co = coroutine.wrap (f)

print (co ()) -- one
print (co ()) -- two
print (co ()) -- three
