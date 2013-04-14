function async(cb) {
  setTimeout(function() {
    cb("hello");
  }, 1000);
}

async(function(msg) {
  console.log(msg);
});

