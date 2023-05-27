const axios = require("axios");
const { performance } = require("perf_hooks");

const url = "https://eoobxe7m89qj9cl.m.pipedream.net";

const startTime = performance.now();

axios
  .get(url)
  .then((response) => {
    const endTime = performance.now();
    const duration = endTime - startTime;
    console.log("GET request successful!");
    console.log("Response:", response.data);
    console.log(`Request duration: ${duration} milliseconds`);
  })
  .catch((error) => {
    console.error("Error occurred:", error.message);
  });
