<script>
    import logo from "$lib/assets/images/logo.svg";
</script>

<nav class="navbar">
    <div class="row-segment-wave grid-nav">
      <div class="nav1">
        <a href="/home"><img class="logo" alt="" src={logo} width="446px" height="50px"></a>
      </div>
      <a href="../home/about_us" class="nav-element nav2">O nás</a>
      <a href="../home/register" class="nav-element nav3">Registrace</a>
      <a href="../home/archive" class="nav-element nav4">Archiv</a>
      <a href="../home/rules" class="nav-element nav5">Pravidla</a>
    </div>
</nav>

<style>
.navbar{
    background-color: var(--color-blue);
    padding-top: 30px;
    padding-bottom: 30px;
    font-family: "Fredoka";
}
.nav-element{
    text-decoration: none;
    color: white;
    font-size: 25px;
    font-weight: 600;
    margin: auto;
}
.nav-element::before {
    content: "";
    display: inline-block;
    width: 10px;
    height: 10px;
    background-color: white;
    border-radius: 2px;
    margin-inline: 5px;
    transform: translateY(-20px);
    opacity: 0;
    transition: transform 0.1s ease-out, opacity 0.1s ease-out;
}
.nav-element::after {
    content: "";
    display: inline-block;
    width: 10px;
    height: 10px;
    background-color: white;
    border-radius: 2px;
    margin-inline: 5px;
    transform: translateY(20px);
    opacity: 0;
    transition: transform 0.1s ease-out, opacity 0.1s ease-out;
}
.nav-element:hover::after{
    transform: translateY(-10px);
    opacity: 1;
    transition: transform 0.1s ease-out, opacity 0.1s ease-out;
}
.nav-element:hover::before{
    transform: translateY(10px);
    opacity: 1;
    transition: transform 0.1s ease-out, opacity 0.1s ease-out;
}
.grid-nav{
    display: grid;
    grid-template: 'n1 n2 n3 n4 n5';
    text-align: center;
}
.nav1 {grid-area: n1}
.nav2 {grid-area: n2}
.nav3 {grid-area: n3}
.nav4 {grid-area: n4}
.nav5 {grid-area: n5}

@media (max-width: 1100px){
    .grid-nav{
        grid-template:
            'n1 n1 n1 n1'
            'n2 n3 n4 n5'
    }
    .logo{
        padding-bottom: 20px;
        margin-bottom: 20px;
    }
}

@media (max-width: 700px){
    .grid-nav{
        grid-template:
            'n1'
            'n2'
            'n3'
            'n4'
            'n5' ;
    }
    .nav-element{
        padding-top: 10px;
        padding-bottom: 10px;
        border-top: var(--color-lightblue) 2px solid;
    }
    .logo{
        margin-bottom: 20px;
        width: 80%;
        max-width: 446px;
        height: auto;
    }
    .nav2{
        border: none;
    }
}
</style>