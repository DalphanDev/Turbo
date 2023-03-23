const { spawn } = require("child_process");
const path = require("path");
const numWorkers = 1; // Adjust based on your system resources and workload
const taskQueue = [];
const workers = [];

let myClient = null;

const command = {
  command: "new_client",
  proxy: "207.90.213.151:15413:egvrca423:qhYCz8388o",
};

const command2 = {
  command: "do",
  url: "https://www.dalphan.dev/",
  method: "GET",
  headers: "",
  body: "",
  clientID: myClient,
};

const command3 = {
  command: "do",
  url: "https://www.dalphan.dev/",
  method: "POST",
  headers: "",
  body: "",
};

const bridgePath = path.join(__dirname + "/bridge/bridge.exe");

console.log(bridgePath);

for (let i = 0; i < numWorkers; i++) {
  const worker = spawn(bridgePath);
  worker.isIdle = true;

  worker.stdin.write(JSON.stringify(command) + "\n");

  worker.stdout.on("data", (data) => {
    console.log(data);
    console.log(`Worker ${i}: ${data}`);

    const parsedData = JSON.parse(data);

    if (parsedData.command === "new_client") {
      myClient = parsedData.clientID;
      console.log(`Fetched clientID: [${myClient}]`);
      worker.stdin.write(
        JSON.stringify({
          command: "do",
          clientID: myClient,
          url: "https://dalphan.myshopify.com/",
          body: "",
          headers: "",
          method: "GET",
        }) + "\n"
      );
    }
  });

  worker.stderr.on("data", (data) => {
    console.error(`Worker ${i}: ${data}`);
  });

  worker.on("spawn", () => {
    console.log("Worker has spawned! Running [bridge.exe]");
  });

  workers.push(worker);
}
