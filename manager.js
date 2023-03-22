const { spawn } = require("child_process");
const path = require("path");
const numWorkers = 1; // Adjust based on your system resources and workload
const taskQueue = [];
const workers = [];

const command = {
  command: "new_client",
  proxy: "207.90.213.151:15413:egvrca423:qhYCz8388o",
};

const bridgePath = path.join(__dirname + "/bridge/bridge.exe");
console.log(bridgePath);

for (let i = 0; i < numWorkers; i++) {
  const worker = spawn(bridgePath);
  worker.isIdle = true;

  worker.stdin.write(JSON.stringify(command) + "\n");

  worker.stdout.on("data", (data) => {
    console.log(`Worker ${i}: ${data}`);
  });

  worker.stderr.on("data", (data) => {
    console.error(`Worker ${i}: ${data}`);
  });

  worker.on("spawn", () => {
    console.log("Worker has spawned! Running [bridge.exe]");
  });

  workers.push(worker);
}
