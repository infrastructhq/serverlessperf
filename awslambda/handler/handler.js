const fibonacci = require("fibonacci");

exports.perf = function(event, context, callback) {
  callback(null, {
    cpu_ms: fibonacci.iterate(1000).ms
  });
};
