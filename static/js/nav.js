function togglenav() {
    var x = document.getElementById("hamburger");
    if (x.className === "pg-nav") {
      x.className += " responsive";
    } else {
      x.className = "pg-nav";
    }
  }
