function strip(s) {
    return s.replace(/\s+/g, '');
}

function initCookie2FA(name, value) {
    const d = new Date();
    d.setTime(d.getTime() + (5 * 60 * 1000));
    let expires = "expires=" + d.toUTCString();
    document.cookie = name + "=" + value + ";" + expires + ";path=/";
}

function getCookie2FA(cname) {
    let name = cname + "=";
    let decodedCookie = decodeURIComponent(document.cookie);
    let ca = decodedCookie.split(';');
    for(let i = 0; i <ca.length; i++) {
        let c = ca[i].trim();
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}

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

async function fetchKey(endpoint, aborter) {
    const response = await fetch(`/login/${endpoint}`, {
        signal: aborter
    });
    const text = await response.text();
    if (response.ok) {
        return text;
    } else {
        throw new Error(text);
    }
}

async function verifyCode(endpoint, hash, input, aborter) {
    const response = await fetch(`/login/${endpoint}?hash=${encodeURIComponent(hash)}&input=${encodeURIComponent(input)}`, {
        signal: aborter.signal,
    });
    const text = await response.text();

    if (response.ok) {
        console.log("authentication code success")
        return text;
    } else {
        console.log("authentication code failure")
        throw new Error(text);
    }
}

async function retrieveCode() {
    var inputValues = Array.from(document.querySelectorAll('.symbol')).map(input => input.value);
    var input = inputValues.toString().replaceAll(",", "");
    const handleIncorrectAnimation = () => {
        document.getElementById("verify").classList.remove("incorrect");
    };

    if (strip(input).length < 5) {
        return none;
    }

    var controller = new AbortController();
    for (const endpoint of['verifykey']) {
        const responseText = await verifyCode(endpoint, getCookie2FA("vivian2FA"), input.toUpperCase(), controller);
        const results = JSON.parse(responseText);

        if (!results) {
            document.getElementById("verify").classList.add("incorrect");
            document.getElementById("verify").addEventListener("animationend", handleIncorrectAnimation, {
                once: true
            });
            errorMessage("incorrect code");
            console.log(results);
        } else {
            if (document.getElementById('error')) {
                document.getElementById('error').style.visibility = 'false';
                document.getElementById('error').style.display = 'none';
            }
            console.log(results);
        }
    }
}

function createVerificationElement() {
    if (!document.getElementById("verifyid")) {
        var container = document.getElementById("verify");
        var verificationDiv = document.createElement("div");
        verificationDiv.className = "verification";
        verificationDiv.id = "verifyid";

        for (var i = 1; i <= 5; i++) {
            var input = document.createElement("input");
            input.id = "code" + i;
            input.type = "code";
            input.className = "symbol";
            //input.maxLength = 1;
            input.setAttribute("oninput", "focusNextInput(this, event)");
            input.required = true;

            verificationDiv.appendChild(input);
        }

        container.appendChild(verificationDiv);
        document.getElementById("code1").focus();
    }
}

function showPassword() {
    var x = document.getElementById("pass");
    var showButton = document.getElementById('show_hide');
    if (x.type === "password") {
        x.type = "text";
        showButton.innerText = "hide";
    } else {
        x.type = "password";
        showButton.innerText = "show";
    }
}

function errorMessage(msg) {
    if (!document.getElementById("error")) {
        var div = document.createElement("div");
        div.id = "error"
        div.className = "errormsg";
        div.innerText = msg;

        document.getElementById("main").appendChild(div);
    } else {
        var div = document.getElementById("error");
        div.innerText = msg;
    }
}

function main() {
    const pass = document.getElementById('pass');
    const email = document.getElementById('email');
    const enterButton = document.getElementById('enter');
    const showButton = document.getElementById('show_hide');
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

    const maxAttempts = 10;
    var attempts = 0;
    var ifVerifiedBasic = false;

    //TODO: only perform 2FA if enabled, this is for testing
    const login = async() => {
        if (ifVerifiedBasic == true) {
            return;
        }
        if (attempts >= maxAttempts) {
            return;
        }
        if (controller !== undefined) {
            controller.abort();
        }

        controller = new AbortController();

        try {
            for (const endpoint of['login']) {
                email.disabled = true;
                pass.disabled = true;
                const responseText = await loginResponse(endpoint, email.value, pass.value, controller);
                const results = JSON.parse(responseText);

                if (results == null || results === false) {
                    attempts++;
                    var attemptslog = ": attempts left: ".concat(maxAttempts - attempts);
                    pass.classList.add("incorrect");
                    email.classList.add("incorrect");

                    pass.addEventListener("animationend", handleIncorrectAnimation, {
                        once: true
                    });
                    email.addEventListener("animationend", handleIncorrectAnimation, {
                        once: true
                    });

                    if (attempts < 5) {
                        errorMessage("invalid credentials");
                    } else if (attempts >= 5 && attempts < maxAttempts) {
                        errorMessage("invalid credentials".concat(attemptslog));
                    } else {
                        errorMessage("too many login attempts");
                        pass.value = "";
                        email.value = "";
                        pass.innerText = "";
                        email.innerText = "";
                        email.disabled = true;
                        pass.disabled = true;
                        return;
                    }
                } else {
                    if (document.getElementById('error')) {
                        document.getElementById('error').style.visibility = 'false';
                        document.getElementById('error').style.display = 'none';
                    }
                    createVerificationElement();
                    ifVerifiedBasic = true;
                    for (const endpoint of['generatekey']) {
                        fetchKey(endpoint, controller.signal).then((v) => {
                            const results = JSON.parse(v);
                            if (results == null || results.length == 0) {
                                return;
                            }
                            document.cookie = initCookie2FA("vivian2FA", results) 
                        })
                    }

                    //clears email and pass divs to replace with verification element
                    //email.style.display = 'none';
                    //email.style.visibility = 'hidden';
                    //pass.style.display = 'none';
                    //pass.style.visibility = 'hidden';
                    enterButton.innerText = 'verify';
                }
            }
            email.disabled = false;
            pass.disabled = false;
        } catch (error) {
            errorMessage("server error: cannot validate")
            console.error(error);
        }

        if (ifVerifiedBasic == true) {
            email.disabled = true;
            pass.disabled = true;
        }
    };

    //pass.addEventListener('keydown', (e) => {
    //    if (e.key == 'Backspace') {
    //        pass.innerText = "";
    //        pass.value = "";
    //    }
    //})

    document.addEventListener('keypress', (e) => {
        if (!ifVerifiedBasic) {
            if (e.key == 'Enter' && pass.validity.valid && email.validity.valid) {
                login();
            }
        } else {
            if (e.key == 'Enter') {
                retrieveCode();
            }
        }
    });

    enterButton.addEventListener('click', () => {
        if (!ifVerifiedBasic) {
            login();
        } else {
            retrieveCode();
        }
    });

    showButton.addEventListener('click', showPassword);
}

document.addEventListener('DOMContentLoaded', () => {
    main(); 
});
