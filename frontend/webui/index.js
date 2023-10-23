function strip(s) {
  return s.replace(/\s+/g, '');
}

async function search(endpoint, query, aborter) {
  const response = await fetch(`/${endpoint}?q=${query}`, {signal: aborter});
  const text = await response.text();
  if (response.ok) {
    return text;
  } else {
    throw new Error(text);
  }
}

function main() {
  const query = document.getElementById('pass');
  const button = document.getElementById('but');
  const loader = document.getElementById('test');
  const timer = document.getElementById('time');
  query.focus();
  button.hidden = true;

  let controller; 
  let pending = 0; 
  const displayed = new Set(); 
  const perform_search = () => {
    if (controller != undefined) {
      controller.abort();
    }
    controller = new AbortController();

    for (const endpoint of ['echo']) {
      if (pending == 0) {
        loader.hidden = false;
      }
      pending++;

      var start_time = performance.now();
      search(endpoint, query.value, controller.signal).then((v) => {
        const results = JSON.parse(v);
        if (results == null || results.length == 0) {
          loader.innerHTML = "red";
        } else {
          loader.innerHTML = v;
        }
      }).finally(() => {
        pending--;
        if (pending == 0) {
          //loader.hidden = true;
        }
      });
      var end_time = performance.now();
      timer.innerHTML = (end_time - start_time);
    }
  }

  button.addEventListener('click', perform_search);

  query.addEventListener('keypress', (e) => {
    if (e.key == 'Enter' && strip(query.value) != "") {
      perform_search();
    }
  });

  query.addEventListener('input', (e) => {
    if (strip(query.value) == "") {
      button.disabled = true;
    } else {
      button.disabled = false;
    }
  });
}

document.addEventListener('DOMContentLoaded', main);