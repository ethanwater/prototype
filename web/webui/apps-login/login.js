import { strip } from "../apps-echo/echo.js";

async function loginResponse(endpoint, email, password, aborter) {
    const response = await fetch(`/${endpoint}?q=${encodeURIComponent(email)}&p=${encodeURIComponent(password)}`, {
        signal: aborter.signal,
    });
    const text = await response.text();

    if (response.ok) {
        return text;
    } else {
        throw new Error(text);
    }
}

function main() {
    const pass = document.getElementById('pass');
    const email = document.getElementById('email');
    email.value = localStorage.getItem('email');
    const enterButton = document.getElementById('enter');
    const inputs = document.querySelectorAll('input');

    let controller;

    inputs.forEach((input) => {
        input.setAttribute('autocomplete', 'off');
        input.setAttribute('autocorrect', 'off');
        input.setAttribute('autocapitalize', 'off');
        input.setAttribute('spellcheck', false);
    });

    const handleIncorrectAnimation = () => {
        pass.classList.remove("incorrect");
        email.classList.remove("incorrect");
    };

    const login = async () => {
        if (controller !== undefined) {
            controller.abort();
        }

        controller = new AbortController();

        try {
            for (const endpoint of ['login']) {
                const responseText = await loginResponse(endpoint, email.value, pass.value, controller);
                const results = JSON.parse(responseText);

                if (results == null || results === false) {
                    pass.classList.add("incorrect");
                    email.classList.add("incorrect");

                    pass.addEventListener("animationend", handleIncorrectAnimation, { once: true });
                    email.addEventListener("animationend", handleIncorrectAnimation, { once: true });
                } else {
                    localStorage.setItem("email", email.value);
                    window.location.assign("../apps-echo/index.html");
                }
            }
        } catch (error) {
            console.error(error);
        }
    };

    document.addEventListener('keypress', (e) => {
        if (e.key == 'Enter' && pass.validity.valid && email.validity.valid) {
            login();
        }
    });

    email.addEventListener('input', () => {
        if (strip(email.value) !== "" && strip(pass.value) !== "") {
            enterButton.disabled = false;
        } else {
            enterButton.disabled = true;
        }
    });

    enterButton.addEventListener('click', login);
}

document.addEventListener('DOMContentLoaded', main);
