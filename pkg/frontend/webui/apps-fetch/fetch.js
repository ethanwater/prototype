async function fetchUsers(endpoint, aborter) {
  const response = await fetch(`/${endpoint}`, {signal: aborter});
  const text = await response.text();
  if (response.ok) {
    return text;
  } else {
    throw new Error(text);
  }
}

function emoji_span(emoji) {
    const span = document.createElement('span');
    span.innerText = emoji;
    span.classList.add('emoji');
    span.addEventListener('click', () => {
      if (navigator.clipboard) {
        navigator.clipboard.writeText(emoji);
      }
    });
    return span;
  }


function main() {
    const host = document.getElementById('host');
    host.innerText = window.location.hostname;
    const emojis = document.getElementById('emojis');
    const loader = document.getElementById('loader_box');
  
    let controller; 
    let pending = 0; 
    const displayed = new Set(); 
    const fetch = () => {
      if (controller != undefined) {
        controller.abort();
      }
      controller = new AbortController();

      while (emojis.children.length > 1) {
        emojis.children[0].remove();
      }
      displayed.clear();
  
      for (const endpoint of ['fetch']) {
        if (pending == 0) {
          loader.hidden = false;
        }
        pending++;
  
        fetchUsers(endpoint, controller.signal).then((v) => {
          const results = JSON.parse(v);
          if (results == null || results.length == 0) {
            return;
          }
          for (let emoji of results) {
            if (!displayed.has(emoji.Name)) {
              displayed.add(emoji.Name);
              emojis.insertBefore(emoji_span(emoji.Name), loader);
            }
          }
        }).finally(() => {
          pending--;
          if (pending == 0) {
            loader.hidden = true;
          }
        });
      }
    }

    const addApp = document.getElementById('addApp');
    const echoApp = document.getElementById('echoApp');
  
    const addapp = () => {
      window.location.assign("../apps-adduser/index.html");
    }
    addApp.addEventListener('click', addapp);
  
    const echoapp = () => {
      window.location.assign("../apps-echo/index.html");
    }
    echoApp.addEventListener('click', echoapp);
    fetch();
}
  
document.addEventListener('DOMContentLoaded', main);