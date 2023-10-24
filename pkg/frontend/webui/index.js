function strip(s) {
  return s.replace(/\s+/g, '');
}

async function echoResponse(endpoint, query, aborter) {
  const response = await fetch(`/${endpoint}?q=${query}`, {signal: aborter});
  const text = await response.text();
  if (response.ok) {
    return text;
  } else {
    throw new Error(text);
  }
}

async function add(endpoint, query, aborter) {
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
  const addButton = document.getElementById('add');
  const echoButton = document.getElementById('echo');
  const loader = document.getElementById('test');
  const timer = document.getElementById('time');

  query.focus();
  addButton.disabled = true;
  echoButton.disabled = true;

  let controller; 
  let pending = 0; 
  const add_user = () => {
    if(controller != undefined) {
      controller.abort();
    }
    controller = new AbortController();

    for (const endpoint of ['add']) {
      add(endpoint, query.value, controller.signal).then((x) => {
        JSON.parse(x);
      })
    }
  }

  const echo = () => {
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
      echoResponse(endpoint, query.value, controller.signal).then((v) => {
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

  addButton.addEventListener('click', add_user);
  echoButton.addEventListener('click', echo);

  //query.addEventListener('keypress', (e) => {
  //  if (e.key == 'Enter' && strip(query.value) != "") {
  //    perform_search();
  //  }
  //});

  query.addEventListener('input', (e) => {
    if (strip(query.value) == "") {
      addButton.disabled = true;
      echoButton.disabled = true;
    } else {
      addButton.disabled = false;
      echoButton.disabled = false;
    }
  });
}

document.addEventListener('DOMContentLoaded', main);