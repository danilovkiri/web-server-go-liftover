var coll = document.getElementsByClassName("collapsibleButton");
var i;

for (i = 0; i < coll.length; i++) {
    coll[i].addEventListener("mouseover", function() {
        this.classList.toggle("activeButton");
        var content = this.nextElementSibling;
        if (content.style.maxHeight){
            content.style.maxHeight = null;
        } else {
            content.style.maxHeight = content.scrollHeight + "px";
        }
    });
    coll[i].addEventListener("mouseout", function() {
        this.classList.toggle("activeButton");
        var content = this.nextElementSibling;
        if (content.style.maxHeight){
            content.style.maxHeight = null;
        } else {
            content.style.maxHeight = content.scrollHeight + "px";
        }
    });
}