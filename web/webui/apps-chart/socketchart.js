var socket = new WebSocket("wss://localhost:2695/wscalls", "protocolTwo");

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
  const data = [
    { deployment: 1, count: event.data },
    { deployment: 2, count: 100 },
  ];

  acquisitionsChart.data.labels = data.map(row => "deployment: " + row.deployment);
  acquisitionsChart.data.datasets[0].data = data.map(row => row.count);
  
  acquisitionsChart.update();
};

socket.onclose = function(event) {
  console.log("WebSocket connection closed.");
};

socket.onerror = function(event) {
  console.error("WebSocket error:", event);
};

window.addEventListener('beforeunload', function(event) {
  if (socket.readyState === WebSocket.OPEN) {
      socket.close();
  }
});

