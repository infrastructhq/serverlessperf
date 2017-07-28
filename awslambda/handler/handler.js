const fibonacci = require("fibonacci");
const speedtest = require("speedtest-net");

exports.perf = function(event, context, callback) {
  speedtest({ maxTime: 2000 }).on("data", data => {
    callback(null, {
      cpu_ms: fibonacci.iterate(1000).ms,
      network_mbs: data.speeds.download
    });
  });
};
