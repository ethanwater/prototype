var socket = new WebSocket("wss://localhost:2695/ws", "protocolOne");
var call = 0;

// Create the chart instance
var acquisitionsChart = new Chart(document.getElementById('acquisitions'), {
  type: 'bar',
  data: {
    labels: [],
    datasets: [{
      label: 'socket calls',
      data: []
    }]
  }
});

socket.onopen = function(event) {
  console.log("WebSocket connection established.");
};

socket.onmessage = function(event) {
  call++;
  const data = [
    { deployment: 1, count: call },
    { deployment: 2, count: 100 },
  ];

  // Update chart data
  acquisitionsChart.data.labels = data.map(row => "deployment: " + row.deployment);
  acquisitionsChart.data.datasets[0].data = data.map(row => row.count);
  
  // Update the chart
  acquisitionsChart.update();
};

socket.onclose = function(event) {
  console.log("WebSocket connection closed.");
};

socket.onerror = function(event) {
  console.error("WebSocket error:", event);
};
