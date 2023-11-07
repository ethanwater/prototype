import {strip} from "../apps-echo/echo"

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
    var email = (<HTMLInputElement>document.getElementById('email'));
    const host = localStorage.getItem('email');
    if (host != null) {
      email.value = host;
    }

    var enterButton = (<HTMLInputElement>document.getElementById('enter'));
    var pass = (<HTMLInputElement>document.getElementById('pass'));
    var attempts = (<HTMLInputElement>document.getElementById('atp'));
    var loginAttempts = 0
  
    const inputs = document.querySelectorAll('input');
		
    inputs.forEach(input => {
      input.setAttribute('autocomplete', 'off')
      input.setAttribute('autocorrect', 'off')
      input.setAttribute('autocapitalize', 'off')
      input.setAttribute('spellcheck', 'off')
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
						localStorage.setItem("email", email.value)
            window.location.assign("../apps-echo/index.html");
          }
        });
      }
    }

		document.addEventListener('keypress', (p) => {
			if (p.key == 'Enter') {
				login();
			}
		});
  
    email.addEventListener('input', (e) => {
      if (strip(email.value) != "" && strip(pass.value)!="") {
        enterButton.disabled = false;
      } else {
        enterButton.disabled = true;
      }
    });

    enterButton.addEventListener('click', login);
  }
  
  document.addEventListener('DOMContentLoaded', main);
