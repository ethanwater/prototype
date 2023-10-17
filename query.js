//query endpoint "/user_query?=user"
async function search(endpoint, query, abort) {
    const response = await fetch(`/${endpoint}?q=${query}`, {signal: abort});
    const text = await response.text();
    if (response.ok) {
        return text;
    } else {
        throw new Error(input);
    }
}

var test = document.getElementById("test");
var input = document.getElementById("pass").value;
var controller = new AbortController();

function GetValue() {
    search("/github_user_query", input.value, controller.signal)
    test.innerHTML = input;
}