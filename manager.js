const { spawn } = require("child_process");
const path = require("path");
const { performance } = require("perf_hooks");
const numWorkers = 1; // Adjust based on your system resources and workload
const taskQueue = [];
const workers = [];

let myClient = null;

const command = {
  command: "new_client",
  proxy: "",
};

const command2 = {
  command: "do",
  url: "https://eoobxe7m89qj9cl.m.pipedream.net/",
  method: "GET",
  headers: "",
  body: "",
  clientID: myClient,
};

const command3 = {
  command: "do",
  url: "https://www.dalphan.dev/",
  method: "POST",
  headers: {
    "User-Agent": "Inside Peanut Butter",
  },
  body: "",
};

const bridgePath = path.join(__dirname + "/bridge/bridge.exe");

console.log(bridgePath);

for (let i = 0; i < numWorkers; i++) {
  const worker = spawn(bridgePath);
  worker.isIdle = true;

  const startTime = performance.now();
  worker.stdin.write(JSON.stringify(command) + "\n");

  worker.stdout.on("data", (data) => {
    console.log(data);
    console.log(`Worker ${i}: ${data}`);

    const parsedData = JSON.parse(data);

    console.log(parsedData);

    if (parsedData.command === "do") {
      const endTime = performance.now();
      const duration = endTime - startTime;
      console.log(`Request duration: ${duration} milliseconds`);
    }

    if (parsedData.command === "new_client") {
      myClient = parsedData.clientID;
      console.log(`Fetched clientID: [${myClient}]`);
      worker.stdin.write(
        JSON.stringify({
          command: "do",
          clientID: myClient,
          url: "https://eoobxe7m89qj9cl.m.pipedream.net/",
          body: "",
          headers: {
            "User-Agent": "Inside Peanut Butter",
          },
          method: "GET",
        }) + "\n"
      );
    }
  });

  worker.stderr.on("data", (data) => {
    console.error(`Error at Worker ${i}: ${data}`);
  });

  worker.on("spawn", () => {
    console.log("Worker has spawned! Running [bridge.exe]");
  });

  workers.push(worker);
}
