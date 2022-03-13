let content = document.getElementById("canvas");
let wrapper = document.getElementById("wrapper");
let exhaustFlame = document.getElementById("exhaust-flame");
let exhaustFumes1 = document.getElementById("exhaust-fumes-1");
let exhaustFumes2 = document.getElementById("exhaust-fumes-2");
let exhaustFumes3 = document.getElementById("exhaust-fumes-3");
let exhaustFumes4 = document.getElementById("exhaust-fumes-4");
let exhaustFumes5 = document.getElementById("exhaust-fumes-5");
let line1 = document.getElementById("line1");
let line2 = document.getElementById("line2");
let line3 = document.getElementById("line3");
let line4 = document.getElementById("line4");


let btn = document.getElementById("submitButton");
btn.addEventListener("click", function() {
    content.style = 'transform: rotate(90deg);';
    wrapper.style = 'animation: bounce 0.5s infinite;';
    exhaustFlame.style = 'visibility: visible;';
    exhaustFumes1.style = 'visibility: visible;';
    exhaustFumes2.style = 'visibility: visible;';
    exhaustFumes3.style = 'visibility: visible;';
    exhaustFumes4.style = 'visibility: visible;';
    exhaustFumes5.style = 'visibility: visible;';
    line1.style = 'visibility: visible;';
    line2.style = 'visibility: visible;';
    line3.style = 'visibility: visible;';
    line4.style = 'visibility: visible;';
});