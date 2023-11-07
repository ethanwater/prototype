import { strip } from "../apps-echo/echo.js";

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
    const pass = document.getElementById('pass');
    const passquery = document.querySelector("input#pass");
    const email = document.getElementById('email');
    const emailquery = document.querySelector("input#email");
	email.value = localStorage.getItem('email');

    const enterButton = document.getElementById('enter');
    const attempts = document.getElementById('atp');
    const inputs = document.querySelectorAll('input');
		
    var loginAttempts = 0;
    let controller; 

    inputs.forEach(input => {
        input.setAttribute('autocomplete', 'off')
        input.setAttribute('autocorrect', 'off')
        input.setAttribute('autocapitalize', 'off')
        input.setAttribute('spellcheck', false)
    })
  
    const login = () => {
        if (controller != undefined) {
            controller.abort();
        }
        controller = new AbortController();
        const endpoint = 'login'
        loginResponse(endpoint, email.value + "&" + pass.value, controller.signal).then((v) => {
            const results = JSON.parse(v);
            if (results == null || results == false ) {
                loginAttempts ++;
                attempts.innerText = "attempts left: " + (3 - loginAttempts);
                if (loginAttempts == 3) {
                    attempts.innerText = "locked";
                    enterButton.disabled = true;
                }
                passquery.classList.add("incorrect");
                emailquery.classList.add("incorrect");

                passquery.addEventListener("animationend", (e) => {
                    passquery.classList.remove("incorrect");
                    emailquery.classList.remove("incorrect");
                });
            } else {
			    localStorage.setItem("email", email.value)
                window.location.assign("../apps-echo/index.html");
            }
        });
    }

    document.addEventListener('keypress', (e) =>  {
        if (loginAttempts < 3) {
            if (e.key == 'Enter') {
                login();
            }
        }
    })

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