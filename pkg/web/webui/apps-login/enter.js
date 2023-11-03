async function loginResponse(endpoint, query, aborter) {
    const response = await fetch(`/${endpoint}?${query}`, {signal: aborter});
    const text = await response.text();
    if (response.ok) {
      return text;
    } else {
      throw new Error(text);
    }
  }


function main() {
    const enterButton = document.getElementById('enter');
    const email = document.getElementById('email');
    const pass = document.getElementById('pass');
    const attempts = document.getElementById('atp');
    var loginAttempts = 0
  
    const inputs = document.querySelectorAll('input');

  
    inputs.forEach(input => {
      input.setAttribute('autocomplete', 'off')
      input.setAttribute('autocorrect', 'off')
      input.setAttribute('autocapitalize', 'off')
      input.setAttribute('spellcheck', false)
    })
  
    let controller; 

    const login = () => {
      if (controller != undefined) {
        controller.abort();
      }
      controller = new AbortController();
  
      for (const endpoint of ['login']) {
        loginResponse(endpoint, email.value + "&" + pass.value, controller.signal).then((v) => {
          const results = JSON.parse(v);
          if (results == null || results.length == 0 || results == false ) {
            loginAttempts ++;
            attempts.innerText = "attempts left: " + (3 - loginAttempts);
            if (loginAttempts == 3) {
              enterButton.disabled = true;
            }
          } else {
            window.location.assign("../apps-adduser/index.html");
          }
        }).finally(() => {
          pending--;
          if (pending == 0) {
          }
        });
      }
    }

    email.addEventListener('input', (e) => {
      if (strip(email.value) != ""  & strip(pass.value)!="") {
        enterButton.disabled = false;
      } else {
        enterButton.disabled = true;
      }
    });
    

    enterButton.addEventListener('click', login);
  }
  
  document.addEventListener('DOMContentLoaded', main);