<script>

</script>

<!-- <h1>Правила и законы социалистического парламента ГХД</h1> -->
<div id="main" class="f-18">
        <div class="white-section">
          <h1 class="rules-title">Přehled pravidel</h1>     
          <div class="buttons">
            <a href="../home/rules/details"><button class="rules-button">Podrobná pravidla</button></a>
            <a href="../home/rules/instructions"><button class="rules-button">Pokyny k organizaci</button></a>
          </div>
          <div class="rules-list">
            <h2>Kdo může soutěžit?</h2>
            <p class="blue-left">
              V soutěžích mohou soutěžit žáci 8. a 9. tříd základních škol.
            </p>
            <p class="blue-left">
              Každé soutěžní družstvo se skládá z maximálně 5 členů. Členy týmu lze po domluvě přes email vystřídat.
            </p>
          </div>
        </div>
        <div class="gray-section">
          <div class="rules-list">
            <h2>Povolené pomůcky</h2>
            <p class="blue-left">
              V online kole bude každý tým potřebovat minimálně jeden počítač s připojením k internetu, ale doporučujeme jeden počítač pro každého člena týmu.
            </p>
            <p class="blue-left">
              Matematické, fyzikální a chemické tabulky pro základní nebo střední školy
            </p>
            <p class="blue-left">
              Kalkulačky bez připojení k internetu
            </p>
            <p class="blue-left">
              Psací a rýsovací pomůcky, pomocné papíry
            </p>
            <p class="blue-left">
              Vše ostatní je zakázáno.
            </p>
          </div>
        </div>
        <div class="white-section">
          <div class="rules-list">
            <h2>Jak probíhají soutěže?</h2>
            <p class="blue-left">
              Soutěže probíhají ve dvou kolech – online a prezenčním.
            </p>
            <p class="blue-left">
              V každém kole dostane tým na začátku určitý obnos DopplerCoinů. Vyhraje ten, který jich bude mít na konci hry nejvíce.
            </p>
            <div class="blue-left">
              <span>
                Za DopplerCoiny si tým může kupovat úlohy tří různých obtížností:
                <ul>  
                  <li>Lehká: úlohy, pro které stačí "dosadit do vzorečku", odměna za ně je ale velmi nízká.</li>
                  <li>Střední: úlohy, které přibližně odpovídají obtížnosti úloh v testech z matematiky, resp. fyziky.</li>
                  <li>Těžká: úlohy, u kterých je potřeba se opravdu zamyslet. Jsou sice obtížnější, ale zato zajímavější a odměna za jejich vyřešení je větší.</li>
                </ul>
                  Čím je úloha těžší, tím vyšší je i její cena.
              </span>
            </div>
            <p class="blue-left">
              Za vyřešení úlohy dostane tým odměnu. Čím je úloha těžší, tím větší je odměna.
            </p>
            <p class="blue-left">
              Pokud tým úlohu není schopen vyřešit, může jí prodat zpět, ale za méně než za kolik ji koupil.
            </p>
          </div>
        </div>
        <div class="gray-section">
          <div class="rules-list">
            <h2>Jak budou vypadat online kola?</h2> 
            <p class="blue-left">
              Online kola proběhnou protřednictvím herního serveru na adrese, kterou obdrží týmy emailem před začátkem.
            </p>
            <!-- <p class="blue-left" style="color: red; font-weight: bold;">
              Pokud vám do 24 hodin od registrace nepřijde potvrzovací mail, prosíme zkontrolujte složku spam, popř. správnost zadané emailové adresy
            </p> -->
            <!-- <p>Harmonogram soutěží bude následující:
                <p>9:45 - Spuštění úvodního streamu na <a href="https://www.youtube.com/channel/UC8hsjufLjcGySi79GDzqJGg" target="_blank" class="link">YouTube</a></p>
                <p>10:00 - Začátek soutěže</p>
                <p>12:00 - Konec soutěže</p>
                <p>12:30 - Vyhlášení výsledku na streamu</p>
            </p> -->
            <p class="blue-left">
              Online kola proběhnou 2. a 4. prosince 2025.
            </p>
          </div> 
        </div>
        <div class="white-section">
          <div class="rules-list">
            <h2>Jak budou probíhat prezenční kola?</h2>
            <p class="blue-left">
              Prezenční kola proběhnou v prostorách Gymnázia Christiana Dopplera.
            </p>
            <p class="blue-left">
              Harmonogram i pravidla budou velmi podobná jako v online kole.
            </p>
            <p class="blue-left">
              Přesnější instrukce týmy obdrží v emailu po postupu.
            </p>
            <p class="blue-left">
              Prezenční kola proběhnou 11. a 16. prosince 2025, zhruba týden po online kolech.
            </p>
          </div> 
        </div>
      </div>
<style>
.rules-button{
    cursor: pointer;
    font-family: 'Lexend';
    font-size: 20px;
    font-weight: 500;
    padding: 10px;  
    padding-inline: 12px;
    color: white;
    border-radius: 10px;
    background-color: var(--color-blue);
    padding-top: 7px;
    width: 250px;
    margin: auto;
    transition: all 0.3s ease-out;
}

.buttons{
    display: flex;
    justify-content: start;
    gap: 2vw;
    padding-bottom: 20px;
    padding-top: 10px;
}

.rules-button:hover{
    background-color: #3c1fe2;
    box-shadow: 0px 0px 20px 0px #907cff7e;
    transition: all 0s;
}

.rules-list p, .rules-list div{
    margin-left: 50px;
    padding-top: 15px;
    padding-bottom: 15px;
    margin-bottom: 30px;
}

.rules-list h2{
    font-size: 30px;
}

.rules-title{
    margin-left: -5%;
}

@media (max-width: 700px){
    .rules-list p, .rules-list div, .rules-title{
        margin-left: 0;
    }
    .buttons{
        flex-direction: column;
        gap: 30px;
    }
}

    #main{
    font-family: 'Lexend';
}
.f-18{
    font-size: 18px;
}
.gray-section{
    padding-left: 20%;
    background-color: var(--color-light-gray);
    padding-top: 20px;
    padding-bottom: 40px;
    padding-right: 15%;
}
.white-section{
    padding-left: 20%;
    background-color: white;
    padding-top: 20px;
    padding-bottom: 40px;
    padding-right: 15%;
}
    .blue-left {
    border-left: 4px var(--color-blue) solid;
    border-radius: 4px;
    padding-left: 10px;
}
</style>