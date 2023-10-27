function main() {
    const enterButton = document.getElementById('enter');
    
    const enter = () => {
        window.location.href = '../apps-echo/index.html';
    }

    enterButton.addEventListener('click', enter)
}

document.addEventListener('DOMContentLoaded', main);