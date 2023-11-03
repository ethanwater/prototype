async function fetchUsers(endpoint, aborter) {
  const response = await fetch(`/${endpoint}`, {signal: aborter});
  const text = await response.text();
  if (response.ok) {
    return text;
  } else {
    throw new Error(text);
  }
}

function account_span(account) {
    const span = document.createElement('span');
    span.innerText = account;
    span.classList.add('account');
    span.addEventListener('click', () => {
      if (navigator.clipboard) {
        navigator.clipboard.writeText(account);
      }
    });
    return span;
  }


function main() {
    const host = document.getElementById('host');
    host.innerText = window.location.hostname;
    const accounts = document.getElementById('accounts');
    const loader = document.getElementById('loader_box');
  
    let controller; 
    let pending = 0; 
    const displayed = new Set(); 
    const fetch = () => {
      if (controller != undefined) {
        controller.abort();
      }
      controller = new AbortController();

      while (accounts.children.length > 1) {
        accounts.children[0].remove();
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
          for (let account of results) {
            if (!displayed.has(account.Alias)) {
              displayed.add(account.Alias);
              accounts.insertBefore(account_span(account.Alias), loader);
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