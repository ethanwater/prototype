function main() {
    const enterButton = document.getElementById('enter');
    
    const enter = () => {
        window.location.href = '../apps-echo/echo.html';
    }

    enterButton.addEventListener('click', enter)
}

document.addEventListener('DOMContentLoaded', main);