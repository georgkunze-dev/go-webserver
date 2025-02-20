mybutton = document.getElementById("myBtn");

// Wenn 20px nach unten gescrollt ist, wird der bitton sichtbar
window.onscroll = function() {scrollFunction()};

function scrollFunction() {
  if (document.body.scrollTop > 100 || document.documentElement.scrollTop > 100) {
    mybutton.style.display = "block";
  } else {
    mybutton.style.display = "none";
  }
}

// Wenn der button angeklickt wird, wird nach oben gescrollt
function topFunction() {
  document.body.scrollTop = 0; 
  document.documentElement.scrollTop = 0;
}